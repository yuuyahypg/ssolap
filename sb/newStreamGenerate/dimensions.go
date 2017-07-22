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
    measures map[string]string
    temporalDimension []interface{}
    geoCoder *GeoCoder
}

type DimensionTable struct {
    table map[string]map[string]string
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)
const layoutMinute = "2006/01/02 15:04"
const layoutHour = "2006/01/02 15"
const layoutDay = "2006/01/02"
const layoutMonth = "2006/01"
const layoutYear = "2006/01"

// create star schema from dimension.json
func NewStarSchema(js *simplejson.Json, topology string) (*StarSchema, error) {
    dimensionTypes := map[string]string{}
    dimensionTables := map[string]*DimensionTable{}
    measures := map[string]string{}
    temporalDimension := []interface{}{}
    var geoCoder *GeoCoder

    dimensions, err := js.Get("dimensions").Array()
    if err != nil {
        panic(err)
    }

    // fact tableとdimension tableを結合するために必要な情報を保存する
    for _, v := range dimensions {
        if dimension, isMap := v.(map[string]interface{}); isMap != false {
            switch dimension["type"].(string) {
            case "normal":
                input := dimension["input"].(string)
                dimensionTypes[input] = "normal"
                dimensionTables[input] = createDimensionMap(dimension, topology)
            case "spatial":
                name := dimension["name"].(string)
                dimensionTypes[name] = "spatial"
                geoCoder, _ = ConnectDB(topology)
            case "temporal":
                input := dimension["input"].(string)
                dimensionTypes[input] = "temporal"
                temporalDimension = dimension["outputs"].([]interface{})
            }
        }
    }

    fact, err := js.Get("fact").Map()
    if err != nil {
        panic(err)
    }

    ms, isArray := fact["measures"].([]interface{})
    if isArray == false {
        panic("missing dimensions.json")
    }

    for _, m := range ms {
        if measure, isMap := m.(map[string]interface{}); isMap != false {
            measures[measure["name"].(string)] = measure["type"].(string)
        }
    }

    return &StarSchema{
        dimensionTypes: dimensionTypes,
        dimensionTables: dimensionTables,
        measures: measures,
        temporalDimension: temporalDimension,
        geoCoder: geoCoder,
    }, nil
}

// create dimension tables from csv file
func createDimensionMap(dimension map[string]interface{}, topology string) *DimensionTable {
    dimensionMap := map[string]map[string]string{}
    var inputIndex int

    for i, column := range dimension["csv"].([]interface{}) {
        if column.(string) == dimension["input"].(string) {
            inputIndex = i
        }
    }

    picks := pick(dimension["outputs"].([]interface{}), dimension["csv"].([]interface{}))

    // df, err := os.Open("./config/dimensionLevels/" + dimension["name"].(string) + ".csv")
    df, err := os.Open("./config/" + topology + "/dimensionTables/" + dimension["name"].(string) + ".csv")
    if err != nil {
        fmt.Println("not exist dimension table file")
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

        for column, i := range picks {
            columns[column] = record[i]
        }

        dimensionMap[record[inputIndex]] = columns
    }
    df.Close()

    return &DimensionTable{
        table: dimensionMap,
    }
}

func pick(outputs []interface{}, csv []interface{}) map[string]int {
    picks := map[string]int{}
    for i, column := range csv {
        for _, output := range outputs {
            if column == output {
                picks[output.(string)] = i
                break
            }
        }
    }
    return picks
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
            newTuple["lon"] = tuple.Data["lon"]
            newTuple["lat"] = tuple.Data["lat"]

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
                case "day":
                    newTuple["day"] = data.String(tsJst.Format(layoutDay))
                case "month":
                    newTuple["month"] = data.String(tsJst.Format(layoutMonth))
                case "year":
                    newTuple["year"] = data.String(tsJst.Format(layoutYear))
                }
            }
            newTuple[k] = t
        }
    }

    for k, _ := range s.measures {
        newTuple[k] = tuple.Data[k]
    }

    return &newTuple, nil
}
