package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir(app.staticPath))
	router.Handler(
		http.MethodGet,
		"/static/*filepath",
		http.StripPrefix("/static", fileServer),
	)

	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave)
	dynamicThen := dynamicMiddleware.ThenFunc

	// TODO Improve these, httprouter won't allow confilicting routes /snippets/:id + /snippets/new etc
	router.Handler(http.MethodGet, "/", dynamicThen(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamicThen(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamicThen(app.snippetCreateForm))
	router.Handler(http.MethodPost, "/snippet/create", dynamicThen(app.snippetCreate))

	router.Handler(http.MethodGet, "/user/signup", dynamicThen(app.userSignupForm))
	router.Handler(http.MethodPost, "/user/signup", dynamicThen(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", dynamicThen(app.userLoginForm))
	router.Handler(http.MethodPost, "/user/login", dynamicThen(app.userLogin))
	router.Handler(http.MethodPost, "/user/logout", dynamicThen(app.userLogout))

	requestMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return requestMiddleware.Then(router)
}
