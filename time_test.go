package kiicli

import "testing"

func Test_timeFromStringInUTC(t *testing.T) {
	v, _ := timeFromStringInUTC("a")
	expected := int64(-62135596800)
	if v.Unix() != expected {
		t.Errorf("expected %v, but %v", expected, v.Unix())
	}
}
