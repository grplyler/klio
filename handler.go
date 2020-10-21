package klio

import (
	"encoding/json"
	"fmt"
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
	// defer log.Printf("Connection from %s closing.", conn.RemoteAddr().String())
	// defer conn.Close()
	client := conn.RemoteAddr().String()
	fmt.Println(conn)
	// notify := make(chan error)

	// Process Messages
	var packet map[string]interface{}
	for {

		d := json.NewDecoder(conn)
		err := d.Decode(&packet)

		fmt.Println("Packet:", packet)

		if err != nil {
			fmt.Println("ERRROR:", err)
			if err == io.EOF {
				log.Println("Transmission finished with EOF. Closing.")
				break
			} else {
				log.Printf("Error: %s", err)
				break
			}
		}

		log.Printf("Got Packet: %s\n", packet)

		if packet["msg"] == "exit" {
			log.Printf("%s requested disconnect.\n", client)
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

func (k *Klio) ClientHandleConnection(conn net.Conn) {
	d := json.NewDecoder(conn)
	var packet map[string]interface{}
	derr := d.Decode(&packet)
	if derr != nil {
		log.Fatalln(derr)
	}
	fmt.Println(packet)

}
