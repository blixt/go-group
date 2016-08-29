package group

import (
	"os"
	"testing"
)

func TestDeeplyNested(t *testing.T) {
	root := NewGroup("gcloud")
	preview := root.Sub("preview")
	app := preview.Sub("app")
	deploy := app.Sub("deploy")

	args := []string{"gcloud", "preview", "app", "deploy"}
	g := root.Parse(args[1:])
	if g != deploy {
		t.Errorf("unexpected group %v", g)
	}
}

func TestGlobal(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"git", "-v", "clone", "-b", "dev", "git.example.com:repo.git"}

	verbose := Flag.Bool("v", false, "verbose")
	clone := Sub("clone")
	branch := clone.Flag.String("b", "master", "branch")

	g := Parse()
	if !*verbose {
		t.Errorf("unexpected verbose value %v", *verbose)
	}
	if g != clone {
		t.Errorf("unexpected group %v", g)
	}
	if *branch != "dev" {
		t.Errorf("unexpected branch value %v", *branch)
	}
}
