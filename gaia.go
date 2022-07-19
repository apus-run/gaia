package gaia

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/apus-run/gaia/v1/log"
	"github.com/apus-run/gaia/v1/registry"
	"github.com/apus-run/gaia/v1/transport"
)

type App interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

type Gaia struct {
	opts     options
	ctx      context.Context
	cancel   func()
	mu       sync.Mutex
	instance *registry.ServiceInstance
}

// New create an application lifecycle manager.
func New(opts ...Option) *Gaia {
	o := defaultOptions()

	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}

	if o.logger != nil {
		log.SetLogger(o.logger)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &Gaia{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

// ID returns app instance id.
func (a *Gaia) ID() string { return a.opts.id }

// Name returns service name.
func (a *Gaia) Name() string { return a.opts.name }

// Version returns app version.
func (a *Gaia) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *Gaia) Metadata() map[string]string { return a.opts.metadata }

// Endpoint returns endpoints.
func (a *Gaia) Endpoint() []string {
	if a.instance != nil {
		return a.instance.Endpoints
	}
	return nil
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *Gaia) Run() error {
	// build service instance
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()
	eg, ctx := errgroup.WithContext(NewContext(a.ctx, a))
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx, cancel := context.WithTimeout(NewContext(a.opts.ctx, a), a.opts.stopTimeout)
			defer cancel()
			return srv.GracefullyStop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(NewContext(a.opts.ctx, a))
		})
	}
	wg.Wait()

	// register service
	if a.opts.registry != nil {
		c, cancel := context.WithTimeout(ctx, a.opts.registryTimeout)
		defer cancel()
		if err := a.opts.registry.Register(c, instance); err != nil {
			return err
		}
	}

	// watch signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-quit:
			return a.Stop()
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Stop gracefully stops the application.
func (a *Gaia) Stop() error {
	// deregister instance
	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()
	if a.opts.registry != nil && instance != nil {
		ctx, cancel := context.WithTimeout(NewContext(a.ctx, a), a.opts.registryTimeout)
		defer cancel()
		if err := a.opts.registry.Deregister(ctx, instance); err != nil {
			return err
		}
	}

	// cancel app
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func (a *Gaia) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0, len(a.opts.endpoints))
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range a.opts.servers {
			if r, ok := srv.(transport.Endpointer); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, e.String())
			}
		}
	}
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s App) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s App, ok bool) {
	s, ok = ctx.Value(appKey{}).(App)
	return
}
