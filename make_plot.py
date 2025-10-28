#!/usr/bin/env python3
# Read the CSV produced by run_bench.sh and generate vector plots with error bars.
#
# Usage:
#   python3 make_plot.py results.csv
#
# Output:
#   throughput_vs_nodes.pdf
#   throughput_vs_nodes.svg
#
# The plot shows mean throughput with std-dev error bars for PUT and GET
# across N in {1,2,4,8,16}.

import sys
import csv
from collections import defaultdict, OrderedDict

import matplotlib.pyplot as plt

if len(sys.argv) < 2:
    print("Usage: python3 make_plot.py results.csv", file=sys.stderr)
    sys.exit(2)

csv_path = sys.argv[1]

# Data[N][PHASE] -> list of throughputs
data = defaultdict(lambda: defaultdict(list))

with open(csv_path, newline='') as f:
    reader = csv.DictReader(f)
    for row in reader:
        N = int(row["N"])
        phase = row["PHASE"]
        thr = float(row["throughput_ops_per_sec"])
        data[N][phase].append(thr)

# Stable order of Ns
Ns = sorted(data.keys())

def mean_std(vals):
    if not vals:
        return 0.0, 0.0
    m = sum(vals) / len(vals)
    var = sum((x - m)**2 for x in vals) / len(vals)
    return m, var**0.5

# Prepare series
put_means, put_stds, get_means, get_stds = [], [], [], []

for N in Ns:
    m_put, s_put = mean_std(data[N].get("PUT", []))
    m_get, s_get = mean_std(data[N].get("GET", []))
    put_means.append(m_put)
    put_stds.append(s_put)
    get_means.append(m_get)
    get_stds.append(s_get)

# Plot (matplotlib, single chart, no custom colors/styles)
plt.figure()
plt.errorbar(Ns, put_means, yerr=put_stds, fmt='-o', label='PUT throughput')
plt.errorbar(Ns, get_means, yerr=get_stds, fmt='-s', label='GET throughput')
plt.xlabel('Number of nodes in network')
plt.ylabel('Throughput (ops/sec)')
plt.title('DHT Throughput vs. Cluster Size (mean ± std)')
plt.xticks(Ns)
plt.grid(True, which='both', axis='both', linestyle='--', linewidth=0.5)
plt.legend()
plt.tight_layout()
plt.savefig('/mnt/data/throughput_vs_nodes.pdf')
plt.savefig('/mnt/data/throughput_vs_nodes.svg')

print("Saved: /mnt/data/throughput_vs_nodes.pdf and .svg")
