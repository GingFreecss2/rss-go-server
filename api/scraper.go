package api

import (
	// Import necessary packages for database operations, logging, string manipulation, synchronization, and time management.
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	// Import internal database package for interacting with the database.
	"github.com/GingFreecss2/rss-go-server/internal/database"

	// Import UUID package for generating unique identifiers.
	"github.com/google/uuid"
)

// StartScraping periodically fetches and processes RSS feeds using multiple goroutines.
func StartScraping(
	db *database.Queries, // Connection to the database.
	concurrency int, // Number of concurrent goroutines to use for scraping.
	timeBetweenRequest time.Duration, // Time to wait between scraping each RSS feed.
) {
	// Log a message indicating the scraping schedule.
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	// Create a ticker that triggers periodically.
	ticker := time.NewTicker(timeBetweenRequest)

	// Continuously scrape RSS feeds until the program is terminated.
	for ; ; <-ticker.C {
		// Fetch the next batch of feeds to scrape.
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency), // Limit the number of feeds to fetch.
		)

		if err != nil {
			log.Println("Error fetching feeds:", err)
			continue // Skip to the next iteration if there's an error.
		}

		// Create a wait group to track the completion of goroutines.
		wg := &sync.WaitGroup{}

		// Start a goroutine to scrape each feed concurrently.
		for _, feed := range feeds {
			wg.Add(1) // Increment the wait group counter before starting the goroutine.

			go scrapeFeed(db, wg, feed) // Start the goroutine to scrape the feed.
		}

		// Wait for all goroutines to finish before proceeding.
		wg.Wait()
	}
}

// scrapeFeed fetches an RSS feed, extracts relevant data, and stores it in the database.
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	// Defer decrementing the wait group counter to signal the goroutine's completion.
	defer wg.Done()

	// Mark the feed as fetched to prevent duplicate processing.
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)

	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return // Exit the goroutine if there's an error.
	}

	// Fetch the RSS feed using the `UrlToFeed` function from `rss.go`.
	rssFeed, err := UrlToFeed(feed.Url)

	if err != nil {
		log.Println("Error fetching feed:", err)
		return // Exit the goroutine if there's an error.
	}

	// Process each item in the RSS feed.
	for _, item := range rssFeed.Channel.Item {
		// Initialize a NullString to store the description, which may be empty.
		description := sql.NullString{}

		// Set the description value if it's not empty.
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// Parse the publication date from the RSS feed format.
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			log.Printf("Error parsing date %v with err %v", item.PubDate, err)
			continue // Skip the item if the date is invalid.
		}

		// Create a new post record in the database with the extracted information.
		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				continue // Skip the item if it already exists in the database.
			}
			log.Println("error creating post: ", err) // Log an error if the post creation fails.
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
