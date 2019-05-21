package truapi

import (
	"net/http"

	"github.com/TruStory/truchain/x/truapi/render"
)

// HandleStatistics dumps metrics per user per day basis (NOT 'accumulated till date' basis as in HandleMetrics).
func (ta *TruAPI) HandleStatistics(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		render.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	address := r.FormValue("address")
	if address == "" {
		http.Error(w, "must provide a valid address", http.StatusBadRequest)
		return
	}
	from := r.FormValue("from")
	if from == "" {
		render.Error(w, r, "provide a valid from date", http.StatusBadRequest)
		return
	}
	to := r.FormValue("to")
	if to == "" {
		render.Error(w, r, "provide a valid to date", http.StatusBadRequest)
		return
	}

	userMetrics, err := ta.DBClient.AggregateUserMetricsByAddressBetweenDates(address, from, to)
	if err != nil {
		panic(err)
	}

	render.JSON(w, r, userMetrics, http.StatusOK)

}
