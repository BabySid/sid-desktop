package base

const (
	SuccessCode = 1
)

var (
	charNeedQuota = map[rune]int{
		' ':  1,
		'&':  1,
		'(':  1,
		')':  1,
		'[':  1,
		']':  1,
		'{':  1,
		'}':  1,
		'^':  1,
		'=':  1,
		';':  1,
		'!':  1,
		'\'': 1,
		'+':  1,
		',':  1,
		'`':  1,
		'~':  1}
)

func ExecExplorer(params []string) error {
	c := NewCmd()
	c.SetSuccessCode(SuccessCode)
	return c.Run("explorer", params)
}

func ExecApp(appFullPath string) error {
	path := ""
	for _, ch := range appFullPath {
		if _, ok := charNeedQuota[ch]; ok {
			path += string('^') + string(ch)
		} else {
			path += string(ch)
		}
	}

	c := NewCmd()
	return c.Run("cmd", []string{"/c", path})
}
