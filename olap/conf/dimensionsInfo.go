package conf

import (
    "github.com/bitly/go-simplejson"

    "io/ioutil"
    "fmt"
)

type DimensionsInfo struct {
    DimensionsInfo [][]string
    SumType string
}

func NewDimensionsInfo() *DimensionsInfo {
    dimensionsInfo := [][]string{}
    dimensions := getDimensionsFromJson()
    var sumType string

    for _, v := range dimensions {
        if dimension, isMap := v.(map[string]interface{}); isMap != false {
            if dimType := dimension["type"].(string); dimType != "measure" {
                levels := dimension["levels"].([]interface{})
                dimensionsInfo = append(dimensionsInfo, newDimensionInfo(levels))
            } else {
                measure := []string{"count"}
                sumType = dimension["value"].(string)
                if sumType != "none" {
                    measure = append(measure, dimension["name"].(string))
                }
                dimensionsInfo = append(dimensionsInfo, measure)
            }
        }
    }

    return &DimensionsInfo{
        DimensionsInfo: dimensionsInfo,
        SumType: sumType,
    }
}

func newDimensionInfo(levels []interface{}) []string {
    dimensionInfo := []string{}
    for _, v := range levels {
        dimensionInfo = append(dimensionInfo, v.(string))
    }

    return dimensionInfo
}

func getDimensionsFromJson() []interface{} {
    dimensionFile, err := ioutil.ReadFile("./config/dimension.json")
    if err != nil {
        fmt.Println("not exist dimension.json")
        panic(err)
    }

    js, err := simplejson.NewJson(dimensionFile)
    if err != nil {
        panic(err)
    }

    dimensions, err := js.Get("dimensions").Array()
    if err != nil {
        panic(err)
    }

    return dimensions
}
