package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"

	"github.com/craigpastro/mtls-example/middleware"
	pb "github.com/craigpastro/mtls-example/protos/api/v1"
	"github.com/craigpastro/mtls-example/server"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	caCert     = "certs/ca.pem"
	serverCert = "certs/server.pem"
	serverKey  = "certs/server-key.pem"
)

type Config struct {
	ServerAddr string `split_words:"true" default:"localhost:9090"`
}

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal("error reading config", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	run(ctx, config)
}

func run(ctx context.Context, config Config) {
	creds := credentials.NewTLS(configureTLS())
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(middleware.Authenticate),
	)
	pb.RegisterServiceServer(s, server.NewServer())

	lis, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting server on %s", config.ServerAddr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func configureTLS() *tls.Config {
	data, err := ioutil.ReadFile(caCert)
	if err != nil {
		log.Fatalf("error reading '%s': %v", caCert, err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(data) {
		log.Fatal("error adding ca cert to cert pool")
	}

	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatalf("error loading '%s' and '%s': %v", serverCert, serverKey, err)
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13,
	}
}
