CREATE TABLE IF NOT EXISTS socials (
    id SERIAL PRIMARY KEY,        
    user_id INTEGER NOT NULL,                                              
    social_media VARCHAR(255) NOT NULL,
    social_url TEXT NOT NULL,                      
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_on TIMESTAMP DEFAULT NOW(),
    updated_on TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_social_media UNIQUE (user_id, social_media) 
);
