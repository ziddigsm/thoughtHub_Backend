ThoughtHub Backend
A Go-based backend service for ThoughtHub, a blogging platform with features for post management, commenting, search functionality, and more.

Features
Blog Management: Create, read, update, and delete blog posts
Comment System: Post and manage comments on blogs
Search Functionality: Find blogs by keywords
API Key Authentication: Rotating daily API keys for secure access
Rate Limiting: Prevents API abuse with configurable request limits
User Management: Handle user registration and authentication
CORS Support: Configured for cross-origin requests


Blog Service
POST /create_blog - Create a new blog post

GET /get_blogs - Retrieve all blogs

GET /up_likes - Increment likes on a blog

POST /post_comment - Add a comment to a blog

DELETE /delete_blog_by_id - Delete a specific blog

🔍 Search Service
GET /search_blogs - Search blogs by keywords

👥 Additional Services
User management endpoints

Menu-related endpoints
