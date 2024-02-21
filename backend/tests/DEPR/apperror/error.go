package apperror

import "encoding/json"

type MarshalJSON interface {
	UnmarshalJSON(b []byte) error
	MarshalJSON() ([]byte, error)
}

type AppError struct {
	Err     error           `json:"error"`
	Comment string          `json:"comment"`
	Details json.RawMessage `json:"details,omitempty"`
}

func (e *AppError) MarshalJSON() ([]byte, error) {
	m, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (e *AppError) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, e)
	return err
}

func NewAppError(err error, comment string, details json.RawMessage) AppError {
	return AppError{
		Err:     err,
		Comment: comment,
		Details: details,
	}
}
