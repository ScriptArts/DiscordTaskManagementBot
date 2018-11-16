package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitialize(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir = filepath.Join(dir, "test_tmp")
	os.Setenv("DISCORD_TASK_MANAGEMENT_DIR", dir)

	err = initialize()
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitializeInvalid(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir = filepath.Join(dir, "test_tmp", "invalid_data")
	os.Setenv("DISCORD_TASK_MANAGEMENT_DIR", dir)

	err = initialize()
	if err == nil {
		t.Fatal(err)
	}
}
