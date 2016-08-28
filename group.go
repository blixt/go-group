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
	if arguments[0] != g.name {
		panic("group: unexpected argument")
	}
	if err := g.Flag.Parse(arguments[1:]); err != nil {
		panic(err)
	}
	subarg := g.Flag.Arg(0)
	if sub, ok := g.subs[subarg]; ok {
		// Ignore everything before the first instance of the subargument.
		for i, a := range arguments {
			if i > 0 && a == subarg {
				arguments = arguments[i:]
				break
			}
		}
		g.sub = sub
		arguments[0] = subarg
		sub.Parse(arguments)
		return sub
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
	return CommandLine.Parse(os.Args)
}

func Sub(name string) *Group {
	return CommandLine.Sub(name)
}

func Subs() []string {
	return CommandLine.Subs()
}
