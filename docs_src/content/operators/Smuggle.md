---
title: "Smuggle"
weight: 10
---

```go
func Smuggle(fn, expectedValue any) TestDeep
```

Smuggle operator allows to change data contents or mutate it into
another type before stepping down in favor of generic comparison
process. Of course it is a [smuggler operator]({{% ref "operators#smuggler-operators" %}}). So *fn* is a function
that must take one parameter whose type must be convertible to the
type of the compared value.

As convenient shortcuts, *fn* can be a `string` specifying a
fields-path through structs, maps & slices, or `any` other type, in
this case a simple cast is done (see below for details).

*fn* must return at least one value. These value will be compared as is
to *expectedValue*, here integer 28:

```go
td.Cmp(t, "0028",
  td.Smuggle(func(value string) int {
    num, _ := strconv.Atoi(value)
    return num
  }, 28),
)
```

or using an other [TestDeep operator]({{% ref "operators" %}}), here [`Between`]({{% ref "Between" %}})(28, 30):

```go
td.Cmp(t, "0029",
  td.Smuggle(func(value string) int {
    num, _ := strconv.Atoi(value)
    return num
  }, td.Between(28, 30)),
)
```

*fn* can return a second boolean value, used to tell that a problem
occurred and so stop the comparison:

```go
td.Cmp(t, "0029",
  td.Smuggle(func(value string) (int, bool) {
    num, err := strconv.Atoi(value)
    return num, err == nil
  }, td.Between(28, 30)),
)
```

*fn* can return a third `string` value which is used to describe the
test when a problem occurred (false second boolean value):

```go
td.Cmp(t, "0029",
  td.Smuggle(func(value string) (int, bool, string) {
    num, err := strconv.Atoi(value)
    if err != nil {
      return 0, false, "string must contain a number"
    }
    return num, true, ""
  }, td.Between(28, 30)),
)
```

Instead of returning (X, `bool`) or (X, `bool`, `string`), *fn* can
return (X, [`error`](https://pkg.go.dev/builtin#error)). When a problem occurs, the returned [`error`](https://pkg.go.dev/builtin#error) is
non-`nil`, as in:

```go
td.Cmp(t, "0029",
  td.Smuggle(func(value string) (int, error) {
    num, err := strconv.Atoi(value)
    return num, err
  }, td.Between(28, 30)),
)
```

Which can be simplified to:

```go
td.Cmp(t, "0029", td.Smuggle(strconv.Atoi, td.Between(28, 30)))
```

Imagine you want to compare that the Year of a date is between 2010
and 2020:

```go
td.Cmp(t, time.Date(2015, time.May, 1, 1, 2, 3, 0, time.UTC),
  td.Smuggle(func(date time.Time) int { return date.Year() },
    td.Between(2010, 2020)),
)
```

In this case the data location forwarded to next test will be
something like "DATA.MyTimeField<smuggled>", but you can act on it
too by returning a [`SmuggledGot`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SmuggledGot) struct (by value or by address):

```go
td.Cmp(t, time.Date(2015, time.May, 1, 1, 2, 3, 0, time.UTC),
  td.Smuggle(func(date time.Time) SmuggledGot {
    return SmuggledGot{
      Name: "Year",
      Got:  date.Year(),
    }
  }, td.Between(2010, 2020)),
)
```

then the data location forwarded to next test will be something like
"DATA.MyTimeField.Year". The "." between the current path (here
"DATA.MyTimeField") and the returned Name "Year" is automatically
added when Name starts with a Letter.

Note that [`SmuggledGot`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SmuggledGot) and [`*SmuggledGot`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SmuggledGot) returns are treated
equally, and they are only used when *fn* has only one returned value
or when the second boolean returned value is true.

Of course, all cases can go together:

```go
// Accepts a "YYYY/mm/DD HH:MM:SS" string to produce a time.Time and tests
// whether this date is contained between 2 hours before now and now.
td.Cmp(t, "2020-01-25 12:13:14",
  td.Smuggle(func(date string) (*SmuggledGot, bool, string) {
    date, err := time.Parse("2006/01/02 15:04:05", date)
    if err != nil {
      return nil, false, `date must conform to "YYYY/mm/DD HH:MM:SS" format`
    }
    return &SmuggledGot{
      Name: "Date",
      Got:  date,
    }, true, ""
  }, td.Between(time.Now().Add(-2*time.Hour), time.Now())),
)
```

or:

```go
// Accepts a "YYYY/mm/DD HH:MM:SS" string to produce a time.Time and tests
// whether this date is contained between 2 hours before now and now.
td.Cmp(t, "2020-01-25 12:13:14",
  td.Smuggle(func(date string) (*SmuggledGot, error) {
    date, err := time.Parse("2006/01/02 15:04:05", date)
    if err != nil {
      return nil, err
    }
    return &SmuggledGot{
      Name: "Date",
      Got:  date,
    }, nil
  }, td.Between(time.Now().Add(-2*time.Hour), time.Now())),
)
```

Smuggle can also be used to access a struct field embedded in
several struct layers.

```go
type A struct{ Num int }
// func (a *A) String() string { return fmt.Sprintf("Num is %d", a.Num) }
type B struct{ As map[string]*A }
type C struct{ B B }
got := C{B: B{As: map[string]*A{"foo": {Num: 12}}}}

// Tests that got.B.A.Num is 12
td.Cmp(t, got,
  td.Smuggle(func(c C) int {
    return c.B.As["foo"].Num
  }, 12))
```

As brought up above, a fields-path can be passed as *fn* value
instead of a function pointer. Using this feature, the [`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp)
call in the above example can be rewritten as follows:

```go
// Tests that got.B.As["foo"].Num is 12
td.Cmp(t, got, td.Smuggle("B.As[foo].Num", 12))
```

In addition, simple public methods can also be called like in:

```go
td.Cmp(t, got, td.Smuggle("B.As[foo].String()", "Num is 12"))
```

Allowed methods must not take `any` parameter and must return one
value or a value and an [`error`](https://pkg.go.dev/builtin#error). For the latter case, if the method
returns a non-`nil` [`error`](https://pkg.go.dev/builtin#error), the comparison fails. The comparison also
fails if a panic occurs or if a method cannot be called. No private
fields should be traversed before calling the method. For fun,
consider a more complex example involving [`reflect`](https://pkg.go.dev/reflect) and chaining
method calls:

```go
got := reflect.Valueof(&C{B: B{As: map[string]*A{"foo": {Num: 12}}}})
td.Cmp(t, got, td.Smuggle("Elem().Interface().B.As[foo].String()", "Num is 12"))
```

Contrary to [`JSONPointer`]({{% ref "JSONPointer" %}}) operator, private fields can be followed
and public methods on public fields can be called. Arrays, slices
and maps work using the index/key inside square brackets (e.g. [12]
or [foo]). Maps work only for simple key types (`string` or numbers),
without "" when using strings (e.g. [foo]).

Behind the scenes, a temporary function is automatically created to
achieve the same goal, but adds some checks against `nil` values and
auto-dereferences interfaces and pointers, even on several levels,
like in:

```go
type A struct{ N any }
num := 12
pnum := &num
td.Cmp(t, A{N: &pnum}, td.Smuggle("N", 12))
```

Last but not least, a simple type can be passed as *fn* to operate
a cast, handling specifically strings and slices of bytes:

```go
td.Cmp(t, `{"foo":1}`, td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":1}`)))
// or equally
td.Cmp(t, `{"foo":1}`, td.Smuggle(json.RawMessage(nil), td.JSON(`{"foo":1}`)))
```

converts on the fly a `string` to a [json.RawMessage](https://pkg.go.dev/encoding/json#RawMessage) so [`JSON`]({{% ref "JSON" %}}) operator
can parse it as JSON. This is mostly a shortcut for:

```go
td.Cmp(t, `{"foo":1}`, td.Smuggle(
  func(r json.RawMessage) json.RawMessage { return r },
  td.JSON(`{"foo":1}`)))
```

except that for strings and slices of bytes (like here), it accepts
[`io.Reader`](https://pkg.go.dev/io#Reader) interface too:

```go
var body io.Reader
// …
td.Cmp(t, body, td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":1}`)))
// or equally
td.Cmp(t, body, td.Smuggle(json.RawMessage(nil), td.JSON(`{"foo":1}`)))
```

This last example allows to easily inject body content into JSON
operator.

The difference between Smuggle and [`Code`]({{% ref "Code" %}}) operators is that [`Code`]({{% ref "Code" %}})
is used to do a final comparison while Smuggle transforms the data
and then steps down in favor of generic comparison
process. Moreover, the type accepted as input for the function is
more lax to facilitate the writing of tests (e.g. the function can
accept a `float64` and the got value be an `int`). See examples. On the
other hand, the output type is strict and must match exactly the
expected value type. The fields-path `string` *fn* shortcut and the
cast feature are not available with [`Code`]({{% ref "Code" %}}) operator.

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of only parameter of
*fn*. For the case where *fn* is a fields-path, it is always
`any`, as the type can not be known in advance.

> See also [`Code`]({{% ref "Code" %}}), [`JSONPointer`]({{% ref "JSONPointer" %}}) and [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten).



> See also [<i class='fas fa-book'></i> Smuggle godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Smuggle).

### Examples

{{%expand "Convert example" %}}```go
	t := &testing.T{}

	got := int64(123)

	ok := td.Cmp(t, got,
		td.Smuggle(func(n int64) int { return int(n) }, 123),
		"checks int64 got against an int value")
	fmt.Println(ok)

	ok = td.Cmp(t, "123",
		td.Smuggle(
			func(numStr string) (int, bool) {
				n, err := strconv.Atoi(numStr)
				return n, err == nil
			},
			td.Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = td.Cmp(t, "123",
		td.Smuggle(
			func(numStr string) (int, bool, string) {
				n, err := strconv.Atoi(numStr)
				if err != nil {
					return 0, false, "string must contain a number"
				}
				return n, true, ""
			},
			td.Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = td.Cmp(t, "123",
		td.Smuggle(
			func(numStr string) (int, error) { //nolint: gocritic
				return strconv.Atoi(numStr)
			},
			td.Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Short version :)
	ok = td.Cmp(t, "123",
		td.Smuggle(strconv.Atoi, td.Between(120, 130)),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Lax example" %}}```go
	t := &testing.T{}

	// got is an int16 and Smuggle func input is an int64: it is OK
	got := int(123)

	ok := td.Cmp(t, got,
		td.Smuggle(func(n int64) uint32 { return uint32(n) }, uint32(123)))
	fmt.Println("got int16(123) → smuggle via int64 → uint32(123):", ok)

	// Output:
	// got int16(123) → smuggle via int64 → uint32(123): true

```{{% /expand%}}
{{%expand "Auto_unmarshal example" %}}```go
	t := &testing.T{}

	// Automatically json.Unmarshal to compare
	got := []byte(`{"a":1,"b":2}`)

	ok := td.Cmp(t, got,
		td.Smuggle(
			func(b json.RawMessage) (r map[string]int, err error) {
				err = json.Unmarshal(b, &r)
				return
			},
			map[string]int{
				"a": 1,
				"b": 2,
			}))
	fmt.Println("JSON contents is OK:", ok)

	// Output:
	// JSON contents is OK: true

```{{% /expand%}}
{{%expand "Cast example" %}}```go
	t := &testing.T{}

	// A string containing JSON
	got := `{ "foo": 123 }`

	// Automatically cast a string to a json.RawMessage so td.JSON can operate
	ok := td.Cmp(t, got,
		td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":123}`)))
	fmt.Println("JSON contents in string is OK:", ok)

	// Automatically read from io.Reader to a json.RawMessage
	ok = td.Cmp(t, bytes.NewReader([]byte(got)),
		td.Smuggle(json.RawMessage{}, td.JSON(`{"foo":123}`)))
	fmt.Println("JSON contents just read is OK:", ok)

	// Output:
	// JSON contents in string is OK: true
	// JSON contents just read is OK: true

```{{% /expand%}}
{{%expand "Complex example" %}}```go
	t := &testing.T{}

	// No end date but a start date and a duration
	type StartDuration struct {
		StartDate time.Time
		Duration  time.Duration
	}

	// Checks that end date is between 17th and 19th February both at 0h
	// for each of these durations in hours

	for _, duration := range []time.Duration{48 * time.Hour, 72 * time.Hour, 96 * time.Hour} {
		got := StartDuration{
			StartDate: time.Date(2018, time.February, 14, 12, 13, 14, 0, time.UTC),
			Duration:  duration,
		}

		// Simplest way, but in case of Between() failure, error will be bound
		// to DATA<smuggled>, not very clear...
		ok := td.Cmp(t, got,
			td.Smuggle(
				func(sd StartDuration) time.Time {
					return sd.StartDate.Add(sd.Duration)
				},
				td.Between(
					time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
					time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC))))
		fmt.Println(ok)

		// Name the computed value "ComputedEndDate" to render a Between() failure
		// more understandable, so error will be bound to DATA.ComputedEndDate
		ok = td.Cmp(t, got,
			td.Smuggle(
				func(sd StartDuration) td.SmuggledGot {
					return td.SmuggledGot{
						Name: "ComputedEndDate",
						Got:  sd.StartDate.Add(sd.Duration),
					}
				},
				td.Between(
					time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
					time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC))))
		fmt.Println(ok)
	}

	// Output:
	// false
	// false
	// true
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Interface example" %}}```go
	t := &testing.T{}

	gotTime, err := time.Parse(time.RFC3339, "2018-05-23T12:13:14Z")
	if err != nil {
		t.Fatal(err)
	}

	// Do not check the struct itself, but its stringified form
	ok := td.Cmp(t, gotTime,
		td.Smuggle(func(s fmt.Stringer) string {
			return s.String()
		},
			"2018-05-23 12:13:14 +0000 UTC"))
	fmt.Println("stringified time.Time OK:", ok)

	// If got does not implement the fmt.Stringer interface, it fails
	// without calling the Smuggle func
	type MyTime time.Time
	ok = td.Cmp(t, MyTime(gotTime),
		td.Smuggle(func(s fmt.Stringer) string {
			fmt.Println("Smuggle func called!")
			return s.String()
		},
			"2018-05-23 12:13:14 +0000 UTC"))
	fmt.Println("stringified MyTime OK:", ok)

	// Output:
	// stringified time.Time OK: true
	// stringified MyTime OK: false

```{{% /expand%}}
{{%expand "Field_path example" %}}```go
	t := &testing.T{}

	type Body struct {
		Name  string
		Value any
	}
	type Request struct {
		Body *Body
	}
	type Transaction struct {
		Request
	}
	type ValueNum struct {
		Num int
	}

	got := &Transaction{
		Request: Request{
			Body: &Body{
				Name:  "test",
				Value: &ValueNum{Num: 123},
			},
		},
	}

	// Want to check whether Num is between 100 and 200?
	ok := td.Cmp(t, got,
		td.Smuggle(
			func(tr *Transaction) (int, error) {
				if tr.Body == nil ||
					tr.Body.Value == nil {
					return 0, errors.New("Request.Body or Request.Body.Value is nil")
				}
				if v, ok := tr.Body.Value.(*ValueNum); ok && v != nil {
					return v.Num, nil
				}
				return 0, errors.New("Request.Body.Value isn't *ValueNum or nil")
			},
			td.Between(100, 200)))
	fmt.Println("check Num by hand:", ok)

	// Same, but automagically generated...
	ok = td.Cmp(t, got, td.Smuggle("Request.Body.Value.Num", td.Between(100, 200)))
	fmt.Println("check Num using a fields-path:", ok)

	// And as Request is an anonymous field, can be simplified further
	// as it can be omitted
	ok = td.Cmp(t, got, td.Smuggle("Body.Value.Num", td.Between(100, 200)))
	fmt.Println("check Num using an other fields-path:", ok)

	// Note that maps and array/slices are supported
	got.Body.Value = map[string]any{
		"foo": []any{
			3: map[int]string{666: "bar"},
		},
	}
	ok = td.Cmp(t, got, td.Smuggle("Body.Value[foo][3][666]", "bar"))
	fmt.Println("check fields-path including maps/slices:", ok)

	// Output:
	// check Num by hand: true
	// check Num using a fields-path: true
	// check Num using an other fields-path: true
	// check fields-path including maps/slices: true

```{{% /expand%}}
## CmpSmuggle shortcut

```go
func CmpSmuggle(t TestingT, got, fn , expectedValue any, args ...any) bool
```

CmpSmuggle is a shortcut for:

```go
td.Cmp(t, got, td.Smuggle(fn, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

If *t* is a [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) then its Config field is inherited.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSmuggle godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSmuggle).

### Examples

{{%expand "Convert example" %}}```go
	t := &testing.T{}

	got := int64(123)

	ok := td.CmpSmuggle(t, got, func(n int64) int { return int(n) }, 123,
		"checks int64 got against an int value")
	fmt.Println(ok)

	ok = td.CmpSmuggle(t, "123", func(numStr string) (int, bool) {
		n, err := strconv.Atoi(numStr)
		return n, err == nil
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = td.CmpSmuggle(t, "123", func(numStr string) (int, bool, string) {
		n, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, false, "string must contain a number"
		}
		return n, true, ""
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = td.CmpSmuggle(t, "123", func(numStr string) (int, error) { //nolint: gocritic
		return strconv.Atoi(numStr)
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Short version :)
	ok = td.CmpSmuggle(t, "123", strconv.Atoi, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Lax example" %}}```go
	t := &testing.T{}

	// got is an int16 and Smuggle func input is an int64: it is OK
	got := int(123)

	ok := td.CmpSmuggle(t, got, func(n int64) uint32 { return uint32(n) }, uint32(123))
	fmt.Println("got int16(123) → smuggle via int64 → uint32(123):", ok)

	// Output:
	// got int16(123) → smuggle via int64 → uint32(123): true

```{{% /expand%}}
{{%expand "Auto_unmarshal example" %}}```go
	t := &testing.T{}

	// Automatically json.Unmarshal to compare
	got := []byte(`{"a":1,"b":2}`)

	ok := td.CmpSmuggle(t, got, func(b json.RawMessage) (r map[string]int, err error) {
		err = json.Unmarshal(b, &r)
		return
	}, map[string]int{
		"a": 1,
		"b": 2,
	})
	fmt.Println("JSON contents is OK:", ok)

	// Output:
	// JSON contents is OK: true

```{{% /expand%}}
{{%expand "Cast example" %}}```go
	t := &testing.T{}

	// A string containing JSON
	got := `{ "foo": 123 }`

	// Automatically cast a string to a json.RawMessage so td.JSON can operate
	ok := td.CmpSmuggle(t, got, json.RawMessage{}, td.JSON(`{"foo":123}`))
	fmt.Println("JSON contents in string is OK:", ok)

	// Automatically read from io.Reader to a json.RawMessage
	ok = td.CmpSmuggle(t, bytes.NewReader([]byte(got)), json.RawMessage{}, td.JSON(`{"foo":123}`))
	fmt.Println("JSON contents just read is OK:", ok)

	// Output:
	// JSON contents in string is OK: true
	// JSON contents just read is OK: true

```{{% /expand%}}
{{%expand "Complex example" %}}```go
	t := &testing.T{}

	// No end date but a start date and a duration
	type StartDuration struct {
		StartDate time.Time
		Duration  time.Duration
	}

	// Checks that end date is between 17th and 19th February both at 0h
	// for each of these durations in hours

	for _, duration := range []time.Duration{48 * time.Hour, 72 * time.Hour, 96 * time.Hour} {
		got := StartDuration{
			StartDate: time.Date(2018, time.February, 14, 12, 13, 14, 0, time.UTC),
			Duration:  duration,
		}

		// Simplest way, but in case of Between() failure, error will be bound
		// to DATA<smuggled>, not very clear...
		ok := td.CmpSmuggle(t, got, func(sd StartDuration) time.Time {
			return sd.StartDate.Add(sd.Duration)
		}, td.Between(
			time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC)))
		fmt.Println(ok)

		// Name the computed value "ComputedEndDate" to render a Between() failure
		// more understandable, so error will be bound to DATA.ComputedEndDate
		ok = td.CmpSmuggle(t, got, func(sd StartDuration) td.SmuggledGot {
			return td.SmuggledGot{
				Name: "ComputedEndDate",
				Got:  sd.StartDate.Add(sd.Duration),
			}
		}, td.Between(
			time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC)))
		fmt.Println(ok)
	}

	// Output:
	// false
	// false
	// true
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Interface example" %}}```go
	t := &testing.T{}

	gotTime, err := time.Parse(time.RFC3339, "2018-05-23T12:13:14Z")
	if err != nil {
		t.Fatal(err)
	}

	// Do not check the struct itself, but its stringified form
	ok := td.CmpSmuggle(t, gotTime, func(s fmt.Stringer) string {
		return s.String()
	}, "2018-05-23 12:13:14 +0000 UTC")
	fmt.Println("stringified time.Time OK:", ok)

	// If got does not implement the fmt.Stringer interface, it fails
	// without calling the Smuggle func
	type MyTime time.Time
	ok = td.CmpSmuggle(t, MyTime(gotTime), func(s fmt.Stringer) string {
		fmt.Println("Smuggle func called!")
		return s.String()
	}, "2018-05-23 12:13:14 +0000 UTC")
	fmt.Println("stringified MyTime OK:", ok)

	// Output:
	// stringified time.Time OK: true
	// stringified MyTime OK: false

```{{% /expand%}}
{{%expand "Field_path example" %}}```go
	t := &testing.T{}

	type Body struct {
		Name  string
		Value any
	}
	type Request struct {
		Body *Body
	}
	type Transaction struct {
		Request
	}
	type ValueNum struct {
		Num int
	}

	got := &Transaction{
		Request: Request{
			Body: &Body{
				Name:  "test",
				Value: &ValueNum{Num: 123},
			},
		},
	}

	// Want to check whether Num is between 100 and 200?
	ok := td.CmpSmuggle(t, got, func(tr *Transaction) (int, error) {
		if tr.Body == nil ||
			tr.Body.Value == nil {
			return 0, errors.New("Request.Body or Request.Body.Value is nil")
		}
		if v, ok := tr.Body.Value.(*ValueNum); ok && v != nil {
			return v.Num, nil
		}
		return 0, errors.New("Request.Body.Value isn't *ValueNum or nil")
	}, td.Between(100, 200))
	fmt.Println("check Num by hand:", ok)

	// Same, but automagically generated...
	ok = td.CmpSmuggle(t, got, "Request.Body.Value.Num", td.Between(100, 200))
	fmt.Println("check Num using a fields-path:", ok)

	// And as Request is an anonymous field, can be simplified further
	// as it can be omitted
	ok = td.CmpSmuggle(t, got, "Body.Value.Num", td.Between(100, 200))
	fmt.Println("check Num using an other fields-path:", ok)

	// Note that maps and array/slices are supported
	got.Body.Value = map[string]any{
		"foo": []any{
			3: map[int]string{666: "bar"},
		},
	}
	ok = td.CmpSmuggle(t, got, "Body.Value[foo][3][666]", "bar")
	fmt.Println("check fields-path including maps/slices:", ok)

	// Output:
	// check Num by hand: true
	// check Num using a fields-path: true
	// check Num using an other fields-path: true
	// check fields-path including maps/slices: true

```{{% /expand%}}
## T.Smuggle shortcut

```go
func (t *T) Smuggle(got, fn , expectedValue any, args ...any) bool
```

Smuggle is a shortcut for:

```go
t.Cmp(got, td.Smuggle(fn, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Smuggle godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Smuggle).

### Examples

{{%expand "Convert example" %}}```go
	t := td.NewT(&testing.T{})

	got := int64(123)

	ok := t.Smuggle(got, func(n int64) int { return int(n) }, 123,
		"checks int64 got against an int value")
	fmt.Println(ok)

	ok = t.Smuggle("123", func(numStr string) (int, bool) {
		n, err := strconv.Atoi(numStr)
		return n, err == nil
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = t.Smuggle("123", func(numStr string) (int, bool, string) {
		n, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, false, "string must contain a number"
		}
		return n, true, ""
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	ok = t.Smuggle("123", func(numStr string) (int, error) { //nolint: gocritic
		return strconv.Atoi(numStr)
	}, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Short version :)
	ok = t.Smuggle("123", strconv.Atoi, td.Between(120, 130),
		"checks that number in %#v is in [120 .. 130]")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Lax example" %}}```go
	t := td.NewT(&testing.T{})

	// got is an int16 and Smuggle func input is an int64: it is OK
	got := int(123)

	ok := t.Smuggle(got, func(n int64) uint32 { return uint32(n) }, uint32(123))
	fmt.Println("got int16(123) → smuggle via int64 → uint32(123):", ok)

	// Output:
	// got int16(123) → smuggle via int64 → uint32(123): true

```{{% /expand%}}
{{%expand "Auto_unmarshal example" %}}```go
	t := td.NewT(&testing.T{})

	// Automatically json.Unmarshal to compare
	got := []byte(`{"a":1,"b":2}`)

	ok := t.Smuggle(got, func(b json.RawMessage) (r map[string]int, err error) {
		err = json.Unmarshal(b, &r)
		return
	}, map[string]int{
		"a": 1,
		"b": 2,
	})
	fmt.Println("JSON contents is OK:", ok)

	// Output:
	// JSON contents is OK: true

```{{% /expand%}}
{{%expand "Cast example" %}}```go
	t := td.NewT(&testing.T{})

	// A string containing JSON
	got := `{ "foo": 123 }`

	// Automatically cast a string to a json.RawMessage so td.JSON can operate
	ok := t.Smuggle(got, json.RawMessage{}, td.JSON(`{"foo":123}`))
	fmt.Println("JSON contents in string is OK:", ok)

	// Automatically read from io.Reader to a json.RawMessage
	ok = t.Smuggle(bytes.NewReader([]byte(got)), json.RawMessage{}, td.JSON(`{"foo":123}`))
	fmt.Println("JSON contents just read is OK:", ok)

	// Output:
	// JSON contents in string is OK: true
	// JSON contents just read is OK: true

```{{% /expand%}}
{{%expand "Complex example" %}}```go
	t := td.NewT(&testing.T{})

	// No end date but a start date and a duration
	type StartDuration struct {
		StartDate time.Time
		Duration  time.Duration
	}

	// Checks that end date is between 17th and 19th February both at 0h
	// for each of these durations in hours

	for _, duration := range []time.Duration{48 * time.Hour, 72 * time.Hour, 96 * time.Hour} {
		got := StartDuration{
			StartDate: time.Date(2018, time.February, 14, 12, 13, 14, 0, time.UTC),
			Duration:  duration,
		}

		// Simplest way, but in case of Between() failure, error will be bound
		// to DATA<smuggled>, not very clear...
		ok := t.Smuggle(got, func(sd StartDuration) time.Time {
			return sd.StartDate.Add(sd.Duration)
		}, td.Between(
			time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC)))
		fmt.Println(ok)

		// Name the computed value "ComputedEndDate" to render a Between() failure
		// more understandable, so error will be bound to DATA.ComputedEndDate
		ok = t.Smuggle(got, func(sd StartDuration) td.SmuggledGot {
			return td.SmuggledGot{
				Name: "ComputedEndDate",
				Got:  sd.StartDate.Add(sd.Duration),
			}
		}, td.Between(
			time.Date(2018, time.February, 17, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.February, 19, 0, 0, 0, 0, time.UTC)))
		fmt.Println(ok)
	}

	// Output:
	// false
	// false
	// true
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Interface example" %}}```go
	t := td.NewT(&testing.T{})

	gotTime, err := time.Parse(time.RFC3339, "2018-05-23T12:13:14Z")
	if err != nil {
		t.Fatal(err)
	}

	// Do not check the struct itself, but its stringified form
	ok := t.Smuggle(gotTime, func(s fmt.Stringer) string {
		return s.String()
	}, "2018-05-23 12:13:14 +0000 UTC")
	fmt.Println("stringified time.Time OK:", ok)

	// If got does not implement the fmt.Stringer interface, it fails
	// without calling the Smuggle func
	type MyTime time.Time
	ok = t.Smuggle(MyTime(gotTime), func(s fmt.Stringer) string {
		fmt.Println("Smuggle func called!")
		return s.String()
	}, "2018-05-23 12:13:14 +0000 UTC")
	fmt.Println("stringified MyTime OK:", ok)

	// Output:
	// stringified time.Time OK: true
	// stringified MyTime OK: false

```{{% /expand%}}
{{%expand "Field_path example" %}}```go
	t := td.NewT(&testing.T{})

	type Body struct {
		Name  string
		Value any
	}
	type Request struct {
		Body *Body
	}
	type Transaction struct {
		Request
	}
	type ValueNum struct {
		Num int
	}

	got := &Transaction{
		Request: Request{
			Body: &Body{
				Name:  "test",
				Value: &ValueNum{Num: 123},
			},
		},
	}

	// Want to check whether Num is between 100 and 200?
	ok := t.Smuggle(got, func(tr *Transaction) (int, error) {
		if tr.Body == nil ||
			tr.Body.Value == nil {
			return 0, errors.New("Request.Body or Request.Body.Value is nil")
		}
		if v, ok := tr.Body.Value.(*ValueNum); ok && v != nil {
			return v.Num, nil
		}
		return 0, errors.New("Request.Body.Value isn't *ValueNum or nil")
	}, td.Between(100, 200))
	fmt.Println("check Num by hand:", ok)

	// Same, but automagically generated...
	ok = t.Smuggle(got, "Request.Body.Value.Num", td.Between(100, 200))
	fmt.Println("check Num using a fields-path:", ok)

	// And as Request is an anonymous field, can be simplified further
	// as it can be omitted
	ok = t.Smuggle(got, "Body.Value.Num", td.Between(100, 200))
	fmt.Println("check Num using an other fields-path:", ok)

	// Note that maps and array/slices are supported
	got.Body.Value = map[string]any{
		"foo": []any{
			3: map[int]string{666: "bar"},
		},
	}
	ok = t.Smuggle(got, "Body.Value[foo][3][666]", "bar")
	fmt.Println("check fields-path including maps/slices:", ok)

	// Output:
	// check Num by hand: true
	// check Num using a fields-path: true
	// check Num using an other fields-path: true
	// check fields-path including maps/slices: true

```{{% /expand%}}
