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
