package errors

const (
	Success             = 200
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	Conflict            = 409
	Validation          = 422
	RateLimit           = 429
	TooManyRequests     = 429
	InternalError       = 500
	NotImplemented      = 501
	
	// 错误类型
	InvalidRequestError  = 400
	InvalidTokenError    = 401
	ValidationError      = 400
	InternalServerError  = 500
)

var errorMessages = map[int]string{
	BadRequest:      "Request param error",
	Unauthorized:    "Unauthorized access",
	Forbidden:       "Access Denied",
	NotFound:        "Resource does not exist",
	Conflict:        "Resource already exists",
	Validation:      "Validation failed",
	RateLimit:       "Rate limit exceeded",
	InternalError:   "Internal server error",
	NotImplemented:  "Feature not implemented",
}
