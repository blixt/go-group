group package
=============

[![GoDoc](https://godoc.org/github.com/blixt/go-group?status.svg)](https://godoc.org/github.com/blixt/go-group)

The group package makes it easier to handle multiple command line groups
with the `flags` package.


Example #1
----------

```go
package main

import "fmt"
import "github.com/blixt/go-group"

// group.Flag lets you create global flags (same as flag package).
var verbose = group.Flag.Bool("v", false, "Output more")
// Keep references to sub-commands to parse arguments and flags later.
var help = group.Sub("help")
var clone = group.Sub("clone")
// Each sub-command has its own FlagSet for specific flags.
var branch = clone.Flag.String("branch", "master", "The branch to clone")

func main() {
  // Parse the sub-command used (and any flags along the way).
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


Example #2
----------

It's also possible to create deeply nested sub-commands.

```go
package main

import "fmt"
import "github.com/blixt/go-group"

var preview = group.Sub("preview")
// When the "app" group is stable, change `preview` to `group` below.
var app = preview.Sub("app")
var deploy = app.Sub("deploy")

func main() {
  switch group.Parse() {
  case deploy:
    fmt.Println("Deploying...")
  default:
    fmt.Println("Unsupported command.")
  }
}
```