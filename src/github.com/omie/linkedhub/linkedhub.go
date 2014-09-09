/*
    User : https://api.github.com/users/Omie
        returns a dict
        has repos_url : https://api.github.com/users/Omie/repos
    Repos : https://api.github.com/users/Omie/repos
        returns a list of dict
        has collaborators_url : https://api.github.com/repos/Omie/configfiles/collaborators
    Collaborators : https://api.github.com/repos/Omie/configfiles/collaborators
        returns a list of dict
        has repos_url for each user
*/

package main

import (
        "encoding/json"
        "io/ioutil"
        "net/http"
        "strings"
        "log"
)

var visitedUsers, visitedRepos = make(map[string]string), make(map[string]string)
var visitedColab  = map[string]string{}

//because Math.Min is for float64
func min(a, b int) int {
    if a <= b {
        return a
    }
    return b
}

func getData(url string) ([]byte, error) {
        log.Println("--- reached getData for ", url)

        resp, err := http.Get(url)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }

        return body, nil
}
//TODO: declare custom type for User use Tag

//TODO: may be returning only single value back is much better than entire data structure
func jsonToMap(jsondata []byte) (map[string]interface{}, error) {
        log.Println("--- reached jsonToMap")

        var f map[string]interface{}
        if err := json.Unmarshal(jsondata, &f); err != nil {
            return nil, err
        }

        return f, nil
}

func jsonToList(jsondata []byte) ([]interface{}, error) {
        log.Println("--- reached jsonToList")

        var f []interface{}
        if err := json.Unmarshal(jsondata, &f); err != nil {
            return nil, err
        }

        return f, nil
}

func getReposURL(username string) (string, error) {
        log.Println("--- reached getReposURL for ", username)

        data, err := getData("https://api.github.com/users/" + username)
        if err != nil {
            return "", err
        }

        json, err := jsonToMap(data)
        if err != nil {
            return "", err
        }

        str := json["repos_url"].(string)
        return str, nil
}

func processCollaborators(collabURL string) {
        log.Println("--- reached processCollaborators")

        data, err := getData(collabURL)
        if err != nil {
            return
        }

        jsonList, err := jsonToList(data)
        if err != nil {
            return
        }

        //for each collaborator
        for _, v := range jsonList {
            m, ok := v.(map[string]interface{})
            if !ok {
                log.Println("collab not ok")
                break
            }

            //handle user if not previously listed
            tempUser := m["login"].(string)
            if _, exists := visitedUsers[tempUser]; exists {
                continue
            }

            //We found new user in network
            log.Println("User : ", tempUser)
            visitedUsers[tempUser] = tempUser
            tempRepoURL := m["repos_url"].(string)

            //make a call to processRepo(tempRepoURL)
            processRepos(tempRepoURL)
            visitedRepos[tempRepoURL] = tempRepoURL
        } //end for
}

func processRepos(repoURL string) {
        log.Println("--- reached processRepos for ", repoURL)

        data, err := getData(repoURL) //get a list of repositories
        if err != nil {
            log.Println("err while getting data", err)
            return
        }

        jsonList, err := jsonToList(data)
        if err != nil {
            log.Println("err while converting list", err)
            return
        }

        m := min(len(jsonList), 2)
        jsonList = jsonList[:m] //limit to only 2 entries for time being

        //for each repo
        for _, v := range jsonList {
            m, ok := v.(map[string]interface{})
            if !ok {
                log.Println("err while converting to map")
                return
            }

            //handle collabs only if this repo is NOT previously visited
            tempCollabsURL := m["collaborators_url"].(string)
            idx := strings.Index(tempCollabsURL, "{")
            //use bytes package for serious string manipulation. much faster
            collabURL := tempCollabsURL[:idx]

            if repo, exists := visitedColab[collabURL]; exists {
                log.Println("already visited repo: ", repo )
                continue
            }

            visitedColab[collabURL] = collabURL
            processCollaborators(collabURL)

        } //end for
} //end processRepos

func main() {

    //get username from command line
    var u string = "Omie"
    //fmt.Println("Enter github username: ")
    //fmt.Scanln(&u)

    repoURL, err := getReposURL(u)
    if err != nil {
        log.Println("error while getting repo url for: ", u)
        return
    }

    processRepos(repoURL)
}

