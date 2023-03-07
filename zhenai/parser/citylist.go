package parser

import (
	"crawler-distributed/config"
	"crawler-distributed/model"
	"log"
	"regexp"
)

const cityListRe = `<a href="(http://localhost:8080/mock/www.zhenai.com/zhenghun/[a-zA-Z0-9]+)"[^>]*>([^>]*)</a>`

// ParseCityList 城市列表解析器
func ParseCityList(contents []byte, _ string) model.ParserResult {

	log.Printf("ParseCityList is running...\n")

	re := regexp.MustCompile(cityListRe)
	all := re.FindAllSubmatch(contents, -1)

	result := model.ParserResult{}
	for _, m := range all {
		// 只向数据库中存储有价值的数据，citylist 的名字除外
		//result.Items = append(result.Items, "City "+string(m[2]))
		result.Requests = append(result.Requests, model.Request{
			Url:    string(m[1]),
			Parser: model.NewFuncParser(ParseCity, config.ParseCity),
		})
	}

	return result
}
