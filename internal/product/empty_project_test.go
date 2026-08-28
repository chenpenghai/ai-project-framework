package product

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"

	"github.com/chenpenghai/ai-project-framework/internal/graph"
	"github.com/chenpenghai/ai-project-framework/internal/scan"
)

func TestEmptyProjectIsStaticAndEmpty(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "empty-project"))

	var files []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode().Perm()&0o111 != 0 {
			t.Fatalf("empty project contains executable file: %s", path)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)

	want := []string{".apf/project.yaml", "AGENTS.md"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("empty-project files = %#v, want %#v", files, want)
	}

	copyRoot := t.TempDir()
	for _, rel := range want {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dst := filepath.Join(copyRoot, filepath.FromSlash(rel))
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	snapshot, err := (scan.Scanner{}).Scan(copyRoot)
	if err != nil {
		t.Fatal(err)
	}
	projects, modules := 0, 0
	for _, node := range snapshot.Nodes {
		switch node.Kind {
		case graph.NodeProject:
			projects++
		case graph.NodeModule:
			modules++
		}
	}
	if projects != 0 || modules != 0 {
		t.Fatalf("empty project must scan as empty: projects=%d modules=%d", projects, modules)
	}
}
