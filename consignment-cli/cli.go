package main

import (
	pb "shippy/consignment-service/proto/consignment"
	"io/ioutil"
	"encoding/json"
	"errors"
	"google.golang.org/grpc"
	"log"
	"context"
	"os"
)

const(
	ADDRESS = "localhost:50055"
	DEFAULT_INFO_FILE = "consignment.json"
)

//读取consignment.json中记录的货物信息
func parseFile(filename string)(*pb.Consignment,error){
	data,err := ioutil.ReadFile(filename)
	if err !=nil {
		return nil,err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data,&consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main(){
	//连接到grpc服务器
	conn,err := grpc.Dial(ADDRESS,grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer conn.Close()

	//初始化grpc 客户端
	client := pb.NewShippingServiceClient(conn)

	//在命令行中指定新的货物信息 json文件
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}

	//解析货物信息
	consigment,err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error:%v",err)
	}

	//调用RPC
	//将货物存储到我们自己的仓库里
	resp,err := client.CreateConsignment(context.Background(),consigment)
	if err != nil {
		log.Fatalf("create consignment error:%v",err)
	}

	//新货物是否托运成功
	log.Printf("created: %t",resp.Created)

	// 列出目前所有托运的货物
	resp, err = client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments: %v", err)
	}
	for _, c := range resp.Consignments {
		log.Printf("%+v", c)
	}



}