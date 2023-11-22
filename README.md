# RSS Aggregator

RSS Aggregator is a backend server written in Golang that aggregates blog posts from RSS feeds. The server allows users to add different RSS feeds to its database and it will automatically collect all of the posts for those feeds, download them, and save them in the database for later viewing.

## Project Structure

The project is structured as follows:

- `main.go`: The entry point of the application. It sets up the server, connects to the database, and starts the scraping process.
- `api/`: Contains the handlers and middleware for the API endpoints.
- `internal/`: Contains the internal packages for the application, including the database and authentication logic.
- `models/`: Contains the data models for the application.
- `sql/`: Contains the SQL queries and schema for the database.
- `utils/`: Contains utility functions for the application.
- `go.mod` and `go.sum`: The Go module files that manage the project's dependencies.

## Higher Level Architecture

The following diagram provides a high-level overview of the project's architecture:

![Higher Level Architecture](/architecture.png)

In this diagram:

- The `main` package contains the `main.go` file, which is the entry point of the application.
- The `api` package contains the handler files (`handler_*.go`), the `middleware_auth.go`, `rss.go`, and `scraper.go` files, which are responsible for handling the API requests and scraping RSS feeds.
- The `internal` package contains the `auth` and `database` packages, which are responsible for handling authentication and interacting with the database.
- The `Database` is where the feeds, posts, users, and feed follow relationships are stored.
- The `main.go` file calls the handler files and the `rss.go` and `scraper.go` files to handle the API requests and scrape the RSS feeds. The handler files and the `rss.go` and `scraper.go` files then interact with the `Database`.


## Getting Started

To get started with the project, clone the repository and navigate to the project directory:

`bash git clone https://github.com/GingFreecss2/rss-go-server.git cd rss-go-server`


Then, set up the environment variables in a `.env` file:

`bash PORT=8080 DB_URL=postgres://username:password@localhost:5432/database`


Finally, start the server:

`bash go run main.go`


The server will start and begin scraping the RSS feeds.

## API Endpoints

The API provides several endpoints for managing users, feeds, and feed follows:

- `GET /v1/error`: Endpoint to trigger an error for testing purposes.
- `GET /v1/feed_follows`: Retrieve feed follow relationships (requires authentication).
- `GET /v1/posts`: Retrieve posts for a user (requires authentication).
- `POST /v1/feeds`: Create a new feed (requires authentication).
- `POST /v1/feed_follows`: Create a new feed follow relationship (requires authentication).
- `DELETE /v1/feed_follows/{feedFollowID}`: Delete a feed follow relationship (requires authentication).

For more details on each endpoint and their request/response formats, refer to the API code.

## Contributing

Contributions to the project are welcome. Please feel free to open an issue or submit a pull request.