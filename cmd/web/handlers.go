package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"binai.net/internal/constants"
	"binai.net/internal/models"
	"binai.net/internal/validator"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	firstPage := 1
	pageSize := 5
	companyName := app.companyName
	var input struct {
		Search string
		models.Filters
	}

	// v := validator.New()

	qs := r.URL.Query()

	r.ParseForm()
	search := r.Form.Get("search") // search value
	price := r.Form.Get("price")   // price
	if price == "" {
		price = "0"
	}
	regions := r.Form["regions"] // regions
	regionsStr := strings.Join(regions, ", ")

	startDate := r.Form.Get("startDate") // start date
	endDate := r.Form.Get("endDate")     // end date

	pageSizeStr := r.Form.Get("pageSize") // page_size
	if pageSizeStr == "" {
		pageSizeStr = "5000"
	}
	page_size, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		app.serverError(w, err)
	}

	pageSize = page_size

	input.Search = app.readString(qs, "search", search)
	input.Filters.Page = app.readInt(qs, "page", firstPage)
	input.Filters.PageSize = app.readInt(qs, "page_size", pageSize)
	input.Filters.Price = app.readString(qs, "price", price)
	input.Filters.Regions = app.readString(qs, "regions", regionsStr)
	input.Filters.StartDate = app.readString(qs, "start-date", startDate)
	input.Filters.EndDate = app.readString(qs, "end-date", endDate)

	input.Filters.Sort = app.readString(qs, "sort", "-createdate")

	input.Filters.SortSafelist = []string{
		"id",
		"advert_id",
		"name",
		"status",
		"createdate",
		"organizer",
		"price",
		"url",
		"lottype",
		"startdate",
		"enddate",
		"linkdownloadfile",

		"-id",
		"-advert_id",
		"-name",
		"-status",
		"-createdate",
		"-organizer",
		"-price",
		"-url",
		"-lottype",
		"-startdate",
		"-enddate",
		"-linkdownloadfile",
	}

	lots, metadata, err := app.lots.GetRelevantLotList(companyName, input.Search, input.Filters)

	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}

	fmt.Println(companyName)

	totalRecords, err := app.lots.GetLots("aetz")
	if err != nil {
		log.Print("err", err)
	}

	fmt.Println(metadata)
	data := app.newTemplateData(r)
	data.Lots = lots
	data.CountLots = len(totalRecords)
	data.Metadata = metadata
	data.Userdata.Name = companyName
	data.Regions = constants.Regions
	data.PageSize = constants.PageSize
	data.Title = "АЭТЗ"

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	lot, err := app.lots.GetRelevantLotById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//mockData := web.LoadMockData()
	//lot, _ := web.GetMockDateById(mockData, id)

	data := app.newTemplateData(r)
	data.Lot = lot
	data.Lot.Now = time.Now()
	data.Title = data.Lot.Name

	app.render(w, http.StatusOK, "view.html", data)
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	app.errorLog.Print("Error")

	err := app.decodePostForm(r, &form)
	if err != nil {
		fmt.Println("decode")
		app.errorLog.Printf("Error decoding form: %v, Request path: %s", err, r.URL.Path)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		stringToLog := fmt.Sprintf("%s", err)
		err = app.LogToFile("./logs/user_login/", stringToLog)
		if err != nil {
			fmt.Println(err)
		}
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Это поле не может быть пустым")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Введите валидную почту")
	form.CheckField(validator.NotBlank(form.Password), "password", "Это поле не может быть пустым")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Почта или пароль неправильный")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// web.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	id := app.sessionManager.Get(r.Context(), "authenticatedUserID")
	userId, ok := id.(int)
	if !ok {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	user, err := app.users.UserInfo(userId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data.Userdata.Name = user.Name
	data.Userdata.Email = user.Email

	file, err := os.Open("urls.txt")
	currentDir, _ := os.Getwd()
	if err != nil {
		fmt.Println("Текущая директория:", currentDir)
		fmt.Println("Ошибка urls файла:", err)
	}
	defer file.Close()

	var lots []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lots = append(lots, line)
		}
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	data.Userdata.Lots = lots

	app.render(w, http.StatusOK, "profile.html", data)
}

func (app *application) Statistics(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "statistics.html", data)
}

// func (web *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
// 	data := web.newTemplateData(r)

// 	data.Form = snippetCreateForm{
// 		Expires: 365,
// 	}

// 	web.render(w, http.StatusOK, "create.html", data)
// }

// type snippetCreateForm struct {
// 	Title               string `form:"title"`
// 	Content             string `form:"content"`
// 	Expires             int    `form:"expires"`
// 	validator.Validator `form:"-"`
// }

// func (web *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
// 	var form snippetCreateForm

// 	err := web.decodePostForm(r, &form)
// 	if err != nil {
// 		web.clientError(w, http.StatusBadRequest)
// 		return
// 	}

// 	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
// 	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
// 	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

// 	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

// 	if !form.Valid() {
// 		data := web.newTemplateData(r)
// 		data.Form = form

// 		web.render(w, http.StatusUnprocessableEntity, "create.html", data)
// 		return
// 	}

// 	id, err := web.snippets.Insert(form.Title, form.Content, form.Expires)
// 	if err != nil {
// 		web.serverError(w, err)
// 		return
// 	}

// 	web.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

// 	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
// }
