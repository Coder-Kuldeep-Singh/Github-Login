package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://github.com/"
)

var (
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
)

// App to store our http.Client
type App struct {
	Client *http.Client
}

// AuthenticityToken to store authenticity_token value
type AuthenticityToken struct {
	Token string
}

// Project to store the list of repositories
// scraped from gitlab account.
type Project struct {
	Name string
}

// This function will scrape the value of hidden input
// authenticity_token from gitlab signin page.
func (app *App) getToken() AuthenticityToken {
	loginURL := baseURL + "login"
	client := app.Client
	response, err := client.Get(loginURL)
	CheckError("Error fetching response", err)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	CheckError("Error loading HTTP response body", err)

	token, _ := document.Find("input[name='authenticity_token']").Attr("value")
	authenticityToken := AuthenticityToken{
		Token: token,
	}
	return authenticityToken
}

// This function will login to the website using
// the credentials username, password and authenticity_token.
func (app *App) login() {
	client := app.Client
	authenticityToken := app.getToken()
	loginURL := baseURL + "login"

	data := url.Values{
		"authenticity_token": {authenticityToken.Token},
		"login":              {username},
		"password":           {password},
	}
	response, err := client.PostForm(loginURL, data)
	CheckError("Failed while sending data for login", err)
	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)
	CheckError("Error to take response of the page", err)
}

func (app *App) getProjects() []Project {

	projectsURL := baseURL + "mavensingh"
	client := app.Client

	response, err := client.Get(projectsURL)
	CheckError("Error fetching response", err)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	CheckError("Error loading HTTP response body", err)

	var projects []Project

	document.Find(".day").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("data-count")
		project := Project{
			Name: name,
		}
		projects = append(projects, project)
	})
	return projects
}

func CheckError(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// we are creating a cookiejar to store cookies
// required while logging into the website
func main() {
	jar, _ := cookiejar.New(nil)
	app := App{
		Client: &http.Client{Jar: jar},
	}
	app.login()
	projects := app.getProjects()
	for _, project := range projects {
		fmt.Printf("Commits --> %s\n", project.Name)
	}
}
