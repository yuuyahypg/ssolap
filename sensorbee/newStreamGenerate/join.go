package newStreamGenerate

import (
    "gopkg.in/sensorbee/sensorbee.v0/bql/udf"
    "gopkg.in/sensorbee/sensorbee.v0/core"
    "github.com/bitly/go-simplejson"

    "io/ioutil"
    "fmt"
)

type Joiner struct {
    schema *StarSchema
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
    dimensionFile, err := ioutil.ReadFile("./config/dimension.json")
    if err != nil {
        fmt.Println("not exist dimension.json")
        panic(err)
    }

    js, err := simplejson.NewJson(dimensionFile)
    if err != nil {
        panic(err)
    }

    schema, err := NewStarSchema(js)
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
