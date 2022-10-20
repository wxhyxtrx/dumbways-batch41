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

	route.HandleFunc("/Error", page404).Methods("Get")

	fmt.Println("Sever Running on Port 5000")
	http.ListenAndServe("localhost:5000", route)
}

type ValueBlog struct {
	Title      string
	StartDate  string
	EndDate    string
	Deskripsi  string
	React      string
	Node       string
	Next       string
	Typescript string
	Gambar     string
	Durasi     string
	Penulis    string
}

var Blog = []ValueBlog{
	{
		Title:      "Dumbways Day 8 Week 2",
		StartDate:  "27 November 2022",
		EndDate:    "1 Januari 2022",
		Deskripsi:  "Sedang mempelajari cara penggunaan Golang dalam Structur Interface Manipulatuin",
		React:      "React JS",
		Node:       "Node JS",
		Typescript: "Typescript",
		Next:       "Next JS",
		Gambar:     "https://niagaspace.sgp1.digitaloceanspaces.com/blog/wp-content/uploads/2021/12/08144522/apa-itu-programmer-1.jpg",
		Durasi:     "30 Menit yang lalu",
		Penulis:    "Wahyu Tricahyo | Dumbways Batch 41",
	},
	{
		Title:      "Dumbways Day 8 Week 2 | Testing 1",
		StartDate:  "27 November 1999",
		EndDate:    "27 November 2022",
		Deskripsi:  "Sedang mempelajari cara penggunaan Golang dalam Structur Interface Manipulatuin",
		React:      "React JS",
		Node:       "Node JS",
		Typescript: "Typescript",
		Next:       "Next JS",
		Gambar:     "https://niagaspace.sgp1.digitaloceanspaces.com/blog/wp-content/uploads/2021/12/08144522/apa-itu-programmer-1.jpg",
		Durasi:     "30 Menit yang lalu",
		Penulis:    "Wahyu Tricahyo | Dumbways Batch 41",
	},
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
func page404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/404.html")
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

	respData := map[string]interface{}{
		"Blog": Blog,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}
func detailBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/my-blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	var DataBlog = ValueBlog{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	for index, data := range Blog {
		if id == index {
			DataBlog = ValueBlog{
				Title:      data.Title,
				StartDate:  data.StartDate,
				EndDate:    data.EndDate,
				Deskripsi:  data.Deskripsi,
				React:      data.React,
				Node:       data.Node,
				Next:       data.Next,
				Typescript: data.Typescript,
				Gambar:     data.Gambar,
				Durasi:     data.Durasi,
				Penulis:    data.Penulis,
			}
		}
	}

	data := map[string]interface{}{
		"BlogView": DataBlog,
	}

	fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}
