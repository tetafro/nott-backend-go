package response

// Error is a response error.
type Error string

// String returns string representation of the error.
func (err Error) String() string {
	return string(err)
}
