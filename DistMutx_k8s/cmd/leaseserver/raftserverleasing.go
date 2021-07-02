package main

import (
	"SemaforoDistribuido/pkg/common"
	"fmt"
	"github.com/tidwall/uhaha"
	"log"
	"net/rpc"
)

const (
	NOBODY = ""
)

type Mutexes struct {
	mtxs map[string]*mutex
}

type mutex struct {
	common.Queue
	inside string
}

func (m *mutex) canLock() bool {
	return m.inside == NOBODY
}

func (m *mutex) release() {
	if m.Queue.HasNext() {
		m.inside = m.Queue.Remove().(string)
		sendWakeUp(m.inside, "")

	} else {
		fmt.Printf("Nobody inside\n")
		m.inside = NOBODY
	}
}

func sendWakeUp(dest string, idOp string) {

	client, err := rpc.DialHTTP("tcp", string(dest))

	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	//fmt.Printf("Calling RPC\n")
	var tmp struct{}
	err = client.Call(
		"SemaphoreRPC.WakeUp",
		idOp,
		&tmp,
	)
	if err != nil {
		log.Printf("Send error, %v", err)
	}
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
			Queue:  common.Queue{},
			inside: "",
		}
		mutexes.mtxs[idMutex] = mtx
	}

	if mtx.canLock() {
		sendWakeUp(client, "")
		fmt.Printf("Lock '%v' acquired by '%v'\n", idMutex, client)
		mtx.inside = client
		return common.GO, nil
	} else {
		if mtx.inside != client {
			mtx.Add(client)
		}

		return common.QUEUED, nil
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
	conf.Name = "Distributed Mutexes"
	conf.AddWriteCommand("wait", cmdWAIT)
	conf.AddWriteCommand("signal", cmdSIGNAL)

	uhaha.Main(conf)
}
