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

/*  types used for marhsalling data to json */
type node struct {
    Name string `json:"name"`
    Group int `json:"group"`
    Image string `json:"image"`
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

// keep track of visited URLs and previously known users
var visited = make(map[string]string)
var knownUsers = make(map[string]int)

// gets set to current limit in runtime
var requestsLeft = 60

// github a/c credentials
var username, password string

// determines how deep to crawl
var maxDepth int

// holds list of users found as nodes
// node-to-node connections using 0 based index
// as required for d3
var nodes []node
var connections []connection

// get data from remote url and return unparsed output
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
            log.Println("error in http request: ", err)
            return nil, err
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Println("error reading request body: ", err)
            return nil, err
        }

        return body, nil
}

// determine current API limit
func getAPILimit() (int, error) {
    jsonData, err := getData("https://api.github.com/rate_limit")
    if err != nil {
        return 0, err
    }

    var limitData ghlib.GhLimit
    if err := json.Unmarshal(jsonData, &limitData); err != nil {
        return 0, err
    }

    limit := limitData.Rate.Remaining
    if limit <= 10 {
        return 0, errors.New("Too few of API calls left. Not worth it.")
    }

    return limit, nil
}

// get User details from API, retrieve repos_url for the user and
// return the same
func getReposURL(username string) (string, error) {
        log.Println("--- reached getReposURL for ", username)

        userJSONData, err := getData("https://api.github.com/users/" + username)
        if err != nil {
            return "", err
        }

        var user ghlib.GhUser
        if err := json.Unmarshal(userJSONData, &user); err != nil {
            return "", err
        }
        return user.ReposUrl, nil
}

// repo contributors is a list of maps [ {}, {}]
// retrieve the list, parse json, for each user:
//      process user's repositories to move to further depth
func processContributors(contribURL string, currentDepth int, parent int) {
        log.Println("--- reached processContributors for ", contribURL)
        if _, exists := visited[contribURL]; exists {
            log.Println("--- skipped ", contribURL)
            return
        }
        visited[contribURL] = contribURL

        jsonData, err := getData(contribURL)
        if err != nil {
            log.Println("error while getting contributors data: ", err)
            return
        }

        var contributors []*ghlib.GhUser
        err = json.Unmarshal(jsonData, &contributors)
        if err != nil {
            log.Println("error while unmarshalling contributors: ", err)
            return
        }
        //for each contributor
        for _, contributor := range contributors {
            //if user is already visited, just mark a connection and move to next
            tempUser := contributor.Login
            if nodeIdx, exists := knownUsers[tempUser]; exists {
                connections = append(connections, connection{parent, nodeIdx, 1})
                continue
            }
            //We found new user in network
            /* this might slow down
            for t:=0; t<=currentDepth; t++ {
                fmt.Print("\t|")
            }
            */
            fmt.Println(tempUser)

            //push to nodes list and connection list
            nodes = append(nodes, node{tempUser, 1, contributor.AvatarUrl})
            nodeIdx := len(nodes)-1
            connections = append(connections, connection{parent, nodeIdx, 1})

            knownUsers[tempUser] = nodeIdx
            tempRepoURL := contributor.ReposUrl

            //get repositories of this new user
            processRepos(tempRepoURL, currentDepth+1, nodeIdx)

        } //end for
}

// process a list of repositories
// for each repository, find and process collaborators
func processRepos(repoURL string, currentDepth int, parent int) {
        log.Println("--- reached processRepos for ", repoURL)
        if currentDepth > maxDepth {
            log.Println("maxDepth reached")
            return
        }

        if _, exists := visited[repoURL]; exists {
            return
        }
        visited[repoURL] = repoURL

        repoData, err := getData(repoURL) //get a list of repositories
        if err != nil {
            log.Println("error while getting repo list data: ", err)
            return
        }

        var repoList []*ghlib.GhRepository
        err = json.Unmarshal(repoData, &repoList)
        if err != nil {
            log.Println("error while unmarshalling repo list: ", err)
            return
        }

        for _, repo := range repoList {
            contribURL := repo.ContributorsUrl
            processContributors(contribURL, currentDepth, parent)
        }

} //end processRepos

func handleUserInput() {
    fmt.Println("Enter Github credentials -")
    fmt.Print("username(Email address): ")
    fmt.Scanln(&username)

    fmt.Print("password: ")
    fmt.Scanln(&password)

    fmt.Print("Max depth: ")
    fmt.Scanln(&maxDepth)
}

func dumpD3Json() {
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

func main() {
    f, err := os.OpenFile("/tmp/linkedhub.log", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("Could not open file for logging")
        return
    }
    defer f.Close()

    log.SetOutput(f)
    //log.SetOutput(ioutil.Discard)

    handleUserInput()

    //find out current API limit
    limit, err := getAPILimit()
    if err != nil {
        fmt.Println("error while getting api limit: ", err)
        return
    }
    requestsLeft = limit
    fmt.Println("requests left for this hour: ", requestsLeft)

    //get username from command line
    var u string
    fmt.Println("Enter github username to start with: ")
    fmt.Scanln(&u)

    repoURL, err := getReposURL(u)
    if err != nil {
        log.Println("error while getting repo url for: ", u)
        return
    }

    processRepos(repoURL, 0, 0)
    fmt.Println("dumping ", string(len(nodes)), " nodes to json")
    dumpD3Json()
}


