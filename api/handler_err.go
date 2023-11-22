package api

import (
	"net/http"

	"github.com/GingFreecss2/rss-go-server/utils"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 400, "Something went wrong")
}
