package kiicli

import "testing"

func Test_converLogFormat(t *testing.T) {
	k := convertLogFormat("${time} [${level}]")
	expected := "{{.time}} [{{.level}}]"
	if k != expected {
		t.Errorf("expected %v, but %v", expected, k)
	}

	k = convertLogFormat("{{.time}} [{{.level}}]")
	expected = "{{.time}} [{{.level}}]"
	if k != expected {
		t.Errorf("expected %v, but %v", expected, k)
	}

	k = convertLogFormat("${level}")
	expected = "{{.level}}"
	if k != expected {
		t.Errorf("expected %v, but %v", expected, k)
	}

	k = convertLogFormat("${entry-name}")
	expected = `{{index . "entry-name"}}`
	if k != expected {
		t.Errorf("expected %v, but %v", expected, k)
	}
}
