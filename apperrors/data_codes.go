package apperrors

const (
	// ErrorDataSerialization is error type returned if data.ObjectID was not serialized successfully
	ErrorDataSerialization = ErrorNamespaceData + ":Serialization"
	// ErrorDataValidation is error type returned if data.Ref is not pass validation
	ErrorDataValidation = ErrorNamespaceData + ":Validation"
)
