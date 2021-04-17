package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://hnhiring.com/technologies/vue"

	content, err := loadPage(url)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}

func loadPage(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("URL [%v]: %v", url, err)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error procesando cotenido: %v", err)
	}

	return content, nil
}
