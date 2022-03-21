package base

import (
	"fmt"
)

func True(cond bool, a ...interface{}) {
	TrueF(cond, fmt.Sprint(a...))
}

func False(cond bool, a ...interface{}) {
	TrueF(!cond, fmt.Sprint(a...))
}

func TrueF(cond bool, format string, a ...interface{}) {
	if !cond {
		if a == nil || len(a) == 0 {
			panic(format)
		} else {
			panic(fmt.Sprintf(format, a...))
		}
	}
}

func FalseF(cond bool, format string, a ...interface{}) {
	TrueF(!cond, format, a...)
}

func AssertHere() {
	TrueF(false, "CANNOT run here")
}
