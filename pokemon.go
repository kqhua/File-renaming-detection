package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	re "regexp"
	"strings"
)

func main() {
	//testMoveDiscovery()
	//testRenameDiscovery()
	// out := gitRenameDiff("rename-test")
	// fmt.Println(string(out))
	discoverDiff()
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

func testMoveDiscovery() {
	fmt.Println("Testing moving a file and detecting change")
	fmt.Println()
	ogName := "first_name.txt"
	moveFile(ogName, true)
	discoverMove()
	fmt.Println("Adding and making a test commit to save movement changes")
	gitAdd()
	gitCommit()
	moveFile(ogName, false)
	discoverMove()
	fmt.Println("Undoing test commit")
	gitUndoCommit()
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

// moveFile is a wrapper function to move a file in/out of the move folder
func moveFile(fileName string, in bool) {
	var first, second string
	move_path := filepath.Join("move", fileName)
	if in {
		fmt.Println("Moving " + fileName + " into move")
		first = fileName
		second = move_path
	} else {
		fmt.Println("Moving " + fileName + " out of move")
		first = move_path
		second = fileName
	}
	fmt.Println()
	err := os.Rename(first, second)
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
	} else {
		fmt.Println("File renamed not detected. Files differed too much")
	}
	fmt.Println("Running git reset")
	fmt.Println()
	gitReset()
}

func discoverMove() {
	fmt.Println("Running git add")
	gitAdd()
	out := gitStatus()
	renamedLine := re.MustCompile(`(renamed.*?)(?:\r|\n|\r\n)`)
	found := renamedLine.Find(out)
	if found != nil {
		fmt.Printf("Discovered line: %q \n", found)
		printMoveDiscovery(string(found))
	} else {
		fmt.Println("File renamed not detected.")
	}
	fmt.Println("Running git reset")
	fmt.Println()
	gitReset()

}

func discoverDiff() {
	fmt.Println("Running git add")
	gitAdd()
	out := gitRenameDiff("rename-test")
	renameFrom := re.MustCompile(`(rename from.*?)(?:\r|\n|\r\n)`)
	foundFrom := renameFrom.FindAll(out, -1)
	renameTo := re.MustCompile(`(rename to.*?)(?:\r|\n|\r\n)`)
	foundTo := renameTo.FindAll(out, -1)
	if foundFrom != nil {
		fmt.Printf("%q\n", foundFrom)
		fmt.Printf("%q\n", foundTo)
	} else {
		fmt.Println("File renamed not detected.")
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

func printMoveDiscovery(renameLine string) {
	words := strings.Fields(renameLine)
	ogFilename := words[1]
	newFilename := words[3]

	if strings.Contains(newFilename, "move/") {
		fmt.Println(ogFilename, "was moved into the move directory")
	}

	if strings.Contains(ogFilename, "move/") {
		fmt.Println(newFilename, "was moved out of the move directory")
	}
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
}

func gitReset() {
	cmd := exec.Command("git", "reset")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
}

func gitCommit() {
	cmd := exec.Command("git", "commit", "-m", "'test_commit' ")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
}

func gitUndoCommit() {
	cmd := exec.Command("git", "reset", "HEAD~")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
}

func gitRenameDiff(branch string) []byte {
	cmd := exec.Command("git", "diff", "--diff-filter=R", branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	return out
}

func hey() {
	fmt.Println("hey")
}

func heyy() {
	fmt.Println("heyy")
}
