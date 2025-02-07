package jsonutil

import (
	"encoding/json"
	"math"
	"strings"
)

type Option struct {
	Indent uint
	UseTab bool
}

const tabEquivalentSpaceCount float64 = 4

func Json(target any, o Option) ([]byte, error) {
	indent := o.Indent
	if indent == 0 {
		return json.Marshal(o)
	}
	useTab := o.UseTab
	spacer := " "
	if useTab {
		indent = uint(math.Ceil(float64(indent) / tabEquivalentSpaceCount))
		spacer = "	"
	}
	return json.MarshalIndent(target, "", strings.Repeat(spacer, int(indent)))
}
