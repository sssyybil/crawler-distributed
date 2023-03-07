package model

import (
	"crawler-distributed/config"
)

// Request 请求信息
type Request struct {
	Url    string
	Parser Parser
}

// ParserResult 解析结果
type ParserResult struct {
	Requests []Request
	Items    []Item
}

// Processor 处理器
type Processor func(Request) (ParserResult, error)

// Parser 解析器
type Parser interface {
	Parse(contents []byte, url string) ParserResult
	Serialize() (name string, args string)
}

// NilParser Parser 解析器实现一
type NilParser struct{}

func (NilParser) Parse(contents []byte, url string) ParserResult {
	return ParserResult{}
}

func (NilParser) Serialize() (name string, args string) {
	return config.NilParser, ""
}

// ParserFunc 根据传入的内容和 URL 即可具有解析功能。
type ParserFunc func(contents []byte, url string) ParserResult

// FuncParser Parser 解析器实现二
type FuncParser struct {
	parser ParserFunc
	name   string
}

func NewFuncParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		name:   name,
	}
}

func (f *FuncParser) Parse(contents []byte, url string) ParserResult {
	return f.parser(contents, url)
}

func (f *FuncParser) Serialize() (name string, args string) {
	return f.name, ""
}
