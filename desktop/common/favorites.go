package common

import (
	"fmt"
	"github.com/sahilm/fuzzy"
	"strings"
)

type FavoritesList struct {
	favors []Favorites
}

func NewFavoritesList() *FavoritesList {
	return &FavoritesList{
		favors: make([]Favorites, 0),
	}
}

func (s *FavoritesList) Find(name string) *FavoritesList {
	matches := fuzzy.FindFrom(name, s)

	rs := NewFavoritesList()
	for _, match := range matches {
		rs.favors = append(rs.favors, s.favors[match.Index])
	}

	return rs
}

func (s *FavoritesList) UpdateFavorites(d Favorites) {
	for _, app := range s.favors {
		if app.ID == d.ID {
			app = d
			return
		}
	}
}

func (s *FavoritesList) String(i int) string {
	return s.favors[i].Name + s.favors[i].Url + strings.Join(s.favors[i].Tags, " ")
}

func (s *FavoritesList) Len() int {
	return len(s.favors)
}

func (s *FavoritesList) Set(d []Favorites) {
	if d == nil {
		return
	}
	s.favors = d
}

func (s *FavoritesList) Append(d Favorites) {
	s.favors = append(s.favors, d)
}

func (s *FavoritesList) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.favors), len(s.favors))
	for i := range s.favors {
		rs[i] = s.favors[i]
	}
	return rs
}

func (s *FavoritesList) GetFavorites() []Favorites {
	return s.favors
}

func (s *FavoritesList) Debug() {
	for _, fav := range s.favors {
		fmt.Println(fav.ID, fav.Name, fav.Url, fav.Tags, fav.CreateTime, fav.AccessTime)
	}
}

const (
	FavorTagSep = ";"
)

type Favorites struct {
	ID         int64    `json:"-"`
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	Tags       []string `json:"tags"`
	CreateTime int64    `json:"create_time"`
	AccessTime int64    `json:"access_time"`
}
