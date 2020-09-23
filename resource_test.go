package regip_test

import (
	"regip"
	"testing"
)

func TestInterfaces(t *testing.T) {
	var _ regip.Resource = &regip.Record{}
	var _ regip.Resource = &regip.User{}
	var _ regip.Resource = regip.Country{}
	var _ regip.Resource = regip.IndexRecord{}
	var _ regip.Resource = &regip.Trigram{}
}
