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

//Функция аналог Parse
func TestingParse(address string) (AddressRange, error) {
	var result AddressRange

	numbers := strings.Split(address, ".")

	if strings.Contains(numbers[3], "/") {
		mask := strings.Split(numbers[3], "/")[1]

		if mask == "24" {
			part := toBytes(numbers[:3])

			result.Min = net.IP(addBytes(part, []byte{0}))
			result.Max = net.IP(addBytes(part, []byte{255}))
		} else if mask == "16" {
			part := toBytes(numbers[:2])

			result.Min = net.IP(addBytes(part, []byte{0, 0}))
			result.Max = net.IP(addBytes(part, []byte{255, 255}))
		} else if mask == "8" {
			part := toBytes(numbers[:1])

			result.Min = net.IP(addBytes(part, []byte{0, 0, 0}))
			result.Max = net.IP(addBytes(part, []byte{255, 255, 255}))
		} else {
			//Проверить что min правильный
			result.Min = net.IP([]byte{128, 0, 0, 0})
			result.Max = net.IP([]byte{255, 255, 255, 255})
		}

	} else if strings.Contains(numbers[3], "*") {
		if numbers[0] == "*" {
			//Проверить что min правильный
			result.Min = net.IP([]byte{0, 0, 0, 0})
			result.Max = net.IP([]byte{255, 255, 255, 255})
		} else if numbers[1] == "*" {
			part := toBytes(numbers[:1])

			result.Min = net.IP(addBytes(part, []byte{0, 0, 0}))
			result.Max = net.IP(addBytes(part, []byte{255, 255, 255}))
		} else if numbers[2] == "*" {
			part := toBytes(numbers[:2])

			result.Min = net.IP(addBytes(part, []byte{0, 0}))
			result.Max = net.IP(addBytes(part, []byte{255, 255}))
		} else {
			part := toBytes(numbers[:3])

			result.Min = net.IP(addBytes(part, []byte{0}))
			result.Max = net.IP(addBytes(part, []byte{255}))
		}
	} else if strings.Contains(numbers[3], "-") {
		//Получаем элементы после последнего "."
		lastParts := strings.Split(numbers[3], "-")

		lastNumberStr := lastParts[1]
		lastNumber, _ := strconv.Atoi(lastNumberStr)
		maxLastByte := byte(lastNumber % 256)

		strBytes := make([]string, 4)
		strBytes = append(numbers[:3], lastParts[0])

		part := toBytes(numbers[:3])

		result.Min = net.IP(toBytes(strBytes))
		result.Max = net.IP(addBytes(part, []byte{maxLastByte}))
	} else {
		part := toBytes(numbers[:4])

		result.Min = net.IP(part)
		result.Max = net.IP(part)
	}

	return result, nil
}

func toBytes(slice []string) []byte {
	res := make([]byte, len(slice))

	for i, s := range slice {
		num, _ := strconv.Atoi(s)

		res[i] = byte(num)
	}

	return res
}

func addBytes(slice []byte, add []byte) []byte {
	return append(slice, add...)
}