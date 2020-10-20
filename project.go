package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var secretMap = make(map[string]string)

var patternMap = make(map[string]string)
var patternStr string
var inputFilePath string
var outputFilePath string

var token = "dp.pt.Lz1zQrXCBXtQVeYsdeZE5XY3O6KVqMazwfYxFXB6"
var api = "api.doppler.com/v3/configs/config/secret"
var proj = "358632fea11"
var config = "dev"

func main() {
	patternMap["dollar"] = ".*\\$[a-zA-Z0-9_&*@#]+.*$"
	patternMap["dollar-curly"] = ".*\\${[a-zA-Z0-9_&*@#]+}.*$"
	patternMap["handlebars"] = "^{{[a-zA-Z0-9_&*@#]+}}.*$"
	patternMap["dollar-handlebars"] = ".*\\${{[a-zA-Z0-9_&*@#]+}}.*$"

	if len(os.Args) > 1 {
		if os.Args[1] == "input" {
			inputFilePath = os.Args[2]
		} else {
			fmt.Println("please pass the input as 1st argument")
			os.Exit(1)
		}
	} else {
		fmt.Println("please provide the input path")
		os.Exit(1)
	}

	if len(os.Args) > 3 {
		if os.Args[3] == "output" {
			outputFilePath = os.Args[4]
		} else if os.Args[3] == "pattern" {
			patternStr = patternMap[os.Args[4]]

			outputFilePath = inputFilePath
		}
	} else {
		outputFilePath = inputFilePath
	}

	if len(os.Args) > 5 {
		patternStr = patternMap[os.Args[6]]
	}

	fmt.Println("Input ", inputFilePath)
	fmt.Println("Output ", outputFilePath)
	fmt.Println("Pattern ", patternStr)

	if len(patternStr) == 0 {
		patternStr = patternMap["dollar-curly"]
		fmt.Println("Pattern ", patternStr)
	}

	// Check if directory or file
	name := inputFilePath
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// do directory stuff
		fmt.Println("directory")
	case mode.IsRegular():
		// do file stuff
		fmt.Println("file")
	}

	words, err := scanWords(inputFilePath)
	if err != nil {
		panic(err)
	}

	for _, word := range words {
		fmt.Println("Word " + word)
		trimWord := trimWord(word)

		secretMap[word] = callDopplerAPI(trimWord)
		//fmt.Println(secretMap[word])
	}

	err1 := filepath.Walk(inputFilePath, replaceSecret)
	if err1 != nil {
		panic(err1)
	}
}

func callDopplerAPI(secretVal string) string {
	type value struct {
		Raw      string
		Computed string
	}

	type result struct {
		Name  string
		Value value
	}

	urlString := fmt.Sprintf("https://%v@%v?project=%v&config=%v&name=%v", token, api, proj, config, secretVal)

	resp, err := http.Get(urlString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var data result
	json.Unmarshal(text, &data)
	//fmt.Printf("Results: %v\n", data.Value.Computed)
	//fmt.Println(string(text))

	return data.Value.Computed
}

func getURL(url string) string {
	token := "dp.pt.Lz1zQrXCBXtQVeYsdeZE5XY3O6KVqMazwfYxFXB6"
	domain := "api.doppler.com/v3/configs/config/secret?project="

	var urlString string
	urlString += "https://"
	urlString += token
	urlString += "@"
	urlString += domain

	return urlString
}

func trimWord(word string) string {
	// Make a Regex to say we only want letters
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(word, "")

	//fmt.Printf("A string of %s becomes %s \n", word, processedString)
	return processedString
}

func replaceSecret(inputPath string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil
	}

	fmt.Println("map:", secretMap)

	matched, err := filepath.Match("*.txt", fi.Name())

	if err != nil {
		//panic(err)
		fmt.Println(err)
		return err
	}

	if matched {
		read, err := ioutil.ReadFile(inputPath)
		if err != nil {
			panic(err)
		}

		for key, value := range secretMap {
			//fmt.Printf("Key: %s\tValue: %v\n", key, value)
			if value != "" {
				newContents := strings.Replace(string(read), key, value, -1)
				read = []byte(newContents)
			}
		}

		err = ioutil.WriteFile(outputFilePath, read, 0)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func scanWords(path string) ([]string, error) {
	r, _ := regexp.Compile(patternStr)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			//fmt.Println(r.FindAllStringIndex(scanner.Text(), -1))
			newWords := r.FindAllString(scanner.Text(), -1)
			for _, word := range newWords {
				words = append(words, word)
			}
		}
	}

	return words, nil
}
