# Conventional Commit Message Parser

Parse conventional commit message with this parser.

## Usage

```go
package main

import (
    "fmt"
    convcommit "github.com/chonla/conv-commit"
)

func main() {
    commitMessage := "feat(api): This is a commit message"
    result, err := convcommit.Parse(commitMessage)

    if err != nil {
        fmt.Println(result)
    } else {
        fmt.Println(err)
    }
}
```

## Conventional commit?

See [Conventional Commit](https://www.conventionalcommits.org/)

## License

[MIT](./LICENSE)