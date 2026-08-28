package graph

import "testing"

func TestDependencyCyclesFindsProjectSCCs(t *testing.T) {
	snapshot := Snapshot{
		Nodes: []Node{
			{ID: "project:a", Kind: NodeProject},
			{ID: "project:b", Kind: NodeProject},
			{ID: "project:c", Kind: NodeProject},
		},
		Edges: []Edge{
			{From: "project:a", To: "project:b", Kind: EdgeDependsOn},
			{From: "project:b", To: "project:a", Kind: EdgeDependsOn},
			{From: "project:b", To: "project:c", Kind: EdgeDependsOn},
		},
	}

	cycles := DependencyCycles(snapshot)
	if len(cycles) != 1 {
		t.Fatalf("cycles = %#v, want one", cycles)
	}
	if len(cycles[0]) != 2 || cycles[0][0] != "project:a" || cycles[0][1] != "project:b" {
		t.Fatalf("cycle = %#v, want [project:a project:b]", cycles[0])
	}
}
