package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func readStream(r io.Reader) ([]byte, error) {
	// This function reads all bytes from the provided io.Reader and returns them as a byte slice.
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r) // Copy data from r to buf
	if err != nil {
		return nil, err // If there's an error during copying, return it
	}
	return buf.Bytes(), nil // Return the read bytes
}

func byteCount(fileName string) (int64, error) {
	// This function returns the size of the specified file in bytes.
	info, err := os.Stat(fileName) // Get file information, like size, modification time, etc.
	if err != nil {
		return 0, err // If the file doesn't exist or can't be accessed, return the error
	}
	return info.Size(), nil // Return the size of the file
}

func lineCount(text string) int {
	// This function counts the number of lines in a given text.
	lines := strings.Split(text, "\n") // Split the text into lines based on newline characters

	// If the text ends with a newline, the last element will be an empty string.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		return len(lines) - 1 // Exclude the empty line caused by the trailing newline
	}

	return len(lines) // Otherwise, return the total count of lines
}

func wordCount(text string) int {
	// This function counts the number of words in the given text.
	return len(strings.Fields(text)) // Split by whitespace and count the resulting fields
}

func charCount(text string) int {
	// This function counts the number of characters in the given text.
	return len(text) // Simply return the length of the text
}

func sswc(args []string, stream io.Reader) (string, error) {
	// Create a new flag set to parse command-line arguments for the sswc function
	flagSet := flag.NewFlagSet("cwcc", flag.ContinueOnError)

	// Define boolean flags for counting bytes, lines, words, and characters
	countByteFlag := flagSet.Bool("c", false, "Display byte count")
	countLineFlag := flagSet.Bool("l", false, "Display line count")
	countWordFlag := flagSet.Bool("w", false, "Display word count")
	countCharFlag := flagSet.Bool("m", false, "Display character count")

	// Parse the command-line arguments to see which flags are set
	if err := flagSet.Parse(args); err != nil {
		return "", err // Return an error if argument parsing fails
	}

	// Get the non-flag arguments (like file names)
	fileArgs := flagSet.Args()

	// Default to the first argument as the file name
	var fileName string
	if len(fileArgs) > 0 {
		fileName = fileArgs[0]
	}

	// Initialize the data slice to hold the file's contents
	var data []byte
	var err error

	if fileName != "" { // If a file name is provided
		data, err = os.ReadFile(fileName) // Read the file
		if err != nil {
			return "", err // Return an error if the file can't be read
		}
	} else if stream != nil { // If no file is provided, use the stream (like standard input)
		data, err = readStream(stream)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("no input provided") // No file or stream to read from
	}

	contents := string(data) // Convert the byte slice to a string for further processing

	var result []string // A slice to accumulate the results

	// Add the appropriate counts to the result based on the flags set
	if *countByteFlag {
		byteCountValue, byteCountErr := byteCount(fileName) // Get the byte count of the file
		if byteCountErr != nil {
			return "", byteCountErr // Return an error if byte count fails
		}
		result = append(result, fmt.Sprintf("%d", byteCountValue)) // Add byte count to result
	}

	if *countCharFlag {
		result = append(result, fmt.Sprintf("%d", charCount(contents))) // Add character count
	}

	if *countLineFlag {
		result = append(result, fmt.Sprintf("%d", lineCount(contents))) // Add line count
	}

	if *countWordFlag {
		result = append(result, fmt.Sprintf("%d", wordCount(contents))) // Add word count
	}

	if len(result) == 0 { // If no specific flags were set, default to line, word, and byte count
		result = append(result, fmt.Sprintf("%d", lineCount(contents)))
		result = append(result, fmt.Sprintf("%d", wordCount(contents)))
		byteCountValue, byteCountErr := byteCount(fileName) // Get byte count
		if byteCountErr != nil {
			return "", byteCountErr
		}
		result = append(result, fmt.Sprintf("%d", byteCountValue))
	}

	return strings.Join(result, " "), nil // Join the results with a space and return them
}

func main() {
	input := os.Stdin                // Standard input, used if no file is provided
	args := os.Args[1:]              // Command-line arguments, excluding the program name
	result, err := sswc(args, input) // Call the sswc function with the given arguments
	if err != nil {
		panic(err) // Handle errors by panicking (in a real program, this might be handled differently)
	}
	fmt.Println(result) // Print the result to the console
}
