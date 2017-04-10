package sender

import (
    "gopkg.in/sensorbee/sensorbee.v0/bql"
    "gopkg.in/sensorbee/sensorbee.v0/core"
    "gopkg.in/sensorbee/sensorbee.v0/data"
    "github.com/BurntSushi/toml"

    "encoding/csv"
    "strings"
    "strconv"
    "io"
    "os"
    "os/user"
    "bytes"
    "time"
    "fmt"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// sensorbee source plugin

type Connector struct {
    dir string
}

func (c *Connector) GenerateStream(ctx *core.Context, writer core.Writer) error {
    fmt.Println("send start")
    // file open
    buffer := bytes.NewBufferString(c.dir)
    //buffer.WriteString("test.csv")
    file, err := os.Open(buffer.String())
    if err != nil {
        panic(err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    time.Sleep(10 * time.Second)


    fmt.Println(time.Now())
    // read sensor data
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        id, _ := strconv.Atoi(record[0])
        lon, _ := strconv.ParseFloat(record[4], 64)
        lat, _ := strconv.ParseFloat(record[5], 64)

        tuple := core.NewTuple(data.Map{
            "id":        data.Int(id),
            "age":       data.String(record[7]),
            "gender":    data.String(record[6]),
            "purpose":   data.String(record[10]),
            "transport": data.String(record[13]),
            "work":      data.String(record[9]),
            "lon":       data.Float(lon),
            "lat":       data.Float(lat),
            "timestamp": data.Timestamp(time.Now().In(jst)),
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
type Config struct {
    Source SourceDir
}

type SourceDir struct {
    Path string `toml:"path"`
}

type SourceGetter struct {
}

func (s *SourceGetter) CreateSource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Source, error) {
    var config Config
    usr, _ := user.Current()
    _, err := toml.DecodeFile("./config/senderConfig.toml", &config)
    if err != nil {
          panic(err)
    }

    dir := strings.Replace(config.Source.Path,  "~", usr.HomeDir, 1)

    src := &Connector{
      dir: dir,
    }

    return core.ImplementSourceStop(src), nil
}
