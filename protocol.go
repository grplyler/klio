package klio

import (
	"fmt"
	"net"
)

var ALLOWED_FORMATS = []string{"json", "msgpack"}

type Protocol struct {
	Name   string
	Events map[string]ProtocolHandler
	Format string
}

type Context struct {
	Klio    *Klio
	Event   string
	Conn    net.Conn
	Client  string
	Message map[string]interface{}
}

type ProtocolError struct {
	Protocol *Protocol
	Event    string
	Message  string
}

type ProtocolHandler func(c *Context)

func (p *Protocol) AddHandler(event string, handler ProtocolHandler) {
	fmt.Println("Adding Protocol Handler:", event)
	p.Events[event] = handler
}

// Check that we can handle this event type
func (p *Protocol) Validate(ctx *Context) *ProtocolError {

	// Simple nil map check
	if p.Events[ctx.Event] == nil {
		return &ProtocolError{
			Event:   ctx.Event,
			Message: fmt.Sprintf("Invalid event '%s' from '%s'", ctx.Event, ctx.Client),
		}
	}

	return nil
}
