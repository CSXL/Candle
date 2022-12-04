package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/CSXL/Candle/api"
	input "github.com/CSXL/Candle/util/input"
)

// We can mock these functions in our tests
var osExit = os.Exit
var NewScanner input.ScannerInterface = input.NewScanner()

// When true the main function will exit after the first command
var TESTING_MODE = false

func main() {
	var COMMAND_PROVIDED = false
	if len(os.Args) > 1 {
		COMMAND_PROVIDED = true
	}

	scanner := NewScanner

	for {
		fmt.Print("--> ")
		var args []string
		if COMMAND_PROVIDED {
			args = os.Args[1:]
		} else {
			scanner.Scan()
			args = strings.Split(scanner.Text(), " ")
		}

		switch args[0] {
		case "about":
			fmt.Println("Candle CLI")
			fmt.Println("CSX Labs")
			fmt.Println("Made w/ <3 by @absozero and @ecsbeats")
		case "get-info":
			stock_data := api.GetStockData(args[1])
			// Prints first 10 records
			var i int = 0
			for key, value := range stock_data {
				if i == 10 {
					break
				}
				fmt.Println(key, value)
				i++
			}

		case "exit":
			fmt.Println("Exiting...")
			osExit(0)
		default:
			fmt.Printf("Invalid command: %s\n", args[0])
		}

		if COMMAND_PROVIDED || TESTING_MODE {
			break
		}
	}
}
