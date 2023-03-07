package wclient

import (
	"context"
	"crawler-distributed/model"
	"crawler-distributed/pb"
	"crawler-distributed/worker/parser"
	"fmt"
)

func CreateProcessor(ctx context.Context, clientChan *chan pb.CrawlServiceClient) model.Processor {

	return func(request model.Request) (model.ParserResult, error) {
		grpcClient := <-*clientChan
		sRequest := parser.SerializeRequest(request)
		parserResult, err := grpcClient.Process(ctx, sRequest)
		if err != nil {
			return model.ParserResult{}, fmt.Errorf("[worker.client.CreateProcessor] error call worker server: %v", err)
		}

		result := parser.DeserializeParserResult(parserResult)
		return result, nil
	}
}
