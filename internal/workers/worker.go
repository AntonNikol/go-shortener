package workers

import (
	"context"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"log"
	"sync"
	"time"
)

type task struct {
	userID  string
	listIDS []string
}

type Job struct {
	ch      chan task
	ctx     context.Context
	mu      sync.Mutex
	buffer  map[string][]string
	storage *repositories.Repository
}

func NewBatchPostponeRemover(ctx context.Context, storage *repositories.Repository, period time.Duration, chCap int) *Job {
	j := Job{
		ch:      make(chan task, chCap),
		ctx:     ctx,
		storage: storage,
	}

	j.purgeBuffer()

	go j.init(period)
	return &j
}

func (j *Job) init(period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		select {
		case req := <-j.ch:
			j.addToBuffer(req)
		case <-ticker.C:
			j.removeBatchStorage()
		case <-j.ctx.Done():
			ticker.Stop()
			close(j.ch)
			j.removeBatchStorage()
		}
	}
}

func (j *Job) Remove(userID string, listIDS []string) {
	j.ch <- task{userID: userID, listIDS: listIDS}
}

func (j *Job) addToBuffer(t task) {
	j.mu.Lock()

	for _, v := range t.listIDS {
		log.Printf("addToBuffer: userId %s link id %s", t.userID, v)
		j.buffer[t.userID] = append(j.buffer[t.userID], v)
	}
	j.mu.Unlock()
}

func (j *Job) removeBatchStorage() {
	j.mu.Lock()
	for userID, linkList := range j.buffer {
		err := (*j.storage).Delete(j.ctx, linkList, userID)
		if err != nil {
			log.Printf("unable delete itemsIDS %v", err)
		}
		log.Printf("removeBatchStorage: user %s %v", userID, linkList)
	}
	j.purgeBuffer()
	j.mu.Unlock()
}

func (j *Job) purgeBuffer() {
	j.buffer = make(map[string][]string)
}
