package test_test

import (
	"fmt"
	"go-crawler/db"
	"go-crawler/models"
	"go-crawler/redis"
	"go-crawler/routers"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"unsafe"

	"github.com/alicebob/miniredis"
	"github.com/imroc/req"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var Ts *httptest.Server

var token string
var accessToken string

func Params(method, api string, params url.Values) *http.Response {
	req, _ := http.NewRequest(method, Ts.URL+"/"+api, nil)
	if len(params) > 0 {
		req.URL.RawQuery = params.Encode()
	}
	resp, _ := http.DefaultClient.Do(req)
	return resp
}

func headers(token string) req.Header {
	return req.Header{
		"Accept":        "application/json",
		"Authorization": token,
	}
}

func Get(body interface{}, s, token string) *req.Resp {
	resp, _ := req.Get(Ts.URL+s, headers(token), req.BodyJSON(&body))
	return resp
}

func Post(body interface{}, s, token string) *req.Resp {
	resp, _ := req.Post(Ts.URL+s, headers(token), req.BodyJSON(&body))
	return resp
}

func Put(body interface{}, s, token string) *req.Resp {
	resp, _ := req.Put(Ts.URL+s, headers(token), req.BodyJSON(&body))
	return resp
}

func Delete(body interface{}, s, token string) *req.Resp {
	resp, _ := req.Delete(Ts.URL+s, headers(token), req.BodyJSON(&body))
	return resp
}

func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func BoxIn() *gorm.DB {
	tmp := db.DB
	if db.DB != nil {
		fmt.Println("formal", tmp)
		db.DB = db.DB.Begin()
	}
	return tmp
}

func BoxOut(dbCache *gorm.DB) {
	if db.DB != nil {
		db.DB.Rollback()
	}
	db.DB = dbCache
}

func TestGoCrawler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go crawler Suite")
}

var dbcache *gorm.DB
var MockRedis *miniredis.Miniredis
var _ = BeforeSuite(func() {
	var err error
	MockRedis, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	redis.TestConnectRedis("redis://" + MockRedis.Addr())

	tags := []string{"中国", "国家", "新闻"}
	redis.Del("keywords")
	for _, tag := range tags {
		redis.SAdd("keywords", tag)
	}

	db.Open("test")
	db.Migrate("test", &models.Blog{}, &models.Blogger{}, &models.Tag{}, &models.Category{})

	dbcache = BoxIn()
	Ts = httptest.NewServer(routers.SetupRouter())
})

var _ = AfterSuite(func() {
	BoxOut(dbcache)
	MockRedis.Close()
	db.CleanTablesData()
	db.Close()
})
