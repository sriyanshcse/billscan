package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	billscan "github.com/okcredit/billscan/api/go"
	"net/http"
	"sync"
)

func New(api billscan.APIServer) http.Handler {
	return &router{api: api}
}

var _ http.Handler = &router{}

type router struct {
	api        billscan.APIServer
	distinctID string
	chiRouter  chi.Router
	once       sync.Once
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.once.Do(func() {

		// handler
		h := &handlers{api: router.api}

		// setup chi router
		r := chi.NewRouter()

		// global middlewares
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(middleware.StripSlashes)
		r.Use(middleware.RealIP)
		r.Use(middleware.RequestLogger(NewFilteredLogger([]string{"/"})))

		// routes
		r.Method(http.MethodGet, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// api v1
		r.Route("/v1", func(router chi.Router) {
			router.Method(http.MethodGet, "/users/{user_id}/contacts", h.ListContacts())
			router.Method(http.MethodGet, "/users/{user_id}/contacts/{contact_id}", h.GetContact())
			router.Method(http.MethodPost, "/users/{user_id}/contacts", h.AddContact())
			router.Method(http.MethodPut, "/users/{user_id}/contacts/{contact_id}", h.UpdateContact())
			router.Method(http.MethodDelete, "/users/{user_id}/contacts/{contact_id}", h.DeleteContact())
			router.Method(http.MethodGet, "/users/{user_id}", h.GetUser())
			router.Method(http.MethodGet, "/users/find/{mobile}", h.GetUserByMobile())
			router.Method(http.MethodPost, "/users", h.CreateUser())
			router.Method(http.MethodPut, "/users/{user_id}", h.UpdateUser())
			router.Method(http.MethodDelete, "/users/{user_id}", h.DeleteUser())
			router.Method(http.MethodGet, "/bills/{bill_id}", h.GetBill())
			router.Method(http.MethodPost, "/bills", h.CreateBill())
			router.Method(http.MethodPut, "/bills/{bill_id}", h.UpdateBill())
			router.Method(http.MethodDelete, "/bills/{bill_id}", h.DeleteBill())
			router.Method(http.MethodGet, "/users/{user_id}/contacts/{contact_id}/bills", h.ListBills())
			router.Method(http.MethodPost, "/login", h.Login())
		})

		router.chiRouter = r
	})

	router.chiRouter.ServeHTTP(w, r)
}
