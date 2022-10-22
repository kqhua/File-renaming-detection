package main

import (
	"fmt"
	"io"
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

// Diff functions

func parseDiffDiscovery(diffs [][]byte) ([]string, []string) {
	var filtered_froms []string
	var filtered_tos []string

	for i, _ := range diffs {
		currDiff := string(diffs[i])
		currDiffList := strings.Fields(currDiff)
		from := currDiffList[2][2:]
		to := currDiffList[3][2:]

		filtered_froms = append(filtered_froms, from)
		filtered_tos = append(filtered_tos, to)
	}

	fmt.Println(filtered_froms)
	fmt.Println(filtered_tos)
	return filtered_froms, filtered_tos
}

func discoverDiff() {
	fmt.Println("Running git add")
	gitAdd()
	out := gitRenameDiff("rename-test")
	diffLines := re.MustCompile(`(diff --git a/.* b/.*)(?:\r|\n|\r\n)`)
	found := diffLines.FindAll(out, -1)

	if found != nil {
		parseDiffDiscovery(found)
	} else {
		fmt.Println("File renamed not detected.")
	}
	fmt.Println("Running git reset")
	fmt.Println()
	gitReset()
}

func findDiff() ([]string, []string) {
	fmt.Println("Running git add")
	gitAdd()
	out := gitRenameDiff("rename-test")
	diffLines := re.MustCompile(`(diff --git a/.* b/.*)(?:\r|\n|\r\n)`)
	found := diffLines.FindAll(out, -1)
	var ogNames, newNames []string
	if found != nil {
		ogNames, newNames = parseDiffDiscovery(found)
	} else {
		fmt.Println("File renamed not detected.")
	}
	fmt.Println("Running git reset")
	fmt.Println()
	gitReset()
	return ogNames, newNames
}

func findInArray(names []string, target string) int {
	for i, name := range names {
		if name == target {
			return i
		}
	}
	return -1
}

func checkRename(image string) string {
	ogNames, newNames := findDiff()
	if ogNames == nil && newNames == nil {
		return ""
	}

	namePos := findInArray(ogNames, image)
	if namePos == -1 {
		return ""
	}

	return newNames[namePos]
}

// Check if a given image is multi-arch
func checkMultiArch(image string) bool {
	//check image has been renamed
	newName := checkRename(image)
	if newName != "" {
		image = newName
	}
	// Check if image has a platforms.txt file
	file, err := os.Open(fmt.Sprintf("%s/%s", image, "platforms.txt"))
	if err != nil {
		return false
	}
	defer file.Close()

	platformFileRaw, err := io.ReadAll(file)
	if err != nil {
		return false
	}
	platformFile := string(platformFileRaw)

	// Check the file contains both AMD and ARM platforms
	if strings.Contains(platformFile, "linux/amd64") && strings.Contains(platformFile, "linux/arm64") {
		return true
	}

	return false
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

// func makeFile(fileName string) {
// 	f, err := os.Create(fileName)

// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	defer f.Close()

// 	_, err2 := f.WriteString("old falcon\n")

// 	if err2 != nil {
// 		fmt.Println(err2.Error())
// 	}

// 	fmt.Println("done")

// }
