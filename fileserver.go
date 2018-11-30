package main

import (
	"flag"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	port     string
	username string
	password string
	dir      string
	head     = `<html>
	<form action="/" enctype="multipart/form-data" method="post">
		<input type="file" name="files" multiple="multiple" />
		<input type="submit" value="上传" />
	</form>
`
	errFmt = `<font color="red">%s</font><br />
	<input type="button" value="返回" onclick="history.back()">`
)

func init() {
	flag.StringVar(&port, "port", "80", "http port")
	flag.StringVar(&username, "user", "", "username")
	flag.StringVar(&password, "pass", "", "password")
	flag.StringVar(&dir, "dir", "./", "file dir")
	flag.Parse()
}

func uploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, fmt.Sprintf(errFmt, "Error: "+err.Error()))
	}

	var p = "/"
	u, err := url.ParseRequestURI(c.GetHeader("Referer"))
	if err == nil {
		// remove ../
		p = path.Clean("/"+u.Path) + "/"
	}

	files := form.File["files"]
	for _, file := range files {
		err = c.SaveUploadedFile(file, path.Join(dir, p, path.Clean("/"+file.Filename)))
		if err != nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, fmt.Sprintf(errFmt, file.Filename+" 上传失败<br />Error: "+err.Error()))
			break
		}
	}
	c.Redirect(302, p)
}

func writeHead(c *gin.Context) {
	if strings.HasSuffix(c.Request.URL.Path, "/") {
		c.Writer.WriteString(head)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	var authGroup *gin.RouterGroup
	if username != "" {
		authGroup = r.Group("/", gin.BasicAuth(gin.Accounts{
			username: password,
		}))
	} else {
		authGroup = r.Group("/")
	}

	getGroup := authGroup.Group("/", writeHead)

	getGroup.StaticFS("/", gin.Dir(dir, true))
	authGroup.POST("/", uploadFile)

	err := r.Run(":" + port)
	if err != nil {
		fmt.Println(err)
	}
}
