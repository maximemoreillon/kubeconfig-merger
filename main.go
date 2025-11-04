package main

import (
	"flag"
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

	sourcePath := flag.String("source", filepath.Join(kubePath, "config.d"), "Source directory")
	flag.Parse()

	outputPath := filepath.Join(kubePath, "config")

	fmt.Printf("Merging YAML files at %v into %v\n", *sourcePath, outputPath)

	yamlFiles, err := findYAMLFiles(*sourcePath)
	if err != nil {
		fmt.Printf("Error finding YAML files: %v\n", err)
		return
	}

	if len(yamlFiles) == 0 {
		fmt.Printf("No YML file found in %v\n", *sourcePath)
		return
	}

	fmt.Printf("Found %d files \n", len(yamlFiles))

	err = os.Setenv("KUBECONFIG", strings.Join(yamlFiles, ";"))
	if err != nil {
		fmt.Printf("Error setting environment variable: %v\n", err)
		return
	}

	cmd := exec.Command("kubectl", "config", "view", "--merge", "--flatten")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing kubectl command: %v\nOutput: %s", err, output)
		return
	}

	err = os.WriteFile(outputPath, []byte(output), 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		return
	}

	fmt.Printf("Generated %v \n", outputPath)

}