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