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

Searching for [arr](- "results=SearchString(TEXT)") will [return:](- "!VerifyRows")

* [george.harrison](- "?results")
* ringo.starr

