package main


type User struct {
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
    Type string `json:"type"`
    Site_admin bool `json:"site_admin"`
}
