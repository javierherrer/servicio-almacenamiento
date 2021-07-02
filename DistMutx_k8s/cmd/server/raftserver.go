package main

import (
	"SemaforoDistribuido/pkg/common"
	"encoding/gob"
	"fmt"
	"github.com/tidwall/uhaha"
	"log"
	"net"
	"strconv"
)

const (
	NOBODY = ""
)

type Mutexes struct {
	mtxs map[string]*mutex
}

type mutex struct {
	id string
	common.Queue
	inside string
}

func (m *mutex) canLock() bool {
	return m.inside == NOBODY
}

func (m *mutex) release() {
	if m.Queue.HasNext() {
		m.inside = m.Queue.Remove().(string)
		sendWakeUp(m.inside, m.id)

	} else {
		m.inside = NOBODY
	}
}

func sendWakeUp(dest string, idMutex string) {

	conn, err := net.Dial("tcp", dest)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	//var tmp struct{}
	enconder := gob.NewEncoder(conn)
	err = enconder.Encode(common.WakeUpMsg)
	if err != nil {
		log.Printf("Send error, %v", err)
	}
	log.Printf("Lock '%v' acquired by '%v'\n", idMutex, dest)
}

func cmdWAIT(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 3 {
		return nil, uhaha.ErrWrongNumArgs
	}
	mutexes := m.Data().(*Mutexes)
	if mutexes.mtxs == nil {
		mutexes.mtxs = make(map[string]*mutex)
	}
	idMutex := args[1]
	client := args[2]

	mtx, ok := mutexes.mtxs[idMutex]

	if !ok {
		mtx = &mutex{
			id:     idMutex,
			Queue:  common.Queue{},
			inside: "",
		}
		mutexes.mtxs[idMutex] = mtx
	}

	if mtx.canLock() {
		sendWakeUp(client, "")
		log.Printf("Lock '%v' acquired by '%v'\n", idMutex, client)
		mtx.inside = client
		return common.GO , nil
	} else {
		if mtx.inside != client {
			mtx.Add(client)
		}

		return common.QUEUED + " pos:"+ strconv.Itoa(mtx.Size()), nil
	}
}

func cmdSIGNAL(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 3 {
		return nil, uhaha.ErrWrongNumArgs
	}
	mutexes := m.Data().(*Mutexes)
	idMutex := args[1]
	client := args[2]
	mutex, ok := mutexes.mtxs[idMutex]

	if !ok {
		return nil, uhaha.ErrInvalid
	}
	if mutex.inside == client {
		mutex.release()
		fmt.Printf("Lock '%v' released by '%v'\n", idMutex, client)
		return common.RELEASED, nil
	} else {
		return "false", nil
	}
}

func main() {
	var conf uhaha.Config
	conf.InitialData = new(Mutexes)
	conf.Version = "1.0"
	conf.Name = "Distributed Mutexes"
	conf.AddWriteCommand("wait", cmdWAIT)
	conf.AddWriteCommand("signal", cmdSIGNAL)

	uhaha.Main(conf)
}
