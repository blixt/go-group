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

// Keep references to sub-commands to parse arguments and flags later.
var help = group.Sub("help")
var clone = group.Sub("clone")

// Each sub-command has its own FlagSet for specific flags.
var branch = clone.Flag.String("branch", "master", "The branch to clone")

func main() {
	// Parse the sub-command used (and any flags along the way).
	switch group.Parse() {
	case help:
		fmt.Println("This is some help for this tool.")
	case clone:
		repo := clone.Flag.Arg(0)
		if repo == "" {
			fmt.Println("Invalid repository!")
			break
		}
		fmt.Printf("Cloning %s (branch %s)...\n", repo, *branch)
	default:
		fmt.Println("Unrecognized group. Choose one of:", group.Subs())
	}
}
```

### Result

```
$ ./example clone -branch dev git.example.com:myrepo.git
Cloning git.example.com:myrepo.git (branch dev)...

$ ./example help
This is some help for this tool.
```


Example #2
----------

This example shows how to make a deeply nested command and global flags.

```go
package main

import "fmt"
import "github.com/blixt/go-group"

var preview = group.Sub("preview")
var app = preview.Sub("app")
var deploy = app.Sub("deploy")

// Global flag (before any of the command groups):
var verbose = group.Flag.Bool("v", false, "Output more")

func main() {
	switch group.Parse() {
	case deploy:
		fmt.Println("Deploying...")
	default:
		fmt.Println("Unsupported command.")
	}

	if *verbose {
		fmt.Println("And here's a bunch of extra output because you specified -v.")
	}
}
```

(When the `app` command group is no longer in preview, you would just
need to change `preview.Sub("app")` to `group.Sub("app")`.)

### Result

```
$ ./example -v preview app deploy
Deploying...
And here's a bunch of extra output because you specified -v.

$ ./example somethingelse
Unsupported command.
```
