package comment

import (
	"encoding/json"
	"net/http"

	"github.com/seagalputra/talkbox/utils"
)

type (
	InsertCommentReq struct {
		Body       string `json:"body"`
		Attachment string `json:"attachment"`
	}

	InsertCommentRes struct {
		utils.CommonRes
	}

	UpdateCommentReq struct {
		Body       *string
		Attachment *string
	}

	Handler struct {
		FindAllCommentFunc func() ([]Comment, error)
		InsertCommentFunc  func(InsertCommentReq) error
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
// @Param       request body     comment.InsertCommentReq true "Insert comment request body"
// @Success     200     {object} comment.InsertCommentRes
// @Failure     400     {object} utils.ErrorRes
// @Failure     401     {object} utils.ErrorRes
// @Failure     404     {object} utils.ErrorRes
// @Failure     422     {object} utils.ErrorRes
// @Failure     500     {object} utils.ErrorRes
// @Router      /comments [post]
func (f *Handler) Insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request InsertCommentReq
	json.NewDecoder(r.Body).Decode(&request)

	// insert to google sheet
	err := f.InsertCommentFunc(request)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := utils.ErrorRes{
			Status:    utils.ERROR,
			Message:   "Failed to insert a new comment",
			ErrorCode: "ERR_INSERT_FAILED",
		}
		json.NewEncoder(w).Encode(&errorResponse)
		return
	}
	// process the output
}

func (f *Handler) FindAll(w http.ResponseWriter, r *http.Request) {}

func (f *Handler) Delete(w http.ResponseWriter, r *http.Request) {}

func (f *Handler) Update(w http.ResponseWriter, r *http.Request) {}
