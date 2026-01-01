package main

import (
	"fmt"
	"os"

	"github.com/wozaki/alfred-cursor-launcher/internal"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		runList()
	case "open":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: URI not specified\n")
			os.Exit(1)
		}
		runOpen(os.Args[2])
	case "version":
		fmt.Println(version)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s list              - List recent projects\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s open <uri>        - Open a project\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s version           - Show version\n", os.Args[0])
}

func runList() {
	sf := internal.NewScriptFilter()

	store, err := internal.NewProjectStore()
	if err != nil {
		sf.AddErrorItem(err.Error())
		sf.Print()
		os.Exit(1)
	}

	projects, err := store.List()
	if err != nil {
		sf.AddErrorItem(err.Error())
		sf.Print()
		os.Exit(1)
	}

	for _, project := range projects {
		sf.AddItem(project.ToAlfredItem())
	}

	sf.Print()
}

func runOpen(uri string) {
	opener := internal.NewOpener()
	if err := opener.Open(uri); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
