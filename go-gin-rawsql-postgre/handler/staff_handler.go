package handler

import (
	"database/sql"
	"film-rental/model"
	"film-rental/repository"
	token "film-rental/token"
	"film-rental/util"
	"film-rental/validator"
	"time"

	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	token_timeout_in_minute = 1
)

func GetStaffs(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 25
	}

	staffs, count, err := repository.GetAllStaff(page, limit)
	pageCount := math.Ceil(float64(count) / float64(limit))

	pagination := PaginationMeta{
		Limit:      limit,
		Page:       page,
		TotalCount: count,
		TotalPage:  int(pageCount),
	}

	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to get staffs", err)
		return
	}
	writeSuccessWithMeta(c, http.StatusOK, "Success", pagination, staffs)
}

func AddStaff(c *gin.Context) {
	var reqStaff model.CreateStaffRequest
	if err := c.ShouldBindJSON(&reqStaff); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate required fields
	message, err := validateStaffFields(reqStaff)
	if message != "" {
		writeError(c, http.StatusBadRequest, message, err)
		return
	}

	// Check if username already exists
	exists, err := repository.IsUsernameExists(reqStaff.Username)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to check username", err)
		return
	}
	if exists {
		writeError(c, http.StatusConflict, "Username already exists", nil)
		return
	}

	hashed, err := util.HashPassword(reqStaff.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	// Convert request to Staff model
	staff := model.Staff{
		FirstName:  reqStaff.FirstName,
		LastName:   reqStaff.LastName,
		AddressId:  reqStaff.AddressId,
		Email:      reqStaff.Email,
		StoreId:    reqStaff.StoreId,
		Active:     reqStaff.Active,
		Username:   reqStaff.Username,
		Password:   hashed,
		Role:       reqStaff.Role,
		Picture:    reqStaff.Picture,
		LastUpdate: time.Now(),
	}

	id, err := repository.InsertStaff(staff)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to insert staff", err)
		return
	}
	writeSuccess(c, http.StatusCreated, "Success", map[string]any{"id": id})
}

func LoginStaff(jwtMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqStaffInfo model.LoginRequest
		if err := c.ShouldBindJSON(&reqStaffInfo); err != nil {
			writeError(c, http.StatusBadRequest, "Invalid request body", err)
			return
		}
		if err := validator.ValidateString(reqStaffInfo.Username, 3, 30); err != nil {
			writeError(c, http.StatusBadRequest, "Username validator", err)
			return
		}
		if err := validator.ValidateString(reqStaffInfo.Password, 6, 30); err != nil {
			writeError(c, http.StatusBadRequest, "Password validator", err)
			return
		}

		staffRecord, err := repository.GetStaff(reqStaffInfo.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(c, http.StatusUnauthorized, "User not found", err)
				return
			}
			writeError(c, http.StatusInternalServerError, "Failed to log in", err)
			return
		}

		if err := util.CheckPassword(reqStaffInfo.Password, staffRecord.Password); err != nil {
			writeError(c, http.StatusUnauthorized, "Invalid email or password", err)
			return
		}

		// Create access token (short-lived, e.g., 15 minutes)
		accessToken, err := jwtMaker.CreateToken(
			reqStaffInfo.Username,
			staffRecord.Role,
			time.Duration(token_timeout_in_minute)*time.Minute,
			token.TokenTypeAccessToken,
		)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "Failed to create access token", err)
			return
		}

		// Create refresh token (long-lived, e.g., 7 days)
		refreshToken, err := jwtMaker.CreateToken(
			reqStaffInfo.Username,
			staffRecord.Role,
			7*24*time.Hour, // 7 days
			token.TokenTypeRefreshToken,
		)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "Failed to create refresh token", err)
			return
		}

		writeSuccess(c, http.StatusOK, "Success", model.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    token_timeout_in_minute * 60, // in seconds
		})
	}
}

func RefreshToken(jwtMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			writeError(c, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if req.RefreshToken == "" {
			writeError(c, http.StatusBadRequest, "Refresh token is required", nil)
			return
		}

		// Verify the refresh token
		payload, err := jwtMaker.VerifyToken(req.RefreshToken, token.TokenTypeRefreshToken)
		if err != nil {
			writeError(c, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}

		// Check if the user still exists in the database
		_, err = repository.GetStaff(payload.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(c, http.StatusUnauthorized, "User not found", err)
				return
			}
			writeError(c, http.StatusInternalServerError, "Failed to verify user", err)
			return
		}

		// Create new access token
		accessToken, err := jwtMaker.CreateToken(
			payload.Username,
			payload.Role,
			time.Duration(token_timeout_in_minute)*time.Minute, // convert int to time.Duration
			token.TokenTypeAccessToken,
		)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "Failed to create access token", err)
			return
		}

		// Create new refresh token (optional - you can reuse the old one or create a new one)
		refreshToken, err := jwtMaker.CreateToken(
			payload.Username,
			payload.Role,
			7*24*time.Hour, // 7 days
			token.TokenTypeRefreshToken,
		)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "Failed to create refresh token", err)
			return
		}

		writeSuccess(c, http.StatusOK, "Token refreshed successfully", model.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    token_timeout_in_minute * 60, // in seconds
		})
	}
}

// validateStaffFields validates all required staff fields
func validateStaffFields(reqStaff model.CreateStaffRequest) (string, error) {
	if reqStaff.FirstName == "" {
		return "First name is required", nil
	}
	if reqStaff.LastName == "" {
		return "Last name is required", nil
	}
	if reqStaff.Email == "" {
		return "Email is required", nil
	}
	if reqStaff.Username == "" {
		return "Username is required", nil
	}
	if reqStaff.Password == "" {
		return "Password is required", nil
	}
	if reqStaff.Role == "" {
		return "Role is required", nil
	}

	if err := validator.ValidateString(reqStaff.Username, 3, 30); err != nil {
		return "Username validation failed", err
	}

	if err := validator.ValidateString(reqStaff.Password, 6, 30); err != nil {
		return "Password validation failed", err
	}

	if !model.IsValidRole(reqStaff.Role) {
		return "Invalid role. Must be 'admin' or 'user'", nil
	}

	return "", nil
}
