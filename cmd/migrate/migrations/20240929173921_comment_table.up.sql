CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,                         
    user_id INTEGER NOT NULL,                      
    blog_id INTEGER NOT NULL,                      
    is_active BOOLEAN NOT NULL DEFAULT TRUE,       
    comment VARCHAR(100) NOT NULL,                 
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_comment
        FOREIGN KEY (user_id)
        REFERENCES users(id)                       
        ON DELETE CASCADE,                         
    CONSTRAINT fk_blog_comment
        FOREIGN KEY (blog_id)
        REFERENCES blogs(id)                       
        ON DELETE CASCADE                          
);
