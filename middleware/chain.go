package middleware

import "net/http"

type Chain struct {
	handlers []func(http.Handler) http.Handler
}

func NewChain(handlers ...func(http.Handler) http.Handler) *Chain {
	return &Chain{handlers: handlers}
}

// NOTE: Interesting that you can implement a "Thennable" to a type
// NOTE: Saw this Chaining methiod in a book i was reading last week - Programming Microservices in GO
// NOTE: Lets use it as it becomes more visually nicer and we can avoid nesting
func (c *Chain) Then(h http.Handler) http.Handler {
	for i := len(c.handlers) - 1; i >= 0; i-- {
		h = c.handlers[i](h)
	}
	return h
}
