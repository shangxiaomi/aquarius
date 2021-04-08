package main

/*
(1)
$ curl -i http://localhost:9999/
HTTP/1.1 200 OK
Date: Mon, 12 Aug 2019 16:52:52 GMT
Content-Length: 18
Content-Type: text/html; charset=utf-8
<h1>Hello aquarius</h1>

(2)
$ curl "http://localhost:9999/hello?name=aquariusktutu"
hello aquariusktutu, you're at /hello

(3)
$ curl "http://localhost:9999/login" -X POST -d 'username=aquariusktutu&password=1234'
{"password":"1234","username":"aquariusktutu"}

(4)
$ curl "http://localhost:9999/xxx"
404 NOT FOUND: /xxx
*/

import (
	"aquarius"
	"fmt"
	"net/http"
)

func main() {
	e := aquarius.New()
	r := e.Group("/shp")
	{
		r.GET("/", func(c *aquarius.Context) {
			c.HTML(http.StatusOK, "<h1>Hello aquarius</h1>")
		})
		r.GET("/hello", func(c *aquarius.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})

		r.POST("/login", func(c *aquarius.Context) {
			c.JSON(http.StatusOK, aquarius.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}
	e.GET("/hello", func(c *aquarius.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	group2 := e.Group("/middleware")
	group2.Use(func(c *aquarius.Context) {
		fmt.Println("/middleware begin")
		c.Next()
		fmt.Println("/middleware end")
	})
	group2.GET("/", func(c *aquarius.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	group1 := e.Group("/middleware/shp")
	group1.Use(func(c *aquarius.Context) {
		fmt.Println("/middleware/shp begin")
		c.Next()
		fmt.Println("/middleware/shp end")
	})
	group1.GET("/", func(c *aquarius.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	e.Run(":9999")
}
