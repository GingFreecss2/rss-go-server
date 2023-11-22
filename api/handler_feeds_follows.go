package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GingFreecss2/rss-go-server/internal/database"
	"github.com/GingFreecss2/rss-go-server/models"
	"github.com/GingFreecss2/rss-go-server/utils"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *ApiConfig) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Couldn't get feed follow: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, models.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *ApiConfig) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}

	utils.RespondWithJSON(w, 200, struct{}{})
}
