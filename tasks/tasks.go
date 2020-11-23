package tasks

import (
	"fmt"
	"go-crawler/db"
	"go-crawler/models"
	"go-crawler/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/lib/pq"
)

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
	personalReg    = `<div class=\\\"WB_from S_txt2\\\">(.*?)<\\\/div>`

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

//FetchFollows get his follow list
func FetchFollows(userID string) []string {
	c := colly.NewCollector()
	followIDs := []string{}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomUserAgent())
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		followReg, _ := regexp.Compile(followReg)
		doc := followReg.FindString(string(r.Body))

		followHTMLReg, _ := regexp.Compile(followHTMLReg)
		docs := followHTMLReg.FindAllString(doc, -1)
		for _, d := range docs {
			dochtml, _ := htmlquery.Parse(strings.NewReader(replace(d)))
			members := htmlquery.Find(dochtml, followMemberNode)
			for _, member := range members {
				followURL := htmlquery.SelectAttr(member, "usercard")
				follow, _ := regexp.Compile(followIDReg)
				urlID := follow.FindString(followURL)
				if len(urlID) != 0 {
					followIDs = append(followIDs, urlID[3:])
				}
			}
		}

		//next page
		followNext, _ := regexp.Compile(followNextPageReg)
		nextHTMLStr := followNext.FindString(doc)
		nextHTML, _ := htmlquery.Parse(strings.NewReader(replace(nextHTMLStr)))
		url := htmlquery.SelectAttr(htmlquery.FindOne(nextHTML, followNextPageNode), "href")
		if url != "" {
			url = baseURL + url[1:]
			cookies, err := utils.GetCookie()
			if err == nil {
				c.SetCookies(url, cookies)
				c.Visit(url)
			} else {
				log.Println(err)
			}
		}
	})

	url := baseURL + userID + "/follow"
	cookies, err := utils.GetCookie()
	if err == nil {
		c.SetCookies(url, cookies)
		c.Visit(url)
	} else {
		log.Fatal(err)
	}
	return followIDs
}

//Fetch crawler job
func Fetch() {
	bloggers, rs := models.FindBloggers()
	if rs.Error != nil {
		return
	} else if rs.NotFound() {
		return
	}

	for _, blogger := range bloggers {
		FetchBloggerDetails(blogger.ID)
		blogsIndex := FetchPersonalPage(blogger.ID)
		for _, blog := range blogsIndex {
			FetchBlogDetails(blogger.ID, blog)
		}
	}
}

func FetchPersonalPage(userID string) []string {
	c := colly.NewCollector()
	contentIDs := []string{}
	lastest, res := models.FindLastest(userID)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomUserAgent())
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		re := regexp.MustCompile(personalReg)
		matches := re.FindAllSubmatch(r.Body, -1)
		for _, match := range matches {
			detail, _ := htmlquery.Parse(strings.NewReader(replace(string(match[0]))))
			hrefNode := htmlquery.Find(detail, "//a[1]/@href")
			titleNode := htmlquery.Find(detail, "//a[1]/@title")
			if len(hrefNode) > 0 && len(titleNode) > 0 {
				href := htmlquery.InnerText(hrefNode[0])
				hrefRe := regexp.MustCompile(`/([0-9a-zA-Z]+)`)
				idMatches := hrefRe.FindAllSubmatch([]byte(href), -1)
				if len(idMatches) == 2 {
					if string(idMatches[0][1]) == userID {
						timeString := htmlquery.InnerText(titleNode[0])
						timeStamp, _ := time.ParseInLocation("2006-01-02 15:04", timeString, time.Local)
						if res.Found() {
							if timeStamp.After(lastest.Timestamp) {
								contentIDs = append(contentIDs, string(idMatches[1][1]))
							}
						} else {
							contentIDs = append(contentIDs, string(idMatches[1][1]))
						}
					}
				}
			}
		}
	})

	url := baseURL + "u/" + userID + "?is_all=1"
	cookies, err := utils.GetCookie()
	if err == nil {
		c.SetCookies(url, cookies)
		c.Visit(url)
	} else {
		log.Fatal(err)
	}
	return contentIDs
}

func FetchBlogDetails(userID, contentID string) error {
	_, result := models.FindBlogByID(contentID)
	if result.Found() {
		return fmt.Errorf("blog exists")
	}

	var tags [3][]string
	var fetchErr error
	var video, article string
	var imgs pq.StringArray

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomUserAgent())
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		rs := ""
		//		retry := 2
		sreg, _ := regexp.Compile(weiboDetailReg)
		script := sreg.FindString(string(r.Body))

		texthtmlreg, _ := regexp.Compile(textHTMLReg)
		html := texthtmlreg.FindString(script)

		//fetch imgs
		imgReg, _ := regexp.Compile(imgNodeReg)
		imgshtml := imgReg.FindAllString(script, -1)
		detailreg, _ := regexp.Compile("mw690|thumb150") //mh690, thumb150 and orj360 are marks for images of different sizes in image hosting server
		for _, imghtml := range imgshtml {
			img, _ := htmlquery.Parse(strings.NewReader(replace(imghtml)))
			imgurl := htmlquery.FindOne(img, imgNode)
			url := htmlquery.SelectAttr(imgurl, "src")
			url = detailreg.ReplaceAllString(url, "orj360")
			imgs = append(imgs, url)
		}

		rtReg, _ := regexp.Compile(rtNodeReg)
		rthtml := rtReg.FindString(script)
		rt, _ := htmlquery.Parse(strings.NewReader(replace(rthtml)))
		retweet := htmlquery.SelectAttr(htmlquery.FindOne(rt, rtNode), "href")
		if !strings.HasPrefix(retweet, "/") {
			retweet = ""
		}

		//fetch timestamp
		timereg, _ := regexp.Compile(timeNodeReg)
		timehtml := timereg.FindString(script)
		timereg, _ = regexp.Compile(timeReg)
		timeString := timereg.FindString(timehtml)
		if timeString != "" {
			timeString = timeString[7:17]
		}
		unix, _ := strconv.ParseInt(timeString, 10, 64)
		timeStamp := time.Unix(unix, 0)

		detail, _ := htmlquery.Parse(strings.NewReader(replace(html)))
		textnode := htmlquery.Find(detail, textNode)
		if len(textnode) == 0 {
			fetchErr = fmt.Errorf("%s/%s empty content", userID, contentID)
			return
		}

		spreg, _ := regexp.Compile(spReg)
		sphtml := spreg.FindString(script)
		spdetail, _ := htmlquery.Parse(strings.NewReader(replace(sphtml)))
		spnode := htmlquery.FindOne(spdetail, spNode)

		//fetch video
		if htmlquery.FindOne(detail, videoNode) != nil {
			video = htmlquery.SelectAttr(htmlquery.FindOne(detail, videoNode).Parent, "href")
		}

		//fetch article
		if htmlquery.FindOne(detail, articleNode) != nil {
			article = htmlquery.SelectAttr(htmlquery.FindOne(detail, articleNode).Parent, "href")
		}

		for _, text := range textnode {
			rs += text.Data
		}
		rs = strings.TrimSpace(rs)
		if rs == "" {
			fetchErr = fmt.Errorf("%s/%s empty content", userID, contentID)
			return
		}
		if spnode == nil {
			tags[0] = fetchSuperTopic(rs)
		} else {
			//tags[0] = []string{"" + spnode.Data}
			tags[0] = []string{spnode.Data}
		}
		tags[1] = fetchTopic(rs)
		tags[2] = fetchKeyword(strings.ToLower(rs))
		tagIDs := pq.StringArray{}
		for tagType, tagNames := range tags {
			for _, tagName := range tagNames {
				if tagType != models.KEYWORD {
					tag, res := models.CreateTag(tagName, tagType)
					if res.Found() {
						tagIDs = append(tagIDs, tag.ID.String())
					} else {
						fetchErr = res.Err()
						return
					}
				} else {
					tag, res := models.FindTag(tagName, tagType)
					if res.Found() {
						tagIDs = append(tagIDs, tag.ID.String())
					} else {
						fetchErr = res.Err()
						return
					}
				}
			}
		}

		blog := models.Blog{ID: contentID, Content: rs, Tags: tagIDs, BloggerID: userID, Timestamp: timeStamp, Pictures: imgs, Video: video, Article: article, Rt: retweet}
		fetchErr = db.DB.Create(&blog).Error
	})

	url := baseURL + userID + "/" + contentID
	cookies, err := utils.GetCookie()
	if err == nil {
		c.SetCookies(url, cookies)
		c.Visit(url)
	} else {
		return err
	}
	return fetchErr
}

func FetchBloggerDetails(userID string) error {
	var username, avatar string
	var follows int64
	var fetchErr error
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomUserAgent())
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		reUsername := regexp.MustCompile(nameReg)
		usernameM := reUsername.FindAllSubmatch(r.Body, -1)
		if len(usernameM) != 0 {
			username = string(usernameM[0][1])
		} else {
			fetchErr = fmt.Errorf("fail to fetch blogger username")
			return
		}

		reAvatar := regexp.MustCompile(avatarReg)
		avatarM := reAvatar.FindAllSubmatch(r.Body, -1)
		if len(avatarM) != 0 {
			avatarDetail, _ := htmlquery.Parse(strings.NewReader(replace(string(avatarM[0][0]))))
			srcNode := htmlquery.Find(avatarDetail, "//img/@src")
			if len(srcNode) > 0 {
				avatar = htmlquery.InnerText(srcNode[0])
			}
		}
		if avatar == "" {
			fetchErr = fmt.Errorf("fail to fetch blogger avatar")
			return
		}

		reFollower := regexp.MustCompile(followTableReg)
		followerM := reFollower.FindAllSubmatch(r.Body, -1)
		if len(followerM) != 0 {
			followerDetail, _ := htmlquery.Parse(strings.NewReader(replace(string(followerM[0][0]))))
			followerNode := htmlquery.Find(followerDetail, "//strong")
			if len(followerNode) == 3 {
				follows, fetchErr = strconv.ParseInt(htmlquery.InnerText(followerNode[1]), 10, 64)
				if fetchErr != nil {
					return
				}
			}
		}

		blogger, result := models.FindBloggerByID(userID)
		if result.Found() {
			blogger.Username = username
			blogger.Avatar = avatar
			blogger.Follows = follows
			db.DB.Save(&blogger)
		} else {
			fetchErr = result.Err()
			return
		}
	})

	url := baseURL + "u/" + userID
	cookies, err := utils.GetCookie()
	if err == nil {
		c.SetCookies(url, cookies)
		c.Visit(url)
	} else {
		fetchErr = err
	}
	return fetchErr
}

func replace(s string) string {
	rs := strings.Replace(s, "\\n", "\n", -1)
	rs = strings.Replace(rs, "\\t", "\t", -1)
	rs = strings.Replace(rs, "\\r", "\r", -1)
	rs = strings.Replace(rs, "\\", "", -1)
	return rs
}
