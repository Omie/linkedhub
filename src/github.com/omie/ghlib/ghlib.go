package ghlib

/*  types for github API limit check */

type GhLimitResourcesCore struct {
    Limit int `json:"limit"`
    Remaining int `json:"remaining"`
    reset int `json:"reset"`
}

type GhLimitResourcesSearch struct {
    Limit int `json:"limit"`
    Remaining int `json:"remaining"`
    reset int `json:"reset"`
}

type GhLimitResources struct {
    Core GhLimitResourcesCore `json:"core"`
    Search GhLimitResourcesSearch `json:"search"`
}

type GhLimitRate struct {
    Limit int `json:"limit"`
    Remaining int `json:"remaining"`
    reset int `json:"reset"`
}

type GhLimit struct {
    Resources GhLimitResources `json:"resources"`
    Rate GhLimitRate `json:"rate"`
}

/*  Github API types */

type GhUser struct {
    Login string `json:"login"`
    Id int `json:"id"`
    Avatar_url string `json:"avatar_url"`
    Gravatar_id string `json:"gravatar_id"`
    Url string `json:"url"`
    Html_url string `json:"html_url"`
    Followers_url string `json:"followers_url"`
    Following_url string `json:"following_url"`
    Gists_url string `json:"gists_url"`
    Starred_url string `json:"starred_url"`
    Subscriptions_url string `json:"subscriptions_url"`
    Organizations_url string `json:"organizations_url"`
    Repos_url string `json:"repos_url"`
    Events_url string `json:"events_url"`
    Received_events_url string `json:"received_events_url"`
    Usertype string `json:"type"`
    Site_admin bool `json:"site_admin"`
    Name string `json:"name"`
    Company string `json:"company"`
    Blog string `json:"blog"`
    Location string `json:"location"`
    Email string "email"
    Hireable bool `json:"hireable"`
    Bio string `json:"bio"`
    Public_repos int `json:"public_repos"`
    Public_gists int `json:"public_gists"`
    Followers int `json:"followers"`
    Following int `json:"following"`
    Created_at string `json:"created_at"`
    Updated_at string `json:"updated_at"`
}


