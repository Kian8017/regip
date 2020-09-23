package regip

import (
	"strings"
)

func FilterName(want string, inputF *Flow) *Flow {
	return NewFlowFilter(func(n Resource) bool {
		res, ok := n.(Nameable)
		if !ok {
			return false
		}
		return strings.Contains(res.Name(), want)
	})(inputF)
}

/*
func FilterType(nt RecordType, inputF *Flow) *Flow {
	return NewFlowFilter(func(n *Record) bool {
		return n.Type == nt
	})(inputF)
}

func FilterCountry(ci ID, inputF *Flow) *Flow {
	return NewFlowFilter(func(n *Record) bool {
		return n.Country.Equal(ci)
	})(inputF)
}
*/
