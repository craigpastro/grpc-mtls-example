package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const echoScope = "echo:true"

func Authenticate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if p, ok := peer.FromContext(ctx); ok {
		if mtls, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			for _, item := range mtls.State.PeerCertificates {
				for _, organization := range item.Subject.Organization {
					if organization == echoScope {
						log.Println("authenticated")
						return handler(ctx, req)
					}
				}
			}
		}
	}

	log.Println("unauthentication")
	return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
}
