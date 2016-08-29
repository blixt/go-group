package group

import (
	"flag"
	"os"
)

var CommandLine = NewGroup(os.Args[0])
var Flag = CommandLine.Flag

func NewGroup(name string) *Group {
	return &Group{
		Flag: flag.NewFlagSet(name, flag.ExitOnError),
		name: name,
		subs: make(map[string]*Group),
	}
}

type Group struct {
	Flag *flag.FlagSet
	name string
	sub  *Group
	subs map[string]*Group
}

func (g *Group) ActiveSub() *Group {
	if !g.Flag.Parsed() {
		panic("group: ActiveSub called before Parse")
	}
	return g.sub
}

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

func (g *Group) Sub(name string) *Group {
	if g.subs[name] != nil {
		panic("group: duplicate")
	}
	sg := NewGroup(name)
	g.subs[name] = sg
	return sg
}

func (g *Group) Subs() []string {
	subs := make([]string, 0, len(g.subs))
	for s := range g.subs {
		subs = append(subs, s)
	}
	return subs
}

func ActiveSub() *Group {
	return CommandLine.ActiveSub()
}

func Parse() *Group {
	return CommandLine.Parse(os.Args[1:])
}

func Sub(name string) *Group {
	return CommandLine.Sub(name)
}

func Subs() []string {
	return CommandLine.Subs()
}
