package main

import (
	"fmt"
	"os"
)

func listFolders() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
    fmt.Printf("Err: %v", err)
    os.Exit(1)

	}
	path := fmt.Sprintf("%v/projects", homeDir)
	fmt.Print(path)
	entries, err := os.ReadDir(path)
	if err != nil {
    fmt.Printf("Err: %v", err)
    os.Exit(1)
	}
	folders := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}
	return folders
}
