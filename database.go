package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func DatabaseInit(db *sql.DB) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = ""
	cfg.Net = "tcp"
	cfg.Addr = "localhost:3306"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS recordings;")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
		return nil, err
	}

	cfg.DBName = "recordings"

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected to MySQL database successfully")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS album (id INT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255), artist VARCHAR(255), price DECIMAL(5,2));")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
		return nil, err
	}
	return db, nil
}
