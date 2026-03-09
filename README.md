# LLM Dashboard

A full-stack application providing a chat interface powered by a local Large Language Model (LLM) using [Ollama](https://github.com/ollama/ollama). The application allows users to interact with the LLM through a chat-like interface with support for code execution, web search, and background job scheduling.

## Features

- **Chat Interface**: Streamed responses from local LLM models
- **Tool Integration**:
  - Python code execution
  - Web search capabilities
- **User Authentication**: JWT-based auth with registration and login
- **Job Scheduling**: Background task execution for LLM queries
- **File Management**: Upload and manage files for context
- **Rate Limiting**: Protected API endpoints

## Tech Stack

### Backend
- Go
- SQLite (Database)
- Ollama (LLM Integration)
- golang-migrate (Database Migrations)
- Charmbracelet/log (Logging)

### Frontend
- Vue.js
- Pinia (State Management)
- Vite (Build Tool)
- Marked (Markdown Renderer)
- Highlight.js (Syntax Highlighting)
- Vue-Toastify (Notifications)

### Services
- Python Flask (Code Execution Service)

## API Endpoints

### Authentication
- `POST /api/register` - Register new user
- `POST /api/login` - User login
- `POST /api/logout` - User logout
- `POST /api/refresh` - Refresh JWT token
- `GET /api/me` - Get current user

### Chat
- `POST /api/chat` - Send chat message
- `GET /api/stream` - Stream chat responses
- `GET /api/messages` - Get chat history
- `DELETE /api/deleteMessages` - Delete chat messages

### Models
- `GET /api/modelList` - List available models
- `GET /api/modelInfo` - Get model information

### Jobs
- `GET /api/jobList` - List all jobs
- `POST /api/createJob` - Create new job
- `POST /api/updateJob` - Update job
- `POST /api/deleteJob` - Delete job
- `GET /api/viewJob` - Get job details

### Files
- `POST /api/uploadFile` - Upload file
- `GET /api/getFileList` - List files
- `GET /api/getFile` - Get file content
- `POST /api/deleteFile` - Delete file

## Project Structure

```
.
├── backend/
│   ├── api/          # API clients (Ollama, tools)
│   ├── db/           # Database layer
│   ├── jobs/         # Background job processing
│   └── server/       # HTTP server and handlers
├── services/        # Python services (code execution)
└── README.md
```
