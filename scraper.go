package main

import (
	"flag"
	"log"

	"gopkg.in/ini.v1"
)

func main() {
	limit := flag.Int("limit", 0, "how many problems to add to db")
	flag.Parse()

	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("fail to load config: %v", err)
	}
	dataSourceName := cfg.Section("database").Key("dsn").String()

	db := openDB(dataSourceName)
	defer db.Close()

	client := NewClient()

	topics, err := client.GetTopics()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got %d topics", len(topics))
	for _, t := range topics {
		err = db.addTopic(t)
		if err != nil {
			log.Fatalf("error adding topic [%s]: %v", t.Slug, err)
		}
	}

	last, err := db.getLastProblemID()
	if err != nil {
		log.Fatal(err)
	}

	skip, err := db.getProblemCount()
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for {
		_, questions, err := client.GetQuestionList(100, int(skip))
		log.Printf("got %d problems", len(questions))
		if err != nil {
			log.Fatal(err)
		}
		for _, q := range questions {
			if *limit > 0 && count == *limit {
				break
			}
			if q.LeetcodeID <= last || q.PaidOnly {
				continue
			}
			question, err := client.GetQuestion(q.TitleSlug)
			if err != nil {
				log.Fatalf("error getting question [%d. %s]: %v", q.LeetcodeID, q.TitleSlug, err)
			}
			err = db.addProblem(question)
			if err != nil {
				log.Fatalf("error adding problem [%d. %s] to db: %v", question.LeetcodeID, question.Title, err)
			}
			count++
			log.Printf("added problem [%d. %s]", question.LeetcodeID, question.Title)
		}
		if (*limit > 0 && count == *limit) || len(questions) == 0 {
			break
		}
		skip += 100
	}
	log.Printf("added %d problems to db", count)
}
