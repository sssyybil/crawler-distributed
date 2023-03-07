package parser

import (
	"crawler-distributed/config"
	"crawler-distributed/model"
)

type ProfileParser struct {
	userName string
}

func (p *ProfileParser) Parse(contents []byte, url string) model.ParserResult {
	return ParseProfile(contents, url, p.userName)
}

func (p *ProfileParser) Serialize() (name string, args string) {
	return config.ParseProfile, p.userName
}

func NewProfileParser(name string) *ProfileParser {
	return &ProfileParser{userName: name}
}
