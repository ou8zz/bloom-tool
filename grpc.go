package service

import (
	"context"
	"time"

	"github.com/ou8zz/bloom-tool/pb"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var rpcConn *grpc.ClientConn
var bloomClient pb.BloomServiceClient

func InitGRpcClient(addrs []string) {
	var lastErr error
	for _, addr := range addrs {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			lastErr = err
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn.Connect()
		for {
			state := conn.GetState()
			if state == connectivity.Ready {
				break
			}
			if !conn.WaitForStateChange(ctx, state) {
				break
			}
		}
		cancel()

		if conn.GetState() == connectivity.Ready {
			rpcConn = conn
			bloomClient = pb.NewBloomServiceClient(conn)
			log.Printf("grpc client success, state:%s, addr:%s", rpcConn.GetState().String(), rpcConn.CanonicalTarget())
			return
		}

		_ = conn.Close()
		lastErr = context.DeadlineExceeded
	}

	if lastErr != nil {
		log.Fatal(lastErr)
	}
}

func CloseGRpcClient() {
	if rpcConn != nil {
		rpcConn.Close()
	}
}

func Add(key string) {
	if bloomClient == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := bloomClient.Add(ctx, &pb.ItemRequest{Key: key})
	if err != nil {
		log.Printf("Add key %s error: %v", key, err)
		return
	}
}

func Exists(key string) bool {
	if bloomClient == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := bloomClient.Exists(ctx, &pb.ItemRequest{Key: key})
	if err != nil {
		log.Printf("Check key %s error: %v", key, err)
		return false
	}
	return resp.GetExists()
}

func Load(key string) string {
	reply := ""
	if bloomClient == nil {
		return reply
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 单个请求连接超时
	defer cancel()

	resp, err := bloomClient.Load(ctx, &pb.MapRequest{Key: key})
	if err != nil {
		log.Printf("Key %s error: %v", key, err)
		return reply
	}
	return resp.GetValue()
}

func Store(key, val string) {
	if bloomClient == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 单个请求连接超时
	defer cancel()

	resp, err := bloomClient.Store(ctx, &pb.MapRequest{Key: key, Value: val})
	if err != nil || !resp.Exists {
		log.Printf("Key %s value %s error: %v", key, val, err)
		return
	}
}
