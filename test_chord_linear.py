# test_chord_linear.py
import requests
import hashlib
import itertools
import random
import string
import time

SESSION = requests.Session()
TIMEOUT = 5
RETRY = 3

def _req(method, url, **kwargs):
    last_exc = None
    for _ in range(RETRY):
        try:
            return SESSION.request(method, url, timeout=TIMEOUT, **kwargs)
        except requests.RequestException as e:
            last_exc = e
            time.sleep(0.2)
    raise last_exc

def sha256_mod_2m(s: str, m_bits: int) -> int:
    # mirror utils.ConsistentHash: SHA-256 -> big.Int -> mod 2^m
    h = hashlib.sha256(s.encode()).digest()
    num = int.from_bytes(h, "big")
    return num % (1 << m_bits)

def between(id_, start, end):
    # mirror utils.between: (start, end] modulo ring
    if start < end:
        return start < id_ <= end
    return id_ > start or id_ <= end

def normalize_url(entry: str, default_domain: str) -> str:
    """
    Accepts:
      - 'c11-0:41305'  -> http://c11-0.ifi.uit.no:41305
      - 'c11-0.ifi.uit.no:41305' -> http://c11-0.ifi.uit.no:41305
      - 'http://c11-0.ifi.uit.no:41305' -> unchanged
    """
    u = entry.strip().rstrip("/")
    if u.startswith("http://") or u.startswith("https://"):
        return u
    # host[:domain?]:port
    if ":" not in u:
        raise ValueError(f"Expected 'host:port', got '{entry}'")
    host, port = u.split(":", 1)
    if "." not in host:
        host = f"{host}.{default_domain.lstrip('.')}"
    return f"http://{host}:{port}"

def load_network(url):
    r = _req("GET", f"{url}/network")
    r.raise_for_status()
    js = r.json()
    return {
        "self": js["self"],              # {"addr","id"}
        "pred": js["predecessor"],       # {"addr","id"}
        "succ": js["successor"],         # {"addr","id"}
        "m": js["hashlen"],              # HASHLEN (m)
    }

def discover_cluster(raw_urls, default_domain: str):
    """
    raw_urls: list of 'host:port' or full URLs.
    Returns (nodes_sorted_by_id, m_bits)
    """
    nodes = {}
    m_val = None
    for raw in raw_urls:
        url = normalize_url(raw, default_domain)
        info = load_network(url)
        nodes[info["self"]["id"]] = {
            "addr": info["self"]["addr"],
            "id": info["self"]["id"],
            "pred": info["pred"]["id"],
            "succ": info["succ"]["id"],
        }
        m_val = info["m"]
    ordered = [nodes[i] for i in sorted(nodes.keys())]
    return ordered, m_val

def key_for_node(target_node, pred_id, m_bits, prefix="k"):
    # brute force small suffixes until hash falls within (pred, self]
    for i in itertools.count():
        k = f"{prefix}-{i}"
        kid = sha256_mod_2m(k, m_bits)
        if between(kid, pred_id, target_node):
            return k, kid

def any_other_node_addr(nodes, node_id):
    for n in nodes:
        if n["id"] != node_id:
            return n["addr"]
    return nodes[0]["addr"]

def put(addr, key, value: str):
    return _req(
        "PUT",
        f"{addr}/storage/{key}",
        data=value.encode("utf-8"),
        headers={"Content-Type":"text/plain"},
    )

def get(addr, key):
    return _req("GET", f"{addr}/storage/{key}")

def random_value(n=12):
    letters = string.ascii_letters
    return "".join(random.choice(letters) for _ in range(n))


# ---------------- Tests (use 'cluster' fixture from conftest.py) ----------------

def test_put_get_on_responsible_node_success(cluster):
    nodes, m = cluster
    target = nodes[0]
    key, _ = key_for_node(target["id"], target["pred"], m, prefix="resp-ok")
    val = random_value()

    r_put = put(target["addr"], key, val)
    assert r_put.status_code == 200, f"PUT expected 200, got {r_put.status_code}"

    other_addr = any_other_node_addr(nodes, target["id"])
    r_get = get(other_addr, key)
    assert r_get.status_code == 200, f"GET expected 200, got {r_get.status_code}"
    assert r_get.text == val, "Value mismatch after PUT/GET"

def test_put_on_non_responsible_node_forwards_and_gets(cluster):
    nodes, m = cluster
    target = nodes[-1]
    entry = nodes[0] if nodes[0]["id"] != target["id"] else nodes[1]

    key, _ = key_for_node(target["id"], target["pred"], m, prefix="fwd-ok")
    val = random_value()

    r_put = put(entry["addr"], key, val)
    assert r_put.status_code == 200, f"PUT via non-responsible expected 200, got {r_put.status_code}"

    # GET from all nodes should work
    for n in nodes:
        r_get = get(n["addr"], key)
        assert r_get.status_code == 200, f"GET from {n['addr']} expected 200, got {r_get.status_code}"
        assert r_get.text == val, "Value mismatch after forwarded PUT"

def test_get_not_found_responsible_node(cluster):
    nodes, m = cluster
    target = nodes[0]
    key, _ = key_for_node(target["id"], target["pred"], m, prefix="missing-resp")
    r_get = get(target["addr"], key)
    assert r_get.status_code == 404, f"Expected 404, got {r_get.status_code}"

def test_get_not_found_via_forwarding(cluster):
    nodes, m = cluster
    target = nodes[1 % len(nodes)]
    entry = nodes[0]
    key, _ = key_for_node(target["id"], target["pred"], m, prefix="missing-fwd")
    r_get = get(entry["addr"], key)
    assert r_get.status_code == 404, f"Expected 404 via forwarding, got {r_get.status_code}"
