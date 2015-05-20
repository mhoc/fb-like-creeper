package main

import (
  "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// See the readme for steps on generating this access token
	// Essentially you have to create a facebook app then request a debug access token using the graph API explorer in the context of that app
	// Then make sure to grant that debug access token all of the permissions you absolutely can
	DEBUG_ACCESS_TOKEN = ""
)

// Merges two maps
func Merge(m1 map[string]int, m2 map[string]int) map[string]int {
  Combined := make(map[string]int)
  for key, value := range m1 {
    Combined[key] = value
  }
  for key, value := range m2 {
    if _, in := Combined[key]; in {
      Combined[key] += value
    } else {
      Combined[key] = value
    }
  }
  return Combined
}

// Gets and parses the likes of an actual post
func GetAndParseLikes(postId string, overrideURL string) map[string]int {

  // Build the url
  var url string
  if overrideURL == "" {
    url = fmt.Sprintf("https://graph.facebook.com/%v/likes?access_token=%v", postId, DEBUG_ACCESS_TOKEN)
  } else {
    url = overrideURL
  }

  // Execute the request
  res1, _ := http.Get(url)
  res2, _ := ioutil.ReadAll(res1.Body)

  // Convert it into a map
  var postLikes map[string]interface{}
  json.Unmarshal(res2, &postLikes)

  // Check to make sure that there are actually likes on this page
  if len(postLikes["data"].([]interface{})) == 0 {
    return make(map[string]int)
  }

  // Tally up the likes for this post
  talliedLikes := make(map[string]int)
  for _, like := range postLikes["data"].([]interface{}) {
    name := like.(map[string]interface{})["name"].(string)
    if _, in := talliedLikes[name]; in {
      talliedLikes[name] += 1
    } else {
      talliedLikes[name] = 1
    }
  }

  // If there is pagination, recurse and do this again
  if _, in := postLikes["paging"].(map[string]interface{})["next"]; in {
    nextPage := GetAndParseLikes("", postLikes["paging"].(map[string]interface{})["next"].(string))
    return Merge(talliedLikes, nextPage)
  } else {
    return talliedLikes
  }

}

// Gets and parses all of the posts on a user's newsfeed
func GetAndParseFeed(overrideURL string) map[string]int {

  // Build the url
  var url string
  if overrideURL == "" {
    url = fmt.Sprintf("https://graph.facebook.com/me/feed?access_token=%v", DEBUG_ACCESS_TOKEN)
  } else {
    url = overrideURL
  }

  // Execute the request
  res1, _ := http.Get(url)
  res2, _ := ioutil.ReadAll(res1.Body)

  // Convert to a map
  var feedPage map[string]interface{}
  json.Unmarshal(res2, &feedPage)

  // Check to make sure there are posts on this page
  if len(feedPage["data"].([]interface{})) == 0 {
    return make(map[string]int)
  }

  // Iterate over each post on this page and merge into a master list of likes
  totalLikes := make(map[string]int)
  for _, post := range feedPage["data"].([]interface{}) {

    // Log it
    if _, in := post.(map[string]interface{})["message"]; in {
      fmt.Printf("Getting likes for post \"%v\"\n", post.(map[string]interface{})["message"])
    } else {
      fmt.Printf("Getting likes for post \"%v\"\n", post.(map[string]interface{})["story"])
    }

    // Run the likes
    postId := post.(map[string]interface{})["id"].(string)
    postLikes := GetAndParseLikes(postId, "")
    totalLikes = Merge(totalLikes, postLikes)
  }

  // Handle pagination
  if _, in := feedPage["paging"].(map[string]interface{})["next"]; in {
    nextPage := GetAndParseFeed(feedPage["paging"].(map[string]interface{})["next"].(string))
    return Merge(totalLikes, nextPage)
  } else {
    return totalLikes
  }

}

// Prints the largest value from the map then returns the map minus that value
// This is slower than sorting the map but eh
func PrintLargest(likeMap map[string]int) map[string]int {
  largestKey := ""
  largestValue := -1
  for key, value := range likeMap {
    if value > largestValue {
      largestKey = key
      largestValue = value
    }
  }
  fmt.Printf("%v (%v)\n", largestKey, largestValue)
  delete(likeMap, largestKey)
  return likeMap
}

func main() {

  if DEBUG_ACCESS_TOKEN == "" {
    fmt.Println("You need to get a debug access token and put it in the program. Go read the readme.")
    return
  }

  likeMap := GetAndParseFeed("")

  // Iterate over the map to find the largest value,
  // then print and delete it. This is faster than sorting and more memory efficient.
  fmt.Printf("\n\n")
  for len(likeMap) > 0 {
    likeMap = PrintLargest(likeMap)
  }

}
