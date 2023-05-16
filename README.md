# InnoTaxiOrder API

This is a microservice for creating and managing orders.

## Installation

To install and run the service, follow these steps:

Install Go on your system if you haven't already done so.
Clone the repository to your local machine.
Run the following command in the root directory of the project:

    go run ./cmd/main.go

The service should now be running on localhost:8080.

## Code Description

# Project structure

The code for the microservice is organized into several packages:

- cmd/main.go contains the main function for the service.
- internal/app/app.go contains functions which sets up the API routes and starts the server.
- internal/model/ contains the data models for the application. In this case, there is models for GraphQL - graphqlModel.go and main models - model.go.
- internal/repo/elastic/ contains the repository implementation for working with the databases.
- internal/services/ contains the business logic services for the application.
- internal/handlers/ contains the API request handlers for the application. Service provides handlers for creating orders, finding it, cancel and some other features.
- client/ contains implementation of grpc client to user service and driver service.

---

- The file gqlgen.yml contains the description of graphQL.

- Business logic services are located in the internal/services/ package. Service uses repository layer to get data.

### Conclusion

This microservice demonstrates a simple way to create and manage orders using Go. The code is organized into packages, making it easy to maintain and extend. The `OrderService` and `OrderHandler` objects provide the business logic and API endpoints, respectively.