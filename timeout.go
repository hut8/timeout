package timeout

import (
	"context"
	"sync"
	"time"
)

type Timeout struct {
	d      time.Duration
	timer  *time.Timer
	ctx    context.Context
	cancel context.CancelFunc
	once   sync.Once
	// about last time write was done
	lastTime  time.Time
	lastBytes int
}

func NewTimeout(d time.Duration, ctx context.Context) *Timeout {
	ctx, cancel := context.WithCancel(ctx)
	return &Timeout{
		d:      d,
		timer:  nil,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *Timeout) initialize() {
	t.timer = time.AfterFunc(t.d, t.doTimeout)
}

func (t *Timeout) doTimeout() {
	t.cancel()
}

func (t *Timeout) Cancel() {
	t.timer.Stop()
}

func (t *Timeout) Write(data []byte) (int, error) {
	t.once.Do(t.initialize)
	t.timer.Reset(t.d)
	t.lastTime = time.Now()
	t.lastBytes = len(data)
	return len(data), nil
}
