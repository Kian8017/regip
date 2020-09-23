package regip

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Trigram struct {
	Combi string `json:"combi"`
	Ids   []ID   `json:"ids"`
	MetaID
}

func NewTrigram(c string, Ids []ID) *Trigram {
	var t Trigram
	t.Combi = c
	t.Ids = Ids
	t.Sort()
	return &t
}

func (t *Trigram) String() string {
	return fmt.Sprintf("Trigram:{Combi:%s(Len:%d),NumIds:%d}", t.Combi, len(t.Combi), len(t.Ids))
}

func (t *Trigram) Sort() {
	// Trigram IDs are required to be in ascending order
	sort.Slice(t.Ids, func(i, j int) bool {
		return t.Ids[i].LessThan(t.Ids[j])
	})
}

func (t *Trigram) Have(i ID) bool {
	ind := sort.Search(len(t.Ids), func(j int) bool {
		return t.Ids[j].GreaterThanEquals(i)
	})
	if ind < len(t.Ids) && t.Ids[ind].Equal(i) {
		return true
	} else {
		return false
	}
}

func (t *Trigram) AddIds(n []ID) {
	var toAdd []ID
	for _, cur := range n {
		if !t.Have(cur) {
			toAdd = append(toAdd, cur)
		}
	}
	if len(toAdd) > 0 {
		t.Ids = append(t.Ids, toAdd...)
		t.Sort()
	}
}

func (t *Trigram) MarshalBinary() []byte {
	panic("MARSHALLING TRIGRAM " + t.String())
	buf := make([]byte, 1+len(t.Combi)+(len(t.Ids)*ID_LENGTH)) // combi may not necessarily be 3 bytes (Unicode)
	// first byte is length of combi
	buf[0] = uint8(len(t.Combi))
	// then combi
	copy(buf[1:], t.Combi)
	start := len(t.Combi) + 1
	// then the Ids themselves
	for i := 0; i < len(t.Ids); i++ {
		copy(buf[start+(ID_LENGTH*i):], t.Ids[i])
	}

	return buf
}

func (t *Trigram) MarshalString() (string, error) {
	t.M_ID = t.ID()
	raw, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// The only thing that matters for ID's sake is the combi
func (t *Trigram) ID() ID     { return NewID(RT_trigram, []byte(t.Combi)) }
func (t *Trigram) Type() byte { return RT_trigram }

func UnmarshalTrigramBinary(b []byte) (*Trigram, error) {
	var t Trigram
	lc := uint8(b[0])
	meta := int(lc) + 1 // Extra for combi length
	// Total length minus meta / ID_Length

	if (len(b)-meta)%ID_LENGTH != 0 {
		// FIXME: DEBUG
		panic(fmt.Sprintf("LC IS %d, tested is %d, total is %s, tlen is %d", lc, len(b)-meta, string(b), len(b)))
		return nil, ErrInvalidLength
	}
	t.Combi = string(b[1 : 1+lc])

	// DEBUG
	if len(t.Combi) != 3 {
		panic(fmt.Sprintf("COMBI LENGTH IS %d, not 3, combi is '%s'", len(t.Combi), t.Combi))
	}

	t.Ids = SplitIds(b[meta:])
	t.Sort()
	return &t, nil
}

func UnmarshalTrigramText(raw string) (*Trigram, error) {
	var t Trigram
	err := json.Unmarshal([]byte(raw), &t)
	if err != nil {
		return nil, err
	}
	t.Sort()
	return &t, nil
}
