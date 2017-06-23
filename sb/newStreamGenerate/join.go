package newStreamGenerate

import (
    "gopkg.in/sensorbee/sensorbee.v0/bql/udf"
    "gopkg.in/sensorbee/sensorbee.v0/core"
    "github.com/bitly/go-simplejson"
    "github.com/BurntSushi/toml"

    "io/ioutil"
    "fmt"
)

type Joiner struct {
    schema *StarSchema
}

type TopologyConfig struct {
    Topology Tconf
}

type Tconf struct {
    Name string `toml:"name"`
}

func (j *Joiner) Process(ctx *core.Context, tuple *core.Tuple, w core.Writer) error {
    newData, _ := j.schema.JoinDimensions(tuple)
    newTuple := core.NewTuple(*newData)

    if err := w.Write(ctx, newTuple); err != nil {
      return err
    }
    return nil
}

func (j *Joiner) Terminate(ctx *core.Context) error {
    if j.schema.geoCoder != nil {
        j.schema.geoCoder.Close()
    }
    return nil
}

func CreateJoiner(decl udf.UDSFDeclarer, inputStream string) (udf.UDSF, error) {
    var tc TopologyConfig
    _, err := toml.DecodeFile("./config/topology.toml", &tc)
    if err != nil {
        fmt.Println("not exist: ./config/topology.toml")
        panic(err)
    }

    dimensionFile, err := ioutil.ReadFile("./config/" + tc.Topology.Name + "/dimensions.json")
    if err != nil {
        fmt.Println("invalid path: ./config/" + tc.Topology.Name + "/dimensions.json")
        panic(err)
    }

    js, err := simplejson.NewJson(dimensionFile)
    if err != nil {
        panic(err)
    }

    schema, err := NewStarSchema(js, tc.Topology.Name)
    if err != nil {
        panic(err)
    }

    if err := decl.Input(inputStream, nil); err != nil {
        return nil, err
    }

    return &Joiner{
        schema: schema,
    }, nil
}
