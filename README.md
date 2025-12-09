# Gator CLI

## Overview:
This is a demo Golang CLI tool which will fetch RSS feeds for users and display them in the terminal. This tool uses Postgresql for a database backend. with SQLC for generating type-safe SQL code and Gator for managing the database migrations. 

## Installation:
* Use `docker-compose.yml` file at the root of the repository to generate a local environment. Remember to change the password to something secure!
* Add suitable connection string into a `.gatorconfig.json` file in your home directory with the following syntax (including suitable connection string):
```
{"db_url":"postgres_connection_string","current_user_name":"kahya"}
```
* To install the Gator CLI run: `go install github.com/wmag19/gator`

## Commands:

### User Management
- `register <name>` - Create a new user and set as current
- `login <name>` - Switch to an existing user
- `users` - List all users (current user marked)
- `reset` - Delete all users from database

### Feed Management
- `addfeed <name> <url>` - Add a new RSS feed and automatically follow it
- `feeds` - List all feeds with their URLs and creators
- `follow <url>` - Follow an existing feed
- `following` - List all feeds you're following
- `unfollow <url>` - Unfollow a feed

### Reading Posts
- `browse [limit]` - View recent posts from followed feeds (default limit: 2)
- `agg <duration>` - Start feed aggregator to fetch new posts at specified intervals (e.g., "30s", "1m", "1h")