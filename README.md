# Web Forum Project

This project is a web forum that allows users to communicate, create posts, comment, like/dislike posts, and filter content. It is built using SQLite for database management and Docker for containerization. The forum supports user authentication, session management, and various filtering mechanisms.

---

## Features

1. **User Authentication**:
   - Users can register with an email, username, and password.
   - Passwords are encrypted before being stored in the database.
   - Users can log in and maintain a session using cookies with an expiration date.
   - UUIDs are used for session management (Bonus).

2. **Communication**:
   - Registered users can create posts and associate them with one or more categories.
   - Users can comment on posts.
   - Posts and comments are visible to all users (registered and non-registered).

3. **Likes and Dislikes**:
   - Registered users can like or dislike posts and comments.
   - The number of likes and dislikes is visible to all users.

4. **Filtering**:
   - Users can filter posts by:
     - Categories (subforums).
     - Posts created by the logged-in user.
     - Posts liked by the logged-in user.

5. **Database**:
   - SQLite is used to store data (users, posts, comments, categories, likes, dislikes).
   - The database is structured using an Entity-Relationship Diagram (ERD).

6. **Docker**:
   - The application is containerized using Docker for easy deployment and testing.

---

## Technologies Used

- **Backend**: Go (Golang)
- **Database**: SQLite
- **Authentication**: Cookies, UUID (Bonus), bcrypt for password encryption
- **Containerization**: Docker
- **Testing**: Unit tests (recommended)

---

## Setup Instructions

### Prerequisites

1. Install [Docker](https://docs.docker.com/get-docker/).
2. Install [Go](https://golang.org/doc/install) (if testing locally without Docker).

### Steps to Run the Project

1. **Clone the Repository**:
```bash
    git clone https://github.com/A-fethi/Forum.git
    cd Forum
```

2. **Build and Run the Docker Container**:
```bash
    docker build -t web-forum .
    docker run -p 8080:8080 web-forum
```

3. **Access the Forum**:

- Open your browser and go to http://localhost:8080.

4. **Database Initialization**:

- The SQLite database will be initialized automatically when the application runs for the first time.

- Sample data can be inserted using the INSERT queries provided in the code.

---

## Database Schema

The database consists of the following tables:

1. Users:

- id, email, username, password (encrypted), created_at

2. Posts:

- id, user_id, title, content, created_at

3. Comments:

- id, post_id, user_id, content, created_at

4. Categories:

- id, name

5. Post_Categories:

- post_id, category_id

6. Likes:

- id, user_id, post_id, comment_id, type (like/dislike)

---

## API Endpoints

1. Authentication:

- POST /register - Register a new user.

- POST /login - Log in and create a session.

- POST /logout - Log out and clear the session.

2. Posts:

- GET /posts - Get all posts (filterable by category, user, or liked posts).

- POST /posts - Create a new post (requires authentication).

- GET /posts/{id} - Get a specific post with comments.

3. Comments:

- POST /posts/{id}/comments - Add a comment to a post (requires authentication).

4. Likes/Dislikes:

- POST /posts/{id}/like - Like a post (requires authentication).

- POST /posts/{id}/dislike - Dislike a post (requires authentication).

- POST /comments/{id}/like - Like a comment (requires authentication).

- POST /comments/{id}/dislike - Dislike a comment (requires authentication).

---

## Error Handling

- The application handles HTTP errors (e.g., 404, 500) and displays user-friendly messages.

- Database errors (e.g., duplicate email) are caught and returned as JSON responses.

---

## Bonus Features

1. Password Encryption:

- Passwords are encrypted using bcrypt.

2. UUID for Sessions:

- UUIDs are used to manage user sessions securely.

3. Docker Optimization:

- Multi-stage Docker builds are used to reduce the final image size.

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.