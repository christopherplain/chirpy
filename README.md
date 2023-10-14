# Chirpy

Chirpy is a Go microblogging web server built as an exploration exercise. This project is not production read at this time. However, feel free to take and use what you can.

## Requirements

All development and testing was done using Go `1.21.1`. Other versions may work but have not been tested.

## Running the App

Run the following command in the terminal to launch the server:

```shell
% go build -o chirpy && ./chirpy
```

Use the "debug" flag to clear the database:

```shell
% go build -o chirpy && ./chirpy --debug
```

The server is configured by default to listen on port 8080.

## Boot.dev

I created Chirpy while working through the coursework at [Boot.dev](https://boot.dev). Boot.dev is a great resource for learning backend development and is worth checking out if you're interested in building similar projects.
