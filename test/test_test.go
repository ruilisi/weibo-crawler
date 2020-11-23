package test_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test", func() {

	Describe("test", func() {
		var result map[string]interface{}
		Context("Add blogger", func() {
			It("success", func() {
				params := map[string][]string{
					"bloggers": {"5467852665"},
				}
				res := Post(params, "/add_bloggers", "")
				res.ToJSON(&result)
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				Ω(fmt.Sprint(result["status"])).To(Equal("success"))
			})
		})
		Context("set keywords & category", func() {
			It("success", func() {
				params := map[string][]string{
					"keywords": {"中国", "国家", "新闻"},
				}
				res := Post(params, "/tags/set_keywords", "")
				res.ToJSON(&result)
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				Ω(fmt.Sprint(result["status"])).To(Equal("success"))

				res = Post(nil, "/tags/keywords", "")
				res.ToJSON(&result)
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				Ω(fmt.Sprint(result["status"])).To(Equal("success"))

				// param := map[string]interface{}{
				// "category": "0:[中国, 国家]/中国\n1:[新闻]/新闻\n",
				// }
				// res = Post(param, "/category/set", "")
				// res.ToJSON(&result)
				// Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				// Ω(fmt.Sprint(result["status"])).To(Equal("success"))
			})

		})
		Context("task", func() {
			It("success", func() {
				res := Get(nil, "/task", "")
				res.ToJSON(&result)
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				Ω(fmt.Sprint(result["status"])).To(Equal("success"))
			})
		})
		Context("query", func() {
			It("success", func() {
				params := map[string]interface{}{
					"amount": "10",
				}
				res := Post(params, "/query_blogs", "")
				res.ToJSON(&result)
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				Ω(fmt.Sprint(result["status"])).To(Equal("success"))
				rs := result["blogs"].([]interface{})
				Ω(len(rs)).To(Equal(10))
				Ω(rs[0].(map[string]interface{})["BloggerID"]).To(Equal("5467852665"))
				Ω(rs[0].(map[string]interface{})["Follows"]).Should(BeNumerically(">=", 11898602.0))
				Ω(rs[0].(map[string]interface{})["Avatar"]).Should(HavePrefix("https://tvax2.sinaimg.cn/crop.0.0.996.996.180/005Y2yg1ly8gdi7zkua63j30ro0ro407.jpg"))
			})
		})
		// Context("query hot topics", func() {
		// It("success", func() {
		// res := Get(nil, "category/query", "")
		// res.ToJSON(&result)
		// Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
		// Ω(fmt.Sprint(result["status"])).To(Equal("success"))
		// rs := result["categories"].([]interface{})
		// Ω(len(rs)).To(Equal(3))
		// })
		// })
	})
})
