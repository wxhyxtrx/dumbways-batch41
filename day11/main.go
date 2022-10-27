package main

import (
	"context"
	"day11/connection"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	route := mux.NewRouter()

	connection.Connect_DB()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("Get")

	route.HandleFunc("/Project", project).Methods("Get")
	route.HandleFunc("/Project", addProject).Methods("POST")
	route.HandleFunc("/Project/{id}", deleteBlog).Methods("Get")

	route.HandleFunc("/Blog", blog).Methods("Get")
	route.HandleFunc("/Detail/{id}", detailBlog).Methods("Get")

	route.HandleFunc("/Delete/{id}", deleteBlog).Methods("Get")

	route.HandleFunc("/Update/{id}", pageEdit).Methods("Get")
	route.HandleFunc("/Update-Blog/{id}", editBlog).Methods("POST")

	route.HandleFunc("/Contact", contact).Methods("Get")

	route.HandleFunc("/register", pageRegister).Methods("Get")
	route.HandleFunc("/registered", register).Methods("Post")

	route.HandleFunc("/login", pageLogin).Methods("Get")
	route.HandleFunc("/log-in", login).Methods("POST")

	route.HandleFunc("/Logout", logout).Methods("Get")

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
	IdUser           string
	postat           time.Time
	Durasi           string
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var Data = MetaData{}

type MetaData struct {
	Title     string
	isLogin   bool
	Username  string
	FlashData string
}

// var Blog = []ValueBlog{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// SESSION
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["isLogin"] != true {
		Data.isLogin = false
	} else {
		Data.isLogin = session.Values["isLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}
	fm := session.Flashes("message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	fmt.Println(flashes)
	Data.FlashData = strings.Join(flashes, "")

	query := "SELECT kd_project, nm_project, start_date, end_date, deskripsi, teknologi, gambar,post_at FROM project "
	rows, _ := connection.Conn.Query(context.Background(), query)

	var result []ValueBlog

	for rows.Next() {
		var each = ValueBlog{}
		var err = rows.Scan(&each.Id, &each.Title, &each.start_date, &each.end_date, &each.Deskripsi, &each.Teknologi, &each.Gambar, &each.postat)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.IdUser = "Siapa Aja Dah"
		each.Durasi = each.postat.Format("2 January 2006")

		each.Format_startdate = each.start_date.Format("2 January 2006")
		each.Format_enddate = each.end_date.Format("2 January 2006")

		result = append(result, each) //untuk memasukkan data yang ada di dalam tabel database ke dalam array

	}
	// fmt.Println(result)

	respData := map[string]interface{}{
		"Project": result,
		"Data":    Data,
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

// bagian project
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
	var inputTehnik = []string{r.PostForm.Get("react"), r.PostForm.Get("node"), r.PostForm.Get("next"), r.PostForm.Get("typescript")}
	var inputGambar = r.PostForm.Get("gambar")
	_, err = connection.Conn.Exec(context.Background(), "insert into project(nm_project,start_date,end_date,deskripsi,teknologi, gambar) values ($1,$2,$3,$4,$5,$6)", inputTitle, inputStartdate, inputEnddate, inputDeskripsi, inputTehnik, inputGambar)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}
	http.Redirect(w, r, "/Project", http.StatusMovedPermanently)
}

// bagian blog
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
		each.IdUser = "Penulisnya aku"
		// each.Gambar = "https://media.sproutsocial.com/uploads/2017/02/10x-featured-social-media-image-size.png"
		each.Durasi = "1 menit yang lalu"

		each.Format_startdate = each.start_date.Format("2 January 2006")
		each.Format_enddate = each.end_date.Format("2 January 2006")

		result = append(result, each)
	}

	// fmt.Println(result)

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

	var DataBlog = ValueBlog{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(),
		"SELECT  nm_project, start_date, end_date, deskripsi, teknologi, gambar, post_at FROM project WHERE kd_project=$1", id).Scan(&DataBlog.Title, &DataBlog.start_date, &DataBlog.end_date, &DataBlog.Deskripsi, &DataBlog.Teknologi, &DataBlog.Gambar, &DataBlog.postat)

	// fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	DataBlog.Format_startdate = DataBlog.start_date.Format("2 January 2006")
	DataBlog.Format_enddate = DataBlog.end_date.Format("2 January 2006")
	DataBlog.Durasi = DataBlog.postat.Format("2 January 2006")

	DataBlog.IdUser = "Wahyu Tricahyo"

	fmt.Println(DataBlog.Deskripsi)

	data := map[string]interface{}{
		"BlogView": DataBlog,
	}

	fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}
func deleteBlog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM project WHERE kd_project=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}
	http.Redirect(w, r, "/Blog", http.StatusFound)
}

// bagian update
func pageEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/update-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	} else {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		updateData := ValueBlog{}

		err = connection.Conn.QueryRow(context.Background(), "SELECT kd_project, nm_project, start_date, end_date, deskripsi, teknologi, gambar, post_at FROM project WHERE kd_project=$1", id).Scan(&updateData.Id, &updateData.Title, &updateData.start_date, &updateData.end_date, &updateData.Deskripsi, &updateData.Teknologi, &updateData.Gambar, &updateData.postat)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Message :" + err.Error()))
			return
		}
		updateData = ValueBlog{
			Id:               updateData.Id,
			Title:            updateData.Title,
			Deskripsi:        updateData.Deskripsi,
			Format_startdate: updateData.start_date.Format("2006-01-02"),
			Format_enddate:   updateData.end_date.Format("2006-01-02"),
			Teknologi:        updateData.Teknologi,
			Gambar:           updateData.Gambar,
		}
		respData := map[string]interface{}{
			"DataEdit": updateData,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)
	}
}
func editBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	inputTitle := r.PostForm.Get("namaproject")
	inputStartdate := r.PostForm.Get("startdate")
	inputEnddate := r.PostForm.Get("enddate")
	inputDeskripsi := r.PostForm.Get("deskripsi")
	inputTehnik := []string{r.PostForm.Get("react"), r.PostForm.Get("node"), r.PostForm.Get("next"), r.PostForm.Get("typescript")}
	inputGambar := r.PostForm.Get("gambar")

	// UPDATE PROJECT TO POSTGRESQL
	_, err = connection.Conn.Exec(context.Background(), `UPDATE project
		SET "nm_project"=$1, "start_date"=$2, "end_date"=$3, "deskripsi"=$4, "teknologi"=$5, "gambar"=$6
		WHERE "kd_project"=$7`, inputTitle, inputStartdate, inputEnddate, inputDeskripsi, inputTehnik, inputGambar, id)
	// ERROR HANDLING INSERT PROJECT TO POSTGRESQL
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
func pageRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}
func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	usernamae := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	passwordhash, _ := bcrypt.GenerateFromPassword([]byte(password), 4)
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(username, email, password) VALUES($1,$2,$3)", usernamae, email, passwordhash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}
func pageLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/login.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fm := session.Flashes("message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")
	respData := map[string]interface{}{
		"Login": Data,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var email = r.PostForm.Get("email")
	var password = r.PostForm.Get("password")

	var dataUser = User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(
		&dataUser.Id, &dataUser.Name, &dataUser.Email, &dataUser.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(password))
	if err != nil {
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		session.AddFlash("Password Salah", "message")
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Values["isLogin"] = true
	session.Values["Name"] = dataUser.Name
	session.Options.MaxAge = 10800

	session.AddFlash("Succesfully Login!", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
