package data

// Namespace is object namespace
type Namespace string

// NewNamespace creates namespace with random uuid
func NewNamespace(id string) Namespace {
	if id == "" {
		id = NewNamespaceURN(NewCorrelationID())
	}

	return Namespace(id)
}
