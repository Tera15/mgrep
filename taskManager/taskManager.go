package taskmanager

type Task struct {
	Path string
}
type TaskManager struct {
	Tasks chan Task
}

func (t *TaskManager) Add(task Task) {
	t.Tasks <- task
}

func (t *TaskManager) Get() Task {
	outTask := <-t.Tasks
	return outTask
}

func (t *TaskManager) Size() int {
	return len(t.Tasks)
}

func (t *TaskManager) NewTask(path string) Task {
	return Task{path}
}

func (t *TaskManager) Finalize(numWorkers int) {
	finalTask := t.NewTask("")
	for i := 0; i < numWorkers; i++ {
		t.Tasks <- finalTask
	}
}

func NewManager(buffer int) TaskManager {
	return TaskManager{make(chan Task, buffer)}
}
