package gohttp

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	obj := make(map[string][]string)
	obj["name"] = []string{"张三"}
	obj["age"] = []string{"23"}
	url := BuildGetUrl("http://httpbin.org/get", obj)
	t.Log(url)
	getreq, _ := RequestBuilder().Header("head-k", "head-v").Url(url).Build()
	client := HttpClientBuilder().Build()
	resp, err := client.Execute(getreq)
	if err == nil {
		str, _ := resp.BodyString()
		t.Log(str)
	} else {
		t.Error(err)
	}
}

func TestPost(t *testing.T) {
	postbody := FormBodyBuilder().AddParam("name", "test").AddParam("age", "23").Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").
		Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().Build()
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		t.Log(str)
	} else {
		t.Error(err)
	}
}

func TestJson(t *testing.T) {
	obj := make(map[string]string)
	obj["name"] = "test"
	obj["age"] = "23"
	data, _ := json.Marshal(obj)
	jsonbody := JsonBodyBuilder().Json(data).Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").
		Url("http://httpbin.org/post").Post(jsonbody).Build()
	client := HttpClientBuilder().Build()
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		t.Log(str)
	} else {
		t.Error(err)
	}
}

func TestFile(t *testing.T) {
	filebody := MuiltBodyBuilder().AddFile("test", "/Users/twogoods/python/haha.txt").
		AddParam("key", "twogoods").Build()
	req, _ := RequestBuilder().Url("http://127.0.0.1:8888/upload").
		Post(filebody).Build()
	client := HttpClientBuilder().Build()
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		t.Log(str)
	} else {
		t.Error(err)
	}
}

func TestSetCookie(t *testing.T) {
	cookie := &http.Cookie{}
	cookie.Name = "username"
	cookie.Value = "twogoods"
	postbody := FormBodyBuilder().AddParam("name", "test").Build()
	req, _ := RequestBuilder().Header("Cookie", "pid=123234").
		Cookie(cookie).Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().Build()
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		t.Log(str)
	} else {
		t.Error(err)
	}
}

func TestSubmit(t *testing.T) {
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
		t.Log("return:" + str)
	}
}

func TestAsync(t *testing.T) {
	postbody := FormBodyBuilder().AddParam("name", "test").
		AddParam("age", "23").Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").
		Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().Build()
	client.Async(req,
		func(resp *TGHttpResponse) {
			result, _ := resp.BodyString()
			t.Log(result)
		},
		func(err error) {
			t.Log("test err")
			t.Error(err)
		})
	time.Sleep(10 * time.Second)
}

func BenchmarkPost(b *testing.B) {
	postbody := FormBodyBuilder().AddParam("name", "test").Build()
	req, err := RequestBuilder().Header("head-k", "head-v").
		Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().Build()
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		b.Log(str)
	} else {
		b.Error(err)
	}
}
