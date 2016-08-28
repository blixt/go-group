group package
=============

The group package makes it easier to handle multiple command line groups
with the `flags` package.


Example
-------

```go
package main

import "fmt"
import "github.com/blixt/go-group"

func main() {
  verbose := group.Flag.Bool("v", false, "Output more")

  help := group.Sub("help")

  clone := group.Sub("clone")
  branch := clone.Flag.String("branch", "master", "The branch to clone")
  clone.Flag.StringVar(branch, "b", "master", "The branch to clone (shorthand)")

  switch group.Parse() {
  case help:
    category := help.Flag.Arg(0)
    if category != "" {
      fmt.Println("TODO: Write some help about", category)
    } else {
      fmt.Println("TODO: Write some general help for this tool")
    }
  case clone:
    repo := clone.Flag.Arg(0)
    if repo == "" {
      fmt.Println("Invalid repository!")
      return
    }
    fmt.Printf("Cloning %s (branch %s)...\n", repo, *branch)
  default:
    fmt.Println("Unrecognized group. Choose one of:", group.Subs())
  }

  if *verbose {
    fmt.Println("And here's a bunch of extra output because you specified -v.")
  }
}
```