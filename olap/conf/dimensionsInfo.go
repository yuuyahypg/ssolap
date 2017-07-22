package conf

import (
    "github.com/bitly/go-simplejson"

    "io/ioutil"
    "fmt"
)

type DimensionsInfo struct {
    DimensionsInfo [][][]string
    Measures []map[string]string
    TemporalInput string
    TemporalUnit string
}

func NewDimensionsInfo(topology string) *DimensionsInfo {
    dimensionsInfo := [][][]string{}
    measures := []map[string]string{}
    dimensions, fact := getDimensionsFromJson(topology)
    var input string
    var unit string

    for _, v := range dimensions {
        if dimension, isMap := v.(map[string]interface{}); isMap != false {
            levels := dimension["rollUp"].([]interface{})
            dimensionsInfo = append(dimensionsInfo, newDimensionInfo(levels))
            if dimType := dimension["type"].(string); dimType == "temporal" {
                input = dimension["input"].(string)
                first := levels[0].([]interface{})
                unit = first[0].(string)
            }
            //if dimType := dimension["type"].(string); dimType != "measure" {
                //levels := dimension["levels"].([]interface{})
                //dimensionsInfo = append(dimensionsInfo, newDimensionInfo(levels))
            //} else {
                //measure := []string{"count"}
                //sumType = dimension["value"].(string)
                //if sumType != "none" {
                    //measure = append(measure, dimension["name"].(string))
                //}
                //dimensionsInfo = append(dimensionsInfo, measure)
            //}
        }
    }

    if mJson, isArray := fact["measures"].([]interface{}); isArray != false {
        for _, measure := range mJson {
            if m, isMap := measure.(map[string]interface{}); isMap != false {
                newMeasure := map[string]string{}
                for k, v := range m {
                    newMeasure[k] = v.(string)
                }
                measures = append(measures, newMeasure)
            }
        }
    }

    return &DimensionsInfo{
        DimensionsInfo: dimensionsInfo,
        Measures: measures,
        TemporalInput: input,
        TemporalUnit: unit,
    }
}

func newDimensionInfo(levels []interface{}) [][]string {
    dimensionInfo := [][]string{}
    for _, level := range levels {
        l := level.([]interface{})
        dim := []string{}
        for _, v := range l {
            dim = append(dim, v.(string))
        }
        dimensionInfo = append(dimensionInfo, dim)
    }

    return dimensionInfo
}

func getDimensionsFromJson(topology string) ([]interface{}, map[string]interface{}) {
    dimensionFile, err := ioutil.ReadFile("./config/" + topology + "/dimensions.json")
    if err != nil {
        fmt.Println("not exist dimensions.json")
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

    fact, err := js.Get("fact").Map()
    if err != nil {
        panic(err)
    }

    return dimensions, fact
}
