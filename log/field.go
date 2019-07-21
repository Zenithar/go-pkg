package log

// Field declare logger attributes
type Field struct {
	Name  string
	Value interface{}
}

// -----------------------------------------------------------------------------

// Error is a field builder for log attributes
func Error(err error) Field {
	return Field{
		Name:  "error",
		Value: err,
	}
}

// String is a field builder for string value
func String(name, value string) Field {
	return Field{
		Name:  name,
		Value: value,
	}
}
