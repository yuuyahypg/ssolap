package server

import (
    "bytes"
    "github.com/gin-gonic/gin"
    "github.com/bitly/go-simplejson"

    "github.com/yuuyahypg/ssolap/olap/conf"
  	"github.com/yuuyahypg/ssolap/olap/buffer"

    //"io/ioutil"
    //"fmt"
)

func SetRoutes(e *gin.Engine, conf *conf.Conf, buf *buffer.RegisteredBuffer) {
    dimensionFile, err := Asset("config/dimensions.json")
    if err != nil {
        panic(err)
    }

    js, err := simplejson.NewJson(dimensionFile)
    if err != nil {
        panic(err)
    }
    SetApiDimensions(e, js)
    SetApiRequest(e, js, conf, buf)
}

func SetApiDimensions(e *gin.Engine, js *simplejson.Json) {
  dimensions, err := js.Get("dimensions").Array()
  if err != nil {
      panic(err)
  }

  fact, err := js.Get("fact").Map()
  if err != nil {
      panic(err)
  }

  e.GET("/api/dimensions", func(c *gin.Context) {
      c.JSON(200, gin.H{
          "dimensions": dimensions,
          "fact": fact,
      })
  })
}

func SetApiRequest(e *gin.Engine, js *simplejson.Json, conf *conf.Conf, buf *buffer.RegisteredBuffer) {
  fact, err := js.Get("fact").Map()
  if err != nil {
      panic(err)
  }

  dimensions, _ := fact["dimensions"].([]interface{})

  e.GET("/api/request", func(c *gin.Context) {
      query := bytes.Buffer{}
      for _, dimension := range dimensions {
          ds, _ := dimension.(string)
          query.WriteString(c.Query(ds))
          query.WriteString(";")
      }

      ob := buf.GetResult(query.String(), conf)
      result := []map[string]interface{}{}
      for _, tuple := range ob.Buffer {
        result = append(result, tuple)
      }

      c.JSON(200, gin.H{
          "tuples": result,
      })
  })
}
