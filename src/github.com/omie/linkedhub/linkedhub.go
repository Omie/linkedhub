/*
    User : https://api.github.com/users/Omie
        returns a dict
        has repos_url : https://api.github.com/users/Omie/repos
    Repos : https://api.github.com/users/Omie/repos
        returns a list of dict
        has contributors_url : https://api.github.com/repos/Omie/configfiles/contributors
    Contributors : https://api.github.com/repos/Omie/configfiles/contributors
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
        "log"
        "errors"
        "github.com/omie/ghlib"
)

type node struct {
    Name string `json:"name"`
    Group int `json:"group"`
}

type connection struct {
    Source int `json:"source"`
    Target int `json:"target"`
    Value int `json:"value"`
}

type graphdata struct {
    Nodes []node `json:"nodes"`
    Connections []connection `json:"links"`
}

var visited = make(map[string]string)

var requestsLeft int = 60

var username, password string

var maxDepth int = 0

var nodes []node
var connections []connection

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

func processContributors(contribURL string, currentDepth int, parent int) {
        log.Println("--- reached processContributors for ", contribURL)
        if _, exists := visited[contribURL]; exists {
            log.Println("--- skipped ", contribURL)
            return
        }
        visited[contribURL] = contribURL

        jsonData, err := getData(contribURL)
        if err != nil {
            return
        }

        var contributors []*ghlib.GhUser
        err = json.Unmarshal(jsonData, &contributors)
        if err != nil {
            log.Println("Error while parsing contributors: ", err)
            return
        }
        //for each contributor
        for _, contributor := range contributors {
            //handle user if not previously listed
            tempUser := contributor.Login
            if _, exists := visited[tempUser]; exists {
                continue
            }
            //We found new user in network
            for t:=0; t<=currentDepth; t++ {
                fmt.Print("\t|")
            }
            fmt.Print(tempUser, "\n")
            nodes = append(nodes, node{tempUser, 1})
            nodeIdx := len(nodes)-1
            connections = append(connections, connection{parent, nodeIdx, 1})

            visited[tempUser] = tempUser
            tempRepoURL := contributor.ReposUrl

            //make a call to processRepo(tempRepoURL)
            processRepos(tempRepoURL, currentDepth+1, nodeIdx)
        } //end for
}

func processRepos(repoURL string, currentDepth int, parent int) {
        log.Println("--- reached processRepos for ", repoURL)
        if currentDepth > maxDepth {
            log.Println("maxDepth reached")
            return
        }

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
            contribURL := repo.ContributorsUrl
            log.Println(contribURL)
            processContributors(contribURL, currentDepth, parent)
        }

} //end processRepos

func main() {
    f, err := os.OpenFile("/tmp/linkedhub.log", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("Could not open file for logging")
        return
    }
    defer f.Close()

    log.SetOutput(f)

    fmt.Println("Enter github credentials")
    fmt.Print("username: ")
    fmt.Scanln(&username)
    fmt.Print("password: ")
    fmt.Scanln(&password)
    fmt.Print("Max depth: ")
    fmt.Scanln(&maxDepth)

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

    processRepos(repoURL, 0, 0)

    fw, err := os.OpenFile("graph.json", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
    if err != nil {
        fmt.Println("Could not open file for writing json")
        return
    }
    defer fw.Close()

    var gdata = &graphdata {
        Nodes: nodes, 
        Connections: connections,
    }
    var toWrite []byte
    toWrite, err = json.Marshal(gdata)
    if err != nil {
        fmt.Println("Error marshalling data: ", err)
    }
    fw.Write(toWrite)

}


