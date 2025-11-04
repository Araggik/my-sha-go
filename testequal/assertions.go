//go:build !solution

package testequal

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	isAssert := checkEqual(expected, actual)

	if !isAssert {
		raiseError(t, msgAndArgs, false)
	}

	return isAssert
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	isAssert := !checkEqual(expected, actual)

	if !isAssert {
		raiseError(t, msgAndArgs, false)
	}

	return isAssert
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	isAssert := checkEqual(expected, actual)

	if !isAssert {
		raiseError(t, msgAndArgs, true)
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	isAssert := !checkEqual(expected, actual)

	if !isAssert {
		raiseError(t, msgAndArgs, true)
	}
}

func raiseError(t T, msgAndArgs []any, isFail bool) {
	t.Helper()

	l := len(msgAndArgs)

	var str string
	var ok bool

	if l > 0 {
		str, ok = msgAndArgs[0].(string)

		if ok {
			t.Errorf(str, msgAndArgs[1:l]...)

			if isFail {
				t.FailNow()
			}
		}
	}

	if !ok {
		t.Errorf("")

		if isFail {
			t.FailNow()
		}
	}
}

func receiveIsEqualForSimpleType(ex any, a any, ok bool) bool {
	if ok {
		return ex == a
	} else {
		return false
	}
}

func receiveIsEqualForIntSlice(ex []int, a []int) bool {
	if len(ex) != len(a) {
		return false
	} else if ex == nil && a == nil {
		return true
	} else if (ex == nil && a != nil) || (a == nil && ex != nil) {
		return false
	}

	isEqual := true

	l := len(ex)

	for i := 0; i < l && isEqual; i++ {
		isEqual = ex[i] == a[i]
	}

	return isEqual
}

func receiveIsEqualForStrSlice(ex []string, a []string) bool {
	if len(ex) != len(a) {
		return false
	} else if ex == nil && a == nil {
		return true
	} else if (ex == nil && a != nil) || (a == nil && ex != nil) {
		return false
	}

	isEqual := true

	l := len(ex)

	for i := 0; i < l && isEqual; i++ {
		isEqual = ex[i] == a[i]
	}

	return isEqual
}

func receiveIsEqualForByteSlice(ex []byte, a []byte) bool {
	if len(ex) != len(a) {
		return false
	} else if ex == nil && a == nil {
		return true
	} else if (ex == nil && a != nil) || (a == nil && ex != nil) {
		return false
	}

	isEqual := true

	l := len(ex)

	for i := 0; i < l && isEqual; i++ {
		isEqual = ex[i] == a[i]
	}

	return isEqual
}

func receiveIsEqualForMap(ex map[string]string, a map[string]string) bool {
	if len(ex) != len(a) {
		return false
	} else if ex == nil && a == nil {
		return true
	} else if (ex == nil && a != nil) || (a == nil && ex != nil) {
		return false
	}

	isEqual := true

loop:
	for k, v := range ex {
		aVal, ok := a[k]

		isEqual = ok && (aVal == v)

		if !isEqual {
			break loop
		}
	}

	return isEqual
}

func checkEqual(expected any, actual any) bool {
	isEqual := true

	switch ex := expected.(type) {
	case int:
		a, ok := actual.(int)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case string:
		a, ok := actual.(string)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case int8:
		a, ok := actual.(int8)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case int16:
		a, ok := actual.(int16)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case int32:
		a, ok := actual.(int32)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case int64:
		a, ok := actual.(int64)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case uint8:
		a, ok := actual.(uint8)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case uint16:
		a, ok := actual.(uint16)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case uint32:
		a, ok := actual.(uint32)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)
	case uint64:
		a, ok := actual.(uint64)

		isEqual = receiveIsEqualForSimpleType(ex, a, ok)

	case []int:
		a, ok := actual.([]int)

		if ok {
			isEqual = receiveIsEqualForIntSlice(ex, a)
		} else {
			isEqual = false
		}
	case []string:
		a, ok := actual.([]string)

		if ok {
			isEqual = receiveIsEqualForStrSlice(ex, a)
		} else {
			isEqual = false
		}
	case []byte:
		a, ok := actual.([]byte)

		if ok {
			isEqual = receiveIsEqualForByteSlice(ex, a)
		} else {
			isEqual = false
		}
	case map[string]string:
		a, ok := actual.(map[string]string)

		if ok {
			isEqual = receiveIsEqualForMap(ex, a)
		} else {
			isEqual = false
		}
	default:
		isEqual = false
	}

	return isEqual
}
