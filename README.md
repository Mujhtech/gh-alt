# GH ALT

A GitHub alternative built with Go for the backend and React + Vite for the frontend.

## Features

- **User Authentication**: Register and login functionality with JWT tokens
- **Repository Management**: Create and view repositories
- **Commit History**: Track and view commits for each repository
- **File Browser**: View files and their content in repositories
- **Modern UI**: Dark-themed GitHub-inspired interface

## Tech Stack

### Backend
- Go (Golang)
- Gorilla Mux (HTTP router)
- SQLite (Database)
- JWT (Authentication)
- bcrypt (Password hashing)

### Frontend
- React
- Vite
- React Router (Routing)
- Axios (HTTP client)

## Getting Started

### Prerequisites
- Go 1.24 or higher
- Node.js 18 or higher
- npm or yarn

### Backend Setup

1. Navigate to the project root:
```bash
cd gh-alt
```

2. Install Go dependencies:
```bash
go mod tidy
```

3. Build and run the backend:
```bash
go run backend/cmd/main.go
```

The backend server will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Repositories
- `GET /api/repositories` - Get all public repositories
- `POST /api/repositories` - Create a new repository (requires auth)
- `GET /api/repositories/:id` - Get repository details
- `GET /api/repositories/:id/files` - Get repository files
- `GET /api/repositories/:id/commits` - Get repository commits
- `POST /api/repositories/:id/commit` - Create a new commit (requires auth)

## Project Structure

```
gh-alt/
├── backend/
│   ├── cmd/
│   │   └── main.go              # Main entry point
│   ├── internal/
│   │   ├── api/
│   │   │   └── handlers.go      # HTTP handlers
│   │   ├── database/
│   │   │   └── database.go      # Database initialization
│   │   ├── middleware/
│   │   │   └── middleware.go    # CORS and Auth middleware
│   │   └── models/
│   │       └── models.go        # Data models
│   └── pkg/
│       └── auth/
│           └── auth.go          # Authentication utilities
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   └── Navbar.jsx       # Navigation bar
│   │   ├── context/
│   │   │   └── AuthContext.jsx  # Auth context provider
│   │   ├── pages/
│   │   │   ├── Home.jsx         # Repository list page
│   │   │   ├── Login.jsx        # Login page
│   │   │   ├── Register.jsx     # Register page
│   │   │   └── Repository.jsx   # Repository detail page
│   │   ├── services/
│   │   │   └── api.js           # API client
│   │   ├── App.jsx              # Main app component
│   │   └── main.jsx             # Entry point
│   └── package.json
├── go.mod
├── go.sum
└── README.md
```

## Development

### Running Backend Tests
```bash
go test ./...
```

### Building for Production

#### Backend
```bash
go build -o server ./backend/cmd/main.go
./server
```

#### Frontend
```bash
cd frontend
npm run build
```

The built files will be in `frontend/dist/`

## License

MIT License
