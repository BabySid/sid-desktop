package common

import (
	"fmt"
	"github.com/BabySid/proto/sodor"
	"github.com/sahilm/fuzzy"
	"strings"
)

type ThomasInfosWrapper struct {
	Instance *sodor.ThomasInfos
}

func NewThomasInfosWrapper(infos *sodor.ThomasInfos) *ThomasInfosWrapper {
	return &ThomasInfosWrapper{
		Instance: infos,
	}
}

func (s *ThomasInfosWrapper) Find(name string) *ThomasInfosWrapper {
	matches := fuzzy.FindFrom(name, s)

	rs := sodor.ThomasInfos{}
	rs.ThomasInfos = make([]*sodor.ThomasInfo, 0)
	for _, match := range matches {
		rs.ThomasInfos = append(rs.ThomasInfos, s.Instance.ThomasInfos[match.Index])
	}

	return NewThomasInfosWrapper(&rs)
}

func (s *ThomasInfosWrapper) String(i int) string {
	return fmt.Sprintf("%s %s %d %s",
		s.Instance.ThomasInfos[i].Name,
		s.Instance.ThomasInfos[i].Host, s.Instance.ThomasInfos[i].Port,
		strings.Join(s.Instance.ThomasInfos[i].Tags, " "))
}

func (s *ThomasInfosWrapper) Len() int {
	return len(s.Instance.ThomasInfos)
}

func (s *ThomasInfosWrapper) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.Instance.ThomasInfos), len(s.Instance.ThomasInfos))
	for i := range s.Instance.ThomasInfos {
		rs[i] = s.Instance.ThomasInfos[i]
	}
	return rs
}

const (
	ThomasTagSep = ";"
)
