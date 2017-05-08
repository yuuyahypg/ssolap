package server

import (
	"log"
	"os"

	"html/template"
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir, AssetInfo, root}
	return &binaryFileSystem{
		fs,
	}
}

func Run() {
  r := gin.Default()
  r.HTMLRender = loadTemplates("react.html")

  r.Use(func(c *gin.Context) {
		id, _ := uuid.NewV4()
		c.Set("uuid", id)
	})

  SetRoutes(r)

  react := NewReact(
		"assets/js/bundle.js",
		true,
		r,
	)

	r.GET("/", react.Handle)

  r.Use(static.Serve("/", BinaryFileSystem("assets")))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	r.Run(":" + port)
}

func loadTemplates(list ...string) multitemplate.Render {
	r := multitemplate.New()

	for _, x := range list {
		templateString, err := Asset("server/templates/" + x)
		if err != nil {
			log.Fatal(err)
		}

		tmplMessage, err := template.New(x).Parse(string(templateString))
		if err != nil {
			log.Fatal(err)
		}

		r.Add(x, tmplMessage)
	}

	return r
}
