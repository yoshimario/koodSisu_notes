package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/eiannone/keyboard"
	"runtime"
)

//Defining menu options as constants to prevent accidental modifications.
const (
	Show 	= "1. Show notes"
	Search 	= "2. Search notes."
	Add 	= "3. Add a note."
	Delete 	= "4. Delete a note"
	Exit 	= "5. Exit."
)
var contentSlice []string
var filename string
var password string

var options = []string {
	Show,
	Search,
	Add,
	Delete,
	Exit,
}



func main() {
	//checking if user provided filename.
	if len(os.Args) < 2 {
		fmt.Println("No file name provided. Usage: go run main.go <filename> [password]")
		return
	}
	
	filename = os.Args[1]
	//checking for password.
	if len(os.Args) > 2 {
		password = os.Args[2]
	}
	//checking if filename ends in .txt if not adding it.
	if !strings.HasSuffix(filename, ".txt") {
		filename += ".txt"
	}
	//checking if file exists.
	if _, err := os.Stat(filename); err == nil {
		//exists reading it.
		readFile(filename, password)
	} else if os.IsNotExist(err) {
		//doesen't exist. creating it as empty file.
		createFile(filename)
	} else {
		fmt.Println("Error checking file: ", err)
	}
	CLIinterface()
	//saving content map to file on exit
	defer saveToFile()
}

func createFile(filename string) {
	//Create empty file.
	err := os.WriteFile(filename, []byte{}, 0644)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}

	fmt.Println("File created successfully.")
}

func readFile(filename, password string) {
	//Reading the content from file.
	encodedContent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file: ", err) 
		return
	}

	var content string
	if password != "" {
		//decoding content. 
		decodedContent, err := base64.StdEncoding.DecodeString(string(encodedContent))
		if err != nil {
			fmt.Println("Error decoding content: ", err)
			return
		}

		//decrypt the content
		content = xorEncryptDecrypt(string(decodedContent), password)
	} else {
		//reading as plain text.
		content = string(encodedContent)
	}

	//storing contents to global map.
	contentSlice = strings.Split(content, "\n")	
	
}

func saveToFile() {
	output := strings.Join(contentSlice, "\n")
	if len(output) > 0 {
		//encrypting content if there is password.
		if password != "" {
			encryptedContent := xorEncryptDecrypt(output, password)
			output = base64.StdEncoding.EncodeToString([]byte(encryptedContent))
		}
	}

	//writing the content in file.
	err := os.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error saving file: ", err)
		return
	}

	fmt.Println("Data saved to file on exit.")
}

func xorEncryptDecrypt(input, key string) string {
	output := make([]byte, len(input))

	for i := 0; i < len(input); i++ {
		output[i] = input[i] ^ key[i%len(key)]
	}

	return string(output)
}
//this function drives the program. Here we handle user inputs and navigate the menu.
func CLIinterface() {
	currentSelection := 0
	// opening keyboard import for inputs.
	if err := keyboard.Open(); err != nil {
		//some error handling incase keyboard fails to open.
		fmt.Println("Error opening keyboard:", err) 
	}
	/*making sure keyboard input closes when we are done running CLIinterface.
	this wont be executed until we are leaving CLIinterface()*/
	defer keyboard.Close() 

	//loop to display menu and handle inputs.
	for {
		clearTerminal()
		displayMenu(options, currentSelection)

		char, key, err := keyboard.GetKey()
		//some more error handling for keyboard inputs.
		if err != nil {
			fmt.Println("Error reading key: ", err)
		}

		//handling inputs in switch statement
		switch key {
		case keyboard.KeyArrowUp:
			currentSelection--
			if currentSelection < 0 { //returning to bottom
				currentSelection = len(options) - 1
			}
		case keyboard.KeyArrowDown:
			currentSelection++
			if currentSelection >= len(options) { //returning to top
				currentSelection = 0
			}
		case keyboard.KeyEnter: //selection
			input := options[currentSelection] //getting the user input
			//checking if exit selected.
			if input == Exit {
				fmt.Println("\nExiting program.")
				return
			}
			executeCommand(input)
		case keyboard.KeyEsc: //exit
			fmt.Println("\nExiting program.")
			return
		default: 
			//handling also numeric keys.
			if char >= '1' && char <= '5' {
				currentSelection = int(char - '1') //converting char to int
			} else {
				fmt.Println("\nInvalid key.")
				fmt.Println("\nNavigation: 1-5 or up and down arrowkeys.\nSelect command: Enter key.\nQuit progam: Escape key.")
				fmt.Println("Press any key to continue...")
				keyboard.GetKey()
			}
		}
	}
}
//here we can handle the user input. I.E what command to run.
func executeCommand(command string) {
	//handling the commands in switch statement.
	switch command {
		case Show:
			//do something
			fmt.Printf("\nselected: %s\n", command)
			fmt.Println("Press any key to continue...")
			_, _, _ =keyboard.GetKey()
		case Search:
			//do something
			fmt.Printf("\nselected: %s\n", command)
			fmt.Println("Press any key to continue...")
			_, _, _ =keyboard.GetKey()
		case Add:
			//do something
			fmt.Printf("\nselected: %s\n", command)
			fmt.Println("Press any key to continue...")
			_, _, _ =keyboard.GetKey()
		case Delete:
			//do something
			fmt.Printf("\nselected: %s\n", command)
			fmt.Println("Press any key to continue...")
			_, _, _ =keyboard.GetKey()
		case Exit: {
			fmt.Println("Saving file, exiting program...")
			return
		}
	}
	
}

//function for clearing the terminal.
func clearTerminal() {
	var cmd *exec.Cmd
	
	//clearing terminal based on OS
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cis")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

//function to display menu in the interface.
func displayMenu(options[]string, currentSelection int) {
	fmt.Println("\nSelect operation:")

	for i, option := range options {
		if i == currentSelection {
			fmt.Printf("\033[1;30;47m %s\033[0m\n", option)
		} else {
			fmt.Printf(" %s\n", option)
		}
	}

	fmt.Println("Navigation: 1-5 or up and down arrowkeys.\nSelect command: Enter key.\nQuit progam: Escape key.")
}

