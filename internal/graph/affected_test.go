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
	assertStrings(t, got, want)
}

func TestAffectedProjectsIncludesWorkspaceMembersGovernedByChangedFile(t *testing.T) {
	snapshot := Snapshot{
		Nodes: []Node{
			{ID: "project:.", Kind: NodeProject},
			{ID: "project:app", Kind: NodeProject},
			{ID: "project:shared", Kind: NodeProject},
			{ID: "file:go.work", Kind: NodeFile},
		},
		Edges: []Edge{
			{From: "project:.", To: "file:go.work", Kind: EdgeContains},
			{From: "file:go.work", To: "project:shared", Kind: EdgeGoverns},
			{From: "project:app", To: "project:shared", Kind: EdgeDependsOn},
		},
	}

	got := AffectedProjects(snapshot, []string{"go.work"})
	want := []string{"project:.", "project:app", "project:shared"}
	assertStrings(t, got, want)
}

func assertStrings(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got = %#v, want %#v", got, want)
		}
	}
}
