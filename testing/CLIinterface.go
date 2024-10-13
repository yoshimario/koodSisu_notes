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
				fmt.Println("\n%sExiting program.%s", Yellow, Reset)
				return
			}
			executeCommand(input)
		case keyboard.KeyEsc: //exit
			fmt.Println("\n%sExiting program.%s", Yellow, Reset)
			return
		default: 
			//handling also numeric keys.
			if char >= '1' && char <= '5' {
				currentSelection = int(char - '1') //converting char to int
			} else {
				fmt.Println("\n%sInvalid key.%s", Red, Reset)
				fmt.Println("\n%sNavigation: 1-5 or up and down arrowkeys.\nSelect command: Enter key.\nQuit progam: Escape key.%s", Yellow, Reset)
				fmt.Println("\n%sPress any key to continue...%s", Yellow, Reset)
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
			fmt.Println("\nPress any key to continue...")
			_, _, _ =keyboard.GetKey()
		case Search:
			//do something
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Println("\nPress any key to continue...")
			_, _, _ =keyboard.GetKey()
		case Add:
			//do something
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Println("\nPress any key to continue...")
			_, _, _ =keyboard.GetKey()
		case Delete:
			//do something
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Println("%s\nPress any key to continue...%s", Yellow, Reset)
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
	fmt.Println("\n%sSelect operation:%s", Yellow)

	for i, option := range options {
		if i == currentSelection {
			fmt.Printf("%s%s %s%s\n", GreenBg, Green, option, Reset)
		} else {
			fmt.Printf("%s %s%s\n", Green, option, Reset)
		}
	}

	fmt.Println("%sNavigation: 1-5 or up and down arrowkeys.\nSelect command: Enter key.\nQuit progam: Escape key.%s", Yellow, Reset)
}


