package main

import (
	"fmt"
	"os"
	"os/exec"
	re "regexp"
	"strings"
)

func main() {
	testRenameDiscovery()
}

// Test Functions

// testRenameDiscovery renames a file and attempts to discover the changes. Reverts files to og state after
func testRenameDiscovery() {
	fmt.Println("Testing renaming file and detecting change")
	fmt.Println()
	ogName := "first_name.txt"
	newName := "second_name.txt"
	renameFile(ogName, newName)
	discoverRename()
	renameFile(newName, ogName)
}

// Util Functions

// renameFile is a wrapper function to rename a file
func renameFile(fileName string, newName string) {
	fmt.Println("Renaming " + fileName + " to " + newName)
	fmt.Println()
	err := os.Rename(fileName, newName)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// discoverRename stages new changes in the repo with git add, searches for renamed line and prints it if found, then unstages changes.
func discoverRename() {
	fmt.Println("Running git add")
	gitAdd()
	out := gitStatus()
	renamedLine := re.MustCompile(`(renamed.*?)(?:\r|\n|\r\n)`)
	found := renamedLine.Find(out)
	if found != nil {

		fmt.Printf("Discovered line: %q \n", found)
		printRenameDiscovery(string(found))

	}
	fmt.Println("Running git reset")
	fmt.Println()
	gitReset()
}

func printRenameDiscovery(renameLine string) {
	words := strings.Fields(renameLine)
	ogFilename := words[1]
	newFilename := words[3]

	fmt.Println("Old name:", ogFilename)
	fmt.Println("New name:", newFilename)
}

func discoverModified() {
	out := gitStatus()
	renamedLine := re.MustCompile(`(modified.*?)(?:\r|\n|\r\n)`)
	found := renamedLine.FindAll(out, -1)
	if found != nil {
		fmt.Printf("%q\n", found)
	}
}

// Git functions
func gitStatus() []byte {
	cmd := exec.Command("git", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	return out
}

func gitAdd() {
	cmd := exec.Command("git", "add", ".")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	//fmt.Printf("combined out:\n%s\n", string(out))
}

func gitReset() {
	cmd := exec.Command("git", "reset")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	//fmt.Printf("combined out:\n%s\n", string(out))
}
