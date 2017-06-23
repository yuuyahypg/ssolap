package conf

import (
    "strings"
    //"fmt"
)

type RegisteredQuery struct {
    Query map[int][]string
}

func NewRegisteredQuery(vertices *Vertices, dimInfo *DimensionsInfo) *RegisteredQuery {
    query := map[int][]string{}

    for k, v := range vertices.VertexMap {
        if v.Reference == -1 {
            attribute := registQuery(v, dimInfo)
            query[k] = attribute
        }
    }

    return &RegisteredQuery{
        Query: query,
    }
}

func registQuery(vertex *Vertex, dimInfo *DimensionsInfo) []string {
    attributes := []string{}
    dimensions := strings.Split(vertex.Dimension, ";")
    for i, dimension := range dimensions[:(len(dimensions) - 1)] {
        for _, levels := range dimInfo.DimensionsInfo[i] {
            flag := false
            for _, level := range levels {
              if level == dimension {
                  flag = true
              }

              if flag {
                  attributes = append(attributes, level)
              }
            }

            if flag {
              break
            }
        }
    }

    return attributes
}
