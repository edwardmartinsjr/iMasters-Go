package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//"anaconda"
	//"net/url"
)

func main() {

	lst1 := afterClassifier(getTermClassifier(getTwitterSentimentClassifier("./twittersentimentclassifier.json")))

	// for _, item := range lst1 {
	// 	fmt.Print(item.Term + "	" + strconv.Itoa(item.FreqDist) + "\n")
	// }

	_, lst2 := getTermClassifier(getTwitterSentimentClassifier("./twitter.json"))
	//_, lst2 := getTermClassifier(twitterStreamingAPI())

	for i, item := range lst2 {
		var count = 0
		for _, itemTerm := range item.Term {
			for _, itemClassifier := range lst1 {
				if itemTerm == itemClassifier.Term {
					//Aqui definimos qual frequencia de distribuição (score) será considerado como um sentimento positivo ou negativo
					if itemClassifier.FreqDist >= 0 {
						count++
					} else {
						count--
					}
				}
			}

		}
		if count >= 0 {
			lst2[i].Classifier = "positive"
		} else {
			lst2[i].Classifier = "negative"
		}
	}

	_, originalTweetList := getTermClassifier(getTwitterSentimentClassifier("./twittersentimentclassifier.json"))

	outLine("Tweets originais:", originalTweetList)
	outLine("Tweets classificados:", lst2)

}

// func twitterStreamingAPI() []TwitterSentimentClassifier {
// 	//Mais sobre: https://github.com/ChimeraCoder/anaconda

// 	anaconda.SetConsumerKey("xxx")                                                                                     //Consumer Key
// 	anaconda.SetConsumerSecret("xxx")                                                             //Consumer Secret
// 	client := anaconda.NewTwitterApi("xxx", "xxx") //Access Token, Access Token Secret

// 	// setando os parametros utilizando url.Values
// 	v := url.Values{}
// 	v.Set("count", "30") // ou v.Set("locations", "<Locations>")
// 	result, err := client.GetSearch("golang", nil)//buscar por tweets que contenham o termo “golang”
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	// Ao menos que exista algo estranho, devemos ter ao menos 2 tweets
// 	if len(result.Statuses) < 2 {
// 		fmt.Printf("Esperado 2 ou mais tweets, foram encontrados %d", len(result.Statuses))
// 		os.Exit(1)
// 	}

// 	twitterSentimentClassifier := make([]TwitterSentimentClassifier, len(result.Statuses))

// 	// verificar a existência de tweet vazio
// 	for i, tweet := range result.Statuses {
// 		twitterSentimentClassifier[i].Tweet = tweet.Text
// 	}

// 	return twitterSentimentClassifier
// }

func outLine(text string, lst []TermClassifier) {
	fmt.Print(" \n" + text + " \n")
	for _, itemClassifier := range lst {
		for _, itemTerm := range itemClassifier.Term {
			fmt.Print(itemTerm + " ")
		}
		fmt.Print(" - " + itemClassifier.Classifier + " \n")
	}
}

func getTwitterSentimentClassifier(file string) []TwitterSentimentClassifier {
	//Realiza a leitura do arquivo json
	raw, err := ioutil.ReadFile(file)

	//Tratamento de erros padrão.
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var twitterSentimentClassifier []TwitterSentimentClassifier

	//Unmarshal do conteúdo do arquivo json para um tipo struct TwitterSentimentClassifier
	json.Unmarshal(raw, &twitterSentimentClassifier)
	return twitterSentimentClassifier
}

func getTermClassifier(twitterSentimentClassifier []TwitterSentimentClassifier) (int, []TermClassifier) {

	termClassifier := make([]TermClassifier, len(twitterSentimentClassifier))
	var generalCount int
	for i, item := range twitterSentimentClassifier {

		//Primeiro vamos fazer através da função Fields o split da sentença por espaços
		tweet := strings.Fields(item.Tweet)

		//Criamos um slice do tipo Term do tamanho máximo dos splits do nosso tweet.
		term2Classifier := make([]string, len(tweet))

		var count int
		for j, termTweet := range tweet {
			//Estamos considerando apenas palavras maiores que três caracteres para serem consideradas como termos válidos
			//Utilizamos rune para prevenir caracteres especiais, acentos, caracteres asiaticos e também emogis
			if len([]rune(termTweet)) >= 3 {
				term2Classifier[j] = strings.ToLower(termTweet)
				count++
			}
		}

		termClassifier[i].Term = make([]string, count)
		count = 0
		for k, termTweetClassifier := range term2Classifier {
			//Realizamos um ajuste do tamanho da slice final de termos
			if term2Classifier[k] != "" {
				termClassifier[i].Term[count] = termTweetClassifier
				count++
				generalCount++
			}
		}

		termClassifier[i].Classifier = item.Classifier
	}

	return generalCount, termClassifier
}

func afterClassifier(generalCount int, termClassifier []TermClassifier) []TermClassified {
	//Esta slice receberá todos os termos identificados na slice anterior
	termClassified := make([]TermClassified, generalCount)

	var count int
	var countTermAfterClassifier int
	for _, item := range termClassifier {
		for _, itemTerm := range item.Term {

			//Antes de aplicar o score em um termo verificamos se ele já não fora identificado anteriormente.
			//Caso este termo já tenha sido identificado apenas contabilizamos a frequencia de distribuição (score)
			var found bool
			for _, itemTermForCompareBeforeInsert := range termClassified {
				if itemTermForCompareBeforeInsert.Term == itemTerm {
					found = true
					break
				}
			}
			if !found {
				termClassified[count].Term = itemTerm
				countTermAfterClassifier++
				//Agora iremos aplicar a frequencia de distribuição (score) de cada termo em relação ao sentimento que demos em cada um dos tweets
				for _, itemForCompare := range termClassifier {
					for _, itemTermToCompare := range itemForCompare.Term {

						if itemTerm == itemTermToCompare {
							if itemForCompare.Classifier == "positive" {
								termClassified[count].FreqDist = termClassified[count].FreqDist + 1
							} else if itemForCompare.Classifier == "negative" {
								termClassified[count].FreqDist = termClassified[count].FreqDist - 1
							}

						}
					}
				}
			}

			count++
		}
	}

	// Removendo registros vazios
	termAfterClassifierClassified := make([]TermClassified, countTermAfterClassifier)
	var countAfterClassifierClassified int
	for _, itemTermClassified := range termClassified {
		if itemTermClassified.Term != "" {
			termAfterClassifierClassified[countAfterClassifierClassified] = itemTermClassified
			countAfterClassifierClassified++
		}

	}

	return termAfterClassifierClassified
}

//TwitterSentimentClassifier : Tweets para treinar o classificador.
type TwitterSentimentClassifier struct {
	Tweet      string
	Classifier string
}

//TermClassified : Termos classificados.
type TermClassified struct {
	Term     string
	FreqDist int
}

//TermClassifier : Classificador dos Termos dos tweets para treinar o classificador
type TermClassifier struct {
	Term       []string
	Classifier string
}
