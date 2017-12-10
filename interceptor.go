package winter

type Interceptor interface {
	before(Request, Response) bool
	after(Response, Request) Exception
	done(Response, Request, error)
}

type Validator interface {
}
