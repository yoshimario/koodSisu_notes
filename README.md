# Notes Tool (CLI)

The **Notes Tool** is a command-line interface (CLI) application that allows users to manage notes in a simple text file. It offers features like showing, searching, adding, and deleting notes, with optional password protection using XOR encryption for securing your notes.

---

## Features

- **Show Notes**: View the list of saved notes.
- **Search Notes**: Search for specific notes within the collection (future implementation).
- **Add a Note**: Add a new note to the list.
- **Delete a Note**: Remove a note from the list.
- **Encryption**: Optional password-based encryption for securing notes.
- **Intuitive Navigation**: Use arrow keys or number keys to navigate the menu.

---
## Requirements

- **Go Programming Language**: You need to have [Go](https://golang.org/doc/install) installed on your machine to run this program.

---

## Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/notes-tool.git
    ```

2. **Navigate to the project directory**:
    ```bash
    cd notes-tool
    ```

3. **Build the Go program**:
    ```bash
    go build main.go
    ```

---

Usage
After building the program, you can run it using the following command in your terminal:

bash
./main <filename> [password]
Parameters:
<filename>: The name of the text file where notes are saved. If the file doesn't exist, it will be created.
[password] (optional): An optional password for encrypting/decrypting the file content using XOR encryption.
Example Commands:
Run the program without encryption:

bash
./main notes.txt
This will create a file called notes.txt (if it doesn't exist) and allow you to manage notes in plain text.
Run the program with encryption:

bash
./main secrets.txt mypassword
This will create a file called secrets.txt (if it doesn't exist) and encrypt/decrypt the notes with the password mypassword.
Navigating the Menu
Use the Up and Down arrow keys to navigate between menu options.
Press the Enter key to select the highlighted option.
You can also press number keys (1-5) to select the respective option.
Press Esc to exit the program at any time.
Menu Options:
Show Notes:

Displays the current list of notes stored in the file.
Search Notes (not yet implemented):

Will search for specific notes based on a keyword (future implementation).
Add a Note:

Prompts you to enter a new note and adds it to the file.
Delete a Note:

Prompts you to select and remove a note from the file.
Exit:

Exits the program and automatically saves changes to the file.
Example:
1. Adding a New Note
When you run the tool and choose Add a Note from the menu:

bash
3. Add a note.
You will be prompted to input the text for the note. For example:

arduino
Enter note text: "Buy groceries"
The note will be added to the collection and saved to the file.

2. Viewing Notes
When selecting Show Notes:

bash
1. Show notes
You will see a list of notes in the file:

csharp
001 - Buy groceries
002 - Complete homework
003 - Meeting with Bob
3. Deleting a Note
When selecting Delete a Note, you will be prompted to choose a note by number:

arduino
Enter note number to delete: 2
The second note will be removed from the list.

Encryption
The program uses XOR encryption if a password is provided. This simple encryption method protects the content of your notes by encoding and decoding them using the provided password. The encrypted data is stored as Base64-encoded strings in the file.

Important:
Use the same password each time to open an encrypted file.
If you forget your password, you won't be able to decrypt your notes.
Known Issues / Features to Implement
Search function is not yet implemented.
Improvement: Add more advanced encryption options or stronger algorithms.
UX Improvements: Colors and more user-friendly interface.
License
This project is open-source and available under the MIT License.