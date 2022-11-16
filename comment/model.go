package comment

import (
	"time"

	"github.com/rs/xid"
)

type (
	CommentColumns []string

	Comment struct {
		ID           string `db:"id"`
		ParentID     string `db:"parent_id"`
		Body         string `db:"body"`
		Attachment   string `db:"attachment"`
		LikeCount    int    `db:"like_count"`
		DislikeCount int    `db:"dislike_count"`
		ModeratedBy  string `db:"moderated_by"`
		CreatedAt    string `db:"created_at"`
		UpdatedAt    string `db:"updated_at"`
		DeletedAt    string `db:"deleted_at"`
	}
)

var (
	Columns = CommentColumns{
		"id",
		"parent_id",
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

func New(parentID, body, attachment, moderatedBy string, likeCount, dislikeCount int) Comment {
	guid := xid.New()
	now := time.Now().Format(time.RFC3339)

	return Comment{
		ID:           guid.String(),
		ParentID:     parentID,
		Body:         body,
		Attachment:   attachment,
		LikeCount:    likeCount,
		DislikeCount: dislikeCount,
		ModeratedBy:  moderatedBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func Save(commentReq InsertCommentReq) error {
	panic("Not implemented yet!")
}

func FindAll() ([]Comment, error) {
	panic("Not implemented yet!")
}
