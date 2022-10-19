package main

import (
	"fmt"
	"os"
	"os/exec"
	re "regexp"
)

func main() {
	fmt.Println("Testing renaming file and detecting change")

	filename := "data.txt"
	// makeFile(filename)
	newName := "pokemon.txt"
	// renameFile(filename, newName)

	renameFile(filename, newName)

	//gitStatus()

}

func makeFile(fileName string) {
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer f.Close()

	_, err2 := f.WriteString("old falcon\n")

	if err2 != nil {
		fmt.Println(err2.Error())
	}

	fmt.Println("done")
}

func renameFile(fileName string, newName string) {
	fmt.Println("renaming " + fileName + " to " + newName)
	err := os.Rename(fileName, newName)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func lsa() {
	cmd := exec.Command("ls", "-lah")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	if ok, _ := re.Match("pokemon.txt", out); ok {
		err = nil
		fmt.Println("found pokemon")
	}
}

func gitStatus() {
	cmd := exec.Command("git", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	if ok, _ := re.Match("pokemon.txt", out); ok {
		err = nil
		fmt.Println("found pokemon")
	}
}
