package regip

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Message struct {
	Id      int         `json:"id"`
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
}

func NewMessage(i int, t MessageType, p string) *Message {
	return &Message{Id: i, Type: t, Payload: p}
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) String() string {
	return fmt.Sprintf("Message(%s-%d){%s}", m.Type, m.Id, m.Payload)
}

func UnmarshalMessage(raw []byte) (*Message, error) {
	var mess Message
	err := json.Unmarshal(raw, &mess)
	return &mess, err
}

var ErrUnknownType = errors.New("unknown message type")

type MessageType string

const (
	// Meta
	MT_stop     MessageType = "stop"
	MT_ok                   = "ok"
	MT_fail                 = "fail"
	MT_exists               = "exists"
	MT_noexists             = "exists"
	MT_test                 = "testing"

	// Errors
	MT_errnotimplemented    MessageType = "notimplemented"
	MT_errnotauth                       = "notauthorized"
	MT_errinvalidformatting             = "invalidformatting"

	// Resources
	MT_record      MessageType = "record"
	MT_country                 = "country"
	MT_user                    = "user"
	MT_indexrecord             = "indexrecord"
	MT_trigram                 = "trigram"

	// API: Commands
	MT_ping      MessageType = "ping"
	MT_pong                  = "pong"
	MT_list                  = "list"
	MT_new                   = "new"
	MT_get                   = "get"
	MT_del                   = "del"
	MT_fullindex             = "fullindex"

	// API: Authorization
	MT_login MessageType = "login"
)
