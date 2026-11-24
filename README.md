# envios-demoras-api

API en Go para recibir los eventos de envíos del operador logístico y consultar el estado y
si un envío está demorado.

## Cómo correrlo

```
go run ./cmd/api
```

Levanta en el puerto 8080 y crea `envios.db` (SQLite).

## Endpoints

- `POST /eventos` — recibe un evento:
  ```json
  {"shipment_id": "123", "status": "EN_TRANSITO", "event_date": "2026-11-24T10:00:00Z", "promised_date": "2026-11-26T18:00:00Z"}
  ```
  La fecha prometida viene en el primer evento; si no viene después, se mantiene la anterior.
- `GET /envios/{id}` — estado actual y `demorado`.
- `GET /envios/{id}/historial` — eventos del envío.

## Estados

CREADO, EN_TRANSITO, EN_DISTRIBUCION, ENTREGADO, CANCELADO.

## Demora

Un envío está demorado si no está entregado y pasaron más de 24 hs de la fecha prometida
(dejé un margen para no generar alertas de más).

## Tests

```
go test ./...
```
