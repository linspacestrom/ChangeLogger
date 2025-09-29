package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linspacestrom/ChangeLogger/internal/domain"
	"github.com/linspacestrom/ChangeLogger/internal/services"
)

type Server struct {
	engine *gin.Engine
	svc    *services.ProjectService
}

func NewRoutes(svc *services.ProjectService) *Server {
	e := gin.Default()
	s := &Server{engine: e, svc: svc}
	s.RegisterRoutes()
	return s
}

func (s *Server) RegisterRoutes() {
	s.engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	s.engine.GET("/projects", s.listProjects)
	s.engine.GET("/projects/:id", s.getProject)
	s.engine.POST("/projects", s.createProject)
	s.engine.PATCH("/projects/:id", s.updateProject)
	s.engine.DELETE("/projects/:id", s.deleteProject)
}

func (s *Server) Handler() http.Handler { return s.engine }

// listProjects godoc
// @Summary      Get list of projects
// @Description  Returns all projects
// @Tags         projects
// @Produce      json
// @Success      200  {array}   domain.Project
// @Failure      500  {object}  map[string]string
// @Router       /projects [get]
func (s *Server) listProjects(c *gin.Context) {
	projects, err := s.svc.ListProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// getProject godoc
// @Summary 	Get project by ID
// @Description Return project by ID
// @Tags 		projects
// Produce 		json
// @Param 		id   path	string true "Project ID"
// @Success 	200 {object}  domain.Project
// @Failure 	404 {object}  map[string]string
// @Router 		/projects/{id}  [get]
func (s *Server) getProject(c *gin.Context) {
	id := c.Param("id")
	project, err := s.svc.GetProject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project with id %s not found", id)})
		return
	}

	c.JSON(http.StatusOK, project)
}

// createProject godoc
// @Summary      Create project
// @Description  Creates a new project with given title
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project  body      domain.ProjectCreateOrUpdate  true  "Project data"
// @Success      201  {object}  domain.Project
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects [post]
func (s *Server) createProject(c *gin.Context) {
	var reqData domain.ProjectCreateOrUpdate

	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := s.svc.CreateProject(c.Request.Context(), reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// updateProject godoc
// @Summary      Update project
// @Description  Update project with given title
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id       path      string                       true  "Project ID"
// @Param        project  body      domain.ProjectCreateOrUpdate true  "Project data"
// @Success      200  {object}  domain.Project
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id} [patch]
func (s *Server) updateProject(c *gin.Context) {
	id := c.Param("id")

	var reqData domain.ProjectCreateOrUpdate

	if err := c.ShouldBindBodyWithJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := s.svc.UpdateProject(c.Request.Context(), id, reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

// deleteProject godoc
// @Summary      Delete project
// @Description  Delete project by id
// @Tags         projects
// @Param 		 id   path	string true "Project ID"
// @Success      204  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /projects [delete]
func (s *Server) deleteProject(c *gin.Context) {
	id := c.Param("id")
	err := s.svc.DeleteProject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project with id %s not found", id)})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": fmt.Sprintf("Project with id %s succsses deleted", id)})

}
