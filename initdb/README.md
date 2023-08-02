# Import CSV into database

Imports people from a given CSV file into database by making POST requests to the titanic API. The input is validated for correct types. There is no deduplication for existing people.

## Requirements

- Go installed
- This directory should be on your GOPATH

## Usage

```
go run initdb.go ./titanic.csv <titanic-host> <titanic-port>
```
