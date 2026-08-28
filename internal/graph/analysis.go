package graph

import "sort"

// DependencyCycles returns strongly connected project components in the
// declared dependency graph. Each returned component represents at least one
// dependency cycle. The output is deterministic.
func DependencyCycles(snapshot Snapshot) [][]string {
	projects := map[string]bool{}
	for _, node := range snapshot.Nodes {
		if node.Kind == NodeProject {
			projects[node.ID] = true
		}
	}

	adj := map[string][]string{}
	selfLoop := map[string]bool{}
	for _, edge := range snapshot.Edges {
		if edge.Kind != EdgeDependsOn || !projects[edge.From] || !projects[edge.To] {
			continue
		}
		adj[edge.From] = append(adj[edge.From], edge.To)
		if edge.From == edge.To {
			selfLoop[edge.From] = true
		}
	}
	for id := range adj {
		sort.Strings(adj[id])
	}

	index := 0
	indices := map[string]int{}
	lowlink := map[string]int{}
	onStack := map[string]bool{}
	var stack []string
	var cycles [][]string

	var strongConnect func(string)
	strongConnect = func(v string) {
		indices[v] = index
		lowlink[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range adj[v] {
			if _, seen := indices[w]; !seen {
				strongConnect(w)
				if lowlink[w] < lowlink[v] {
					lowlink[v] = lowlink[w]
				}
			} else if onStack[w] && indices[w] < lowlink[v] {
				lowlink[v] = indices[w]
			}
		}

		if lowlink[v] != indices[v] {
			return
		}
		var component []string
		for {
			w := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			onStack[w] = false
			component = append(component, w)
			if w == v {
				break
			}
		}
		if len(component) > 1 || (len(component) == 1 && selfLoop[component[0]]) {
			sort.Strings(component)
			cycles = append(cycles, component)
		}
	}

	ids := make([]string, 0, len(projects))
	for id := range projects {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if _, seen := indices[id]; !seen {
			strongConnect(id)
		}
	}
	sort.Slice(cycles, func(i, j int) bool {
		return cycles[i][0] < cycles[j][0]
	})
	return cycles
}
