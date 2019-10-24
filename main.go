package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	color.Yellow("Insira o produto:")
	product, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	product = strings.TrimRight(product, "\n")
	color.Green("O produto digitado foi %s \n", product)
	ter := openFile()
	treatString(product, ter)
	defer ter.Close()
}

func openFile() *os.File {

	f, err := os.OpenFile("categorization.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

func treatString(s string, f *os.File) {
	s = strings.Replace(s, ",", "", -1)
	arr := make([]string, 10)
	arr = strings.Split(s, " ")
	m := getEverythingJSON(f)
	writeIt := compareJSON(m, &arr)
	jsonString, _ := json.Marshal(writeIt)
	f.Write(jsonString)

}

func getEverythingJSON(f *os.File) *map[string][]string {
	filer, _ := ioutil.ReadFile("categorization.json")
	c := make(map[string][]string)
	err := json.Unmarshal(filer, &c)
	if err != nil {
		return &c
	}

	return &c
}

func compareJSON(m *map[string][]string, sm *[]string) map[string][]string {

	arri := *m

	keys := []string{}

	for key := range *m {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		saveMap := askForAnswer(*sm, &arri)
		return saveMap
	}

	missing := []string{}
	finding := []string{}
	for _, t := range keys {
		for _, a := range arri[t] {
			for _, s := range *sm {
				if s == a {
					if t != "ignore" {
						retrieveAnswer(s, t)
					}
					finding = append(finding, s)
				} else {
					if !contains(missing, s) {
						missing = append(missing, s)
					}

				}
			}
		}
	}

	rMissing := []string{}
	for _, teste := range missing {
		if !contains(finding, teste) {
			rMissing = append(rMissing, teste)
		}
	}

	saveMap := askForAnswer(rMissing, &arri)
	return saveMap

}

func askForAnswer(s []string, m *map[string][]string) map[string][]string {

	result := *m

	for _, each := range s {
		result = getInfoFromBash(each, result)
	}
	return result
}

func getInfoFromBash(s string, m map[string][]string) map[string][]string {
	color.Yellow("categoria da palavra %s ", s)
	reader := bufio.NewReader(os.Stdin)

	category, err := reader.ReadString('\n')
	if err != nil {
		panic("erro na hora de ler a linha da categoria")
	}

	category = strings.TrimRight(category, "\n")
	if category != "" {
		m[category] = append(m[category], s)
	} else {
		m["ignore"] = append(m["ignore"], s)
	}

	return m
}

func retrieveAnswer(word string, categorie string) {
	color.Green("%s ------------- %s", word, categorie)
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
