package parser

import (
	"crawler-distributed/model"
	"log"
	"regexp"
	"strconv"
)

var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>(\d+)岁</td>`)
var heightRe = regexp.MustCompile(`<td><span class="label">身高：</span>(\d+)CM</td>`)
var incomeRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
var weightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">(\d+)KG</span></td>`)
var genderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var xinZuoRe = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)
var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var educationRe = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var occupationRe = regexp.MustCompile(`<td><span class="label">职业：</span><span field="">([^<]+)</span></td>`)
var huKouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var houseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var carRe = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)

// 匹配用户详情下的【猜你喜欢】部分的用户
var guessRe = regexp.MustCompile(` <a class="exp-user-name"[^>]*href="(http://localhost:8080/mock/album.zhenai.com/u/\d+)">([^>]+)</a>`)

// 匹配 URL 中的 用户ID
var idUrlRe = regexp.MustCompile(`http://localhost:8080/mock/album.zhenai.com/u/(\d+)`)

func ParseProfile(contents []byte, url string, name string) model.ParserResult {
	profile := model.Profile{Name: name}

	if age, err := strconv.Atoi(extractString(contents, ageRe)); err == nil {
		profile.Age = age
	}
	if height, err := strconv.Atoi(extractString(contents, heightRe)); err == nil {
		profile.Height = height
	}
	if weight, err := strconv.Atoi(extractString(contents, weightRe)); err == nil {
		profile.Weight = weight
	}
	profile.Income = extractString(contents, incomeRe)
	profile.Gender = extractString(contents, genderRe)
	profile.Car = extractString(contents, carRe)
	profile.Education = extractString(contents, educationRe)
	profile.HuKou = extractString(contents, huKouRe)
	profile.House = extractString(contents, houseRe)
	profile.Marriage = extractString(contents, marriageRe)
	profile.Occupation = extractString(contents, occupationRe)
	profile.XinZuo = extractString(contents, xinZuoRe)

	result := model.ParserResult{
		Items: []model.Item{
			{
				Url:     url,
				Id:      extractString([]byte(url), idUrlRe),
				Type:    "zhenai",
				Payload: profile,
			},
		},
	}

	matches := guessRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, model.Request{
			Url:    string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	log.Printf("[ParseProfile] result:%v\n", result)
	return result
}

// extractString 提取正则表达式匹配后的结果中的第一个内容
func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
