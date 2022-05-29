package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Database struct {
	*sql.DB
}

func openDB(dataSourceName string) *Database {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &Database{DB: db}
}

func (db *Database) getLastProblemID() (int32, error) {
	var id int32
	err := db.QueryRow("SELECT COALESCE(MAX(leetcode_id), 0) FROM problems").Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (db *Database) addTopic(t *Topic) error {
	_, err := db.Exec("INSERT INTO topics (slug, name) VALUES ($1, $2) ON CONFLICT DO NOTHING", t.Slug, t.Name)

	return err
}

func (db *Database) addProblem(p *Question) error {
	_, err := db.Exec("INSERT INTO problems (leetcode_id, title, title_slug, content, difficulty, likes, dislikes, total_accepted, total_submission, similar_questions, hints) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		p.LeetcodeID, p.Title, p.TitleSlug, p.Content, p.Difficulty, p.Likes, p.Dislikes, p.TotalAccepted, p.TotalSubmission, p.SimilarQuestions, p.Hints)

	if err != nil {
		return err
	}

	for _, topic := range p.TopicTags {
		_, err = db.Exec("INSERT INTO problem_topic (problem_id, topic_slug) VALUES ($1, $2)", p.LeetcodeID, topic)
		if err != nil {
			return err
		}
	}

	return nil
}
