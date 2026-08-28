package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/chenpenghai/ai-project-framework/internal/graph"
	"github.com/chenpenghai/ai-project-framework/internal/scan"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "scan":
		if err := runScan(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "apf:", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func runScan(args []string) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	jsonOut := fs.Bool("json", false, "emit the full structure snapshot as JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}

	snapshot, err := (scan.Scanner{}).Scan(root)
	if err != nil {
		return err
	}
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(snapshot)
	}
	printSummary(snapshot)
	return nil
}

func printSummary(s graph.Snapshot) {
	counts := map[graph.NodeKind]int{}
	for _, n := range s.Nodes {
		counts[n.Kind]++
	}
	fmt.Println("AI Project Framework — structure scan")
	fmt.Println("root:", s.Root)
	fmt.Printf("git: %t (%d changed files)\n", s.Git.Available, len(s.Git.ChangedFiles))
	fmt.Printf("projects: %d\n", counts[graph.NodeProject])
	fmt.Printf("modules: %d\n", counts[graph.NodeModule])
	fmt.Printf("files: %d\n", counts[graph.NodeFile])

	if counts[graph.NodeProject] > 0 {
		fmt.Println("\nprojects:")
		for _, n := range s.Nodes {
			if n.Kind == graph.NodeProject {
				fmt.Printf("  - %s [%s] %s\n", n.Name, n.Metadata["ecosystem"], n.Path)
			}
		}
	}
	if counts[graph.NodeModule] > 0 {
		fmt.Println("\nexplicit modules:")
		for _, n := range s.Nodes {
			if n.Kind == graph.NodeModule {
				fmt.Printf("  - %s %s\n", n.Name, n.Path)
			}
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: apf scan [--json] [repository]")
}
