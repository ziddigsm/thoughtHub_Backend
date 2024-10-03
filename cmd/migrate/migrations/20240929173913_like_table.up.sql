CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,                        
    likes INTEGER NOT NULL DEFAULT 0,             
    blog_id INTEGER NOT NULL,                     
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_blog
        FOREIGN KEY (blog_id)
        REFERENCES blogs(id)                       
        ON DELETE CASCADE                          
);
