package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/MattP2001/envios-demoras-api/internal/store"
)

// Tolerancia antes de marcar un envio como demorado (supuesto: 24hs).
const tolerancia = 24 * time.Hour

type Handler struct {
	Store *store.Store
}

type eventoRequest struct {
	ShipmentID   string     `json:"shipment_id"`
	Status       string     `json:"status"`
	EventDate    time.Time  `json:"event_date"`
	PromisedDate *time.Time `json:"promised_date"`
}

type envioResponse struct {
	*store.Shipment
	Demorado bool `json:"demorado"`
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /eventos", h.recibirEvento)
	mux.HandleFunc("GET /envios/{id}", h.getEnvio)
	mux.HandleFunc("GET /envios/{id}/historial", h.getHistorial)
	return mux
}

func (h *Handler) recibirEvento(w http.ResponseWriter, r *http.Request) {
	var req eventoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json invalido", http.StatusBadRequest)
		return
	}
	if req.ShipmentID == "" || req.Status == "" {
		http.Error(w, "faltan datos", http.StatusBadRequest)
		return
	}
	if req.EventDate.IsZero() {
		req.EventDate = time.Now()
	}
	if err := h.Store.SaveEvent(req.ShipmentID, req.Status, req.EventDate, req.PromisedDate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getEnvio(w http.ResponseWriter, r *http.Request) {
	sh, err := h.Store.GetShipment(r.PathValue("id"))
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "no existe", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(envioResponse{Shipment: sh, Demorado: EstaDemorado(sh, time.Now())})
}

func (h *Handler) getHistorial(w http.ResponseWriter, r *http.Request) {
	evs, err := h.Store.GetEvents(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(evs)
}

// EstaDemorado: si no se entrego y paso la fecha prometida mas la tolerancia.
func EstaDemorado(sh *store.Shipment, ahora time.Time) bool {
	if sh.Status == "ENTREGADO" || sh.PromisedDate == nil {
		return false
	}
	return ahora.After(sh.PromisedDate.Add(tolerancia))
}
