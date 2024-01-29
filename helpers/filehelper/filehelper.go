package filehelper // Check if the input is a file path

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

func IsFilePath(input string) (bool, string) {
	// Check if the file exists in the current directory
	if _, err := os.Stat(input); err == nil {
		absPath, _ := filepath.Abs(input)
		return true, absPath
	}
	return false, ""
}

// Read file contents and return as an array of strings
func ReadFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var contents []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contents = append(contents, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}

func ReadStdin() []string {
	var inputArray []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputArray = append(inputArray, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading standard input:", err)
	}
	return inputArray
}

func PrintStructFields(s interface{}) {
	v := reflect.ValueOf(s)

	// Make sure the input is a struct
	if v.Kind() != reflect.Struct {
		fmt.Println("Input is not a struct")
		return
	}

	// Iterate over fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		fmt.Printf("%s: %v\n", field.Name, value.Interface())
	}
}
