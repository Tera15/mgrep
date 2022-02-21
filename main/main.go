package main

import (
	"flag"
	"fmt"
	task "mgrep/task"
	taskmanager "mgrep/taskManager"
	"os"
	"sync"
)

func getDirectory(t *taskmanager.TaskManager, path string) {
	dir, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("error opening dir at path: ", path, "\n", err)
		return
	}

	for _, entry := range dir {
		fmt.Println("hi")
		newPath := path + "/" + entry.Name()
		fmt.Println(entry.Type())
		if entry.IsDir() {
			getDirectory(t, newPath)
		} else {
			np := t.NewTask(newPath)
			t.Add(np)
		}
	}
}

func main() {

	flag.Parse()
	searchTerm := flag.Arg(1)
	path := flag.Arg(2)

	tm := taskmanager.NewManager(100)
	rChan := make(chan task.Result, 112)

	var workerWg sync.WaitGroup
	var displayWg sync.WaitGroup

	workers := 10

	workerWg.Add(1)
	go func() {
		defer workerWg.Done()
		getDirectory(&tm, path)
		tm.Finalize(workers)
	}()

	for i := 0; i < workers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for {
				workEntry := tm.Get()
				if workEntry.Path != "" {
					workerResult := task.ReadFileContents(workEntry.Path, searchTerm)
					if workerResult != nil {
						for _, r := range workerResult.List {
							rChan <- r
						}
					}
				} else {
					return
				}
			}
		}()
	}

	blockWorkersWg := make(chan struct{})
	go func() {
		workerWg.Wait()
		close(blockWorkersWg)
	}()

	displayWg.Add(1)
	go func() {
		for {
			select {
			case r := <-rChan:
				fmt.Printf("Line: %v\nLine number:%v\nPath: %v\n", r.Line, r.LineNumber, r.Path)
			case <-blockWorkersWg:
				if len(rChan) == 0 {
					fmt.Println("DONE BITCH")
					displayWg.Done()
					return
				}
			}
		}
	}()
	displayWg.Wait()
}
