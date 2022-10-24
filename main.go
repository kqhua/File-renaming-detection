package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	re "regexp"
	"strings"
)

func main() {
	multi := checkMultiArch("platforms.txt")
	fmt.Println(multi)
}

// Diff functions

func parseDiffDiscovery(diffs [][]byte) ([]string, []string) {
	var filtered_froms []string
	var filtered_tos []string

	for i, _ := range diffs {
		currDiff := string(diffs[i])
		fmt.Println(currDiff)
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

func findDiff() ([]string, []string) {
	fmt.Println("Running git add")
	gitAdd()
	out := gitRenameDiff("main")
	//fmt.Println(string(out))
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
		//fmt.Println("looking for", target, "in", name)
		if strings.Contains(name, target) {
			return i
		}
	}
	return -1
}

func checkRename(image string) string {
	ogNames, newNames := findDiff()
	if ogNames == nil && newNames == nil {
		fmt.Println("No diff found")
		return ""
	}

	//use ognames if comapring branch to main, newNames if main to branch
	namePos := findInArray(newNames, image)
	if namePos == -1 {
		fmt.Println("couldn't find ogName")
		return ""
	}

	return newNames[namePos]
}

// Check if a given image is multi-arch image string will be "platform.txt" as funcitons will detect the folder rename and txt movement
func checkMultiArch(image string) bool {
	//check image has been renamed
	newName := checkRename(image)
	fmt.Println(newName)
	if newName != "" {
		image = newName
	}
	file, err := os.Open(image)
	if err != nil {
		return false
	}
	defer file.Close()

	platformFileRaw, err := io.ReadAll(file)
	if err != nil {
		return false
	}
	platformFile := string(platformFileRaw)
	//fmt.Println(platformFile)
	// Check the file contains both AMD and ARM platforms
	if strings.Contains(platformFile, "linux/amd64") && strings.Contains(platformFile, "linux/arm64") {
		return true
	}

	return false
}

// Git functions

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

func gitRenameDiff(branch string) []byte {
	cmd := exec.Command("git", "diff", "--diff-filter=R", branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	return out
}
