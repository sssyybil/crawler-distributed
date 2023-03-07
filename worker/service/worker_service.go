package service

import (
	"bufio"
	"context"
	"crawler-distributed/config"
	"crawler-distributed/model"
	"crawler-distributed/pb"
	"crawler-distributed/worker/parser"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"time"
)

type CrawlService struct {
	pb.UnimplementedCrawlServiceServer
}

func NewCrawlService() *CrawlService {
	return &CrawlService{}
}

func (CrawlService) Process(ctx context.Context, req *pb.SerializedRequest) (*pb.SerializedParserResult, error) {
	// 解析 Request 请求
	request, err := parser.DeserializeRequest(req)
	if err != nil {
		return nil, err
	}

	parserResult, err := Worker(request)
	if err != nil {
		return nil, err
	}

	result := parser.SerializeParserResult(parserResult)
	return result, nil
}

/*
用于限制请求的速度，防止请求过快被限流，多个 worker 会抢这一个 rateLimiter，抢到后才能进行后续的 Get 请求。
也可以写成：var rateLimiter = time.Tick(100 * time.Millisecond)
*/
var rateLimiter = time.Tick(time.Second / config.Qps)

func Worker(r model.Request) (model.ParserResult, error) {
	log.Printf("💤 Fetching %s", r.Url)

	data, err := Fetch(r.Url)
	if err != nil {
		log.Printf("[worker.service.Worker] error fetching url %s, %v", r.Url, err)
		return model.ParserResult{}, nil
	}
	return r.Parser.Parse(data, r.Url), nil
}

func Fetch(url string) ([]byte, error) {
	<-rateLimiter

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[worker.service.Fetch] wrong status code: %d", resp.StatusCode)
	}

	// 获取网页编码，并将其转换成 UTF-8 编码
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	data, err := io.ReadAll(utf8Reader)
	if err != nil {
		return nil, fmt.Errorf("[worker.service.Fetch] cannot read web content: %v", err)
	}
	return data, nil
}

// determineEncoding 使用 golang.org/x/net/html 库检测网页编码
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	// 调用 peek 函数可使流能重复读，检测网页编码只需前 1024 个字节即可
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}

	// TODO 返回值含义有待确认
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
