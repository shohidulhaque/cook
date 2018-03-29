package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/crypto/bcrypt"
	ps "local.proj/Cook/parser"
)

type entity struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

type parent struct {
	Body struct {
		Entity []entity `json:"entity"`
	} `json:"body"`
}

//Never Liked Global variables but until I think of a workaround
var newfileTimings map[string]string
var oldfileTimings map[string]string
var hashJSONold parent
var hashJSONnew parent
var tagList []string
var fileList map[string]string

//Stop Go from throwing warnings if a variable is not used
func doNothing(str string) {
	//Go is badass
}

//Simple Error Checker
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Generate the list of files to be compiled
func generateFileList(tag string) {
	details := ps.FileDetails[tag]

	_, err := os.Stat(details.File)
	checkErr(err)

	fileList[tag] = details.File

	if details.Deps == nil {
		return
	}

	for _, name := range details.Deps {
		generateFileList(name)
	}
}

//Function for executing and debugging exec.Cmd
func checkCommand(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
}

//Generating hash from timestamp
func hashTime(timeStamp string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(timeStamp), 14)
	return string(bytes), err
}

//Comparing hashes of the current timestamp with the previous one
func checkTimeStamp(timeStamp string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(timeStamp))
	return err == nil
}
