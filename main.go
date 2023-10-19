package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
    Title string
    Body []byte
}

var templates = template.Must(template.ParseFiles("html/edit.html", "html/view.html", "html/login.html", "html/wrong.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    rawId := r.URL.Path[len("/view/"):]
    id, err := strconv.ParseInt(rawId, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    user, err := LoadUser(id)
    if err != nil {
        http.Redirect(w, r, "/edit/" + rawId, http.StatusFound)
        return
    }
    renderTemplate(w, "view", user)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    rawId := r.URL.Path[len("/edit/"):]
    id, err := strconv.ParseInt(rawId, 10, 64)
    if err != nil {
        fmt.Println(err)
    }

    user, err := LoadUser(id)
    if err != nil {
        user = &User{}
    }
    renderTemplate(w, "edit", user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "login", nil)
}

func loginCheckHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")


    user, _ := LoginUser(username, password)
    if user.Uid > 0 {
        renderTemplate(w, "view", user)
        return
    }
    user, _ = FindUserByPassword(password)
    renderTemplate(w, "wrong", user)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    rawId := r.URL.Path[len("/edit/"):]
    id, err := strconv.ParseInt(rawId, 10, 64)
    if err != nil {
        fmt.Println(err)
    }

    username := r.FormValue("username")
    password := r.FormValue("password")
    if id > 0 {
        user, _ := LoadUser(id)
        user.Username = username
        user.Password = password
        user.save()
    } else {
        user := &User{Username: username, Password: password}
        id, err = user.save()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

    }
    http.Redirect(w, r, "/view/" + strconv.FormatInt(id, 10), http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, user *User) {
    err := templates.ExecuteTemplate(w, tmpl+".html", user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    http.HandleFunc("/save/", saveHandler)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/login/check", loginCheckHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
