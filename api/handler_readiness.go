package api

import (
	"net/http"

	"github.com/GingFreecss2/rss-go-server/utils"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, struct{}{})
}
