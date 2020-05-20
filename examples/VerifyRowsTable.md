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

First and last names [broken down](- "res=BreakDownNames()") will [return:](- "!VerifyRows")

| [First](- "?res.First") | [Last](- "$res.Last") |
|-------------------|-----|
| john       | lennon |
| ringo       | starr |
| george   | harrison |
| paul    | mccartney |


