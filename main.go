package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/eiannone/keyboard"
	"runtime"
	"time"
	"strconv"
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
	Cyan 	= "\033[36m"
)
var options = []string {
	Show,
	Search,
	Add,
	Delete,
	Exit,
}
type Note struct {
	Text	string
	Tags 	[]string
	Time 	time.Time
}
var contentSlice []Note
var filename string
var password string
//struct for note data.


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
		fmt.Printf("%sError creating file: %s%s\n", Red, err, Reset)
		return
	}

	fmt.Printf("%sFile created successfully.%s\n", Yellow, Reset)
}

func readFile(filename, password string) {
	//Reading the content from file.
	encodedContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("%sError reading file: %s%s\n", Red, err, Reset) 
		return
	}

	var content string
	if password != "" {
		//decrypting content
		decryptedContent := xorEncryptDecrypt(string(encodedContent), password)
		//decoding content. 
		decodedContent, err := base64.StdEncoding.DecodeString(decryptedContent)
		if err != nil {
			fmt.Printf("%sError decoding content: %s%s\n", Red, err, Reset)
			return
		}

		//decrypt the content
		//content = xorEncryptDecrypt(string(decodedContent), password)
		content = string(decodedContent)
	} else {
		//reading as plain text.
		content = string(encodedContent)
	}
	// splitting contents into individual notes.
	lines := strings.Split(content, "\n")	

	for _, line := range lines {
		//skipping empty lines.
		if line == "" {
			continue
		}

		// Extracting time, text, and tags from line.
		//format should be "DD/MM/YY Hour:Minute - text [tag1, tag2]"

		parts := strings.SplitN(line, " - ", 2)
		if len(parts) < 2 {
			continue //invalid format
		}

		//parsing time, text and tags.
		timeStr := parts[0]
		textAndTags := parts[1]

		//parse tags
		var tags []string
		if idx := strings.Index(textAndTags, "["); idx != -1 {
			if endIdx := strings.Index(textAndTags, "]"); endIdx != -1 {
				tagsStr := textAndTags[idx+1 : endIdx] // Get content inside brackets
				tags = strings.Split(tagsStr, ", ")
				textAndTags = textAndTags[:idx] // Remove tags from text.
			}
		}

		//creating time object
		timeParsed, err := time.Parse("02/01/2006 15:04", timeStr)
		if err != nil {
			fmt.Printf("%sError parsing time: %s%s\n", Red, err, Reset)
			continue //skipping this note.
		}

		//creating note and appending it to a slice.
		note := Note{
			Text: textAndTags,
			Tags: tags,
			Time: timeParsed,
		}
		contentSlice = append(contentSlice, note)
	}
	
}
/*EncodedContent := base64.StdEncoding.EncodeToString([]byte(output))
			encryptedContent := xorEncryptDecrypt(EncodedContent, password)
			output = encryptedContent*/
//saving the contentSlice into a text file.
func saveToFile() {
	var output strings.Builder

	for _, note := range contentSlice {
		//formatting time
		timeStr := note.Time.Format("02/01/2006 15:04")
		//joining tags with commas
		tagsStr := strings.Join(note.Tags, ", ")
		//creating line in desired format.
		if tagsStr != "" {
			output.WriteString(fmt.Sprintf("%s - %s [%s]\n",timeStr, note.Text, tagsStr))
		} else {
			output.WriteString(fmt.Sprintf("%s - %s\n", timeStr, note.Text)) //incase there are no tags.
		}
		
	
	}
	var outputString string
	//encrypting content if there is a password.
	if password != "" {
		encodedContent := base64.StdEncoding.EncodeToString([]byte(output.String()))
		encryptedContent := xorEncryptDecrypt(encodedContent, password)
		outputString = encryptedContent
	} else {
		outputString = output.String()
	}

	//writing the content in to a file.
	err := os.WriteFile(filename, []byte(outputString), 0644)
	if err != nil {
		fmt.Printf("%sError saving file: %s%s\n", Red, err, Reset)
		return
	}

	fmt.Printf("\n%s\nData saved successfully.%s\n", Yellow, Reset)
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
				clearTerminal()
				fmt.Printf("\n%sInvalid key.%s\n", Red, Reset)
				fmt.Printf("\n%sNavigation: 1-5 or up and down arrowkeys.\nSelect command: Enter key.\nQuit progam: Escape key.%s\n", Yellow, Reset)
				fmt.Printf("\n%sPress any key to continue...%s\n", Yellow, Reset)
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
			//Showing all the notes here:
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			if len(contentSlice) == 0 {
				fmt.Println("No notes available.")
			} else {
				for i := 0; i < len(contentSlice); i++ {
					PrintNote(i) //Use the PrintNote function to print each note.
				}
			}
			fmt.Printf("\n%sPress any key to continue...%s\n", Yellow, Reset)
    		_, _, _ = keyboard.GetKey()
		case Search:
			//do something
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			promptMessage := fmt.Sprintf("%sSearch notes with index or tags (comma separated): %s", Yellow, Reset)
			fmt.Print(promptMessage)

			var input string
			var tags []string
			for {
				char, key, err := keyboard.GetKey()
				if err != nil {
					fmt.Printf("%sError reading key: %s%s\n", Red, err, Reset)
					
					continue
				}

				if key == keyboard.KeyEnter {
					if input == "" {
						fmt.Printf("%sNo input provided. Please enter an index or tags.%s\n", Red, Reset)
						fmt.Print(promptMessage)
						continue
					}

					//checking if input is index
					if index, err := strconv.Atoi(input); err == nil && index > 0 {
						fmt.Println()
						ReadNotes(index, nil) //search by index
					} else {
						//split input into tags and trim whitespace
						tags = strings.Split(input, ",")
						for i := range tags {
							tags[i] = strings.TrimSpace(tags[i])
						}
						fmt.Println()
						ReadNotes(0, tags) //search by tags
					}
					break
				} else if key == keyboard.KeyBackspace2 {
					if len(input) > 0 {
						input = input[:len(input)-1] //remove last character
						fmt.Print("\r" + strings.Repeat(" ", 50) + "\r") //clear the line
						// Reprint prompt and current input
						fmt.Print(fmt.Sprintf("%s%s%s", promptMessage, Reset, input))
					}
				} else {
					input += string(char)
					fmt.Print(string(char))
				}
			}

			fmt.Printf("\n%sPress any key to continue...%s\n", Yellow, Reset)
			_, _, _ =keyboard.GetKey()
		case Add:
			//adding a note here.
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Printf("%sEnter the note: %s", Yellow, Reset)

			var noteText string
   		 	for {
        		char, key, err := keyboard.GetKey()
        		if err != nil {
            		fmt.Printf("Error reading key: %s\n", err)
            		continue
        		}

        		if key == keyboard.KeyEnter {
            		break // Exit loop on Enter
        		} else if key == keyboard.KeyBackspace2 { // Handle backspace
					
            		if len(noteText) > 0 {
                		noteText = noteText[:len(noteText)-1] // Remove last character
                		fmt.Print("\r" + strings.Repeat(" ", 50) + "\r") // Clear the line
                		fmt.Printf("%sEnter note text: %s" + noteText, Yellow, Reset) // Reprint prompt and current text
            		}
        		} else if key == keyboard.KeySpace { // Handle space
            		noteText += " " // Append space to noteText
            		fmt.Print(" ")   // Echo the space
        		} else {
            		noteText += string(char) // Append char to noteText
            		fmt.Print(string(char))   // Echo the character
        		}
    		}

    		// Clear the line after input
    		fmt.Print("\r\033[K") // Clears the current line

    		// Now we ask for tags
    		fmt.Printf("%sEnter tags (comma separated) if needed: %s", Yellow, Reset)
    		var tagsText string 
    		for {
        		char, key, err := keyboard.GetKey()
        		if err != nil {
            		fmt.Printf("%sError reading key: %s%s\n", Red, err, Reset)
            		continue
        		}
        
        		if key == keyboard.KeyEnter {
            		break // Exiting loop on Enter
        		} else if key == keyboard.KeyBackspace2 {
            		if len(tagsText) > 0 {
                		tagsText = tagsText[:len(tagsText)-1] // Removes last character
                		fmt.Print("\r" + strings.Repeat(" ", 50) + "\r") // Clear the line
                		fmt.Print("Enter tags (comma separated) if needed: " + tagsText) // Reprint prompt and current text
            		}
        		} else if key == keyboard.KeySpace { // Handle space
            		tagsText += " " // Append space to tagsText
            		fmt.Print(" ")   // Echo the space
        		} else {
            		tagsText += string(char) // Appending char to tagsText
            		fmt.Print(string(char))   // Echoing the character
        		}
    		}
    		// Clear the line after input
    		fmt.Print("\r\033[K") // Clear the current line

    		// Splitting tags by comma.
    		tags := strings.Split(tagsText, ",")
    		for i := range tags {
        		tags[i] = strings.TrimSpace(tags[i]) // Trimming whitespace from each tag.
    		}
    		// Add note with tags
    		
			if len(noteText) > 0 {
				AddNote(noteText, tags...)
			} else {
				fmt.Printf("%sNote cannot be empty.%s", Red, Reset)
			}
			fmt.Printf("\n%sPress any key to continue...%s\n", Yellow, Reset)
			_, _, _ =keyboard.GetKey()

		case Delete:
			//do something
			clearTerminal()
			fmt.Printf("\n%sselected: %s%s\n", Yellow, command, Reset)
			fmt.Printf("\n%sEnter the number of note to remove or 0 to cancel: %s", Yellow, Reset)
			var indexInput string
			for {
				char, key, err := keyboard.GetKey()
				if err != nil {
					fmt.Printf("%sError reading key: %s%s\n", Red, err, Reset)
					continue
				}

				if key == keyboard.KeyEnter {
					if indexInput == "0" {
						fmt.Printf("\n%sOperation cancelled.%s", Yellow, Reset)
						break //exiting if cancelled
					}

					//trying to convert input to an index
					index, err := strconv.Atoi(indexInput)
					if err != nil || index <= 0 {
						fmt.Printf("%sInvalid input. Please enter a valid index.%s\n", Red, Reset)
						indexInput = "" //reset input
						continue
					}

					//deleting the note
					err = DeleteNote(index)
					if err != nil {
						fmt.Printf("%s%s%s", Red, err.Error(), Reset)
					}
					break
				} else if key == keyboard.KeyBackspace2 {
					if len(indexInput) > 0 {
						indexInput = indexInput[:len(indexInput)-1] //remove last character
						fmt.Print("\r" + strings.Repeat(" ", 50) + "\r") //clear the line
						fmt.Print(fmt.Sprintf("%sEnter the index of the note to delete (0 to cancel): %s%s", Yellow, Reset, indexInput)) // Reprint prompt and current input
            
					} 
				} else  if char >= '0' && char <= '9' {//only appending numbers.
					indexInput += string(char) //appending char to indexInput
					fmt.Print(string(char)) //echoing character
				}
			}
			fmt.Printf("\n%sPress any key to continue...%s\n", Yellow, Reset)
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
//adding a new note to the contentSlice
func AddNote(text string, tags ...string) {
	loc, err := time.LoadLocation("Europe/Helsinki") //adding local timezone
	if err != nil {
		fmt.Printf("%sError loading location: %s%s\n", Red, err, Reset)
		return
	}

	//get current time to variable.
	currentTime := time.Now().In(loc)
	note := Note{
		Text: text,
		Tags: tags,
		Time: currentTime,
	}
	contentSlice = append(contentSlice, note)
}
//reading notes 
func ReadNotes(index int, tags []string) {
	//if index is given checking if valid input. 
	if index > 0 {
		if index-1 >= len(contentSlice) {
			fmt.Printf("%sError: index out of range%s\n", Red, Reset)
			return
		}
		PrintNote(index - 1)
		return
	}
	//printing notes if there are tags given by the user
	if len(tags) > 0 {
		found := false //flagging found notes.
		for i, note := range contentSlice {
			for _, tag := range tags {
				if contains(note.Tags, tag) {
					PrintNote(i)
					found = true
					break;
				}
			}
		}
		if !found {
			fmt.Printf("%sNo notes found with provided tags.%s\n", Red, Reset)
		}
		return
	}
}

func PrintNote(index int) {
	note, err := GetNote(index) //using GetNote function to retrieve the note.
	if err != nil {
		fmt.Printf("%sError retrieving the note: %s%s\n", Red, err, Reset)
	}
	//formatting the note in a more readable form.
	indexStr := fmt.Sprintf("%03d", index+1)

	//formatting the time
	timeStr := note.Time.Format("02/01/2006 15:04")

	//Printing the note
	fmt.Printf("%s%s: %s - %s%s\n", Cyan, indexStr, Yellow+timeStr, Green+note.Text, Reset)
}
// helper to print note by reference
func PrintNoteByIndex(index int) {
	if index < 0 || index >= len(contentSlice) {
		fmt.Printf("%sError: index out of range%s\n", Red, Reset)
	}

	PrintNote(index)
}

//contains function for tag checking.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
//function for deleting note by index
func DeleteNote(index int) error {
	//adjusting index
	index--
	// error handling for invalid index
	if index < 0 || index >= len(contentSlice) {
		return fmt.Errorf("%s\nIndex out of range%s", Red, Reset)
	}
	// Removing note by slicing the contentSlice.
	contentSlice = append(contentSlice[:index], contentSlice[index+1:]...)
	return nil
}

//helper function for getting a singular note from slice based on index. error handling incase of panic!
func GetNote(index int) (*Note, error) {
	//error handling incase of invalid index.
	if index < 0 || index >= len(contentSlice) {
		return nil, fmt.Errorf("index out of range")
	}
	return &contentSlice[index], nil
}
