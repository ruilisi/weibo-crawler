## router
set listening port
```go
router.Run(fmt.Sprintf(":%v", conf.GetEnv("HTTP_PORT")))
```
set router group
```go
groups := router.Group("/group")
```
set middleware
```go
groups.Use(middleware.Auth("group"))
```
