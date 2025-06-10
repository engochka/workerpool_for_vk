package workerpool

import (
	"fmt"
	"sync"
)

type Pool struct {
	Tasks       []*Task
	Workers     []*Worker
	concurrency int
	collector   chan *Task
	quit        chan struct{}
	mu          sync.Mutex
	wg          sync.WaitGroup
}

// NewPool возвращает экземпляр Pool
func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		collector:   make(chan *Task, 1000),
		quit:        make(chan struct{}),
	}
}

// AddTask добавляет задачу в пул
func (p *Pool) AddTask(task *Task) {
	p.collector <- task
}

// AddWorker добавляет нового воркера
func (p *Pool) AddWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	workerID := len(p.Workers) + 1
	worker := NewWorker(p.collector, workerID)
	p.Workers = append(p.Workers, worker)
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		worker.StartBackground()
	}()
	fmt.Printf("Worker %d added\n", workerID)
}

// RemoveWorker удаляет последнего воркера
func (p *Pool) RemoveWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Workers) == 0 {
		fmt.Println("Должен быть хотя бы 1 воркер")
		return
	}

	lastWorker := p.Workers[len(p.Workers)-1]
	lastWorker.Stop()
	p.Workers = p.Workers[:len(p.Workers)-1]
	fmt.Printf("Worker %d removed\n", lastWorker.ID)
}

// RunBackground запускает бесконечный цикл в ожидании задачи и/или воркера
func (p *Pool) RunBackground() {
	for i := 1; i <= p.concurrency; i++ {
		p.AddWorker()
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	<-p.quit
}

// Stop останавливает пул
func (p *Pool) Stop() {
	close(p.quit)
	p.wg.Wait()
}
