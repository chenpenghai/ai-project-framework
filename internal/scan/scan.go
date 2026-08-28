package scan

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chenpenghai/ai-project-framework/internal/graph"
)

type manifestSpec struct {
	Ecosystem  string
	Confidence graph.Confidence
}

var manifestSpecs = map[string]manifestSpec{
	"package.json":        {"node", graph.ConfidenceDeclared},
	"go.mod":              {"go", graph.ConfidenceDeclared},
	"go.work":             {"go-workspace", graph.ConfidenceDeclared},
	"Cargo.toml":          {"rust", graph.ConfidenceDeclared},
	"pyproject.toml":      {"python", graph.ConfidenceDeclared},
	"setup.py":            {"python", graph.ConfidenceDeclared},
	"setup.cfg":           {"python", graph.ConfidenceDeclared},
	"requirements.txt":    {"python", graph.ConfidenceInferredHigh},
	"Pipfile":             {"python", graph.ConfidenceInferredHigh},
	"pom.xml":             {"maven", graph.ConfidenceDeclared},
	"build.gradle":        {"gradle", graph.ConfidenceDeclared},
	"build.gradle.kts":    {"gradle", graph.ConfidenceDeclared},
	"settings.gradle":     {"gradle-workspace", graph.ConfidenceDeclared},
	"settings.gradle.kts": {"gradle-workspace", graph.ConfidenceDeclared},
	"CMakeLists.txt":      {"cmake", graph.ConfidenceDeclared},
	"meson.build":         {"meson", graph.ConfidenceDeclared},
	"composer.json":       {"php", graph.ConfidenceDeclared},
	"Gemfile":             {"ruby", graph.ConfidenceInferredHigh},
	"pubspec.yaml":        {"dart", graph.ConfidenceDeclared},
	"Package.swift":       {"swift", graph.ConfidenceDeclared},
	"mix.exs":             {"elixir", graph.ConfidenceDeclared},
}

var fallbackSkipDirs = map[string]bool{
	".git": true, "node_modules": true, ".venv": true, "venv": true,
	"dist": true, "build": true, "target": true, ".idea": true, ".vscode": true,
}

type Scanner struct{}

func (Scanner) Scan(root string) (graph.Snapshot, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return graph.Snapshot{}, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return graph.Snapshot{}, err
	}
	if !info.IsDir() {
		return graph.Snapshot{}, fmt.Errorf("scan root is not a directory: %s", abs)
	}

	gitRoot, gitOK := detectGitRoot(abs)
	if gitOK {
		abs = gitRoot
	}

	files, err := listFiles(abs, gitOK)
	if err != nil {
		return graph.Snapshot{}, err
	}
	changed := []string(nil)
	if gitOK {
		changed = gitChangedFiles(abs)
	}

	snapshot := graph.Snapshot{
		Version: 1,
		Root:    filepath.ToSlash(abs),
		Git: graph.GitState{
			Available:    gitOK,
			ChangedFiles: changed,
		},
	}

	repoID := "repository:."
	snapshot.Nodes = append(snapshot.Nodes, graph.Node{
		ID: repoID, Kind: graph.NodeRepository, Name: filepath.Base(abs), Source: "filesystem", Confidence: graph.ConfidenceObserved,
	})

	projects := discoverProjects(abs, files)
	modules := discoverModules(abs, files)

	for _, p := range projects {
		snapshot.Nodes = append(snapshot.Nodes, p.Node)
	}
	for _, m := range modules {
		snapshot.Nodes = append(snapshot.Nodes, m.Node)
	}
	for _, rel := range files {
		snapshot.Nodes = append(snapshot.Nodes, graph.Node{
			ID: "file:" + rel, Kind: graph.NodeFile, Name: filepath.Base(rel), Path: rel, Source: "filesystem", Confidence: graph.ConfidenceObserved,
			Metadata: fileMetadata(rel),
		})
	}

	// Project hierarchy is derived from manifest locations.
	for _, p := range projects {
		parent := repoID
		if q := nearestProjectAncestor(p.Path, projects); q != nil {
			parent = q.Node.ID
		}
		snapshot.Edges = append(snapshot.Edges, graph.Edge{From: parent, To: p.Node.ID, Kind: graph.EdgeContains, Source: "manifest", Confidence: graph.ConfidenceDerived})
	}

	// Explicit logical modules attach to the nearest project, otherwise repo.
	for _, m := range modules {
		parent := repoID
		if p := nearestProjectForPath(m.Path, projects); p != nil {
			parent = p.Node.ID
		}
		snapshot.Edges = append(snapshot.Edges, graph.Edge{From: parent, To: m.Node.ID, Kind: graph.EdgeContains, Source: "MODULE.md", Confidence: graph.ConfidenceDerived})
	}

	// Each file belongs to the nearest explicit module, else nearest project, else repo.
	for _, rel := range files {
		parent := repoID
		if m := nearestModuleForPath(rel, modules); m != nil {
			parent = m.Node.ID
		} else if p := nearestProjectForPath(rel, projects); p != nil {
			parent = p.Node.ID
		}
		snapshot.Edges = append(snapshot.Edges, graph.Edge{From: parent, To: "file:" + rel, Kind: graph.EdgeContains, Source: "derived-containment", Confidence: graph.ConfidenceDerived})
	}

	sortSnapshot(&snapshot)
	return snapshot, nil
}

type locatedNode struct {
	Path string
	Node graph.Node
}

func discoverProjects(root string, files []string) []locatedNode {
	type aggregate struct {
		path       string
		name       string
		confidence graph.Confidence
		manifests  []string
		ecosystems []string
	}
	byPath := map[string]*aggregate{}

	for _, rel := range files {
		base := filepath.Base(rel)
		spec, ok := manifestSpecs[base]
		if !ok && !strings.HasSuffix(strings.ToLower(base), ".csproj") && !strings.HasSuffix(strings.ToLower(base), ".sln") {
			continue
		}
		if !ok {
			if strings.HasSuffix(strings.ToLower(base), ".csproj") {
				spec = manifestSpec{"dotnet", graph.ConfidenceDeclared}
			} else {
				spec = manifestSpec{"dotnet-workspace", graph.ConfidenceDeclared}
			}
		}
		dir := filepath.ToSlash(filepath.Dir(rel))
		a := byPath[dir]
		if a == nil {
			a = &aggregate{path: dir, confidence: spec.Confidence}
			byPath[dir] = a
		}
		a.manifests = append(a.manifests, rel)
		a.ecosystems = append(a.ecosystems, spec.Ecosystem)
		if confidenceRank(spec.Confidence) > confidenceRank(a.confidence) {
			a.confidence = spec.Confidence
		}
		if name := projectName(filepath.Join(root, filepath.FromSlash(rel)), base, dir); name != "" {
			if a.name == "" || base == "package.json" || base == "go.mod" {
				a.name = name
			}
		}
	}

	out := make([]locatedNode, 0, len(byPath))
	for _, a := range byPath {
		sort.Strings(a.manifests)
		a.ecosystems = sortedUnique(a.ecosystems)
		if a.name == "" {
			if a.path == "." {
				a.name = filepath.Base(root)
			} else {
				a.name = filepath.Base(a.path)
			}
		}
		out = append(out, locatedNode{Path: a.path, Node: graph.Node{
			ID: "project:" + a.path, Kind: graph.NodeProject, Name: a.name, Path: a.path,
			Source: "project-manifest", Confidence: a.confidence, Evidence: append([]string(nil), a.manifests...),
			Metadata: map[string]string{
				"ecosystems": strings.Join(a.ecosystems, ","),
				"manifests":  strings.Join(a.manifests, ","),
			},
		}})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Node.ID < out[j].Node.ID })
	return out
}

func confidenceRank(c graph.Confidence) int {
	switch c {
	case graph.ConfidenceDeclared:
		return 4
	case graph.ConfidenceObserved:
		return 3
	case graph.ConfidenceDerived:
		return 3
	case graph.ConfidenceInferredHigh:
		return 2
	case graph.ConfidenceCandidate:
		return 1
	default:
		return 0
	}
}

func sortedUnique(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	sort.Strings(in)
	return unique(in)
}

func discoverModules(root string, files []string) []locatedNode {
	var out []locatedNode
	for _, rel := range files {
		if filepath.Base(rel) != "MODULE.md" {
			continue
		}
		dir := filepath.ToSlash(filepath.Dir(rel))
		name := moduleName(filepath.Join(root, filepath.FromSlash(rel)), dir)
		out = append(out, locatedNode{Path: dir, Node: graph.Node{
			ID: "module:" + dir, Kind: graph.NodeModule, Name: name, Path: dir, Source: "MODULE.md", Confidence: graph.ConfidenceDeclared, Evidence: []string{rel},
			Metadata: map[string]string{"declaration": rel},
		}})
	}
	return dedupeLocated(out)
}

func dedupeLocated(in []locatedNode) []locatedNode {
	byID := map[string]locatedNode{}
	for _, n := range in {
		if _, exists := byID[n.Node.ID]; !exists {
			byID[n.Node.ID] = n
		}
	}
	out := make([]locatedNode, 0, len(byID))
	for _, n := range byID {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Node.ID < out[j].Node.ID })
	return out
}

func nearestProjectAncestor(path string, projects []locatedNode) *locatedNode {
	var best *locatedNode
	for i := range projects {
		p := &projects[i]
		if p.Path == path {
			continue
		}
		if containsPath(p.Path, path) && (best == nil || pathDepth(p.Path) > pathDepth(best.Path)) {
			best = p
		}
	}
	return best
}

func nearestProjectForPath(path string, projects []locatedNode) *locatedNode {
	var best *locatedNode
	for i := range projects {
		p := &projects[i]
		if containsPath(p.Path, path) && (best == nil || pathDepth(p.Path) > pathDepth(best.Path)) {
			best = p
		}
	}
	return best
}

func nearestModuleForPath(path string, modules []locatedNode) *locatedNode {
	var best *locatedNode
	for i := range modules {
		m := &modules[i]
		if containsPath(m.Path, path) && (best == nil || pathDepth(m.Path) > pathDepth(best.Path)) {
			best = m
		}
	}
	return best
}

func containsPath(parent, child string) bool {
	if parent == "." {
		return true
	}
	parent = strings.TrimSuffix(filepath.ToSlash(parent), "/")
	child = filepath.ToSlash(child)
	return child == parent || strings.HasPrefix(child, parent+"/")
}

func pathDepth(path string) int {
	if path == "." || path == "" {
		return 0
	}
	return strings.Count(filepath.ToSlash(path), "/") + 1
}

func fileMetadata(rel string) map[string]string {
	ext := strings.ToLower(filepath.Ext(rel))
	if ext == "" {
		return nil
	}
	return map[string]string{"extension": ext}
}

func projectName(path, base, dir string) string {
	if base == "package.json" {
		data, err := os.ReadFile(path)
		if err == nil {
			var x struct {
				Name string `json:"name"`
			}
			if json.Unmarshal(data, &x) == nil && strings.TrimSpace(x.Name) != "" {
				return x.Name
			}
		}
	}
	if base == "go.mod" {
		f, err := os.Open(path)
		if err == nil {
			defer f.Close()
			s := bufio.NewScanner(f)
			for s.Scan() {
				line := strings.TrimSpace(s.Text())
				if strings.HasPrefix(line, "module ") {
					return strings.TrimSpace(strings.TrimPrefix(line, "module "))
				}
			}
		}
	}
	if dir == "." {
		return filepath.Base(filepath.Dir(path))
	}
	return filepath.Base(dir)
}

func moduleName(path, dir string) string {
	data, err := os.ReadFile(path)
	if err == nil {
		s := bufio.NewScanner(strings.NewReader(string(data)))
		for s.Scan() {
			line := strings.TrimSpace(s.Text())
			if strings.HasPrefix(line, "module:") {
				if v := strings.TrimSpace(strings.TrimPrefix(line, "module:")); v != "" {
					return strings.Trim(v, "\"'")
				}
		}
	}
	if dir == "." {
		return filepath.Base(filepath.Dir(path))
	}
	return filepath.Base(dir)
}

func detectGitRoot(root string) (string, bool) {
	cmd := exec.Command("git", "-C", root, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	p := strings.TrimSpace(string(out))
	if p == "" {
		return "", false
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", false
	}
	return abs, true
}

func listFiles(root string, gitOK bool) ([]string, error) {
	if gitOK {
		cmd := exec.Command("git", "-C", root, "ls-files", "-co", "--exclude-standard", "-z")
		out, err := cmd.Output()
		if err == nil {
			parts := strings.Split(string(out), "\x00")
			files := make([]string, 0, len(parts))
			for _, p := range parts {
				if p == "" {
					continue
				}
				files = append(files, filepath.ToSlash(p))
			}
			sort.Strings(files)
			return unique(files), nil
		}
	}

	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			if path != root && fallbackSkipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return unique(files), nil
}

func gitChangedFiles(root string) []string {
	seen := map[string]struct{}{}
	commands := [][]string{
		{"diff", "--name-only", "--relative"},
		{"diff", "--cached", "--name-only", "--relative"},
		{"ls-files", "--others", "--exclude-standard"},
	}
	for _, args := range commands {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		out, err := cmd.Output()
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				seen[filepath.ToSlash(line)] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

func unique(in []string) []string {
	if len(in) < 2 {
		return in
	}
	out := in[:0]
	var prev string
	for i, s := range in {
		if i == 0 || s != prev {
			out = append(out, s)
		}
		prev = s
	}
	return out
}

func sortSnapshot(s *graph.Snapshot) {
	sort.Slice(s.Nodes, func(i, j int) bool { return s.Nodes[i].ID < s.Nodes[j].ID })
	sort.Slice(s.Edges, func(i, j int) bool {
		if s.Edges[i].From != s.Edges[j].From {
			return s.Edges[i].From < s.Edges[j].From
		}
		if s.Edges[i].To != s.Edges[j].To {
			return s.Edges[i].To < s.Edges[j].To
		}
		return s.Edges[i].Kind < s.Edges[j].Kind
	})
}
