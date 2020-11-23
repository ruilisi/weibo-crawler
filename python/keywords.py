from gensim.models import KeyedVectors
from gensim import models
from sklearn.cluster import DBSCAN 

import numpy as np
import pandas as pd
import matplotlib as mpl
import matplotlib.pyplot as plt
import hdbscan


#  corpus.txt is a massive dataset from Tencent AI Lab Embedding Corpus for Chinese Words and Phrases
#  link refer to https://ai.tencent.com/ailab/nlp/zh/embedding.html 
#  first time run the following two lines
#  wv = KeyedVectors.load_word2vec_format('./corpus/corpus.txt', binary=False)
#  wv.save_word2vec_format('./corpus/corpus.bin',binary=True)

wv = KeyedVectors.load_word2vec_format('./corpus/corpus.bin', binary=True)

with open('keywords.txt', 'r') as f:
    raw_dataset = f.readlines()

dataset = []
valueset = []
for i in raw_dataset:
    i = i.strip()
    if i in wv:
        dataset.append(wv[i])
        valueset.append(i) 
    else:
       raw_dataset.remove(i+"\n")

#print(valueset,len(raw_dataset),len(dataset))


#  res = []
#  for eps in np.arange(0.001,1,0.05):
    #  for min_samples in range(2,10):
        #  clustering = DBSCAN(eps = eps, min_samples = min_samples).fit(dataset)
        #  n_clustering = len([i for i in set(clustering.labels_) if i == -1])
        #  outliners = np.sum(np.where(clustering.labels_ == -1, 1,0))
        #  res.append({'eps':eps,'min_samples':min_samples,'n_clusters':n_clustering,'outliners':outliners})
#
#  df = pd.DataFrame(res)
#  print(df.loc[df.n_clusters == 3, :])
#
#  clustering = DBSCAN(eps=0.05, min_samples=2).fit(dataset)

clustering = hdbscan.HDBSCAN(min_cluster_size=2)
cluster_labels = clustering.fit_predict(dataset)

dis = []
for data in dataset:
    d = np.linalg.norm(data)
    dis.append(d)

# print(dis)
# print(clustering.labels_)
print(cluster_labels)

setLabels = sorted(list(set(cluster_labels)))
dic = {}
for i in setLabels:
    dic[i] = []

rs = []
for i in range(0,len(dis)):
    rs.append([valueset[i],dis[i],cluster_labels[i]])
print(rs)

for i in rs:
    dic[i[2]].append(i[0]) 
print(dic)

fig = plt.figure()
s=[100, 103, 150, 200]

with open("dict.txt","w") as fw:
    for key,value in dic.items():
       fw.writelines('{key}:{value}\n'.format(key = key, value = value))
       fw.flush()

ax1 = fig.add_subplot(1, 2, 1)
ax1.scatter(range(1,len(dis)+1),dis,marker='o')
ax1.set_title('dataset')
#
ax2 = fig.add_subplot(1, 2, 2)
ax2.scatter(range(1,len(dis)+1), dis,c=cluster_labels,marker='o')
ax2.set_title('result')

plt.show()

