# StreamFlix - Personal Video Streaming Platform

StreamFlix is a self-hosted video streaming platform that allows you to manage and stream your personal video collection through a modern web interface. It features a Netflix-inspired UI, real-time video streaming, genre-based organization, and an intuitive upload system.

## 🚀 Features

- **Modern Web Interface**: Netflix-style UI with smooth animations and responsive design
- **Video Streaming**: HTTP range request support for efficient streaming
- **Genre Organization**: Categorize videos by genres
- **Live Upload**: Real-time upload progress with drag-and-drop support
- **Cover Images**: Custom thumbnail support for videos
- **Metadata Management**: Store and display video metadata (title, description, year, etc.)
- **Responsive Design**: Works seamlessly on desktop and mobile devices

## 🛠️ Tech Stack

### Backend (Go)
- **Net/HTTP**: Core HTTP server and routing
- **Database**: PostgreSQL with `lib/pq` driver
- **File System**: Local storage for videos and cover images
- **Middleware**: Custom middleware for logging, CORS, and error handling

### Frontend
- **Pure JavaScript**: No framework dependencies
- **HTML5**: Semantic markup
- **CSS3**: Modern styling with Flexbox and Grid
- **Media Features**: HTML5 video player with custom controls

### Development Tools
- Go 1.23.0
- PostgreSQL
- Modern web browser with HTML5 video support

## 🏗️ System Architecture

### Directory Structure
```plaintext
streamer/
├── api/           # HTTP API handlers
├── config/        # Application configuration
├── db/           # Database operations
├── models/       # Data models
├── services/     # Business logic
├── utils/        # Utility functions
└── server/       # Server setup and static files
```

### Components

1. **API Layer** (`api/`)
   - Routes registration
   - Request handling
   - Response formatting

2. **Service Layer** (`services/`)
   - Video streaming logic
   - File upload handling
   - Metadata management

3. **Database Layer** (`db/`)
   - PostgreSQL connection
   - CRUD operations
   - Schema management

4. **Static Files** (`server/static/`)
   - Frontend application
   - Styles and scripts
   - Media assets

## 🌊 Data Flow

1. **Video Streaming**
```plaintext
Client Request → HTTP Server → Video Service → File System → Streamed Response
```

2. **Video Upload**
```plaintext
File Upload → Upload Service → File System Storage → Database Entry → Response
```

3. **Video Listing**
```plaintext
Client Request → Database Query → Video Service → JSON Response
```

## 📝 Database Schema

### Videos Table
```sql
CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    genre VARCHAR(100),
    release_year INTEGER,
    cover_image_path VARCHAR(255),
    file_path VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    duration INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Genres Table
```sql
CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);
```

## 🚀 Getting Started

### Prerequisites
- Go 1.23.0 or later
- PostgreSQL
- Modern web browser

### Configuration

1. **Database Setup**
```sql
CREATE DATABASE streamflix;
CREATE USER xyz WITH PASSWORD 'xyz';
GRANT ALL PRIVILEGES ON DATABASE streamflix TO xyz;
```

2. **Environment Configuration**
```go
// db/db.go
var (
    DB_HOST     = "localhost"
    DB_PORT     = 5432
    DB_NAME     = "streamflix"
    DB_USER     = "xyz"
    DB_PASSWORD = "xyz"
)
```

### Installation

1. Clone the repository
```bash
git clone https://github.com/DevMaan707/streamer.git
cd streamer
```

2. Install dependencies
```bash
go mod download
```

3. Create necessary directories
```bash
mkdir videos covers
```

4. Build and run
```bash
go build
./streamer -port=5101 -videos=./videos -covers=./covers
```

### Usage

1. Access the web interface at `http://localhost:5101`
2. Upload videos through the upload tab
3. Browse and stream videos through the main interface

## 🔒 Security Considerations

- Implements path traversal protection
- Validates file types and sizes
- Sanitizes filenames
- Uses prepared SQL statements

## 🔧 Configuration Options

- `-port`: HTTP server port (default: 5101)
- `-videos`: Video storage directory (default: ./videos)
- `-covers`: Cover images directory (default: ./covers)
- `-max-upload`: Maximum upload size in MB (default: 1024)

## 📈 Performance Features

- Chunked video streaming
- HTTP range request support
- Database connection pooling
- Efficient file copying
- Client-side caching support

## 🛠️ Development

### Building from Source
```bash
go build -o streamer main.go
```

### Running Tests
```bash
go test ./...
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
