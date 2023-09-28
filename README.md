# Chirpy

Chirpy is a web server built as a guided learning exercise. Feel free to take what you can, but keep in mind this is far from production ready code.

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

Chirpy is a guided learning exercise I created while working through [Boot.dev](https://boot.dev) coursework. Boot.dev is a great resource for learning backend development. Please check out their courses if you're at all interested in learning to build similar projects.
