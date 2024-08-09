-- Удаление таблиц, если они существуют
DROP TABLE IF EXISTS Apartment_House;
DROP TABLE IF EXISTS Apartment;
DROP TABLE IF EXISTS House;

-- Создание таблицы House
CREATE TABLE House (
    house_id SERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    year_built INT NOT NULL,
    developer TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_apartment_added TIMESTAMP
);

-- Создание таблицы Apartment
CREATE TABLE Apartment (
    apartment_id SERIAL PRIMARY KEY,
    apartment_number INT,
    price INT CHECK (price >= 0),
    rooms INT CHECK (rooms >= 0),
    house_id INT,
    FOREIGN KEY (house_id) REFERENCES House(house_id) ON DELETE CASCADE
);

-- Создание таблицы Apartment_House
CREATE TABLE Apartment_House (
    apartment_house_id SERIAL PRIMARY KEY,
    apartment_id INT,
    house_id INT,
    status TEXT DEFAULT 'created',
    FOREIGN KEY (apartment_id) REFERENCES Apartment(apartment_id) ON DELETE CASCADE,
    FOREIGN KEY (house_id) REFERENCES House(house_id) ON DELETE CASCADE
);
