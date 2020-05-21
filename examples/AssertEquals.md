# assert-equals command


## overview
Create a spec file AssertEquals.md containing:

```markdown
[Hello World!](- "?GetGreeting()")
```


Create a fixture file AssertEquals_test.go that implements the **GetGreeting()** function:

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

When run, the output specification will show:

![output](./images/AssertEquals.png)

## examples

[Hello World!](- "?GetGreeting()")

Given my name is [Carl](- "name")

[Hello Carl!](- "?GetPersonalisedGreeting(name)")

Given your names are [Ryan](- "ryan") and [Liam](- "liam")

[Hello Ryan and Liam!](- "?GetMultipleGreeting(ryan, liam)")