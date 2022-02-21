package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Line       string
	LineNumber int
	Path       string
}
type Results struct {
	List []Result
}

func ReadFileContents(path string, target string) *Results {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file", err)
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	fmt.Println(scanner)
	lineNumber := 1
	results := Results{make([]Result, 0)}

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), target) {
			res := Result{scanner.Text(), lineNumber, path}
			results.List = append(results.List, res)
		}
		lineNumber += 1
	}

	return &results
}
