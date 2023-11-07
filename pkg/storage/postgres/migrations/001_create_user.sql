-- +goose Up
CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,
  email TEXT NOT NULL,
  name TEXT NOT NULL
);
