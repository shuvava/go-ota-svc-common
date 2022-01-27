package data

import "github.com/google/uuid"

// CorrelationID wrapper on the top of github.com/google/uuid
type CorrelationID uuid.UUID

func (c CorrelationID) String() string {
	return uuid.UUID(c).String()
}

// NewCorrelationID creates a new CorrelationID
func NewCorrelationID() CorrelationID {
	id := uuid.New()
	return CorrelationID(id)
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
