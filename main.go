package main

import (
	"fmt"
	"github.com/chzyer/readline"
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
	rl, err := readline.NewEx(
		&readline.Config{
			Prompt:            fmt.Sprintf("%s> ", address),
			HistoryFile:       "./.history",
			HistorySearchFold: true,
		},
	)
	if err != nil {
		log.Fatalf("Error creating readline instance: %v", err)
	}
	defer rl.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	rl.HistoryEnable()

	for {
		select {
		case <-sigChan:
			fmt.Println("\nReceived SIGINT, exiting...")
			return
		default:
			line, err := rl.Readline()
			if err != nil {
				if err.Error() == "interrupt" {
					return
				}
				fmt.Printf("Error reading input: %v\n", err)
				continue
			}

			line = strings.TrimSpace(line)

			if line == "exit" || line == "quit" {
				fmt.Println("Exiting...")
				return
			}

			strs := strings.Split(line, " ")
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
