package httputil

// ResponseBody represents http response body
type ResponseBody struct {
	Title string      `json:"title"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}
