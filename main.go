package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	s := New()
	err := s.Run()
	if err != nil {
		log.Fatal().Err(err)
	}
}

// Message is a struct for the message
type Message struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Data  string `json:"message"`
}

// Server is a struct for the server. Just a listener.
type Server struct {
	ln net.Listener
}

func New() *Server {
	return &Server{}
}

// Run starts the server. It listens on the given network and address and starts a goroutine for each connection.
func (s *Server) Run() error {
	err := s.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	for {
		// get connection from a client
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.HandleConnection(conn)
	}
}

// Listen fills the Server's listener.
func (s *Server) Listen(network, addr string) error {
	ln, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	s.ln = ln
	return nil
}

// HandleConnection handles a connection. It reads the message. If the message is "START", it sends a Message as a stream of bytes to the client.
func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Error().Err(err)
	}

	log.Debug().Msgf("request's message: %s", message)

	if message == "START" || message == "START\n" || message == "START\r" || message == "START\r\n" {
		newmessage := "Starting notification service..."
		conn.Write([]byte(newmessage + "\n")) // send back
	} else {
		newmessage := "Goodbye, client!"
		conn.Write([]byte(newmessage + "\n")) // send back
		return
	}

	for {
		messageJSON := Message{
			Type:  "push",
			Title: "Current time",
			Data:  fmt.Sprint(time.Unix(time.Now().Unix(), 0).Format(time.RFC3339)),
		}
		message, _ := json.Marshal(messageJSON)
		_, err := conn.Write(append(message, '\n'))
		if err != nil {
			return
		}
		log.Debug().Msgf("Message sent to %s", conn.RemoteAddr().String())
		time.Sleep(time.Second)
	}
}
