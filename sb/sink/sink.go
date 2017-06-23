package sink

import (
    "gopkg.in/sensorbee/sensorbee.v0/bql"
    "gopkg.in/sensorbee/sensorbee.v0/core"
    "gopkg.in/sensorbee/sensorbee.v0/data"

    "github.com/yuuyahypg/ssolap/server"

    "fmt"
    //"time"
)

type Sink struct {
    server *server.Server
}

func (s *Sink) Write(ctx *core.Context, t *core.Tuple) error {
    s.server.RecieveTuple(t)
    //s.count = s.count + 1
    //if s.count >= 100000 {
    //    fmt.Println(time.Now())
    //}
    return nil
}

func (s *Sink) Close(ctx *core.Context) error {
    s.server.Close()
    return nil
}

func Create(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Sink, error) {
    s := server.Run()
    fmt.Println("success create server")
    return &Sink{
        server: s,
    }, nil
}
