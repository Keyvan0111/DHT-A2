# conftest.py
import pytest
import json
import os
from test_chord_linear import discover_cluster

def pytest_addoption(parser):
    parser.addoption(
        "--urls",
        action="store",
        default=None,
        help='JSON array of node addresses WITHOUT prefixes, e.g. '
             '\'["c11-0:41305","c11-6:35399"]\'',
    )
    parser.addoption(
        "--domain",
        action="store",
        default=os.environ.get("DOMAIN_SUFFIX", "ifi.uit.no"),
        help='Domain suffix to append when missing (default: "ifi.uit.no"). '
             'Env override: DOMAIN_SUFFIX',
    )

@pytest.fixture(scope="session")
def cluster(pytestconfig):
    urls_json = pytestconfig.getoption("urls")
    if not urls_json:
        # also allow env var for convenience
        urls_json = os.environ.get("CLUSTER_URLS")
    if not urls_json:
        raise RuntimeError(
            'Pass --urls \'["host:port", ...]\' or set CLUSTER_URLS env var.'
        )

    domain = pytestconfig.getoption("domain")
    raw_urls = json.loads(urls_json)
    nodes, m_bits = discover_cluster(raw_urls, domain)
    assert len(nodes) >= 2, "Need at least 2 nodes to test forwarding"
    return nodes, m_bits
