package main

import (
	"fmt"
	proto "github.com/zhaoxin-BF/proto-apis/golang/stream/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	for i := 0; i <= 1000; i++ {
		order := proto.Order{
			Id:      int32(i),
			OrderSn: time.Now().Format("20060102150405") + "order_sn",
			Date:    time.Now().Format("2006-01-02 15:04:05"),
		}
		err := resp.Send(&proto.OrderListResponse{
			Order: &order,
		})
		if err != nil {
			fmt.Println("err send msg~")
			return err
		}
		time.Sleep(time.Second)
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
	md, ok := metadata.FromIncomingContext(resp.Context())
	if ok {
		// 从元数据中获取设备ID
		deviceID := md.Get("x-everai-device-id")
		if len(deviceID) > 0 {
			fmt.Println("Device ID:", deviceID[0])
		}
	}
	for {
		//time.Sleep(2 * time.Second)
		err := resp.Send(&proto.SumDataResponse{Result: 1})
		if err != nil {
		}
		//time.Sleep(1 * time.Second)

		res, err := resp.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			return nil
		}
		fmt.Printf("%+v\n", res)
	}
}
