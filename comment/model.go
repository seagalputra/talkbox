package comment

import (
	"context"
	"log"
	"time"

	"github.com/rs/xid"
	"github.com/seagalputra/talkbox/config"
)

type (
	CommentColumns []string

	Comment struct {
		ID           string  `db:"id" json:"id"`
		ParentID     *string `db:"parent_id" json:"parent_id"`
		PostID       string  `db:"post_id" json:"post_id"`
		Body         string  `db:"body" json:"body"`
		Attachment   *string `db:"attachment" json:"attachment"`
		LikeCount    int     `db:"like_count" json:"like_count"`
		DislikeCount int     `db:"dislike_count" json:"dislike_count"`
		ModeratedBy  string  `db:"moderated_by" json:"moderated_by"`
		CreatedAt    string  `db:"created_at" json:"created_at"`
		UpdatedAt    string  `db:"updated_at" json:"updated_at"`
		DeletedAt    *string `db:"deleted_at" json:"-"`
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
	store := config.GetSheetDB(SHEET_NAME, Columns)
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
		log.Printf("Save: %v", err)
		return nil, err
	}

	return &comment, nil
}

func FindAll(commentReq ListCommentReq) []Comment {
	store := config.GetSheetDB(SHEET_NAME, Columns)
	defer store.Close(context.Background())

	var comments []Comment
	err := store.Select(&comments).Exec(context.Background())
	if err != nil {
		log.Printf("FindAll: %v", err)
		return []Comment{}
	}

	return comments
}

func Update(id string, commentReq UpdateCommentReq) (*Comment, error) {
	panic("Not implemented yet!")
}

func Remove(id string) error {
	panic("Not implemented yet!")
}
