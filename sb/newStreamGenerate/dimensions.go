package newStreamGenerate

import (
    "gopkg.in/sensorbee/sensorbee.v0/core"
    "gopkg.in/sensorbee/sensorbee.v0/data"
    "github.com/bitly/go-simplejson"

    "io"
    "os"
    "fmt"
    "time"
    "encoding/csv"
)

// example StarSchema
// dimensionTypes = {
//     "gender": "normal"
// }
// dimensionTables = {
//     "1": { "gender": "男" },
//     "2": { "gender": "女" },
//     "9": { "gender": "不明" },
// }

type StarSchema struct {
    dimensionTypes map[string]string
    dimensionTables map[string]*DimensionTable
    temporalDimension []interface{}
    aggFunc string
    geoCoder *GeoCoder
}

type DimensionTable struct {
    table map[string]map[string]string
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)
const layoutMinute = "2006/01/02 15:04"
const layoutHour = "2006/01/02 15"

// create star schema from dimension.json
func NewStarSchema(js *simplejson.Json) (*StarSchema, error) {
    dimensionTypes := map[string]string{}
    dimensionTables := map[string]*DimensionTable{}
    temporalDimension := []interface{}{}
    var aggFunc string
    var geoCoder *GeoCoder
    dimensions, err := js.Get("dimensions").Array()
    if err != nil {
        panic(err)
    }

    // fact tableとdimension tableを結合するために必要な情報を保存する
    for _, v := range dimensions {
        if dimension, isMap := v.(map[string]interface{}); isMap != false {
            name := dimension["name"].(string)

            switch dimension["type"].(string) {
            case "normal":
                dimensionTypes[name] = "normal"
                dimensionTables[name] = createDimensionMap(dimension)
            case "spatial":
                dimensionTypes[name] = "spatial"
                geoCoder, _ = ConnectDB()
            case "temporal":
                dimensionTypes[name] = "temporal"
                temporalDimension = dimension["levels"].([]interface{})
            case "measure":
                dimensionTypes[name] = "measure"
                aggFunc = dimension["func"].(string)
            }
        }
    }

    return &StarSchema{
        dimensionTypes: dimensionTypes,
        dimensionTables: dimensionTables,
        temporalDimension: temporalDimension,
        aggFunc: aggFunc,
        geoCoder: geoCoder,
    }, nil
}

// create dimension tables from csv file
func createDimensionMap(dimension map[string]interface{}) *DimensionTable {
    dimensionMap := map[string]map[string]string{}
    levels := dimension["levels"].([]interface{})
    length := len(levels)

    // df, err := os.Open("./config/dimensionLevels/" + dimension["name"].(string) + ".csv")
    df, err := os.Open("./config/dimensionTables/" + dimension["name"].(string) + ".csv")
    if err != nil {
        fmt.Println("not exist dimension level file")
        panic(err)
    }

    reader := csv.NewReader(df)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        }

        columns := map[string]string{}

        for i := 0; i < length; i++ {
            //columns[levels[i].(string)] = record[i]
            columns[levels[0].(string)] = record[1]
        }

        dimensionMap[record[0]] = columns
    }
    df.Close()
    
    return &DimensionTable{
        table: dimensionMap,
    }
}

func (s *StarSchema) JoinDimensions(tuple *core.Tuple) (*data.Map, error) {
    newTuple := data.Map{}

    for k, v := range s.dimensionTypes {
        switch v {
        case "normal":
            i, _ := tuple.Data[k]
            id, _ := data.AsString(i)
            for columnName, value := range s.dimensionTables[k].table[id] {
                newTuple[columnName] = data.String(value)
            }
        case "spatial":
            id, _ := tuple.Data["id"]
            lon, _ := tuple.Data["lon"]
            lat, _ := tuple.Data["lat"]
            idi, _ := data.AsInt(id)
            lonF, _ := data.AsFloat(lon)
            latF, _ := data.AsFloat(lat)

            region2, region1, city, prefecture, err := s.geoCoder.GeoCoding(int(idi), lonF, latF)

            if err != nil {
                fmt.Println(err)
            }

            newTuple["region2"] = data.String(region2)
            newTuple["region1"] = data.String(region1)
            newTuple["city"] = data.String(city)
            newTuple["prefecture"] = data.String(prefecture)

        case "temporal":
            t, _ := tuple.Data[k]
            ts, _ := data.AsTimestamp(t)
            tsJst := ts.In(jst)
            for _, level := range s.temporalDimension {
                level_s, _ := level.(string)
                switch level_s {
                case "minute":
                    newTuple["minute"] = data.String(tsJst.Format(layoutMinute))
                case "hour":
                    newTuple["hour"] = data.String(tsJst.Format(layoutHour))
                }
            }
            newTuple[k] = t
        case "measure":
            if s.aggFunc != "count" {
                newTuple[k] = tuple.Data[k]
            }
        }
    }

    return &newTuple, nil
}
