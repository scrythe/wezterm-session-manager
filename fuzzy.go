package main

import (
	"fmt"
	"github.com/sahilm/fuzzy"
	"os"
)

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

func fuzzyFind() {
	data := []string{"game.cpp", "moduleNameResolver.ts", "my name is_Ramsey"}
	matches := fuzzy.Find("cp", data)
	for _, match := range matches {
		for i := 0; i < len(match.Str); i++ {
			if contains(i, match.MatchedIndexes) {
				fmt.Printf("<%s>", string(match.Str[i]))
			} else {
				fmt.Print(string(match.Str[i]))
			}
		}
		fmt.Println()
	}
}

func listFolders() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	path := fmt.Sprintf("%v/projects", homeDir)
	fmt.Print(path)
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Print(err)
		return
	}
	folders := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}
	for _, folder := range folders {
		fmt.Println(folder)
	}
}
