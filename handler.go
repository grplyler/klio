package klio

import (
	"encoding/json"
	"io"
	"log"
	"net"
)

func (k *Klio) Dispatch(event string, ctx *Context) {
	log.Printf("Dispatch Handling Event: '%s'\n", event)

	// Call Protocol Event Handler Message
	handler := ctx.Klio.Proto.Events[event]
	handler(ctx) // todo: handle acks

}

func (k *Klio) HandleConnection(conn net.Conn) {
	defer log.Printf("Connection from %s closing.", conn.RemoteAddr().String())
	defer conn.Close()
	client := conn.RemoteAddr().String()
	// notify := make(chan error)

	// Process Messages
	var packet map[string]interface{}
	for {

		d := json.NewDecoder(conn)
		err := d.Decode(&packet)
		if err != nil {
			if err == io.EOF {
				log.Println("Transmission finished with EOF. Closing.")
				break
			} else {
				log.Println("Error %s packet:", err)
				break
			}
		}

		log.Printf("Got Packet: %s\n", packet)

		if packet["msg"] == "exit" {
			log.Println("%s requested disconnect.", client)
			break
		}

		event := packet["msg"].(string)

		// Create Content
		context := &Context{
			Klio:    k,
			Event:   event,
			Message: packet,
			Client:  client,
		}

		// Check that this is a valid message we handle handle
		verr := k.Proto.Validate(context)
		if verr != nil {
			log.Println(verr.Message)
		} else {
			// Dispatch
			k.Dispatch(event, context)
		}

	}

	log.Println("Got Connection from:", conn.LocalAddr().String())
}
