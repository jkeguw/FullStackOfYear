package errors

const (
	Success         = 200
	BadRequest      = 400
	Unauthorized    = 401
	Forbidden       = 403
	NotFound        = 404
	TooManyRequests = 429
	InternalError   = 500
)

var errorMessages = map[int]string{
	BadRequest:      "Request param error",
	Unauthorized:    "Unauthorized access",
	Forbidden:       "Access Denied",
	NotFound:        "Resource does not exist",
	TooManyRequests: "Too many requests",
	InternalError:   "Internal server error",
}
