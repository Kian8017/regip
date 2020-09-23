/* Copyright (c) 2020 Kian Musser.
 * This file is part of regip.
 *
 * regip is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * regip is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with regip.  If not, see <https://www.gnu.org/licenses/>.
 */

package regip

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

const ID_LENGTH = 16

var ErrInvalidLength error = errors.New("invalid length for ID")

// First byte is the ID type (allows prefix searches), the next 15 are a truncated sha256 hash
type ID []byte // Length: 16

func (i ID) String() string { return hex.EncodeToString(i) }

func ParseHex(inp string) (ID, error) {
	if len(inp) != ID_LENGTH*2 {
		return nil, ErrInvalidLength
	}
	dec, err := hex.DecodeString(inp)
	if err != nil {
		return nil, err
	}
	return ID(dec), nil
}

func NewID(f byte, val []byte) ID {
	ret := make([]byte, ID_LENGTH)
	ret[0] = f
	// Generate hash
	valHash := sha256.Sum256(val)
	copy(ret[1:], valHash[:ID_LENGTH-1])
	return ret
}

func (i ID) SetType(f byte) { i[0] = f }

func (i ID) Equal(other ID) bool             { return bytes.Equal(i, other) }
func (i ID) LessThan(other ID) bool          { return bytes.Compare(i, other) < 0 }
func (i ID) GreaterThanEquals(other ID) bool { return bytes.Compare(i, other) >= 0 }

func (i ID) MarshalString() (string, error) { return i.String(), nil }
func (i ID) MarshalText() ([]byte, error)   { return []byte(i.String()), nil }

func (i *ID) UnmarshalText(raw []byte) error {
	j, err := ParseHex(string(raw))
	//fmt.Println("UNMARSHAL TEXT ON ID CALLED, raw is: ", string(raw))
	//fmt.Println("\tERROR IS ", err, " RESULT IS", j)
	*i = j
	return err
}

// To match resource interface
func (i ID) MarshalBinary() []byte { return i }
func (i ID) Type() byte            { return RT_id }
func (i ID) ID() ID                { return i }

func UnmarshalIDBinary(b []byte) (ID, error) {
	if len(b) != ID_LENGTH {
		return ID{}, ErrInvalidLength
	}
	return ID(b), nil
}
func UnmarshalIDText(raw string) (ID, error) { return ParseHex(raw) }

func SplitIds(val []byte) []ID {
	if len(val)%ID_LENGTH != 0 {
		panic(fmt.Sprintf("id.Split: incorrect length %d -- can't split", len(val)))
	}
	num := len(val) / ID_LENGTH
	ret := make([]ID, 0, num)
	for i := 0; i < len(val); i += ID_LENGTH {
		ret = append(ret, val[i:i+ID_LENGTH])
	}
	return ret
}
