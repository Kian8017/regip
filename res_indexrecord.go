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
	"fmt"
)

type IndexRecord struct {
	NameID ID
}

func NewIndexRecord(ni ID) IndexRecord {
	var ir IndexRecord
	// LESSON: If you don't init, it won't copy b\c copy will see it's length 0 and not copy anything
	ir.NameID = make([]byte, ID_LENGTH)
	// Copy
	copy(ir.NameID, ni)
	ir.NameID.SetType(RT_indexRecord)
	return ir
}

func (ir IndexRecord) String() string {
	return fmt.Sprintf("IndexRecord:{NameID:(%s)}", ir.NameID.String())
}

func (ir IndexRecord) MarshalBinary() []byte          { return []byte(ir.NameID) }
func (ir IndexRecord) MarshalString() (string, error) { return ir.NameID.String(), nil }

func (ir IndexRecord) ID() ID     { return ir.NameID }
func (ir IndexRecord) Type() byte { return RT_indexRecord }

func UnmarshalIndexRecordBinary(b []byte) (IndexRecord, error) {
	if len(b) != ID_LENGTH {
		return IndexRecord{}, ErrInvalidLength
	}
	return NewIndexRecord(ID(b)), nil
}

func UnmarshalIndexRecordText(raw string) (IndexRecord, error) {
	i, err := ParseHex(raw)
	if err != nil {
		return IndexRecord{}, err
	}
	return NewIndexRecord(i), nil
}
