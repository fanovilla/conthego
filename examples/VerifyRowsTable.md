# Partial Matches

Username searches return partial matches, i.e. all usernames containing the search string are returned.

### Example

Given these users:

| [Username][setup] |
|-------------------|
| john.lennon       |
| ringo.starr       |
| george.harrison   |
| paul.mccartney    |

[setup]: - "SetUpUser(TEXT)"

Searching for [arr](- "results=BreakDownNamesWith(TEXT)") and having names broken down will return:

| [ ][t1] [First][first] | [Last][last] |
|-------------------------|---------------|
| ringo                   |               |
| george                  |               |

[t1]: - "result:results"
[first]: - "?result.First"
[last]:  - "$result.Last"


