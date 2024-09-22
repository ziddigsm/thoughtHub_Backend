CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,        
    mail VARCHAR(255) UNIQUE NOT NULL,  
    name VARCHAR(100) NOT NULL,         
    username VARCHAR(20),               
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_on TIMESTAMP DEFAULT NOW(),
    updated_on TIMESTAMP DEFAULT NOW() 
);