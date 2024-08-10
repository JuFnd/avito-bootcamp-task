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

-- Вставка тестовых данных в таблицу House
INSERT INTO House (address, year_built, developer, last_apartment_added) VALUES
('123 Main St', 2000, 'Builder A', '2024-08-01 10:00:00'),
('456 Elm St', 2010, 'Builder B', '2024-08-02 11:00:00'),
('789 Oak St', 2015, NULL, '2024-08-03 12:00:00');

-- Вставка тестовых данных в таблицу Apartment
INSERT INTO Apartment (apartment_number, price, rooms, house_id) VALUES
(101, 150000, 3, 1),
(102, 200000, 4, 1),
(201, 250000, 2, 2),
(202, 300000, 3, 2),
(301, 180000, 2, 3),
(302, 220000, 3, 3);

-- Вставка тестовых данных в таблицу Apartment_House
INSERT INTO Apartment_House (apartment_id, house_id, status) VALUES
(1, 1, 'approved'),
(2, 1, 'on moderation'),
(3, 2, 'approved'),
(4, 2, 'declined'),
(5, 3, 'approved'),
(6, 3, 'created');
