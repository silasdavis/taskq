package processor

import (
	"time"

	"gopkg.in/msgqueue.v1"
)

type Queuer interface {
	Name() string
	Processor() *Processor
	Add(msg *msgqueue.Message) error
	Call(args ...interface{}) error
	CallOnce(dur time.Duration, args ...interface{}) error
	ReserveN(n int) ([]msgqueue.Message, error)
	Release(*msgqueue.Message, time.Duration) error
	Delete(msg *msgqueue.Message) error
	DeleteBatch(msg []*msgqueue.Message) error
	Purge() error
}
