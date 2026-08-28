package scan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chenpenghai/ai-project-framework/internal/graph"
)

func TestScanDiscoversProjectsExplicitModulesAndContainment(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"name":"demo-root"}`)
	write(t, root, "src/order/MODULE.md", "---\nmodule: order\n---\n")
	write(t, root, "src/order/refund.js", "export function refund() {}\n")
	write(t, root, "services/worker/go.mod", "module example.com/worker\n\ngo 1.23\n")
	write(t, root, "services/worker/main.go", "package main\n")

	s, err := (Scanner{}).Scan(root)
	if err != nil {
		t.Fatal(err)
	}

	assertNode(t, s.Nodes, "project:.", graph.NodeProject, "demo-root")
	assertNode(t, s.Nodes, "project:services/worker", graph.NodeProject, "example.com/worker")
	assertNode(t, s.Nodes, "module:src/order", graph.NodeModule, "order")
	assertEdge(t, s.Edges, "project:.", "module:src/order", graph.EdgeContains)
	assertEdge(t, s.Edges, "module:src/order", "file:src/order/refund.js", graph.EdgeContains)
	assertEdge(t, s.Edges, "project:.", "project:services/worker", graph.EdgeContains)
	assertEdge(t, s.Edges, "project:services/worker", "file:services/worker/main.go", graph.EdgeContains)
}

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertNode(t *testing.T, nodes []graph.Node, id string, kind graph.NodeKind, name string) {
	t.Helper()
	for _, n := range nodes {
		if n.ID == id {
			if n.Kind != kind || n.Name != name {
				t.Fatalf("node %s = kind %s name %q, want %s %q", id, n.Kind, n.Name, kind, name)
			}
			return
		}
	}
	t.Fatalf("missing node %s", id)
}

func assertEdge(t *testing.T, edges []graph.Edge, from, to string, kind graph.EdgeKind) {
	t.Helper()
	for _, e := range edges {
		if e.From == from && e.To == to && e.Kind == kind {
			return
		}
	}
	t.Fatalf("missing edge %s -[%s]-> %s", from, kind, to)
}
