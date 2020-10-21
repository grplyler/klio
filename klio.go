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

func NewKlio() *Klio {
	return &Klio{}
}

func (k *Klio) AddProtocol(name, format string) {
	if Contains(ALLOWED_FORMATS, format) == false {
		log.Fatalf("Unsupported protocol message format: %s\n", format)
	}
	var events = make(map[string]ProtocolHandler, 0)
	var proto = &Protocol{Name: name, Format: format, Events: events}
	k.Proto = proto
}

func (k *Klio) On(msg string, handler ProtocolHandler) {
	k.Proto.Events[msg] = handler
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
	if len(k.Proto.Events) == 0 {
		log.Fatalf("Protocol '%s' has no message handlers.\n", k.Proto.Name)
	}

	// Create Listening Socket
	socket, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalln(err)
	}
	// defer socket.Close()

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
	fmt.Fprintf(conn, "{\"msg\": \"hi\"}")

	// reader := bufio.NewReader(conn)
	// buf := make([]byte, 1024)
	// reader.Read(buf)
	// fmt.Println(buf)

	// Todo: Wait for messages

	// e.Encode("{\"msg\": \"hi\"}")

	// Todo: Make concurrent and use channels

	k.HandleConnection(conn)
}
