# Properties for Notetool

Added some "properties" for this task. Will update this if needed when we discuss this task as a group. 

## Functions

* main() - need to read arguments(Args) from commandline for text file's name.

* InitializeTxtFile() check if valid input in args. Help if needed. Create/load text file. Can return true if filename is valid or false if not.

* CLI - this can be run as it's own function or in main(). Basicly forever loop that does something when user gives input: 1: show notes, 2: add note, 3: delete note, 4: exit program. Need some error handling when user gives invalid key.

* ReadNotes() - reads text file (probably into a map or slice) and prints out the result for user. this function can be void as in doesent return anything.

* AddNote() - adds user input into text file. 

* DeleteNote() - removes note from text file based on integer that user inputs. Need to read text file into a map or slice and remove said note. Need to also raise all indexes to the removed one. For example note in 1 is removed 2 is now 1 3 is 2 etc...

check out the CLI layout from Usage.

## Markdown

We can do this as we progress. Requirements for .md below.

    Explains what the tool does
    Explains the usage of the tool, with examples
    Explains how the data is stored
