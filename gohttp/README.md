# gohttp
像okhttp一样在go中使用httpclient

# install
`go get github.com/twogoods/golib/gohttp`  

# quickstart
## GET
```
	//create request
	postbody := FormBodyBuilder().AddParam("name", "kobe").AddParam("age", "23").Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").
		Url("http://httpbin.org/post").Post(postbody).Build()
	// create httpclient
	client := HttpClientBuilder().Build()
	//do request
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		fmt.Println(str)
	} else {
		fmt.Println(err)
	}
```

## POST
```
	postbody := FormBodyBuilder().AddParam("name", "test").AddParam("age", "23").Build()
	// set header here
	req, _ := RequestBuilder().Header("head-k", "head-v").
		Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().Build()
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		fmt.Println(str)
	} else {
		fmt.Println(err)
	}
```

## Upload file
```
	filebody := MuiltBodyBuilder().AddFile("test", "/Users/twogoods/haha.txt").
		AddParam("key", "twogoods").Build()
	req, _ := RequestBuilder().Url("http://127.0.0.1:8888/upload").
		Post(filebody).Build()
	client := HttpClientBuilder().Build()
	client.Execute(req)
```

## Cookie
auto manage cookie,just use `EnableCookie(true)`
```
	postbody := FormBodyBuilder().AddParam("name", "test").AddParam("age", "23").Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").
		Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).
		EnableCookie(true).Build()
	client.Execute(req)

	// set cookie in request like this
	cookie := &http.Cookie{}
	cookie.Name = "username"
	cookie.Value = "twogoods"
	postbody := FormBodyBuilder().AddParam("name", "test").Build()
	req, _ := RequestBuilder().Header("Cookie", "pid=123234").
		Cookie(cookie).Url("http://httpbin.org/post").Post(postbody).Build()
```

## submit
submit the request,wait the request finish
```
	postbody := FormBodyBuilder().AddParam("name", "test").AddParam("age", "23").Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").
		Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).
		EnableCookie(true).Build()
	resp, err := client.Submit(req).Get()
	if err != nil {
		t.Log(err)
	} else {
		str, _ := resp.BodyString()
		fmt.Println("return:" + str)
	}
```

## Async
use async thread to request
```
	postbody := FormBodyBuilder().AddParam("name", "test").AddParam("age", "23").Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).EnableCookie(true).Build()
	client.Async(req,
		func(resp *TGHttpResponse) {
			result, _ := resp.BodyString()
			t.Log(result)
		},
		func(err error) {
			t.Log("test err")
			t.Error(err)
		})
```
for more sample see `tghttp_test.go`
