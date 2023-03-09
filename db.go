package main

import (
	"database/sql"
	"fmt"
)

func dbConnect() error {
	db, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		return err
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'articles';").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		fmt.Println("No articles table, creating...")
		err = dbCreateTable()
		if err != nil {
			return err
		}
	}

	return nil
}

func dbCreateTable() error {
	stmtQuery := "CREATE TABLE IF NOT EXISTS articles (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, title TEXT, content TEXT);"
	stmt, err := db.Prepare(stmtQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("Database created successfully!")
	return nil
}

func dbCreateArticle(article *Article) error {
	query, err := db.Prepare("INSERT INTO articles(title, content) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(article.Title, article.Content)

	if err != nil {
		return err
	}

	return nil
}

func dbGetArticle(articleID string) (*Article, error) {
	query, err := db.Prepare("SELECT id, title, content FROM articles WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	result := query.QueryRow(articleID)
	article := new(Article)
	err = result.Scan(
		&article.ID,
		&article.Title,
		&article.Content,
	)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func dbGetAllArticles() ([]*Article, error) {
	query, err := db.Prepare("SELECT id, title, content FROM articles")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, 0)
	for result.Next() {
		data := new(Article)
		err := result.Scan(
			&data.ID,
			&data.Title,
			&data.Content,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, data)
	}
	return articles, nil
}

func dbUpdateArticle(article *Article) (bool, error) {
	query, err := db.Prepare("UPDATE articles SET (title, content) = (?, ?) where id = ?")
	if err != nil {
		return false, err
	}
	defer query.Close()

	result, err := query.Exec(&article.Title, &article.Content, &article.ID)
	if err != nil {
		return false, err
	}

	updated, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if updated == 1 {
		return true, nil
	}

	return false, nil
}

func dbDeleteArticle(articleID string) (bool, error) {
	query, err := db.Prepare("DELETE FROM articles WHERE id = ?")
	if err != nil {
		return false, err
	}
	defer query.Close()

	result, err := query.Exec(articleID)
	if err != nil {
		return false, err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if deleted == 1 {
		return true, nil
	}

	return false, nil
}
