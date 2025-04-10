
# ThoughtHub Backend

A Go-based backend service for ThoughtHub, a blogging platform with features for post management, commenting, search functionality, and more.

---

## 🚀 Features

- **Blog Management**: Create, read, update, and delete blog posts  
- **Comment System**: Post and manage comments on blogs  
- **Search Functionality**: Find blogs by keywords  
- **API Key Authentication**: Rotating daily API keys for secure access  
- **Rate Limiting**: Prevents API abuse with configurable request limits  
- **User Management**: Handle user registration and authentication  
- **CORS Support**: Configured for cross-origin requests  

---

## 📁 Project Structure

```
thoughtHub_Backend/
├── cmd/
│   ├── api/
│   │   └── api.go          # Main API server setup
├── service/
│   ├── blog/
│   │   └── routes.go       # Blog service endpoints
│   ├── search/
│   │   └── route.go        # Search functionality
│   ├── users/              # User management
│   └── menu/               # Menu-related endpoints
├── utils/
│   └── utils.go            # Utility functions and middleware
├── .env                    # Environment variables
└── go.mod                  # Go module dependencies
```

---

## 🛠️ Prerequisites

- Go **1.23.1** or higher  
- **PostgreSQL** database  
- Environment variables configured in `.env`

---

## 📦 Installation

### 1. Clone the repository:

```bash
git clone https://github.com/ziddigsm/thoughtHub_Backend.git
cd thoughtHub_Backend
```

### 2. Install dependencies:

```bash
go mod tidy
```

### 3. Set up `.env`

Configure the environment variables based on the provided example.

### 4. Run the application:

```bash
make run
```

---

## ⚙️ Configuration

### Environment Variables

Set the following in your `.env` file:

```env
DB_HOST=your-database-host
DB_PORT=5432
DB_USER=your-database-username
DB_PASSWORD=your-database-password
DB_NAME=your-database-name
DB_SSLMODE=disable
DB_CONNECTION_STRING=your-database-connection-string

# API keys
API_KEY_0=your-api-key
API_KEY_1=your-api-key
API_KEY_2=your-api-key
API_KEY_3=your-api-key
API_KEY_4=your-api-key
API_KEY_5=your-api-key
API_KEY_6=your-api-key
```

---

## 📡 API Endpoints

> All endpoints are prefixed with `/api/v1`

### 📝 Blog Service

- `POST /create_blog` - Create a new blog post  
- `GET /get_blogs` - Retrieve all blogs  
- `GET /up_likes` - Increment likes on a blog  
- `POST /post_comment` - Add a comment to a blog  
- `DELETE /delete_blog_by_id` - Delete a specific blog  

### 🔍 Search Service

- `GET /search_blogs` - Search blogs by keywords  

### 👥 Additional Services

- User management endpoints  
- Menu-related endpoints  

---

## 🔐 Authentication & Security

### API Key Authentication

Uses a **rotating API key** system based on a logic.  
Include the key in the request header:

```http
X-API-Key: valid-api-key
```

### Rate Limiting

- **Default Rate**: 15 requests per second  
- **Burst Capacity**: 23 requests  

---

## 🌐 CORS Configuration

CORS is enabled for:

- `http://localhost:3000` (development)  
- `https://thoughthub.live` (production)  

---

## 📚 Dependencies

- `github.com/gorilla/mux` - HTTP router and URL matcher  
- `github.com/gorilla/handlers` - CORS and logging middleware  
- `gorm.io/gorm` - ORM library for database operations  
- `golang.org/x/time/rate` - Rate limiting implementation  
- `github.com/joho/godotenv` - Environment variable management  

---

## ⚠️ Error Handling

Standardized error responses in JSON format:

```json
{
  "message": "Error description"
}
```

---

## 🧪 Development Guidelines

- Ensure all endpoint handlers validate API keys and apply rate limiting  
- Follow Go best practices for error handling  
- Use utility functions for consistent response formatting  
- Keep the API design RESTful  

---

## 🔮 Future Improvements

- Enhanced authentication with JWT  
- Improved search with filtering options  
- Caching layer for frequently accessed data  
- Analytics for tracking API usage  

---

© 2025 ThoughtHub Backend
