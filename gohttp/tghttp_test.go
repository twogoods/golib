package gohttp

import (
	"encoding/json"
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
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).EnableCookie(true).Build()
	resp, err := client.Execute(getreq)
	if err == nil {
		str, _ := resp.BodyString()
		t.Log(str)
	} else {
		t.Error(err)
	}
}

func TestPost(t *testing.T) {
	postbody := FormBodyBuilder().AddParam("name", "kobe").AddParam("age", "23").Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).EnableCookie(true).Build()
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
	req, _ := RequestBuilder().Header("head-k", "head-v").Url("http://httpbin.org/post").Post(jsonbody).Build()
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).EnableCookie(true).Build()
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		t.Log(str)
	} else {
		t.Error(err)
	}
}

func TestFile(t *testing.T) {
	filebody := MuiltBodyBuilder().AddFile("test", "/Users/twogoods/python/haha.txt").AddParam("key", "twogoods").Build()
	req, _ := RequestBuilder().Url("http://127.0.0.1:8888/upload").Post(filebody).Build()
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
	postbody := FormBodyBuilder().AddParam("name", "kobe").AddParam("age", "23").Build()
	req, _ := RequestBuilder().Header("head-k", "head-v").Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).EnableCookie(true).Build()
	resp, err := client.Submit(req).Get()
	if err != nil {
		t.Log(err)
	} else {
		str, _ := resp.BodyString()
		t.Log("return:" + str)
	}
	if err == nil {
		str, _ := resp.BodyString()
		t.Log(str)
	} else {
		t.Error(err)
	}
}

func TestAsync(t *testing.T) {
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
	time.Sleep(10 * time.Second)
}

func BenchmarkPost(b *testing.B) {
	postbody := FormBodyBuilder().AddParam("name", "kobe").AddParam("age", "23").Build()
	req, err := RequestBuilder().Header("head-k", "head-v").Url("http://httpbin.org/post").Post(postbody).Build()
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).EnableCookie(true).Build()
	resp, err := client.Execute(req)
	if err == nil {
		str, _ := resp.BodyString()
		b.Log(str)
	} else {
		b.Error(err)
	}
}
func Testcrawler(t *testing.T) {
	loginurl := "http://www.ycjw.zjut.edu.cn//logon.aspx?__EVENTTARGET=&__EVENTARGUMENT=&__VIEWSTATE=dDwtMTU2MDM2OTk5Nzt0PDtsPGk8MT47PjtsPHQ8O2w8aTwzPjtpPDEzPjs%2BO2w8dDw7bDxpPDE%2BO2k8Mz47aTw1PjtpPDc%2BO2k8OT47aTwxMT47aTwxMz47aTwxNT47aTwxNz47PjtsPHQ8cDxwPGw8QmFja0ltYWdlVXJsOz47bDxodHRwOi8vd3d3Lnljancuemp1dC5lZHUuY24vL2ltYWdlcy9iZy5naWY7Pj47Pjs7Pjt0PHA8cDxsPEJhY2tJbWFnZVVybDs%2BO2w8aHR0cDovL3d3dy55Y2p3LnpqdXQuZWR1LmNuLy9pbWFnZXMvYmcxLmdpZjs%2BPjs%2BOzs%2BO3Q8cDxwPGw8QmFja0ltYWdlVXJsOz47bDxodHRwOi8vd3d3Lnljancuemp1dC5lZHUuY24vL2ltYWdlcy9iZzEuZ2lmOz4%2BOz47Oz47dDxwPHA8bDxCYWNrSW1hZ2VVcmw7PjtsPGh0dHA6Ly93d3cueWNqdy56anV0LmVkdS5jbi8vaW1hZ2VzL2JnMS5naWY7Pj47Pjs7Pjt0PHA8cDxsPEJhY2tJbWFnZVVybDs%2BO2w8aHR0cDovL3d3dy55Y2p3LnpqdXQuZWR1LmNuLy9pbWFnZXMvYmcxLmdpZjs%2BPjs%2BOzs%2BO3Q8cDxwPGw8QmFja0ltYWdlVXJsOz47bDxodHRwOi8vd3d3Lnljancuemp1dC5lZHUuY24vL2ltYWdlcy9iZzEuZ2lmOz4%2BOz47Oz47dDxwPHA8bDxCYWNrSW1hZ2VVcmw7PjtsPGh0dHA6Ly93d3cueWNqdy56anV0LmVkdS5jbi8vaW1hZ2VzL2JnMS5naWY7Pj47Pjs7Pjt0PHA8cDxsPEJhY2tJbWFnZVVybDs%2BO2w8aHR0cDovL3d3dy55Y2p3LnpqdXQuZWR1LmNuLy9pbWFnZXMvYmcxLmdpZjs%2BPjs%2BOzs%2BO3Q8cDxwPGw8QmFja0ltYWdlVXJsOz47bDxodHRwOi8vd3d3Lnljancuemp1dC5lZHUuY24vL2ltYWdlcy9iZzEuZ2lmOz4%2BOz47Oz47Pj47dDx0PDt0PGk8Mz47QDwtLeeUqOaIt%2Bexu%2BWeiy0tO%2BaVmeW4iDvlrabnlJ87PjtAPC0t55So5oi357G75Z6LLS075pWZ5biIO%2BWtpueUnzs%2BPjs%2BOzs%2BOz4%2BOz4%2BO2w8SW1nX0RMOz4%2Bqmizg8nuU1ebhUFzNA%2Fqu71sECk%3D&Img_DL.x=23&Img_DL.y=12" +
		"&Cbo_LX=" + "%d1%a7%c9%fa" +
		"&Txt_UserName=" + "201203870310" +
		"&Txt_Password=" + "682481lu"
	t.Log(loginurl)
	req, err := RequestBuilder().Url(loginurl).Build()
	//create client
	client := HttpClientBuilder().SetTimeout(20*time.Second, 20*time.Second).EnableCookie(true).Build()
	resp, err := client.Execute(req)
	if err == nil {
		resp.BodyString()
		t.Log("success:")
	} else {
		t.Log(err)
	}

	t.Log("resp cookies:", resp.Response.Cookies())

	postbody := FormBodyBuilder().AddParam("__EVENTTARGET", "").
		AddParam("Cbo_LX", "").
		AddParam("__VIEWSTATE", "dDw5NjExNjI1OTE7dDw7bDxpPDE+Oz47bDx0PDtsPGk8NT47aTw3PjtpPDk+O2k8MTU+O2k8MTk+O2k8MjE+Oz47bDx0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w85a2m55Sf566h55CGXD7mn6Xor6Lns7vnu59cPuWtpueUn+ivvuihqOafpeivojs+Pjs+Ozs+Oz4+O3Q8cDxwPGw8VGV4dDs+O2w8MjAxNi0wNC0xNCAwODoyMjowOTs+Pjs+Ozs+O3Q8dDw7dDxpPDQxPjtAPDE5OTkvMjAwMCgxKTsxOTk5LzIwMDAoMik7MjAwMC8yMDAxKDEpOzIwMDAvMjAwMSgyKTsyMDAxLzIwMDIoMSk7MjAwMS8yMDAyKDIpOzIwMDIvMjAwMygxKTsyMDAyLzIwMDMoMik7MjAwMy8yMDA0KDEpOzIwMDMvMjAwNCgyKTsyMDA0LzIwMDUoMSk7MjAwNC8yMDA1KDIpOzIwMDUvMjAwNigxKTsyMDA1LzIwMDYoMik7MjAwNi8yMDA3KDEpOzIwMDYvMjAwNygyKTsyMDA3LzIwMDgoMSk7MjAwNy8yMDA4KDIpOzIwMDgvMjAwOSgxKTsyMDA4LzIwMDkoMik7MjAwOS8yMDEwKDEpOzIwMDkvMjAxMCgyKTsyMDEwLzIwMTEoMSk7MjAxMC8yMDExKDIpOzIwMTEvMjAxMigxKTsyMDExLzIwMTIoMik7MjAxMi8yMDEzKDEpOzIwMTIvMjAxMygyKTsyMDEzLzIwMTQoMSk7MjAxMy8yMDE0KDIpOzIwMTQvMjAxNSgxKTsyMDE0LzIwMTUoMik7MjAxNS8yMDE2KDEpOzIwMTUvMjAxNigyKTsyMDE2LzIwMTcoMSk7MjAxNi8yMDE3KDIpOzIwMTcvMjAxOCgxKTsyMDE3LzIwMTgoMik7MjAxOC8yMDE5KDEpOzIwMTgvMjAxOSgyKTsyMDk4LzIwOTkoMSk7PjtAPDE5OTkvMjAwMCgxKTsxOTk5LzIwMDAoMik7MjAwMC8yMDAxKDEpOzIwMDAvMjAwMSgyKTsyMDAxLzIwMDIoMSk7MjAwMS8yMDAyKDIpOzIwMDIvMjAwMygxKTsyMDAyLzIwMDMoMik7MjAwMy8yMDA0KDEpOzIwMDMvMjAwNCgyKTsyMDA0LzIwMDUoMSk7MjAwNC8yMDA1KDIpOzIwMDUvMjAwNigxKTsyMDA1LzIwMDYoMik7MjAwNi8yMDA3KDEpOzIwMDYvMjAwNygyKTsyMDA3LzIwMDgoMSk7MjAwNy8yMDA4KDIpOzIwMDgvMjAwOSgxKTsyMDA4LzIwMDkoMik7MjAwOS8yMDEwKDEpOzIwMDkvMjAxMCgyKTsyMDEwLzIwMTEoMSk7MjAxMC8yMDExKDIpOzIwMTEvMjAxMigxKTsyMDExLzIwMTIoMik7MjAxMi8yMDEzKDEpOzIwMTIvMjAxMygyKTsyMDEzLzIwMTQoMSk7MjAxMy8yMDE0KDIpOzIwMTQvMjAxNSgxKTsyMDE0LzIwMTUoMik7MjAxNS8yMDE2KDEpOzIwMTUvMjAxNigyKTsyMDE2LzIwMTcoMSk7MjAxNi8yMDE3KDIpOzIwMTcvMjAxOCgxKTsyMDE3LzIwMTgoMik7MjAxOC8yMDE5KDEpOzIwMTgvMjAxOSgyKTsyMDk4LzIwOTkoMSk7Pj47bDxpPDMyPjs+Pjs7Pjt0PHA8cDxsPFRleHQ7PjtsPOWtpueUnzrljaLosarluIXnmoTor77ooajkv6Hmga87Pj47Pjs7Pjt0PEAwPHA8cDxsPFBhZ2VDb3VudDtfIUl0ZW1Db3VudDtfIURhdGFTb3VyY2VJdGVtQ291bnQ7RGF0YUtleXM7PjtsPGk8MT47aTw0PjtpPDQ+O2w8aTwyMTUwNj47aTwxNzE4Mz47aTwxMjA1MT47aTwxMzg2MT47Pjs+Pjs+Ozs7Ozs7Ozs7Oz47bDxpPDA+Oz47bDx0PDtsPGk8MT47aTwyPjtpPDM+O2k8ND47PjtsPHQ8O2w8aTwwPjtpPDE+O2k8Mj47aTwzPjtpPDQ+O2k8NT47PjtsPHQ8O2w8aTwxPjs+O2w8dDxwPHA8bDxUZXh0Oz47bDzlpKflrabnlJ/ogYzkuJrlj5HlsZXkuI7lsLHkuJrmjIflr7zvvIjkuIvvvInihaA66YeR6K+X5Y2XIOiDoeWtn+adsCA7Pj47Pjs7Pjs+Pjt0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w85a2m55Sf5aSEOz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPDAuNTs+Pjs+Ozs+Oz4+O3Q8O2w8aTwxPjs+O2w8dDxwPHA8bDxUZXh0Oz47bDw4Oz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPFxlOz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPOWFrOmAieivvjs+Pjs+Ozs+Oz4+Oz4+O3Q8O2w8aTwwPjtpPDE+O2k8Mj47aTwzPjtpPDQ+O2k8NT47PjtsPHQ8O2w8aTwxPjs+O2w8dDxwPHA8bDxUZXh0Oz47bDzmnLrlmajkurrmjqfliLbihaA656a56ZGr54eaIOasp+ael+aelyA7Pj47Pjs7Pjs+Pjt0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w85L+h5oGv5bel56iL5a2m6ZmiOz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPDM7Pj47Pjs7Pjs+Pjt0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w8NDg7Pj47Pjs7Pjs+Pjt0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w8MS0xMuWRqDrmmJ/mnJ8yKDMtNCkgIOS7geWSjDIwMlw7MS0xMuWRqDrmmJ/mnJ81KDMtNCkgIOS7geWSjDIwNlw7Oz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPOmZkOmAieivvjs+Pjs+Ozs+Oz4+Oz4+O3Q8O2w8aTwwPjtpPDE+O2k8Mj47aTwzPjtpPDQ+O2k8NT47PjtsPHQ8O2w8aTwxPjs+O2w8dDxwPHA8bDxUZXh0Oz47bDznp5Hlrabmlrnms5XorrrihaA65pyx5pix5rW3ICA7Pj47Pjs7Pjs+Pjt0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w85pS/5rK75LiO5YWs5YWx566h55CG5a2m6ZmiOz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPDI7Pj47Pjs7Pjs+Pjt0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w8MzI7Pj47Pjs7Pjs+Pjt0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w8MS0xNuWRqDrmmJ/mnJ8zKDYtNykgIOWNmuaYk0M0MDRcOzs+Pjs+Ozs+Oz4+O3Q8O2w8aTwxPjs+O2w8dDxwPHA8bDxUZXh0Oz47bDzlv4Xkv67or747Pj47Pjs7Pjs+Pjs+Pjt0PDtsPGk8MD47aTwxPjtpPDI+O2k8Mz47aTw0PjtpPDU+Oz47bDx0PDtsPGk8MT47PjtsPHQ8cDxwPGw8VGV4dDs+O2w85L2T6IKyOuabueerueiQjSAgOz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPOS9k+WGm+mDqDs+Pjs+Ozs+Oz4+O3Q8O2w8aTwxPjs+O2w8dDxwPHA8bDxUZXh0Oz47bDwxOz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPDMyOz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPDEtNDrmmJ/mnJ81KDYtNikgOz4+Oz47Oz47Pj47dDw7bDxpPDE+Oz47bDx0PHA8cDxsPFRleHQ7PjtsPOmZkOmAieivvjs+Pjs+Ozs+Oz4+Oz4+Oz4+Oz4+O3Q8QDA8cDxwPGw8UGFnZUNvdW50O18hSXRlbUNvdW50O18hRGF0YVNvdXJjZUl0ZW1Db3VudDtEYXRhS2V5czs+O2w8aTwwPjtpPC0xPjtpPC0xPjtsPD47Pj47Pjs7Ozs7Ozs7Ozs+Ozs+Oz4+Oz4+Oz5jXZNEkKNQcbzqr3qilRfxGP9xGg==").
		AddParam("Cbo_Xueqi", "2015/2016(2)").AddParam("Button1", "按课程查询").Build()

	get, _ := RequestBuilder().Url("http://www.ycjw.zjut.edu.cn//stdgl/cxxt/Web_Std_XQKB.aspx").Post(postbody).Build()
	getresp, err := client.Execute(get)
	if err == nil {
		str, _ := getresp.BodyString()
		//getresp.BodyString()
		t.Log("success:", str)
	} else {
		t.Log("err:", err)
	}
}
