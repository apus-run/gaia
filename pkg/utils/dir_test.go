package utils

import (
	"bou.ke/monkey"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWritableDir(t *testing.T) {
	filename := ""
	dir := MakeDirectory(filename)
	assert.Equal(t, GetWorkDir(), dir)

	filename = "./test.txt"
	dir = MakeDirectory(filename)
	exp, _ := filepath.Abs(filename)
	assert.Equal(t, exp, dir)

	filename = "./none/existed/test.txt"
	exp, _ = filepath.Abs(filename)
	dir = MakeDirectory(filename)
	os.RemoveAll("./none")
	assert.Equal(t, exp, dir)

	filename = "~/none/existed/test.txt"
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.TempDir()
	}
	exp = filepath.Join(home, "none/existed/test.txt")
	dir = MakeDirectory(filename)
	os.RemoveAll(home + "/none")
	assert.Equal(t, exp, dir)
}

func TestGetWorkDirFail(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(os.Getwd, func() (string, error) {
		return "", fmt.Errorf("error")
	})

	path := GetWorkDir()
	home, err := os.UserHomeDir()
	assert.Nil(t, err)
	assert.Equal(t, path, home)

	monkey.Patch(os.UserHomeDir, func() (string, error) {
		return "", fmt.Errorf("error")
	})

	path = GetWorkDir()
	assert.Equal(t, path, os.TempDir())

}

func TestMakeDirectoryFail(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch(os.UserHomeDir, func() (string, error) {
		return "", fmt.Errorf("error")
	})

	filename := "~/test.txt"
	result := MakeDirectory(filename)
	assert.Equal(t, result, filepath.Join(os.TempDir(), filename[2:]))

	monkey.Unpatch(os.UserHomeDir)

	monkey.Patch(filepath.Abs, func(path string) (string, error) {
		return "", fmt.Errorf("error")
	})
	filename = "../test.txt"
	result = MakeDirectory(filename)
	assert.Equal(t, result, filepath.Join(GetWorkDir(), "test.txt"))

	monkey.Unpatch(filepath.Abs)

	monkey.Patch(os.MkdirAll, func(string, os.FileMode) error {
		return fmt.Errorf("error")
	})

	filename = "/not/existed/test.txt"
	result = MakeDirectory(filename)
	assert.Equal(t, result, filepath.Join(GetWorkDir(), "test.txt"))

	monkey.Unpatch(os.MkdirAll)

}
