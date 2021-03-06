package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFile(t *testing.T) {
	var (
		path = filepath.Join(os.TempDir(), "test_config")
		file = filepath.Join(path, "test.json")
	)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(path)

	if err := os.MkdirAll(path, 0700); err != nil {
		t.Error(err)
	}

	testSource(t, file)
	testSource(t, path)
}

func testSource(t *testing.T, path string) {
	t.Log(path)
	s := NewSource(path)
	ss, err := s.Load()
	if err != nil {
		t.Error(err)
	}
	for _, f := range ss {
		t.Logf("文件名: %s", f)
	}
}
