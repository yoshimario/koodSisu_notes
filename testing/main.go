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
	"time"
)

//Defining menu options as constants to prevent accidental modifications.
const (
	Show 	= "1. Show notes"
	Search 	= "2. Search notes."
	Add 	= "3. Add a note."
	Delete 	= "4. Delete a note"
	Exit 	= "5. Exit."
	//some colours for interface.
	Reset 	= "\033[0m"
	Red 	= "\033[31m"
	Green 	= "\033[32m"
	GreenBg = "\033[42m"
	Yellow 	= "\033[33m"
	Black 	= "\033[30m"
)
var options = []string {
	Show,
	Search,
	Add,
	Delete,
	Exit,
}
var contentSlice []string
var filename string
var password string


func main() {
	//checking if user provided filename.
	if len(os.Args) < 2 {
		fmt.Printf("%sNo file name provided. %sUsage: %sgo run main.go <filename> [password]%s\n", Red, Yellow, Green, Reset)
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
		fmt.Printf("%sError checking file: %s%s\n", Red, err, Reset)
	}
	//greeting the user and pausing the program so the user can read the welcome.
	fmt.Printf("%sWelcome to the notes tool!%s\n", Yellow, Reset)
	time.Sleep(2 * time.Second)
	//saving content map to file on exit
	defer saveToFile()
	CLIinterface()
}

func createFile(filename string) {
	//Create empty file.
	err := os.WriteFile(filename, []byte{}, 0644)
	if err != nil {
		fmt.Printf("%sError creating file: %s\n", Red, err, Reset)
		return
	}

	fmt.Printf("%sFile created successfully.%s\n", Yellow, Reset)
}

func readFile(filename, password string) {
	//Reading the content from file.
	encodedContent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%sError reading file: %s%s\n", Red, err, Reset) 
		return
	}

	var content string
	if password != "" {
		//decoding content. 
		decodedContent, err := base64.StdEncoding.DecodeString(string(encodedContent))
		if err != nil {
			fmt.Printf("%sError decoding content: %s%s\n", Red, err, Reset)
			return
		}

		//decrypt the content
		content = xorEncryptDecrypt(string(decodedContent), password)
	} else {
		//reading as plain text.
		content = string(encodedContent)
	}

	//storing contents to global slice.
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
		fmt.Printf("%sError saving file: %s\n", Red, err, Reset)
		return
	}

	fmt.Printf("%s\nData saved successfully.%s\n", Yellow, Reset)
}
/*XOR encryption for the data. Password works as key.
 XOR works as an and/or encryption method in binary level. For example:
 we have char 'a' that we want to encrypt with key 'b'
 'a' ascii value is 97 and in binary it's 01100001
 'b' ascii value is 98 and in binary it's 01100010
 XOR will now compare the binary values bit by bit if they are equal or not.
 if bits are equal resulting bit is 0 and if not then 1
 so to encrypt a with b results: 00000011
*/
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
		fmt.Printf("%sError opening keyboard: %s%s\n", Red, err, Reset) 
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
			fmt.Printf("%sError reading key: %s%s\n", Red, err, Reset)
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
				clearTerminal()
				fmt.Printf("\n%sExiting program.%s\n", Yellow, Reset)
				return
			}
			executeCommand(input)
		case keyboard.KeyEsc: //exit
			fmt.Printf("\n%sExiting program.%s\n", Yellow, Reset)
			return
		default: 
			//handling also numeric keys.
			if char >= '1' && char <= '5' {
				currentSelection = int(char - '1') //converting char to int
			} else {
				fmt.Printf("\n%sInvalid key.%s\n", Red, Reset)
				fmt.Printf("\n%sNavigation: 1-5 or up and down arrowkeys.\nSelect command: Enter key.\nQuit progam: Escape key.%s\n", Yellow, Reset)
				fmt.Printf("%sPress any key to continue...%s\n", Yellow, Reset)
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
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Printf("%sPress any key to continue...%s\n", Yellow, Reset)
			_, _, _ =keyboard.GetKey()
		case Search:
			//do something
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Printf("%sPress any key to continue...%s\n", Yellow, Reset)
			_, _, _ =keyboard.GetKey()
		case Add:
			//do something
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Printf("%sPress any key to continue...%s\n", Yellow, Reset)
			_, _, _ =keyboard.GetKey()
		case Delete:
			//do something
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Printf("%sPress any key to continue...%s\n", Yellow, Reset)
			_, _, _ =keyboard.GetKey()
		
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
	fmt.Printf("\n%sSelect operation:%s\n", Yellow, Reset)

	for i, option := range options {
		if i == currentSelection {
			fmt.Printf("%s%s %s%s\n", GreenBg, Green, option, Reset)
		} else {
			fmt.Printf("%s %s%s\n", Green, option, Reset)
		}
	}

	fmt.Printf("%s\nNavigation: 1-5 or up and down arrowkeys.\nSelect command: Enter key.\nQuit progam: Escape key.%s\n", Yellow, Reset)
}


