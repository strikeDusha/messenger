package database

import (
	"crypto/sha256"
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var salt = os.Getenv("SALT")

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
		receiver TEXT REFERENCES users(un),
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

func (s *Storage) AddMessage(m *Message) error {
	q := "INSERT INTO messages (sender ,receiver ,content,sent) VALUES (?,?,?,?)"
	st, err := s.db.Prepare(q)
	if err != nil {
		return err
	}
	defer st.Close()

	_, err = st.Exec(m.Sender, m.Receiver, m.Text, m.Sent)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetChat(se string, re string) ([]*MsgMini, error) {

	row, err := s.db.Query("SELECT text,sent FROM users WHERE sender  = ? AND receiver  = ?", se, re)
	if err != nil {
		return nil, err
	}

	var msg []*MsgMini
	defer row.Close()
	for row.Next() {
		var u MsgMini
		if err := row.Scan(&u.Text, &u.Time); err != nil {
			return nil, err
		}
		msg = append(msg, &u)
	}
	if err := row.Err(); err != nil {
		log.Fatal(err)
	}
	return msg, nil
}
