// Package rest provides HTTP handlers for the application.
package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/alexm/fuzzy-builder/internal/models"
	"github.com/alexm/fuzzy-builder/internal/repositories"
	"github.com/alexm/fuzzy-builder/internal/services"
)

type Server struct {
	users    *repositories.UserRepository
	passwd   services.PasswordHasher
	jwt      services.JWTIssuer
	projects *repositories.ProjectRepository
	defects  *repositories.DefectRepository
	comments *repositories.CommentRepository
	attach   *repositories.AttachmentRepository
}

func NewServer(users *repositories.UserRepository, projects *repositories.ProjectRepository, defects *repositories.DefectRepository, comments *repositories.CommentRepository, attach *repositories.AttachmentRepository, passwd services.PasswordHasher, jwt services.JWTIssuer) *Server {
	return &Server{users: users, projects: projects, defects: defects, comments: comments, attach: attach, passwd: passwd, jwt: jwt}
}

func (s *Server) Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/health", s.handleHealth).Methods(http.MethodGet)
	r.HandleFunc("/api/auth/register", s.handleRegister).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/login", s.handleLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/users/me", s.handleMe).Methods(http.MethodGet)

	// Projects (protected, manager create/update/delete)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(s.jwtAuth)
	api.HandleFunc("/projects", s.handleProjectsList).Methods(http.MethodGet)
	api.HandleFunc("/projects", s.handleProjectCreate).Methods(http.MethodPost)
	api.HandleFunc("/projects/{id}", s.handleProjectUpdate).Methods(http.MethodPut)
	api.HandleFunc("/projects/{id}", s.handleProjectDelete).Methods(http.MethodDelete)

	// Defects
	api.HandleFunc("/defects", s.handleDefectsList).Methods(http.MethodGet)
	api.HandleFunc("/defects", s.handleDefectCreate).Methods(http.MethodPost)
	api.HandleFunc("/defects/{id}/status", s.handleDefectUpdateStatus).Methods(http.MethodPut)
	api.HandleFunc("/defects/{id}/comments", s.handleDefectAddComment).Methods(http.MethodPost)
	api.HandleFunc("/defects/{id}/comments", s.handleDefectListComments).Methods(http.MethodGet)
	api.HandleFunc("/defects/{id}/attachments", s.handleDefectAddAttachment).Methods(http.MethodPost)
	api.HandleFunc("/defects/{id}/attachments", s.handleDefectListAttachments).Methods(http.MethodGet)
	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	if req.Role != "engineer" && req.Role != "manager" && req.Role != "observer" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}
	// hash password, create user
	hash, err := s.passwd.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "hash error", http.StatusInternalServerError)
		return
	}
	user := &models.User{Email: req.Email, PasswordHash: hash, FullName: req.FullName, Role: models.Role(req.Role)}
	if err := s.users.Create(r.Context(), user); err != nil {
		http.Error(w, "cannot create user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	user, err := s.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := s.passwd.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := s.jwt.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	// For the first pass, read token from Authorization: Bearer <token>
	auth := r.Header.Get("Authorization")
	if len(auth) < 8 || auth[:7] != "Bearer " {
		http.Error(w, "missing bearer token", http.StatusUnauthorized)
		return
	}
	claims, _, err := s.jwt.ParseToken(auth[7:])
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"claims": claims})
}

func (s *Server) handleProjectsList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := int64(20)
	offset := int64(0)
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			offset = n
		}
	}
	items, err := s.projects.List(r.Context(), int32(limit), int32(offset))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list projects")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type projectCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) handleProjectCreate(w http.ResponseWriter, r *http.Request) {
	// Only manager can create in later iteration; for now we'll allow all authenticated
	var req projectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name required")
		return
	}
	p := &models.Project{Name: req.Name, Description: req.Description, CreatedBy: 0}
	if err := s.projects.Create(r.Context(), p); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot create")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

type projectUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) handleProjectUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req projectUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	p := &models.Project{ID: id, Name: req.Name, Description: req.Description}
	if err := s.projects.Update(r.Context(), p); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot update")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handleProjectDelete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.projects.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot delete")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDefectsList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var f repositories.DefectFilter
	if v := q.Get("project_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.ProjectID = &n
		}
	}
	if v := q.Get("assigned_to"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.AssignedTo = &n
		}
	}
	if v := q.Get("status"); v != "" {
		s := models.DefectStatus(v)
		f.Status = &s
	}
	if v := q.Get("priority"); v != "" {
		p := models.DefectPriority(v)
		f.Priority = &p
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			f.Limit = int32(n)
		}
	} else {
		f.Limit = 20
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			f.Offset = int32(n)
		}
	}
	items, err := s.defects.List(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list defects")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type defectCreateRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ProjectID   int64   `json:"project_id"`
	AssignedTo  *int64  `json:"assigned_to"`
	Priority    string  `json:"priority"`
	DueDate     *string `json:"due_date"`
}

func (s *Server) handleDefectCreate(w http.ResponseWriter, r *http.Request) {
	var req defectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Title == "" || req.ProjectID == 0 || req.Priority == "" {
		writeError(w, http.StatusBadRequest, "missing fields")
		return
	}
	priority := models.DefectPriority(req.Priority)
	d := &models.Defect{Title: req.Title, Description: req.Description, ProjectID: req.ProjectID, AssignedTo: req.AssignedTo, Status: models.DefectStatusNew, Priority: priority, CreatedBy: 0}
	if err := s.defects.Create(r.Context(), d); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot create defect")
		return
	}
	writeJSON(w, http.StatusCreated, d)
}

type defectUpdateStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) handleDefectUpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req defectUpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	status := models.DefectStatus(req.Status)
	if err := s.defects.UpdateStatus(r.Context(), id, status); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot update status")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type defectAddCommentRequest struct {
	Text string `json:"text"`
}

func (s *Server) handleDefectAddComment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	defectID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req defectAddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		writeError(w, http.StatusBadRequest, "text required")
		return
	}
	c := &models.Comment{DefectID: defectID, UserID: 0, Text: req.Text}
	if err := s.comments.Create(r.Context(), c); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot add comment")
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (s *Server) handleDefectListComments(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	defectID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	q := r.URL.Query()
	var limit int32 = 20
	var offset int32 = 0
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			limit = int32(n)
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			offset = int32(n)
		}
	}
	items, err := s.comments.ListByDefect(r.Context(), defectID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list comments")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type defectAddAttachmentRequest struct {
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
}

func (s *Server) handleDefectAddAttachment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	defectID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req defectAddAttachmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.Filename) == "" || strings.TrimSpace(req.Filepath) == "" {
		writeError(w, http.StatusBadRequest, "filename and filepath required")
		return
	}
	a := &models.Attachment{DefectID: defectID, Filename: req.Filename, Filepath: req.Filepath, UploadedBy: 0}
	if err := s.attach.Create(r.Context(), a); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot add attachment")
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

func (s *Server) handleDefectListAttachments(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	defectID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	q := r.URL.Query()
	var limit int32 = 20
	var offset int32 = 0
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			limit = int32(n)
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			offset = int32(n)
		}
	}
	items, err := s.attach.ListByDefect(r.Context(), defectID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list attachments")
		return
	}
	writeJSON(w, http.StatusOK, items)
}
