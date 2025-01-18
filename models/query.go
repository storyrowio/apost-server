package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"strconv"
	"strings"
)

type Query struct {
	Keyword     string `form:"keyword"`
	Sort        string `form:"sort"`
	Limit       string `form:"limit"`
	Page        string `form:"page"`
	UserId      string `form:"user"`
	ProjectId   string `form:"project"`
	WorkspaceId string `form:"workspace"`
	Assigned    string `form:"assigned"`
	Completed   string `form:"completed"`
}

type Pagination struct {
	Count int64 `json:"count"`
	Page  int64 `json:"page"`
	Pages int64 `json:"pages"`
	Skip  int64 `json:"skip"`
	Limit int64 `json:"limit"`
}

type Result struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Query      Query       `json:"query"`
}

func (q Query) GetPagination(count int64) Pagination {
	pages := int64(0)
	page := int64(1)
	limit := int64(0)

	if q.Limit != "" {
		limit, _ = strconv.ParseInt(q.Limit, 10, 64)
		pages = int64(math.Round(float64(count) / float64(limit)))
	}

	if q.Page != "" {
		page, _ = strconv.ParseInt(q.Page, 10, 64)
	}

	skip := page*limit - limit

	return Pagination{
		Count: count,
		Page:  page,
		Pages: pages,
		Skip:  skip,
		Limit: limit,
	}
}

func (q Query) GetOptions() *options.FindOptions {
	opts := options.Find()

	limit, _ := strconv.ParseInt(q.Limit, 10, 64)
	opts.SetLimit(limit)

	if q.Sort != "" {
		splitted := strings.Split(q.Sort, ",")

		key := splitted[0]
		value, _ := strconv.ParseInt(splitted[1], 10, 64)

		opts.SetSort(bson.D{{key, value}})
	}

	page := int64(1)
	if q.Page != "" {
		page, _ = strconv.ParseInt(q.Page, 10, 64)
	}

	skip := page*limit - limit
	opts.SetSkip(skip)

	return opts
}

func (q Query) GetQueryFind() bson.M {
	query := bson.M{}

	if q.Keyword != "" {
		regex := primitive.Regex{
			Pattern: "^" + q.Keyword,
			Options: "i",
		}
		query["name"] = regex
		query["title"] = regex
	}

	if q.ProjectId != "" {
		query["projectId"] = q.ProjectId
	}

	if q.WorkspaceId != "" {
		query["workspaceId"] = q.WorkspaceId
	}

	return query
}
