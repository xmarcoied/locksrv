package server

import (
	"bufio"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/xmarcoied/locksrv/lock"
)

// TCPServer defines a locking server in TCP protocol
type TCPServer struct {
	ln      *net.TCPListener
	port    string
	Manager lock.Manager
}

// NewTCPServer creates a new instance of TCPServer datastructure
func NewTCPServer(p string, l lock.Manager) TCPServer {
	return TCPServer{
		port:    p,
		Manager: l,
	}
}

// ListenAndServer ...
func (s TCPServer) ListenAndServer() error {
	tcpaddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		log.Error(err)
		return err
	}

	s.ln, err = net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Start listening to tcp port", tcpaddr)

	err = s.acceptConnections()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s TCPServer) acceptConnections() error {
	for {
		tcpconn, err := s.ln.AcceptTCP()
		if err != nil {
			log.Error(err)
			return err
		}

		log.Info("New tcp connection", tcpconn.RemoteAddr())
		go s.client(tcpconn)
	}
}

// Shutdown close all tcp connections
func (s TCPServer) Shutdown() {}

// client handles an incoming tcp connection
func (s TCPServer) client(c *net.TCPConn) {
	// var clientID string = uuid.New().String()
	var clientID string = c.RemoteAddr().String()
	for {
		msg, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Errorf("Client [%v] stopped listening due error: [%v]", c.RemoteAddr(), err)
			// Sending release all
			s.Manager.ReleaseResources(clientID)
			break
		}

		srvMsg := serverMessage(msg)

		log.Infof("Received msg [%s] from client [%v]", srvMsg.String(), c.RemoteAddr())
		if srvMsg.Valid() == false {
			log.Warnf("Received msg: [%s] from client [%v] isn't valid", srvMsg.String(), c.RemoteAddr())
		} else {
			if srvMsg.Action() == "lock" {
				// response := s.Manager.LockResource(clientID, srvMsg.Resource())
				go func(c *net.TCPConn) {
					response := s.Manager.LockResource(clientID, srvMsg.Resource())
					c.Write([]byte(response + "\n"))
				}(c)
				// c.Write([]byte(response + "\n"))
			}

			if srvMsg.Action() == "unlock" {
				go func(c *net.TCPConn) {
					response := s.Manager.ReleaseResource(clientID, srvMsg.Resource())
					c.Write([]byte(response + "\n"))
				}(c)
				// response := s.Manager.ReleaseResource(clientID, srvMsg.Resource())
				// c.Write([]byte(response + "\n"))
			}

		}
	}
}
