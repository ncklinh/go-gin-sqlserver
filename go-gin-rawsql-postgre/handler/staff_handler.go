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
	var staff model.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		writeError(c, http.StatusBadRequest, "", err)
		return
	}
	hashed, err := util.HashPassword(staff.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to insert staff", err)
		return
	}
	staff.Password = hashed

	id, err := repository.InsertStaff(staff)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to insert staff", err)
		return
	}
	writeSuccess(c, http.StatusCreated, "Success", map[string]any{"id": id})
}

func LoginStaff(jwtMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqStaffInfo model.Staff
		if err := c.ShouldBindJSON(&reqStaffInfo); err != nil {
			writeError(c, http.StatusBadRequest, "", err)
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
		accessToken, err := jwtMaker.CreateToken(
			reqStaffInfo.Username,
			time.Hour,
			token.TokenTypeAccessToken,
		)

		if err != nil {
			writeError(c, http.StatusInternalServerError, "", err)
			return
		}
		writeSuccess(c, http.StatusOK, "Success", gin.H{"accessToken": accessToken})
	}
}
