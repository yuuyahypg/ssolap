package buffer

import (
    "github.com/yuuyahypg/ssolap/olap/conf"

    "bytes"
    "strings"
)

type OndemandBuffer struct {
    Buffer map[string]map[string]interface{}
}

func (rb *RegisteredBuffer) GetResult(query string, config *conf.Conf) *OndemandBuffer {
    buf := &OndemandBuffer{Buffer: map[string]map[string]interface{}{}}
    index := config.Ver.ReverseMap[query]
    ref := config.Ver.VertexMap[index].Reference

    rb.mutex.Lock()
    if ref == -1 {
        buf.registeredQuery(rb, index)
    } else {
        buf.ondemandQuery(rb, ref, query, config)
    }
    rb.mutex.Unlock()

    return buf
}

func (buf *OndemandBuffer) registeredQuery(rb *RegisteredBuffer, index int) {
    for i := 0; i < len(rb.RegiBuff[index]); i++ {
        for _, data := range rb.RegiBuff[index][i] {
            dimString := bytes.Buffer{}
            value := map[string]interface{}{}
            if rb.SumType == "none" {
              for j := 0; j < len(data) - 1; j++ {
                  dimString.WriteString(data[j].(string))
                  dimString.WriteString(";")
                  value[rb.RegiQuery.Query[index][j]] = data[j]
              }

              value["count"] = data[len(data) - 1]
              buf.Buffer[dimString.String()] = value
            } else {
                for j := 0; j < len(data) - 2; j++ {
                    dimString.WriteString(data[j].(string))
                    dimString.WriteString(";")
                    value[rb.RegiQuery.Query[index][j]] = data[j]
                }
                value["count"] = data[len(data) - 2]
                value[rb.Sum] = data[len(data) - 1]
                buf.Buffer[dimString.String()] = value
            }
        }
    }
}

func (buf *OndemandBuffer) ondemandQuery(rb *RegisteredBuffer, ref int, queryString string, config *conf.Conf) {
    position, query := parseString(queryString, config.DimInfo, config.Query.Query[ref])
    for i := 0; i < len(rb.RegiBuff[ref]); i++ {
        for _, data := range rb.RegiBuff[ref][i] {
            dimString := bytes.Buffer{}
            value := map[string]interface{}{}
            for k, pos := range position {
                dimString.WriteString(data[pos].(string))
                dimString.WriteString(";")
                value[query[k]] = data[pos]
            }

            if v, ok := buf.Buffer[dimString.String()]; ok {
                if rb.SumType == "int" {
                    v["count"] = v["count"].(int) + data[len(data) - 2].(int)
                    v[rb.Sum] = v[rb.Sum].(int) + data[len(data) - 1].(int)
                } else if rb.SumType == "float" {
                    v["count"] = v["count"].(int) + data[len(data) - 2].(int)
                    v[rb.Sum] = v[rb.Sum].(float64) + data[len(data) - 1].(float64)
                } else {
                    v["count"] = v["count"].(int) + data[len(data) - 1].(int)
                }
            } else {
                if rb.SumType == "none" {
                    value["count"] = data[len(data) - 1]
                    buf.Buffer[dimString.String()] = value
                } else {
                    value["count"] = data[len(data) - 2]
                    value[rb.Sum] = data[len(data) - 1]
                    buf.Buffer[dimString.String()] = value
                }
            }
        }
    }
}

func parseString(query string, dimInfo *conf.DimensionsInfo, rq []string) ([]int, []string) {
    result := []int{}
    attribute := []string{}
    dimensions := strings.Split(query, ";")
    for k, v := range dimensions {
        flag := false
        for _, level := range dimInfo.DimensionsInfo[k] {
            if v == level {
                flag = true
            }

            if flag {
                attribute = append(attribute, level)
            }
        }
    }

    for _, v := range attribute {
        for i := 0; i < len(rq); i++ {
            if v == rq[i] {
                result = append(result, i)
                break
            }
        }
    }

    return result, attribute
}
