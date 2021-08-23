package main

import (
	"encoding/json"
	"fmt"
	"github/gorilla/mux"
	"io/ioutil"

	"log"
	"net/http"

	"time"
)

const month = time.Hour * 24 * 30

var topLangs map[string][]string

func home(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, topLangs)
}

func getRepos() error {
	topLangs = make(map[string][]string)
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/search/repositories?q=created:>%v&sort=stars&order=desc", time.Now().Add(month*-1).Format("2006-01-02")))
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var results map[string]interface{}
	err = json.Unmarshal(data, &results)
	if err != nil {
		return err
	}
	items := results["items"].([]interface{})
	for i, item := range items {
		if i == 100 {
			break
		}

		repo := item.(map[string]interface{})
		if repo["language"] != nil {
			lang := repo["language"].(string)
			topLangs[lang] = append(topLangs[lang], repo["html_url"].(string))
		}
	}
	return nil
}

func main() {
	err := getRepos()
	if err != nil {
		fmt.Print(err)
		return
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8081", router))
}
