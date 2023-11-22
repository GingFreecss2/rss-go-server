package api

// Import the necessary packages for encoding/decoding XML,
// input/output operations, making HTTP requests, and time-related functions.

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// Define the RSSFeed struct to represent an RSS feed.
type RSSFeed struct {
	// Channel represents the root element of an RSS feed.
	Channel struct {
		// Title is the title of the RSS feed.
		Title string `xml:"title"`
		// Link is the URL of the RSS feed.
		Link string `xml:"link"`
		// Description provides a brief overview of the RSS feed.
		Description string `xml:"description"`
		// Language specifies the language of the RSS feed.
		Language string `xml:"language"`
		// Item represents a single RSS item.
		Item []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// Define the RSSItem struct to represent a single item in an RSS feed.
type RSSItem struct {
	// Title is the title of the RSS item.
	Title string `xml:"title"`
	// Link is the URL of the RSS item.
	Link string `xml:"link"`
	// PubDate indicates the publication date of the RSS item.
	PubDate string `xml:"pubDate"`
	// Description provides a summary of the RSS item.
	Description string `xml:"description"`
}

// UrlToFeed parses an RSS feed from a given URL
// and returns an RSSFeed struct containing the parsed data.
func UrlToFeed(url string) (*RSSFeed, error) {
	// Create an HTTP client with a timeout of 10 seconds.
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	// Send an HTTP GET request to the specified URL.
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	// Defer closing the response body to avoid resource leaks.
	defer resp.Body.Close()

	// Read the entire response body into a byte slice.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Initialize an empty RSSFeed struct.
	rssFeed := RSSFeed{}

	// Unmarshal the XML data into the RSSFeed struct using the xml.Unmarshal function.
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, err
	}

	// Return the parsed RSSFeed struct or an error if parsing failed.
	return &rssFeed, nil
}
