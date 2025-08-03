//go:build !solution

package hotelbusiness

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	var result = []Load{}

	//Упорядоченный массив по дате, в котором хранится изменение гостей
	var guestChange = []Load{}

	dateAndGuestCountMap := make(map[int]int)

	for _, v := range guests {
		dateAndGuestCountMap[v.CheckInDate]++
		dateAndGuestCountMap[v.CheckOutDate]--
	}

	for k, v := range dateAndGuestCountMap {
		n := len(guestChange)

		i := 0

		//Поиск индекса в сортированном массиве по StartDate
		for ; i < n && guestChange[i].StartDate < k; i++ {
		}

		newElem := Load{
			StartDate:  k,
			GuestCount: v,
		}

		if i != n {
			guestChange = append(guestChange, guestChange[n-1])

			//Сдвиг элементов справа при вставке
			for j := n - 1; j > i; j-- {
				guestChange[j] = guestChange[j-1]
			}

			guestChange[i] = newElem
		} else {
			guestChange = append(guestChange, newElem)
		}
	}

	//Создание result
	var guestCount int

	for _, v := range guestChange {
		if v.GuestCount != 0 {
			guestCount += v.GuestCount

			result = append(result, Load{
				StartDate:  v.StartDate,
				GuestCount: guestCount,
			})
		}
	}

	return result
}
