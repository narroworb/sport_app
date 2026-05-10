package analytics

import (
	"context"
	"fmt"
	"os"
	"time"

	analyticsv1 "github.com/narroworb/core_api/gen/analytics/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultAddr = "analytics_service:50051"
)

type Client struct {
	conn   *grpc.ClientConn
	client analyticsv1.AnalyticsServiceClient
}

func NewClient() (*Client, error) {
	addr := os.Getenv("ANALYTICS_GRPC_ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("dial analytics grpc %q: %w", addr, err)
	}

	return &Client{
		conn:   conn,
		client: analyticsv1.NewAnalyticsServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) Raw() analyticsv1.AnalyticsServiceClient {
	return c.client
}
