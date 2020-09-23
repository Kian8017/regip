package regip

import (
	"encoding/json"
	"fmt"
)

// RECORD
// Record is the primary struct
// In the database, it's stored like this:
// [C][C][C][C] [T] [N][N][N]...
// Where C are the bytes that represent the country, T is type, and N the record
type Record struct {
	Name     string     `json:"name"`
	NameType RecordType `json:"type"`
	Place    ID         `json:"country"`
	MetaID
}

func (r *Record) String() string {
	return fmt.Sprintf("Record{ID:(%s),Place:(%s),Type:(%d)}|%s", r.ID().String(), r.Place.String(), r.NameType, r.Name)
}

func NewRecord(n string, t RecordType, c ID) *Record {
	return &Record{Name: n, NameType: t, Place: c}
}

// Marshal encodes the record into its byte array
func (r *Record) MarshalBinary() []byte {
	buf := make([]byte, len(r.Name)+ID_LENGTH+1) // The record + id + 1 for metadata
	copy(buf, r.Place)
	buf[ID_LENGTH] = byte(r.NameType) // ID_LENGTH = index after country id
	// Add record
	copy(buf[ID_LENGTH+1:], []byte(r.Name))
	return buf
}

func (r *Record) MarshalString() (string, error) {
	r.M_ID = r.ID()
	raw, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (r *Record) ID() ID     { return NewID(byte(RT_record), r.MarshalBinary()) }
func (r *Record) Type() byte { return RT_record }

func UnmarshalRecordBinary(b []byte) (*Record, error) {
	var n Record
	n.Place = b[:ID_LENGTH]
	n.NameType = RecordType(b[ID_LENGTH])
	n.Name = string(b[ID_LENGTH+1:])
	return &n, nil // FIXME add error checking
}

func UnmarshalRecordText(raw string) (*Record, error) {
	var n Record
	err := json.Unmarshal([]byte(raw), &n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

type RecordType uint8

const (
	Name    RecordType = 1
	Place              = 2
	Various            = 4
)

func ParseRecordType(i string) (RecordType, bool) {
	switch i {
	case "name":
		return Name, true
	case "place":
		return Place, true
	case "various":
		return Various, true
	default:
		return RecordType(0), false
	}
}
