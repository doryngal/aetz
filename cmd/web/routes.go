package main

import (
	"net/http"

	"binai.net/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/statistics", dynamic.ThenFunc(app.Statistics))
	router.Handler(http.MethodGet, "/lot/:id", dynamic.ThenFunc(app.snippetView))
	//router.Handler(http.MethodGet, "/user/profile", protected.ThenFunc(app.Profile))
	//router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(web.snippetCreate))
	// router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(web.snippetCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
