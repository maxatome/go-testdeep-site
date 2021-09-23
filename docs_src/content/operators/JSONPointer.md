---
title: "JSONPointer"
weight: 10
---

```go
func JSONPointer(pointer string, expectedValue interface{}) TestDeep
```

JSONPointer is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It takes the JSON
representation of data, gets the value corresponding to the [`JSON`]({{< ref "JSON" >}})
pointer *pointer* (as RFC 6901 specifies it) and compares it to
*expectedValue*.

[`Lax`]({{< ref "Lax" >}}) mode is automatically enabled to simplify numeric tests.

JSONPointer does its best to convert back the [`JSON`]({{< ref "JSON" >}}) pointed data to
the type of *expectedValue* or to the type behind the
*expectedValue* operator, if it is an operator. Allowing to do
things like:

```go
type Item struct {
  Val  int   `json:"val"`
  Next *Item `json:"next"`
}
got := Item{Val: 1, Next: &Item{Val: 2, Next: &Item{Val: 3}}}

td.Cmp(t, got, td.JSONPointer("/next/next", Item{Val: 3}))
td.Cmp(t, got, td.JSONPointer("/next/next", &Item{Val: 3}))
td.Cmp(t,
  got,
  td.JSONPointer("/next/next",
    td.Struct(Item{}, td.StructFields{"Val": td.Gte(3)})),
)

got := map[string]int64{"zzz": 42} // 42 is int64 here
td.Cmp(t, got, td.JSONPointer("/zzz", 42))
td.Cmp(t, got, td.JSONPointer("/zzz", td.Between(40, 45)))
```

Of course, it does this conversion only if the expected type can be
guessed. In the case the conversion cannot occur, data is compared
as is, in its freshly unmarshalled [`JSON`]({{< ref "JSON" >}}) form (so as `bool`, `float64`,
`string`, `[]interface{}`, `map[string]interface{}` or simply `nil`).

Note that as any [TestDeep operator]({{< ref "operators" >}}) can be used as *expectedValue*,
[`JSON`]({{< ref "JSON" >}}) operator works out of the box:

```go
got := json.RawMessage(`{"foo":{"bar": {"zip": true}}}`)
td.Cmp(t, got, td.JSONPointer("/foo/bar", td.JSON(`{"zip": true}`)))
```

It can be used with structs lacking json tags. In this case, fields
names have to be used in [`JSON`]({{< ref "JSON" >}}) pointer:

```go
type Item struct {
  Val  int
  Next *Item
}
got := Item{Val: 1, Next: &Item{Val: 2, Next: &Item{Val: 3}}}

td.Cmp(t, got, td.JSONPointer("/Next/Next", Item{Val: 3}))
```

Contrary to [`Smuggle`]({{< ref "Smuggle" >}}) operator and its fields-path feature, only
public fields can be followed, as private ones are never (un)marshalled.

There is no JSONHas nor JSONHasnt operators to only check a [`JSON`]({{< ref "JSON" >}})
pointer exists or not, but they can easily be emulated:

```go
JSONHas := func(pointer string) td.TestDeep {
  return td.JSONPointer(pointer, td.Ignore())
}

JSONHasnt := func(pointer string) td.TestDeep {
  return td.Not(td.JSONPointer(pointer, td.Ignore()))
}
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method always returns `nil` as the expected type cannot be
guessed from a [`JSON`]({{< ref "JSON" >}}) pointer.


> See also [<i class='fas fa-book'></i> JSONPointer godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#JSONPointer).

### Examples

{{%expand "Rfc6901 example" %}}```go
	t := &testing.T{}

	got := json.RawMessage(`
{
   "foo":  ["bar", "baz"],
   "":     0,
   "a/b":  1,
   "c%d":  2,
   "e^f":  3,
   "g|h":  4,
   "i\\j": 5,
   "k\"l": 6,
   " ":    7,
   "m~n":  8
}`)

	expected := map[string]interface{}{
		"foo": []interface{}{"bar", "baz"},
		"":    0,
		"a/b": 1,
		"c%d": 2,
		"e^f": 3,
		"g|h": 4,
		`i\j`: 5,
		`k"l`: 6,
		" ":   7,
		"m~n": 8,
	}
	ok := td.Cmp(t, got, td.JSONPointer("", expected))
	fmt.Println("Empty JSON pointer means all:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/foo`, []interface{}{"bar", "baz"}))
	fmt.Println("Extract `foo` key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/foo/0`, "bar"))
	fmt.Println("First item of `foo` key slice:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/`, 0))
	fmt.Println("Empty key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/a~1b`, 1))
	fmt.Println("Slash has to be escaped using `~1`:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/c%d`, 2))
	fmt.Println("% in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/e^f`, 3))
	fmt.Println("^ in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/g|h`, 4))
	fmt.Println("| in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/i\j`, 5))
	fmt.Println("Backslash in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/k"l`, 6))
	fmt.Println("Double-quote in key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/ `, 7))
	fmt.Println("Space key:", ok)

	ok = td.Cmp(t, got, td.JSONPointer(`/m~0n`, 8))
	fmt.Println("Tilde has to be escaped using `~0`:", ok)

	// Output:
	// Empty JSON pointer means all: true
	// Extract `foo` key: true
	// First item of `foo` key slice: true
	// Empty key: true
	// Slash has to be escaped using `~1`: true
	// % in key: true
	// ^ in key: true
	// | in key: true
	// Backslash in key: true
	// Double-quote in key: true
	// Space key: true
	// Tilde has to be escaped using `~0`: true

```{{% /expand%}}
{{%expand "Struct example" %}}```go
	t := &testing.T{}

	// Without json tags, encoding/json uses public fields name
	type Item struct {
		Name  string
		Value int64
		Next  *Item
	}

	got := Item{
		Name:  "first",
		Value: 1,
		Next: &Item{
			Name:  "second",
			Value: 2,
			Next: &Item{
				Name:  "third",
				Value: 3,
			},
		},
	}

	ok := td.Cmp(t, got, td.JSONPointer("/Next/Next/Name", "third"))
	fmt.Println("3rd item name is `third`:", ok)

	ok = td.Cmp(t, got, td.JSONPointer("/Next/Next/Value", td.Gte(int64(3))))
	fmt.Println("3rd item value is greater or equal than 3:", ok)

	ok = td.Cmp(t, got,
		td.JSONPointer("/Next",
			td.JSONPointer("/Next",
				td.JSONPointer("/Value", td.Gte(int64(3))))))
	fmt.Println("3rd item value is still greater or equal than 3:", ok)

	ok = td.Cmp(t, got, td.JSONPointer("/Next/Next/Next/Name", td.Ignore()))
	fmt.Println("4th item exists and has a name:", ok)

	// Struct comparison work with or without pointer: &Item{…} works too
	ok = td.Cmp(t, got, td.JSONPointer("/Next/Next", Item{
		Name:  "third",
		Value: 3,
	}))
	fmt.Println("3rd item full comparison:", ok)

	// Output:
	// 3rd item name is `third`: true
	// 3rd item value is greater or equal than 3: true
	// 3rd item value is still greater or equal than 3: true
	// 4th item exists and has a name: false
	// 3rd item full comparison: true

```{{% /expand%}}
{{%expand "Has_hasnt example" %}}```go
	t := &testing.T{}

	got := json.RawMessage(`
{
  "name": "Bob",
  "age": 42,
  "children": [
    {
      "name": "Alice",
      "age": 16
    },
    {
      "name": "Britt",
      "age": 21,
      "children": [
        {
          "name": "John",
          "age": 1
        }
      ]
    }
  ]
}`)

	// Has Bob some children?
	ok := td.Cmp(t, got, td.JSONPointer("/children", td.Len(td.Gt(0))))
	fmt.Println("Bob has at least one child:", ok)

	// But checking "children" exists is enough here
	ok = td.Cmp(t, got, td.JSONPointer("/children/0/children", td.Ignore()))
	fmt.Println("Alice has children:", ok)

	ok = td.Cmp(t, got, td.JSONPointer("/children/1/children", td.Ignore()))
	fmt.Println("Britt has children:", ok)

	// The reverse can be checked too
	ok = td.Cmp(t, got, td.Not(td.JSONPointer("/children/0/children", td.Ignore())))
	fmt.Println("Alice hasn't children:", ok)

	ok = td.Cmp(t, got, td.Not(td.JSONPointer("/children/1/children", td.Ignore())))
	fmt.Println("Britt hasn't children:", ok)

	// Output:
	// Bob has at least one child: true
	// Alice has children: false
	// Britt has children: true
	// Alice hasn't children: true
	// Britt hasn't children: false

```{{% /expand%}}
## CmpJSONPointer shortcut

```go
func CmpJSONPointer(t TestingT, got interface{}, pointer string, expectedValue interface{}, args ...interface{}) bool
```

CmpJSONPointer is a shortcut for:

```go
td.Cmp(t, got, td.JSONPointer(pointer, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpJSONPointer godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpJSONPointer).

### Examples

{{%expand "Rfc6901 example" %}}```go
	t := &testing.T{}

	got := json.RawMessage(`
{
   "foo":  ["bar", "baz"],
   "":     0,
   "a/b":  1,
   "c%d":  2,
   "e^f":  3,
   "g|h":  4,
   "i\\j": 5,
   "k\"l": 6,
   " ":    7,
   "m~n":  8
}`)

	expected := map[string]interface{}{
		"foo": []interface{}{"bar", "baz"},
		"":    0,
		"a/b": 1,
		"c%d": 2,
		"e^f": 3,
		"g|h": 4,
		`i\j`: 5,
		`k"l`: 6,
		" ":   7,
		"m~n": 8,
	}
	ok := td.CmpJSONPointer(t, got, "", expected)
	fmt.Println("Empty JSON pointer means all:", ok)

	ok = td.CmpJSONPointer(t, got, `/foo`, []interface{}{"bar", "baz"})
	fmt.Println("Extract `foo` key:", ok)

	ok = td.CmpJSONPointer(t, got, `/foo/0`, "bar")
	fmt.Println("First item of `foo` key slice:", ok)

	ok = td.CmpJSONPointer(t, got, `/`, 0)
	fmt.Println("Empty key:", ok)

	ok = td.CmpJSONPointer(t, got, `/a~1b`, 1)
	fmt.Println("Slash has to be escaped using `~1`:", ok)

	ok = td.CmpJSONPointer(t, got, `/c%d`, 2)
	fmt.Println("% in key:", ok)

	ok = td.CmpJSONPointer(t, got, `/e^f`, 3)
	fmt.Println("^ in key:", ok)

	ok = td.CmpJSONPointer(t, got, `/g|h`, 4)
	fmt.Println("| in key:", ok)

	ok = td.CmpJSONPointer(t, got, `/i\j`, 5)
	fmt.Println("Backslash in key:", ok)

	ok = td.CmpJSONPointer(t, got, `/k"l`, 6)
	fmt.Println("Double-quote in key:", ok)

	ok = td.CmpJSONPointer(t, got, `/ `, 7)
	fmt.Println("Space key:", ok)

	ok = td.CmpJSONPointer(t, got, `/m~0n`, 8)
	fmt.Println("Tilde has to be escaped using `~0`:", ok)

	// Output:
	// Empty JSON pointer means all: true
	// Extract `foo` key: true
	// First item of `foo` key slice: true
	// Empty key: true
	// Slash has to be escaped using `~1`: true
	// % in key: true
	// ^ in key: true
	// | in key: true
	// Backslash in key: true
	// Double-quote in key: true
	// Space key: true
	// Tilde has to be escaped using `~0`: true

```{{% /expand%}}
{{%expand "Struct example" %}}```go
	t := &testing.T{}

	// Without json tags, encoding/json uses public fields name
	type Item struct {
		Name  string
		Value int64
		Next  *Item
	}

	got := Item{
		Name:  "first",
		Value: 1,
		Next: &Item{
			Name:  "second",
			Value: 2,
			Next: &Item{
				Name:  "third",
				Value: 3,
			},
		},
	}

	ok := td.CmpJSONPointer(t, got, "/Next/Next/Name", "third")
	fmt.Println("3rd item name is `third`:", ok)

	ok = td.CmpJSONPointer(t, got, "/Next/Next/Value", td.Gte(int64(3)))
	fmt.Println("3rd item value is greater or equal than 3:", ok)

	ok = td.CmpJSONPointer(t, got, "/Next", td.JSONPointer("/Next",
		td.JSONPointer("/Value", td.Gte(int64(3)))))
	fmt.Println("3rd item value is still greater or equal than 3:", ok)

	ok = td.CmpJSONPointer(t, got, "/Next/Next/Next/Name", td.Ignore())
	fmt.Println("4th item exists and has a name:", ok)

	// Struct comparison work with or without pointer: &Item{…} works too
	ok = td.CmpJSONPointer(t, got, "/Next/Next", Item{
		Name:  "third",
		Value: 3,
	})
	fmt.Println("3rd item full comparison:", ok)

	// Output:
	// 3rd item name is `third`: true
	// 3rd item value is greater or equal than 3: true
	// 3rd item value is still greater or equal than 3: true
	// 4th item exists and has a name: false
	// 3rd item full comparison: true

```{{% /expand%}}
{{%expand "Has_hasnt example" %}}```go
	t := &testing.T{}

	got := json.RawMessage(`
{
  "name": "Bob",
  "age": 42,
  "children": [
    {
      "name": "Alice",
      "age": 16
    },
    {
      "name": "Britt",
      "age": 21,
      "children": [
        {
          "name": "John",
          "age": 1
        }
      ]
    }
  ]
}`)

	// Has Bob some children?
	ok := td.CmpJSONPointer(t, got, "/children", td.Len(td.Gt(0)))
	fmt.Println("Bob has at least one child:", ok)

	// But checking "children" exists is enough here
	ok = td.CmpJSONPointer(t, got, "/children/0/children", td.Ignore())
	fmt.Println("Alice has children:", ok)

	ok = td.CmpJSONPointer(t, got, "/children/1/children", td.Ignore())
	fmt.Println("Britt has children:", ok)

	// The reverse can be checked too
	ok = td.Cmp(t, got, td.Not(td.JSONPointer("/children/0/children", td.Ignore())))
	fmt.Println("Alice hasn't children:", ok)

	ok = td.Cmp(t, got, td.Not(td.JSONPointer("/children/1/children", td.Ignore())))
	fmt.Println("Britt hasn't children:", ok)

	// Output:
	// Bob has at least one child: true
	// Alice has children: false
	// Britt has children: true
	// Alice hasn't children: true
	// Britt hasn't children: false

```{{% /expand%}}
## T.JSONPointer shortcut

```go
func (t *T) JSONPointer(got interface{}, pointer string, expectedValue interface{}, args ...interface{}) bool
```

JSONPointer is a shortcut for:

```go
t.Cmp(got, td.JSONPointer(pointer, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.JSONPointer godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.JSONPointer).

### Examples

{{%expand "Rfc6901 example" %}}```go
	t := td.NewT(&testing.T{})

	got := json.RawMessage(`
{
   "foo":  ["bar", "baz"],
   "":     0,
   "a/b":  1,
   "c%d":  2,
   "e^f":  3,
   "g|h":  4,
   "i\\j": 5,
   "k\"l": 6,
   " ":    7,
   "m~n":  8
}`)

	expected := map[string]interface{}{
		"foo": []interface{}{"bar", "baz"},
		"":    0,
		"a/b": 1,
		"c%d": 2,
		"e^f": 3,
		"g|h": 4,
		`i\j`: 5,
		`k"l`: 6,
		" ":   7,
		"m~n": 8,
	}
	ok := t.JSONPointer(got, "", expected)
	fmt.Println("Empty JSON pointer means all:", ok)

	ok = t.JSONPointer(got, `/foo`, []interface{}{"bar", "baz"})
	fmt.Println("Extract `foo` key:", ok)

	ok = t.JSONPointer(got, `/foo/0`, "bar")
	fmt.Println("First item of `foo` key slice:", ok)

	ok = t.JSONPointer(got, `/`, 0)
	fmt.Println("Empty key:", ok)

	ok = t.JSONPointer(got, `/a~1b`, 1)
	fmt.Println("Slash has to be escaped using `~1`:", ok)

	ok = t.JSONPointer(got, `/c%d`, 2)
	fmt.Println("% in key:", ok)

	ok = t.JSONPointer(got, `/e^f`, 3)
	fmt.Println("^ in key:", ok)

	ok = t.JSONPointer(got, `/g|h`, 4)
	fmt.Println("| in key:", ok)

	ok = t.JSONPointer(got, `/i\j`, 5)
	fmt.Println("Backslash in key:", ok)

	ok = t.JSONPointer(got, `/k"l`, 6)
	fmt.Println("Double-quote in key:", ok)

	ok = t.JSONPointer(got, `/ `, 7)
	fmt.Println("Space key:", ok)

	ok = t.JSONPointer(got, `/m~0n`, 8)
	fmt.Println("Tilde has to be escaped using `~0`:", ok)

	// Output:
	// Empty JSON pointer means all: true
	// Extract `foo` key: true
	// First item of `foo` key slice: true
	// Empty key: true
	// Slash has to be escaped using `~1`: true
	// % in key: true
	// ^ in key: true
	// | in key: true
	// Backslash in key: true
	// Double-quote in key: true
	// Space key: true
	// Tilde has to be escaped using `~0`: true

```{{% /expand%}}
{{%expand "Struct example" %}}```go
	t := td.NewT(&testing.T{})

	// Without json tags, encoding/json uses public fields name
	type Item struct {
		Name  string
		Value int64
		Next  *Item
	}

	got := Item{
		Name:  "first",
		Value: 1,
		Next: &Item{
			Name:  "second",
			Value: 2,
			Next: &Item{
				Name:  "third",
				Value: 3,
			},
		},
	}

	ok := t.JSONPointer(got, "/Next/Next/Name", "third")
	fmt.Println("3rd item name is `third`:", ok)

	ok = t.JSONPointer(got, "/Next/Next/Value", td.Gte(int64(3)))
	fmt.Println("3rd item value is greater or equal than 3:", ok)

	ok = t.JSONPointer(got, "/Next", td.JSONPointer("/Next",
		td.JSONPointer("/Value", td.Gte(int64(3)))))
	fmt.Println("3rd item value is still greater or equal than 3:", ok)

	ok = t.JSONPointer(got, "/Next/Next/Next/Name", td.Ignore())
	fmt.Println("4th item exists and has a name:", ok)

	// Struct comparison work with or without pointer: &Item{…} works too
	ok = t.JSONPointer(got, "/Next/Next", Item{
		Name:  "third",
		Value: 3,
	})
	fmt.Println("3rd item full comparison:", ok)

	// Output:
	// 3rd item name is `third`: true
	// 3rd item value is greater or equal than 3: true
	// 3rd item value is still greater or equal than 3: true
	// 4th item exists and has a name: false
	// 3rd item full comparison: true

```{{% /expand%}}
{{%expand "Has_hasnt example" %}}```go
	t := td.NewT(&testing.T{})

	got := json.RawMessage(`
{
  "name": "Bob",
  "age": 42,
  "children": [
    {
      "name": "Alice",
      "age": 16
    },
    {
      "name": "Britt",
      "age": 21,
      "children": [
        {
          "name": "John",
          "age": 1
        }
      ]
    }
  ]
}`)

	// Has Bob some children?
	ok := t.JSONPointer(got, "/children", td.Len(td.Gt(0)))
	fmt.Println("Bob has at least one child:", ok)

	// But checking "children" exists is enough here
	ok = t.JSONPointer(got, "/children/0/children", td.Ignore())
	fmt.Println("Alice has children:", ok)

	ok = t.JSONPointer(got, "/children/1/children", td.Ignore())
	fmt.Println("Britt has children:", ok)

	// The reverse can be checked too
	ok = t.Cmp(got, td.Not(td.JSONPointer("/children/0/children", td.Ignore())))
	fmt.Println("Alice hasn't children:", ok)

	ok = t.Cmp(got, td.Not(td.JSONPointer("/children/1/children", td.Ignore())))
	fmt.Println("Britt hasn't children:", ok)

	// Output:
	// Bob has at least one child: true
	// Alice has children: false
	// Britt has children: true
	// Alice hasn't children: true
	// Britt hasn't children: false

```{{% /expand%}}
