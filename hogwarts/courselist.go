//go:build !solution

package hogwarts


//Test
//Test2
func GetCourseList(prereqs map[string][]string) []string {
	var result = []string{}

	setForCheck := make(map[string]struct{})

	n := len(prereqs)

	//Поиск курсов без зависимых курсов
	for k, v := range prereqs {
		if len(v) == 0 {
			setForCheck[k] = struct{}{}

			result = append(result, k)
		}

		for _, c := range v {
			_, ok := prereqs[c]
			_, ok2 := setForCheck[c]

			if !ok && !ok2 {
				setForCheck[c] = struct{}{}

				result = append(result, c)
			}
		}
	}

	if len(setForCheck) == 0 && n != 0 {
		panic("Нет курса без зависимых курсов")
	}

	//Подсчет всех курсов
	allCourseSet := make(map[string]struct{})

	for k, v := range prereqs {
		allCourseSet[k] = struct{}{}

		for _, c := range v {
			allCourseSet[c] = struct{}{}
		}
	}

	allCourseN := len(allCourseSet)

	for len(setForCheck) != allCourseN {
		addedCourseSet := make(map[string]struct{})

		for k, v := range prereqs {
			if _, ok := setForCheck[k]; !ok {
				isAdded := true

				l := len(v)

				for i := 0; isAdded && i < l; i++ {
					if _, ok := setForCheck[v[i]]; !ok {
						isAdded = false
					}
				}

				if isAdded {
					addedCourseSet[k] = struct{}{}
				}
			}
		}

		if len(addedCourseSet) == 0 {
			panic("Циклическая зависимость")
		} else {
			for k := range addedCourseSet {
				setForCheck[k] = struct{}{}

				result = append(result, k)
			}
		}
	}

	return result
}
