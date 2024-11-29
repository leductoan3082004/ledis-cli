package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"ledis-cli/http_utils"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	host string
	port string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ledis-cli",
		Short: "A simple Redis-like CLI",
		Run: func(cmd *cobra.Command, args []string) {
			connectAndRunInteractiveMode()
		},
	}

	rootCmd.PersistentFlags().StringVar(&host, "host", "127.0.0.1", "Redis server host")
	rootCmd.PersistentFlags().StringVar(&port, "port", "6379", "Redis server port")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func connectAndRunInteractiveMode() {
	address := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Connected to %s\n", address)

	_, err := http_utils.GetHttpClient().Get(fmt.Sprintf("http://%s/ping", address), map[string]string{})
	if err != nil {
		fmt.Printf("Error cannot connect to address %s: %v\n", address, err)
		return
	}
	startPrompt(address)
}

func startPrompt(address string) {
	scanner := bufio.NewScanner(os.Stdin)
	prompt := fmt.Sprintf("%s> ", address)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	for {
		select {
		case <-sigChan:
			fmt.Println("\nReceived SIGINT, exiting...")
			return

		default:
			fmt.Print(prompt)
			if scanner.Scan() {
				input := scanner.Text()

				if input == "exit" || input == "quit" {
					fmt.Println("Exiting...")
					return
				}

				strs := strings.Split(input, " ")
				if len(strs) < 1 {
					fmt.Println("\nPlease enter the command")
					continue
				}

				command := strs[0]
				args := strs[1:]

				resp, err := http_utils.GetHttpClient().Post(
					fmt.Sprintf("http://%s", address), map[string]string{}, map[string]interface{}{
						"command": command,
						"args":    args,
					},
				)

				if err != nil {
					fmt.Printf("%v\n", err)
					continue
				}

				fmt.Printf("%v\n", resp)
			}
		}
	}
}
