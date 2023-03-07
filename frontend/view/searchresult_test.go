package view

import (
	"crawler-distributed/frontend/model"
	crawlerModel "crawler-distributed/model"
	"os"
	"testing"
)

func TestSearchResultView_Render(t *testing.T) {
	item := crawlerModel.Item{
		Url:  "http://localhost:8080/mock/album.zhenai.com/u/3903982005871861481",
		Id:   "3903982005871861481",
		Type: "zhenai",
		Payload: crawlerModel.Profile{
			Name:       "一身傲气如你*",
			Age:        24,
			Height:     140,
			Weight:     257,
			Income:     "8001-10000元",
			Gender:     "男",
			XinZuo:     "天蝎座",
			Occupation: "测试工程师",
			Marriage:   "未婚",
			House:      "无房",
			HuKou:      "沈阳市",
			Education:  "硕士",
			Car:        "有车",
		},
	}
	var items []*crawlerModel.Item
	for i := 0; i < 10; i++ {
		items = append(items, &item)
	}

	outFile, err := os.Create("template.test.html")
	data := model.SearchResult{
		Hits:  100,
		Items: items,
	}

	view := CreateSearchResultView("template.html")
	err = view.Render(outFile, data)
	if err != nil {
		panic(err)
	}
}
