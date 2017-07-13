package main

import (
	"fmt"
	"reflect"

	algorithmia "github.com/algorithmiaio/algorithmia-go"
)

var apiKey = "ALGORITHMIA_API_KEY"
var client = algorithmia.NewClient(apiKey, "")

func main() {
	input := "http://www.isitnude.com.s3-website-us-east-1.amazonaws.com/assets/images/sample/young-man-by-the-sea.jpg"
	//input := "http://az616578.vo.msecnd.net/files/2016/12/30/6361871495693720411264482533_friendship.jpg"

	algo, _ := client.Algo("algo://sfw/NudityDetection/1.1.6")
	resp, err := algo.Pipe(input)
	if err != nil {
		fmt.Println(resp, err)
		return
	}

	response := resp.(*algorithmia.AlgoResponse)

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		fmt.Printf("Convert Error")
	}

	for _, item := range result {
		r := reflect.TypeOf(item).Kind()
		switch r {
		case reflect.Float64:
			if item.(float64) < 0.85 {
				fmt.Printf("Uncertain: %v\n", item)
			} else {
				fmt.Printf("Certain: %v\n", item)
			}
		case reflect.String:
			if item.(string) == "true" {
				fmt.Println("Nude")
			} else {
				fmt.Println("Not nude")
			}
		default:
		}
	}

}
