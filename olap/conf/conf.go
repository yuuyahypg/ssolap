package conf

import (
    "github.com/BurntSushi/toml"
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

func NewConf() *Conf {
    vertices := NewVertices()
    dimensionsInfo := NewDimensionsInfo()
    registeredQuery := NewRegisteredQuery(vertices, dimensionsInfo)

    var config Conf
    _, err := toml.DecodeFile("./config/olapConfig.toml", &config)
    if err != nil {
        panic(err)
    }

    config.Ver = vertices
    config.DimInfo = dimensionsInfo
    config.Query = registeredQuery


    return &config
}
