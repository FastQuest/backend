package classroom

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"flashquest/internal/appcontext"
	"flashquest/pkg/apiresp"

	"github.com/gorilla/mux"
)

type Handler struct{ service Service }

func NewHandler(service Service) *Handler { return &Handler{service: service} }

// CreateClassroom creates a classroom for the authenticated professor or admin.
// @Summary Create a classroom
// @Tags Classrooms
// @Accept json
// @Produce json
// @Param classroom body CreateClassroomRequest true "Classroom name"
// @Success 201 {object} ClassroomResponse
// @Failure 400 {object} apiresp.ErrorResponse "Invalid JSON body or classroom name"
// @Failure 401 {object} apiresp.ErrorResponse "Unauthenticated"
// @Failure 403 {object} apiresp.ErrorResponse "Only professors and admins can create classrooms"
// @Failure 500 {object} apiresp.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /classrooms [post]
func (h *Handler) CreateClassroom(w http.ResponseWriter, r *http.Request) {
	var request CreateClassroomRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	userID, role, ok := authContext(w, r)
	if !ok {
		return
	}
	response, err := h.service.CreateClassroom(userID, role, request.Name)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, response, http.StatusCreated)
}

// ListClassrooms lists classrooms created by or actively joined by the authenticated user.
// @Summary List the user's classrooms
// @Tags Classrooms
// @Produce json
// @Success 200 {array} ClassroomResponse
// @Failure 401 {object} apiresp.ErrorResponse "Unauthenticated"
// @Failure 500 {object} apiresp.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /classrooms [get]
func (h *Handler) ListClassrooms(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := authContext(w, r)
	if !ok {
		return
	}
	response, err := h.service.ListClassrooms(userID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, response, http.StatusOK)
}

// UpdateClassroomStatus activates or archives a classroom.
// @Summary Update classroom status
// @Tags Classrooms
// @Accept json
// @Param id path int true "Classroom ID"
// @Param status body UpdateClassroomStatusRequest true "New classroom status (active or archived)"
// @Success 204 "Classroom status updated"
// @Failure 400 {object} apiresp.ErrorResponse "Invalid JSON body, classroom ID, or status"
// @Failure 401 {object} apiresp.ErrorResponse "Unauthenticated"
// @Failure 403 {object} apiresp.ErrorResponse "Only the classroom owner or an admin can update its status"
// @Failure 404 {object} apiresp.ErrorResponse "Classroom not found"
// @Failure 500 {object} apiresp.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /classrooms/{id}/status [put]
func (h *Handler) UpdateClassroomStatus(w http.ResponseWriter, r *http.Request) {
	var request UpdateClassroomStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	userID, role, ok := authContext(w, r)
	if !ok {
		return
	}
	classroomID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := h.service.UpdateClassroomStatus(userID, role, classroomID, request.Status); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AddParticipant adds a user as an active participant in a classroom.
// @Summary Add a classroom participant
// @Tags Classrooms
// @Accept json
// @Produce json
// @Param id path int true "Classroom ID"
// @Param participant body AddParticipantRequest true "User to add"
// @Success 201 {object} ParticipantResponse
// @Failure 400 {object} apiresp.ErrorResponse "Invalid JSON body, classroom ID, or user ID"
// @Failure 401 {object} apiresp.ErrorResponse "Unauthenticated"
// @Failure 403 {object} apiresp.ErrorResponse "Only the classroom owner or an admin can add participants, and the classroom must be active"
// @Failure 404 {object} apiresp.ErrorResponse "Classroom not found"
// @Failure 409 {object} apiresp.ErrorResponse "User is already a participant"
// @Failure 500 {object} apiresp.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /classrooms/{id}/participants [post]
func (h *Handler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	var request AddParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	userID, role, ok := authContext(w, r)
	if !ok {
		return
	}
	classroomID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	response, err := h.service.AddParticipant(userID, role, classroomID, request.UserID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, response, http.StatusCreated)
}

// ListParticipants lists participants in a classroom.
// @Summary List classroom participants
// @Tags Classrooms
// @Produce json
// @Param id path int true "Classroom ID"
// @Success 200 {array} ParticipantResponse
// @Failure 400 {object} apiresp.ErrorResponse "Invalid classroom ID"
// @Failure 401 {object} apiresp.ErrorResponse "Unauthenticated"
// @Failure 403 {object} apiresp.ErrorResponse "Only the classroom owner or an admin can list participants"
// @Failure 404 {object} apiresp.ErrorResponse "Classroom not found"
// @Failure 500 {object} apiresp.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /classrooms/{id}/participants [get]
func (h *Handler) ListParticipants(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := authContext(w, r)
	if !ok {
		return
	}
	classroomID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	response, err := h.service.ListParticipants(userID, role, classroomID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, response, http.StatusOK)
}

// CreateActivity creates an activity from a question set and saves an immutable snapshot.
// @Summary Create a classroom activity
// @Tags Classrooms
// @Accept json
// @Produce json
// @Param id path int true "Classroom ID"
// @Param activity body CreateActivityRequest true "Activity title and question set"
// @Success 201 {object} ActivityResponse
// @Failure 400 {object} apiresp.ErrorResponse "Invalid JSON body, classroom ID, title, or question set ID"
// @Failure 401 {object} apiresp.ErrorResponse "Unauthenticated"
// @Failure 403 {object} apiresp.ErrorResponse "Only the classroom owner or an admin can create activities; classroom or question set is inactive or inaccessible"
// @Failure 404 {object} apiresp.ErrorResponse "Classroom or usable question set not found"
// @Failure 500 {object} apiresp.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /classrooms/{id}/activities [post]
func (h *Handler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var request CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	userID, role, ok := authContext(w, r)
	if !ok {
		return
	}
	classroomID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	response, err := h.service.CreateActivity(userID, role, classroomID, request.QuestionSetID, request.Title)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, response, http.StatusCreated)
}

// ListActivities lists classroom activities with their snapshotted questions and options.
// @Summary List classroom activities
// @Tags Classrooms
// @Produce json
// @Param id path int true "Classroom ID"
// @Success 200 {array} ActivityDetailResponse
// @Failure 400 {object} apiresp.ErrorResponse "Invalid classroom ID"
// @Failure 401 {object} apiresp.ErrorResponse "Unauthenticated"
// @Failure 500 {object} apiresp.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /classrooms/{id}/activities [get]
func (h *Handler) ListActivities(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := authContext(w, r)
	if !ok {
		return
	}
	classroomID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	response, err := h.service.ListActivities(userID, role, classroomID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, response, http.StatusOK)
}

func authContext(w http.ResponseWriter, r *http.Request) (uint, string, bool) {
	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		apiresp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return 0, "", false
	}
	role, _ := appcontext.GetRole(r.Context())
	return userID, role, true
}

func pathID(w http.ResponseWriter, r *http.Request, name string) (uint, bool) {
	value, err := strconv.ParseUint(mux.Vars(r)[name], 10, 32)
	if err != nil || value == 0 {
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid ID")
		return 0, false
	}
	return uint(value), true
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
	case errors.Is(err, ErrInvalidStatus):
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
	case errors.Is(err, ErrForbidden):
		apiresp.WriteError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrInvalidQuestionSet), errors.Is(err, ErrQuestionSetEmpty):
		apiresp.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, ErrQuestionSetForbidden):
		apiresp.WriteError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
	case errors.Is(err, ErrDuplicateParticipant):
		apiresp.WriteError(w, http.StatusConflict, "CONFLICT", err.Error())
	case errors.Is(err, ErrInactiveClassroom), errors.Is(err, ErrInactiveActivity), errors.Is(err, ErrNotActiveParticipant):
		apiresp.WriteError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
	default:
		apiresp.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Request failed")
	}
}

func writeJSON(w http.ResponseWriter, payload any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
