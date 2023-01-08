package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
	"github.com/gin-gonic/gin"
)

//go:embed templates
var tfs embed.FS

const length = 16

func main() {
	r := setupRouter()
	log.Fatal(r.Run())
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	t := template.Must(template.ParseFS(tfs, "templates/*"))
	r.SetHTMLTemplate(t)
	r.GET("/", generate)
	return r
}

func generate(c *gin.Context) {
	l := c.DefaultQuery("l", strconv.Itoa(length))
	mime := c.DefaultQuery("mime", "html")

	ln, err := strconv.Atoi(l)
	if err != nil {
		log.Println("incorrect length, converted to", length)
		ln = length
	}

	pwd := generator.Generate(ln)

	if mime == "json" {
		c.JSON(http.StatusOK, gin.H{
			"length":   len(pwd),
			"password": pwd,
		})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"length":   len(pwd),
		"password": pwd,
	})
}
