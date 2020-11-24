package tool

import (
	"bytes"
	"github.com/robertkrimen/otto"
	"regexp"
	"strconv"
	"strings"
)

func GetCFEmailPayload(str string) string {
	s := strings.Split(str, "data-cfemail=")
	if len(s) > 1 {
		s = strings.Split(s[1], "\"")
		str = s[1]
		return str
	}
	return ""
}

// Remove cloud flare email protection
func CFEmailDecode(a string) (s string) {
	if a == "" {
		return
	}
	var e bytes.Buffer
	r, _ := strconv.ParseInt(a[0:2], 16, 0)
	for n := 4; n < len(a)+2; n += 2 {
		i, _ := strconv.ParseInt(a[n-2:n], 16, 0)
		e.WriteString(string(i ^ r))
	}
	return e.String()
}

// Get result var of a js script
func ScriptGet(js string, varname string) string {
	vm := otto.New()
	_, err := vm.Run(js)
	if err != nil {
		return ""
	}
	if value, err := vm.Get(varname); err == nil {
		if v, err := value.ToString(); err == nil {
			return v
		}
	}
	return ""
}

// replace location with varname and remove window
func ScriptReplace(js string, varname string) string {
	strs := strings.Split(js, ";")
	varWindow := ""
	varLocation := ""
	bound := len(strs)

	if len(js) < 2 {
		return js
	}
	for i, _ := range strs {
		//replace location
		if strings.Contains(strs[i], "location") {
			strarr := strings.Split(strs[i], " = ")
			if len(strarr) == 2 {
				varLocation = strarr[0]
				strs[i] = ""
			} else {
				re3, err := regexp.Compile("location.*?[]]") // location[_jzvXT]
				if err == nil {
					strs[i] = re3.ReplaceAllLiteralString(strs[i], varname)
				}
			}
		}
		if varLocation != "" && strings.Contains(strs[i], varLocation) {
			re3, err := regexp.Compile(varLocation + ".*?[]]") // _LoKlO[_jzvXT]
			if err == nil {
				strs[i] = re3.ReplaceAllLiteralString(strs[i], varname)
			}
		}
		// remove window
		if strings.Contains(strs[i], "window") {
			varWindow = strings.Split(strs[i], " = window")[0]
			strs[i] = ""
		}
	}

	if varWindow != "" {
		for i, _ := range strs {
			if strings.Contains(strs[i], varWindow) {
				bound = i
				break
			}
		}
	}
	js = strings.Join(strs[:bound], ";")
	return js
}