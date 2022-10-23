package main

import (
	"context"
	"day9/connection"
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

	connection.Connect_DB()

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
	Id               int
	Title            string
	start_date       time.Time
	end_date         time.Time
	Format_startdate string
	Format_enddate   string
	Deskripsi        string
	Teknologi        []string
	Gambar           string
	Penulis          string
	Durasi           string
}

var Blog = []ValueBlog{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	query := "SELECT kd_project, nm_project, start_date, end_date, deskripsi, teknologi, gambar FROM project "
	rows, _ := connection.Conn.Query(context.Background(), query)

	// query := "SELECT kd_project, nm_project, start_date, end_date, deskripsi, teknologi, gambar FROM project"
	// rows, _ := connection.Conn.Query(context.Background(), query)

	var result []ValueBlog
	fmt.Println(result)
	for rows.Next() {
		var each = ValueBlog{}
		var err = rows.Scan(&each.Id, &each.Title, &each.start_date, &each.end_date, &each.Deskripsi, &each.Teknologi, &each.Gambar)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Penulis = "Siapa Aja Dah"
		each.Durasi = "2 menit yang lalu"

		each.Format_startdate = each.start_date.Format("2 January 2022")
		each.Format_enddate = each.end_date.Format("2 January 2022")

		result = append(result, each) //untuk memasukkan data yang ada di dalam tabel database ke dalam array

	}
	fmt.Println(result)

	respData := map[string]interface{}{
		"Project": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
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

	// var inputTitle = r.PostForm.Get("namaproject")
	// var inputStartdate = r.PostForm.Get("startdate")
	// var inputEnddate = r.PostForm.Get("enddate")
	// var inputDeskripsi = r.PostForm.Get("deskripsi")
	// // var inputReact = r.PostForm.Get("react")
	// // var inputNode = r.PostForm.Get("node")
	// // var inputNext = r.PostForm.Get("next")
	// var inputTehnik = r.PostForm.Get("typescript")
	// var inputGambar = r.PostForm.Get("gambar")
	// // var inputDurasi = time.Now().String()
	// var inputPenulis = "Dumbways Batch 41"

	newBlog := ValueBlog{
		// Title:     inputTitle,
		// StartDate: inputStartdate,
		// EndDate:   inputEnddate,
		// Deskripsi: inputDeskripsi,
		// tehnik:    inputTehnik,
		// Gambar:    inputGambar,
		// // Durasi:     inputDurasi,
		// Penulis: inputPenulis,
	}
	Blog = append(Blog, newBlog)
	http.Redirect(w, r, "/Project", http.StatusMovedPermanently)
}

func blog(w http.ResponseWriter, r *http.Request) {

	//bagian tampilan
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}
	//ini bagian database//
	query := "SELECT kd_project, nm_project, start_date, end_date, deskripsi, teknologi, gambar FROM project"
	rows, _ := connection.Conn.Query(context.Background(), query)

	var result []ValueBlog //mendeklarasikan array dari ValueBlog

	for rows.Next() {
		var each = ValueBlog{}

		var err = rows.Scan(&each.Id, &each.Title, &each.start_date, &each.end_date, &each.Deskripsi, &each.Teknologi, &each.Gambar)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Penulis = "Penulisnya aku"
		// each.Gambar = "https://media.sproutsocial.com/uploads/2017/02/10x-featured-social-media-image-size.png"
		each.Durasi = "1 menit yang lalu"

		each.Format_startdate = each.start_date.Format("2 January 2006")
		each.Format_enddate = each.end_date.Format("2 January 2006")

		result = append(result, each)
	}

	fmt.Println(result)

	respData := map[string]interface{}{
		"Blog": result,
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

	query := "SELECT kd_project, nm_project, start_date, end_date, deskripsi, image FROM project"
	rows, _ := connection.Conn.Query(context.Background(), query)

	var DataBlog = ValueBlog{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	for rows.Next() {
		if id == DataBlog.Id {
			DataBlog = ValueBlog{
				Title:      DataBlog.Title,
				start_date: DataBlog.start_date,
				end_date:   DataBlog.end_date,
				Deskripsi:  DataBlog.Deskripsi,
				Gambar:     DataBlog.Gambar,
				Durasi:     DataBlog.Durasi,
				Penulis:    DataBlog.Penulis,
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
