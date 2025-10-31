package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findYAMLFiles(root string) ([]string, error) {
	var yamlFiles []string

	// filepath.WalkDir walks the file tree rooted at root
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			// Return the error to stop the walk or handle it as needed
			return e
		}

		// Check if it's a file and has the correct extension
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(d.Name()))
			if ext == ".yaml" || ext == ".yml" {
				yamlFiles = append(yamlFiles, s)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return yamlFiles, nil
}


func main () {

	// TODO: Have this customizable
	HOME := os.Getenv("HOME")
	kubePath := filepath.Join(HOME, ".kube")
	confiDotDPath := filepath.Join(kubePath, "config.d")
	outputPath := filepath.Join(kubePath, "config")

	yamlFiles, err := findYAMLFiles(confiDotDPath)
	if err != nil {
		log.Fatalf("Error finding YAML files: %v", err)
	}

	err = os.Setenv("KUBECONFIG", strings.Join(yamlFiles, ":"))
	if err != nil {
		log.Fatalf("Error setting environment variable: %v", err)
	}

	cmd := exec.Command("kubectl", "config", "view", "--merge", "--flatten")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing kubectl command: %v\nOutput: %s", err, output)
	}

	err = os.WriteFile(outputPath, []byte(output), 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

}