package scan

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chenpenghai/ai-project-framework/internal/graph"
)

type goProjectInfo struct {
	projectID string
	manifest  string
	module    string
	requires  []string
}

type nodeProjectInfo struct {
	projectID    string
	manifest     string
	packageName  string
	dependencies map[string]string
}

func discoverDeclaredDependencies(root string, projects []locatedNode) []graph.Edge {
	var goProjects []goProjectInfo
	var nodeProjects []nodeProjectInfo

	for _, project := range projects {
		for _, manifest := range project.Node.Evidence {
			switch filepath.Base(manifest) {
			case "go.mod":
				module, requires := readGoMod(filepath.Join(root, filepath.FromSlash(manifest)))
				if module != "" {
					goProjects = append(goProjects, goProjectInfo{
						projectID: project.Node.ID,
						manifest:  manifest,
						module:    module,
						requires:  requires,
					})
				}
			case "package.json":
				name, dependencies := readPackageJSON(filepath.Join(root, filepath.FromSlash(manifest)))
				if name != "" {
					nodeProjects = append(nodeProjects, nodeProjectInfo{
						projectID:    project.Node.ID,
						manifest:     manifest,
						packageName:  name,
						dependencies: dependencies,
					})
				}
			}
		}
	}

	goOwners := map[string]string{}
	for _, project := range goProjects {
		goOwners[project.module] = project.projectID
	}
	nodeOwners := map[string]string{}
	for _, project := range nodeProjects {
		nodeOwners[project.packageName] = project.projectID
	}

	type edgeKey struct {
		from string
		to   string
	}
	edges := map[edgeKey]graph.Edge{}
	add := func(from, to, evidence string) {
		if from == "" || to == "" || from == to {
			return
		}
		key := edgeKey{from: from, to: to}
		edge, exists := edges[key]
		if !exists {
			edge = graph.Edge{
				From:       from,
				To:         to,
				Kind:       graph.EdgeDependsOn,
				Source:     "manifest-dependency",
				Confidence: graph.ConfidenceDeclared,
			}
		}
		edge.Evidence = append(edge.Evidence, evidence)
		edges[key] = edge
	}

	for _, project := range goProjects {
		for _, required := range project.requires {
			if owner, ok := goOwners[required]; ok {
				add(project.projectID, owner, project.manifest+": require "+required)
			}
		}
	}

	for _, project := range nodeProjects {
		for dependency := range project.dependencies {
			if owner, ok := nodeOwners[dependency]; ok {
				add(project.projectID, owner, project.manifest+": dependency "+dependency)
			}
		}
	}

	out := make([]graph.Edge, 0, len(edges))
	for _, edge := range edges {
		sort.Strings(edge.Evidence)
		edge.Evidence = unique(edge.Evidence)
		out = append(out, edge)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].From != out[j].From {
			return out[i].From < out[j].From
		}
		return out[i].To < out[j].To
	})
	return out
}

func readGoMod(path string) (string, []string) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	var module string
	var requires []string
	inRequireBlock := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(stripLineComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "module ") {
			module = firstField(strings.TrimSpace(strings.TrimPrefix(line, "module ")))
			continue
		}
		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}
		if strings.HasPrefix(line, "require ") {
			if required := firstField(strings.TrimSpace(strings.TrimPrefix(line, "require "))); required != "" {
				requires = append(requires, required)
			}
			continue
		}
		if inRequireBlock {
			if required := firstField(line); required != "" {
				requires = append(requires, required)
			}
		}
	}
	sort.Strings(requires)
	return module, unique(requires)
}

func readPackageJSON(path string) (string, map[string]string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}
	var manifest struct {
		Name                 string            `json:"name"`
		Dependencies         map[string]string `json:"dependencies"`
		DevDependencies      map[string]string `json:"devDependencies"`
		PeerDependencies     map[string]string `json:"peerDependencies"`
		OptionalDependencies map[string]string `json:"optionalDependencies"`
	}
	if json.Unmarshal(data, &manifest) != nil {
		return "", nil
	}
	all := map[string]string{}
	mergeStringMap(all, manifest.Dependencies)
	mergeStringMap(all, manifest.DevDependencies)
	mergeStringMap(all, manifest.PeerDependencies)
	mergeStringMap(all, manifest.OptionalDependencies)
	return strings.TrimSpace(manifest.Name), all
}

func mergeStringMap(dst, src map[string]string) {
	for key, value := range src {
		dst[key] = value
	}
}

func stripLineComment(line string) string {
	if i := strings.Index(line, "//"); i >= 0 {
		return line[:i]
	}
	return line
}

func firstField(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
