# CSV Import Search Service

A Go service designed to import CSV data and provide efficient search capabilities via a RESTful API.

## Features

- Import CSV files into structured data models.
- Expose RESTful endpoints for searching and retrieving data.
- Modular architecture with clear separation of concerns.
- Built with Go for performance and concurrency.

## Project Structure

```
CSV_Import_Search_Service/
├── cmd/ # Entry point of the application
├── internal/ # Core application logic
├── pkg/models/ # Data models and structures
├── go.mod # Go module file
├── go.sum # Go module checksums
└── .gitignore # Git ignore file
```

1. Clone the repository:

   ```bash
   git clone https://github.com/diogoamaralse/CSV_Import_Search_Service.git
   cd CSV_Import_Search_Service
   go run ./cmd/main.go


endpoints:

GET -> /api/v1/user?email=diogo.amaral@gmail.com
POST -> /api/v1/user [multipart] key file in a CSV format