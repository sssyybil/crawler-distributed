package controller

import (
	"context"
	"crawler-distributed/config"
	"crawler-distributed/frontend/model"
	"crawler-distributed/frontend/view"
	cmodel "crawler-distributed/model"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SearchResultHandler struct {
	view   view.SearchResultView
	client *elasticsearch.Client
}

// CreateSearchResultHandler 初始化 SearchResultHandler
func CreateSearchResultHandler(template string) SearchResultHandler {
	client, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{config.ElasticSearchAddr},
		})
	if err != nil {
		log.Fatalf("[CreateSearchResultHandler] error: %v\n", err)
	}
	return SearchResultHandler{
		view:   view.CreateSearchResultView(template),
		client: client,
	}
}

// url 格式： http://localhost:8888/search?q=男 已购房&from=10
func (h SearchResultHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 获取 url 中参数 q 后面的内容
	query := strings.TrimSpace(request.FormValue("q"))

	// 获取 url 中的分页值，也就是放在 from 后面的数值
	from, err := strconv.Atoi(request.FormValue("from"))
	if err != nil {
		from = 0
	}

	log.Printf("当前检索条件为：%s，从第 %v 页开始检索\n", query, from)

	var page model.SearchResult
	if page, err = h.getSearchResult(query, from); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	if err = h.view.Render(writer, page); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func (h SearchResultHandler) getSearchResult(query string, from int) (model.SearchResult, error) {
	var result model.SearchResult
	es := h.client
	resp, err := es.Search(
		es.Search.WithIndex(config.ElasticSearchIndexWithGrpc),
		es.Search.WithQuery(query),
		es.Search.WithFrom(from),
		es.Search.WithContext(context.Background()),
	)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	var r map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %v", err)
		return result, err
	}

	var items []*cmodel.Item
	for _, hit := range r["hits"].(map[string]any)["hits"].([]any) {
		before := hit.(map[string]any)["_source"].(map[string]any)
		items = append(items, convert(before))
	}

	result.Query = query
	result.Hits = r["hits"].(map[string]any)["total"].(map[string]any)["value"].(float64)
	result.Start = from
	result.Items = items
	result.PrevFrom = result.Start - len(result.Items)
	result.NextFrom = result.Start + len(result.Items)

	return result, nil
}

func convert(before map[string]any) *cmodel.Item {
	return &cmodel.Item{
		Url:     before["url"].(string),
		Id:      before["id"].(string),
		Type:    before["type"].(string),
		Payload: before["profile"],
	}
}
