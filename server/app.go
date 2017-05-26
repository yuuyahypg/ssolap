package server

import (
	"log"
	"os"
	//"fmt"

	"html/template"
	"net/http"
	"strings"

	"gopkg.in/sensorbee/sensorbee.v0/core"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"

	"github.com/yuuyahypg/ssolap/olap/conf"
	"github.com/yuuyahypg/ssolap/olap/buffer"
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

type Server struct {
    Conf *conf.Conf
    Buffer *buffer.RegisteredBuffer
}

func (s *Server) RecieveTuple(tuple *core.Tuple) {
    s.Buffer.AddTuple(tuple)
}

func (s *Server) Close() {
    s.Buffer.DeleteSchedule.Stop()
}

func Run() *Server {
	config := conf.NewConf()
	buf := buffer.NewRegisteredBuffer(config)

  r := gin.Default()
  r.HTMLRender = loadTemplates("react.html")

  r.Use(func(c *gin.Context) {
		id, _ := uuid.NewV4()
		c.Set("uuid", id)
	})

  SetRoutes(r, config, buf)

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
	go r.Run(":" + port)

	server := &Server{
			Conf: config,
			Buffer: buf,
	}

	return server
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
