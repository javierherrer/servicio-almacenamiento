package semaphoreclient

import (
	"SemaforoDistribuido/pkg/common"
	"encoding/gob"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net"
)

type Semaphore struct {
	serverHost string
	me         string
	wakeUp     chan struct{}
	idOp       string
	signaled bool
}

// Returns a new semaphore
func New(host, cluster string) *Semaphore {
	sem := Semaphore{
		serverHost: cluster,
		me:         host,
		wakeUp:     make(chan struct{}, 1),
	}
	sem.updateLeader()
	go sem.startTCPServer()
	return &sem

}

// Calls the raft server to get the lock of the specified distributed mutex
// Returns a channel to wait on
func (s *Semaphore) Wait(semaphore string) (chan struct{}, error) {
	s.updateLeader()
	c, err := redis.Dial("tcp", s.serverHost)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	s.signaled = false // No se nos ha avisado aún
	r, err := redis.String(c.Do("wait", semaphore, s.me))

	if err != nil {
		return nil, err
	}
	log.Printf("%v: Got response: %s\n",s.me, r)
	return s.wakeUp, nil
}

// Calls the raft server to release the lock of the specified distributed mutex
func (s *Semaphore) Signal(semaphore string) error{
	s.updateLeader()
	c, err := redis.Dial("tcp", s.serverHost)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	defer c.Close()

	r, err := redis.String(c.Do("signal", semaphore, s.me))
	if err != nil {
		log.Printf("%v\n", err)
	}
	fmt.Printf("%v: Got response:  %s\n", s.me, r)
	return nil
}

// Starts the tcp server to wait for the signal of the dist. mutex
func (s *Semaphore) startTCPServer() {
	listener, err := net.Listen("tcp", s.me)

	if err != nil {
		log.Printf("listen tcp error: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept tcp error: %v", err)
		}
		s.processConn(conn)
	}

}

// Checks if the msg was a wakeup msg
func (s *Semaphore) processConn(conn net.Conn) {
	var msg string

	decoder := gob.NewDecoder(conn)
	decoder.Decode(&msg)
	if msg == common.WakeUpMsg {

		if s.signaled {
			return
		}
		s.wakeUp <- struct{}{}

		s.signaled = true
	}
}



// Calls the raft server to release the lock of the specified distributed mutex
func (s *Semaphore) updateLeader() {
	c, err := redis.Dial("tcp", s.serverHost)

	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	defer c.Close()

	leader, err := redis.String(
		c.Do("RAFT", "LEADER"))
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	s.serverHost = leader
}

func (s *Semaphore) Me() string {
	return s.me
}