package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var users = make(chan string, 5)
var repos = make(chan string, 5)
var collabs = make(chan string, 5)

var visitedUsers, visitedRepos = make(map[string]string), make(map[string]string)
var visitedColab  = make(map[string]string)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getData(url string) []byte {
    //fmt.Println("--- reached getData for ", url)
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	return body
}

func JsonToMap(jsondata []byte) map[string]interface{} {
    //fmt.Println("--- reached JsonToMap")
	var f interface{}
	err := json.Unmarshal(jsondata, &f)
	check(err)

	m := f.(map[string]interface{})
	return m
}

func JsonToList(jsondata []byte) []interface{} {
    //fmt.Println("--- reached JsonToList")
	var f interface{}
	err := json.Unmarshal(jsondata, &f)
	check(err)

	m := f.([]interface{})
	return m
}

func getReposURL(username string) string {
    //fmt.Println("--- reached getReposURL for ", username)
	data := getData("https://api.github.com/users/" + username)
	json := JsonToMap(data)
	str := json["repos_url"].(string)
	return str
}

func processCollaborators() {
	for {
        //fmt.Println("--- reached processCollaborators")
		collab_url := <-collabs
        //fmt.Println(collab_url)
		data := getData(collab_url)
		json := JsonToList(data)
		//for each collaborator
		for _, v := range json {
			if m, ok := v.(map[string]interface{}); ok {

				tempu := m["login"].(string)
				//handle user if not previously listed
				if _, ok = visitedUsers[tempu]; !ok {
					users <- tempu
                    visitedUsers[tempu] = tempu

                    tempr := m["repos_url"].(string)
					repos <- tempr
                    visitedRepos[tempr] = tempr
				} else {
                    //fmt.Println("Already printed: ", tempu)
                }
			}else{
                fmt.Println("collab not ok")
            }
		}
	} //end for
}

func processRepos() {
	for {
        //fmt.Println("--- reached processRepos")
		repo := <-repos
		data := getData(repo) //get a list of repositories
		json := JsonToList(data)
		//for each repo
		for _, v := range json {
			if m, ok := v.(map[string]interface{}); ok {
				//handle collabs only if this repo is NOT previously visited
                temps := m["collaborators_url"].(string)
                idx := strings.Index(temps, "{")
                collab_url := temps[:idx]
				if _, ok := visitedColab[collab_url]; !ok {
                    visitedColab[collab_url] = collab_url
                    collabs <- collab_url
				} else {
                    //fmt.Println("already visited repo: ", repo)
                }
			}else{
                //fmt.Println("repos not ok")
            }
		}
	}
}

func processUsers() {
	for {
        //fmt.Println("--- reached processUsers")
		user := <-users
		fmt.Println(user)
	}
}

func main() {

	//get username from command line
    var u string
    fmt.Println("Enter github username: ")
    fmt.Scanln(&u)

	users <- u
	url := getReposURL(u)
	repos <- url

    visitedUsers[u] = u

	go processUsers()
	go processRepos()
	go processCollaborators()

	fmt.Scanln()
	fmt.Println("Done")

}
