package httplayer

import (
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

func (p *Pool) Run() {
	for i := 0; i < int(config.ClientOptions.PoolWorkers); i++ {
		p.wg.Add(1)
		go p.doWork()
	}

	for {
		select {
		case err := <-p.errors:
			log.Printf("Client: error from channel: %s\n", err)
		default:
			p.wg.Wait()
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
	p.wg.Done()
}

func (p *Pool) Add(r *resty.Request) {
	p.queue <- r
}
