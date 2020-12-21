package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	res, err := http.Get("https://vercel-func-mu.vercel.app/index_small.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(len(bytes))
}

func TestSearch(t *testing.T) {
	fmt.Println("212325964a6633aad9460a9a8301093e" == "212325964a6633aad9460a9a8301093e")
}
