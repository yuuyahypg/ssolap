package sender

import (
    "gopkg.in/sensorbee/sensorbee.v0/bql"
    "gopkg.in/sensorbee/sensorbee.v0/core"
    "gopkg.in/sensorbee/sensorbee.v0/data"

    "github.com/bitly/go-simplejson"
    "encoding/json"
    "io/ioutil"
    "time"
    "fmt"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// sensorbee source plugin

type Connector struct {
}

func (c *Connector) GenerateStream(ctx *core.Context, writer core.Writer) error {
    // file open
    rf, err := ioutil.ReadFile("/Users/yuuya/data/logs/20170509/1479253956_exp.json")
    if err != nil {
        panic(err)
    }

    js, err := simplejson.NewJson(rf)

    arr, err := js.Array()
    if err != nil {
        panic(err)
    }

    time.Sleep(10 * time.Second)

    fmt.Println("send start")
    //fmt.Println(time.Now())
    // read sensor data

    for _, t := range arr {
        d, _ := t.(map[string]interface{})
        lo, _ := d["lo"].(json.Number)
        la, _ := d["la"].(json.Number)
        tm := d["tm"].(json.Number)
        id := 1
        lon, _ := lo.Float64()
        lat, _ := la.Float64()
        unix, _ := tm.Int64()

        tuple := core.NewTuple(data.Map{
            "id":        data.Int(id),
            "lon":       data.Float(lon),
            "lat":       data.Float(lat),
            "timestamp": data.Timestamp(time.Unix(unix, 0).In(jst)),
            //"timestamp": data.Timestamp(time.Now().In(jst)),
        })

        if err := writer.Write(ctx, tuple); err != nil {
          return err
        }
    }

    return nil
}

func (c *Connector) Stop(ctx *core.Context) error {
    // This method will be implemented by utility functions.
    return nil
}

///////////////////////////////////////////////////////////////////////////////

type SourceGetter struct {
}

func (s *SourceGetter) CreateSource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Source, error) {
    src := &Connector{}

    return core.ImplementSourceStop(src), nil
}
