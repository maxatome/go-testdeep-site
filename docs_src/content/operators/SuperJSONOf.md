---
title: "SuperJSONOf"
weight: 10
---

```go
func SuperJSONOf(expectedJSON any, params ...any) TestDeep
```

SuperJSONOf operator allows to compare the JSON representation of
data against *expectedJSON*. Unlike JSON operator, marshaled data
must be a JSON object/map (aka {…}). *expectedJSON* can be a:

- `string` containing JSON data like `{"fullname":"Bob","age":42}`
- `string` containing a JSON filename, ending with ".json" (its
  content is [`os.ReadFile`](https://pkg.go.dev/os#ReadFile) before unmarshaling)
- `[]byte` containing JSON data
- [`encoding/json.RawMessage`](https://pkg.go.dev/encoding/json#RawMessage) containing JSON data
- [`io.Reader`](https://pkg.go.dev/io#Reader) stream containing JSON data (is [`io.ReadAll`](https://pkg.go.dev/io#ReadAll) before
  unmarshaling)


JSON data contained in *expectedJSON* must be a JSON object/map
(aka {…}) too. During a match, each expected entry should match in
the compared map. But some entries in the compared map may not be
expected.

```go
type MyStruct struct {
  Name string `json:"name"`
  Age  int    `json:"age"`
  City string `json:"city"`
}
got := MyStruct{
  Name: "Bob",
  Age:  42,
  City: "TestCity",
}
td.Cmp(t, got, td.SuperJSONOf(`{"name": "Bob", "age": 42}`))  // succeeds
td.Cmp(t, got, td.SuperJSONOf(`{"name": "Bob", "zip": 666}`)) // fails, miss "zip"
```

*expectedJSON* JSON value can contain placeholders. The *params* are
for any placeholder parameters in *expectedJSON*. *params* can contain
[TestDeep operators]({{< ref "operators" >}}) as well as raw values. A placeholder can be
numeric like `$2` or named like `$name` and always references an item
in *params*.

Numeric placeholders reference the n'th "operators" item (starting
at 1). Named placeholders are used with [`Tag`]({{< ref "Tag" >}}) operator as follows:

```go
td.Cmp(t, gotValue,
  SuperJSONOf(`{"fullname": $name, "age": $2, "gender": $3}`,
    td.Tag("name", td.HasPrefix("Foo")), // matches $1 and $name
    td.Between(41, 43),                  // matches only $2
    "male"))                             // matches only $3
```

Note that placeholders can be double-quoted as in:

```go
td.Cmp(t, gotValue,
  td.SuperJSONOf(`{"fullname": "$name", "age": "$2", "gender": "$3"}`,
    td.Tag("name", td.HasPrefix("Foo")), // matches $1 and $name
    td.Between(41, 43),                  // matches only $2
    "male"))                             // matches only $3
```

It makes no difference whatever the underlying type of the replaced
item is (= double quoting a placeholder matching a number is not a
problem). It is just a matter of taste, double-quoting placeholders
can be preferred when the JSON data has to conform to the JSON
specification, like when used in a ".json" file.

SuperJSONOf does its best to convert back the JSON corresponding to a
placeholder to the type of the placeholder or, if the placeholder
is an operator, to the type behind the operator. Allowing to do
things like:

```go
td.Cmp(t, gotValue,
  td.SuperJSONOf(`{"foo":$1}`, []int{1, 2, 3, 4}))
td.Cmp(t, gotValue,
  td.SuperJSONOf(`{"foo":$1}`, []any{1, 2, td.Between(2, 4), 4}))
td.Cmp(t, gotValue,
  td.SuperJSONOf(`{"foo":$1}`, td.Between(27, 32)))
```

Of course, it does this conversion only if the expected type can be
guessed. In the case the conversion cannot occur, data is compared
as is, in its freshly unmarshaled JSON form (so as `bool`, `float64`,
`string`, `[]any`, `map[string]any` or simply `nil`).

Note *expectedJSON* can be a `[]byte`, an [`encoding/json.RawMessage`](https://pkg.go.dev/encoding/json#RawMessage), a
JSON filename or a [`io.Reader`](https://pkg.go.dev/io#Reader):

```go
td.Cmp(t, gotValue, td.SuperJSONOf("file.json", td.Between(12, 34)))
td.Cmp(t, gotValue, td.SuperJSONOf([]byte(`[1, $1, 3]`), td.Between(12, 34)))
td.Cmp(t, gotValue, td.SuperJSONOf(osFile, td.Between(12, 34)))
```

A JSON filename ends with ".json".

To avoid a legit "$" `string` prefix causes a bad placeholder [`error`](https://pkg.go.dev/builtin#error),
just double it to escape it. Note it is only needed when the "$" is
the first character of a `string`:

```go
td.Cmp(t, gotValue,
  td.SuperJSONOf(`{"fullname": "$name", "details": "$$info", "age": $2}`,
    td.Tag("name", td.HasPrefix("Foo")), // matches $1 and $name
    td.Between(41, 43)))                 // matches only $2
```

For the "details" key, the raw value "`$info`" is expected, no
placeholders are involved here.

Note that [`Lax`]({{< ref "Lax" >}}) mode is automatically enabled by SuperJSONOf operator to
simplify numeric tests.

Comments can be embedded in JSON data:

```go
td.Cmp(t, gotValue,
  td.SuperJSONOf(`
{
  // A guy properties:
  "fullname": "$name",  // The full name of the guy
  "details":  "$$info", // Literally "$info", thanks to "$" escape
  "age":      $2        /* The age of the guy:
                           - placeholder unquoted, but could be without
                             any change
                           - to demonstrate a multi-lines comment */
}`,
    td.Tag("name", td.HasPrefix("Foo")), // matches $1 and $name
    td.Between(41, 43)))                 // matches only $2
```

Comments, like in go, have 2 forms. To quote the Go language specification:

- line comments start with the character sequence // and stop at the
  end of the line.
- multi-lines comments start with the character sequence /* and stop
  with the first subsequent character sequence */.


Other JSON divergences:

- ',' can precede a '}' or a ']' (as in go);
- strings can contain non-escaped \n, \r and \t;
- raw strings are accepted (r{raw}, r!raw!, …), see below;
- int_lit & float_lit numbers as defined in go spec are accepted;
- numbers can be prefixed by '+'.


Most operators can be directly embedded in SuperJSONOf without requiring
any placeholder. If an operators does not take `any` parameter, the
parenthesis can be omitted.

```go
td.Cmp(t, gotValue,
  td.SuperJSONOf(`
{
  "fullname": HasPrefix("Foo"),
  "age":      Between(41, 43),
  "details":  SuperMapOf({
    "address": NotEmpty, // () are optional when no parameters
    "car":     Any("Peugeot", "Tesla", "Jeep") // any of these
  })
}`))
```

Placeholders can be used `any`where, even in operators parameters as in:

```go
td.Cmp(t, gotValue, td.SuperJSONOf(`{"fullname": HasPrefix($1)}`, "Zip"))
```

A few notes about operators embedding:

- [`SubMapOf`]({{< ref "SubMapOf" >}}) and [`SuperMapOf`]({{< ref "SuperMapOf" >}}) take only one parameter, a JSON object;
- the optional 3rd parameter of [`Between`]({{< ref "Between" >}}) has to be specified as a `string`
  and can be: "[]" or "BoundsInIn" (default), "[[" or "BoundsInOut",
  "]]" or "BoundsOutIn", "][" or "BoundsOutOut";
- not all operators are embeddable only the following are: [`All`]({{< ref "All" >}}),
  [`Any`]({{< ref "Any" >}}), [`ArrayEach`]({{< ref "ArrayEach" >}}), [`Bag`]({{< ref "Bag" >}}), [`Between`]({{< ref "Between" >}}), [`Contains`]({{< ref "Contains" >}}),
  [`ContainsKey`]({{< ref "ContainsKey" >}}), [`Empty`]({{< ref "Empty" >}}), [`First`]({{< ref "First" >}}), [`Grep`]({{< ref "Grep" >}}), [`Gt`]({{< ref "Gt" >}}), [`Gte`]({{< ref "Gte" >}}),
  [`HasPrefix`]({{< ref "HasPrefix" >}}), [`HasSuffix`]({{< ref "HasSuffix" >}}), [`Ignore`]({{< ref "Ignore" >}}), [`JSONPointer`]({{< ref "JSONPointer" >}}), [`Keys`]({{< ref "Keys" >}}),
  [`Last`]({{< ref "Last" >}}), [`Len`]({{< ref "Len" >}}), [`Lt`]({{< ref "Lt" >}}), [`Lte`]({{< ref "Lte" >}}), [`MapEach`]({{< ref "MapEach" >}}), [`N`]({{< ref "N" >}}), [`NaN`]({{< ref "NaN" >}}), [`Nil`]({{< ref "Nil" >}}),
  [`None`]({{< ref "None" >}}), [`Not`]({{< ref "Not" >}}), [`NotAny`]({{< ref "NotAny" >}}), [`NotEmpty`]({{< ref "NotEmpty" >}}), [`NotNaN`]({{< ref "NotNaN" >}}), [`NotNil`]({{< ref "NotNil" >}}),
  [`NotZero`]({{< ref "NotZero" >}}), [`Re`]({{< ref "Re" >}}), [`ReAll`]({{< ref "ReAll" >}}), [`Set`]({{< ref "Set" >}}), [`SubBagOf`]({{< ref "SubBagOf" >}}), [`SubMapOf`]({{< ref "SubMapOf" >}}),
  [`SubSetOf`]({{< ref "SubSetOf" >}}), [`SuperBagOf`]({{< ref "SuperBagOf" >}}), [`SuperMapOf`]({{< ref "SuperMapOf" >}}), [`SuperSetOf`]({{< ref "SuperSetOf" >}}), [`Values`]({{< ref "Values" >}})
  and [`Zero`]({{< ref "Zero" >}}).


It is also possible to embed operators in JSON strings. This way,
the JSON specification can be fulfilled. To avoid collision with
possible strings, just prefix the first operator name with
"$^". The previous example becomes:

```go
td.Cmp(t, gotValue,
  td.SuperJSONOf(`
{
  "fullname": "$^HasPrefix(\"Foo\")",
  "age":      "$^Between(41, 43)",
  "details":  "$^SuperMapOf({
    \"address\": NotEmpty, // () are optional when no parameters
    \"car\":     Any(\"Peugeot\", \"Tesla\", \"Jeep\") // any of these
  })"
}`))
```

As you can see, in this case, strings in strings have to be
escaped. Fortunately, newlines are accepted, but unfortunately they
are forbidden by JSON specification. To avoid too much escaping,
raw strings are accepted. A raw `string` is a "r" followed by a
delimiter, the corresponding delimiter closes the `string`. The
following raw strings are all the same as "foo\\bar(\"zip\")!":

- r'foo\bar"zip"!'
- r,foo\bar"zip"!,
- r%foo\bar"zip"!%
- r(foo\bar("zip")!)
- r{foo\bar("zip")!}
- r[foo\bar("zip")!]
- r<foo\bar("zip")!>


So non-bracketing delimiters use the same character before and
after, but the 4 sorts of ASCII brackets (round, angle, square,
curly) all nest: r[x[y]z] equals "x[y]z". The end delimiter cannot
be escaped.

With raw strings, the previous example becomes:

```go
td.Cmp(t, gotValue,
  td.SuperJSONOf(`
{
  "fullname": "$^HasPrefix(r<Foo>)",
  "age":      "$^Between(41, 43)",
  "details":  "$^SuperMapOf({
    r<address>: NotEmpty, // () are optional when no parameters
    r<car>:     Any(r<Peugeot>, r<Tesla>, r<Jeep>) // any of these
  })"
}`))
```

Note that raw strings are accepted `any`where, not only in original
JSON strings.

To be complete, $^ can prefix an operator even outside a
`string`. This is accepted for compatibility purpose as the first
operator embedding feature used this way to embed some operators.

So the following calls are all equivalent:

```go
td.Cmp(t, gotValue, td.SuperJSONOf(`{"id": $1}`, td.NotZero()))
td.Cmp(t, gotValue, td.SuperJSONOf(`{"id": NotZero}`))
td.Cmp(t, gotValue, td.SuperJSONOf(`{"id": NotZero()}`))
td.Cmp(t, gotValue, td.SuperJSONOf(`{"id": $^NotZero}`))
td.Cmp(t, gotValue, td.SuperJSONOf(`{"id": $^NotZero()}`))
td.Cmp(t, gotValue, td.SuperJSONOf(`{"id": "$^NotZero"}`))
td.Cmp(t, gotValue, td.SuperJSONOf(`{"id": "$^NotZero()"}`))
```

As for placeholders, there is no differences between `$^NotZero` and
"`$^NotZero`".

Tip: when an [`io.Reader`](https://pkg.go.dev/io#Reader) is expected to contain JSON data, it
cannot be tested directly, but using the [`Smuggle`]({{< ref "Smuggle" >}}) operator simply
solves the problem:

```go
var body io.Reader
// …
td.Cmp(t, body, td.Smuggle(json.RawMessage{}, td.SuperJSONOf(`{"foo":1}`)))
// or equally
td.Cmp(t, body, td.Smuggle(json.RawMessage(nil), td.SuperJSONOf(`{"foo":1}`)))
```

[`Smuggle`]({{< ref "Smuggle" >}}) reads from body into an [`encoding/json.RawMessage`](https://pkg.go.dev/encoding/json#RawMessage) then
this buffer is unmarshaled by SuperJSONOf operator before the comparison.

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the `map[string]any` type.

> See also [`JSON`]({{< ref "JSON" >}}), [`JSONPointer`]({{< ref "JSONPointer" >}}) and [`SubJSONOf`]({{< ref "SubJSONOf" >}}).


> See also [<i class='fas fa-book'></i> SuperJSONOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperJSONOf).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	ok := td.Cmp(t, got, td.SuperJSONOf(`{"age":42,"fullname":"Bob","gender":"male"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = td.Cmp(t, got, td.SuperJSONOf(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = td.Cmp(t, got, td.SuperJSONOf(`
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42,    /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
  "gender":   "male" // The gender!
}`))
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"fullname":"Bob","gender":"male","details":{}}`))
	fmt.Println("check got with details field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with details field: false

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	ok := td.Cmp(t, got,
		td.SuperJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			42, "Bob Foobar", "male"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"age": $1, "fullname": $2, "gender": $3}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar"),
			td.NotEmpty()))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"age": "$1", "fullname": "$2", "gender": "$3"}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar"),
			td.NotEmpty()))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"age": $age, "fullname": $name, "gender": $gender}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.HasSuffix("Foobar")),
			td.Tag("gender", td.NotEmpty())))
	fmt.Println("check got with named placeholders:", ok)

	ok = td.Cmp(t, got,
		td.SuperJSONOf(`{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`))
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = os.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.Cmp(t, got,
		td.SuperJSONOf(filename,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.Cmp(t, got,
		td.SuperJSONOf(file,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
## CmpSuperJSONOf shortcut

```go
func CmpSuperJSONOf(t TestingT, got, expectedJSON any, params []any, args ...any) bool
```

CmpSuperJSONOf is a shortcut for:

```go
td.Cmp(t, got, td.SuperJSONOf(expectedJSON, params...), args...)
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


> See also [<i class='fas fa-book'></i> CmpSuperJSONOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSuperJSONOf).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	ok := td.CmpSuperJSONOf(t, got, `{"age":42,"fullname":"Bob","gender":"male"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = td.CmpSuperJSONOf(t, got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42,    /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
  "gender":   "male" // The gender!
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"fullname":"Bob","gender":"male","details":{}}`, nil)
	fmt.Println("check got with details field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with details field: false

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	ok := td.CmpSuperJSONOf(t, got, `{"age": $1, "fullname": $2, "gender": $3}`, []any{42, "Bob Foobar", "male"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"age": $1, "fullname": $2, "gender": $3}`, []any{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"age": "$1", "fullname": "$2", "gender": "$3"}`, []any{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"age": $age, "fullname": $name, "gender": $gender}`, []any{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("gender", td.NotEmpty())})
	fmt.Println("check got with named placeholders:", ok)

	ok = td.CmpSuperJSONOf(t, got, `{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = os.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.CmpSuperJSONOf(t, got, filename, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.CmpSuperJSONOf(t, got, file, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
## T.SuperJSONOf shortcut

```go
func (t *T) SuperJSONOf(got, expectedJSON any, params []any, args ...any) bool
```

SuperJSONOf is a shortcut for:

```go
t.Cmp(got, td.SuperJSONOf(expectedJSON, params...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SuperJSONOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SuperJSONOf).

### Examples

{{%expand "Basic example" %}}```go
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	ok := t.SuperJSONOf(got, `{"age":42,"fullname":"Bob","gender":"male"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = t.SuperJSONOf(got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = t.SuperJSONOf(got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42,    /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
  "gender":   "male" // The gender!
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = t.SuperJSONOf(got, `{"fullname":"Bob","gender":"male","details":{}}`, nil)
	fmt.Println("check got with details field:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with details field: false

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	ok := t.SuperJSONOf(got, `{"age": $1, "fullname": $2, "gender": $3}`, []any{42, "Bob Foobar", "male"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = t.SuperJSONOf(got, `{"age": $1, "fullname": $2, "gender": $3}`, []any{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = t.SuperJSONOf(got, `{"age": "$1", "fullname": "$2", "gender": "$3"}`, []any{td.Between(40, 45), td.HasSuffix("Foobar"), td.NotEmpty()})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = t.SuperJSONOf(got, `{"age": $age, "fullname": $name, "gender": $gender}`, []any{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("gender", td.NotEmpty())})
	fmt.Println("check got with named placeholders:", ok)

	ok = t.SuperJSONOf(got, `{"age": $^NotZero, "fullname": $^NotEmpty, "gender": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got with operator shortcuts: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		City     string `json:"city"`
		Zip      int    `json:"zip"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
		City:     "TestCity",
		Zip:      666,
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = os.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := t.SuperJSONOf(got, filename, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = t.SuperJSONOf(got, file, []any{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
