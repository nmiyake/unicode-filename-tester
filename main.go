package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
)

func main() {
	verbose := flag.Bool("v", false, "print verbose output")
	flag.Parse()

	success, err := runTest(*verbose)
	if err != nil {
		if *verbose {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	if !success {
		if *verbose {
			fmt.Println("Failed: Unicode file names were normalized")
		}
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("Success: Unicode file names were not normalized")
	}
}

func runTest(verbose bool) (bool, error) {
	tmpDir, err := ioutil.TempDir(".", "")
	if err != nil {
		return false, fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Printf("failed to remove %s in defer: %v\n", tmpDir, err)
		}
	}()

	// ö.txt
	composed := string('\u00F6') + ".txt"
	if err := ioutil.WriteFile(path.Join(tmpDir, composed), []byte("composed"), 0755); err != nil {
		return false, fmt.Errorf("failed to write composed: %v", err)
	}

	// o¨.txt
	decomposed := string('\u006F') + string('\u0308') + ".txt"
	if err := ioutil.WriteFile(path.Join(tmpDir, decomposed), []byte("decomposed"), 0755); err != nil {
		return false, fmt.Errorf("failed to write decomposed: %v", err)
	}

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return false, fmt.Errorf("failed to list entries in directory %s: %v", tmpDir, err)
	}
	fileNames := make([]string, len(files))
	for i := range files {
		fileNames[i] = files[i].Name()
	}
	if verbose {
		fmt.Printf("Number of files:\n")
		fmt.Printf("\tExpected: %d\n", 2)
		fmt.Printf("\tGot:      %d\n", len(files))

		fmt.Printf("Files:\n")
		fmt.Printf("\tExpected: %v\n", []string{decomposed, composed})
		fmt.Printf("\tGot:      %v\n", fileNames)
	}

	composedBytes, err := ioutil.ReadFile(path.Join(tmpDir, composed))
	if err != nil {
		return false, fmt.Errorf("failed to read composed: %v", err)
	}
	if verbose {
		fmt.Printf("Content of %s (\\u00F6.txt):\n", composed)
		fmt.Printf("\tExpected: %s\n", "composed")
		fmt.Printf("\tGot:      %s\n", string(composedBytes))
	}

	decomposedBytes, err := ioutil.ReadFile(path.Join(tmpDir, decomposed))
	if err != nil {
		return false, fmt.Errorf("failed to read decomposed: %v", err)
	}
	if verbose {
		fmt.Printf("Content of %s (\\u006F\\u0308.txt):\n", decomposed)
		fmt.Printf("\tExpected: %s\n", "decomposed")
		fmt.Printf("\tGot:      %s\n", string(decomposedBytes))
	}

	return !reflect.DeepEqual(composedBytes, decomposedBytes), nil
}
