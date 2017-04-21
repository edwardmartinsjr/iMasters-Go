package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/liviosoares/go-watson-sdk/watson"
	"github.com/liviosoares/go-watson-sdk/watson/conversation"
)

func main() {
	config := watson.Config{
		Credentials: watson.Credentials{
			Username: "Username",
			Password: "Password",
			Url:      "https://gateway.watsonplatform.net/conversation/api",
		},
	}
	client, err := conversation.NewClient(config)
	if err != nil {
		log.Printf("Error: [%s]", err)
	}

	var messageResponse conversation.MessageResponse
	log.Printf("Welcome to Dropabot")

	messageResponse, err = client.Message("Workspace ID", "", nil)
	log.Printf("Dropabot:  %s", messageResponse.Output)

	reader := bufio.NewReader(os.Stdin)
	for {

		// lendo entrada do terminal
		texto, _ := reader.ReadString('\n')
		log.Printf("You: %s", texto)

		texto = strings.Replace(texto, "\n", "", -1)
		messageResponse, err = client.Message("Workspace ID", texto, messageResponse.Context)
		if err != nil {
			log.Printf("Error: [%s]", err)
		}

		// escrevendo a resposta do bot no terminal
		log.Printf("Dropabot:  %s", messageResponse.Output)
	}
}
