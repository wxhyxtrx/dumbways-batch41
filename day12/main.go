package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"day12/connection"
	"day12/middleware"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	route := mux.NewRouter()

	connection.Connect_DB()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/temp-upload/").Handler(http.StripPrefix("/temp-upload/", http.FileServer(http.Dir("./temp-upload/"))))

	route.HandleFunc("/", home).Methods("Get")

	route.HandleFunc("/Project", project).Methods("Get")
	route.HandleFunc("/Project", middleware.UploadFile(addProject)).Methods("POST")
	route.HandleFunc("/Project/{id}", deleteBlog).Methods("Get")

	route.HandleFunc("/Blog", blog).Methods("Get")
	route.HandleFunc("/Detail/{id}", detailBlog).Methods("Get")

	route.HandleFunc("/Delete/{id}", deleteBlog).Methods("Get")

	route.HandleFunc("/Update/{id}", pageEdit).Methods("Get")
	route.HandleFunc("/Update-Blog/{id}", editBlog).Methods("POST")

	route.HandleFunc("/contact", contact).Methods("Get")

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
	IdUser           int
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
	UserId    int
	Title     string
	IsLogin   bool
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

	//==============================================================
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Username"].(string)
		Data.UserId = session.Values["UserID"].(int)
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

	// ===============================================================================
	var query string

	if !Data.IsLogin {
		query = "SELECT * FROM project"
		rows, _ := connection.Conn.Query(context.Background(), query)
		var result []ValueBlog
		for rows.Next() {
			var each = ValueBlog{}
			var err = rows.Scan(&each.Id, &each.Title, &each.start_date, &each.end_date, &each.Deskripsi, &each.Teknologi, &each.Gambar, &each.postat, &each.IdUser)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			each.Format_startdate = each.start_date.Format("2 January 2006")
			each.Format_enddate = each.end_date.Format("2 January 2006")

			result = append(result, each) //untuk memasukkan data yang ada di dalam tabel database ke dalam array

		}
		respData := map[string]interface{}{
			"Project": result,
			"Data":    Data,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)

	} else {
		query = "SELECT kd_project, nm_project, start_date, end_date, deskripsi, teknologi, gambar,post_at FROM project WHERE user_id =$1"
		rows, _ := connection.Conn.Query(context.Background(), query, Data.UserId)

		var result []ValueBlog

		for rows.Next() {
			var each = ValueBlog{}
			var err = rows.Scan(&each.Id, &each.Title, &each.start_date, &each.end_date, &each.Deskripsi, &each.Teknologi, &each.Gambar, &each.postat)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			each.Format_startdate = each.start_date.Format("2 January 2006")
			each.Format_enddate = each.end_date.Format("2 January 2006")

			result = append(result, each) //untuk memasukkan data yang ada di dalam tabel database ke dalam array

		}
		respData := map[string]interface{}{
			"Project": result,
			"Data":    Data,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)
	}
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	// CHECK LOGIN STATUS
	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Username"].(string)
		// Data.UserID = session.Values["UserID"].(int)
	}

	response := map[string]interface{}{
		"Data": Data,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message" + err.Error()))
		return
	}
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.Username = session.Values["Username"].(string)
		Data.IsLogin = session.Values["IsLogin"].(bool)
	}

	respData := map[string]interface{}{
		"Data": Data,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}
func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Username"].(string)
		Data.UserId = session.Values["UserID"].(int)
	}

	var inputTitle = r.PostForm.Get("namaproject")
	var inputStartdate = r.PostForm.Get("startdate")
	var inputEnddate = r.PostForm.Get("enddate")
	var inputDeskripsi = r.PostForm.Get("deskripsi")
	var inputTehnik = []string{r.PostForm.Get("react"), r.PostForm.Get("node"), r.PostForm.Get("next"), r.PostForm.Get("typescript")}

	dataContext := r.Context().Value("dataFile")
	var inputGambar = dataContext.(string)

	_, err = connection.Conn.Exec(context.Background(), "insert into project(nm_project,start_date,end_date,deskripsi,teknologi, gambar, user_id) values ($1,$2,$3,$4,$5,$6,$7)", inputTitle, inputStartdate, inputEnddate, inputDeskripsi, inputTehnik, inputGambar, Data.UserId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}
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
	query := "SELECT * FROM project"
	rows, _ := connection.Conn.Query(context.Background(), query)

	var result []ValueBlog

	for rows.Next() {
		var each = ValueBlog{}

		var err = rows.Scan(&each.Id, &each.Title, &each.start_date, &each.end_date, &each.Deskripsi, &each.Teknologi, &each.Gambar, &each.postat, &each.IdUser)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Format_startdate = each.start_date.Format("2 January 2006")
		each.Format_enddate = each.end_date.Format("2 January 2006")

		result = append(result, each)
	}

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Username"].(string)
		Data.UserId = session.Values["UserID"].(int)
	}

	respData := map[string]interface{}{
		"Blog": result,
		"Data": Data,
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

	fmt.Println(DataBlog.Deskripsi)

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	// CHECK LOGIN STATUS
	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Username"].(string)
		// Data.UserID = session.Values["UserID"].(int)
	}

	respData := map[string]interface{}{
		"BlogView": DataBlog,
		"Data":     Data,
	}

	// fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
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

func pageEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/update-project.html")

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	// CHECK LOGIN STATUS
	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Username"].(string)
		// Data.UserID = session.Values["UserID"].(int)
	}

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
			"Data":     Data,
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

	dataContext := r.Context().Value("dataFile")
	inputGambar := dataContext.(string)

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

//=====================================================

// SESSION

//====================================================

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

	session.Values["IsLogin"] = true
	session.Values["Username"] = dataUser.Name
	session.Values["UserID"] = dataUser.Id
	session.Options.MaxAge = 10800

	session.AddFlash("Succesfully Login!", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Options.MaxAge = -1
	session.Values["IsLogin"] = false
	session.Values["UserID"] = 0
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
