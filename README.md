# Book Tracker API

A RESTful API built in Go to manage a book tracker application, allowing users to add, list, update, and delete books.

## Setup

### Prerequisites
- Go 1.22 or later
- Git

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/mikaelh82/book-tracker.git
   cd book-tracker
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:8080` (or the port specified via the `BACKEND_PORT` environment variable). (TODO: Continue to mention environment variables in the instruction but also in the .env.example so its easy for the user to know what env variables they need to setup)

### Running Tests (TODO: PLEASE BE MORE SPECIFIC HERE LATER)
Run all tests:
```bash
go test ./...
```

## Requirements Checklist

Track progress on the project requirements:

- [ ] **HTTP Server & API Endpoints**
  - [ ] Create a basic HTTP server using `net/http`
  - [ ] Implement `POST /books` to add a new book
  - [ ] Implement `GET /books` to retrieve all books
  - [ ] Implement `PUT /books/{id}` to update a book’s status
  - [ ] Implement `DELETE /books/{id}` to delete a book
- [ ] **Book Data Structure**
  - [ ] Define `Book` struct with `ID`, `Title`, `Author`, `Status`
  - [ ] Use UUIDs for book IDs (`github.com/google/uuid`)
  - [ ] Validate input data
  - [ ] Store data in memory (slice or map) or in an in-process db like duckdb (lets try duckdb as this was mentioned in discussions)
- [ ] **List Filtering and Sorting**
  - [ ] Sort books by title (A-Z) by default
  - [ ] Filter books by status (e.g., `GET /books?status=reading`)
- [ ] **Extra Features (Choose at least one)**
  - [ ] Add pagination (`limit` and `offset` query parameters)
  - [ ] Implement request logging middleware
  - [ ] Support graceful shutdown
  - [ ] Add `/metrics` endpoint for book counts
- [ ] **Testing**
  - [ ] Write unit tests for book validation
  - [ ] Write tests for creating a book (`POST /books`)
  - [ ] Write tests for updating a book’s status (`PUT /books/{id}`)
  - [ ] Write tests for filtering books by status (`GET /books?status=reading`)