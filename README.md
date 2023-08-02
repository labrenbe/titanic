# Titanic
Titanic is a service that allows users to query and analyse data from passengers that were aboard the Titanic via a REST API.
The application requires a connection to a postgresql database to store data.

Features include:
- listing all passengers
- showing personal data of a specific passenger
- adding a new passenger
- updating personal data of a passenger
- deleting a person from the list of passengers 

Personal data of passengers that can be accessed:
- name
- sex
- age
- number of siblings or spouses aboard
- number of parents or children aboard
- fare
- passenger class
- survived (true/false)

## API documentation
A detailed documentation on the correct usage of the titanic API can be found in [app/swagger.yml](https://github.com/labrenbe/titanic/blob/main/app/swagger.yml).

## Local deployment with docker-compose

Build the app container
```
cd app
docker build -t titanic .
```
In the root directory run docker-compose
```
docker-compose up -d
```
This will deploy the titanic server and a postgresql that is used to store data. The titanic server will start listening on port 8080.

## Fill database with prepared data
To fill the database with data from a csv file on you local computer run (requires Go installed):
```
go run ./initdb/initdb.go ./initdb/titanic.csv localhost 8080
```

## Deployment to Kubernetes
To run titanic in production you might want to deploy it to kubernetes.
Instructions on how to deploy titanic on kubernetes can be found in [k8s](https://github.com/labrenbe/titanic/tree/main/k8s).

## Areas of improvement
- Use govulncheck instead of trivy
- Integration tests in the pipeline
- Use postgres operator
- Use grafana operator instead of built-in grafana
- Implement OpenTelemetry
- Use loki for logging
