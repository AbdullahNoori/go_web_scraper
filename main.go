
package main


// package main and import the required library.
import {

	"fmt"
	"io/ioutil"
	"strings"
	"log"
	"net/http"
	// "net/http/cookiejar"
	"net/url"

}

go get -u github.com/gocolly/colly/...


// global constant baseURL to store base url 
// of the website and variables username and 
// password to store gitlab username and password
const (
	baseURL = "https://gitlab.com"
)

var (
	username = "your gitlab username"
	password = "your gitlab password"
)


// Struct App to store our http. Client, AuthenticityToken
//  to store authenticity_token value and Project to store 
// the list of repositories scraped from git
type App struct {
	Client *http.Client
}

type AuthenticityToken struct {
	Token string
}

type Project struct {
	Name string
}

// function to storing the value of the token 
// in AuthenticityToken struct and returning it.
func (app *App) getToken() AuthenticityToken {
	loginURL := baseURL + "/users/sign_in"
	client := app.Client

	response, err := client.Get(loginURL)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	token, _ := document.Find("input[name='authenticity_token']").Attr("value")

	authenticityToken := AuthenticityToken{
		Token: token,
	}

	return authenticityToken
}

// function will login to the website using the 
// credentials username, password and authenticity_token

func (app *App) login() {
	client := app.Client

	authenticityToken := app.getToken()

	loginURL := baseURL + "/users/sign_in"

	data := url.Values{
		"authenticity_token": {authenticityToken.Token},
		"user[login]":        {username},
		"user[password]":     {password},
	}

	response, err := client.PostForm(loginURL, data)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

// scrape list of projects and store in projects
//  variable and at the end we are looping through 
// the projects array and printing our project name to the console
func main() {
	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	app.login()
	projects := app.getProjects()

	for index, project := range projects {
		fmt.Printf("%d: %s\n", index+1, project.Name)
	}
}