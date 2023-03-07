package service

import (
	"bytes"
	"context"
	"crawler-distributed/pb"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"google.golang.org/grpc/codes"
	"log"
)

type ItemSaverService struct {
	client *elasticsearch.Client
	Index  string
	pb.UnimplementedItemSaverServiceServer
}

func NewItemServer(esClient *elasticsearch.Client, index string) *ItemSaverService {
	return &ItemSaverService{
		client: esClient,
		Index:  index,
	}
}

func (server *ItemSaverService) Save(ctx context.Context, req *pb.ItemSaverRequest) (*pb.ItemSaverResponse, error) {
	log.Println("[service.Save] receive a save-item request")

	item := req.Item
	err := SaveItemToEs(server.client, item, server.Index)
	if err != nil {
		return &pb.ItemSaverResponse{Msg: codes.Internal.String()}, errors.New(fmt.Sprintf("[service.Save] error save item to elasticsearch: %v", err))
	}
	return &pb.ItemSaverResponse{Msg: codes.OK.String()}, nil
}

// SaveItemToEs 将数据存储到 ElasticSearch 中
func SaveItemToEs(client *elasticsearch.Client, item *pb.Item, index string) error {
	data, err := json.Marshal(item)
	if err != nil {
		log.Fatalf("[SaveItemToEs] An error Occured %v", err)
	}

	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}
	// 当获取的「用户ID」不为空时，使用「用户ID」作为 DocumentID
	if item.Id != "" {
		req.DocumentID = item.Id
	}

	// 发送并执行请求
	resp, err := req.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("An error Occured %v", err)
	}
	defer resp.Body.Close()
	if resp.IsError() {
		log.Printf("[SaveItemToEs] [%s] Error", resp.Status())
		return err
	}

	// 解析响应结果并打印到控制台
	var r map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return err
	}
	return nil
}
