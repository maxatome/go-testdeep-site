+++
title = "Using JSON"
weight = 70
+++

JSON is a first class citizen in go-testdeep world thanks to its
specific operators: [`JSON`]({{% ref "JSON" %}}),
[`SubJSONOf`]({{% ref "SubJSONOf" %}}),
[`SuperJSONOf`]({{% ref "SuperJSONOf" %}}) and
[`JSONPointer`]({{% ref "JSONPointer" %}}).

```go
import (
  "testing"
  "time"

  "github.com/maxatome/go-testdeep/td"
)

func TestCreateRecord(tt *testing.T) {
  t := td.NewT(tt)

  before := time.Now().Truncate(time.Second)
  record, err := CreateRecord("Bob", 23)

  if t.CmpNoError(err) {
    t = t.RootName("RECORD") // Use RECORD instead of DATA in failure reports
    t.Cmp(record, td.JSON(`
{
  "Name":      "Bob",
  "Age":       23,
  "Id":        NotZero(), // comments and operators allowed!
  "CreatedAt": $1
}`,
      td.Between(before, time.Now()),
    ),
      "Newly created record")
  }
}
```

Test it in playground: https://play.golang.org/p/pUC-RMPWyhu
