package iprange

import (
	"testing"
)

func FuzzIpRange(f *testing.F) {
	testcases := []string{
		"192.168.1.1",
		"192.168.1.1/24",
		"10.1.2.3/16",
		"192.168.1.*",
		"192.168.1.10-20",
		"192.168.10-20.1",
		"0-255.1.1.1",
		"1-2.3-4.5-6.7-8",
		"192.168.10-20.*/25",
		//Тесты с ошибкой
		"0",
		"192.168.10",
		"0.0.0.0/70",
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, a string) {
		_, _ = Parse(a)
	})
}
