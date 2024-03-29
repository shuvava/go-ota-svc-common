package data_test

import (
	"encoding/json"
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
		newID, err := data.CorrelationIDFromString(str)
		if err != nil {
			t.Errorf("CorrelationID.FromString() returned error: %v", err)
		}
		if id != newID {
			t.Error("CorrelationID.String() and FromString() did not return the same value")
		}
	})
	t.Run("CorrelationID should be serializable to valid JSON", func(t *testing.T) {
		//var id *data.CorrelationID
		id := data.NewCorrelationID()
		b, err := json.Marshal(&id)
		if err != nil {
			t.Errorf("json.Marshal returned error: %v", err)
		}
		if string(b) == "" || len(b) != 38 {
			t.Error("JSON string is not valid")
		}
		var newID data.CorrelationID
		err = json.Unmarshal(b, &newID)
		if err != nil {
			t.Errorf("json.Unmarshal returned error: %v", err)
		}
		if id != newID {
			t.Error("CorrelationID.String() and FromString() did not return the same value")
		}
	})
}

func TestChildCorrelationID(t *testing.T) {
	t.Run("NewChildCorrelationID should create consistent CorrelationID", func(t *testing.T) {
		parent := data.NewCorrelationID()
		id := data.NewCorrelationID()
		child1 := data.NewChildCorrelationID(parent, id.String())
		child2 := data.NewChildCorrelationID(parent, id.String())
		if child1 != child2 {
			t.Error("NewChildCorrelationID is not consistent")
		}
	})
}
