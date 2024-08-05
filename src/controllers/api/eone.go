package api

import (
	"net/http"

	"github.com/vladanan/vezbamo4/src/controllers/clr"
	"github.com/vladanan/vezbamo4/src/models"
)

type EoneHandler struct {
	db models.DB
}

func NewEoneHandler(db models.DB) *EoneHandler {
	return &EoneHandler{db: db}
}

func (h *EoneHandler) HandleGetBilling(w http.ResponseWriter, r *http.Request) error {

	// l := clr.GetELRfunc2()

	data, err := h.db.GetBilling(r)
	if err != nil {
		return err
	}
	if data != nil {
		return clr.WriteJSON(w, 200, data)
	} else {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

	// w.Write(db.GetQuestions())
	// return db.GetQuestions()
}
