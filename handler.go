package klio

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

func (k *Klio) Dispatch(event string, ctx *Context) {
	log.Printf("Dispatch Handler Event: '%s'\n", event)

	verr := k.Proto.Validate(ctx)
	if verr != nil {
		log.Println("Error during protocol validation:", verr)
	} else {
		// Call Protocol Event Handler Message
		handler := ctx.Klio.Proto.Handlers[event]
		if handler == nil {
			log.Printf("Protocol Error: No handler defined for %s", event)
		} else {
			handler(ctx)
		}
	}

}

func (k *Klio) HandleConnection(conn net.Conn, exit chan bool) {
	defer log.Printf("Connection from %s closing.", conn.RemoteAddr().String())
	defer conn.Close()
	client := conn.RemoteAddr().String()

	// Send Startup Message
	context := &Context{
		Klio:   k,
		Event:  "_connect",
		Client: client,
		Conn:   conn,
	}

	// Send startup message
	k.Dispatch("_connect", context)

	// Process Messages
	var packet map[string]interface{}
	for {

		// Todo: Grab Raw Data asnd Store in Context
		d := json.NewDecoder(conn)
		err := d.Decode(&packet)

		if err != nil {
			fmt.Println("ERRROR:", err)
			if err == io.EOF {
				log.Println("Transmission finished with EOF. Closing.")
				exit <- true
				break
			} else {
				log.Printf("Error: %s", err)
				break
			}
		}

		log.Printf("Got Packet: %s\n", packet)

		if packet["msg"] == "exit" {
			log.Printf("%s requested disconnect.\n", client)
			exit <- true
			break
		}

		event := packet["msg"].(string)

		// Create Content
		context := &Context{
			Klio:    k,
			Event:   event,
			Message: packet,
			Client:  client,
			Conn:    conn,
		}

		// Check that this is a valid message we handle handle
		k.Dispatch(event, context)

	}

	log.Println("Got Connection from:", conn.LocalAddr().String())
}
