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
	conthego.RunSpec(conthego.NewFixture(t, &FixtureAssertEquals{}))
}

type FixtureAssertEquals struct {
}

func (f FixtureAssertEquals) GetGreeting() string {
	return "Hello World!"
}
```

## Supported Commands

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
[ ](- "var1=Hello('literal')")

isTrue or assertEquals
[World](- "?Hello()")
[World](- "?var1")
[World](- "?var1.prop")

directives
[ ](- "!expectedtofail")

```

## Tips

* struct properties to assert on must be exported (dependent on json marshal behaviour)