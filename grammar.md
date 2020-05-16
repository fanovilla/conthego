
set
[World](- "var1")

echo
[](- "$var1")
[](- "$var1.prop2")

exec
[](- "Hello()")
[](- "var1=Hello(TEXT)")
[](- "var2=Hello(var1)")
[](- "var1=Hello('literal')")

assert equals or isTrue
[World](- "?Hello()")
[World](- "?var1")
