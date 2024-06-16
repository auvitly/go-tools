package lab

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type Context struct {
	Context context.Context
}

// NewContext - instantiating a builder for a context.
func NewContext() *Context {
	return &Context{
		Context: context.Background(),
	}
}

// NewIncomingContext - instantiating a builder for a incoming context.
func NewIncomingContext(md metadata.MD) *Context {
	return &Context{
		Context: metadata.NewIncomingContext(context.Background(), md),
	}
}

// NewOutgoingContext - instantiating a builder for a incoming context.
func NewOutgoingContext(md metadata.MD) *Context {
	return &Context{
		Context: metadata.NewOutgoingContext(context.Background(), md),
	}
}

// WithValue - add value by key to context.
func (c *Context) WithValue(key, value any) *Context {
	if c.Context == nil {
		c.Context = context.Background()
	}

	c.Context = context.WithValue(c.Context, key, value)

	return c
}
