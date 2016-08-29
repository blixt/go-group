/*
Package group implements parsing for command-line groups.

This package is a light-weight wrapper for the built-in flag package, adding
support for branching based on specific arguments.

This declares command group clone.
	import "github.com/blixt/go-group"
	var clone = group.Sub("clone")
If you like, you can bind flags that will only be parsed
if the group is specified in the arguments list.
	var branch = clone.Flag.String("branch", "master", "branch to clone")
Multiple levels of grouping is also supported.
	var codereview = group.Sub("codereview")
	var submit = codereview.Sub("submit")

After all groups and flags are defined, call
	group.Parse()
to parse the command line into the defined groups and flags.

Groups may then be used directly to access parsed flags and arguments. For
convenience, the Parse function returns the deepest matched group, which works
well with the switch statement:
	switch group.Parse() {
	case clone:
		repo := clone.Flag.Arg(0)
		fmt.Println("Cloning", repo, "on branch", *branch)
	case submit:
		fmt.Println("Submitting code review...")
	default:
		fmt.Println("unsupported command")
	}
*/
package group

import (
	"flag"
	"os"
)

// CommandLine is the root command, parsed from os.Args.
// The top-level functions such as ActiveSub, Parse, and
// so on are wrappers for the methods of CommandLine.
var CommandLine = NewGroup(os.Args[0])

// Flag is the FlagSet for the root command, used to specify global flags.
var Flag = CommandLine.Flag

// NewGroup returns a new command group with the specified name.
func NewGroup(name string) *Group {
	return &Group{
		Flag: flag.NewFlagSet(name, flag.ExitOnError),
		name: name,
		subs: make(map[string]*Group),
	}
}

// A Group represents a command group.
type Group struct {
	Flag *flag.FlagSet
	name string
	sub  *Group
	subs map[string]*Group
}

// ActiveSub returns the sub-command under this group which was specified
// in the argument list given to Parse, or nil if none was specified. This
// method will always return nil until Parse has been called.
func (g *Group) ActiveSub() *Group {
	return g.sub
}

// Parse parses command groups and flags from the argument
// list, which should not include the group name.
func (g *Group) Parse(arguments []string) *Group {
	if err := g.Flag.Parse(arguments); err != nil {
		panic(err)
	}
	subarg := g.Flag.Arg(0)
	if sub, ok := g.subs[subarg]; ok {
		// Ignore everything before the first instance of the subargument.
		for i, a := range arguments {
			if a == subarg {
				arguments = arguments[i+1:]
				break
			}
		}
		g.sub = sub
		return sub.Parse(arguments)
	}
	return g
}

// Sub defines a new sub-command group under this group. The return value
// is another Group which can itself have additional sub-commands.
func (g *Group) Sub(name string) *Group {
	if g.subs[name] != nil {
		panic("group: duplicate")
	}
	sg := NewGroup(name)
	g.subs[name] = sg
	return sg
}

// Subs returns a slice of commands available directly under this group.
func (g *Group) Subs() []string {
	subs := make([]string, 0, len(g.subs))
	for s := range g.subs {
		subs = append(subs, s)
	}
	return subs
}

// ActiveSub returns the sub-command that was specified on the
// command line, or nil if none was specified. This function will
// always return nil until Parse has been called.
func ActiveSub() *Group {
	return CommandLine.ActiveSub()
}

// Parse parses the command-line arguments from os.Args[1:].
// Must be called after all groups and flags are defined and
// before they are accessed by the program.
func Parse() *Group {
	return CommandLine.Parse(os.Args[1:])
}

// Sub defines a new sub-command group. The return value is
// a Group which can itself have additional sub-commands.
func Sub(name string) *Group {
	return CommandLine.Sub(name)
}

// Subs returns a slice of available sub-commands
// at the top level of the command.
func Subs() []string {
	return CommandLine.Subs()
}
