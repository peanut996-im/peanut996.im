package http

const (
	// Method
	HTTP_METHOD_GET    string = "GET"
	HTTP_METHOD_POST   string = "POST"
	HTTP_METHOD_PUT    string = "PUT"
	HTTP_METHOD_PATCH  string = "PATCH"
	HTTP_METHOD_DELETE string = "DELETE"
	HTTP_METHOD_HEAD   string = "HEAD"

	// Error
	NEW_REQUEST_ERROR     string = "New http request err: %v, url: %v"
	DO_REQUEST_ERROR      string = "Do http request err: %v, url: %v"
	DO_GET_REQUEST_ERROR  string = "Do get http request err: %v, url: %v"
	DO_POST_REQUEST_ERROR string = "Do post http request err: %v, url: %v"
	READ_RESP_BODY_ERROR  string = "Read resp body err: %v, url: %v"
	MARSHAL_JSON_ERROR    string = "Marshal json  err: %v"
	UNMARSHAL_JSON_ERROR  string = "Unmarshal json  err: %v"
)
