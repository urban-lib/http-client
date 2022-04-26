package httpClient

type Response struct {
	Status int
	Body   []byte
	Error  error
}
