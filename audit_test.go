package mysupplyaudit

import "testing"

func TestEverything(t *testing.T) {

	s, err := NewSupplier("")
	if err != nil {
		t.Error(err)
		return
	}
	err = s.DoAudit(-1)
	if err != nil {
		t.Error(err)
		return
	}

}
