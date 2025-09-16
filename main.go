package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"
)

type RepositoriesJson struct {
	Repositories []string `json:"repositories"`
}

func main() {
	repositoriesBytes, err := os.ReadFile("repositories.json")

	if err != nil {
		fmt.Println("os.ReadFile(): ", err.Error())
		return
	}

	var repositoriesJson RepositoriesJson
	err = json.Unmarshal(repositoriesBytes, &repositoriesJson)

	if err != nil {
		fmt.Println("json.Unmarshal(): ", err.Error())
		return
	}

	repositories := repositoriesJson.Repositories
	length := len(repositories)

	if length == 0 {
		fmt.Println("'repositories.json' is empty")
		return
	}

	if slices.Contains(repositories, "username-or-organization/repository") {
		fmt.Println("Remove 'username-or-organization/repository' from 'repositories.json'")
		return
	}

	parentDirectory := time.Now().Format("2006-01-02 15-04-05 MST")
	err = os.RemoveAll(parentDirectory)

	if err != nil {
		fmt.Println("os.RemoveAll(): ", err.Error())
		return
	}

	err = os.Mkdir(parentDirectory, 0755)

	if err != nil {
		fmt.Println("os.Mkdir(): ", err.Error())
		return
	}

	err = os.Chdir(parentDirectory)

	if err != nil {
		fmt.Println("os.Chdir(parentDirectory): ", err.Error())
		return
	}

	github := "https://github.com/"

	if len(os.Args) > 1 && os.Args[1] == "ssh" {
		github = "git@github.com:"
	}

	for index, repository := range repositories {
		before, after, found := strings.Cut(repository, "/")

		if before == "" || after == "" || !found {
			fmt.Println("Invalid repository: '" + repository + "'")
			return
		}

		fmt.Println("[" + strconv.Itoa(index+1) + "/" + strconv.Itoa(length) + "] " + repository)
		childDirectory := before + " " + after
		command := exec.Command("git", "clone", "--recursive", github+repository, childDirectory)
		err = command.Run()

		if err != nil {
			fmt.Println("command.Run(): ", err.Error())
			return
		}

		err = os.Chdir(childDirectory)

		if err != nil {
			fmt.Println("os.Chdir(childDirectory): ", err.Error())
			return
		}

		command = exec.Command("git", "log", "--reverse", "--format=%as")
		date, err := command.Output()

		if err != nil {
			fmt.Println("command.Output(): ", err.Error())
			return
		}

		err = os.Chdir("..")

		if err != nil {
			fmt.Println("os.Chdir('..'): ", err.Error())
			return
		}

		err = os.Rename(childDirectory, string(date)[:10]+" "+childDirectory)

		if err != nil {
			fmt.Println("os.Rename(): ", err.Error())
			return
		}
	}
}
