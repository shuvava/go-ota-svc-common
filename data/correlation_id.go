package data

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// CorrelationID wrapper on the top of github.com/google/uuid
type CorrelationID uuid.UUID

// CorrelationIDNil is a nil value of CorrelationID
var CorrelationIDNil = CorrelationID(uuid.Nil)

// String converts CorrelationID to string
func (c CorrelationID) String() string {
	return uuid.UUID(c).String()
}

// MarshalJSON custom json serialization func
func (c *CorrelationID) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON custom json deserialization func
func (c *CorrelationID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	id, err := CorrelationIDFromString(s)
	if err != nil {
		return err
	}
	*c = id
	return nil
}

// NewCorrelationID creates a new CorrelationID
func NewCorrelationID() CorrelationID {
	id := uuid.New()
	return CorrelationID(id)
}

// CorrelationIDFromString creates a new CorrelationID from a string
func CorrelationIDFromString(s string) (CorrelationID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return CorrelationIDNil, apperrors.CreateError(apperrors.ErrorDataSerialization, "failed to parse provided string", err)
	}
	return CorrelationID(id), nil
}

// IsNil returns true if the CorrelationID is nil
func (c *CorrelationID) IsNil() bool {
	return c == nil || uuid.UUID(*c) == uuid.Nil
}

// NewChildCorrelationID create a new CorrelationID within a parent namespace
func NewChildCorrelationID(parent CorrelationID, id string) CorrelationID {
	parentID := uuid.UUID(parent)
	if id == "" {
		id = uuid.New().String()
	}
	childID := uuid.NewSHA1(parentID, []byte(id))
	return CorrelationID(childID)
}
