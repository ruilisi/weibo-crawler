## tasks 

### tags.go 
fetch different types of tags including topic, supertopic and customized keywords

### tasks.go

Functions are based on go [colly framework](https://github.com/gocolly/colly).Here's the simple example from official document.
```go 
func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}
```
See [examples](https://github.com/gocolly/colly/tree/master/_examples) for more detailed examples.


#### consts
```go 
const (
	baseURL = "https://weibo.com/"

	followTableReg = `<table class=\\\"tb_counter\\\"(.*?)<\\\/table>`
	avatarReg      = `<p class=\\\"photo_wrap\\\">(.*?)<\\\/p>`
	nameReg        = `<h1 class=\\\"username\\\">(.*?)<\\\/h1>`

	weiboDetailReg = `<script>FM\.view\(\{"ns":"pl\.content\.weiboDetail\.index".*?</script>`
	textHTMLReg    = `<div class=\\\"WB_text\sW_f14.*?<\\\/div>`
	textNode       = `//div[@class='WB_text W_f14']//text()`
	spReg          = `<div class=\\\"WB_tag_rec \\\".*?<\\\/div>`
	spNode         = `//a[@class="W_autocut"]/text()`

	timeNodeReg = `<div class=\\\"WB_from\sS_txt2.*?<\\\/div>`
	timeReg     = `date=\\".*?\\"`
	imgNodeReg  = `<img src=\\\"\\\/\\\/wx\d.sinaimg.cn.*?>`
	imgNode     = `//img`
	videoNode   = `//i[@class="W_ficon ficon_cd_video"]`
	articleNode = `//i[@class="W_ficon ficon_cd_longwb"]`
	rtNodeReg   = `<a class=\\\"S_txt2\\\" target=\\\"_blank.*?>`
	rtNode      = `//a`

	followReg          = `<script>FM\.view\(\{"ns":"pl\.relation\.myFollow\.index".*?</script>`
	followHTMLReg      = `<a target=\\"_blank\\" action-type=.*?<\\/a>`
	followMemberNode   = `//a[@node-type="screen_name"]`
	followIDReg        = `id=\d+`
	followNextPageReg  = `<a bpfilter=\\"page\\" class=\\"page next S_txt1 S_line1\\".*?>`
	followNextPageNode = `//a[@class="page next S_txt1 S_line1"]`
)
```
Consts are regular expressions and xpath selectors to fetch details we need in website.
Due to page rendering method of Weibo website, we can't get html details by crawler unless the scripts loading completed.
What we can get is only scripts.
The const weiboDetailReg is the regular expression that can match the script which about to load, that contains all useful information.
The algorithm to fetch information is as follows. Using regular expressions to find matching script or html data.Then use xpath selectors to get attribute or content of the html node.Follows, text, timestamp, image, video and article are the target we want. 

#### fetchFollows
Fetch all follows of the specified weibo user.
This function is used in seed.go, using the cookies that we have saved before.
GetCookie() is the function that we query redis to get cookies which was set before.
[Start url](https://weibo.com/2028810631/follow) consists of the base url, the user ID and "/follow".(you propably can't access the link ahead because of lacking cookies)
Once accessed, the crawler will search for script that contains the web information then find the matching html node and next page url, finally get the follows' id for futher use.   

#### fetchPersonalPage
Fetch all weibos on first page of the specified user.
This function uses visitor cookies generated in '/utils/cookie.go'.
[Start url](https://weibo.com/u/2028810631?is_all=1) consists base url, "/u", user id and "?is_all=1" so that blogs will order by time.
In case that the blog is too long that some contents not included in the scripts, the crawler only finds detailed blog url in this page.
Considering the timeliness of the blogs, the crawler only fetches first page so that visitor cookies are enough. 
There are about 16 weibos at the homepage of the user.

#### fetchDetails
Fetch blogger id, avaar, name, time, images, video, tags, article, origin weibo if forward.
This funcitons uses visitor cookies generated in '/utils/cookie.go'.
[Url](https://weibo.com/2028810631/JuhgE9y9K) consists base url, user id, blog id.
Tags are supertopics, topics and keywords from blog text or database. They will be stored as the feature of the blog for different categories.
If the blog forwarded, the origin blog will also be recorded in database.



