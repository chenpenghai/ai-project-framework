package consumer

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestNewProjectCreatesOnlyFrameworkControlFiles(t *testing.T) {
	target := filepath.Join(t.TempDir(), "consumer")
	if err := NewProject(target); err != nil {
		t.Fatal(err)
	}

	var got []string
	err := filepath.WalkDir(target, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(target, path)
		if err != nil {
			return err
		}
		got = append(got, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)

	want := []string{".apf/project.yaml", ".gitignore", "AGENTS.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("consumer files = %#v, want %#v", got, want)
	}

	for rel, wantContent := range Files {
		data, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != wantContent {
			t.Fatalf("%s content changed unexpectedly", rel)
		}
	}
}

func TestNewProjectRefusesNonEmptyDirectory(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "user.txt"), []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(target); err == nil {
		t.Fatal("expected non-empty target to be rejected")
	}
	data, err := os.ReadFile(filepath.Join(target, "user.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "keep me" {
		t.Fatal("existing user file was modified")
	}
}
