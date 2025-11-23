package iprange

import (
	"net"
	"strconv"
	"strings"
	"testing"
)

func FuzzIpRange(t *testing.F) {

}

func TestIpRangePossibilities(t *testing.T) {
	res, err := Parse("192.168.255.0/1")
	t.Logf("res: %v, err: %v \n", res, err)
}

// TODO: Функция аналог Parse
func TestingParse(address string) (AddressRange, error) {
	var result AddressRange

	numbers := strings.Split(address, ".")

	if strings.Contains(numbers[3], "/") {
		mask := strings.Split(numbers[3], "/")[2]

		if mask == "24" {
			result.Min = net.IP(toBytes(numbers[:3]))
		} else if mask == "16" {

		} else if mask == "8" {

		} else {

		}

	} else if strings.Contains(numbers[3], ".") {

	} else if strings.Contains(numbers[3], "-") {

	} else {

	}

	return AddressRange{}, nil
}

func toBytes(slice []string) []byte {
	res := make([]byte, 4)

	for i, s := range slice {
		num, _ := strconv.Atoi(s)

		res[i] = byte(num)
	}

	return res
}

// TODO: функция добавляющая байты к слайсу байтов
func addByte(slice []byte, add []byte) []byte {

}
