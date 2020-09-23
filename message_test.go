package regip_test

import (
	"regip"
	"testing"
)

func messageMatch(t *testing.T, m1, m2 *regip.Message) {
	if m1.Id != m2.Id {
		t.Errorf("IDs don't match, want '%d, got '%d", m1.Id, m2.Id)
	}
	if m1.Type != m2.Type {
		t.Errorf("Types don't match, want '%s, got %s'", m1.Type, m2.Type)
	}
	if m1.Payload != m2.Payload {
		t.Errorf("Payloads don't match, want '%s', got '%s'", m1.Payload, m2.Payload)
	}
}

func TestMessageInit(t *testing.T) {
	var _ regip.Message
}

func TestNewMessage(t *testing.T) {
	_ = regip.NewMessage(10, regip.MT_test, "")
}

func TestMarshal(t *testing.T) {
	messages := []struct {
		id      int
		typ     string
		payload string
		json    string
	}{
		{1, "testing", "", `{"id":1,"type":"testing","payload":""}`},
		{859, "other.some", "notapay", `{"id":859,"type":"other.some","payload":"notapay"}`},
	}

	for _, m := range messages {
		mes := regip.NewMessage(m.id, regip.MessageType(m.typ), m.payload)
		mar, err := mes.Marshal()
		if err != nil {
			t.Errorf("Couldn't marshal message with %#v", m)
		}
		if string(mar) != m.json {
			t.Errorf("Marshal message failed, want %s, got %s", m.json, string(mar))
		}
	}
}

func TestString(t *testing.T) {
	m := regip.NewMessage(10, regip.MT_test, "")
	s := m.String()
	if s == "" {
		t.Errorf("Message.String() returned empty string")
	}
}

func TestUnmarshalNormal(t *testing.T) {
	m := regip.NewMessage(10, regip.MT_test, "")
	raw, err := m.Marshal()
	if err != nil {
		t.Fatalf("Couldn't marshal message %s", m.String())
	}
	unm, err := regip.UnmarshalMessage(raw)
	if err != nil {
		t.Fatalf("Couldn't unmarshal message %s", raw)
	}
	messageMatch(t, m, unm)
}

func TestUnmarshalBroken(t *testing.T) {
	// Flawed JSON
	raw := []byte(`{"id":1, "type":"testing,"payload":""}`)
	_, err := regip.UnmarshalMessage(raw)
	if err == nil {
		t.Error("Unmarshal succeeded", err)
	}
}
