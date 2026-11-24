package handlers

import (
	"testing"
	"time"

	"github.com/MattP2001/envios-demoras-api/internal/store"
)

func TestEstaDemorado(t *testing.T) {
	prometida := time.Date(2026, 11, 20, 18, 0, 0, 0, time.UTC)
	sh := &store.Shipment{Status: "EN_TRANSITO", PromisedDate: &prometida}

	if EstaDemorado(sh, prometida.Add(2*time.Hour)) {
		t.Error("no deberia estar demorado dentro de la tolerancia")
	}
	if !EstaDemorado(sh, prometida.Add(30*time.Hour)) {
		t.Error("deberia estar demorado")
	}
	sh.Status = "ENTREGADO"
	if EstaDemorado(sh, prometida.Add(30*time.Hour)) {
		t.Error("entregado no es demorado")
	}
}
