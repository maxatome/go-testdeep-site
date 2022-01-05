---
title: "JSON"
weight: 10
---

```go
func JSON(expectedJSON interface{}, params ...interface{}) TestDeep
```

[`JSON`]({{< ref "JSON" >}}) operator allows to compare the JSON representation of data
against *expectedJSON*. *expectedJSON* can be a:

- `string` containing JSON data like `{"fullname":"Bob","age":42}`
- `string` containing a JSON filename, ending with ".json" (its
  content is [`ioutil.ReadFile`](https://pkg.go.dev/ioutil/#ReadFile) before unmarshaling)
- `[]byte` containing JSON data
- [`io.Reader`](https://pkg.go.dev/io/#Reader) stream containing JSON data (is [`ioutil.ReadAll`](https://pkg.go.dev/ioutil/#ReadAll) before
  unmarshaling)


*expectedJSON* JSON value can contain placeholders. The *params*
are for any placeholder parameters in *expectedJSON*. *params* can
contain [TestDeep operators]({{< ref "operators" >}}) as well as raw values. A placeholder can
be numeric like `$2` or named like `$name` and always references an
item in *params*.

Numeric placeholders reference the n'th "operators" item (starting
at 1). Named placeholders are used with [`Tag`]({{< ref "Tag" >}}) operator as follows:

```go
td.Cmp(t, gotValue,
  td.JSON(`{"fullname": $name, "age": $2, "gender": $3}`,
    td.Tag("name", td.HasPrefix("Foo")), // matches $1 and $name
    td.Between(41, 43),                  // matches only $2
    "male"))                             // matches only $3
```

Note that placeholders can be double-quoted as in:

```go
td.Cmp(t, gotValue,
  td.JSON(`{"fullname": "$name", "age": "$2", "gender": "$3"}`,
    td.Tag("name", td.HasPrefix("Foo")), // matches $1 and $name
    td.Between(41, 43),                  // matches only $2
    "male"))                             // matches only $3
```

It makes no difference whatever the underlying type of the replaced
item is (= double quoting a placeholder matching a number is not a
problem). It is just a matter of taste, double-quoting placeholders
can be preferred when the JSON data has to conform to the JSON
specification, like when used in a ".json" file.

[`JSON`]({{< ref "JSON" >}}) does its best to convert back the [`JSON`]({{< ref "JSON" >}}) corresponding to a
placeholder to the type of the placeholder or, if the placeholder
is an operator, to the type behind the operator. Allowing to do
things like:

```go
td.Cmp(t, gotValue, td.JSON(`{"foo":$1}`, []int{1, 2, 3, 4}))
td.Cmp(t, gotValue,
  td.JSON(`{"foo":$1}`, []interface{}{1, 2, td.Between(2, 4), 4}))
td.Cmp(t, gotValue, td.JSON(`{"foo":$1}`, td.Between(27, 32)))
```

Of course, it does this conversion only if the expected type can be
guessed. In the case the conversion cannot occur, data is compared
as is, in its freshly unmarshaled [`JSON`]({{< ref "JSON" >}}) form (so as `bool`, `float64`,
`string`, `[]interface{}`, `map[string]interface{}` or simply `nil`).

Note *expectedJSON* can be a `[]byte`, JSON filename or [`io.Reader`](https://pkg.go.dev/io/#Reader):

```go
td.Cmp(t, gotValue, td.JSON("file.json", td.Between(12, 34)))
td.Cmp(t, gotValue, td.JSON([]byte(`[1, $1, 3]`), td.Between(12, 34)))
td.Cmp(t, gotValue, td.JSON(osFile, td.Between(12, 34)))
```

A JSON filename ends with ".json".

To avoid a legit "$" `string` prefix causes a bad placeholder [`error`](https://pkg.go.dev/builtin/#error),
just double it to escape it. Note it is only needed when the "$" is
the first character of a `string`:

```go
td.Cmp(t, gotValue,
  td.JSON(`{"fullname": "$name", "details": "$$info", "age": $2}`,
    td.Tag("name", td.HasPrefix("Foo")), // matches $1 and $name
    td.Between(41, 43)))                 // matches only $2
```

For the "details" key, the raw value "`$info`" is expected, no
placeholders are involved here.

Note that [`Lax`]({{< ref "Lax" >}}) mode is automatically enabled by [`JSON`]({{< ref "JSON" >}}) operator to
simplify numeric tests.

Comments can be embedded in JSON data:

```go
td.Cmp(t, gotValue,
  td.JSON(`
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


Other [`JSON`]({{< ref "JSON" >}}) divergences:

- ',' can precede a '}' or a ']' (as in go);
- int_lit & float_lit numbers as defined in go spec are accepted;
- numbers can be prefixed by '+'.


Most operators can be directly embedded in [`JSON`]({{< ref "JSON" >}}) without requiring
any placeholder.

```go
td.Cmp(t, gotValue,
  td.JSON(`
{
  "fullname": HasPrefix("Foo"),
  "age":      Between(41, 43),
  "details":  SuperMapOf({
    "address": NotEmpty(),
    "car":     Any("Peugeot", "Tesla", "Jeep") // any of these
  })
}`))
```

Placeholders can be used anywhere, even in operators parameters as in:

```go
td.Cmp(t, gotValue, td.JSON(`{"fullname": HasPrefix($1)}`, "Zip"))
```

A few notes about operators embedding:

- [`SubMapOf`]({{< ref "SubMapOf" >}}) and [`SuperMapOf`]({{< ref "SuperMapOf" >}}) take only one parameter, a JSON object;
- the optional 3rd parameter of [`Between`]({{< ref "Between" >}}) has to be specified as a `string`
  and can be: "[]" or "[`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind)" (default), "[[" or "[`BoundsInOut`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind)",
  "]]" or "[`BoundsOutIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind)", "][" or "[`BoundsOutOut`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind)";
- not all operators are embeddable only the following are;
- [`All`]({{< ref "All" >}}), [`Any`]({{< ref "Any" >}}), [`ArrayEach`]({{< ref "ArrayEach" >}}), [`Bag`]({{< ref "Bag" >}}), [`Between`]({{< ref "Between" >}}), [`Contains`]({{< ref "Contains" >}}), [`ContainsKey`]({{< ref "ContainsKey" >}}), [`Empty`]({{< ref "Empty" >}}), [`Gt`]({{< ref "Gt" >}}),
  [`Gte`]({{< ref "Gte" >}}), [`HasPrefix`]({{< ref "HasPrefix" >}}), [`HasSuffix`]({{< ref "HasSuffix" >}}), [`Ignore`]({{< ref "Ignore" >}}), JSONPointer, [`Keys`]({{< ref "Keys" >}}), [`Len`]({{< ref "Len" >}}), [`Lt`]({{< ref "Lt" >}}), [`Lte`]({{< ref "Lte" >}}),
  [`MapEach`]({{< ref "MapEach" >}}), [`N`]({{< ref "N" >}}), [`NaN`]({{< ref "NaN" >}}), [`Nil`]({{< ref "Nil" >}}), [`None`]({{< ref "None" >}}), [`Not`]({{< ref "Not" >}}), [`NotAny`]({{< ref "NotAny" >}}), [`NotEmpty`]({{< ref "NotEmpty" >}}), [`NotNaN`]({{< ref "NotNaN" >}}), [`NotNil`]({{< ref "NotNil" >}}),
  [`NotZero`]({{< ref "NotZero" >}}), [`Re`]({{< ref "Re" >}}), [`ReAll`]({{< ref "ReAll" >}}), [`Set`]({{< ref "Set" >}}), [`SubBagOf`]({{< ref "SubBagOf" >}}), [`SubMapOf`]({{< ref "SubMapOf" >}}), [`SubSetOf`]({{< ref "SubSetOf" >}}), [`SuperBagOf`]({{< ref "SuperBagOf" >}}),
  [`SuperMapOf`]({{< ref "SuperMapOf" >}}), [`SuperSetOf`]({{< ref "SuperSetOf" >}}), [`Values`]({{< ref "Values" >}}) and [`Zero`]({{< ref "Zero" >}}).


Operators taking no parameters can also be directly embedded in
JSON data using `$^OperatorName` or "`$^OperatorName`" notation. They
are named shortcut operators (they predate the above operators embedding
but they subsist for compatibility):

```go
td.Cmp(t, gotValue, td.JSON(`{"id": $1}`, td.NotZero()))
```

can be written as:

```go
td.Cmp(t, gotValue, td.JSON(`{"id": $^NotZero}`))
```

or

```go
td.Cmp(t, gotValue, td.JSON(`{"id": "$^NotZero"}`))
```

As for placeholders, there is no differences between `$^NotZero` and
"`$^NotZero`".

The allowed shortcut operators follow:

- [`Empty`]({{< ref "Empty" >}})    → `$^Empty`
- [`Ignore`]({{< ref "Ignore" >}})   → `$^Ignore`
- [`NaN`]({{< ref "NaN" >}})      → `$^NaN`
- [`Nil`]({{< ref "Nil" >}})      → `$^Nil`
- [`NotEmpty`]({{< ref "NotEmpty" >}}) → `$^NotEmpty`
- [`NotNaN`]({{< ref "NotNaN" >}})   → `$^NotNaN`
- [`NotNil`]({{< ref "NotNil" >}})   → `$^NotNil`
- [`NotZero`]({{< ref "NotZero" >}})  → `$^NotZero`
- [`Zero`]({{< ref "Zero" >}})     → `$^Zero`


[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of the *expectedJSON*
[`json.Unmarshal`](https://pkg.go.dev/json/#Unmarshal)'ed. So it can be `bool`, `string`, `float64`,
`[]interface{}`, `map[string]interface{}` or `interface{}` in case
*expectedJSON* is "null".


> See also [<i class='fas fa-book'></i> JSON godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#JSON).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := td.Cmp(t, got, td.JSON(`{"age":42,"fullname":"Bob"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = td.Cmp(t, got, td.JSON(`{"fullname":"Bob","age":42}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = td.Cmp(t, got, td.JSON(`
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`))
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.Cmp(t, got, td.JSON(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with gender field:", ok)

	ok = td.Cmp(t, got, td.JSON(`{"fullname":"Bob"}`))
	fmt.Println("check got with fullname only:", ok)

	ok = td.Cmp(t, true, td.JSON(`true`))
	fmt.Println("check boolean got is true:", ok)

	ok = td.Cmp(t, 42, td.JSON(`42`))
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = td.Cmp(t, got, td.JSON(`null`))
	fmt.Println("check nil got is null:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with gender field: false
	// check got with fullname only: false
	// check boolean got is true: true
	// check numeric got is 42: true
	// check nil got is null: true

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := &testing.T{}

	type Person struct {
		Fullname string    `json:"fullname"`
		Age      int       `json:"age"`
		Children []*Person `json:"children,omitempty"`
	}

	got := &Person{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.Cmp(t, got, td.JSON(`{"age": $1, "fullname": $2}`, 42, "Bob Foobar"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"age": $1, "fullname": $2}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar")))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"age": "$1", "fullname": "$2"}`,
			td.Between(40, 45),
			td.HasSuffix("Foobar")))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"age": $age, "fullname": $name}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.HasSuffix("Foobar"))))
	fmt.Println("check got with named placeholders:", ok)

	got.Children = []*Person{
		{Fullname: "Alice", Age: 28},
		{Fullname: "Brian", Age: 22},
	}
	ok = td.Cmp(t, got,
		td.JSON(`{"age": $age, "fullname": $name, "children": $children}`,
			td.Tag("age", td.Between(40, 45)),
			td.Tag("name", td.HasSuffix("Foobar")),
			td.Tag("children", td.Bag(
				&Person{Fullname: "Brian", Age: 22},
				&Person{Fullname: "Alice", Age: 28},
			))))
	fmt.Println("check got w/named placeholders, and children w/go structs:", ok)

	ok = td.Cmp(t, got,
		td.JSON(`{"age": Between($1, $2), "fullname": HasSuffix($suffix), "children": Len(2)}`,
			40, 45,
			td.Tag("suffix", "Foobar")))
	fmt.Println("check got w/num & named placeholders:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got w/named placeholders, and children w/go structs: true
	// check got w/num & named placeholders: true

```{{% /expand%}}
{{%expand "Embedding example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.Cmp(t, got, td.JSON(`{"age": NotZero(), "fullname": NotEmpty()}`))
	fmt.Println("check got with simple operators:", ok)

	ok = td.Cmp(t, got, td.JSON(`{"age": $^NotZero, "fullname": $^NotEmpty}`))
	fmt.Println("check got with operator shortcuts:", ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "age":      Between(40, 42, "]]"), // in ]40; 42]
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar")  // ← comma is optional here
  )
}`))
	fmt.Println("check got with complex operators:", ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "age":      Between(40, 42, "]["), // in ]40; 42[ → 42 excluded
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar"),
  )
}`))
	fmt.Println("check got with complex operators:", ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "age":      Between($1, $2, $3), // in ]40; 42]
  "fullname": All(
    HasPrefix($4),
    HasSuffix("bar")  // ← comma is optional here
  )
}`,
		40, 42, td.BoundsOutIn,
		"Bob"))
	fmt.Println("check got with complex operators, w/placeholder args:", ok)

	// Output:
	// check got with simple operators: true
	// check got with operator shortcuts: true
	// check got with complex operators: true
	// check got with complex operators: false
	// check got with complex operators, w/placeholder args: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.Cmp(t, got,
		td.JSON(filename,
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
		td.JSON(file,
			td.Tag("name", td.HasPrefix("Bob")),
			td.Tag("age", td.Between(40, 45)),
			td.Tag("gender", td.Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
## CmpJSON shortcut

```go
func CmpJSON(t TestingT, got, expectedJSON interface{}, params []interface{}, args ...interface{}) bool
```

CmpJSON is a shortcut for:

```go
td.Cmp(t, got, td.JSON(expectedJSON, params...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpJSON godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpJSON).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := td.CmpJSON(t, got, `{"age":42,"fullname":"Bob"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = td.CmpJSON(t, got, `{"fullname":"Bob","age":42}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = td.CmpJSON(t, got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = td.CmpJSON(t, got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with gender field:", ok)

	ok = td.CmpJSON(t, got, `{"fullname":"Bob"}`, nil)
	fmt.Println("check got with fullname only:", ok)

	ok = td.CmpJSON(t, true, `true`, nil)
	fmt.Println("check boolean got is true:", ok)

	ok = td.CmpJSON(t, 42, `42`, nil)
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = td.CmpJSON(t, got, `null`, nil)
	fmt.Println("check nil got is null:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with gender field: false
	// check got with fullname only: false
	// check boolean got is true: true
	// check numeric got is 42: true
	// check nil got is null: true

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := &testing.T{}

	type Person struct {
		Fullname string    `json:"fullname"`
		Age      int       `json:"age"`
		Children []*Person `json:"children,omitempty"`
	}

	got := &Person{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.CmpJSON(t, got, `{"age": $1, "fullname": $2}`, []interface{}{42, "Bob Foobar"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = td.CmpJSON(t, got, `{"age": $1, "fullname": $2}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar")})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = td.CmpJSON(t, got, `{"age": "$1", "fullname": "$2"}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar")})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = td.CmpJSON(t, got, `{"age": $age, "fullname": $name}`, []interface{}{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar"))})
	fmt.Println("check got with named placeholders:", ok)

	got.Children = []*Person{
		{Fullname: "Alice", Age: 28},
		{Fullname: "Brian", Age: 22},
	}
	ok = td.CmpJSON(t, got, `{"age": $age, "fullname": $name, "children": $children}`, []interface{}{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("children", td.Bag(
		&Person{Fullname: "Brian", Age: 22},
		&Person{Fullname: "Alice", Age: 28},
	))})
	fmt.Println("check got w/named placeholders, and children w/go structs:", ok)

	ok = td.CmpJSON(t, got, `{"age": Between($1, $2), "fullname": HasSuffix($suffix), "children": Len(2)}`, []interface{}{40, 45, td.Tag("suffix", "Foobar")})
	fmt.Println("check got w/num & named placeholders:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got w/named placeholders, and children w/go structs: true
	// check got w/num & named placeholders: true

```{{% /expand%}}
{{%expand "Embedding example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := td.CmpJSON(t, got, `{"age": NotZero(), "fullname": NotEmpty()}`, nil)
	fmt.Println("check got with simple operators:", ok)

	ok = td.CmpJSON(t, got, `{"age": $^NotZero, "fullname": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	ok = td.CmpJSON(t, got, `
{
  "age":      Between(40, 42, "]]"), // in ]40; 42]
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar")  // ← comma is optional here
  )
}`, nil)
	fmt.Println("check got with complex operators:", ok)

	ok = td.CmpJSON(t, got, `
{
  "age":      Between(40, 42, "]["), // in ]40; 42[ → 42 excluded
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar"),
  )
}`, nil)
	fmt.Println("check got with complex operators:", ok)

	ok = td.CmpJSON(t, got, `
{
  "age":      Between($1, $2, $3), // in ]40; 42]
  "fullname": All(
    HasPrefix($4),
    HasSuffix("bar")  // ← comma is optional here
  )
}`, []interface{}{40, 42, td.BoundsOutIn, "Bob"})
	fmt.Println("check got with complex operators, w/placeholder args:", ok)

	// Output:
	// check got with simple operators: true
	// check got with operator shortcuts: true
	// check got with complex operators: true
	// check got with complex operators: false
	// check got with complex operators, w/placeholder args: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := td.CmpJSON(t, got, filename, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = td.CmpJSON(t, got, file, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
## T.JSON shortcut

```go
func (t *T) JSON(got, expectedJSON interface{}, params []interface{}, args ...interface{}) bool
```

[`JSON`]({{< ref "JSON" >}}) is a shortcut for:

```go
t.Cmp(got, td.JSON(expectedJSON, params...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.JSON godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.JSON).

### Examples

{{%expand "Basic example" %}}```go
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := t.JSON(got, `{"age":42,"fullname":"Bob"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = t.JSON(got, `{"fullname":"Bob","age":42}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = t.JSON(got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = t.JSON(got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with gender field:", ok)

	ok = t.JSON(got, `{"fullname":"Bob"}`, nil)
	fmt.Println("check got with fullname only:", ok)

	ok = t.JSON(true, `true`, nil)
	fmt.Println("check boolean got is true:", ok)

	ok = t.JSON(42, `42`, nil)
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = t.JSON(got, `null`, nil)
	fmt.Println("check nil got is null:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with gender field: false
	// check got with fullname only: false
	// check boolean got is true: true
	// check numeric got is 42: true
	// check nil got is null: true

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := td.NewT(&testing.T{})

	type Person struct {
		Fullname string    `json:"fullname"`
		Age      int       `json:"age"`
		Children []*Person `json:"children,omitempty"`
	}

	got := &Person{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := t.JSON(got, `{"age": $1, "fullname": $2}`, []interface{}{42, "Bob Foobar"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = t.JSON(got, `{"age": $1, "fullname": $2}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar")})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = t.JSON(got, `{"age": "$1", "fullname": "$2"}`, []interface{}{td.Between(40, 45), td.HasSuffix("Foobar")})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = t.JSON(got, `{"age": $age, "fullname": $name}`, []interface{}{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar"))})
	fmt.Println("check got with named placeholders:", ok)

	got.Children = []*Person{
		{Fullname: "Alice", Age: 28},
		{Fullname: "Brian", Age: 22},
	}
	ok = t.JSON(got, `{"age": $age, "fullname": $name, "children": $children}`, []interface{}{td.Tag("age", td.Between(40, 45)), td.Tag("name", td.HasSuffix("Foobar")), td.Tag("children", td.Bag(
		&Person{Fullname: "Brian", Age: 22},
		&Person{Fullname: "Alice", Age: 28},
	))})
	fmt.Println("check got w/named placeholders, and children w/go structs:", ok)

	ok = t.JSON(got, `{"age": Between($1, $2), "fullname": HasSuffix($suffix), "children": Len(2)}`, []interface{}{40, 45, td.Tag("suffix", "Foobar")})
	fmt.Println("check got w/num & named placeholders:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true
	// check got w/named placeholders, and children w/go structs: true
	// check got w/num & named placeholders: true

```{{% /expand%}}
{{%expand "Embedding example" %}}```go
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := t.JSON(got, `{"age": NotZero(), "fullname": NotEmpty()}`, nil)
	fmt.Println("check got with simple operators:", ok)

	ok = t.JSON(got, `{"age": $^NotZero, "fullname": $^NotEmpty}`, nil)
	fmt.Println("check got with operator shortcuts:", ok)

	ok = t.JSON(got, `
{
  "age":      Between(40, 42, "]]"), // in ]40; 42]
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar")  // ← comma is optional here
  )
}`, nil)
	fmt.Println("check got with complex operators:", ok)

	ok = t.JSON(got, `
{
  "age":      Between(40, 42, "]["), // in ]40; 42[ → 42 excluded
  "fullname": All(
    HasPrefix("Bob"),
    HasSuffix("bar"),
  )
}`, nil)
	fmt.Println("check got with complex operators:", ok)

	ok = t.JSON(got, `
{
  "age":      Between($1, $2, $3), // in ]40; 42]
  "fullname": All(
    HasPrefix($4),
    HasSuffix("bar")  // ← comma is optional here
  )
}`, []interface{}{40, 42, td.BoundsOutIn, "Bob"})
	fmt.Println("check got with complex operators, w/placeholder args:", ok)

	// Output:
	// check got with simple operators: true
	// check got with operator shortcuts: true
	// check got with complex operators: true
	// check got with complex operators: false
	// check got with complex operators, w/placeholder args: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := td.NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := t.JSON(got, filename, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = t.JSON(got, file, []interface{}{td.Tag("name", td.HasPrefix("Bob")), td.Tag("age", td.Between(40, 45)), td.Tag("gender", td.Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
