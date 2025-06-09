package interfaces

type Validator interface {
	Validate(i interface{}) error
}
