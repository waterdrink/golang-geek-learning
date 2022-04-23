package error

type ErrorCode int

const (
	OK              ErrorCode = 0
	InvalidArgument ErrorCode = 1
	Internal        ErrorCode = 3
)
