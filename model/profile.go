package model

// Item 爬取的结果信息，将存储到 ElasticSearch 中
type Item struct {
	Url     string // 被爬取的 URL 地址
	Id      string // 用户 ID
	Type    string // ElasticSearch 中的文档的存储类型
	Payload any    // 被爬取的有效信息
}

type Profile struct {
	Name       string
	Gender     string
	Age        int
	Height     int
	Weight     int
	Income     string
	Marriage   string
	Education  string
	Occupation string
	HuKou      string
	XinZuo     string
	House      string
	Car        string
	CommonInfo string
}
