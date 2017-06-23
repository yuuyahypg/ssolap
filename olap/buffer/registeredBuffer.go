package buffer

import (
    "gopkg.in/sensorbee/sensorbee.v0/core"
    "gopkg.in/sensorbee/sensorbee.v0/data"
    "github.com/robfig/cron"

    "github.com/yuuyahypg/ssolap/olap/conf"

    "bytes"
    "time"
    "sync"
    "fmt"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type RegisteredBuffer struct {
    RegiBuff map[int][]map[string]map[string]interface{}
    RegiQuery *conf.RegisteredQuery
    Measures *conf.DimensionsInfo
    topTime time.Time
    ioi int
    DeleteSchedule *cron.Cron
    mutex *sync.Mutex
}

func NewRegisteredBuffer(config *conf.Conf) *RegisteredBuffer {
    regiBuff := map[int][]map[string]map[string]interface{}{}
    for k, _ := range config.Query.Query {
        buff := []map[string]map[string]interface{}{}
        regiBuff[k] = buff
    }

    t := time.Now().In(jst)

    c := cron.New()

    rb := &RegisteredBuffer{
        RegiBuff: regiBuff,
        RegiQuery: config.Query,
        Measures: config.DimInfo,
        topTime: time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, jst),
        ioi: config.Olap.Ioi,
        DeleteSchedule: c,
        mutex: new(sync.Mutex),
    }

    rb.DeleteSchedule.AddFunc("@every 1m", func() { rb.deleteOutOfIoi() })
    rb.DeleteSchedule.Start()

    return rb
}

//func NewRegisteredBuffer(config *conf.Conf) *RegisteredBuffer {
    //regiBuff := map[int][]map[string][]interface{}{}
    //for k, _ := range config.Query.Query {
        //buff := []map[string][]interface{}{}
        //regiBuff[k] = buff
    //}

    //var sum string
    //sType := config.DimInfo.SumType
    //if sType == "none" {
        //sum = "none"
    //} else {
        //sum = config.DimInfo.DimensionsInfo[len(config.DimInfo.DimensionsInfo) - 1][1]
    //}

    //t := time.Now().In(jst)

    //c := cron.New()

    //rb := &RegisteredBuffer{
        //RegiBuff: regiBuff,
        //RegiQuery: config.Query,
        //Sum: sum,
        //SumType: sType,
        //topTime: time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, jst),
        //ioi: config.Olap.Ioi,
        //DeleteSchedule: c,
        //mutex: new(sync.Mutex),
    //}

    //fmt.Println("top time is")
    //fmt.Println(time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, jst))

    //rb.DeleteSchedule.AddFunc("@every 1m", func() { rb.deleteOutOfIoi() })
    //rb.DeleteSchedule.Start()

    //return rb
//}

func (rb *RegisteredBuffer) AddTuple(tuple *core.Tuple) {
    for i, query := range rb.RegiQuery.Query {
        queryString := bytes.Buffer{}
        for _, dim := range query {
            v, _ := tuple.Data[dim]
            vs, _ := data.AsString(v)
            queryString.WriteString(vs)
            queryString.WriteString(";")
        }

        qs := queryString.String()
        t, ok := tuple.Data["timestamp"]
        if !ok {
            break
        }
        ts, _ := data.AsTimestamp(t)
        mins := int(ts.Sub(rb.topTime).Minutes()) % 60

        // 1分ごとに新しいindexを追加する
        rb.mutex.Lock()
        if len(rb.RegiBuff[i]) > mins {
            if d, ok := rb.RegiBuff[i][mins][qs]; ok {
                rb.aggregate(d, tuple)
            } else {
                t := rb.newTuple(query, tuple)
                rb.RegiBuff[i][mins][qs] = t
            }
        } else {
            num := mins - len(rb.RegiBuff[i])
            for j := 0; num >= j; j++ {
                rb.RegiBuff[i] = append(rb.RegiBuff[i], map[string]map[string]interface{}{})
            }

            if d, ok := rb.RegiBuff[i][mins][qs]; ok {
                rb.aggregate(d, tuple)
            } else {
                t := rb.newTuple(query, tuple)
                rb.RegiBuff[i][mins][qs] = t
            }
        }
        rb.mutex.Unlock()
    }
}

func (rb *RegisteredBuffer) aggregate(d map[string]interface{}, tuple *core.Tuple) {
    for _, measure := range rb.Measures.Measures {
        switch measure["type"] {
        case "int":
            v, _ := tuple.Data[measure["name"]]
            vi, _ := data.AsInt(v)
            d[measure["name"]] = d[measure["name"]].(int) + int(vi)
        case "float":
            v, _ := tuple.Data[measure["name"]]
            vf, _ := data.AsFloat(v)
            d[measure["name"]] = d[measure["name"]].(float64) + vf
        }
    }
    d["count"] = d["count"].(int) + 1
}

//func (rb *RegisteredBuffer) aggregate(d []interface{}, tuple *core.Tuple) {
    //length := len(d)

    //switch rb.SumType {
    //case "int":
        //d[length - 2] = d[length - 2].(int) + 1

        //v, _ := tuple.Data[rb.Sum]
        //vi, _ := data.AsInt(v)
        //d[length - 1] = d[length - 1].(int) + int(vi)
    //case "float":
        //d[length - 2] = d[length - 2].(int) + 1

        //v, _ := tuple.Data[rb.Sum]
        //vf, _ := data.AsFloat(v)
        //d[length - 1] = d[length - 1].(float64) + vf
    //case "none":
        //d[length - 1] = d[length - 1].(int) + 1
    //}
//}

func (rb *RegisteredBuffer) newTuple(query []string, tuple *core.Tuple) map[string]interface{} {
    t := map[string]interface{}{}
    for _, dim := range query {
        v, _ := tuple.Data[dim]
        vs, _ := data.AsString(v)
        t[dim] = vs
    }

    t["count"] = 1

    for _, measure := range rb.Measures.Measures {
        switch measure["type"] {
        case "int":
            v, _ := tuple.Data[measure["name"]]
            vi, _ := data.AsInt(v)
            t[measure["name"]] = vi
        case "float":
            v, _ := tuple.Data[measure["name"]]
            vf, _ := data.AsFloat(v)
            t[measure["name"]] = vf
        }
    }

    return t
}

//func (rb *RegisteredBuffer) newTuple(query []string, tuple *core.Tuple) []interface{} {
    //t := []interface{}{}
    //for _, dim := range query {
        //v, _ := tuple.Data[dim]
        //vs, _ := data.AsString(v)
        //t = append(t, vs)
    //}

    //t = append(t, 1)

    //switch rb.SumType {
    //case "int":
      //v, _ := tuple.Data[rb.Sum]
        //vi, _ := data.AsInt(v)
        //t = append(t, vi)
    //case "float":
        //v, _ := tuple.Data[rb.Sum]
        //vf, _ := data.AsFloat(v)
        //t = append(t, vf)
    //}

    //return t
//}

// ioiを超えたデータを削除
func (rb *RegisteredBuffer) deleteOutOfIoi() {
    if (int(time.Now().In(jst).Sub(rb.topTime).Minutes()) % 60) >= rb.ioi {
      rb.mutex.Lock()
      fmt.Println("delete")
      for i, _ := range rb.RegiQuery.Query {
          rb.RegiBuff[i] = rb.RegiBuff[i][1:]
      }
      rb.topTime = rb.topTime.Add(time.Duration(1) * time.Minute)
      rb.mutex.Unlock()
    }
}
