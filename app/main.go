package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trendev/go-pwdgen/generator"
)

const length = 16

func generate(c *gin.Context) {
	l := c.DefaultQuery("l", strconv.Itoa(length))

	ln, err := strconv.Atoi(l)
	if err != nil {
		log.Println("incorrect length, converted to", length)
		ln = length
	}

	pwd := generator.Generate(ln)

	c.JSON(http.StatusOK, gin.H{
		"length":   len(pwd),
		"password": pwd,
	})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", generate)
	return r
}

func main() {
	r := setupRouter()
	log.Fatal(r.Run())
}
