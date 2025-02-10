-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,  -- User's name
    password TEXT NOT NULL,         -- Password
    email TEXT UNIQUE NOT NULL,     -- Email
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- Creation date
);

-- Posts table
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,          -- Username of the user who created the post
    title TEXT NOT NULL,              -- Title of the post
    content TEXT NOT NULL,            -- Content of the post
    image_Content TEXT,
    Categories TEXT,   -- Categories
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    likes INT DEFAULT 0,
    dislikes INT DEFAULT 0,
    FOREIGN KEY (username) REFERENCES users(username)  -- Reference to users table
);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    author STRING NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    likes INTEGER DEFAULT 0,
    dislikes INTEGER DEFAULT 0,
    FOREIGN KEY(post_id) REFERENCES posts(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);

-- User interactions table
CREATE TABLE IF NOT EXISTS user_interactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    item_id INTEGER NOT NULL,
    item_type TEXT NOT NULL,
    action TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, item_id, item_type, action)
);

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    username TEXT NOT NULL UNIQUE,
    session_token TEXT NOT NULL UNIQUE,
    isloggedin BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id)
);