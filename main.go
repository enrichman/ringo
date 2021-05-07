package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/enrichman/ringo/biscuit"
	"github.com/manifoldco/promptui"
)

var secretsFilename string

func main() {
	flag.StringVar(&secretsFilename, "filename", "config/secrets.yml", "the file containing the secrets to decrypt")
	flag.Parse()

	if _, err := os.Open(secretsFilename); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("'config/secrets.yml' file not found!\nYou can use the -filename flag for a different path")
			return
		}
		panic(err)
	}

	if err := biscuit.KmsCallerIdentity(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	list, err := biscuit.List(secretsFilename)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	prompt := promptui.Select{
		Label: "Select a secret to decrypt",
		Items: list,
		Size:  10,
		Searcher: func(input string, index int) bool {
			name := strings.Replace(strings.ToLower(list[index]), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		// this happens when you ^C the prompt
		return
	}

	fmt.Printf("Decrypting %q\n", result)

	res, err := biscuit.Get(secretsFilename, result)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println(res)
}
