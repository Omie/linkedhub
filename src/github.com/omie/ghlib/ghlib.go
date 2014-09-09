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
    AvatarUrl string `json:"avatar_url"`
    GravatarId string `json:"gravatar_id"`
    Url string `json:"url"`
    HtmlUrl string `json:"html_url"`
    FollowersUrl string `json:"followers_url"`
    FollowingUrl string `json:"following_url"`
    GistsUrl string `json:"gists_url"`
    StarredUrl string `json:"starred_url"`
    SubscriptionsUrl string `json:"subscriptions_url"`
    OrganizationsUrl string `json:"organizations_url"`
    ReposUrl string `json:"repos_url"`
    EventsUrl string `json:"events_url"`
    ReceivedEvents_url string `json:"received_events_url"`
    Usertype string `json:"type"`
    SiteAdmin bool `json:"site_admin"`
    Name string `json:"name"`
    Company string `json:"company"`
    Blog string `json:"blog"`
    Location string `json:"location"`
    Email string "email"
    Hireable bool `json:"hireable"`
    Bio string `json:"bio"`
    PublicRepos int `json:"public_repos"`
    PublicGists int `json:"public_gists"`
    Followers int `json:"followers"`
    Following int `json:"following"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

type GhRepository struct {
    Id int `json:"id"`
    Name string `json:"name"`
    FullName string `json:"full_name"`
    Owner GhUser `json:"owner"`
    Private bool `json:"private"`
    HtmlUrl string `json:"html_url"`
    Description string `json:"description"`
    Fork bool `json:"fork"`
    Url string `json:"url"`
    ForksUrl string `json:"forks_url"`
    KeysUrl string `json:"keys_url"`
    CollaboratorsUrl string `json:"collaborators_url"`
    TeamsUrl string `json:"teams_url"`
    HooksUrl string `json:"hooks_url"`
    IssueEvents_url string `json:"issue_events_url"`
    EventsUrl string `json:"events_url"`
    AssigneesUrl string `json:"assignees_url"`
    BranchesUrl string `json:"branches_url"`
    TagsUrl string `json:"tags_url"`
    BlobsUrl string `json:"blobs_url"`
    GitTags_url string `json:"git_tags_url"`
    GitRefs_url string `json:"git_refs_url"`
    TreesUrl string `json:"trees_url"`
    StatusesUrl string `json:"statuses_url"`
    LanguagesUrl string `json:"languages_url"`
    StargazersUrl string `json:"stargazers_url"`
    ContributorsUrl string `json:"contributors_url"`
    SubscribersUrl string `json:"subscribers_url"`
    SubscriptionUrl string `json:"subscription_url"`
    CommitsUrl string `json:"commits_url"`
    GitCommits_url string `json:"git_commits_url"`
    CommentsUrl string `json:"comments_url"`
    IssueComment_url string `json:"issue_comment_url"`
    ContentsUrl string `json:"contents_url"`
    CompareUrl string `json:"compare_url"`
    MergesUrl string `json:"merges_url"`
    ArchiveUrl string `json:"archive_url"`
    DownloadsUrl string `json:"downloads_url"`
    IssuesUrl string `json:"issues_url"`
    PullsUrl string `json:"pulls_url"`
    MilestonesUrl string `json:"milestones_url"`
    NotificationsUrl string `json:"notifications_url"`
    LabelsUrl string `json:"labels_url"`
    ReleasesUrl string `json:"releases_url"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
    PushedAt string `json:"pushed_at"`
    GitUrl string `json:"git_url"`
    SshUrl string `json:"ssh_url"`
    CloneUrl string `json:"clone_url"`
    SvnUrl string `json:"svn_url"`
    Homepage string `json:"homepage"`
    Size int `json:"size"`
    StargazersCount int `json:"stargazers_count"`
    WatchersCount int `json:"watchers_count"`
    Language string `json:"language"`
    HasIssues bool `json:"has_issues"`
    HasDownloads bool `json:"has_downloads"`
    HasWiki bool `json:"has_wiki"`
    ForksCount int `json:"forks_count"`
    MirrorUrl string `json:"mirror_url"`
    OpenIssues_count string `json:"open_issues_count"`
    Forks int `json:"forks"`
    OpenIssues int `json:"open_issues"`
    Watchers int `json:"watchers"`
    DefaultBranch string `json:"default_branch"`
}

type GhRepoList struct {
    Repositories []GhRepository
}

