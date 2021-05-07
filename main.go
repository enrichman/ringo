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

var (
	Version         = "v0.0.0-dev"
	secretsFilename string
)

func main() {
	showVersion := flag.Bool("version", false, "show ringo version")
	flag.StringVar(&secretsFilename, "filename", "config/secrets.yml", "the file containing the secrets to decrypt")
	flag.Parse()

	if *showVersion {
		fmt.Print(Version)
		os.Exit(0)
	}

	if _, err := os.Open(secretsFilename); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("'config/secrets.yml' file not found!\nYou can use the -filename flag for a different path")
			return
		}
		fmt.Print(err)
		os.Exit(1)
	}

	if err := biscuit.KmsCallerIdentity(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	secrets, err := biscuit.List(secretsFilename)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	selectedSecret, err := loadSecretsSelection(secrets)
	if err != nil {
		// this happens when you ^C the prompt
		return
	}

	fmt.Printf("Decrypting %q\n", selectedSecret)
	res, err := biscuit.Get(secretsFilename, selectedSecret)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println(res)
}

func loadSecretsSelection(secrets []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select a secret to decrypt",
		Items: secrets,
		Size:  10,
		Searcher: func(input string, index int) bool {
			name := strings.Replace(strings.ToLower(secrets[index]), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
	}

	_, result, err := prompt.Run()
	return result, err
}
