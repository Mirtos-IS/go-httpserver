package main

import (
	"encoding/json"
	"fmt"
	"httpserver/models"
	"httpserver/views"
	"log"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)


var count int
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func viewHandler(w http.ResponseWriter, r *http.Request) {
    rawId := r.URL.Path[len("/view/"):]
    id, err := strconv.ParseInt(rawId, 10, 64)
    if err != nil {
        fmt.Println(err)
    }
    user, err := models.LoadUser(id)
    if err != nil {
        http.Redirect(w, r, "/edit/" + rawId, http.StatusFound)
        return
    }
    views.View(user).Render(r.Context(), w)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    rawId := r.URL.Path[len("/edit/"):]
    id, err := strconv.ParseInt(rawId, 10, 64)
    if err != nil {
        fmt.Println(err)
    }

    user, err := models.LoadUser(id)
    if err != nil {
        user = &models.User{}
    }
    views.Edit(user).Render(r.Context(), w)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    views.Login().Render(r.Context(), w)
}

func testRHandler(w http.ResponseWriter, r *http.Request) {
    jData, _ := json.Marshal(count)
    count++
    w.Header().Set("Content-Type", "application/json")
    w.Write(jData)
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
    rawId := r.URL.Path[len("/users/"):]

    var id int
    if rawId == "" {
        id = 0
    } else {
        id, _ = strconv.Atoi(rawId)
    }

    users, _ := models.GetUsers(id)

    views.Users(users).Render(r.Context(), w)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    rawId := r.URL.Path[len("/users/get/"):]

    var id int
    if rawId == "" {
        id = 0
    } else {
        id, _ = strconv.Atoi(rawId)
    }

    users, _ := models.GetUsers(id)
    jData, _ := json.Marshal(users)

    w.Header().Set("Content-Type", "application/json")
    w.Write(jData)
}

func loginCheckHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")


    user, _ := models.LoginUser(username, password)
    if user.Uid > 0 {
        views.View(user).Render(r.Context(), w)
        return
    }
    user, _ = models.FindUserByPassword(password)
    views.Wrong(user).Render(r.Context(), w)
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
        user, _ := models.LoadUser(id)
        user.Username = username
        user.Password = password
        user.Save()
    } else {
        user := &models.User{Username: username, Password: password}
        id, err = user.Save()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

    }
    http.Redirect(w, r, "/view/" + strconv.FormatInt(id, 10), http.StatusFound)
}

func main() {
    http.HandleFunc("/save/", saveHandler)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/testR", testRHandler)
    http.HandleFunc("/users/", UsersHandler)
    http.HandleFunc("/users/get", getUsersHandler)
    http.HandleFunc("/login/check", loginCheckHandler)
    //load CSS
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
