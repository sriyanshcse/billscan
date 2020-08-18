package http

import (
	"github.com/go-chi/chi"
	billscan "github.com/okcredit/billscan/api/go"
	"github.com/okcredit/go-common/encoding/json"
	"github.com/okcredit/go-common/errors"
	"github.com/okcredit/go-common/httpx"
	"net/http"
)

type handlers struct {
	api billscan.APIServer
}

func (h *handlers) AddContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.AddContactRequest{UserId: chi.URLParam(r, "user_id")}

		if req.UserId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.api.AddContact(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) UpdateContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.UpdateContactRequest{UserId: chi.URLParam(r, "user_id")}

		if req.UserId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.api.UpdateContact(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) DeleteContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.DeleteContactRequest{UserId: chi.URLParam(r, "user_id"), ContactId: chi.URLParam(r, "contact_id")}

		if req.UserId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
			return
		}

		if req.ContactId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid Contact id"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.api.DeleteContact(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) GetContact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.GetContactRequest{UserId: chi.URLParam(r, "user_id"), ContactId: chi.URLParam(r, "contact_id")}

		if req.UserId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
			return
		}

		if req.ContactId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid Contact id"))
			return
		}

		res, err := h.api.GetContact(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) GetBill() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.GetBillRequest{UserId: chi.URLParam(r, "user_id"), BillId: chi.URLParam(r, "bill_id")}

		//if req.UserId == "" {
		//	_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
		//	return
		//}
		//
		//if req.BillId == "" {
		//	_ = httpx.WriteError(w, errors.From(400, "Invalid bill id"))
		//	return
		//}

		res, err := h.api.GetBill(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) CreateBill() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.CreateBillRequest{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.api.CreateBill(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) UpdateBill() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.UpdateBillRequest{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req.Bill.Id = chi.URLParam(r, "bill_id")
		res, err := h.api.UpdateBill(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) DeleteBill() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.DeleteBillRequest{BillId: chi.URLParam(r, "bill_id")}

		res, err := h.api.DeleteBill(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) ListContacts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.ListContactsRequest{UserId: chi.URLParam(r, "user_id")}

		if req.UserId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
			return
		}

		res, err := h.api.ListContacts(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) ListBills() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.ListBillsRequest{ContactId: chi.URLParam(r, "contact_id"), UserId: chi.URLParam(r, "user_id")}

		if req.UserId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
			return
		}
		if req.ContactId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid contact id"))
			return
		}

		res, err := h.api.ListBills(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) CreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.CreateUserRequest{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.api.CreateUser(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) UpdateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.UpdateUserRequest{UserId: chi.URLParam(r, "user_id")}

		if req.UserId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.api.UpdateUser(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) DeleteUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.DeleteUserRequest{UserId: chi.URLParam(r, "user_id")}

		res, err := h.api.DeleteUser(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.GetUserRequest{UserId: chi.URLParam(r, "user_id")}

		if req.UserId == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid user id"))
			return
		}

		res, err := h.api.GetUser(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) GetUserByMobile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.GetUserByMobileRequest{Mobile: chi.URLParam(r, "mobile")}

		if req.Mobile == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid Mobile"))
			return
		}

		res, err := h.api.GetUserByMobile(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}

func (h *handlers) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &billscan.LoginRequest{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Mobile == "" {
			_ = httpx.WriteError(w, errors.From(400, "Invalid Mobile"))
		}

		res, err := h.api.Login(r.Context(), req)
		if err != nil {
			_ = httpx.WriteError(w, err)
			return
		}

		_ = httpx.WriteJson(w, res)
	})
}
