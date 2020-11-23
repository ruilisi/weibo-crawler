## python 

### save_cookies.py
 Using the [selenium](https://selenium-python.readthedocs.io/) package, simulate a user to visit a weibo page, enter an account password, and log in.
 Manually verify the user's identity via Weibo.
 Get user's cookies and store them for later use.

### keywords.py
 Using corpus.txt to get words vectors by using [gensim](https://radimrehurek.com/gensim/). 
 Open keywords.txt to get all keywords of the crawler exported by database. 
 We will get all keywords vectors, which contains 200-dimension.
 Using HDBSCAN to cluster word vectors

#### DBSCAN
  Density-based spatial clustering of applications with noise ([DBSCAN](https://en.wikipedia.org/wiki/DBSCAN)) is a data clustering algorithm proposed by Martin Ester, Hans-Peter Kriegel, Jörg Sander and Xiaowei Xu in 1996.
It is a density-based clustering non-parametric algorithm: given a set of points in some space, it groups together points that are closely packed together (points with many nearby neighbors), marking as outliers points that lie alone in low-density regions (whose nearest neighbors are too far away).
DBSCAN is one of the most common clustering algorithms and also most cited in scientific literature.

![pic1  DBSCAN](https://pic4.zhimg.com/v2-e65f8b153a4719f4c5fc5d1246877527_r.jpg)

#### HDBSCAN
  [HDBSCAN](https://hdbscan.readthedocs.io/en/latest/basic_hdbscan.html) is a clustering algorithm developed by Campello, Moulavi, and Sander. It extends DBSCAN by converting it into a hierarchical clustering algorithm, and then using a technique to extract a flat clustering based in the stability of clusters.
The goal of this notebook is to give you an overview of how the algorithm works and the motivations behind it. 

![pic2 HDBSCAN](https://hdbscan.readthedocs.io/en/latest/_images/how_hdbscan_works_10_1.png)

```python
clustering = hdbscan.HDBSCAN(min_cluster_size=2)
cluster_labels = clustering.fit_predict(dataset)
```

[![DeXzOU.png](https://s3.ax1x.com/2020/11/18/DeXzOU.png)](https://imgchr.com/i/DeXzOU)

