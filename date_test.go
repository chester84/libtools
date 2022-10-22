package libtools

import (
	"testing"
)

func TestHumanUnixMillis(t *testing.T) {
	var tm int64 = 3670000
	display := HumanUnixMillis(tm)
	if display == "" {
		t.Errorf("get error, display: %s", display)
	} else {
		t.Logf("[HumanUnixMillis] get ret: %s", display)
	}
}
