CREATE TABLE mobile_projects ( 
    id SERIAL PRIMARY KEY, 
    name VARCHAR(100) NOT NULL, 
    description TEXT NOT NULL, 
    img TEXT NOT NULL, 
    price DECIMAL(10,2) NOT NULL, 
    time_develop INTEGER NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
); 
