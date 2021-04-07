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
	"net/http"
)

func main() {
	r := aquarius.New()
	r.GET("/", func(c *aquarius.Context) {
		c.HTML(http.StatusOK, "<h1>Hello aquarius</h1>")
	})
	r.GET("/hello", func(c *aquarius.Context) {
		// expect /hello?name=aquariusktutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *aquarius.Context) {
		c.JSON(http.StatusOK, aquarius.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
