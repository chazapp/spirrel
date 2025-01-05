# spirrel

A Golang REST API server that fetches random articles from Wikipedia and
indexes them in an Elasticsearch instance, randomizing some data in the process.
Use it to generate legitimate data in Elasticsearch and play around with Kibana.

## Dev setup

Clone the repository, build the CLI, [start an Elasticsearch instance with Docker](https://www.elastic.co/guide/en/elasticsearch/reference/current/run-elasticsearch-locally.html), run the CLI:

```bash
$ git clone git@github.com:/chazapp/spirrel
$ go build -o spirrel
$ ./spirrel
NAME:
   spirrel - CLI for generating and storing articles in Elasticsearch

USAGE:
   spirrel [global options] command [command options]

COMMANDS:
   run      Run the API server
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

## Refer to Elasticsearch guide to start a local instance in Docker

$ ./spirrel run -s http://localhost:9200 -k $ELASTICSEARCH_API_KEY
```

Once started, the HTTP API exposes 2 routes:

```
$ curl -X POST http://localhost:8080/createIndex # Creates the `articles` ES index 
$ curl -X POST http://localhost:8080/generate?count=5 # Generates $count articles from random Wikipedia entries and stores them in ES
```

## Production deployment

A Docker container is built by this project's Github Action pipeline.
An Helm Chart is also available for deployment on Kubernetes clusters

// TBD: Documentation for Elastic Operator + Create ES users + Accessing Kibana
