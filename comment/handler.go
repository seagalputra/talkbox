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

		Data Comment `json:"data"`
	}

	UpdateCommentReq struct {
		Body       *string `json:"body"`
		Attachment *string `json:"attachment"`
	}

	ListCommentReq struct {
		Page        string
		Limit       string
		SearchQuery string
	}

	ListCommentRes struct {
		utils.CommonRes

		Data []Comment `json:"data"`
	}

	Handler struct {
		FindAllCommentFunc func(ListCommentReq) []Comment
		InsertCommentFunc  func(string, InsertCommentReq) (*Comment, error)
		DeleteCommentFunc  func(string) error
		UpdateCommentFunc  func(string, UpdateCommentReq) (*Comment, error)
	}
)

func DefaultHandler() Handler {
	return Handler{
		InsertCommentFunc:  Save,
		FindAllCommentFunc: FindAll,
		UpdateCommentFunc:  Update,
		DeleteCommentFunc:  Remove,
	}
}

// Insert
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

// FindAll
// @Summary     Get all comments
// @Description Get all avaiable comments within list
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       page         query    string false "Comment data in page"
// @Param       limit        query    string false "Limit data per request"
// @Param       search_query query    string false "Search query for comment"
// @Success     200          {object} comment.ListCommentRes
// @Failure     400          {object} utils.ErrorRes
// @Failure     401          {object} utils.ErrorRes
// @Failure     404          {object} utils.ErrorRes
// @Failure     500          {object} utils.ErrorRes
// @Router      /comments [get]
func (f *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonEncoder := json.NewEncoder(w)

	// page := r.URL.Query().Get("page")
	// limit := r.URL.Query().Get("limit")
	// searchQuery := r.URL.Query().Get("search_query")

	listCommentReq := ListCommentReq{}
	comments := f.FindAllCommentFunc(listCommentReq)

	response := ListCommentRes{}
	response.CommonRes.Status = utils.SUCCESS
	response.Data = comments

	w.WriteHeader(http.StatusOK)
	jsonEncoder.Encode(&response)
}

func (f *Handler) Delete(w http.ResponseWriter, r *http.Request) {}

func (f *Handler) Update(w http.ResponseWriter, r *http.Request) {}
