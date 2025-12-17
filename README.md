# GH ALT

A GitHub alternative built with Go for the backend and React + Vite for the frontend.

## Features

- **User Authentication**: Register and login functionality with JWT tokens
- **Repository Management**: Create and view repositories
- **Commit History**: Track and view commits for each repository
- **File Browser**: View files and their content in repositories
- **Git Client Support**: Full Git protocol support for push/pull/clone operations
- **Modern UI**: Dark-themed GitHub-inspired interface

## Tech Stack

### Backend
- Go (Golang)
- Gorilla Mux (HTTP router)
- SQLite (Database)
- JWT (Authentication)
- bcrypt (Password hashing)
- go-git (Git operations)
- Git HTTP Smart Protocol

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

### Environment Variables

#### Backend
Create a `.env` file or set these environment variables:
- `JWT_SECRET`: Secret key for JWT token signing (default: dev-secret-key-change-in-production)
- `CORS_ORIGIN`: Allowed CORS origin (default: * for development, set to specific domain in production)
- `DATABASE_PATH`: Path to SQLite database file (default: ./gh-alt.db)
- `REPOS_DIR`: Directory to store Git repositories (default: ./repositories)

Example:
```bash
export JWT_SECRET="your-super-secret-jwt-key"
export CORS_ORIGIN="http://localhost:5173"
export DATABASE_PATH="./data/gh-alt.db"
```

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

## Using Git Client

Once you've created a repository through the web interface, you can use standard Git commands to interact with it:

### Clone a Repository
```bash
git clone http://localhost:8080/git/my-repo
```

### Push to a Repository
```bash
cd my-project
git init
git add .
git commit -m "Initial commit"
git remote add origin http://localhost:8080/git/my-repo
git push -u origin master
```

### Pull from a Repository
```bash
git pull origin master
```

**Note**: The repository name in the Git URL should match the repository name created in the web UI.

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

### Git Protocol
- `GET /git/:repo/info/refs` - Git info/refs (smart HTTP)
- `POST /git/:repo/git-upload-pack` - Git upload pack (clone/fetch)
- `POST /git/:repo/git-receive-pack` - Git receive pack (push)

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
