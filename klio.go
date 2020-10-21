package klio

import (
	"fmt"
	"log"
	"net"
)

type Klio struct {
	Proto *Protocol
	Mode  string
}

type H map[string]string

func NewKlio() *Klio {
	klio := &Klio{}
	var events = make(map[string]ProtocolHandler, 0)

	// Add Default Protocol //todo: decide if we want to do this
	var proto = &Protocol{Name: "default", Format: "json", Handlers: events}
	klio.Proto = proto
	return klio

}

func (k *Klio) AddProtocol(name, format string) {
	if Contains(ALLOWED_FORMATS, format) == false {
		log.Fatalf("Unsupported protocol message format: %s\n", format)
	}
	var events = make(map[string]ProtocolHandler, 0)
	var proto = &Protocol{Name: name, Format: format, Handlers: events}
	k.Proto = proto
}

// Same as klio.Proto.AddHandler
func (k *Klio) On(event string, handler ProtocolHandler) {
	fmt.Println("Adding Protocol Handler:", event)
	k.Proto.Handlers[event] = handler
}

func (k *Klio) Serve(addr string) {

	if k.Mode != "client" {
		k.Mode = "client"
	}

	// Check that we have protocols defined
	if k.Proto == nil {
		log.Println("No protocols defined. Cannot start server.")
		log.Fatalln("Use Klio.NewProtocol(format string)")
	}

	// Check that there is at least one message handler on the protocol
	if len(k.Proto.Handlers) == 0 {
		log.Fatalf("Protocol '%s' has no message handlers.\n", k.Proto.Name)
	}

	// Create Listening Socket
	socket, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer socket.Close()

	exit := make(chan bool, 1)

	// Accept Connections
	for {

		log.Println("Klio Protocol Server Listening on", socket.Addr().String())
		conn, err := socket.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		// Handle Connections
		go k.HandleConnection(conn, exit)
	}
}

// Client
func (k *Klio) Dial(addr string) {

	fmt.Println("Dialing")

	if k.Mode != "client" {
		k.Mode = "client"
	}

	conn, err := net.Dial("tcp4", addr)

	if err != nil {
		log.Fatalln(err)
	}

	exit := make(chan bool)

	go k.HandleConnection(conn, exit)

	<-exit
	fmt.Println("Client Shutdown")

}
