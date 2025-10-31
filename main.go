package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findYAMLFiles(root string) ([]string, error) {
	var yamlFiles []string

	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
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

	
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
	}

	kubePath := filepath.Join(homePath, ".kube")
	configDPath := filepath.Join(kubePath, "config.d") // TODO: Have this customizable
	outputPath := filepath.Join(kubePath, "config")

	fmt.Printf("Merging YAML files at %v into %v\n", configDPath, outputPath)

	yamlFiles, err := findYAMLFiles(configDPath)
	if err != nil {
		fmt.Printf("Error finding YAML files: %v\n", err)
	}

	fmt.Printf("Found %d files \n", len(yamlFiles))

	err = os.Setenv("KUBECONFIG", strings.Join(yamlFiles, ":"))
	if err != nil {
		fmt.Printf("Error setting environment variable: %v\n", err)
	}

	cmd := exec.Command("kubectl", "config", "view", "--merge", "--flatten")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing kubectl command: %v\nOutput: %s", err, output)
	}

	err = os.WriteFile(outputPath, []byte(output), 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
	}

	fmt.Printf("Generated %v \n", outputPath)

}