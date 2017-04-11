package plugin

import (
    "gopkg.in/sensorbee/sensorbee.v0/bql"
    "gopkg.in/sensorbee/sensorbee.v0/bql/udf"

    "github.com/yuuyahypg/ssolap/peopleFlow/sender"
    "github.com/yuuyahypg/ssolap/sb/newStreamGenerate"
    "github.com/yuuyahypg/ssolap/sb/sink"
)

func init() {
    bql.MustRegisterGlobalSourceCreator("sender", &sender.SourceGetter{})
    udf.MustRegisterGlobalUDSFCreator("joinDimension", udf.MustConvertToUDSFCreator(newStreamGenerate.CreateJoiner))
    bql.MustRegisterGlobalSinkCreator("sink", bql.SinkCreatorFunc(sink.Create))
}
