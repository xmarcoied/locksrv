package main

import (
	"flag"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/xmarcoied/locksrv/lock"
	"github.com/xmarcoied/locksrv/server"
)

var tcpport string

func init() {
	// log.SetReportCaller(true)
	flag.StringVar(&tcpport, "tcpport", "5010", "defines port for tcp server")
	flag.Parse()
}

func main() {
	var Manager = lock.New()
	Manager.Init()
	// q := utils.NewQueue()
	// log.Println(q.IsEmpty())

	// marco := utils.NewQueue()
	// marco.Push("Mohsen")

	// // q.Push("Marco")
	// q.Push(marco)

	// ss := q.Pop()
	// log.Println(ss)
	// a, x := ss.(utils.Queue)
	// log.Println(a, x)

	var tcpsrv server.Server = server.NewTCPServer(tcpport, Manager)

	go func() {
		if err := tcpsrv.ListenAndServer(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)

	// Accepting graceful shutdowns when receiving SIGINT (CTRL + C)
	signal.Notify(c, os.Interrupt)

	// Blocking until receiving the signal
	<-c

	// Create a deadline to wait for.
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	// // tcpsrv.Shutdown(ctx)
	log.Warn("shutting down")

	os.Exit(0)
}
