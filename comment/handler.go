package comment

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/seagalputra/talkbox/utils"
)

type (
	InsertCommentReq struct {
		ParentID   *string `json:"parent_id"`
		Body       string  `json:"body"`
		Attachment *string `json:"attachment"`
	}

	InsertCommentRes struct {
		utils.CommonRes

		Data struct {
			ID           string  `json:"id"`
			ParentID     *string `json:"parent_id"`
			PostID       string  `json:"post_id"`
			Body         string  `json:"body"`
			Attachment   *string `json:"attachment"`
			LikeCount    int     `json:"like_count"`
			DislikeCount int     `json:"dislike_count"`
			CreatedAt    string  `json:"created_at"`
			UpdatedAt    string  `json:"updated_at"`
		} `json:"data"`
	}

	UpdateCommentReq struct {
		Body       *string `json:"body"`
		Attachment *string `json:"attachment"`
	}

	Handler struct {
		FindAllCommentFunc func() ([]Comment, error)
		InsertCommentFunc  func(string, InsertCommentReq) (*Comment, error)
		DeleteCommentFunc  func(string) error
		UpdateCommentFunc  func(string, UpdateCommentReq) error
	}
)

func DefaultHandler() Handler {
	return Handler{
		InsertCommentFunc:  Save,
		FindAllCommentFunc: FindAll,
	}
}

// InsertComment
// @Summary     Insert a new comment
// @Description create a new comment based on post slug
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       post_id path     string                   true "Post identifier"
// @Param       request body     comment.InsertCommentReq true "Insert comment request body"
// @Success     200     {object} comment.InsertCommentRes
// @Failure     400     {object} utils.ErrorRes
// @Failure     401     {object} utils.ErrorRes
// @Failure     404     {object} utils.ErrorRes
// @Failure     422     {object} utils.ErrorRes
// @Failure     500     {object} utils.ErrorRes
// @Router      /comments/{post_id} [post]
func (f *Handler) Insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonEncoder := json.NewEncoder(w)
	postID := chi.URLParam(r, "post_id")

	var request InsertCommentReq
	json.NewDecoder(r.Body).Decode(&request)

	// insert to google sheet
	comment, err := f.InsertCommentFunc(postID, request)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := utils.ErrorRes{
			Status:    utils.ERROR,
			Message:   "Failed to insert a comment",
			ErrorCode: "ERR_INSERT_FAILED",
		}
		jsonEncoder.Encode(&errorResponse)
		return
	}

	response := InsertCommentRes{}
	response.CommonRes.Message = "Successfully adding new comment"
	response.CommonRes.Status = utils.SUCCESS
	response.Data.ID = comment.ID
	response.Data.ParentID = comment.ParentID
	response.Data.PostID = comment.PostID
	response.Data.Body = comment.Body
	response.Data.Attachment = comment.Attachment
	response.Data.LikeCount = comment.LikeCount
	response.Data.DislikeCount = comment.DislikeCount
	response.Data.CreatedAt = comment.CreatedAt
	response.Data.UpdatedAt = comment.UpdatedAt

	// process the output
	w.WriteHeader(http.StatusCreated)
	jsonEncoder.Encode(&response)
}

func (f *Handler) FindAll(w http.ResponseWriter, r *http.Request) {}

func (f *Handler) Delete(w http.ResponseWriter, r *http.Request) {}

func (f *Handler) Update(w http.ResponseWriter, r *http.Request) {}
