## test

### install
`go get -u github.com/onsi/ginkgo/ginkgo`

### how to use
run in test folder<br>
ginkgo          test all files<br>
ginkgo -focus=SignIn          test describe named SignIn
```go
var _ = Describe("SignIn", func() {
...  
})
```
