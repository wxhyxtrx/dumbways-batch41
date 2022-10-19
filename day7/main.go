package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("Get")

	route.HandleFunc("/Project", project).Methods("Get")
	route.HandleFunc("/Project", addProject).Methods("POST")

	route.HandleFunc("/Blog", blog).Methods("Get")
	route.HandleFunc("/Detail/{id}", detailBlog).Methods("Get")

	route.HandleFunc("/Contact", contact).Methods("Get")

	fmt.Println("Sever Running on Port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conten-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}
func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Project :" + r.PostForm.Get("namaproject"))
	fmt.Println("Start Date :" + r.PostForm.Get("startdate"))
	fmt.Println("End Date :" + r.PostForm.Get("enddate"))
	fmt.Println("Deskripsi :" + r.PostForm.Get("deskripsi"))
	fmt.Println("React :" + r.PostForm.Get("react"))
	fmt.Println("Node :" + r.PostForm.Get("node"))
	fmt.Println("Next :" + r.PostForm.Get("next"))
	fmt.Println("Typescript :" + r.PostForm.Get("typescript"))
	fmt.Println("File :" + r.PostForm.Get("gambar"))

	http.Redirect(w, r, "/Project", http.StatusMovedPermanently)
}

func blog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}
func detailBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/my-blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	data := map[string]interface{}{
		"Title":   "Dumbways Day 7 Week 2",
		"Content": " Typescript Lorem, ipsum dolor sit amet consectetur adipisicing elit. Similique vero nobis ipsa optio, itaque accusantium voluptates adipisci molestiae exercitationem dignissimos, deleniti omnis atque necessitatibus quasi maiores animi at reiciendis blanditiis. ",
		"Id":      id,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}
