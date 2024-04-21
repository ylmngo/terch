package main

import (
	"os"
	"testing"
)

func TestRenameFile(t *testing.T) {
	TEST_FILES := []string{
		"../uploads/58_bgnet.pdf",
	}

	for _, tf := range TEST_FILES {
		if err := os.Rename(tf, "../uploads/15.pdf"); err != nil {
			t.Fatalf("Unable to rename file: %v\n", err)
		}
	}
}
