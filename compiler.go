package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type data struct {
	Mem  int
	Body string
}

type compilerState struct {
	NumOpenLoops   int
	WritingComment bool
}

func main() {
	// Set CLI flags
	inputFileLoc := flag.String("in", "", "The location of the input file, defaults to empty string. A file is required and will throw if not found")
	numMemoryLocs := flag.Int("mem", 30000, "Number of memory locations to allocate, defaults to 30,000")
	cwd, _ := os.Getwd()
	outputLoc := flag.String("out", filepath.Join(cwd, "out.go"), "The location to output the compiled tempate, defaults to \"out.go\" in the current directory")
	flag.Parse()

	// Check for required commands
	if *inputFileLoc == "" {
		err := errors.New("expected an input file but none was provided, use the in flag, see help for more details")
		processError(err)
	}

	// Parse file contents
	bytes, err := os.ReadFile(*inputFileLoc)
	processError(err)

	// Try and parse the input file
	res, err := processFile(string(bytes))
	processError(err)

	d := data{
		Mem:  *numMemoryLocs,
		Body: res,
	}

	// Create the template and (try to) save
	template, _ := template.ParseFiles("template")
	file, err := os.Create(*outputLoc)

	processError(err)

	template.Execute(file, d)
	file.Close()
}

func isIgnoredRune(r rune, state *compilerState) bool {
	// If we are not writing a comment, ignore all white space
	if !state.WritingComment {
		return unicode.IsSpace(r)
	}

	// If we are writing a comment ignore all white space, expect the space rune specifically
	return unicode.IsSpace(r) && r != ' '
}

func processError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func processFile(input string) (string, error) {
	var sb strings.Builder

	state := &compilerState{
		NumOpenLoops:   0,
		WritingComment: false,
	}

	for _, r := range input {
		// Ignore whitespace etc.
		if isIgnoredRune(r, state) {
			continue
		}

		// Try to process the command and update the tab index as necessary
		didProcessCommand, err := tryProcessCommand(r, &sb, state)

		if err != nil {
			return "", err
		}

		// If we processed the command keep going
		if didProcessCommand {
			state.WritingComment = false
			continue
		}

		// Otherwise start (or continue) writing the comment
		state.WritingComment = writeComment(r, &sb, state)
	}

	// If we happen to be writing a comment at the end of the file make sure we force a newline to prevent catching the closing curly
	if state.WritingComment {
		sb.WriteRune('\n')
	}

	if state.NumOpenLoops > 0 {
		return "", errors.New("invalid syntax, open loop detected at end of program parsing")
	}

	return sb.String(), nil
}

func tryProcessCommand(r rune, sb *strings.Builder, state *compilerState) (bool, error) {
	var command string

	// Get the command for the given operator
	switch r {
	case '>':
		command = "index++"
	case '<':
		command = "index--"
	case '+':
		command = "mem[index]++"
	case '-':
		command = "mem[index]--"
	case '[':
		command = "for mem[index] != 0 {"
	case ']':
		command = "}"
	case ',':
		command = "mem[index] = readChar(reader)"
	case '.':
		command = "fmt.Printf(\"%c\", mem[index])"
	// If we don't match an operator it'll be a "comment" character so return false here
	default:
		return false, nil
	}

	// If we were commenting before start a new line so we don't write operations inside the comment
	if state.WritingComment {
		sb.WriteRune('\n')
	}

	// If we are closing a loop pre-decrement so we can get some tab formatting
	if r == ']' {
		state.NumOpenLoops--
		if state.NumOpenLoops < 0 {
			return false, errors.New("invalid syntax, loop close detected with no associated open loop")
		}
	}

	// Write the command
	writeCommand(sb, state.NumOpenLoops+1, command, true)

	// If we are opening a loop post-increment so we can get some tab formatting
	if r == '[' {
		state.NumOpenLoops++
	}

	return true, nil
}

func writeComment(r rune, sb *strings.Builder, state *compilerState) bool {
	// If this is the first rune in the comment add the comment characters here
	if !state.WritingComment {
		writeCommand(sb, state.NumOpenLoops+1, "// ", false)
	}

	// Write the rune
	sb.WriteRune(r)
	return true
}

func writeCommand(sb *strings.Builder, tabIndex int, command string, withNewLine bool) {
	// Make sure we ad all the tabs for the relevant indent level
	for i := 0; i < tabIndex; i++ {
		sb.WriteRune('\t')
	}

	// Write the command
	sb.WriteString(command)

	// Append a newline if required
	if withNewLine {
		sb.WriteRune('\n')
	}
}
