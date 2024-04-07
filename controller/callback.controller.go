package Controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	rsp "github.com/budiuno/openfinance/middlewares/response"
	"github.com/budiuno/openfinance/models"
)

func CallbackDisbursementHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, repoDisburse disburseSetter) {

	var request models.CallbackDisbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := repoDisburse.UpdateDisbursementStatus(db, request.ReferenceID, request.Status)
	if err != nil {
		fmt.Printf("Error on update disburse status to %s, reference_id %d, error : %v\n", request.Status, request.ReferenceID, err)
	}

	type response struct {
		Message string `json:"message"`
	}

	res := response{
		Message: "OK",
	}

	rsp.RespondWithJSON(w, res, http.StatusOK)
}
