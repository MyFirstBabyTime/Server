package http

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"time"

	"github.com/MyFirstBabyTime/Server/domain"
)

// childrenHandler represent the http handler for children
type childrenHandler struct {
	cUsecase   domain.ChildrenUsecase
	validator  validator
	jwtHandler jwtHandler
}

// jwtHandler is interface of jwt handler
type jwtHandler interface {
	// ParseUUIDFromToken parse token & return token payload and type
	ParseUUIDFromToken(c *gin.Context)
}

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// NewChildrenHandler will initialize the children resources endpoint
func NewChildrenHandler(r *gin.Engine, cu domain.ChildrenUsecase, v validator, jh jwtHandler) {
	h := &childrenHandler{
		cUsecase:   cu,
		validator:  v,
		jwtHandler: jh,
	}

	r.POST("parents/uuid/:parent_uuid/children", h.jwtHandler.ParseUUIDFromToken, h.CreateNewChildren)
}

func (ch *childrenHandler) CreateNewChildren(c *gin.Context) {
	req := new(createNewChildrenRequest)
	if err := ch.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	if c.GetString("uuid") != req.ParentUUID {
		c.JSON(http.StatusForbidden, defaultResp(http.StatusForbidden, 0, "you can't access with that uuid token"))
		return
	}

	chi := &domain.Children{
		ParentUUID: domain.String(req.ParentUUID),
		Name:       domain.String(req.Name),
		Sex:        domain.String(req.Sex),
	}

	if t, err := time.Parse(time.RFC3339, req.Birth); err != nil {
		err = errors.Wrap(err, "failed to parse birth time string")
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	} else {
		chi.Birth = domain.Time(t)
	}

	var profile []byte
	if req.Profile != nil {
		profile = make([]byte, req.Profile.Size)
		file, _ := req.Profile.Open()
		defer func() { _ = file.Close() }()
		_, _ = file.Read(profile)
	} else if req.ProfileBase64 != "" {
		req.ProfileBase64 = string(regexp.MustCompile("^data:image/\\w+;base64,").ReplaceAll([]byte(req.ProfileBase64), []byte("")))
		var err error
		if profile, err = base64.StdEncoding.DecodeString(req.ProfileBase64); err != nil {
			err = errors.Wrap(err, "failed to decode base64 string to byte array")
			c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, err.Error()))
			return
		}
	}

	switch uuid, err := ch.cUsecase.CreateNewChildren(c.Request.Context(), chi, profile); tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusCreated, 0, "succeed to create new children")
		resp["children_uuid"] = uuid
		c.JSON(http.StatusCreated, resp)
	case domain.UsecaseError:
		c.JSON(tErr.Status, defaultResp(tErr.Status, tErr.Code, tErr.Error()))
	default:
		msg := errors.Wrap(err, "CreateNewChildren return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}

// bindRequest method bind *gin.Context to request having BindFrom method
func (ch *childrenHandler) bindRequest(req interface {
	BindFrom(ctx *gin.Context) error
}, c *gin.Context) error {
	if err := req.BindFrom(c); err != nil {
		return errors.Wrap(err, "failed to bind req")
	}
	if err := ch.validator.ValidateStruct(req); err != nil {
		return errors.Wrap(err, "invalid request")
	}
	return nil
}

// defaultResp return response have status, code, message inform
func defaultResp(status, code int, msg string) (resp gin.H) {
	resp = gin.H{}
	resp["status"] = status
	resp["code"] = code
	resp["message"] = msg
	return
}
