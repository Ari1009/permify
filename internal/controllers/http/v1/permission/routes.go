package permission

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/Permify/permify/internal/commands"
	"github.com/Permify/permify/internal/controllers/http/common"
	"github.com/Permify/permify/internal/services"
	"github.com/Permify/permify/pkg/database"
	"github.com/Permify/permify/pkg/errors"
	"github.com/Permify/permify/pkg/logger"
)

var tracer = otel.Tracer("routes")

// permissionRoutes -
type permissionRoutes struct {
	service services.IPermissionService
	logger  logger.Interface
}

// NewPermissionRoutes -
func NewPermissionRoutes(handler *echo.Group, t services.IPermissionService, l logger.Interface) {
	r := &permissionRoutes{t, l}

	h := handler.Group("/permissions")
	{
		h.POST("/check", r.check)
		h.POST("/expand", r.expand)
		// h.POST("/lookup-query", r.lookupQuery)
	}
}

// @Summary     Permission
// @Description check subject is authorized
// @ID          permissions.check
// @Tags  	    Permission
// @Accept      json
// @Produce     json
// @Param       request body CheckRequest true "check subject is authorized"
// @Success     200 {object} CheckResponse
// @Failure     400 {object} common.HTTPErrorResponse
// @Router      /permissions/check [post]
func (r *permissionRoutes) check(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "permissions.check")
	defer span.End()

	request := new(CheckRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}
	v := request.Validate()
	if v != nil {
		return c.JSON(http.StatusUnprocessableEntity, common.ValidationResponse(v))
	}

	var err errors.Error

	if request.Depth == 0 {
		request.Depth = 20
	}

	var response commands.CheckResponse
	response, err = r.service.Check(ctx, request.Subject, request.Action, request.Entity, request.SchemaVersion.String(), request.Depth)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		switch err.Kind() {
		case errors.Database:
			return c.JSON(database.GetKindToHttpStatus(err.SubKind()), common.MResponse(err.Error()))
		case errors.Validation:
			return c.JSON(http.StatusUnprocessableEntity, common.ValidationResponse(err.Params()))
		case errors.Service:
			return c.JSON(http.StatusInternalServerError, common.MResponse(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, common.MResponse(err.Error()))
		}
	}

	return c.JSON(http.StatusOK, CheckResponse{
		Can:            response.Can,
		RemainingDepth: response.RemainingDepth,
		Decisions:      response.Visits,
	})
}

// @Summary     Permission
// @Description expand relationships according to schema
// @ID          permissions.expand
// @Tags  	    Permission
// @Accept      json
// @Produce     json
// @Param       request body ExpandRequest true "expand relationships according to schema"
// @Success     200 {object} ExpandResponse
// @Failure     400 {object} common.HTTPErrorResponse
// @Router      /permissions/expand [post]
func (r *permissionRoutes) expand(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "permissions.expand")
	defer span.End()

	request := new(ExpandRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}
	v := request.Validate()
	if v != nil {
		return c.JSON(http.StatusUnprocessableEntity, common.ValidationResponse(v))
	}

	var err errors.Error

	var response commands.ExpandResponse
	response, err = r.service.Expand(ctx, request.Entity, request.Action, request.SchemaVersion.String())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		switch err.Kind() {
		case errors.Database:
			return c.JSON(database.GetKindToHttpStatus(err.SubKind()), common.MResponse(err.Error()))
		case errors.Validation:
			return c.JSON(http.StatusUnprocessableEntity, common.ValidationResponse(err.Params()))
		case errors.Service:
			return c.JSON(http.StatusInternalServerError, common.MResponse(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, common.MResponse(err.Error()))
		}
	}

	return c.JSON(http.StatusOK, ExpandResponse{
		Tree: response.Tree,
	})
}

// lookupQuery -
func (r *permissionRoutes) lookupQuery(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "permissions.expand")
	defer span.End()

	request := new(LookupQueryRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, &request); err != nil {
		return err
	}
	v := request.Validate()
	if v != nil {
		return c.JSON(http.StatusUnprocessableEntity, common.ValidationResponse(v))
	}

	var err errors.Error

	var response commands.LookupQueryResponse
	response, err = r.service.LookupQuery(ctx, request.EntityType, request.Subject, request.Action, request.SchemaVersion.String())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		switch err.Kind() {
		case errors.Database:
			return c.JSON(database.GetKindToHttpStatus(err.SubKind()), common.MResponse(err.Error()))
		case errors.Validation:
			return c.JSON(http.StatusUnprocessableEntity, common.ValidationResponse(err.Params()))
		case errors.Service:
			return c.JSON(http.StatusInternalServerError, common.MResponse(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, common.MResponse(err.Error()))
		}
	}

	return c.JSON(http.StatusOK, LookupQueryResponse{
		Query: response.Query,
	})
}
