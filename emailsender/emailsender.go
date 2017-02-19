package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
)

func main() {

	//Carrega os destinatários em uma slice
	var emailDestinatarios = getEmailDestinatarios("./emaildestinatarios.json")

	for _, item := range emailDestinatarios {
		//dispara um e-mail para cada um dos destinatários
		disparaEmailDestinatario(item)
	}

}

func disparaEmailDestinatario(emailDestinatario EmailDestinatarios) {
	//Realiza o setup da autorização do servidor de SMTP. Não se esqueça de configuar seu Gmail SMTP server...
	//https://support.google.com/a/answer/176600?hl=en
	//https://support.google.com/accounts/answer/6010255?hl=en

	hostname := "smtp.gmail.com"
	auth := smtp.PlainAuth("", "seuemail@gmail.com", "suasenhagmail", hostname)

	//Criamos um slice do tipo string do tamanho máximo de 1 para receber nosso e-mail destinatário.
	recipients := make([]string, 1)
	recipients[0] = emailDestinatario.Email

	//Veja mais em: https://golang.org/pkg/net/smtp/#SendMail
	err := smtp.SendMail(
		hostname+":25",
		auth, "seuemail@gmail.com",
		recipients,
		/*Mensagem no RFC 822-style*/
		[]byte("Subject:Olá!\n\n Olá "+emailDestinatario.Nome+". Tudo de bom com Go!"))
	if err != nil {
		log.Fatal(err)
	}

}

func getEmailDestinatarios(file string) []EmailDestinatarios {
	//Realiza a leitura do arquivo json
	raw, err := ioutil.ReadFile(file)

	//Tratamento de erros padrão.
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var emailDestinatarios []EmailDestinatarios

	//Unmarshal do conteúdo do arquivo json para um tipo struct EmailDestinatarios
	json.Unmarshal(raw, &emailDestinatarios)
	return emailDestinatarios
}

//EmailDestinatarios : Lista de e-mails dos destinatários
type EmailDestinatarios struct {
	Nome  string
	Email string
}
