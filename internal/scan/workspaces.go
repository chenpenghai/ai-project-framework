package scan

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chenpenghai/ai-project-framework/internal/graph"
)

func discoverWorkspaceEdges(root string, projects []locatedNode) []graph.Edge {
	projectByPath := map[string]string{}
	for _, project := range projects {
		projectByPath[filepath.ToSlash(project.Path)] = project.Node.ID
	}

	var edges []graph.Edge
	for _, project := range projects {
		for _, manifest := range project.Node.Evidence {
			if filepath.Base(manifest) != "go.work" {
				continue
			}
			manifestPath := filepath.Join(root, filepath.FromSlash(manifest))
			manifestDir := filepath.ToSlash(filepath.Dir(manifest))
			for _, usePath := range readGoWorkUses(manifestPath) {
				memberPath := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.FromSlash(manifestDir), filepath.FromSlash(usePath))))
				if memberPath == "" {
					memberPath = "."
				}
				if id, ok := projectByPath[memberPath]; ok && id != project.Node.ID {
					edges = append(edges, graph.Edge{
						From:       "file:" + manifest,
						To:         id,
						Kind:       graph.EdgeGoverns,
						Source:     "go.work",
						Confidence: graph.ConfidenceDeclared,
						Evidence:   []string{manifest + ": use " + usePath},
					})
				}
			}
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		if edges[i].From != edges[j].From {
			return edges[i].From < edges[j].From
		}
		return edges[i].To < edges[j].To
	})
	return edges
}

func readGoWorkUses(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var uses []string
	inUseBlock := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(stripLineComment(scanner.Text()))
		if line == "" {
			continue
		}
		if line == "use (" {
			inUseBlock = true
			continue
		}
		if inUseBlock && line == ")" {
			inUseBlock = false
			continue
		}
		if strings.HasPrefix(line, "use ") {
			if value := firstField(strings.TrimSpace(strings.TrimPrefix(line, "use "))); value != "" {
				uses = append(uses, value)
			}
			continue
		}
		if inUseBlock {
			if value := firstField(line); value != "" {
				uses = append(uses, value)
			}
		}
	}
	sort.Strings(uses)
	return unique(uses)
}
