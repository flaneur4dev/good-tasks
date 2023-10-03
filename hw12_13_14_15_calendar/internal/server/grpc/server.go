package grpcserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/server/grpc/pb"
)

type Server struct {
	port string
	srv  *grpc.Server
	fd   *os.File
}

func New(app application, logPath, port string) (*Server, error) {
	fd, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	s := &Server{port: port, fd: fd}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(s.loggingMiddleware),
	)
	pb.RegisterCalendarServer(srv, newCalendarServer(app))

	s.srv = srv
	return s, nil
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	return s.srv.Serve(l)
}

func (s *Server) Stop() error {
	s.srv.GracefulStop()
	return s.fd.Close()
}

func (s *Server) loggingMiddleware(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	t := time.Now()

	res, err := handler(ctx, req)

	log := fmt.Sprintf("[%s] %s %d",
		t.String(),
		info.FullMethod,
		time.Since(t).Milliseconds(),
	)

	_, werr := s.fd.Write([]byte(log + "\n"))
	if werr != nil {
		fmt.Println("failed to write to logfile: " + werr.Error())
	}

	return res, err
}
