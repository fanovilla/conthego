# Concordion On The Go

A lightweight, opinionated implementation of a Specification By Example (SBE) framework in golang.
Inspired by [concordion.org](https://concordion.org).

## Usage

The package can be imported as `github.com/fanovilla/conthego/conthego`.

See examples in the `examples` directory.

Sample specification `AssertEquals.md`
```markdown
[Hello World!](- "?GetGreeting()")
```

Sample fixture `AssertEquals_test.go`
```go
func TestAssertEquals(t *testing.T) {
    conthego.RunSpec(t, &AssertEqualsFixture{})
}

type AssertEqualsFixture struct {
}

func (f *AssertEqualsFixture) GetGreeting() string {
    return "Hello World!"
}
```

## Supported Commands

* [assert equals](examples/AssertEquals.md)

```
set
[World](- "var1")

echo
[ ](- "$var1")
[ ](- "$var1.prop2")
[ ](- "$Hello()")
[ ](- "$Hello(var1)")

exec
[ ](- "Hello()")
[Blah](- "var1=Hello(TEXT)")
[ ](- "var2=Hello(var1)")

assertTrue or assertEquals
[World](- "?Hello()")
[World](- "?var1")
[World](- "?var1.prop")

directives
[ ](- "!expectedtofail")

```

See [TODOs](TODO.md)

See table examples in the `examples` directory.


## Notes

### specifications
* commands processed via depth-first traversal (e.g. vars can only be used after setting, never before)
* strict markdown tables: at most a single command per column; column commands processed left to right
* empty link text currently require a single space (e.g. use `[ ]` not `[]`)
* echo command uses a `$` prefix (not `c:echo=`)
* assert command uses a `?` prefix (not `?=`)
* set command uses no prefix (no `#`)
* params for method calls uses no prefix (no `#`)
* params for method calls bound as strings; please do conversion as per need inside fixture methods

### fixtures
* fixture methods required to return a value, current limitation
* struct properties to assert on must be exported (dependent on json marshal behaviour)