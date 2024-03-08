package main

import (
	"context"
	"github.com/gin-gonic/gin"
	proto "github.com/zhaoxin-BF/proto-apis/golang/stream/v1"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"strconv"
)

var streamClient proto.StreamServiceClient

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

func orderList(ctx *gin.Context) {
	stream, err := streamClient.OrderList(context.Background(), &proto.OrderSearchParams{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"orders": res})
		log.Println(res)
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
	stream, err := streamClient.SumData(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := 1; i <= 10; i++ {
		err = stream.Send(&proto.SumDataRequest{Number: int32(i)})
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		log.Printf("res number:%d", res.Result)
	}
	stream.CloseSend()
	return
}
