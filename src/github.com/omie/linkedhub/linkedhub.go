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
        "os"
        "fmt"
        "encoding/json"
        "io/ioutil"
        "net/http"
        "strings"
        "log"
        "errors"
        "github.com/omie/ghlib"
)

var visited = make(map[string]string)

var requestsLeft int = 60

var username, password string

//because Math.Min is for float64
func min(a, b int) int {
    if a <= b {
        return a
    }
    return b
}

func getData(url string) ([]byte, error) {
        log.Println("--- reached getData for ", url)

        requestsLeft--
        if requestsLeft < 0 {
            log.Println("--- LIMIT REACHED ")
            return nil, errors.New("limit reached")
        }

        client := &http.Client{}

        /* Authenticate */
        req, err := http.NewRequest("GET", url, nil)
        req.SetBasicAuth(username, password)
        resp, err := client.Do(req)
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

func getApiLimit() (int, error) {
    jsonData, err := getData("https://api.github.com/rate_limit")
    if err != nil {
        return 0, err
    }

    var limitData ghlib.GhLimit
    if err := json.Unmarshal(jsonData, &limitData); err != nil {
        return 0, err
    }
    return limitData.Rate.Remaining, nil
}

func getReposURL(username string) (string, error) {
        log.Println("--- reached getReposURL for ", username)

        userJsonData, err := getData("https://api.github.com/users/" + username)
        if err != nil {
            return "", err
        }

        var user ghlib.GhUser
        if err := json.Unmarshal(userJsonData, &user); err != nil {
            return "", err
        }
        return user.ReposUrl, nil
}

func processCollaborators(collabURL string) {
        log.Println("--- reached processCollaborators for ", collabURL)
        if _, exists := visited[collabURL]; exists {
            log.Println("--- skipped ", collabURL)
            return
        }
        visited[collabURL] = collabURL

        jsonData, err := getData(collabURL)
        if err != nil {
            return
        }

        var collaborators []*ghlib.GhUser
        err = json.Unmarshal(jsonData, &collaborators)
        if err != nil {
            log.Println("Error while parsing collaborators: ", err)
            return
        }
        //for each collaborator
        for _, collaborator := range collaborators {
            //handle user if not previously listed
            tempUser := collaborator.Login
            if _, exists := visited[tempUser]; exists {
                continue
            }
            //We found new user in network
            fmt.Println("User : ", tempUser)
            visited[tempUser] = tempUser
            tempRepoURL := collaborator.ReposUrl

            //make a call to processRepo(tempRepoURL)
            processRepos(tempRepoURL)
        } //end for
}

func processRepos(repoURL string) {
        log.Println("--- reached processRepos for ", repoURL)
        if _, exists := visited[repoURL]; exists {
            log.Println("--- skipped ", repoURL)
            return
        }
        visited[repoURL] = repoURL

        repoData, err := getData(repoURL) //get a list of repositories
        if err != nil {
            log.Println("err while getting data", err)
            return
        }

        var repoList []*ghlib.GhRepository
        err = json.Unmarshal(repoData, &repoList)
        if err != nil {
            log.Println("Error while parsing repo list: ", err)
            return
        }

        //m := min(len(repoList), 2)
        //repoList = repoList[:m] //limit to only 2 entries for time being

        for _, repo := range repoList {
            tempCollabsURL := repo.CollaboratorsUrl
            log.Println(tempCollabsURL)
            idx := strings.Index(tempCollabsURL, "{")
            //use bytes package for serious string manipulation. much faster
            collabURL := tempCollabsURL[:idx]
            processCollaborators(collabURL)
        }

} //end processRepos

func setLogging() error {
    f, err := os.OpenFile("/tmp/linkedhub.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        return err
    }
    defer f.Close()

    log.SetOutput(f)
    return nil
}

func main() {

    err := setLogging()
    if err != nil {
        fmt.Println("Could not open file for logging")
        return
    }

    fmt.Println("Enter github credentials")
    fmt.Print("username: ")
    fmt.Scanln(&username)
    fmt.Print("password: ")
    fmt.Scanln(&password)

    //find out current API limit
    limit, err := getApiLimit()
    if err != nil {
        fmt.Println("error while getting limit: ", err)
        return
    }
    if limit <= 10 {
        fmt.Println("Too few of API calls left. Not worth it.")
        return
    }
    requestsLeft = limit
    fmt.Println("requestsLeft: ", requestsLeft)

    //get username from command line
    var u string
    fmt.Println("Enter github username: ")
    fmt.Scanln(&u)

    repoURL, err := getReposURL(u)
    if err != nil {
        log.Println("error while getting repo url for: ", u)
        return
    }

    processRepos(repoURL)
}

