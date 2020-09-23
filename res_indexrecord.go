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
