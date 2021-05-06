package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/enrichman/ringo/biscuit"
	"github.com/manifoldco/promptui"
)

var secretsFilename string

func main() {
	flag.StringVar(&secretsFilename, "filename", "config/secrets.yml", "the file containing the secrets to decrypt")
	flag.Parse()

	biscuit, err := biscuit.NewClient(secretsFilename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("'config/secrets.yml' file not found!\nYou can use the -filename flag for a different path")
			return
		}
		panic(err)
	}

	list, err := biscuit.List()
	if err != nil {
		panic(err)
	}

	prompt := promptui.Select{
		Label: "Select a secret to decrypt",
		Items: list,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return
	}

	fmt.Printf("Decrypting %q\n", result)
	res, err := biscuit.Get(result)
	if err != nil {
		fmt.Println(res)
		return
	}
	fmt.Println(res)
}
