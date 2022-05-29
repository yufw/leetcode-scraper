package main

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/shurcooL/graphql"
)

type Client struct {
	*graphql.Client
}

type Question struct {
	LeetcodeID       int32
	PaidOnly         bool
	Title            string
	TitleSlug        string
	Content          string
	Difficulty       string
	Likes            int32
	Dislikes         int32
	TopicTags        []string
	TotalAccepted    int64
	TotalSubmission  int64
	SimilarQuestions string
	Hints            []string
}

type QuestionStats struct {
	TotalAccepted   int64 `json:"totalAcceptedRaw"`
	TotalSubmission int64 `json:"totalSubmissionRaw"`
}

type Topic struct {
	Name string
	Slug string
}

func NewClient() *Client {
	client := graphql.NewClient("https://leetcode.com/graphql", nil)

	return &Client{Client: client}
}

func (client *Client) GetQuestion(titleSlug string) (*Question, error) {
	var q struct {
		Question struct {
			QuestionFrontendID graphql.String
			Title              graphql.String
			TitleSlug          graphql.String
			Content            graphql.String
			Difficulty         graphql.String
			Likes              graphql.Int
			Dislikes           graphql.Int
			TopicTags          []struct {
				Name graphql.String
				Slug graphql.String
			}
			Stats            graphql.String
			SimilarQuestions graphql.String
			Hints            []graphql.String
		} `graphql:"question(titleSlug: $titleSlug)"`
	}
	variables := map[string]interface{}{
		"titleSlug": graphql.String(titleSlug),
	}

	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		return nil, err
	}
	leetcodeID, err := strconv.Atoi(string(q.Question.QuestionFrontendID))
	if err != nil {
		return nil, err
	}
	topicTags := make([]string, len(q.Question.TopicTags))
	for i, v := range q.Question.TopicTags {
		topicTags[i] = string(v.Slug)
	}
	var stats QuestionStats
	if err = json.Unmarshal([]byte(string(q.Question.Stats)), &stats); err != nil {
		return nil, err
	}
	hints := make([]string, len(q.Question.Hints))
	for i, v := range q.Question.Hints {
		hints[i] = string(v)
	}

	return &Question{
		LeetcodeID:       int32(leetcodeID),
		Title:            string(q.Question.Title),
		TitleSlug:        string(q.Question.TitleSlug),
		Content:          string(q.Question.Content),
		Difficulty:       string(q.Question.Difficulty),
		Likes:            int32(q.Question.Likes),
		Dislikes:         int32(q.Question.Dislikes),
		TopicTags:        topicTags,
		TotalAccepted:    stats.TotalAccepted,
		TotalSubmission:  stats.TotalSubmission,
		SimilarQuestions: string(q.Question.SimilarQuestions),
		Hints:            hints,
	}, nil
}

func (client *Client) GetQuestionList(limit int, skip int) (int, []*Question, error) {
	var q struct {
		QuestionList struct {
			Total     graphql.Int `graphql:"total: totalNum"`
			Questions []struct {
				FrontendQuestionID graphql.String  `graphql:"frontendQuestionId: questionFrontendId"`
				PaidOnly           graphql.Boolean `graphql:"paidOnly: isPaidOnly"`
				TitleSlug          graphql.String
			} `graphql:"questions: data"`
		} `graphql:"questionList(limit: $limit, skip: $skip, categorySlug: $categorySlug, filters: {})"`
	}
	variables := map[string]interface{}{
		"categorySlug": graphql.String("algorithms"),
		"limit":        graphql.Int(limit),
		"skip":         graphql.Int(skip),
	}

	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		return 0, nil, err
	}

	var questions []*Question
	for _, v := range q.QuestionList.Questions {
		leetcodeID, err := strconv.Atoi(string(v.FrontendQuestionID))
		if err != nil {
			return 0, nil, err
		}
		questions = append(questions, &Question{
			LeetcodeID: int32(leetcodeID),
			PaidOnly:   bool(v.PaidOnly),
			TitleSlug:  string(v.TitleSlug),
		})
	}

	return int(q.QuestionList.Total), questions, nil
}

func (client *Client) GetTopics() ([]*Topic, error) {
	var q struct {
		QuestionTopicTags struct {
			Edges []struct {
				Node struct {
					Name graphql.String
					Slug graphql.String
				}
			}
		}
	}
	err := client.Query(context.Background(), &q, nil)
	if err != nil {
		return nil, err
	}

	topics := make([]*Topic, len(q.QuestionTopicTags.Edges))
	for i, v := range q.QuestionTopicTags.Edges {
		topics[i] = &Topic{
			Name: string(v.Node.Name),
			Slug: string(v.Node.Slug),
		}
	}

	return topics, nil
}
