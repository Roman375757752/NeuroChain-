import pandas as pd
import networkx as nx
import matplotlib.pyplot as plt

# Загрузка транзакций
df = pd.read_csv("transactions.csv")

# Создание направленного графа
G = nx.DiGraph()
for _, row in df.iterrows():
    G.add_edge(row['from'], row['to'], weight=row['amount'])

# Ограничим визуализацию первыми 300 узлами, чтобы не было перегрузки
nodes_to_draw = list(G.nodes)[:300]
G_sub = G.subgraph(nodes_to_draw)

# Построение графа
plt.figure(figsize=(18, 12))
pos = nx.spring_layout(G_sub, k=0.15, iterations=20)
nx.draw_networkx_nodes(G_sub, pos, node_size=40, node_color="skyblue")
nx.draw_networkx_edges(G_sub, pos, edge_color="gray", arrows=True, width=0.5, arrowsize=5)
plt.title("Neurochain Transactions Graph (300 узлов)")
plt.axis('off')
plt.tight_layout()
plt.savefig("transactions.png", dpi=300)
plt.show()
