# Gator

Gator is a simple command line application for aggregating RSS feeds.

It is built using Go, and uses a Postgres database to store feed/post data.

See [boot.dev - Build a Blog Aggregator](https://www.boot.dev/lessons/14b7179b-ced3-4141-9fa5-e67dbc3e5242)


## Requirements

- Go
- Postgres


## Initialisation

DB needs to be created in Postgres:
```
CREATE DATABASE gator;
```

Apply the DB schema with goose:
```
# install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# go into sql/schema dir and run goose against the upversion defs
# postgres_url is the connection string, of the form:
#   protocol://username:password@host:port/database
cd sql/schema
goose postgres "[postgres_url]" up
```

The application expects a JSON configuration file in your HOME directory, `~/.gatorconfig.json`. This tells the application where to connect to the DB (and it is also used to store the active user across executions).

```
{
    "db_url" : "[postgres_url]" ,
}
```

## Usage

The application can be installed and run anywhere:
```
go install
gator [command]
```

Alternatively you can build and run directly from the app dir:
```
go build
./gator [command]
```

Recognised commands:

`register [user_name]`: Create a new user account, and set it to active.

`login [user_name]`: Switch active user account.

`users`: List registered users.

`reset`: Remove all data.

`addfeed [name] [url]`: Subscribe to a new RSS feed with given name and URL. The will also follow the specified feed with the active user.

`feeds`: List all subscribed feeds.

`follow [url]`: Follow the (previously subscribed) feed at URL with the active user.

`unfollow [url]`: Unfollow the feed at URL with the active user.

`following`: List all followed feeds for the active user.

`agg [interval]`: Run the aggregation service. This will poll the most out of date feed at the given interval (e.g. `1m0s` for every minute). This will run indefinitely until you terminate the process (via `Ctrl+C` or whatever)

`browse [num]`: Browse `num` most recent posts; will default to 2 if no `num` specified.
