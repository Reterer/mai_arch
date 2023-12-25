CREATE DATABASE IF NOT EXISTS delivery;

USE delivery;

CREATE TABLE Users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  pass_hash VARCHAR(255) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  UNIQUE (username)
);

CREATE TABLE Items (
  id INT AUTO_INCREMENT PRIMARY KEY,
  data TEXT NOT NULL,
  owner_id INT,
  FOREIGN KEY (owner_id) REFERENCES Users(id)
);

CREATE TABLE Deliveries (
  id INT AUTO_INCREMENT PRIMARY KEY,
  from_user_id INT,
  from_addr TEXT NOT NULL,
  to_user_id INT,
  to_addr TEXT NOT NULL,
  status INT NOT NULL,
  creation_date TIMESTAMP NOT NULL,
  FOREIGN KEY (from_user_id) REFERENCES Users(id),
  FOREIGN KEY (to_user_id) REFERENCES Users(id)
);

CREATE TABLE ItemsDeliveries (
  id INT AUTO_INCREMENT PRIMARY KEY,
  item_id INT,
  delivery_id INT,
  FOREIGN KEY (item_id) REFERENCES Items(id),
  FOREIGN KEY (delivery_id) REFERENCES Deliveries(id),
  UNIQUE KEY (item_id, delivery_id)
);