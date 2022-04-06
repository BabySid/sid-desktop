package common

import "github.com/sahilm/fuzzy"

func FilterOption(filter string, opt []string) (resultList []string) {
	if filter == "" {
		return opt
	}

	rs := make([]string, 0)
	match := fuzzy.Find(filter, opt)
	for _, m := range match {
		rs = append(rs, m.Str)
	}
	return rs
}
