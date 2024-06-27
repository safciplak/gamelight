package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Data struct {
	Name string `json:"name"`
}

func main() {
	dir := "./json_files"

	filePaths := make(chan string)

	var wg sync.WaitGroup

	numWorkers := 5

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(filePaths, &wg)
	}

	go func() {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("error accessing path %q: %v\n", path, err)
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".json" {
				filePaths <- path
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error walking the path %q: %v\n", dir, err)
		}
		close(filePaths)
	}()

	wg.Wait()
}

func worker(filePaths chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range filePaths {
		processFile(path)
	}
}

func processFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error reading file %q: %v\n", path, err)
		return
	}

	var data Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Printf("error unmarshalling JSON from file %q: %v\n", path, err)
		return
	}

	fmt.Printf("Name from file %s: %s\n", path, data.Name)
}
