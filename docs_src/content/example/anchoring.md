+++
title = "Using anchoring"
weight = 60
+++

Operators can also directly be anchored in litterals, still using the
[`td.T` type]({{% ref "functions/td-t" %}}), avoiding the use of the
[`Struct`]({{% ref "Struct" %}}) operator:

```go
import (
  "testing"
  "time"

  "github.com/maxatome/go-testdeep/td"
)

func TestCreateRecord(t *testing.T) {
  assert := td.Assert(t)

  before := time.Now().Truncate(time.Second)
  record := assert.Must(CreateRecord("Bob", 23))

  // Use RECORD instead of DATA in failure reports
  assert.RootName("RECORD").
    Cmp(record,
      &Record{
        Name:      "Bob",
        Age:       23,
        Id:        assert.AnchorT[uint64](td.NotZero()),
        CreatedAt: assert.AnchorT[time.Time](td.Between(before, time.Now())),
      },
      "Newly created record")
}
```

Test it in playground: https://go.dev/play/p/tfo4ASgp8y7

See the
[`AnchorT`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.AnchorT)
method documentation for details. Note that
[`AT`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.AT) method
is also a synonym for `AnchorT`:

```go
          Id:        assert.AT[uint64](td.NotZero()),
          CreatedAt: assert.AT[time.Time](td.Between(before, time.Now())),
```
