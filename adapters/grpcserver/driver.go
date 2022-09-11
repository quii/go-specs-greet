package grpcserver

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string

	connectionOnce sync.Once
	conn           *grpc.ClientConn
}

func (d *Driver) Greet(name string) (string, error) {
	conn, err := d.getConnection()
	if err != nil {
		return "", err
	}

	client := NewGreeterClient(conn)
	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}

func (d *Driver) Close() {
	if d.conn != nil {
		d.conn.Close()
	}
}

func (d *Driver) getConnection() (*grpc.ClientConn, error) {
	var err error
	d.connectionOnce.Do(func() {
		d.conn, err = grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	return d.conn, err
}
