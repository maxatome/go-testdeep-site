---
title: "Catch"
weight: 10
---

```go
func Catch(target, expectedValue any) TestDeep
```

Catch is a [smuggler operator]({{% ref "operators#smuggler-operators" %}}). It allows to copy data in *target* on
the fly before comparing it as usual against *expectedValue*.

*target* must be a non-`nil` pointer and data should be assignable to
its pointed type. If [`BeLax` config flag](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContextConfig.BeLax) is true or called under [`Lax`]({{% ref "Lax" %}})
(and so [`JSON`]({{% ref "JSON" %}})) operator, data should be convertible to its pointer
type.

```go
var id int64
if td.Cmp(t, CreateRecord("test"),
  td.JSON(`{"id": $1, "name": "test"}`, td.Catch(&id, td.NotZero()))) {
  t.Logf("Created record ID is %d", id)
}
```

It is really useful when used with [`JSON`]({{% ref "JSON" %}}) operator and/or [tdhttp](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp) helper.

```go
var id int64
ta := tdhttp.NewTestAPI(t, api.Handler).
  PostJSON("/item", `{"name":"foo"}`).
  CmpStatus(http.StatusCreated).
  CmpJSONBody(td.JSON(`{"id": $1, "name": "foo"}`, td.Catch(&id, td.Gt(0))))
if !ta.Failed() {
  t.Logf("Created record ID is %d", id)
}
```

If you need to only catch data without comparing it, use [`Ignore`]({{% ref "Ignore" %}})
operator as *expectedValue* as in:

```go
var id int64
if td.Cmp(t, CreateRecord("test"),
  td.JSON(`{"id": $1, "name": "test"}`, td.Catch(&id, td.Ignore()))) {
  t.Logf("Created record ID is %d", id)
}
```

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of *expectedValue*,
except if *expectedValue* is a [TestDeep operator]({{% ref "operators" %}}). In this case, it
delegates [`TypeBehind()`]({{% ref "operators#typebehind-method" %}}) to the operator, but if `nil` is returned by
this call, the dereferenced [`reflect.Type`](https://pkg.go.dev/reflect#Type) of *target* is returned.



> See also [<i class='fas fa-book'></i> Catch godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Catch).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	var age int
	ok := td.Cmp(t, got,
		td.JSON(`{"age":$1,"fullname":"Bob"}`,
			td.Catch(&age, td.Between(40, 45))))
	fmt.Println("check got age+fullname:", ok)
	fmt.Println("caught age:", age)

	// Output:
	// check got age+fullname: true
	// caught age: 42

```{{% /expand%}}
