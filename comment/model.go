package comment

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/seagalputra/talkbox/config"
)

type (
	CommentColumns []string

	Comment struct {
		ID           string  `db:"id"`
		ParentID     *string `db:"parent_id"`
		PostID       string  `db:"post_id"`
		Body         string  `db:"body"`
		Attachment   *string `db:"attachment"`
		LikeCount    int     `db:"like_count"`
		DislikeCount int     `db:"dislike_count"`
		ModeratedBy  string  `db:"moderated_by"`
		CreatedAt    string  `db:"created_at"`
		UpdatedAt    string  `db:"updated_at"`
		DeletedAt    *string `db:"deleted_at"`
	}
)

var (
	Columns = CommentColumns{
		"id",
		"parent_id",
		"post_id",
		"body",
		"attachment",
		"like_count",
		"dislike_count",
		"moderated_by",
		"created_at",
		"updated_at",
		"deleted_at",
	}
)

const (
	SHEET_NAME = "comments"
)

func New(parentID *string, postID string, body string, attachment *string, moderatedBy string, likeCount, dislikeCount int) Comment {
	guid := xid.New()
	now := time.Now().Format(time.RFC3339)

	return Comment{
		ID:           guid.String(),
		ParentID:     parentID,
		PostID:       postID,
		Body:         body,
		Attachment:   attachment,
		LikeCount:    likeCount,
		DislikeCount: dislikeCount,
		ModeratedBy:  moderatedBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func Save(postID string, commentReq InsertCommentReq) (*Comment, error) {
	spreadsheetID := config.AppConfig.SpreadsheetID
	store := config.GetSheetDB(spreadsheetID, SHEET_NAME, Columns)
	defer store.Close(context.Background())

	comment := New(
		commentReq.ParentID,
		postID,
		commentReq.Body,
		commentReq.Attachment,
		"user", // TODO: change to real user
		0,
		0,
	)

	err := store.Insert(comment).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func FindAll() ([]Comment, error) {
	panic("Not implemented yet!")
}
