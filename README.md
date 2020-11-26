# MyAnimeList-Go  
[![PkgGoDev](https://pkg.go.dev/badge/github.com/camelva/myanimelist-go)](https://pkg.go.dev/github.com/camelva/myanimelist-go)
[![Build Status](https://travis-ci.com/Camelva/myanimelist-go.svg?branch=main)](https://travis-ci.com/Camelva/myanimelist-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/camelva/myanimelist-go)](https://goreportcard.com/report/github.com/camelva/myanimelist-go)


**MyAnimeList-Go** is a small library to simplify your usage of [MyAnimeList][MAL]' [API]  
  
> According to official API documentation, current version of API labeled as **Beta**. Considering this, **MyAnimeList-Go** can't be called _Stable_ too. But still, at this moment library fully covers all available methods of current API version (v2.0.0).  
  
## Table of contents  
<details>
<summary>click to open</summary> 
  
- [Installation](#installation)  
- [Preparing](#preparing)  
- [Usage](#usage)  
	- [Creating instance](#creating-instance)
	- [Authorization](#authorization)  
		- [Token Expiration](#token-expiration)
		- [Get tokens](#get-tokens)
		- [Set tokens manually](#set-tokens-manually)
	- [Search anime (manga)](#search-anime-manga)
	- [Details about certain anime (manga)](#details-about-certain-anime-manga)
	- [Top anime (manga)](#top-anime-manga)
	- [Seasonal anime](#seasonal-anime)
	- [Anime suggestions](#anime-suggestions)
	- [User information](#user-information)
	- [User anime (manga) list](#user-anime-manga-list)
	- [Update anime (manga) status](#update-anime-manga-status)
	- [Remove entry from list](#remove-entry-from-list)
	- [Forum](#forum)
		- [Forum boards](#forum-boards)
		- [Forum search](#forum-search)
		- [Forum topic information](#forum-topic-information)
	- [Multiple Pages](#multiple-pages)  
	- [Contributing](#contributing)
	- [References](#references)
  
</details>

## Installation
Library was tested to be working at `v1.7+`, but it is recommended to use `v1.11+`
```
import "github.com/camelva/myanimelist-go"  
```
 
## Preparing
Before using API you need to obtain your MyAnimeList client's ID and Secret keys first. It can be achieved by creating new Application. For this, simply head to _Settings => API_. Or just use this [link](https://myanimelist.net/apiconfig). Then, click on **Create ID** button and fill required fields.   
After that, you can find **Client ID** and **Client Secret** fields inside your newly created app. Copy them to safe place - we will need them soon.  
  
![MyAnimeList Application info](../assets/credentials.png?raw=true)  
  
## Usage
### Creating instance
To create new instance you just need to pass `Config` structure to `myanimelist.New()`. Example code below:  
```go  
package main

import (
	"github.com/camelva/myanimelist-go"
	"log"
	"net/http"
	"time"
)

func main() {
	config := myanimelist.Config{
		ClientID: "clientid",
		ClientSecret: "clientsecret",
		RedirectURL: "https://example.com/anime/callback",
		// Optional
		// HTTPClient: *http.Client{Timeout: 5 * time.Second}
		// Logger: *log.Logger{}
	}
	mal, err := myanimelist.New(config)
	if err != nil {
		log.Fatal(err)
	}
	// do stuff
}
```
Here you use **Client ID** and **Client Secret**, obtained on previous step.    
Also make sure you added  your **Redirect URL** to MyAnimeList' application settings, otherwise it will not work.  

_Reference: [New()](https://pkg.go.dev/github.com/camelva/myanimelist-go#New)_

---
### Authorization
Every method of API requires user's **Access Token**, so its good idea to auth as soon as possible.   
MyAnimeList uses OAuth2, so whole process consist of 2 steps: heading user to MyAnimeList's Login page and exchanging received temporaty token for long-term access token.  
  
_Reference: [Auth.LoginURL()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Auth.LoginURL) | [Auth.ExchangeToken()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Auth.ExchangeToken)_ 
  
Example code:   
```go  
package main

import (
	"github.com/camelva/myanimelist-go"
	"log"
	"net/http"
)

// imagine we already made initialization
var mal = &myanimelist.MAL{}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/app", appHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

// Step 1: Heading user to MyAnimeList login page
func loginHandler(w http.ResponseWriter, req *http.Request) {
	// creating authorization URL
	loginPageURL := mal.Auth.LoginURL()
	// and redirecting user there
	http.Redirect(w, req, loginPageURL, http.StatusFound)
}

// Step 2: Exchanging tokens
func callbackHandler(w http.ResponseWriter, req *http.Request) {
	// gathering temporary code from request query
	code := req.FormValue("code")
	// and exchange it for long-lasting access token
	userInfo, err := mal.Auth.ExchangeToken(code)
	if err != nil {
		// handle error
		return
	}
	// optionally you can store tokens in your db or anywhere else
	_, _, _ = userInfo.AccessToken, userInfo.ExpireAt, userInfo.RefreshToken
	http.Redirect(w, req, "/app", http.StatusFound)
}

func appHandler(w http.ResponseWriter, req *http.Request) {
	if mal.Auth.GetTokenInfo().AccessToken == "" {
		http.Redirect(w, req, "/login", http.StatusFound)
	}
	// do some stuff
}
```
  
#### Token expiration
Every user's access tokens have certain time they are valid. Standard, its 1 month (31 day) . You can always check when token will expire by reading `ExpireAt` field of `UserCredentials`.  
If token already expired - you need to ask user to do authorization steps again.
But if token still valid - you can request **token update** and receive new token with fresh duration without even user's interaction.
For this, just call `mal.Auth.RefreshToken()`:  
```go
newCredentials, err := mal.Auth.RefreshToken()  
```

_Reference: [Auth.RefreshToken()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Auth.RefreshToken)_

---
#### Get tokens
You can get your current user's tokens by running:
```go
mal.Auth.GetTokenInfo()
```

_Reference: [Auth.GetTokenInfo()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Auth.GetTokenInfo)_

___
#### Set tokens manually 
If you have your user's tokens saved somewhere and want to continue session without forcing user to log in again - you can set tokens manually. 
```go
mal.Auth.SetTokenInfo(accessToken string, refreshToken string, expireAt time.Time)
```

_Reference: [Auth.SetTokenInfo()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Auth.SetTokenInfo)_

---
### Search anime (manga)
Searching is simple - just use `mal.Anime.Search` or `mal.Manga.Search`with your _query string_ and `PagingSettings` as parameters. These requests are multi-paged, so look at [Multiple pages](#multiple-pages) for additional info.

_Reference:  [Anime.Search()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Anime.Search) | [Manga.Search()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Manga.Search)_

___
### Details about certain anime (manga)
For retrieving detailed info there are `mal.Anime.Details` and `mal.Manga.Details` methods. Both accepts `ID` as first parameter, and, optionally, names of fields to gather. By default, these methods returns `AnimeDetails` (or `MangaDetails`) struct with fields `ID`, `Title` and `MainPicture`. To acquire more fields - you need to explicitly specify them by yourself. You can find list of all _Shared_, _Anime-only_ and _Manga-only_ fields at [Constants](https://pkg.go.dev/github.com/camelva/myanimelist-go#pkg-constants)

_Reference: [Anime.Details()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Anime.Details) | [Manga.Details()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Manga.Details)_

___
### Top anime (manga)
Use `mal.Anime.Top` or `mal.Manga.Top`. First parameter is `RankingType`, second - `PagingSettings` (for more info about paged results see [Multiple pages](#multiple-pages)). There are a couple of different ranks at MyAnimeList, you can find all of them at the official documentation - [Anime ranks](https://myanimelist.net/apiconfig/references/api/v2#operation/anime_ranking_get) and [Manga ranks](https://myanimelist.net/apiconfig/references/api/v2#operation/manga_ranking_get) or at library's documentation [constants section](https://pkg.go.dev/github.com/camelva/myanimelist-go#pkg-constants).

_Reference: [Anime.Top()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Anime.Top) | [Manga.Top()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Manga.Top)_

___
### Seasonal anime
You can get anime list of certain year's season by running 
```go
mal.Anime.Seasonal(year int, season string, sort string, settings PagingSettings)
```
Only `year` and `season` parameters are required, rest are optional. For `season` use one of these constants: `SeasonWinter`, `SeasonSpring`, `SeasonSummer` or `SeasonFall`. 
By passing non-zero `sort` parameter, you can sort result by _score_ or _number of users, added anime to their list_. For this, use either `SortByScore` or `SortByUsersLists` constants.
For additional info about `PagingSettings` see [Multiple pages](#multiple-pages)

_Reference: [Anime.Seasonal()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Anime.Seasonal)_

___
### Anime suggestions
Get anime suggestions for current user by running:
```go
mal.Anime.Suggestions(setting PagingSettings)
```
For additional info about `PagingSettings` see [Multiple pages](#multiple-pages)

_Reference: [Anime.Suggestions()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Anime.Suggestions)_

___
### User information
At the moment, you can acquire information only about current user _(but seems like this API method will support different usernames too)_
```go
mal.User.Info(settings PagingSetting)
```
For additional info about `PagingSettings` see [Multiple pages](#multiple-pages)

_Reference: [User.Info()](https://pkg.go.dev/github.com/camelva/myanimelist-go#User.Info)_

___
### User anime (manga) list
```go
mal.Anime.List.User(username string, status string, sort string, settings PagingSettings)
mal.Manga.List.User(username string, status string, sort string, settings PagingSettings)
```
Here you can use any username to request their anime/manga list. To get current user - pass empty string.
By passing `status` you can filter response to contain only entries with same status. 
You can look at library documentation's [constants section](https://pkg.go.dev/github.com/camelva/myanimelist-go#pkg-constants) to find all _Shared Statuses_, _Anime-only_ and _Manga-only_
Also you can sort result by passing corresponding parameter. List of all available sort constants can also be found at [documentation](https://pkg.go.dev/github.com/camelva/myanimelist-go#pkg-constants) 
For additional info about `PagingSettings` see [Multiple pages](#multiple-pages).

_Reference: [AnimeList.User()](https://pkg.go.dev/github.com/camelva/myanimelist-go#AnimeList.User) | [MangaList.User()](https://pkg.go.dev/github.com/camelva/myanimelist-go#MangaList.User)_

___
### Update anime (manga) status
To add entry to your list or update their statuses you should use corresponding methods: 
```go
mal.Anime.List.Update(config AnimeConfig)
mal.Manga.List.Update(config MangaConfig)
```
_Reference: [AnimeList.Update()](https://pkg.go.dev/github.com/camelva/myanimelist-go#AnimeList.Update) | [MangaList.Update()](https://pkg.go.dev/github.com/camelva/myanimelist-go#MangaList.Update)_

Both of them require appropriate `config` struct. 
It's recommended to create them by running `NewAnimeConfig(id int)` (`NewMangaConfig(id int)`).
These config structures have a bunch of helper methods to set values you want to change. Such as `SetScore`, `SetStatus` and so on. For list of all available methods see documentation reference.

_Reference: [AnimeConfig](https://pkg.go.dev/github.com/camelva/myanimelist-go#AnimeConfig) | [MangaConfig](https://pkg.go.dev/github.com/camelva/myanimelist-go#MangaConfig)_

---
### Remove entry from list
To remove anime from your list use `mal.Anime.List.Remove(id int)` or `mal.Manga.List.Remove(id int)` for manga. 
If there is no entry with such `id` in your list - call still considered as successful and returns no error.

_Reference: [AnimeList.Remove()](https://pkg.go.dev/github.com/camelva/myanimelist-go#AnimeList.Remove) | [MangaList.Remove()](https://pkg.go.dev/github.com/camelva/myanimelist-go#MangaList.Remove)_
 
___
### Forum
#### Forum boards
```go
mal.Forum.Boards()
```
Returns information about all forum categories, boards and sub-boards

_Reference: [Forum.Boards()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Forum.Boards)_

___
#### Forum search
Forum search provided by 
```go
mal.Forum.Search(searchOpts ForumSearchSettings, settings PagingSettings)
```
`ForumSearchSettings` contain all your search parameters. There is no required fields, but you need to fill at least one of them. See documentation for additional info.
For additional info about `PagingSettings` see [Multiple pages](#multiple-pages).

_Reference: [Forum.Search()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Forum.Search) | [ForumSearchSettings](https://pkg.go.dev/github.com/camelva/myanimelist-go#ForumSearchSettings)_

___
#### Forum topic information
To acquire information about certain topic use:
```go
mal.Forum.Topic(id int, settings PagingSettings)
```
For additional info about `PagingSettings` see [Multiple pages](#multiple-pages)

_Reference: [Forum.Topic()](https://pkg.go.dev/github.com/camelva/myanimelist-go#Forum.Topic)_

## Multiple pages
Some requests (such as `Anime.Top()`) returns a lot of data, so result split into multiple pages.  
You can detect such functions by having a special input parameter - struct `PagingSettings` with two fields: `limit` and `offset` (of course you can specify only field you need, or even leave both at zero).   
These functions' returned value (let's call it `PagedResult`) is always struct with, among other things, field `Paging`. This field contains URLs to next and previous pages accordingly. So, for simplified usage, every of such `PagedResult` structures have `Next` and `Prev` methods. Which optionally accepts new `limit` value as parameters.  
  
Example usage:  
```go  
package main

import (
	"github.com/camelva/myanimelist-go"
)

// imagine we already made initialization
var mal = &myanimelist.MAL{}

func main() {
	// lets request 10 top anime, ranked by popularity
	popularAnime, err := mal.Anime.Top(
		myanimelist.RankByPopularity,
		myanimelist.PagingSettings{Limit: 10})
	if err != nil {
		panic(err) // example error handling
	}
    // showed result to user or something else
    // but now we want to get another 10 top anime
    morePopularAnime, err := popularAnime.Next()
    if err != nil {
    	panic(err) // example error handling
    }
    // and now we want more, but only 5
    anotherPopularAnime, err := morePopularAnime.Next(5)
    if err != nil {
    	panic(err) // example error handling
    }
    _ = anotherPopularAnime // do something with result
}  
```  
## Contributing
1.  Fork it (https://github.com/Camelva/myanimelist-go/fork)
2.  Create your feature branch (`git checkout -b feature/fooBar`)
3.  Commit your changes (`git commit -am 'Add some fooBar'`)
4.  Push to the branch (`git push origin feature/fooBar`)
5.  Create a new Pull Request

## References
- [Library Documentation][Doc]  
- MyAnimeList official resources:
	- [MyAnimeList API v2 Documentation][API]
	- [MyAnimeList API Authorization Documentation](https://myanimelist.net/apiconfig/references/authorization)
	- [MyAnimeList API Client Creating Page](https://myanimelist.net/apiconfig)
	- [MyAnimeList API License and Developer Agreement](https://myanimelist.net/static/apiagreement.html)

[MAL]: https://myanimelist.net  
[API]: https://myanimelist.net/apiconfig/references/api/v2  
[Doc]: https://pkg.go.dev/github.com/camelva/myanimelist-go
