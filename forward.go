package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	protocolHttps = "https://"
	protocolHttp  = "http://"
	host          = "xihe.mindspore.cn"
	hostTest      = "xihe2.test.osinfra.cn"
	hostLocal     = "127.0.0.1:8080"
	poolHost      = ".pool1.mindspore.cn"

	typeCloud     = "cloud"
	typeInference = "inference"
	typeEvaluate  = "evaluate"
)

func setCookie(ctx *gin.Context, key string, value string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(999999999) * time.Second),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(ctx.Writer, cookie)
}

func getTypeId(url string) (t, id string) {
	f := func(u string) string {
		r := strings.Split(u, poolHost)[0]
		return strings.Split(r, protocolHttps)[1]
	}

	r := f(url)
	if r == "" {
		logrus.Warnf("cannot split url")

		return
	}

	if strings.Contains(r, typeInference) ||
		strings.Contains(r, typeEvaluate) {

		v := strings.Split(r, "-")

		return v[0], v[1]
	}

	if strings.Contains(r, typeCloud) {
		return typeCloud, strings.Split(r, typeCloud+"-")[1]
	}

	return
}

func forward(ctx *gin.Context, u string) {
	remote, err := url.Parse(u)
	if err != nil {
		logrus.Warnf("cannot parse url")

		ctx.JSON(http.StatusBadRequest, "")
	}

	cookies := ctx.Request.Cookies()
	for _, v := range cookies {
		fmt.Printf("Name: %s, Value: %s\n", v.Name, v.Value)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = ctx.Param("proxyPath")
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func getSetURL(ctx *gin.Context) (string, error) {
	u := ctx.Query("url")
	if u != "" {
		setCookie(ctx, "url", u)

		return u, nil
	}

	cookieURL, err := ctx.Request.Cookie("url")
	if err != nil {
		logrus.Warnf("cannot found url")

		ctx.JSON(http.StatusBadRequest, "")

		return "", errors.New("get url error")
	}

	return cookieURL.Value, nil
}

func proxy(ctx *gin.Context) {
	// get resource url
	u, err := getSetURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")

		return
	}

	// check auth
	path := ctx.Request.URL.Path
	fmt.Printf("path: %v\n", path)
	if path == "/" {
		t, id := getTypeId(u)
		s := NewXiheServer(ctx)
		ok, err := s.AllowedCloud(fmt.Sprintf("%s%s/api/v1/%s/%s", protocolHttps, hostTest, t, id))
		if err != nil {
			logrus.Warnf("internal error: %s", err.Error())
	
			ctx.JSON(http.StatusInternalServerError, "")
	
			return
		}
		if !ok {
			ctx.JSON(http.StatusUnauthorized, "")
	
			return
		}
	}

	// forward
	forward(ctx, u)
}

// func Request(ctx content.Context, url string) (code string, err error) {
// 	url := "https://example.com/api"
// 	cookieValue := "your-cookie-value" // 替换为你的 Cookie 值

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}

// 	// 添加 Cookie 到请求头
// 	cookie := &http.Cookie{Name: "your-cookie-name", Value: cookieValue}
// 	req.AddCookie(cookie)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// 检查响应状态码
// 	if resp.StatusCode != http.StatusOK {
// 		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
// 		return
// 	}

// }