package httplayer

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	pb "github.com/monsterr00/metric-service.gittest_client/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	gaugeMetricType   = "gauge"
	counterMetricType = "counter"
)

type httpAPI struct {
	client      *resty.Client
	app         applayer.App
	workersPool *Pool
	command     chan string
	wg          *sync.WaitGroup
	grpcConn    *grpc.ClientConn
	grpcClient  pb.MetricsClient
}

// New инициализирует уровень app
func New(appLayer applayer.App) *httpAPI {
	api := &httpAPI{
		client:      resty.New(),
		app:         appLayer,
		workersPool: NewPool(),
		command:     make(chan string),
		wg:          &sync.WaitGroup{},
		grpcConn:    nil,
		grpcClient:  nil,
	}

	api.setupClient()
	return api
}

// setupClient устанавливает настройки http-клиента.
func (api *httpAPI) setupClient() {
	api.client.
		SetRetryCount(3).
		SetRetryWaitTime(30 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second)
}

// Engage запускает сбор метрик и другие службы приложения.
func (api *httpAPI) Engage() {
	api.generateCryptoKeys()

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		api.stopServer()
		close(idleConnsClosed)
	}()

	api.startServer()
	api.startServices()

	<-idleConnsClosed
	log.Printf("Client Shutdown gracefully start")
	api.stopServices()
	log.Printf("Client Shutdown gracefully end")
}

// compress сжимает тело запроса.
func (api *httpAPI) compress(body string) (string, error) {
	var err error
	var buf bytes.Buffer
	b := []byte(body)

	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if _, err = gz.Write(b); err != nil {
		log.Printf("Client: compress error: %s\n", err)
		return buf.String(), err
	}
	if err = gz.Close(); err != nil {
		log.Printf("Client: gzip close error: %s\n", err)
		return buf.String(), err
	}

	return buf.String(), nil
}

// SetPrepBatch запускает работу функции пакетной отправки запросов prepBatch
func (api *httpAPI) SetPrepBatch() {
	defer api.wg.Done()
	var status = "work"

	for {
		select {
		case cmd := <-api.command:
			switch cmd {
			case "stop":
				return
			default:
				status = "work"
			}
		default:
			if status == "work" {
				switch config.ClientOptions.GrpcOn {
				case true:
					api.prepBatchGrpc()
				case false:
					api.prepBatch()
				}

				time.Sleep(time.Duration(config.ClientOptions.ReportInterval) * time.Second)
			}
		}
	}
}

// prepBatch разбивает массив отправляемых данных по метрикам на пакеты.
func (api *httpAPI) prepBatch() {
	api.app.LockRW()
	metrics, err := api.app.Metrics()
	api.app.UnlockRW()

	if err != nil {
		log.Printf("Client: error getting gauge metrics %s\n", err)
	}

	var body = "["
	var counter int64

	for _, v := range metrics {
		switch v.MType {
		case gaugeMetricType:
			body += fmt.Sprintf(`{"id":"%s","type":"%s","value":%f},`, v.ID, v.MType, *v.Value)
			counter += 1
		case counterMetricType:
			body += fmt.Sprintf(`{"id":"%s","type":"%s","delta":%d},`, v.ID, v.MType, *v.Delta)
			counter += 1
		}

		if counter == config.ClientOptions.BatchSize {
			if len(body) > 1 {
				originalBody := body[:len(body)-1]
				originalBody += "]"
				api.sendReqToChan(originalBody)
				counter = 0
				body = "["
			}
		}
	}

	//отправляем остатки
	if len(body) > 1 {
		originalBody := body[:len(body)-1]
		originalBody += "]"
		api.sendReqToChan(originalBody)
	}
}

// prepBatchGrpc разбивает массив отправляемых данных по метрикам на пакеты (gRPC).
func (api *httpAPI) prepBatchGrpc() {
	api.app.LockRW()
	metrics, err := api.app.Metrics()
	api.app.UnlockRW()

	if err != nil {
		log.Printf("Client: error getting gauge metrics %s\n", err)
	}

	for _, metric := range metrics {
		var metricGrpc pb.Metric

		metricGrpc.Id = metric.ID
		metricGrpc.MType = metric.MType

		switch metricGrpc.MType {
		case gaugeMetricType:
			metricGrpc.Value = *metric.Value
		case counterMetricType:
			metricGrpc.Delta = *metric.Delta
		}

		if api.grpcClient != nil {
			md := metadata.New(map[string]string{"X-Real-IP": api.realIP()})
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			resp, err := api.grpcClient.AddMetric(ctx, &pb.AddMetricRequest{
				Metric: &metricGrpc,
			})

			if err != nil {
				log.Printf("Client: %s\n", err)
			}
			if resp.Error != "" {
				log.Printf(resp.Error)
			}

			log.Printf("Req sended \n")
		}
	}
}

// sendReqToChan подготовалливает post-запрос и отправляет его в фабрику.
func (api *httpAPI) sendReqToChan(originalBody string) {
	compressedBody, err := api.compress(originalBody)
	if err != nil {
		log.Printf("Client: compress error: %s\n", err)
	}

	req := api.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("HashSHA256", api.signBody(originalBody)).
		SetHeader("X-Real-IP", api.realIP()).
		SetBody(api.encrypt(compressedBody))

	api.workersPool.Add(req)
}

// signBody подписывает тело запроса алгоритмом HMAC, используя SHA-256.
func (api *httpAPI) signBody(body string) string {
	if config.ClientOptions.SignMode {
		h := hmac.New(sha256.New, []byte(config.ClientOptions.Key))
		h.Write([]byte(body))
		return hex.EncodeToString(h.Sum(nil))
	}
	return ""
}

// generateCryptoKeys загружает ключи шифрования из файла или генерирует их
func (api *httpAPI) generateCryptoKeys() {
	filePub, err := os.OpenFile(config.ClientOptions.PublicKeyPath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer filePub.Close()

	scanner := bufio.NewScanner(filePub)
	scanner.Scan()
	publicKeyFile := scanner.Bytes()

	if len(publicKeyFile) > 0 {
		publicKeyFile, err = os.ReadFile(config.ClientOptions.PublicKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		publicKeyBlock, _ := pem.Decode(publicKeyFile)
		publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
		if err != nil {
			log.Fatal(err)
		}

		config.ClientOptions.PublicCryptoKey = publicKey

	} else {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		filePriv, err := os.OpenFile(config.ClientOptions.PrivateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}

		writer := bufio.NewWriter(filePriv)

		if _, err := writer.Write(privateKeyPEM); err != nil {
			log.Fatal(err)
		}
		if err := writer.Flush(); err != nil {
			log.Fatal(err)
		}

		err = filePriv.Close()
		if err != nil {
			log.Fatal(err)
		}

		publicKey := &privateKey.PublicKey
		config.ClientOptions.PublicCryptoKey = publicKey

		publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		})

		writer = bufio.NewWriter(filePub)
		if _, err := writer.Write(publicKeyPEM); err != nil {
			log.Fatal(err)
		}
		if err := writer.Flush(); err != nil {
			log.Fatal(err)
		}
	}
}

// encrypt используется для шифрования исходящих запросов
func (api *httpAPI) encrypt(body string) string {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, config.ClientOptions.PublicCryptoKey, []byte(body))
	if err != nil {
		log.Fatal(err)
	}
	return string(ciphertext)
}

// stopServices останавливает работу всех сервисов агента
func (api *httpAPI) stopServices() {
	api.app.StopMetricGen()
	api.stopReqSend()
}

// stopServer останавливает работу серверов
func (api *httpAPI) stopServer() {
	switch config.ClientOptions.GrpcOn {
	case true:
		api.grpcConn.Close()
	case false:
		api.workersPool.Stop()
	}
}

// startServer запускает работу серверов
func (api *httpAPI) startServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch config.ClientOptions.GrpcOn {
	case true:
		var err error
		api.grpcConn, err = grpc.NewClient(config.ClientOptions.GrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}

		api.grpcClient = pb.NewMetricsClient(api.grpcConn)
	case false:
		api.workersPool.Run(ctx)
	}
}

// startServices запускает работу всех сервисов агента
func (api *httpAPI) startServices() {
	go api.app.SetMetrics()
	go api.SetPrepBatch()
}

// stopReqSend останавливает работу отправки запросов
func (api *httpAPI) stopReqSend() {
	api.wg.Add(1)
	api.command <- "stop"
	api.wg.Wait()
}

// realIP возвращает ip-адрес агента
func (api *httpAPI) realIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("Client: can't get IP address")
	}

	var ip net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println("Client: can't get IP interface addresses")
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil && i.Name == "en0" {
				ip = ipNet.IP
			}
		}
	}

	return ip.String()
}
