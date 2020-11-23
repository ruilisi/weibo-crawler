## utils

### agent.go

```go
func RandomUserAgent() string {
	userAgent := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36",
	}
	return userAgent[rand.Intn(len(userAgent))]
}
```

RandomUserAgent return one of preset user-agent for crawler

### cookie.go

Go's native cookie type differs from the type of cookie you get on the web from save_cookie.py, so it's overridden.
Simulates a browser environment,visit weibo visitor id [generating page](https://passport.weibo.com/visitor/genvisitor), obtains a time-sensitive visitor id, uses that id to access the page, and obtains cookies for that visitor id.




