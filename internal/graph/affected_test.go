package graph

import "testing"

func TestAffectedProjectsIncludesReverseDependencyClosure(t *testing.T) {
	snapshot := Snapshot{
		Nodes: []Node{
			{ID: "project:app", Kind: NodeProject},
			{ID: "project:core", Kind: NodeProject},
			{ID: "project:shared", Kind: NodeProject},
			{ID: "file:shared/value.go", Kind: NodeFile},
		},
		Edges: []Edge{
			{From: "project:shared", To: "file:shared/value.go", Kind: EdgeContains},
			{From: "project:core", To: "project:shared", Kind: EdgeDependsOn},
			{From: "project:app", To: "project:core", Kind: EdgeDependsOn},
		},
	}

	got := AffectedProjects(snapshot, []string{"shared/value.go"})
	want := []string{"project:app", "project:core", "project:shared"}
	if len(got) != len(want) {
		t.Fatalf("affected = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("affected = %#v, want %#v", got, want)
		}
	}
}
