package regip

import (
	"errors"
)

var (
	ResourceNotAGroup = errors.New("resource is not a group")
	UnknownResource   = errors.New("unknown resource")
)

const (
	RT_record      byte = 1
	RT_country          = 2
	RT_indexRecord      = 3
	RT_user             = 4
	RT_trigram          = 5
	RT_id               = 6
)

func StringToResourceType(inp string) (byte, bool) {
	switch inp {
	case "record":
		return RT_record, true
	case "country":
		return RT_country, true
	case "user":
		return RT_user, true
	case "indexrecord":
		return RT_indexRecord, true
	case "trigram":
		return RT_trigram, true
	default:
		return byte(0), false
	}
}

func ResourceTypeToString(inp byte) (string, bool) {
	switch inp {
	case RT_record:
		return "record", true
	case RT_country:
		return "country", true
	case RT_user:
		return "user", true
	case RT_indexRecord:
		return "indexrecord", true
	case RT_trigram:
		return "trigram", true
	default:
		return "", false
	}
}

func ResourceToMessageType(inp byte) (MessageType, bool) {
	switch inp {
	case RT_record:
		return MT_record, true
	case RT_country:
		return MT_country, true
	case RT_user:
		return MT_user, true
	case RT_indexRecord:
		return MT_indexrecord, true
	case RT_trigram:
		return MT_trigram, true
	default:
		return MT_errnotimplemented, false
	}
}

type Resource interface {
	MarshalBinary() []byte
	MarshalString() (string, error)
	Type() byte
	String() string
	ID() ID
}

type Nameable interface {
	Resource
	Name() string
}

func UnmarshalResourceBinary(t byte, raw []byte) (Resource, error) {
	switch t {
	case RT_record:
		return UnmarshalRecordBinary(raw)
	case RT_country:
		return UnmarshalCountryBinary(raw)
	case RT_user:
		return UnmarshalUserBinary(raw)
	case RT_trigram:
		return UnmarshalTrigramBinary(raw)
	case RT_indexRecord:
		return UnmarshalIndexRecordBinary(raw)
	default:
		panic("can't unmarshal resource of type " + string(t))
	}
}

func UnmarshalResourceText(t byte, raw string) (Resource, error) {
	switch t {
	case RT_record:
		return UnmarshalRecordText(raw)
	case RT_country:
		return UnmarshalCountryText(raw)
	case RT_user:
		return UnmarshalUserText(raw)
	case RT_trigram:
		return UnmarshalTrigramText(raw)
	case RT_indexRecord:
		return UnmarshalIndexRecordText(raw)
	default:
		panic("can't unmarshal resource of type " + string(t))
	}
}

type MetaID struct {
	M_ID ID `json:"id"`
}

// FIXME: Implement ResourceGroup
