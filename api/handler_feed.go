package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GingFreecss2/rss-go-server/internal/database"
	"github.com/GingFreecss2/rss-go-server/models"
	"github.com/GingFreecss2/rss-go-server/utils"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Couldn't create a feed: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, models.DatabaseFeedToFeed(feed))
}

func (apiCfg *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, models.DatabaseFeedsToFeeds(feeds))
}
