package notes

import (
	"fmt"
	"os"
	"github.com/eiannone/keyboard"
	"os/exec"
	"runtime"
)

//Defining menu options as constants to prevent accidental modifications.

const (
	Option1 = "1. Show notes"
	Option2 = "2. Search notes."
	Option3 = "3. Add a note."
	Option4 = "4. Delete a note"
	Exit 	= "5. Exit."
)

var options = []string {
	Option1,
	Option2,
	Option3,
	Option4,
	Exit,
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
	if command == Exit { //handling the exit command.
		fmt.Println("Exiting program.")
		return
	}
	fmt.Printf("\nselected: %s\n", command)
	fmt.Println("Press any key to continue...")
	_, _, _ =keyboard.GetKey()
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


