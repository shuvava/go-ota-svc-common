package data_test

import (
	"testing"

	"github.com/shuvava/go-ota-svc-common/data"
)

func TestCorrelationID_IsNil(t *testing.T) {
	t.Run("not initialized CorrelationID should be IsNul", func(t *testing.T) {
		var cid *data.CorrelationID
		if !cid.IsNil() {
			t.Error("CorrelationID.IsNil() returned false for nil")
		}
	})
	t.Run("initialized CorrelationID should not be IsNil", func(t *testing.T) {
		cid := data.NewCorrelationID()
		if cid.IsNil() {
			t.Error("CorrelationID.IsNil() returned true for non-nil")
		}
	})
	t.Run("CorrelationIDNil should be IsNUl", func(t *testing.T) {
		cid := data.CorrelationIDNil
		if !cid.IsNil() {
			t.Error("CorrelationID.IsNil() returned false for nil")
		}
	})
}

func TestCorrelationID_Serialization(t *testing.T) {
	t.Run("CorrelationID should be serializable", func(t *testing.T) {
		id := data.NewCorrelationID()
		str := id.String()
		newId, err := data.CorrelationIDFromString(str)
		if err != nil {
			t.Errorf("CorrelationID.FromString() returned error: %v", err)
		}
		if id != newId {
			t.Error("CorrelationID.String() and FromString() did not return the same value")
		}
	})
}
