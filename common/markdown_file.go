package common

import (
	"github.com/BabySid/gobase"
	"github.com/sahilm/fuzzy"
)

type MarkDownFileList []MarkDownFile

func (s *MarkDownFileList) Find(name string) *MarkDownFileList {
	matches := fuzzy.FindFrom(name, s)

	var rs MarkDownFileList
	for _, match := range matches {
		rs = append(rs, (*s)[match.Index])
	}

	return &rs
}

func (s *MarkDownFileList) String(i int) string {
	return (*s)[i].Name + " " + (*s)[i].Cont
}

func (s *MarkDownFileList) Len() int {
	return len(*s)
}

func (s *MarkDownFileList) Set(d []MarkDownFile) {
	*s = d
}

func (s *MarkDownFileList) Append(d MarkDownFile) {
	*s = append(*s, d)
}

func (s *MarkDownFileList) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(*s), len(*s))
	for i, f := range *s {
		rs[i] = f
	}
	return rs
}

type MarkDownFile struct {
	gobase.TableModel
	Name string `gorm:"not null;size:64" json:"name"`
	Cont string `gorm:"not null;type:text" json:"content"`
}

func (t MarkDownFile) TableName() string {
	return "markdown_files"
}

func (t MarkDownFile) UpdateFields() []string {
	return []string{
		"Name",
		"Cont",
	}
}
