package main

import (
	"fmt"
	proto "github.com/zhaoxin-BF/proto-apis/golang/stream/v1"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type StreamServices struct {
	proto.UnimplementedStreamServiceServer
}

func main() {
	server := grpc.NewServer()

	proto.RegisterStreamServiceServer(server, &StreamServices{})

	lis, err := net.Listen("tcp", "127.0.0.1:9528")
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}

func (services *StreamServices) OrderList(in *proto.OrderSearchParams, resp proto.StreamService_OrderListServer) error {
	for i := 0; i <= 10; i++ {
		order := proto.Order{
			Id:      int32(i),
			OrderSn: time.Now().Format("20060102150405") + "order_sn",
			Date:    time.Now().Format("2006-01-02 15:04:05"),
		}
		err := resp.Send(&proto.OrderListResponse{
			Order: &order,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (services *StreamServices) UploadFile(resp proto.StreamService_UploadFileServer) error {
	for {
		res, err := resp.Recv()
		//接收消息结束，发送结果，并关闭
		if err == io.EOF {
			return resp.SendAndClose(&proto.UploadResponse{
				Msg:     "success",
				Retcode: "0",
			})
		}
		if err != nil {
			return err
		}
		fmt.Println(res)
	}
	return nil
}

func (services *StreamServices) SumData(resp proto.StreamService_SumDataServer) error {
	i := 0
	sum := 0
	for {
		//time.Sleep(1 * time.Second)
		err := resp.Send(&proto.SumDataResponse{Result: int32(sum)})
		if err != nil {
			return err
		}
		res, err := resp.Recv()
		if err == io.EOF {
			return nil
		}
		sum += int(res.Number)
		log.Printf("res:%d, step:%d,sum:%d\r\n", res.Number, i, sum)
		i++
	}
	return nil
}
