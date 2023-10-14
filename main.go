package main

import (
	"io"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"fmt"
)

var printAst = false

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	if printAst {
		ast.Print(os.Stdout, doc)
	}

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func main() {
	r := gin.Default()

	r.Static("static/", "./static")
	r.LoadHTMLGlob("./static/html/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/render", func(c *gin.Context) {
		md, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}

		html := mdToHTML([]byte(gjson.Get(string(md), "data").String()))

		sha := gjson.Get(string(md), "name").String()
		fileName := fmt.Sprintf("./static/html/%s.html", sha)
		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}

		_, err = file.Write(html)
		if err != nil {
			panic(err)
		}

		file.Close()

		c.Redirect(http.StatusFound, fmt.Sprintf("/%x", sha))
	})

	r.GET("/:htmlFileName", func(c *gin.Context) {
		c.HTML(http.StatusOK, fmt.Sprintf("%s.html", c.Params.ByName("htmlFileName")), gin.H{})
	})

	r.Run(":8080")
}
