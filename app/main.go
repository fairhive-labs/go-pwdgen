package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/trendev/go-pwdgen/generator"
)

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		l := 16
		pwd := generator.Generate(l)
		c.JSON(200, gin.H{
			"length":   l,
			"password": pwd,
		})
	})

	log.Fatal(r.Run())
}
