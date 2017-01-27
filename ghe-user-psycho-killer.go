package main

import (
  "os"
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "regexp"
)

var linkListRegex = regexp.MustCompile("<([A-Za-z0-9\\:\\.\\/\\=\\?\\{\\}]*)>; rel=\"([A-Za-z]*)\"")

func main() {
  if os.Getenv("GHE_SERVER") == ""  {
    log.Fatal("Missing environment variable GHE_SERVER")
    return
  }
  if os.Getenv("GHE_WHITE_LIST_ORG") == "" {
    log.Fatal("Missing environment variable GHE_WHITE_LIST_ORG")
    return
  }
  if os.Getenv("GHE_PERSONAL_ACCESS_TOKEN") == "" {
    log.Fatal("Missing environment variable GHE_PERSONAL_ACCESS_TOKEN")
    return
  }

  whiteListUrl := fmt.Sprintf("https://%s/api/v3/orgs/%s/members", os.Getenv("GHE_SERVER"), os.Getenv("GHE_WHITE_LIST_ORG"))
  userUrl := fmt.Sprintf("https://%s/api/v3/users", os.Getenv("GHE_SERVER"))
  pat := os.Getenv("GHE_PERSONAL_ACCESS_TOKEN")

  whiteListUsers  := getUserMap(whiteListUrl, pat)
  allUsers  := getUserList(userUrl, pat)

  suspendCounter := 0
  for _, user := range allUsers {
    if user.Type == "User" {
      _, ok := whiteListUsers[user.Id]
      if !ok {
        suspend(user, pat)
        suspendCounter ++
      }
    }
  }

  log.Printf("Suspended %d users\n", suspendCounter)
}

func suspend(user User, pat string) {
  var url = fmt.Sprintf("https://octodemo.com/api/v3/users/%s/suspended", user.Login)

  log.Printf("Suspend %s with %s\n", user.Login, url)

  req, err := http.NewRequest("PUT", url, nil)
  if err != nil {
    log.Fatal("Failed while building the HTTP client: ", err)
    return
  }

  // Provide authentication
  req.Header.Add("Authorization", fmt.Sprintf("token %s", pat))
  req.Header.Add("Content-Length", "0")

  client := &http.Client{}

  resp, doErr := client.Do(req)
  if doErr != nil {
    log.Fatal("Failed executing HTTP request: ", doErr)
    return
  }

  // Close when method returns
  defer resp.Body.Close()
}

/*
 * Get the users as a map using Id as a key
 */
func getUserMap(url, pat string) map[int]User {
  userArray := getUserList(url, pat)
  var userMap map[int]User = make(map[int]User)

  for _, user := range userArray {
    // Filtering out orgs
    if user.Type == "User" {
      userMap[user.Id] = user
    }
  }

  return userMap
}

/*
 * Get the users (including orgs) in an array
 */
func getUserList(url, pat string) []User {
  log.Printf("Retrieving users with %s", url)
  // Build the request
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatal("Failed while building the HTTP client: ", err)
    return nil
  }

  // Provide authentication
  req.Header.Add("Authorization", fmt.Sprintf("token %s", pat))

  client := &http.Client{}

  resp, doErr := client.Do(req)
  if doErr != nil {
    log.Fatal("Failed executing HTTP request: ", doErr)
    return nil
  }

  // Close when method returns
  defer resp.Body.Close()

  var userArray []User

  // Decode the JSON array
  decodeErr := json.NewDecoder(resp.Body).Decode(&userArray)
  if decodeErr != nil {
  	log.Fatal(decodeErr)
    return nil
  } else {
    // Working the pagination as described in https://developer.github.com/guides/traversing-with-pagination/
    // linkHeader := <https://octodemo.com/api/v3/users?since=35>; rel="next", <https://octodemo.com/api/v3/users{?since}>; rel="first"
    linkHeader := resp.Header.Get("Link")

    if linkHeader != "" {
      linkArray := linkListRegex.FindAllStringSubmatch(linkHeader, -1)
      /*
      linkArray := [["<https://octodemo.com/api/v3/users?since=35>; rel=next", "https://octodemo.com/api/v3/users?since=35", "next"],
       ["<https://octodemo.com/api/v3/users{?since}>; rel="first", "https://octodemo.com/api/v3/users{?since}, "first"]]
      */
      for _, linkElement := range linkArray {
        if linkElement[2] == "next" {
          // Getting the next page of users and appending to the current array
          userArray = append(userArray, getUserList(linkElement[1], pat)...)
          break
        }
	    }
    }
    return userArray
  }
}
