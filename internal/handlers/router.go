package handlers

import (
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

	s.engine.POST("/projects", s.createProject)
	s.engine.GET("/projects", s.listProjects)
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
	}

	c.JSON(http.StatusOK, projects)
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
	}
	project, err := s.svc.CreateProject(c.Request.Context(), reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, project)
}
