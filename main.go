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
	multi := checkMultiArch("test-multi-arch/tag")
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
	out := gitRenameDiff("main")
	diffLines := re.MustCompile(`(diff --git a/.* b/.*)(?:\r|\n|\r\n)`)
	found := diffLines.FindAll(out, -1)
	var ogNames, newNames []string
	if found != nil {
		ogNames, newNames = parseDiffDiscovery(found)
	} else {
		fmt.Println("File renamed not detected.")
	}
	return ogNames, newNames
}

func findInArray(names []string, target string) int {
	for i, name := range names {
		if target == name {
			return i
		}
	}
	return -1
}

func checkRename(image string) string {
	ogNames, newNames := findDiff()
	if ogNames == nil && newNames == nil {
		fmt.Println("No diff found")
		return filepath.Dir(image)
	}
	//use ognames if comapring branch to main, newNames if main to branch
	namePos := findInArray(ogNames, image)
	if namePos == -1 {
		fmt.Println("couldn't find ogName")
		return filepath.Dir(image)
	}

	//return the renamed directory
	return filepath.Dir(newNames[namePos])
}

// Check if a given image is multi-arch image string will be "etc/etc/platform.txt" as funcitons will detect the folder rename and txt movement
func checkMultiArch(image string) bool {
	platform_path := filepath.Join(image, "platforms.txt")

	filet, errt := os.Open(platform_path)
	if errt != nil {
		fmt.Println("Cannot find", platform_path)
	}
	defer filet.Close()

	//check image has been renamed
	dockerfile_path := filepath.Join(image, "Dockerfile")
	image = checkRename(dockerfile_path)
	fmt.Println(image)
	new_platform_path := filepath.Join(image, "platforms.txt")

	file, err := os.Open(new_platform_path)
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
