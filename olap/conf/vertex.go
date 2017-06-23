package conf

import (
    "bytes"
    "strconv"
    "io"
    "os"
    "encoding/csv"
    "fmt"
)

type Vertices struct {
    VertexMap []*Vertex
    ReverseMap map[string]int
}

func NewVertices(topology string) *Vertices {
    vertexMap := []*Vertex{}
    reverseMap := map[string]int{}

    //csvファイルから次元の組み合わせの情報を取得
    df, err := os.Open("./config/" + topology + "/verticesConfig.csv")
    if err != nil {
        fmt.Println("not exist verticesConfig.csv")
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

        index, _ := strconv.Atoi(record[0])

        vertex := newVertex(record)
        vertexMap = append(vertexMap, vertex)
        reverseMap[vertex.Dimension] = index
    }

    df.Close()
    return &Vertices{
        VertexMap: vertexMap,
        ReverseMap: reverseMap,
    }
}

type Vertex struct {
    Reference int
    Dimension string
}

func newVertex(record []string) *Vertex {
    ref, _ := strconv.Atoi(record[1])

    buffer := bytes.Buffer{}
    length := len(record)
    for i := 2; i < length; i++ {
        buffer.WriteString(record[i])
        buffer.WriteString(";")
    }

    return &Vertex{
        Reference: ref,
        Dimension: buffer.String(),
    }
}
