package main

import (
	"SemaforoDistribuido/pkg/semaphoreclient"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	host := os.Args[1]
	cluster := os.Args[2]
	time.Sleep(time.Second * 10 )
	ports := []string{ host+":10000", host+":10001", host+":10002"}
	//var semaphores []*semaphoreclient.Semaphore
	semaphores := make([]*semaphoreclient.Semaphore, 0)
	for _, port := range ports{
		semaphores = append(semaphores, semaphoreclient.New(port, cluster))
	}

	// -----------------------
	for {
		for _, sem := range semaphores {
			go simulateCSAccess(sem)
		}
		time.Sleep(time.Second * 8)
		log.Printf("---------------Another round----------------")
	}
}

func simulateCSAccess(sem *semaphoreclient.Semaphore) {
	rand.Seed(time.Now().UnixNano())
	log.Printf("%v: Esperando SC\n", sem.Me())
	wait, err := sem.Wait("g")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	<-wait
	log.Printf("%v: Entrando SC\n", sem.Me())
	time.Sleep(time.Second * time.Duration(rand.Intn(4)))

	err = sem.Signal("g")
	for err != nil {
		log.Printf("Error: %v", err)
		err = sem.Signal("g")
		time.Sleep(time.Second * 2)
	}
	log.Printf("%v: Liberando SC\n", sem.Me())
}


