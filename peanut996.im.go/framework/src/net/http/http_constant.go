package http

// HttpMethod ...
type HTTPMethod string

const (
	HTTPMethodGET    HTTPMethod = "GET"
	HTTPMethodPOST   HTTPMethod = "POST"
	HTTPMethodPUT    HTTPMethod = "PUT"
	HTTPMethodPATCH  HTTPMethod = "PATCH"
	HTTPMethodDELETE HTTPMethod = "DELETE"
	HTTPMethodHEAD   HTTPMethod = "HEAD"
)
