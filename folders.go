package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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

type Workspace struct {
	Workspace string
}

func containsItem(list []string, item string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}

func listWorkspaces() []string {
	cmd := exec.Command("wezterm", "cli", "list", "--format", "json")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	var workspacesStruct []Workspace
	json.Unmarshal(out, &workspacesStruct)
	workspaces := []string{}
	for _, workspace := range workspacesStruct {
		workspaceName := workspace.Workspace
		if !containsItem(workspaces, workspaceName) {
			workspaces = append(workspaces, workspaceName)
		}
	}
	return workspaces
}

func listAll() []string {
	workspaces := listWorkspaces()
	workspaceLen := len(workspaces)
	folders := listFolders()
	all := make([]string, len(workspaces)+len(folders))
	for i := 0; i < len(all); i++ {
		if i < workspaceLen {
			all[i] = workspaces[i]
			continue
		}
		all[i] = folders[i-workspaceLen]
	}
	return all
}
