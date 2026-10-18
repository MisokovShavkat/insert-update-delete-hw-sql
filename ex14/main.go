package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Book struct {
	ID        int
	Text      string
	Author    string
	Price     int
	Pages     int
	CreatedAt time.Time
}

func ConnDB() *sql.DB {
	ConnStr := "host=localhost port=5432 user=postgres password=1234567890 dbname=books sslmode=disable"

	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		log.Fatal("не удалось создать пул:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("не удалось подключиться:", err)
	}
	return db
}

func AddBook(db *sql.DB) {
	var text, author string
	var price, pages int

	fmt.Print("Введите название книги: ")
	fmt.Scan(&text)

	fmt.Print("Введите автора: ")
	fmt.Scan(&author)

	fmt.Print("Введите цену: ")
	fmt.Scan(&price)

	fmt.Print("Введите количество страниц: ")
	fmt.Scan(&pages)

	query := `
		INSERT INTO books (text, author, price, pages)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(query, text, author, price, pages)
	if err != nil {
		log.Println("Ошибка добавления:", err)
		return
	}

	fmt.Println("Книга успешно добавлена!")
}

func ChangePrice(db *sql.DB) {
	var text string
	var price int

	fmt.Print("Введите название книги: ")
	fmt.Scan(&text)

	fmt.Print("Введите новую цену: ")
	fmt.Scan(&price)

	query := `
		UPDATE books
		SET price = $1
		WHERE text = $2
	`

	result, err := db.Exec(query, price, text)
	if err != nil {
		log.Println("Ошибка изменения цены:", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Ошибка получения количества строк:", err)
		return
	}

	if rows == 0 {
		fmt.Println("Книга не найдена!")
		return
	}

	fmt.Println("Цена книги успешно изменена!")
	fmt.Println("Количество изменённых строк:", rows)
}

func DeleteBook(db *sql.DB) {
	var id int

	fmt.Print("Введите ID книги для удаления: ")
	fmt.Scan(&id)

	query := `
		DELETE FROM books
		WHERE id = $1
	`

	result, err := db.Exec(query, id)
	if err != nil {
		log.Println("Ошибка удаления книги:", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Ошибка получения количества строк:", err)
		return
	}

	if rows == 0 {
		fmt.Println("Книга с таким ID не найдена!")
		return
	}

	fmt.Println("Книга успешно удалена!")
	fmt.Println("Количество удалённых строк:", rows)
}

func AddBookAndGetID(db *sql.DB) {
	var text string
	var author string
	var price int
	var pages int

	fmt.Print("Введите название книги: ")
	fmt.Scan(&text)

	fmt.Print("Введите автора книги: ")
	fmt.Scan(&author)

	fmt.Print("Введите цену книги: ")
	fmt.Scan(&price)

	fmt.Print("Введите количество страниц: ")
	fmt.Scan(&pages)

	query := `
		INSERT INTO books (text, author, price, pages)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int

	err := db.QueryRow(query, text, author, price, pages).Scan(&id)
	if err != nil {
		log.Println("Ошибка добавления книги:", err)
		return
	}

	fmt.Println("Книга успешно добавлена!")
	fmt.Println("ID новой книги:", id)
}
func PrintBooks(db *sql.DB) {
	query := `SELECT id, text, author, price, pages FROM books`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Ошибка получения книг:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, price, pages int
		var text, author string

		err := rows.Scan(&id, &text, &author, &price, &pages)
		if err != nil {
			log.Fatal("Ошибка чтения книги:", err)
		}

		fmt.Println("ID:", id)
		fmt.Println("Название:", text)
		fmt.Println("Автор:", author)
		fmt.Println("Цена:", price)
		fmt.Println("Страниц:", pages)
		fmt.Println("-------------------")
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Ошибка чтения строк:", err)
	}
}

func main() {
	db := ConnDB()
	defer db.Close()

	for {
		var choice int

		fmt.Println("\n===== МЕНЮ КНИГ =====")
		fmt.Println("1. Добавить книгу")
		fmt.Println("2. Изменить цену")
		fmt.Println("3. Удалить книгу")
		fmt.Println("4. Добавить книгу и получить ID")
		fmt.Println("5. Показать все книги")
		fmt.Println("0. Выход")
		fmt.Print("Выберите действие: ")

		_, err := fmt.Scan(&choice)
		if err != nil {
			log.Println("Введите число!")
			return
		}

		switch choice {
		case 1:
			AddBook(db)
		case 2:
			ChangePrice(db)
		case 3:
			DeleteBook(db)
		case 4:
			AddBookAndGetID(db)
		case 5:
			PrintBooks(db)
		case 0:
			fmt.Println("Программа завершена!")
			return
		default:
			fmt.Println("Неверный выбор!")
		}
	}
}
