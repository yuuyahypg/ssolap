package conf

import (
    "strings"
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
    attribute := []string{}
    dimensions := strings.Split(vertex.Dimension, ";")
    for k, v := range dimensions {
        flag := false
        for _, level := range dimInfo.DimensionsInfo[k] {
            if v == level {
                flag = true
            }

            if flag {
                attribute = append(attribute, level)
            }
        }
    }
    return attribute
}
