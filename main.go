package main

// Write a program that creates a shell that responds to the "about", "exit", and "get-info <stock>" commands.
// This program should be written with only standard libraries.
// The "about" command should print out the name of the program "Candle CLI" and the company "CSX Labs".
// The "exit" command should exit the program.
// The "get-info <stock>" command should print out a mock value for the stock.

// Example:
// > about
// Candle CLI
// CSX Labs
// > get-info AAPL
// AAPL - $100
// > exit

// Imports libraries that are initialized
import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/CSXL/Candle/api"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("--> ")
		scanner.Scan()
		args := strings.Split(scanner.Text(), " ")

		switch args[0] {
		case "about":
			fmt.Println("Candle CLI")
			fmt.Println("CSX Labs")
			fmt.Println("Made w/ <3 by @absozero and @ecsbeats")
		case "get-info":
			stock_data := api.GetStockData(args[1])
			// prints first 10 data points (for testing)
			var i int = 0
			for key, value := range stock_data {
				if i == 10 {
					break
				}
				fmt.Println(key, value)
				i++
			}

		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("Invalid command: %s\n", args[0])
		}
	}
}
