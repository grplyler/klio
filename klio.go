package klio

import (
	"log"
	"net"
)

type Klio struct {
	Proto *Protocol
}

func NewKlio() *Klio {
	return &Klio{}
}

func (k *Klio) AddProtocol(name, format string) {
	if Contains(ALLOWED_FORMATS, format) == false {
		log.Fatalln("Unsupported protocol message format: %s", format)
	}
	var events = make(map[string]ProtocolHandler, 0)
	var proto = &Protocol{Name: name, Format: format, Events: events}
	k.Proto = proto
}

func (k *Klio) On(msg string, handler ProtocolHandler) {
	k.Proto.Events[msg] = handler
}

func (k *Klio) Serve(addr string) {

	// Check that we have protocols defined
	if k.Proto == nil {
		log.Println("No protocols defined. Cannot start server.")
		log.Fatalln("Use Klio.NewProtocol(format string)")
	}

	// Check that there is at least one message handler on the protocol
	if len(k.Proto.Events) == 0 {
		log.Fatalf("Protocol '%s' has no message handlers.\n", k.Proto.Name)
	}

	// Create Listening Socket
	socket, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer socket.Close()

	// Accept Connections
	for {

		log.Println("Klio Protocol Server Listening on", socket.Addr().String())
		conn, err := socket.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		// Handle Connections
		go k.HandleConnection(conn)
	}
}
