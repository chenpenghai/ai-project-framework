package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chenpenghai/ai-project-framework/internal/graph"
)

func TestScanDiscoversProjectsModulesContainmentAndDeclaredDependencies(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"name":"demo-root"}`)
	write(t, root, "go.work", "go 1.23\n")
	write(t, root, "src/order/MODULE.md", "---\nmodule: order\n---\n")
	write(t, root, "src/order/refund.js", "export function refund() {}\n")
	write(t, root, "src/billing/MODULE.md", "# Notes\nmodule: not-frontmatter\n")
	write(t, root, "src/billing/calc.js", "export function calc() {}\n")

	write(t, root, "shared/go.mod", "module example.com/shared\n\ngo 1.23\n")
	write(t, root, "shared/shared.go", "package shared\n")
	write(t, root, "services/worker/go.mod", "module example.com/worker\n\ngo 1.23\n\nrequire (\n\texample.com/shared v0.0.0\n)\n")
	write(t, root, "services/worker/main.go", "package main\n")
	write(t, root, "services/python/requirements.txt", "fastapi\n")
	write(t, root, "services/python/main.py", "print('ok')\n")

	write(t, root, "packages/ui/package.json", `{"name":"@demo/ui"}`)
	write(t, root, "packages/ui/index.js", "export const ui = true\n")
	write(t, root, "packages/app/package.json", `{"name":"@demo/app","dependencies":{"@demo/ui":"workspace:*"}}`)
	write(t, root, "packages/app/index.js", "export const app = true\n")

	s, err := (Scanner{}).Scan(root)
	if err != nil {
		t.Fatal(err)
	}

	rootProject := assertNode(t, s.Nodes, "project:.", graph.NodeProject, "demo-root")
	if rootProject.Confidence != graph.ConfidenceDeclared {
		t.Fatalf("root project confidence = %s, want declared", rootProject.Confidence)
	}
	if !strings.Contains(rootProject.Metadata["ecosystems"], "go-workspace") || !strings.Contains(rootProject.Metadata["ecosystems"], "node") {
		t.Fatalf("root project lost polyglot evidence: %#v", rootProject.Metadata)
	}

	assertNode(t, s.Nodes, "project:services/worker", graph.NodeProject, "example.com/worker")
	assertNode(t, s.Nodes, "project:shared", graph.NodeProject, "example.com/shared")
	pythonProject := assertNode(t, s.Nodes, "project:services/python", graph.NodeProject, "python")
	if pythonProject.Confidence != graph.ConfidenceInferredHigh {
		t.Fatalf("requirements-only project confidence = %s, want inferred_high", pythonProject.Confidence)
	}
	assertNode(t, s.Nodes, "module:src/order", graph.NodeModule, "order")
	assertNode(t, s.Nodes, "module:src/billing", graph.NodeModule, "billing")
	assertEdge(t, s.Edges, "project:.", "module:src/order", graph.EdgeContains)
	assertEdge(t, s.Edges, "module:src/order", "file:src/order/refund.js", graph.EdgeContains)
	assertEdge(t, s.Edges, "project:.", "project:services/worker", graph.EdgeContains)
	assertEdge(t, s.Edges, "project:services/worker", "file:services/worker/main.go", graph.EdgeContains)
	assertEdge(t, s.Edges, "project:services/worker", "project:shared", graph.EdgeDependsOn)
	assertEdge(t, s.Edges, "project:packages/app", "project:packages/ui", graph.EdgeDependsOn)
}

func TestNulPathsPreservesValidFilenameCharacters(t *testing.T) {
	got := nulPaths([]byte("dir/a\nb.txt\x00space name.go\x00"))
	if len(got) != 2 || got[0] != "dir/a\nb.txt" || got[1] != "space name.go" {
		t.Fatalf("nulPaths() = %#v", got)
	}
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

func assertNode(t *testing.T, nodes []graph.Node, id string, kind graph.NodeKind, name string) graph.Node {
	t.Helper()
	for _, n := range nodes {
		if n.ID == id {
			if n.Kind != kind || n.Name != name {
				t.Fatalf("node %s = kind %s name %q, want %s %q", id, n.Kind, n.Name, kind, name)
			}
			return n
		}
	}
	t.Fatalf("missing node %s", id)
	return graph.Node{}
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
