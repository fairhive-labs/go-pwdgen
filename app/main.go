package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trendev/go-pwdgen/generator"
)

const length = 16

func main() {
	r := setupRouter()
	log.Fatal(r.Run())
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/json", generateJSON)
	r.GET("/", generateHTML)
	return r
}

func generateJSON(c *gin.Context) {
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

func generateHTML(c *gin.Context) {
	l := c.DefaultQuery("l", strconv.Itoa(length))

	ln, err := strconv.Atoi(l)
	if err != nil {
		log.Println("incorrect length, converted to", length)
		ln = length
	}

	pwd := generator.Generate(ln)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"length":   len(pwd),
		"password": pwd,
	})
}
