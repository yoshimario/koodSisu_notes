package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

// Define a Note structure
type Note struct {
	Text      string
	Tags      []string
	Timestamp time.Time // Change to time.Time
}

// clearScreen clears the terminal screen.
func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin": // Unix-based systems (Linux, macOS)
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println(Red + "Unsupported platform, cannot clear terminal." + Reset)
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(Red + "Please provide the filename as an argument." + Reset)
		return
	}

	// Clear the terminal once at the start of the program
	clearScreen()

	filename := os.Args[1]
	notes, err := loadNotes(filename)
	if err != nil {
		fmt.Println(Red + "Error loading notes: " + err.Error() + Reset)
		return
	}

	for {
		// Color-coded menu options (including text descriptions)
		fmt.Println("\n" + Cyan + "Select operation:" + Reset)
		fmt.Println(Green + "1. Show notes." + Reset)
		fmt.Println(Green + "2. Add a note." + Reset)
		fmt.Println(Green + "3. Search notes." + Reset)
		fmt.Println(Green + "4. Delete a note." + Reset)
		fmt.Println(Green + "5. Exit." + Reset)

		choice := readInput(Yellow + "Enter your choice: " + Reset)

		switch choice {
		case "1":
			showNotes(notes)
		case "2":
			noteText, tags := addNote()
			notes[fmt.Sprintf("%03d", len(notes)+1)] = Note{Text: noteText, Tags: tags, Timestamp: time.Now()}
			err := saveNotes(filename, notes)
			if err != nil {
				fmt.Println(Red + "Error saving note: " + err.Error() + Reset)
			} else {
				fmt.Println(Green + "Note added successfully!" + Reset)
			}
		case "3":
			keyword := readInput(Yellow + "Enter keyword to search: " + Reset)
			searchNotes(notes, keyword)
		case "4":
			deleteNote(notes, filename)
		case "5":
			fmt.Println(Green + "Exiting the program." + Reset)
			return
		default:
			fmt.Println(Red + "Invalid option. Please try again." + Reset)
		}
	}
}

// Function to read input from the user
func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// Function to show all notes
func showNotes(notes map[string]Note) {
	if len(notes) == 0 {
		fmt.Println(Red + "No notes in this collection." + Reset)
		return
	}
	fmt.Println(Cyan + "Notes:" + Reset)
	for id, note := range notes {
		fmt.Printf("%s - %s [Tags: %s] - %s\n", id, note.Text, strings.Join(note.Tags, ", "), note.Timestamp.Format(time.RFC3339))
	}
}

// Function to add a note
func addNote() (string, []string) {
	noteText := readInput(Yellow + "Enter the note text: " + Reset)
	tags := readInput(Yellow + "Enter tags for the note (comma separated): " + Reset)
	tagList := strings.Split(tags, ",")
	for i := range tagList {
		tagList[i] = strings.TrimSpace(tagList[i])
	}
	return noteText, tagList
}

// Function to load notes from a file
func loadNotes(filename string) (map[string]Note, error) {
	notes := make(map[string]Note)

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		// If the file does not exist, return an empty map without error
		return notes, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

// Function to save notes to any specified file
func saveNotes(filename string, notes map[string]Note) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal notes: %v", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("could not write to file: %v", err)
	}

	return nil
}

// Function to search notes by keyword or tag
func searchNotes(notes map[string]Note, keyword string) {
	found := false
	for id, note := range notes {
		if strings.Contains(note.Text, keyword) || contains(note.Tags, keyword) {
			fmt.Printf("%s - %s [Tags: %s] - %s\n", id, note.Text, strings.Join(note.Tags, ", "), note.Timestamp.Format(time.RFC3339))
			found = true
		}
	}

	if !found {
		fmt.Println(Red + "No notes found." + Reset)
	}
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if strings.Contains(elem, item) {
			return true
		}
	}
	return false
}

// Function to delete a note
func deleteNote(notes map[string]Note, filename string) {
	noteID := readInput(Yellow + "Enter the note ID to delete: " + Reset)
	if _, exists := notes[noteID]; !exists {
		fmt.Println(Red + "Note ID not found." + Reset)
		return
	}

	delete(notes, noteID)
	fmt.Println(Green + "Note deleted successfully!" + Reset)

	err := saveNotes(filename, notes)
	if err != nil {
		fmt.Println(Red + "Error saving notes after deletion: " + err.Error() + Reset)
	}
}