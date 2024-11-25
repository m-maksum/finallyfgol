package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryAPI interface {
	AddCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	GetCategoryList(c *gin.Context)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryRepo service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryRepo}
}

func (ct *categoryAPI) AddCategory(c *gin.Context) {
	var newCategory model.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := ct.categoryService.Store(&newCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "add category success"})
}

func (ct *categoryAPI) UpdateCategory(c *gin.Context) {
	//======================================================================
	//ngecek biar ga tabrakan sama middleware
	if !ct.isAuthenticated(c) {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized"})
		return
	}
	//======================================================================

	idCat, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Category ID"})
		return
	}

	var updateCat model.Category
	if err := c.ShouldBindJSON(&updateCat); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	if idCat != updateCat.ID {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Category ID mismatch"})
		return
	}

	err = ct.categoryService.Update(idCat, updateCat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "category update success"})
}

func (ct *categoryAPI) DeleteCategory(c *gin.Context) {
	//=====================================================================
	//ngecek biar ga tabrakan sama middleware
	if !ct.isAuthenticated(c) {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized"})
		return
	}
	//=====================================================================

	idCat, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Category ID"})
		return
	}

	err = ct.categoryService.Delete(idCat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "category delete success"})
}

func (ct *categoryAPI) GetCategoryByID(c *gin.Context) {
	idCat, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	category, err := ct.categoryService.GetByID(idCat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (ct *categoryAPI) GetCategoryList(c *gin.Context) {
	//ngecek biar ga tabrakan sama middleware
	//=====================================================================
	if !ct.isAuthenticated(c) {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized"})
		return
	}
	//=====================================================================

	catList, err := ct.categoryService.GetList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, catList)
}

//function tambahan untuk nanganin masalah middleware
//tapi ini di set hanya berlaku untuk di category aja ges
func (ct *categoryAPI) isAuthenticated(c *gin.Context) bool {
	token, err := c.Cookie("session_token")
	if err != nil || token == "" {
		return false 
	}

	return true
}