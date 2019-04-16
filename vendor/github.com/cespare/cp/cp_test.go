package cp

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestCopyFile(t *testing.T) {
	td := newTestDir(t)
	defer td.remove()

	td.create("a.txt", "contents of a", 0644)
	if err := CopyFile(td.path("b.txt"), td.path("a.txt")); err != nil {
		t.Errorf("CopyFile(b.txt, a.txt): %s", err)
	}
	td.checkContentsMode("b.txt", "contents of a", 0644)

	if err := CopyFile(td.path("b.txt"), td.path("a.txt")); !os.IsExist(err) {
		t.Errorf("CopyFile(b.txt, a.txt) (second time): got %v; want os.IsExist(err)", err)
	}
	td.checkContentsMode("b.txt", "contents of a", 0644)

	if err := CopyFile(td.path("c.txt"), td.path("nonexistent")); !os.IsNotExist(err) {
		t.Errorf("CopyFile(c.txt, nonexistent): got %v; want os.IsNotExist(err)", err)
	}

	td.mkdir("d", 0755)
	if err := CopyFile(td.path("c.txt"), td.path("d")); err != errCopyFileWithDir {
		t.Errorf("CopyFile(c.txt, d): got %v; want errCopyFileWithDir", err)
	}
}

func TestCopyFileOverwrite(t *testing.T) {
	td := newTestDir(t)
	defer td.remove()

	td.create("a.txt", "first", 0644)
	if err := CopyFileOverwrite(td.path("b.txt"), td.path("a.txt")); err != nil {
		t.Errorf("CopyFileOverwite(b.txt, a.txt): %s", err)
	}
	td.checkContentsMode("b.txt", "first", 0644)

	td.create("c.txt", "second", 0600)
	if err := CopyFileOverwrite(td.path("b.txt"), td.path("c.txt")); err != nil {
		t.Errorf("CopyFileOverwite(b.txt, c.txt): %s", err)
	}
	td.checkContentsMode("b.txt", "second", 0600)

	err := CopyFileOverwrite(td.path("b.txt"), td.path("nonexistent"))
	if !os.IsNotExist(err) {
		t.Errorf("CopyFileOverwite(b.txt, nonexistent): got %v; want os.IsNotExist(err)", err)
	}

	td.mkdir("d", 0755)
	if err := CopyFileOverwrite(td.path("b.txt"), td.path("d")); err != errCopyFileWithDir {
		t.Errorf("CopyFileOverwrite(c.txt, d): got %v; want errCopyFileWithDir", err)
	}
}

func TestCopyAll(t *testing.T) {
	td := newTestDir(t)
	defer td.remove()

	td.mkdir("d0", 0755)
	td.mkdir("d0/d1", 0755)
	td.create("d0/a.txt", "a", 0644)
	td.create("d0/d1/b.txt", "b", 0644)
	td.create("d0/d1/c.txt", "c", 0644)

	if err := CopyAll(td.path("target"), td.path("d0")); err != nil {
		t.Fatal(err)
	}

	td.checkAll(
		"target",
		"a.txt", "a",
		"d1/b.txt", "b",
		"d1/c.txt", "c",
	)

	err := CopyAll(td.path("target"), td.path("d0"))
	if !os.IsExist(err) {
		t.Errorf("CopyFile(target, d0) (second time): got %v; want os.IsExist(err)", err)
	}
}

func TestCopyAllOverwrite(t *testing.T) {
	td := newTestDir(t)
	defer td.remove()

	td.mkdir("d0", 0755)
	td.mkdir("d0/d1", 0755)
	td.create("d0/a.txt", "a", 0644)
	td.create("d0/d1/b.txt", "b", 0644)
	td.create("d0/d1/c.txt", "c", 0644)

	td.mkdir("target", 0777)
	td.mkdir("target/d1", 0777)
	td.create("target/d1/b.txt", "bbb", 0660)

	if err := CopyAllOverwrite(td.path("target"), td.path("d0")); err != nil {
		t.Fatal(err)
	}

	td.checkAll(
		"target",
		"a.txt", "a",
		"d1/b.txt", "b",
		"d1/c.txt", "c",
	)
}

func TestCopyDotSlash(t *testing.T) {
	td := newTestDir(t)
	defer td.remove()

	if err := CopyAllOverwrite(td.path("testdata"), "./testdata"); err != nil {
		t.Fatal(err)
	}
	td.checkAll("testdata", "a.txt", "a\n")
}

type testDir struct {
	t   *testing.T
	dir string
}

func newTestDir(t *testing.T) testDir {
	t.Helper()
	dir, err := ioutil.TempDir("", "cp-test-")
	if err != nil {
		t.Fatalf("Cannot create tempdir for test: %s", err)
	}
	return testDir{t, dir}
}

func (td testDir) remove() {
	td.t.Helper()
	if err := os.RemoveAll(td.dir); err != nil {
		td.t.Errorf("Error cleaning up tempdir: %s", err)
	}
}

func (td testDir) path(name string) string {
	return filepath.Join(td.dir, name)
}

func (td testDir) create(name, contents string, perm os.FileMode) {
	td.t.Helper()
	err := ioutil.WriteFile(td.path(name), []byte(contents), perm)
	if err != nil {
		td.t.Fatal(err)
	}
}

func (td testDir) mkdir(name string, perm os.FileMode) {
	td.t.Helper()
	if err := os.Mkdir(td.path(name), perm); err != nil {
		td.t.Fatal(err)
	}
}

func (td testDir) checkContentsMode(name, contents string, perm os.FileMode) {
	td.t.Helper()
	b, err := ioutil.ReadFile(td.path(name))
	if err != nil {
		td.t.Error(err)
		return
	}
	if got := string(b); got != contents {
		td.t.Errorf("for %s, got contents %q; want %q", name, got, contents)
		return
	}
	stat, err := os.Stat(td.path(name))
	if err != nil {
		td.t.Errorf("os.Stat(%s): %s", name, err)
		return
	}
	if got := stat.Mode(); got != perm {
		td.t.Errorf("for %s, got perm %s; want %s", name, got, perm)
	}
}

func (td testDir) checkAll(dir string, nameContents ...string) {
	td.t.Helper()
	if len(nameContents)%2 != 0 {
		panic("bad nameContents pairs")
	}
	names := make(map[string]struct{})
	for i := 0; i < len(nameContents); i += 2 {
		name := nameContents[i]
		names[name] = struct{}{}
		td.checkContentsMode(filepath.Join(dir, name), nameContents[i+1], 0644)
	}

	var all []string
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(td.path(dir), path)
		if err != nil {
			return err
		}
		all = append(all, rel)
		return nil
	}
	if err := filepath.Walk(td.path(dir), walk); err != nil {
		td.t.Fatal(err)
	}
	sort.Strings(all)
	for _, name := range all {
		if _, ok := names[name]; !ok {
			td.t.Errorf("%s unexpectedly contained file %s", dir, name)
		}
	}
}
