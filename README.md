# Concordion On The Go

A lightweight, opinionated implementation of a Specification By Example (SBE) framework in golang.
Inspired by [concordion.org](https://concordion.org).

See examples in the `examples` directory.

## Supported Commands

```
set
[World](- "var1")

echo
[](- "$var1")
[](- "$var1.prop2")
[](- "$Hello()")
[](- "$Hello(var1)")

exec
[](- "Hello()")
[](- "var1=Hello(TEXT)")
[](- "var2=Hello(var1)")
[](- "var1=Hello('literal')")

isTrue or assertEquals
[World](- "?Hello()")
[World](- "?var1")
[World](- "?var1.prop")
```

