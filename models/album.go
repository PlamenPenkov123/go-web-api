package models

import (
	"database/sql"
	"fmt"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func GetAlbums(db *sql.DB) ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("failed to query albums: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return albums, nil
}

func GetAlbumsByArtist(artist string, db *sql.DB) ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
	if err != nil {
		return nil, fmt.Errorf("failed to query albums: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return albums, nil
}

func GetAlbumById(id int, db *sql.DB) (Album, error) {
	var album Album
	err := db.QueryRow("SELECT * FROM album WHERE id = ?", id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return Album{}, nil
		}
	}
	return album, nil
}

func AddAlbum(album Album, db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("failed to add album: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}
	return id, nil
}
func UpdateAlbum(id int, newAlbum Album, db *sql.DB) (int64, error) {
	album, err := GetAlbumById(id, db)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve album: %w", err)
	}
	result, err := db.Exec("UPDATE album SET title = ?, artist = ?, price = ? WHERE id = ?", newAlbum.Title, newAlbum.Artist, newAlbum.Price, album.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to update album: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("no album found with ID %d", album.ID)
	}

	return rowsAffected, nil
}

func DeleteAlbum(id int, db *sql.DB) (int64, error) {
	result, err := db.Exec("DELETE FROM album WHERE id = ?", id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete album: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("no album found with ID %d", id)
	}

	return rowsAffected, nil
}
