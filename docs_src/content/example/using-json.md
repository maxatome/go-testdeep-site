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

func TestCreateRecord(t *testing.T) {
  assert := td.Assert(t)

  before := time.Now().Truncate(time.Second)
  record := td.Must(CreateRecord("Bob", 23))

  assert = assert.RootName("RECORD") // Use RECORD instead of DATA in failure reports
  assert.Cmp(record, td.JSON(`
    {
      "Name":      "Bob",
      "Age":       23,
      "Id":        NotZero(), // comments and operators allowed!
      "CreatedAt": $1
    }`,
    td.Between(before, time.Now())),
    "Newly created record")
}
```

Test it in playground: https://go.dev/play/p/gRR9jGXGQip
