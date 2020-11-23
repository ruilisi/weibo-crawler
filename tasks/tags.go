package tasks

import (
	"go-crawler/redis"
	"regexp"
	"strings"
)

func fetchTopic(s string) []string {
	reg, _ := regexp.Compile("#.*?#")
	rs := reg.FindAllString(s, -1)
	var result []string
	for _, r := range rs {
		result = append(result, r[1:len(r)-1])
	}
	return result
}

func fetchSuperTopic(s string) []string {
	reg, _ := regexp.Compile("\\S+超话")
	rs := reg.FindAllString(s, -1)
	if len(rs) > 1 {
		rs = rs[:1]
	}
	return rs
}

func fetchKeyword(s string) []string {
	rs := []string{}
	for _, keyword := range redis.Smembers("keywords") {
		if strings.Contains(s, keyword) {
			rs = append(rs, keyword)
		}
	}
	return rs
}
