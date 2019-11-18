package mocks

import "testing"

func CmpError(t *testing.T, err error, expErr error) {
	if err != expErr {
		t.Helper()
		t.Errorf("Unexpected error: (%T)%v", err, err)
		t.FailNow()
	}
}
