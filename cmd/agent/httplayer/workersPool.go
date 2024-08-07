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

// Run запускает фабрику.
func (p *Pool) Run(ctx context.Context) {
	for i := 0; i < int(config.ClientOptions.PoolWorkers); i++ {
		p.wg.Add(1)
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

	p.wg.Wait()
}

// doWork отправвляет запросы из очереди на сервер.
func (p *Pool) doWork() {
	defer p.wg.Done()
	for r := range p.queue {
		requestURL := fmt.Sprintf("%s%s%s", "http://", config.ClientOptions.Host, "/updates/")
		req, err := r.Post(requestURL)

		log.Printf("Req status: %d\n", req.StatusCode())

		if err != nil {
			p.errors <- err
		}
	}
}

// Add добавляет post-запрос в очередь.
func (p *Pool) Add(r *resty.Request) {
	p.queue <- r
}

// Stop останавливает работу фабрики.
func (p *Pool) Stop() {
	close(p.queue)
	close(p.errors)
}
