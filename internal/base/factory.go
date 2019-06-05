package base

import (
	"sync"

	"github.com/vmihailenco/taskq"
)

type Factory struct {
	queuesMu sync.RWMutex
	queues   []taskq.Queue
}

func (f *Factory) Add(q taskq.Queue) {
	f.queuesMu.Lock()
	f.queues = append(f.queues, q)
	f.queuesMu.Unlock()
}

func (f *Factory) Queues() []taskq.Queue {
	f.queuesMu.RLock()
	defer f.queuesMu.RUnlock()
	return f.queues
}

func (f *Factory) StartConsumers() error {
	return f.forEachQueue(func(q taskq.Queue) error {
		return q.Consumer().Start()
	})
}

func (f *Factory) CloseConsumers() error {
	return f.forEachQueue(func(q taskq.Queue) error {
		return q.Consumer().Close()
	})
}

func (f *Factory) Close() error {
	return f.forEachQueue(func(q taskq.Queue) error {
		return q.Close()
	})
}

func (f *Factory) forEachQueue(fn func(taskq.Queue) error) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	for _, q := range f.Queues() {
		wg.Add(1)
		go func(q taskq.Queue) {
			defer wg.Done()
			err := fn(q)
			select {
			case errCh <- err:
			default:
			}
		}(q)
	}
	wg.Wait()
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
