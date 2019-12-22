package lock

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/xmarcoied/locksrv/utils"
)

// Manager defines a struct responsible for managing [Clients, Resources, Components]
type Manager struct {
	sync.RWMutex                        // For Read/Write Mutex
	State        map[string]string      // State used for stating which resource takes a client
	ClientHolder map[string]utils.Set   // ClientHolder used for listing resources hold by each client
	ClientQueue  map[string]utils.Queue // ClientQueue used for listing client queue requests for locking each resource
}

type ClientResource struct {
	Client   string
	Resource string
}

// New generates a new Manager
func New() Manager {
	return Manager{
		State:        make(map[string]string),
		ClientHolder: make(map[string]utils.Set),
		ClientQueue:  make(map[string]utils.Queue),
	}
}

var ClientQueue = make(chan ClientResource)
var ReleasedSignal = make(chan string)

// Init fires the background workers needed
func (l Manager) Init() {
	go l.HandleQueuing()
	go l.Signal()
}

func (l *Manager) getLock(resource string) (client string, found bool) {
	l.RLock()
	defer l.RUnlock()

	client, found = l.State[resource]

	return
}

func (l *Manager) setLock(client string, resource string) {
	l.Lock()
	defer l.Unlock()

	l.State[resource] = client

	set, ok := l.ClientHolder[client]
	if ok {
		l.ClientHolder[client] = set.Add(resource)
	} else {
		l.ClientHolder[client] = utils.NewSet().Add(resource)
	}
}

// LockResource tries to lock a resource by a certain client
func (l *Manager) LockResource(client string, resource string) string {
	log.Infof("Lock request from client [%s] to resource [%s]", client, resource)
	holder, found := l.getLock(resource)
	if found == false {
		l.setLock(client, resource)
		return "SUCCESS: Resource is locked"
	} else if holder == client {
		return "ERROR: Resource is already locked"
	}

	ClientQueue <- ClientResource{
		Client:   client,
		Resource: resource,
	}
	return "Info: Lock request is queued"
}

// ReleaseResource tries to release a lock from a resource by a certain client
func (l *Manager) ReleaseResource(client string, resource string) string {
	log.Infof("Release request from client [%s] to resource [%s]", client, resource)
	holder, found := l.getLock(resource)

	if found == false {
		log.Errorf("Resource [%v] is not found", resource)
		return "ERROR: Resource is not found"
	}

	if found && holder == client {
		l.deleteLock(client, resource)
		log.Infof("Resource [%v] released", resource)
		ReleasedSignal <- resource
		return "SUCCESS: Resource is released"
	}

	log.Errorf("Resource [%v] is locked by another resource [%v]", resource, holder)
	return "ERROR: Released is hold by another resource"
}

func (l *Manager) deleteLock(client, resource string) {
	l.Lock()
	defer l.Unlock()
	_, found := l.State[resource]
	if found == true {
		delete(l.State, resource)
	}

	clientHolder, found := l.ClientHolder[client]
	if found == true {
		clientHolder.Remove(resource)
	}
}

// ReleaseResources releases all resources attached with one client
func (l *Manager) ReleaseResources(client string) {
	// Looping over the set
	for r := range l.ClientHolder[client] {
		msg := l.ReleaseResource(client, r.(string))
		log.Println(msg)
	}
}

func (l *Manager) HandleQueuing() {
	for {
		select {
		case c := <-ClientQueue:
			fmt.Println("received", c)
			q, ok := l.ClientQueue[c.Resource]
			if ok {
				l.ClientQueue[c.Resource] = q.Push(c.Client)
			} else {
				l.ClientQueue[c.Resource] = utils.NewQueue().Push(c.Client)
			}
		}
	}
}

func (l *Manager) Signal() {
	for {
		select {
		case c := <-ReleasedSignal:
			fmt.Println("received", c)
			// Fire signal events at here
			q := l.ClientQueue[c]
			if q.IsEmpty() == false {
				log.Println(q, c)
				l.LockResource(q.Pop().(string), c)
			}
		}
	}
}
