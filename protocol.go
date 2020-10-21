package klio

import (
	"encoding/json"
	"fmt"
	"net"
)

var ALLOWED_FORMATS = []string{"json", "msgpack"}

type Protocol struct {
	Name     string
	Handlers map[string]ProtocolHandler
	Format   string
}

type Context struct {
	Klio       *Klio
	Event      string
	Conn       net.Conn
	ClientAddr string
	Message    map[string]interface{}
}

type ProtocolError struct {
	Protocol *Protocol
	Event    string
	Message  string
}

type ProtocolHandler func(c *Context)

func (c *Context) Send(content string) {
	fmt.Fprintf(c.Conn, content)
}

func (c *Context) JSON(object interface{}) {
	enc := json.NewEncoder(c.Conn)
	enc.Encode(object)
}

func (p *Protocol) AddHandler(event string, handler ProtocolHandler) {
	fmt.Println("Adding Protocol Handler:", event)
	p.Handlers[event] = handler
}

// Check that we can handle this event type
func (p *Protocol) Validate(ctx *Context) *ProtocolError {

	// Simple nil map check
	if p.Handlers[ctx.Event] == nil {
		return &ProtocolError{
			Event:   ctx.Event,
			Message: fmt.Sprintf("Invalid event '%s' from '%s'", ctx.Event, ctx.ClientAddr),
		}
	}

	return nil
}

// Todo: make multiformat protocol unpacker
