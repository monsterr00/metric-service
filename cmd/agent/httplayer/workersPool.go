package httplayer

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
)

type Pool struct {
	wg     *sync.WaitGroup
	queue  chan *resty.Request
	errors chan error
}

func NewPool() *Pool {
	return &Pool{
		wg:     &sync.WaitGroup{},
		queue:  make(chan *resty.Request, config.ClientOptions.RateLimit),
		errors: make(chan error, config.ClientOptions.RateLimit),
	}
}

func (p *Pool) Run(ctx context.Context) {
	for i := 0; i < int(config.ClientOptions.PoolWorkers); i++ {
		go p.doWork()
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("Client: channel is closed: Done")
			return
		case err, ok := <-p.errors:
			if !ok {
				log.Printf("Client: channel is closed")
				return
			}
			log.Printf("Client: error from channel: %s\n", err)
		}
	}
}

func (p *Pool) doWork() {
	for r := range p.queue {
		requestURL := fmt.Sprintf("%s%s%s", "http://", config.ClientOptions.Host, "/updates/")
		req, err := r.Post(requestURL)

		log.Printf("Req status: %d\n", req.StatusCode())

		if err != nil {
			p.errors <- err
		}
	}
}

func (p *Pool) Add(r *resty.Request) {
	p.queue <- r
}

func (p *Pool) Stop() {
	close(p.queue)
	close(p.errors)
}
