package olap

import (
    "gopkg.in/sensorbee/sensorbee.v0/core"

    "github.com/yuuyahypg/ssolap/olap/conf"
    "github.com/yuuyahypg/ssolap/olap/buffer"
    pb "github.com/yuuyahypg/ssolap/proto"
    "google.golang.org/grpc"

    "flag"
    "net"
    "strconv"
)

var (
	 addrFlag = flag.String("addr", ":5000", "Address host:post")
)

type Server struct {
    Conf *conf.Conf
    Buffer *buffer.RegisteredBuffer
}

func NewServer() *Server {
    config := conf.NewConf()
    buf := buffer.NewRegisteredBuffer(config)

    lis, err := net.Listen("tcp", *addrFlag)
    if err != nil {
		    panic(err)
	  }

    s := grpc.NewServer()
    server := &Server{
        Conf: config,
        Buffer: buf,
    }

    pb.RegisterOlapServer(s, server)
    go s.Serve(lis)

    return server
}

func (s *Server) RecieveTuple(tuple *core.Tuple) {
    s.Buffer.AddTuple(tuple)
}

func (s *Server) GetResult(req *pb.OlapRequest, st pb.Olap_GetResultServer) error {
    buf := s.Buffer.GetResult(req.Request, s.Conf)

    for _, v := range buf.Buffer {
        tuple := map[string]string{}
        for key, value := range v {
            if key == "count" {
                tuple[key] = strconv.Itoa(value.(int))
            } else if key == s.Buffer.Sum {
                if s.Buffer.SumType == "int" {
                    tuple[key] = strconv.Itoa(value.(int))
                } else {
                    tuple[key] = strconv.FormatFloat(value.(float64), 'E', -1, 64)
                }
            } else {
                tuple[key] = value.(string)
            }
        }

        if err := st.Send(&pb.OlapResult{Tuple: tuple}); err != nil {
            return err
        }
    }
    return nil
}

func (s *Server) Close() {
    s.Buffer.DeleteSchedule.Stop()
}
