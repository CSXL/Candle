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
	"github.com/charmbracelet/lipgloss"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		args := strings.Split(input, " ")

		switch args[0] {
		case "about":
			fmt.Println("Candle CLI")
			fmt.Println("CSX Labs")
		case "get-info":
			api.GetStockData()
		case "stfu":
			var style = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				PaddingTop(2).
				PaddingLeft(4).
				Width(22)

			fmt.Println(style.Render("Hello, kitty."))
		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("Invalid command: %s\n", args[0])
		}
	}
}
