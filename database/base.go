package database

import (
	"crypto/sha256"
	"database/sql"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var salt = os.Getenv("SALT")

type User struct {
	PublicName string `json:"pn"  binding:"required,min=2,max=30"`
	Username   string `json:"un"  binding:"required,min=5,max=20"` //later bind
	Email      string `json:"email"  binding:"required,email"`
	Password   string `json:"password"  binding:"required"` //later bind
}
type Message struct {
	Sender   string
	Reciever string
	Text     string
	Sent     time.Time
}

type Storage struct {
	db *sql.DB
}

func OpenUsers() (*Storage, error) {
	database, err := sql.Open("sqlite", "./DB.sqlite3")
	if err != nil {
		return nil, err
	}
	st, err := database.Prepare(`CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fn TEXT NOT NULL,
    un TEXT UNIQUE NOT NULL,  
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
  	)
	`)
	if err != nil {
		return nil, err
	}
	st.Exec()
	st.Close()
	st2, err := database.Prepare(`CREATE TABLE IF NOT EXISTS messages(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender TEXT REFERENCES users(un),
		reciever TEXT REFERENCES users(un),
		content TEXT NOT NULL ,
		sent TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		return nil, err
	}

	st2.Exec()
	st2.Close()

	return &Storage{db: database}, nil
}

func (s *Storage) AddUser(u *User) error {

	q := "INSERT INTO users (fn,un,email,password) VALUES (?,?,?,?)"
	st, err := s.db.Prepare(q)
	if err != nil {
		return err
	}
	defer st.Close()
	pwd := sha256.New().Sum([]byte(u.Password + salt))
	_, err = st.Exec(u.PublicName, u.Username, u.Email, pwd)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser2Auth(email string, password string) (*User, error) {
	pwd := sha256.New().Sum([]byte(password + salt))
	row := s.db.QueryRow("SELECT fn, un, email, password FROM users WHERE email = ? AND password  = ?", email, pwd)
	var u User
	err := row.Scan(&u.PublicName, &u.Username, &u.Email, &u.Password)
	if err == sql.ErrNoRows {
		return nil, nil // пользователя нет - ошибки тоже
	}
	if err != nil {
		return nil, err //какая то ошибка с базой
	}
	return &u, nil // все заебись - это юзер для получения jwt
}
func (s *Storage) GetByUsername(username string) (*User, error) {

	row := s.db.QueryRow("SELECT fn, un, email, password FROM users WHERE un = ?", username)
	var u User
	err := row.Scan(&u.PublicName, &u.Username, &u.Email, &u.Password)
	if err == sql.ErrNoRows {
		return nil, nil // пользователя нет - ошибки тоже
	}
	if err != nil {
		return nil, err //какая то ошибка с базой
	}
	return &u, nil
}

func (s *Storage) GetAllChatMessages(seUn string, reUn string) int16 {
	return 1488
}
