package conf

import (
    "github.com/BurntSushi/toml"
    //"fmt"
)

type Conf struct {
    Ver *Vertices
    DimInfo *DimensionsInfo
    Query *RegisteredQuery
    Olap OlapConfig
}

type OlapConfig struct {
    Ioi int `toml:"ioi"`
}

func NewConf(topology string) *Conf {
    vertices := NewVertices(topology)
    dimensionsInfo := NewDimensionsInfo(topology)
    registeredQuery := NewRegisteredQuery(vertices, dimensionsInfo)

    var config Conf
    _, err := toml.DecodeFile("./config/" + topology + "/config.toml", &config)
    if err != nil {
        panic(err)
    }

    config.Ver = vertices
    config.DimInfo = dimensionsInfo
    config.Query = registeredQuery


    return &config
}
