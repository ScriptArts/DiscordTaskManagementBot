package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir = filepath.Join(dir, "../", "test_tmp")
	os.Setenv("DISCORD_TASK_MANAGEMENT_DIR", dir)

	err = LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadEnvInvalid(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir = filepath.Join(dir, "../", "test_tmp", "invalid_data")
	os.Setenv("DISCORD_TASK_MANAGEMENT_DIR", dir)

	err = LoadEnv()
	if err == nil {
		t.Fatal(err)
	}
}
