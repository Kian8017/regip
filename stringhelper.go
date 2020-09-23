package regip

import (
	"strings"
)

const DROP string = "\r\n\ufeff" // Drop Byte Order Mark
const REPLACE string = "\t"

func NormalizeString(s string) string {
	ns := strings.ToLower(s)
	var sb strings.Builder
	for _, r := range ns {
		if IsInList(r, REPLACE) {
			sb.WriteString(" ")
		} else if IsInList(r, DROP) {
			continue
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func GenerateTrigrams(s string, pad bool) []string {
	var ns string
	if pad {
		ns = NormalizeString("  " + s + "  ")
	} else {
		ns = NormalizeString(s)
	}
	var ret []string
	for i := 0; i < len(ns)-2; i++ {
		ret = append(ret, ns[i:i+3])
	}
	return ret
}

func IsInList(c rune, list string) bool {
	for _, r := range list {
		if r == c {
			return true
		}
	}
	return false
}
