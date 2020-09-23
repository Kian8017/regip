package regip_test

import (
	"regip"
	"testing"
)

var sample = []regip.ID{
	regip.NewID(regip.RT_record, []byte("test")),
	regip.NewID(regip.RT_record, []byte("other")),
	regip.NewID(regip.RT_record, []byte("what")),
}
var sampleLen = 16*3 + 4

func TestLengthEmpty(t *testing.T) {
	em := regip.NewTrigram("abc", []regip.ID{})
	raw := em.MarshalBinary()
	if len(raw) != 4 {
		t.Errorf("Marshal length incorrect, want %d, got %d", 4, len(raw))
	}
	if uint8(raw[0]) != 3 {
		t.Errorf("Incorrect Combi length, want %d, got %d", 3, raw[0])
	}
}

func TestLengthCouple(t *testing.T) {
	em := regip.NewTrigram("abc", sample)
	raw := em.MarshalBinary()
	if len(raw) != sampleLen {
		t.Errorf("Marshal length incorrect, want %d, got %d", sampleLen, len(raw))
	}
	if uint8(raw[0]) != 3 {
		t.Errorf("Incorrect Combi length, want %d, got %d", 3, raw[0])
	}
}

func TestUnmarshalBinary(t *testing.T) {
	em := regip.NewTrigram("abc", sample)
	raw := em.MarshalBinary()
	rec, err := regip.UnmarshalTrigramBinary(raw)
	if err != nil {
		t.Fatal("UnmarshalBinary: got ", err)
	}
	if len(rec.Ids) != len(sample) {
		t.Fatalf("UnmarshalBinary: length mismatch, want %d, got %d", len(sample), len(rec.Ids))
	}
	for i, r := range rec.Ids {
		if !sample[i].Equal(r) {
			t.Errorf("Unmarshal ID mismatch: expected %s, got %s", sample[i].String(), r.String())
		}
	}
}
