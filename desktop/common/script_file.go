package common

import (
	"fmt"
	"github.com/sahilm/fuzzy"
)

type ScriptFileList struct {
	scripts []ScriptFile
}

func NewScriptFileList() *ScriptFileList {
	return &ScriptFileList{
		scripts: make([]ScriptFile, 0),
	}
}

func (s *ScriptFileList) Find(name string) *ScriptFileList {
	matches := fuzzy.FindFrom(name, s)

	rs := NewScriptFileList()
	for _, match := range matches {
		rs.scripts = append(rs.scripts, s.scripts[match.Index])
	}

	return rs
}

func (s *ScriptFileList) UpdateFavorites(f ScriptFile) {
	for _, file := range s.scripts {
		if file.ID == f.ID {
			file = f
			return
		}
	}
}

func (s *ScriptFileList) String(i int) string {
	return s.scripts[i].Name
}

func (s *ScriptFileList) Len() int {
	return len(s.scripts)
}

func (s *ScriptFileList) Set(d []ScriptFile) {
	if d == nil {
		return
	}
	s.scripts = d
}

func (s *ScriptFileList) Append(d ScriptFile) {
	s.scripts = append(s.scripts, d)
}

func (s *ScriptFileList) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.scripts), len(s.scripts))
	for i := range s.scripts {
		rs[i] = s.scripts[i]
	}
	return rs
}

func (s *ScriptFileList) GetScriptFiles() []ScriptFile {
	return s.scripts
}

func (s *ScriptFileList) Debug() {
	for _, fav := range s.scripts {
		fmt.Println(fav.ID, fav.Name, fav.CreateTime, fav.AccessTime)
	}
}

type ScriptFile struct {
	ID         int64  `json:"-"`
	Name       string `json:"name"`
	Cont       string `json:"content"`
	CreateTime int64  `json:"create_time"`
	AccessTime int64  `json:"access_time"`
}
