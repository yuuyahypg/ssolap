package server

import (
    "strconv"
    "bytes"
    "github.com/gin-gonic/gin"
    "github.com/bitly/go-simplejson"
    "github.com/BurntSushi/toml"

    "github.com/yuuyahypg/ssolap/olap/conf"
  	"github.com/yuuyahypg/ssolap/olap/buffer"

    //"io/ioutil"
    //"fmt"
)

type Config struct {
    Database DbConfig
}

type DbConfig struct {
    Use bool `toml:"use"`
    User string `toml:"user"`
    Name string `toml:"name"`
    Pass string `toml:"pass"`
}

func SetRoutes(e *gin.Engine, conf *conf.Conf, buf *buffer.RegisteredBuffer, topology string) {
    dimensionFile, err := Asset("config/" + topology + "/dimensions.json")
    if err != nil {
        panic(err)
    }

    var config Config
    _, err = toml.DecodeFile("./config/" + topology + "/config.toml", &config)
    if err != nil {
        panic(err)
    }

    js, err := simplejson.NewJson(dimensionFile)
    if err != nil {
        panic(err)
    }
    SetApiDimensions(e, js, config.Database.Use)
    SetApiRequest(e, js, conf, buf)

    if config.Database.Use {
      SetApiGeometry(e, config)
    }
}

func SetApiDimensions(e *gin.Engine, js *simplejson.Json, isDBConnected bool) {
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
          "isDBConnected": isDBConnected,
      })
  })
}

func SetApiRequest(e *gin.Engine, js *simplejson.Json, conf *conf.Conf, buf *buffer.RegisteredBuffer) {
    dimensions := []string{}
    dimensionsJson, err := js.Get("dimensions").Array()
    if err != nil {
        panic(err)
    }

    for _, v := range dimensionsJson {
        if dimensionJson, isMap := v.(map[string]interface{}); isMap != false {
            dimensions = append(dimensions, dimensionJson["name"].(string))
        }
    }

    e.GET("/api/request", func(c *gin.Context) {
        query := bytes.Buffer{}
        for _, dimension := range dimensions {
            query.WriteString(c.Query(dimension))
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

func SetApiGeometry(e *gin.Engine, config Config) {
  db, _ := ConnectDB(config)
  e.GET("/api/geometry", func(c *gin.Context) {
      southWestLon, _ := strconv.ParseFloat(c.Query("southWestLon"), 64)
      southWestLat, _ := strconv.ParseFloat(c.Query("southWestLat"), 64)
      northEastLon, _ := strconv.ParseFloat(c.Query("northEastLon"), 64)
      northEastLat, _ := strconv.ParseFloat(c.Query("northEastLat"), 64)

      var geo *FeatureCollection
      if c.Query("region") == "region2" {
        geo = db.GetBoundedArea(southWestLon, southWestLat, northEastLon, northEastLat)
      } else if c.Query("region") == "city" {
        geo = db.GetBoundedAreaCity(southWestLon, southWestLat, northEastLon, northEastLat)
      } else if c.Query("region") == "prefecture" {
        geo = db.GetBoundedAreaPrefecture(southWestLon, southWestLat, northEastLon, northEastLat)
      }


      c.JSON(200, gin.H{
          "geojson": geo,
      })
  })
}
