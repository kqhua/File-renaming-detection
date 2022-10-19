package main

import (
	"fmt"
	"os"
	"os/exec"
	re "regexp"
)

func main() {
	fmt.Println("Testing renaming file and detecting change")

	//filename := "data.txt"
	// makeFile(filename)
	//newName := "pokemon.txt"
	// renameFile(filename, newName)

	//renameFile(filename, newName)

	discoverModified()

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

func gitStatus() []byte {
	cmd := exec.Command("git", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	//fmt.Printf("combined out:\n%s\n", outString)

	// if ok, _ := re.Match("renamed", out); ok {
	// 	err = nil
	// 	fmt.Println("found")

	// 	// f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// 	// if err != nil {
	// 	// 	fmt.Println(err.Error())
	// 	// }
	// 	// defer f.Close()
	// 	// _, err2 := f.WriteString(outString)

	// 	// if err2 != nil {
	// 	// 	fmt.Println(err2.Error())
	// 	//}

	// 	// temp := strings.Split(outString, "\n")

	// 	// for linenum, line := range temp {
	// 	// 	fmt.Println(linenum, line)
	// 	// }

	// }
	return out
}

func discoverRename() {
	out := gitStatus()

	renamedLine := re.MustCompile(`(renamed.*?)(?:\r|\n|\r\n)`)

	found := renamedLine.Find(out)

	if found != nil {
		fmt.Printf("%q\n", found)
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
