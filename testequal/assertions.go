//go:build !solution

package testequal

func checkEqual(expected, actual interface{}) bool {
	switch exp := expected.(type) {

	case int:
		act, ok := actual.(int)
		if !ok {
			return false
		}
		return exp == act

	case int8:
		act, ok := actual.(int8)
		if !ok {
			return false
		}
		return exp == act

	case int16:
		act, ok := actual.(int16)
		if !ok {
			return false
		}
		return exp == act

	case int32:
		act, ok := actual.(int32)
		if !ok {
			return false
		}
		return exp == act

	case int64:
		act, ok := actual.(int64)
		if !ok {
			return false
		}
		return exp == act

	case uint8:
		act, ok := actual.(uint8)
		if !ok {
			return false
		}
		return exp == act

	case uint16:
		act, ok := actual.(uint16)
		if !ok {
			return false
		}
		return exp == act

	case uint32:
		act, ok := actual.(uint32)
		if !ok {
			return false
		}
		return exp == act

	case uint64:
		act, ok := actual.(uint64)
		if !ok {
			return false
		}
		return exp == act

	case string:
		act, ok := actual.(string)
		if !ok {
			return false
		}
		return exp == act

	case map[string]string:
		act, ok := actual.(map[string]string)
		if !ok || act == nil || exp == nil {
			return false
		}
		if len(act) != len(exp) {
			return false
		}
		for k, v := range exp {
			if !checkEqual(v, act[k]) {
				return false
			}
		}
		return true

	case []int:
		act, ok := actual.([]int)
		if !ok || act == nil || exp == nil {
			return false
		}
		if len(act) != len(exp) {
			return false
		}
		for i := range act {
			if !checkEqual(exp[i], act[i]) {
				return false
			}
		}
		return true

	case []byte:
		act, ok := actual.([]byte)
		if !ok || act == nil || exp == nil {
			return false
		}
		if len(act) != len(exp) {
			return false
		}
		for i := range act {
			if !checkEqual(exp[i], act[i]) {
				return false
			}
		}
		return true

	default:
	}
	return false
}

func testEqualWrap(t T, expected, actual interface{}, need bool, msgAndArgs ...interface{}) bool {
	t.Helper()
	if checkEqual(expected, actual) == need {
		if len(msgAndArgs) > 0 {
			if st, ok := msgAndArgs[0].(string); ok {
				t.Errorf(st, msgAndArgs[1:]...)
			} else {
				panic("not a format string")
			}
		} else {
			t.Errorf("")
		}
		return false
	}
	return true
}

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	return testEqualWrap(t, expected, actual, false, msgAndArgs...)
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	return testEqualWrap(t, expected, actual, true, msgAndArgs...)
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !testEqualWrap(t, expected, actual, false, msgAndArgs...) {
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !testEqualWrap(t, expected, actual, true, msgAndArgs...) {
		t.FailNow()
	}
}
