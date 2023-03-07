package parser

import (
	"crawler-distributed/config"
	"crawler-distributed/model"
	"regexp"
)

var (
	profileRe = regexp.MustCompile(`<a href="(http://localhost:8080/mock/album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	cityUrlRe = regexp.MustCompile(`href="(http://localhost:8080/mock/www.zhenai.com/zhenghun/shanghai/[0-9]+)`)
)

// ParseCity 城市解析器
func ParseCity(contents []byte, url string) model.ParserResult {
	result := model.ParserResult{}

	all := profileRe.FindAllSubmatch(contents, -1)
	for _, m := range all {
		url := string(m[1])
		// 只向数据库中存储有价值的数据，user 的名字除外
		//result.Items = append(result.Items, "User "+name)
		result.Requests = append(result.Requests, model.Request{
			Url:    url,
			Parser: NewProfileParser(string(m[2])),
		})
	}

	matches := cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, model.Request{
			Url:    string(m[1]),
			Parser: model.NewFuncParser(ParseCity, config.ParseCity),
		})
	}

	return result
}
