package zen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	common_gorm "github.com/ngtrvu/zen-go/gorm"
)

const PAGE_SIZE = 24
const DEFAULT_MAC_ADDRESS = "74-E6-E2-10-6B-80"

const (
	CtxUserIDKey  = "userID"
	CtxUserKey    = "user"
	CtxNHSVToken  = "nhsvToken"
	CtxMacAddress = "macAddress"
	CtxIpAddress  = "ipAddress"
	CtxClaimsKey  = "claims"
)

const (
	CtxEmailKey = "email"
	CtxPhoneKey = "phone"
)

type HttpHandler struct {
	JWTSecretKey []byte
	Config       *ZenConfig
}

type FileData struct {
	FileName string
	FileType string
	FileSize int64
	Data     []byte
}

// QueryListOptions represents the query parameters for a list of resources.
type QueryListOptions struct {
	// The page number to return.
	Page int `json:"page"      minimum:"1" validate:"optional" example:"1"`
	// The number of items to return per page.
	PageSize int `json:"page_size" minimum:"1" validate:"optional" example:"20" maximum:"50"`
}

type Response struct {
	// Success represents request is success or not
	Success bool `json:"success"`

	// Data represents response data
	Data interface{} `json:"data"`

	// MetaData represents extra data beside Data. it can be empty
	MetaData interface{} `json:"meta_data"`

	// Error represents error message
	Error string `json:"error" example:""`

	// ErrorCode represents error code. for detail error will show in Error field
	ErrorCode string `json:"error_code" example:""`

	// Pagination represents pagination data
	Pagination interface{} `json:"paging"` // paging
}

type Pagination struct {
	Count     int `json:"count"      example:"10"`
	Page      int `json:"page"       example:"1"`
	PageCount int `json:"page_count" example:"1"`
	PageSize  int `json:"page_size"  example:"20"`
}

func NewHttpHandler(
	cfg *ZenConfig,
) (*HttpHandler, error) {
	var secretKey []byte
	var err error

	if cfg.ECDSAPrivateKeyPath != "" {
		secretKey, err = os.ReadFile(cfg.ECDSAPrivateKeyPath)
		if err != nil {
			return nil, err
		}
		cfg.SecretKey = string(secretKey)
	} else {
		secretKey = []byte(cfg.SecretKey)
	}

	httpHandler := &HttpHandler{
		JWTSecretKey: secretKey,
		Config:       cfg,
	}
	return httpHandler, nil
}

func newFailureResponse(err error) *Response {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return &Response{Success: false, Error: appErr.Message, ErrorCode: appErr.Code}
	} else if err != nil {
		return &Response{Success: false, Error: err.Error(), ErrorCode: ""}
	} else {
		return &Response{Success: false, Error: ErrBadRequest.Message, ErrorCode: ErrBadRequest.Code}
	}
}

// get page_size from the request params. Maximum page_size is 100
func (h *HttpHandler) RequestPage(r *http.Request) int {
	value := r.URL.Query().Get("page")
	page, err := strconv.Atoi(value)
	if err != nil {
		return 1
	}
	if page < 1 {
		page = 1
	}
	return page
}

// get page_size from the request params. Maximum page_size is 100
func (h *HttpHandler) RequestPageSize(r *http.Request) int {
	value := r.URL.Query().Get("page_size")
	page_size, err := strconv.Atoi(value)
	if err != nil {
		return PAGE_SIZE
	}

	if page_size > 100 {
		return 100
	}

	return page_size
}

func (h *HttpHandler) GetGormSearch(r *http.Request, searchFields []*common_gorm.SearchField) common_gorm.Search {
	params := r.URL.Query()
	search := common_gorm.Search{}
	for _, field := range searchFields {
		value := params.Get("search")
		if value == "" {
			continue
		}

		search.AddSearchField(
			&common_gorm.SearchAttribute{
				Field:        field.Field,
				JoinedColumn: field.JoinedColumn,
				Type:         common_gorm.FieldTypeString,
				Operator:     common_gorm.OperatorLike,
				Value:        value,
			},
		)
	}

	return search
}

func (h *HttpHandler) GetDefaultSortFields() []*common_gorm.SortField {
	return []*common_gorm.SortField{{SortBy: "created_at", SortOrder: common_gorm.QuerySortDESC}}
}

func (h *HttpHandler) GetFilteringQueryset(
	r *http.Request,
	searchFields []*common_gorm.SearchField,
	defaultSortFields []*common_gorm.SortField,
) *common_gorm.Query {
	page := h.RequestPage(r)
	pageSize := h.RequestPageSize(r)

	fieldFilters := h.GetGormFilters(r, common_gorm.FilterConfig{})
	search := h.GetGormSearch(r, searchFields)
	orderingFilterFields := h.GetGormSorter(r, defaultSortFields)

	return &common_gorm.Query{
		Limit:      pageSize,
		Offset:     (page - 1) * pageSize,
		Filter:     *fieldFilters,
		Search:     search,
		SortFields: orderingFilterFields,
	}
}

func GetIpAddress(r *http.Request) string {
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Real-IP")
	}
	if ipAddress == "" {
		ipAddress, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return ipAddress
}

func GetMacAddress(r *http.Request) string {
	macAddress := r.Header.Get("X-Device-MAC")
	if macAddress == "" {
		macAddress = DEFAULT_MAC_ADDRESS
	}

	return macAddress
}

func GetNHSVToken(r *http.Request) string {
	return r.Header.Get("Authorization-Nhsv")
}

func (h *HttpHandler) GetGormSorter(
	r *http.Request,
	defaultSortFields []*common_gorm.SortField,
) []*common_gorm.SortField {
	params := r.URL.Query()
	sortFields := []*common_gorm.SortField{}

	if params.Get("ordering") != "" {
		sortFieldKeys := strings.Split(params.Get("ordering"), ",")

		for _, key := range sortFieldKeys {
			field := parseField(key)
			if strings.HasPrefix(field, "-") {
				field = strings.Replace(field, "-", "", 1)
				sortFields = append(
					sortFields,
					&common_gorm.SortField{SortBy: field, SortOrder: common_gorm.QuerySortDESC},
				)
			} else {
				sortFields = append(sortFields, &common_gorm.SortField{SortBy: field, SortOrder: common_gorm.QuerySortASC})
			}
		}
	} else if len(defaultSortFields) > 0 {
		return defaultSortFields
	}

	return sortFields
}

func (h *HttpHandler) GetGormFilters(r *http.Request, config common_gorm.FilterConfig) *common_gorm.Filter {
	params := r.URL.Query()
	filter := &common_gorm.Filter{}
	for key := range params {
		if params.Get(key) == "" {
			continue
		}

		if key == "page" {
			continue
		}

		if key == "page_size" {
			continue
		}

		if key == "ordering" {
			continue
		}

		if key == "search" {
			continue
		}

		field := parseField(key)
		operator := parseOperator(key)
		value := parseValue(params.Get(key), operator)

		filter.AddFilter(
			&common_gorm.FilterAttribute{Field: field, Type: "general", Operator: operator, Value: value},
		)
	}

	return filter
}

func parseOperator(key string) string {
	partKeys := strings.Split(key, "__")

	operator := common_gorm.OperatorEqual
	if len(partKeys) > 1 {
		if partKeys[1] == "in" {
			operator = common_gorm.OperatorIn
		}
	} else if key == "search" {
		operator = common_gorm.OperatorLike
	}

	return operator
}

func parseField(key string) string {
	partKeys := strings.Split(key, "__")
	return partKeys[0]
}

func parseValue(value, operator string) interface{} {
	if operator != common_gorm.OperatorIn && operator != common_gorm.OperatorNotIn {
		return value
	}

	items := strings.Split(value, ",")
	if len(items) == 1 {
		return items[0]
	}

	return items
}

func (ctrl HttpHandler) Success(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Response{Success: true, Data: data})
}

func (ctrl HttpHandler) SuccessWithPagination(w http.ResponseWriter, r *http.Request, data interface{}, total int) {
	page := ctrl.RequestPage(r)
	pageSize := ctrl.RequestPageSize(r)

	render.Status(r, http.StatusOK)
	pagination := Pagination{
		Count:     total,
		Page:      page,
		PageSize:  pageSize,
		PageCount: total/pageSize + 1,
	}

	render.JSON(w, r, Response{Success: true, Data: data, Pagination: pagination})
}

func (ctrl HttpHandler) SuccessWithMetadata(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
	metaData interface{},
	total int,
	page, pageSize int,
) {
	render.Status(r, http.StatusOK)
	pagination := Pagination{
		Count:     total,
		Page:      page,
		PageSize:  pageSize,
		PageCount: total/pageSize + 1,
	}

	render.JSON(w, r, Response{Success: true, Data: data, MetaData: metaData, Pagination: pagination})
}

func (ctrl HttpHandler) SuccessCreated(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, Response{Success: true, Data: data})
}

func (ctrl HttpHandler) SuccessNoContent(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNoContent)
	render.JSON(w, r, Response{Success: true})
}

func (ctrl HttpHandler) ServerError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(newFailureResponse(err))
}

func (ctrl HttpHandler) BadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(newFailureResponse(err))
}

func (ctrl HttpHandler) NotFound(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(newFailureResponse(err))
}

func (ctrl HttpHandler) Unauthorized(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	json.NewEncoder(w).Encode(newFailureResponse(err))
}

func (ctrl HttpHandler) Forbidden(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(newFailureResponse(err))
}

func GetFormFileData(r *http.Request, key string) (fileData *FileData, err error) {
	file, fileHeader, err := r.FormFile(key)
	if err == http.ErrMissingFile {
		return nil, nil
	}

	if err != nil {
		return nil, ErrOpenUploadFile
	}
	defer file.Close()

	if fileHeader.Size > MaxUploadSizeInMegabyteDefault*1024*1024 {
		msg := fmt.Sprintf("image file size exceeds %dMB", MaxUploadSizeInMegabyteDefault)
		return nil, NewAppError("invalid_image_size", msg)
	}

	mimeTypeData := make([]byte, 512)
	_, err = file.Read(mimeTypeData)
	if err != nil {
		return nil, ErrInvalidImg
	}

	filetype := http.DetectContentType(mimeTypeData)
	if filetype != "image/jpeg" && filetype != "image/png" {
		return nil, ErrInvalidImgFormat
	}

	// Read from the start of the file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, ErrInvalidImgFormat
	}

	// init a buffer to read the file content
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		return nil, ErrInvalidImgFormat
	}

	fileData = &FileData{
		FileType: filetype,
		FileName: fileHeader.Filename,
		FileSize: fileHeader.Size,
		Data:     buf.Bytes(),
	}

	return fileData, nil
}
