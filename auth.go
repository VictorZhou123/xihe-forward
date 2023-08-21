package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth interface {
	AllowedCloud(string) (bool, error)
}

type xiheServer struct {
	Ctx *gin.Context
}

func NewXiheServer(
	ctx *gin.Context,
) Auth {
	return &xiheServer{
		Ctx: ctx,
	}
}

func (s xiheServer) AllowedCloud(u string) (ok bool, err error) {
	forward(s.Ctx, u)

	if s.Ctx.Request.Response.Status == fmt.Sprint(http.StatusOK) {
		return true, nil
	}

	return
}
