package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"binai.net/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		// v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

// Filter фильтрует список лотов по цене
func (app *application) Filter(price string, data []models.Lot) []models.Lot {
	var lots []models.Lot

	// Предварительная обработка строки цены фильтра для удаления пробелов
	filterPrice, err := strconv.ParseInt(price, 10, 64)
	if err != nil {
		// Если произошла ошибка преобразования, возвращаем исходный список без фильтрации
		return data
	}

	for _, el := range data {
		// Предварительная обработка строки цены лота для удаления пробелов

		elPrice, err := parsePrice(el.Price)
		if err != nil {
			fmt.Println("elPrice error parse", err)
			// Если произошла ошибка преобразования, пропускаем данный элемент
			continue
		}

		// Сравниваем числовые значения
		if elPrice >= filterPrice {
			lots = append(lots, el)
		}
	}

	return lots
}

func parsePrice(priceStr string) (int64, error) {
	// Удаляем пробелы
	priceStr = strings.ReplaceAll(priceStr, " ", "")

	// Обрезаем дробную часть
	// Найдем индекс точки
	pointIndex := strings.Index(priceStr, ".")
	if pointIndex != -1 {
		priceStr = priceStr[:pointIndex]
	}

	// Преобразуем строку в целое число
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (app *application) LoadMockData() []models.Lot {
	data, err := ioutil.ReadFile("internal/mock/data.json")
	if err != nil {
		return []models.Lot{}
	}

	var lots []models.Lot
	if err := json.Unmarshal(data, &lots); err != nil {
		return []models.Lot{}
	}

	return lots
}

func (app *application) GetMockDateById(mockData []models.Lot, id int) (*models.Lot, error) {
	var res models.Lot
	for _, el := range mockData {
		if el.ID == id {
			res = el
		}
	}

	return &res, nil
}

func (app *application) LogToFile(path, result string) error {
	// Получаем текущую дату
	now := time.Now()
	// Форматируем дату в формате "год-месяц-день"
	formattedDate := now.Format("2006-01-02")
	// Create the directory if it doesn't exist
	err := os.MkdirAll(path, 0o755) // 0755 - разрешения для папки
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// Create a new file to save the downloaded content
	file, err := os.Create(path + formattedDate)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Write the file content from the response body to the newly created file

	_, err = io.WriteString(file, fmt.Sprint(result))
	if err != nil {
		fmt.Println("Error while writing to file:", err)
		return err
	}
	return nil
}
