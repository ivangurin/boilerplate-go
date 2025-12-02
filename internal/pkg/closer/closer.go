package closer

import (
	"context"
	"os"
	"os/signal"
	"sync"

	logger_pkg "boilerplate/internal/pkg/logger"
)

type Closer interface {
	Add(f ...func() error)
	Wait()
	CloseAll()
}

type closer struct {
	ctx      context.Context
	logger   logger_pkg.Logger
	once     sync.Once
	done     chan struct{}
	funcs    []func() error
	shutdown chan os.Signal
	sync.Mutex
}

// os.Interrupt, syscall.SIGINT, syscall.SIGTERM
func New(ctx context.Context, logger logger_pkg.Logger, sig ...os.Signal) Closer {
	closer := &closer{
		ctx:      ctx,
		logger:   logger,
		done:     make(chan struct{}),
		shutdown: make(chan os.Signal, 1),
	}

	if len(sig) > 0 {
		go func() {
			signal.Notify(closer.shutdown, sig...)
			<-closer.shutdown
			signal.Stop(closer.shutdown)
			closer.logger.Infof(ctx, "graceful shutdown started...")
			defer closer.logger.Info(ctx, "graceful shutdown finished")
			closer.CloseAll()
		}()
	}

	return closer
}

func (c *closer) Add(f ...func() error) {
	c.Lock()
	defer c.Unlock()
	c.funcs = append(c.funcs, f...)
}

// nolint: revive
func (c *closer) Wait() {
	if r := recover(); r != nil {
		c.logger.Errorf(c.ctx, "panic while starting app: %v", r)
		c.CloseAll()
		return
	}
	<-c.done
}

func (c *closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.Lock()
		funcs := c.funcs
		c.Unlock()

		for i := len(funcs) - 1; i >= 0; i-- {
			err := c.funcs[i]()
			if err != nil {
				c.logger.Errorf(c.ctx, "close: %s", err.Error())
			}
		}
	})
}
