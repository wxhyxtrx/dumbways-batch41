package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

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
	route.HandleFunc("/Delete/{id}", deleteBlog).Methods("Get")

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

	var inputTitle = r.PostForm.Get("namaproject")
	var inputStartdate = r.PostForm.Get("startdate")
	var inputEnddate = r.PostForm.Get("enddate")
	var inputDeskripsi = r.PostForm.Get("deskripsi")
	var inputReact = r.PostForm.Get("react")
	var inputNode = r.PostForm.Get("node")
	var inputNext = r.PostForm.Get("next")
	var inputTypescript = r.PostForm.Get("typescript")
	var inputGambar = r.PostForm.Get("gambar")
	var inputDurasi = time.Now().String()
	var inputPenulis = "Dumbways Batch 41"

	newBlog := ValueBlog{
		Title:      inputTitle,
		StartDate:  inputStartdate,
		EndDate:    inputEnddate,
		Deskripsi:  inputDeskripsi,
		React:      inputReact,
		Node:       inputNode,
		Next:       inputNext,
		Typescript: inputTypescript,
		Gambar:     inputGambar,
		Durasi:     inputDurasi,
		Penulis:    inputPenulis,
	}
	Blog = append(Blog, newBlog)
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
func deleteBlog(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(index)

	Blog = append(Blog[:index], Blog[index+1:]...)
	fmt.Println(Blog)

	http.Redirect(w, r, "/Blog", http.StatusFound)
}
