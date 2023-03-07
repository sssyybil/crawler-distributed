package view

import (
	"crawler-distributed/frontend/model"
	"html/template"
	"io"
)

type SearchResultView struct {
	template *template.Template
}

// CreateSearchResultView 根据文件名（frontend/view/template.html）创建 SearchResultView
func CreateSearchResultView(filename string) SearchResultView {
	return SearchResultView{
		template: template.Must(template.ParseFiles(filename)),
	}
}

// Render 将 Writer 中读取到的数据写入到 SearchResult 中
func (s SearchResultView) Render(w io.Writer, data model.SearchResult) error {
	return s.template.Execute(w, data)
}
