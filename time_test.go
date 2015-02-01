package kiicli

import "testing"

func Test_timeStringInUTCToLocalTime(t *testing.T) {
	v, _ := timeStringInUTCToLocalTime("a")
	expected := int64(-62135596800)
	if v.Unix() != expected {
		t.Errorf("expected %v, but %v", expected, v.Unix())
	}
}
