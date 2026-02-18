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
  record := td.Must(CreateRecord("Bob", 23))

  // Use RECORD instead of DATA in failure reports
  assert.RootName("RECORD").
    Cmp(record,
      &Record{
        Name:      "Bob",
        Age:       23,
        Id:        assert.Anchor(td.NotZero(), uint64(0)).(uint64),
        CreatedAt: assert.Anchor(td.Between(before, time.Now())).(time.Time),
      },
      "Newly created record")
}
```

Test it in playground: https://go.dev/play/p/2p0gHGuo87O

See the
[`Anchor`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Anchor)
method documentation for details. Note that
[`A`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.A) method
is also a synonym for Anchor.

```go
          Id:        assert.A(td.NotZero(), uint64(0)).(uint64),
          CreatedAt: assert.A(td.Between(before, time.Now())).(time.Time),
```

Thanks to generics, this can also be written as:

```go
          Id:        td.Anchor[uint64](assert, td.NotZero()),
          CreatedAt: td.Anchor[time.Time](assert, td.Between(before, time.Now())),
```
or
```go
          Id:        td.A[uint64](assert, td.NotZero()),
          CreatedAt: td.A[time.Time](assert, td.Between(before, time.Now())),
```
