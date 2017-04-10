package sink

import (
    "gopkg.in/sensorbee/sensorbee.v0/bql"
    "gopkg.in/sensorbee/sensorbee.v0/core"
    "gopkg.in/sensorbee/sensorbee.v0/data"

    //"github.com/yuuyahypg/ssolap/olap"

    "fmt"
    "time"
)

type Sink struct {
    //olapServer *olap.Server
    count int
}

func (s *Sink) Write(ctx *core.Context, t *core.Tuple) error {
    // s.olapServer.RecieveTuple(t)
    s.count = s.count + 1
    if s.count >= 100000 {
        fmt.Println(time.Now())
    }
    return nil
}

func (s *Sink) Close(ctx *core.Context) error {
    //s.olapServer.Close()
    return nil
}

func Create(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Sink, error) {
    //s := olap.NewServer()
    fmt.Println("success create server")
    return &Sink{
        //olapServer: s,
        count: 0,
    }, nil
}
