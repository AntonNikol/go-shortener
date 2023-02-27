package workers

import (
	"context"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"log"
	"sync"
	"time"
)

const (
	workerCount   = 10
	bufferMaxSize = 100
)

type task struct {
	userID  string
	listIDS []string
}

type WorkerPool struct {
	ch      chan task
	ctx     context.Context
	mu      sync.Mutex
	buffer  []task
	storage repositories.Repository
}

func NewBatchPostponeRemover(ctx context.Context, storage repositories.Repository, period time.Duration, chCap int) *WorkerPool {
	w := WorkerPool{
		ch:      make(chan task, chCap),
		ctx:     ctx,
		storage: storage,
	}

	flushCh := make(chan struct{})
	for i := 0; i < workerCount; i++ {
		go func() {
			time.Sleep(1 * time.Second)
			for task := range w.ch {
				w.addToBuffer(task, flushCh)
			}
		}()
	}

	go func() {
		ticker := time.NewTicker(period)
		for {
			select {
			case <-flushCh:
				log.Printf("case 1")
				w.removeBatchStorage()
			case <-ticker.C:
				log.Printf("case 2")
				w.removeBatchStorage()
			case <-w.ctx.Done():
				log.Printf("case 3")
				ticker.Stop()
				w.removeBatchStorage()
			}
		}
	}()

	return &w
}

func (w *WorkerPool) Remove(userID string, listIDS []string) {
	log.Printf("Remove add task to channel")
	w.ch <- task{userID: userID, listIDS: listIDS}
}

func (w *WorkerPool) addToBuffer(t task, flushCh chan struct{}) {
	w.mu.Lock()
	log.Printf("addToBuffer: %+v", t)
	w.buffer = append(w.buffer, t)
	if len(w.buffer) >= bufferMaxSize {
		flushCh <- struct{}{}
	}

	w.mu.Unlock()
}

func (w *WorkerPool) removeBatchStorage() {
	w.mu.Lock()
	log.Printf("removeBatchStorage, buffer len: %v", len(w.buffer))
	for _, task := range w.buffer {

		log.Printf("removeBatchStorage, delete task: %v", task)

		err := w.storage.Delete(w.ctx, task.listIDS, task.userID)
		if err != nil {
			log.Printf("unable delete itemsIDS %v", err)
		}
		log.Printf("removeBatchStorage: user %s %v", task.listIDS, task.userID)
	}
	log.Printf("removeBatchStorage purgeBuffer")
	w.purgeBuffer()
	w.mu.Unlock()
}

func (w *WorkerPool) purgeBuffer() {
	w.buffer = make([]task, 0)
}
