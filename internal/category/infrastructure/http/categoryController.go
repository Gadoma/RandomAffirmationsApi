package http

import (
	"errors"
	"net/http"

	"github.com/gadoma/rafapi/internal/category/domain"
	common "github.com/gadoma/rafapi/internal/common/domain"
	"github.com/gorilla/mux"
)

type CategoryController struct {
	service    domain.CategoryService
	responder  CategoryResponder
	reqHandler *CategoryRequestHandler
}

func NewCategoryController(service domain.CategoryService, responder CategoryResponder, reqHandler *CategoryRequestHandler) *CategoryController {
	return &CategoryController{
		service:    service,
		responder:  responder,
		reqHandler: reqHandler,
	}
}

func (c *CategoryController) RegisterCategoryRoutes(r *mux.Router) {
	r.HandleFunc("/categories", c.handleGetCategories).Methods("GET").Name("getCategories")
	r.HandleFunc("/categories", c.handleCreateCategory).Methods("POST").Name("createCategory")
	r.HandleFunc("/categories/{categoryId:[0-9]+}", c.handleGetCategory).Methods("GET").Name("getCategory")
	r.HandleFunc("/categories/{categoryId:[0-9]+}", c.handleUpdateCategory).Methods("PUT").Name("updateCategory")
	r.HandleFunc("/categories/{categoryId:[0-9]+}", c.handleDeleteCategory).Methods("DELETE").Name("deleteCategory")
}

func (c *CategoryController) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, n, err := c.service.GetCategories(r.Context())

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.responder.RespondSuccessOk(w, categories, n)
}

func (c *CategoryController) handleGetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := c.reqHandler.getCategoryIdParameter(r)

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	categories, err := c.service.GetCategory(r.Context(), id)

	if errors.Is(err, common.ErrorResourceNotFound) {
		c.responder.RespondErrorNotFound(w)
		return
	}

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.responder.RespondSuccessOk(w, categories, 1)
}

func (c *CategoryController) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	au, err := c.reqHandler.getCategoryUpdate(r)

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := c.service.CreateCategory(r.Context(), au)

	if errors.Is(err, domain.ErrorCategoryUpdateInvalidName) {
		c.responder.RespondError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.responder.RespondSuccessOk(w, id, 1)
}

func (c *CategoryController) handleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := c.reqHandler.getCategoryIdParameter(r)

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	au, err := c.reqHandler.getCategoryUpdate(r)

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.UpdateCategory(r.Context(), id, au)

	if errors.Is(err, domain.ErrorCategoryUpdateInvalidName) {
		c.responder.RespondError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.responder.RespondSuccessNoContent(w)
}

func (c *CategoryController) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := c.reqHandler.getCategoryIdParameter(r)

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.DeleteCategory(r.Context(), id)

	if err != nil {
		c.responder.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.responder.RespondSuccessNoContent(w)
}
