package notes

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)
const (
	Reset 	= "\033[0m"
	Red 	= "\033[31m"
	Green 	= "\033[32m"
	GreenBg = "\033[42m"
	Yellow 	= "\033[33m"
	Black 	= "\033[30m"
)
var contentSlice []string
var filename string
var password string

func TempMain() {
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