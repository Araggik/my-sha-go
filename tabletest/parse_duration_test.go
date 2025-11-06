package tabletest

import (
	"testing"
	"time"
	"strconv"
)

func TestParseDurationValid(t *testing.T) {
	var sec int64 = 1000000000
	min := sec * 60
	hour := 60 * min

	tests := []struct {
		input string
		want  time.Duration
	}{
		{"0", 0},
		{"-1.5h", time.Duration(-hour - 30*min)},
	}

	for _, v := range tests {
		result, err := ParseDuration(v.input)

		if err != nil {
			t.Errorf("ParseDuration has error for input: %v", v.input)
		} else if result != v.want {
			t.Errorf("ParseDuration result not equal expected: %v != %v", result, v.want)
		}
	}
}

func TestParseDurationError(t *testing.T) {
	maxIntStr := strconv.FormatInt(1<<63-1, 10)

	tests := []string{
		"",
		".",
		"a",
		maxIntStr + "1",
		//Переполнение int64 на единицу
		maxIntStr[:len(maxIntStr)-1] + "8",
	}

	for _, v := range tests {
		_, err := ParseDuration(v)

		if err == nil {
			t.Errorf("ParseDuration has no error for input: %v", v)
		}
	}
}