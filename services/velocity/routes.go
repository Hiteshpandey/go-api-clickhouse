package velocity

import (
	"encoding/json"
	"net/http"
	"velocityApi/logs"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Handler struct {
	clickhouseConn *driver.Conn
	logger         *logs.ApiLogger
}

func NewHandler(conn *driver.Conn, logger *logs.ApiLogger) *Handler {
	return &Handler{clickhouseConn: conn, logger: logger}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /seeds", h.listSeeds)
	// router.HandleFunc("GET /seeds/{email}", h.fetchSeed)
}

func (h *Handler) listSeeds(w http.ResponseWriter, r *http.Request) {
	exampleData := make(map[string][]string)
	exampleData["data"] = []string{"seed1@gmail.com", "seed2@yahoo.com", "seed3@outlook.com"}
	j, err := json.Marshal(exampleData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// func (h *Handler) fetchSeed(w http.ResponseWriter, r *http.Request) {
// 	email := r.PathValue("email")
// 	exampleData := make(map[string])
// }
