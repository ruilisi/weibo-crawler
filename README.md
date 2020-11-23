# go_crawler

Go_crawler is a crawler project based on golang colly framework to crawl weibo sites and get information. It crawls web content by regular expressions and Xpath selector, spatially transforms keywords using word vector model, and clusters text content by HDBSCAN clustering algorithm.

## Features


**Comprehensive capture of user information  
Multi-dimensional collection of weibo content  
Timed incremental acquisition  
Keyword Cluster Analysis  
Category hotspot sorting**  


Go_crawler is based on following tools

|name|description|
|------|--------|
|[Go](https://github.com/golang/go)|An open source programming language that makes it easy to build simple, reliable, and efficient software.|
|[Python](https://www.python.org/)|Python is an interpreted, high-level and general-purpose programming language.|
|[Gin](https://github.com/gin-gonic/gin)|Web struct based on Go, flexible middleware，strong data binding and outstanding performance.|
|[Ginkgo](https://github.com/onsi/ginkgo)|Ginkgo builds on Go's testing package, allowing expressive Behavior-Driven Development ("BDD") style tests.|
|[Colly](https://github.com/gocolly/colly)|Lightning Fast and Elegant Scraping Framework for Gophers|
|[Postgres](https://www.postgresql.org/)|The world's most advanced open source relational database|
|[Gorm](https://github.com/go-gorm/gorm)|The fantastic ORM library for Golang aims to be developer friendly.|
|[Redis](https://redis.io/)|An open source (BSD licensed), in-memory data structure store, used as a database, cache and message broker.|
|[Docker](https://www.docker.com/)|Docker is a tool designed to make it easier to create, deploy, and run applications by using containers.|
|[Sklearn](https://scikit-learn.org/stable/)|Simple and efficient tools for predictive data analysis|
|[Gensim](https://radimrehurek.com/gensim/)|The fastest library for training of vector embeddings – Python or otherwise.|
|[HDBSCAN](https://hdbscan.readthedocs.io/en/latest/how_hdbscan_works.html)|HDBSCAN is a clustering algorithm extends DBSCAN by converting it into a hierarchical clustering algorithm, and then using a technique to extract a flat clustering based in the stability of clusters.|

## Why not python

Python has more than one set of mature crawler frameworks such as scrapy, pyspider and so on.They have excellent runtime mechanism and powerful capabilities.
But when the anti-crawler mechanism is strong, rewriting the middleware is a very difficult task. And it's not flexible enough to be accessed by a project system .

## Struct
```
.
├── application.yml  
├── args
│   ├── args.go
│   └── cmd.go
├── conf  
│   ├── conf_debug.go
│   ├── conf.go
│   └── conf_release.go
├── controller
│   ├── application.go
│   ├── blogger.go
│   ├── category.go
│   ├── error.go
│   ├── query.go
│   ├── tag.go
│   └── task.go
├── corpus  
│   └── corpus.txt
├── db  
│   └── db.go
├── go.mod
├── go.sum
├── jwt  
│   └── jwt.go
├── main.go
├── Makefile  
├── models  
│   ├── base_model.go
│   ├── blog.go
│   ├── blogger.go
│   ├── category.go
│   ├── tag.go
│   └── user.go
├── python  
│   ├── dict.txt
│   ├── keywords.txt
│   ├── keywords_demo.py
│   └── save_cookies.go
├── README.md
├── redis
│   └── redis.go
├── routers  
│   └── router.go
├── tasks
│   ├── tags.go
│   └── tasks.go
├── test
└── util
    ├── agent.go
    ├── cookie.go
    ├── cookies.txt
    └── util.go
```

## How to use

1. install tools and dependency mentioned above
2. config application.yml, establish connection
3. `go run main.go -db create `
4. `go run main.go -db migrate`
5. `go run main.go` 
6. Add bloggers & keywords  post `/add_bloggers`, `/tags/set_keywords` and  `/tags/cache_keywords`
7. wait 30 minutes or call `/task`(local debug environment)
8. let the bullets fly 
9. Post `/query_blogs` to show datas

if you want to do cluster, post `/tags/keywords`, download corpus, `python keywords.py`, adjust and post `/category/set`
 Get `/category/query` to show hot topics

## Api list

[details ](https://github.com/ruilisi/cs-server/tree/master/controller/README.md)

---

|  API   | CALL  | ROUTER  |  FUNCTION  |
|  :----  | :----   | :----  | :----  |
| [Ping](https://github.com/ruilisi/go-crawler/tree/master/controller/application.go) | GET | /ping | ping |
| [Task](https://github.com/ruilisi/go-crawler/tree/master/controller/task.go) | GET | /task | (auto run every 30 minutes) crawler task |
| [Query_blogs](https://github.com/ruilisi/go-crawler/tree/master/controller/query.go) | POST | /query_blogs | query according to different parameters |
| [Add_bloggers](https://github.com/ruilisi/go-crawler/tree/master/controller/blogger.go) | POST | /add_bloggers | add bloggers in task list |
| [Set_category](https://github.com/ruilisi/go-crawler/tree/master/controller/category.go) | GET | /category/set | set category by clustering result |
| [Set_category_name](https://github.com/ruilisi/go-crawler/tree/master/controller/category.go) | POST | /category/set_name | rename category |
| [Query_category](https://github.com/ruilisi/go-crawler/tree/master/controller/category.go) | GET | /category/query | query category |
| [Query_tags](https://github.com/ruilisi/go-crawler/tree/master/controller/query.go) | GET | /tags/query | query tags |
| [Cache_keywords](https://github.com/ruilisi/go-crawler/tree/master/controller/tag.go) | POST | /tags/cache_keywords | save keywords to redis|
| [Get_keywords](https://github.com/ruilisi/go-crawler/tree/master/controller/query.go) | POST | /tags/keywords | query keywords and write to txt for clustering|
| [Set_keywords](https://github.com/ruilisi/go-crawler/tree/master/controller/tag.go) | POST | /tags/set_keywords | add keywords as tags |

---




