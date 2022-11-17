package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	res, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *res}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	res, err := h.svc.ReadTODO(ctx, int64(req.PrevID), int64(req.Size))
	return &model.ReadTODOResponse{TODOS: res}, err
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	res, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *res}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDS)
	return &model.DeleteTODOResponse{}, err
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var req *model.CreateTODORequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		create, err := h.Create(r.Context(), req)
		if err != nil {
			log.Println(err)
		}
		response, err := json.Marshal(create)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case "PUT":
		var req *model.UpdateTODORequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		update, err := h.Update(r.Context(), req)
		// if reflect.TypeOf(err) == model.ErrNotFound {
		if reflect.DeepEqual(err, new(model.ErrNotFound)) {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			log.Println(err)
		}
		response, err := json.Marshal(update)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case "GET":
		var req *model.ReadTODORequest
		json.NewDecoder(r.Body).Decode(&req)
		read, err := h.Read(r.Context(), req)
		if err != nil {
			log.Println(err)
		}
		response, err := json.Marshal(read)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case "DELETE":
		var req *model.DeleteTODORequest
		json.NewDecoder(r.Body).Decode(&req)
		delete, err := h.Delete(r.Context(), req)
		if err != nil {
			log.Println(err)
		}
		response, err := json.Marshal(delete)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

}
