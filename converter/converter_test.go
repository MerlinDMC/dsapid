package converter

import (
	"testing"
)

func TestIntToInt64(t *testing.T) {
	const (
		in  int   = 1234678
		out int64 = 1234678
	)

	if v := DecodeToInt64(in); v != out {
		t.Errorf("int to int64 conversion failed: expected %d but got %d", out, v)
	}
}

func TestInt64ToInt64(t *testing.T) {
	const (
		in  int64 = 1234678
		out int64 = 1234678
	)

	if v := DecodeToInt64(in); v != out {
		t.Errorf("int64 to int64 conversion failed: expected %d but got %d", out, v)
	}
}

func TestFloat64ToInt64(t *testing.T) {
	const (
		in  float64 = 1234678.0
		out int64   = 1234678
	)

	if v := DecodeToInt64(in); v != out {
		t.Errorf("float64 to int64 conversion failed: expected %d but got %d", out, v)
	}
}

func TestStringToInt64(t *testing.T) {
	const (
		in  string = "1234678"
		out int64  = 1234678
	)

	if v := DecodeToInt64(in); v != out {
		t.Errorf("string to int64 conversion failed: expected %d but got %d", out, v)
	}
}
