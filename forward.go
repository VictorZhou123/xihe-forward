package main

import (
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
	protocol = "https://"
	host     = "xihe.mindspore.cn"
	hostTest = "xihe2.test.osinfra.cn"
	poolHost = ".pool1.mindspore.cn"

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
		return strings.Split(r, protocol)[1]
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

func getSetURL(ctx *gin.Context) string {
	u := ctx.Query("url")
	if u != "" {
		setCookie(ctx, "url", u)

		return u
	}

	cookieURL, err := ctx.Request.Cookie("url")
	if err != nil {
		logrus.Warnf("cannot found url")

		ctx.JSON(http.StatusBadRequest, "")
	}

	return cookieURL.Value
}

func proxy(ctx *gin.Context) {
	// get resource url
	u := getSetURL(ctx)

	// check auth
	t, id := getTypeId(u)
	s := NewXiheServer(ctx)
	ok, err := s.AllowedCloud(fmt.Sprintf("%s%s/api/v1/%s/%s", protocol, hostTest, t, id))
	if err != nil {
		logrus.Warnf("internal error: %s", err.Error())

		ctx.JSON(http.StatusInternalServerError, "")

		return
	}
	if !ok {
		ctx.JSON(http.StatusUnauthorized, "")

		return
	}

	// forward
	forward(ctx, u)
}
