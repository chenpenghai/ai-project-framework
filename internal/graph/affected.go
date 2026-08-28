package graph

import "sort"

// AffectedProjects maps changed files to their owning projects and follows
// reverse declared dependency edges. It does not assume that every repository
// file affects every project.
func AffectedProjects(snapshot Snapshot, changedFiles []string) []string {
	projects := map[string]bool{}
	for _, node := range snapshot.Nodes {
		if node.Kind == NodeProject {
			projects[node.ID] = true
		}
	}

	parent := map[string]string{}
	for _, edge := range snapshot.Edges {
		if edge.Kind == EdgeContains {
			parent[edge.To] = edge.From
		}
	}

	affected := map[string]bool{}
	for _, path := range changedFiles {
		current := "file:" + path
		seen := map[string]bool{}
		for current != "" && !seen[current] {
			seen[current] = true
			owner, ok := parent[current]
			if !ok {
				break
			}
			if projects[owner] {
				affected[owner] = true
				break
			}
			current = owner
		}
	}

	reverse := map[string][]string{}
	for _, edge := range snapshot.Edges {
		if edge.Kind != EdgeDependsOn || !projects[edge.From] || !projects[edge.To] {
			continue
		}
		reverse[edge.To] = append(reverse[edge.To], edge.From)
	}

	queue := make([]string, 0, len(affected))
	for id := range affected {
		queue = append(queue, id)
	}
	sort.Strings(queue)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, dependent := range reverse[current] {
			if affected[dependent] {
				continue
			}
			affected[dependent] = true
			queue = append(queue, dependent)
		}
	}

	out := make([]string, 0, len(affected))
	for id := range affected {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
