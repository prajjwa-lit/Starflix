package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_NAME     = "streamflix"
	DB_USER     = "xyz"
	DB_PASSWORD = "xyz"
)

type Video struct {
	ID          int       `json:"id"`
	Filename    string    `json:"filename"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	ReleaseYear int       `json:"release_year,omitempty"`
	CoverImage  string    `json:"cover_image,omitempty"`
	FilePath    string    `json:"path"`
	FileSize    int64     `json:"size"`
	Duration    int       `json:"duration,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var DB *sql.DB

func Initialize() error {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")
	return nil
}
func EnsureTablesExist() error {
	log.Println("Checking and creating database tables if they don't exist...")
	_, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS videos (
            id SERIAL PRIMARY KEY,
            filename VARCHAR(255) NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            genre VARCHAR(100),
            release_year INTEGER,
            cover_image_path VARCHAR(255),
            file_path VARCHAR(255) NOT NULL,
            file_size BIGINT NOT NULL,
            duration INTEGER,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )
    `)
	if err != nil {
		return fmt.Errorf("failed to create videos table: %w", err)
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS genres (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL UNIQUE
        )
    `)
	if err != nil {
		return fmt.Errorf("failed to create genres table: %w", err)
	}
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM genres").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check genres count: %w", err)
	}

	if count == 0 {
		log.Println("Adding default genres...")
		defaultGenres := []string{
			"Action", "Comedy", "Drama", "Documentary",
			"Horror", "Thriller", "Sci-Fi", "Animation",
		}

		for _, genre := range defaultGenres {
			_, err = DB.Exec("INSERT INTO genres (name) VALUES ($1) ON CONFLICT (name) DO NOTHING", genre)
			if err != nil {
				log.Printf("Warning: Failed to insert genre %s: %v", genre, err)
			}
		}
	}

	log.Println("Database tables verified")
	return nil
}
func TestConnection() error {
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM videos").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}

	log.Printf("Successfully connected to database, videos table has %d records", count)
	return nil
}

func GetAllVideos() ([]Video, error) {
	log.Println("Attempting to query all videos from database...")

	query := `
        SELECT id, filename, title, description, genre, release_year, cover_image_path,
        file_path, file_size, duration, created_at, updated_at
        FROM videos
        ORDER BY created_at DESC
    `

	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	log.Println("Successfully queried videos table")

	var videos []Video
	for rows.Next() {
		var v Video
		var releaseYear, duration sql.NullInt32
		var description, genre, coverImage sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&v.ID, &v.Filename, &v.Title, &description, &genre, &releaseYear, &coverImage,
			&v.FilePath, &v.FileSize, &duration, &createdAt, &updatedAt,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		v.Description = description.String
		v.Genre = genre.String
		v.ReleaseYear = int(releaseYear.Int32)
		v.CoverImage = coverImage.String
		v.Duration = int(duration.Int32)
		v.CreatedAt = createdAt
		v.UpdatedAt = updatedAt

		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	log.Printf("Found %d videos in database", len(videos))
	return videos, nil
}
func GetVideosByGenre(genre string) ([]Video, error) {
	query := `
		SELECT id, filename, title, description, genre, release_year, cover_image_path,
		file_path, file_size, duration, created_at, updated_at
		FROM videos
		WHERE genre = $1
		ORDER BY created_at DESC
	`
	rows, err := DB.Query(query, genre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		var releaseYear, duration sql.NullInt32
		var description, genre, coverImage sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&v.ID, &v.Filename, &v.Title, &description, &genre, &releaseYear, &coverImage,
			&v.FilePath, &v.FileSize, &duration, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}

		v.Description = description.String
		v.Genre = genre.String
		v.ReleaseYear = int(releaseYear.Int32)
		v.CoverImage = coverImage.String
		v.Duration = int(duration.Int32)
		v.CreatedAt = createdAt
		v.UpdatedAt = updatedAt

		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}
func GetAllGenres() ([]Genre, error) {
	query := `SELECT id, name FROM genres ORDER BY name`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var g Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func InsertVideo(video *Video) error {
	query := `
		INSERT INTO videos
		(filename, title, description, genre, release_year, cover_image_path, file_path, file_size, duration)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	var releaseYear, duration interface{}
	if video.ReleaseYear > 0 {
		releaseYear = video.ReleaseYear
	} else {
		releaseYear = nil
	}

	if video.Duration > 0 {
		duration = video.Duration
	} else {
		duration = nil
	}

	return DB.QueryRow(
		query,
		video.Filename,
		video.Title,
		video.Description,
		video.Genre,
		releaseYear,
		video.CoverImage,
		video.FilePath,
		video.FileSize,
		duration,
	).Scan(&video.ID, &video.CreatedAt, &video.UpdatedAt)
}
