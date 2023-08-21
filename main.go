package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	StartWebServer("8000")
}

func StartWebServer(port string) {
	r := gin.Default()

	r.Any("/*proxyPath", proxy)

	r.Run(fmt.Sprintf(":%s", port))
}
