package notes

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var contentSlice []string
var filename string
var password string

func TempMain() {
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