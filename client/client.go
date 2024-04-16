package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	proto "github.com/zhaoxin-BF/proto-apis/golang/stream/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var streamClient proto.StreamServiceClient
var (
	DeviceId = "x-everai-device-id"
)

func main() {
	r := gin.Default()
	r.GET("/testOrderList", orderList)
	r.GET("/testUploadImage", uploadImage)
	r.GET("/testSumData", sumData)
	r.Run(":8080")

}
func init() {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	streamClient = proto.NewStreamServiceClient(connect)
}

func GetConnection() {
	connect, err := grpc.Dial("https://everai.expvent.com.cn:1112", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("success : ", connect)
	streamClient = proto.NewStreamServiceClient(connect)

}

func orderList(ctx *gin.Context) {
	for {
		stream, err := streamClient.OrderList(context.Background(), &proto.OrderSearchParams{})
		if err != nil {
			fmt.Println("get streamClient err")
			time.Sleep(5 * time.Second)
			continue
		}
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				fmt.Println("get receive err")
				break
			}
			ctx.JSON(http.StatusOK, gin.H{"orders": res})
			log.Println(res)
		}
	}
}

func uploadImage(ctx *gin.Context) {
	stream, err := streamClient.UploadFile(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := 1; i <= 10; i++ {
		img := &proto.Image{FileName: "image" + strconv.Itoa(i), File: "file data"}
		images := &proto.ImageList{Image: img}
		err := stream.Send(images)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	//发送完毕 关闭并获取服务端返回的消息
	resp, err := stream.CloseAndRecv()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": resp, "message": "success"})
	log.Println(resp)

}

func sumData(ctx *gin.Context) {
	ctxx := SetDeviceIdIntoCtx(context.Background())
	for {
		stream, err := streamClient.SumData(ctxx)
		if err != nil {
			fmt.Println("get streamClient err")
			time.Sleep(5 * time.Second)
			continue
		}

		for i := 1; ; i++ {
			res, err := stream.Recv()
			if err == io.EOF {
				//break
			}
			if err != nil {
				break
			}
			log.Printf("res number:%d", res.Result)

			err = stream.Send(&proto.SumDataRequest{Number: int32(i)})
			if err == io.EOF {
				//break
			}
			if err != nil {
				break
			}
		}
		fmt.Println("request success!")
		stream.CloseSend()
	}
}

func SetDeviceIdIntoCtx(ctx context.Context) context.Context {
	deviceId := "80808080"
	md := metadata.New(map[string]string{
		DeviceId: deviceId,
	})
	ctx = metadata.NewIncomingContext(ctx, md)
	return metadata.NewOutgoingContext(ctx, md)
}
