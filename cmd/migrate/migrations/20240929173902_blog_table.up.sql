CREATE TABLE IF NOT EXISTS blogs (
    id SERIAL PRIMARY KEY,                        
    user_id INTEGER NOT NULL,                                             
    title VARCHAR(100) NOT NULL,                  
    content TEXT NOT NULL,                        
    is_active BOOLEAN NOT NULL DEFAULT TRUE,      
    blog_image BYTEA,                             
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)                       
        ON DELETE CASCADE                          
);
