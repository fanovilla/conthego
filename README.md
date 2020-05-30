# Concordion On The Go

Conthego. A lightweight implementation of a Specification By Example (SBE) framework in [go](https://golang.org/).
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
import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestAssertEquals(t *testing.T) {
	conthego.RunSpec(t, &AssertEqualsFixture{})
}

type AssertEqualsFixture struct {
}

func (f *AssertEqualsFixture) GetGreeting() string {
	return "Hello World!"
}
```
Sample output:

![output](./examples/images/AssertEquals.png)

## Supported Commands

* [assert-equals](examples/AssertEquals.md)
* [assert-true](examples/AssertTrue.md)
* [set](examples/SetAndEcho.md)
* [echo](examples/SetAndEcho.md)

```
set
[World](- "var1")

echo
[ ](- "$var1")
[ ](- "$var1.prop2")
[ ](- "$Hello()")
[ ](- "$Hello(var1)")
[ ](- "$$rawHtmlToEmbed")


exec
[ ](- "Hello()")
[Blah](- "var1=Hello(TEXT)")
[ ](- "var2=Hello(var1)")

assert-true or assert-equals
[World](- "?Hello()")
[World](- "?var1")
[World](- "?var1.prop")

directives
[ ](- "!ExpectedToFail")

```

See [TODOs](TODO.md)

See table examples in the `examples` directory.


## Notes

### specifications
* commands processed via depth-first traversal (e.g. vars can only be used after setting, never before)
* column commands processed left to right
* empty link text currently require a single space (e.g. use `[ ]` not `[]`)
* echo command uses a `$` prefix (not `c:echo=`)
* assert command uses a `?` prefix (not `?=`)
* set command uses no prefix (no `#`)
* params for method calls uses no prefix (no `#`)
* params for method calls bound as strings; please convert as per need in fixture methods

### fixtures
* fixture methods required to return a value, current limitation
* struct properties to assert on must be exported (dependent on json marshal behaviour)