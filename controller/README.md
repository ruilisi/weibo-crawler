# controller

API list

## GET /ping 

### Response
<div>
<pre>
  success:
  200 OK
  return:
  {
    pong
	}
</pre>
</div>

---
 
## GET /task
Auto run every 30 minutes by cronjob.You can call the interface manually to complete a task once. 
This task is about to request pre-stored blogger data, get everyone's homepage blog posts, and only crawl new posts in page each time.  

### Response 
<div>
<pre>
  success:
  200 OK
  return:
  {
    "status": "succes"
	}
 </pre>
</div>

---

## POST /query_blogs

**request**  
Query blogs by parameters such as interval, keywords and category.
Due to category is the result of clustering keywords, they are mutually exclusive. 
Interval and tags are intersection, topic, super_topic and keywords are union.
Amount is the demanding amount of the result.
Last time is the record of last response result oldest blog.

**response**  
Tags contains super topic, topic and the preset keywords.
Content is the main body of the post.
Pictuers are url array list of the pictures.
Video and article are the short url of the resource. 
Rt is refered to retweet, the blog id of the origin post.

### Request
<div>
<pre>
Query Parameters
{
    "start_time":"2020-11-13 12:37:01",
    "end_time":"2020-11-15 18:37:01",
    "last_time":"2020-11-15 18:37:01",
    "topic":[],
    "super_topic":[],
    "keywords":[],
    "category":"",
    "amount":"2"
}
</pre>
</div>
### Response
<div>
<pre>
  success:
  200 OK
  return:
  {
    "blogs": [
        {
            "BlogId": "Ju2uG3zW8",
            "BloggerId": "3225714712",
            "Username": "玩Switch的呆呆兽",
            "Follows": 1606115,
            "Avatar": "https://tvax3.sinaimg.cn/crop.0.0.776.776.180/c0448018ly8ghw0wkv7zpj20lk0lkgpf.jpg?KID=imgbed,tva&Expires=1605670255&ssig=J0qetsL6p0",
            "Tags": [
                "精灵宝可梦超话",
                "宝可梦剑盾",
                "宝可梦配信",
                "领取"
            ],
            "Content": "不买电影票也能领到萨戮德配信！（当然不是朋友送的那种）配信领取方法一览：°不买电影票也能直接领到萨戮德配信！配信领取...#宝可梦剑盾# #宝可梦配信# ​​​​",
            "Timestamp": "2020-11-15T18:34:40+08:00",
            "Pictures": null,
            "Video": "",
            "Article": "http://t.cn/A6GCoXxc",
            "Rt": ""
        },
        {
            "BlogId": "Ju2tlAN3a",
            "BloggerId": "1367610417",
            "Username": "电玩口袋社区",
            "Follows": 951058,
            "Avatar": "https://tvax1.sinaimg.cn/crop.0.0.1002.1002.180/51841431ly8gdmfad6j4vj20ru0ruq43.jpg?KID=imgbed,tva&Expires=1605670120&ssig=%2Bn6W5fXf9l",
            "Tags": [
                "发售"
            ],
            "Content": "@显卡吧的那些事 在？出来官宣4090Ti的发售日期 ​​​​",
            "Timestamp": "2020-11-15T18:31:23+08:00",
            "Pictures": [
                "//wx2.sinaimg.cn/orj360/51841431ly1gkq1d17uvij20tw1aztn2.jpg"
            ],
            "Video": "",
            "Article": "",
            "Rt": ""
        },
    ],
    "status": "success"
}
</pre>
</div>

---


## POST /add_bloggers

Add blogger to task list 

### Request
<div>
<pre>
Query Parameters
{
  "bloggers":["2396370381","1206323513"]
}
</pre>
</div>
### Response
<div>
<pre>
  success:
  200 OK
  return:
  {
    "status":"success"
  }
</pre>
</div>

---


## GET /category/set

Get clustering result from request and save to database.
number ahead means the category, -1 means noise. 
"/dlc" is the manually labeled name.


### Request
<div>
<pre>
Query Parameters
{
    "category": "-1:['a', 'b', 'c, 'd']/未分类\n0:['dlc', '追加内容', '扩展内容']/dlc"
}
</pre>
</div>
### Response
<div>
<pre>
  success:
  200 OK
  return:
  {
    "status":"success"
  }
</pre>
</div>

---

## GET /category/query

Get category from database and the hit of last 24 hours sort by hit.

### Response
<div>
<pre>
  success:
  200 OK
  return:
  {
    "categories": [
        {
            "Id": 4,
            "Name": "明星娱乐",
            "Count": 157
        },
        {
            "Id": 11,
            "Name": "时事新闻",
            "Count": 106
        },
        {
            "Id": 13,
            "Name": "体育电竞",
            "Count": 83
        },
        ...
    ],
    "status": "success"
}
</pre>
</div>

---

## GET /category/set_name

rename category names by index.

### Request 

<div>
<pre>
Query Parameters  
{
  "name":["明星娱乐","时事新闻"...]
}
</pre>
</div>

### Response
<div>
<pre>
  success:
  200 OK
  return:
  {
    "status": "success"
  }
</pre>
</div>

---

## GET /tags/query

Query all tags including super topic, topic and keywords from database.  
**response**  
Tagtype 0 represents super topic like "IG超话", 1 represents hashtag like "#英雄联盟#", 2 represents keywords we preset.
Category id is the classfication of the tags result.

### Response
<div>
<pre>
  success:
  200 OK
  return:
  {
    "status": "success",
    "tags": [
        {
            "ID": "08013ca9-fefa-4205-bf46-0d0c09809fe8",
            "CreatedAt": "2020-11-17T13:09:44.78047+08:00",
            "UpdatedAt": "2020-11-17T13:09:44.78047+08:00",
            "Name": "英雄联盟奥德赛皮肤",
            "TagType": 1,
            "CategoryId": 0
        },
        {
            "ID": "cea66587-6bdd-46c7-bbe0-768b4f40501a",
            "CreatedAt": "2020-11-17T13:10:02.99029+08:00",
            "UpdatedAt": "2020-11-17T13:10:02.99029+08:00",
            "Name": "如何看待开发者反抗苹果税",
            "TagType": 1,
            "CategoryId": 0
        }
        ...
    ]
  }
</pre>
</div>

---

##  POST /tags/cache_keywords
  
Save or update redis keywords from database

### Response 
<div>
<pre>
  success:
  200 OK
  return:
  {
    "status": "success"
  }
</pre>
</div>

---

## POST /tags/keywords

Query keywords from database and write to /python/keywords.txt for clustering

### Response 
<div>
<pre>
  success:
  200 OK
  return:
  {
    "status": "success"
    "tags": [
        {
            "ID": "08013ca9-fefa-4205-bf46-0d0c09809fe8",
            "CreatedAt": "2020-11-17T13:09:44.78047+08:00",
            "UpdatedAt": "2020-11-17T13:09:44.78047+08:00",
            "Name": "安卓",
            "TagType": 2,
            "CategoryId": 0 
        },
        {
            "ID": "cea66587-6bdd-46c7-bbe0-768b4f40501a",
            "CreatedAt": "2020-11-17T13:10:02.99029+08:00",
            "UpdatedAt": "2020-11-17T13:10:02.99029+08:00",
            "Name": "苹果",
            "TagType": 2,
            "CategoryId": 0 
        }
        ...
    ]
  }
</pre>
</div>

---

## POST /tags/set_keywords

Add preset keywords as features to describe the post
Chinese and English words are recommended

### Request 

<div>
<pre>
Query Parameters  
{
  "keywords":["娱乐","新闻"...]
}
</pre>
</div>

### Response
<div>
<pre>
  success:
  200 OK
  return:
  {
    "status": "success"
  }
</pre>
</div>

---

