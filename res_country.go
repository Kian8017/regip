package regip

import (
	"encoding/json"
	"fmt"
	"strings"
)

// COUNTRY
// Country holds the name and abbreviation of the Country
// On disk, the first byte holds the length of the abbreviation, then comes the abbreviation, then the rest is the name
type Country struct {
	C_abbr string `json:"abbr"`
	C_name string `json:"name"`
	MetaID
}

func NewCountry(a, n string) Country {
	var c Country
	// FIXME: enforce lowercase everywhere
	// (switch to upper case and/or title case in frontend)
	c.C_abbr = strings.ToLower(a)
	c.C_name = strings.ToLower(n)
	return c
}

func (c Country) String() string {
	return fmt.Sprintf("Country{ID:(%s),Abbr:(%s)}|%s", c.ID().String(), c.C_abbr, c.C_name)
}

func (c Country) MarshalBinary() []byte {
	// Just to make sure
	c.C_abbr = strings.ToLower(c.C_abbr)
	c.C_name = strings.ToLower(c.C_name)

	bufV := make([]byte, len(c.C_abbr)+len(c.C_name)+1) // 1 for abbreviation length
	bufV[0] = byte(uint8(len(c.C_abbr)))
	copy(bufV[1:len(c.C_abbr)+1], []byte(c.C_abbr))
	copy(bufV[1+len(c.C_abbr):], []byte(c.C_name))
	return bufV
}

func (c Country) MarshalString() (string, error) {
	c.M_ID = c.ID()
	raw, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(raw), err
}

func (c Country) ID() ID     { return NewID(byte(RT_country), c.MarshalBinary()) }
func (c Country) Type() byte { return RT_country }

// Helper functions
func UnmarshalCountryBinary(b []byte) (Country, error) {
	var c Country
	abbrLen := uint8(b[0])
	c.C_abbr = string(b[1 : abbrLen+1])
	c.C_name = string(b[1+abbrLen:])
	// FIXME: check for bounds errors
	return c, nil
}

func UnmarshalCountryText(raw string) (Country, error) {
	var c Country
	err := json.Unmarshal([]byte(raw), &c)
	if err != nil {
		return Country{}, err
	}
	return c, nil
}
