package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trendev/go-pwdgen/generator"
)

func generate(c *gin.Context) {
	l := 16
	pwd := generator.Generate(l)
	c.JSON(http.StatusOK, gin.H{
		"length":   l,
		"password": pwd,
	})
}

func main() {

	r := gin.Default()
	r.GET("/", generate)
	log.Fatal(r.Run())
}
