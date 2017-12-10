package winter

type Response interface {
	Write([]byte) (int, error)
}

type ResponseFormat interface {
}
