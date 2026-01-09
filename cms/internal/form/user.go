package form

import (
	"fmt"
	"sort"
	"strings"

	"data_backend/internal/app"
	"data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/database"
)

// LoginRequest 登录请求信息
type LoginRequest struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=2,max=100"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
}

// UserLoginInfo 登录信息返回
type UserLoginInfo struct {
	ID         uint32   `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Token      string   `json:"token"`
	Permission []string `json:"permission"`
	IsAdmin    bool     `json:"is_admin"`
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	UserName   string   `form:"user_name" binding:"required,min=2,max=100"`
	Name       string   `form:"name" binding:"required,min=2,max=100"`
	Email      string   `form:"email" binding:"required,email,min=2,max=100"`
	Password   string   `form:"password" binding:"required,min=6,max=20"`
	IsLock     bool     `form:"is_lock"`
	RoleIDList []uint32 `form:"role_id_list[]"`
}

// UserListRequest 获取用户列表请求
type UserListRequest struct {
	Name     string `form:"name" binding:"max=100"`
	UserName string `form:"user_name" binding:"max=100"`
}

func (q *UserListRequest) Parse() (queryParams database.QueryWhereGroup) {
	if q.Name != "" {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "name LIKE ?",
			Value:  []interface{}{"%" + q.Name + "%"},
		})
	}
	if q.UserName != "" {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "user_name LIKE ?",
			Value:  []interface{}{"%" + q.UserName + "%"},
		})
	}

	return queryParams
}

// User 用户列表返回信息
type User struct {
	ID            uint32      `json:"id"`
	UserName      string      `json:"user_name"`
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	IsLock        uint8       `json:"is_lock"`
	LastLogonTime string      `json:"last_logon_time"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	Role          []*dao.Role `json:"roles,omitempty"`
}

func UserFormat(data []*dao.User, orderBy app.OrderBy) (userList []User, err error) {
	for i := 0; i < len(data); i++ {
		lastLogonTimeStr := "-"
		if data[i].LastLogonTime != nil {
			lastLogonTimeStr = data[i].LastLogonTime.Format(pkg.DATE_TIME_FORMAT)
		}
		userList = append(userList, User{
			ID:            data[i].ID,
			UserName:      data[i].UserName,
			Name:          data[i].Name,
			Email:         data[i].Email,
			IsLock:        data[i].IsLock,
			LastLogonTime: lastLogonTimeStr,
			CreatedAt:     data[i].CreatedAt.Format(pkg.DATE_TIME_FORMAT),
			UpdatedAt:     data[i].UpdatedAt.Format(pkg.DATE_TIME_FORMAT),
			Role:          data[i].Role,
		})
	}

	if orderBy.Field != "" {
		sort.Slice(userList, func(i, j int) bool {
			if orderBy.Order == "desc" || orderBy.Order == "DESC" {
				i, j = j, i
			}
			if orderBy.Field == "roles" {
				return len(userList[i].Role) < len(userList[j].Role)
			} else if orderBy.Field == "created_at" {
				return userList[i].CreatedAt < userList[j].CreatedAt
			}
			return false
		})
	}

	return userList, nil
}

// UserUpdateRequest 更新用户信息请求
type UserUpdateRequest struct {
	Name       string   `form:"name" binding:"required,min=2,max=100"`
	Email      string   `form:"email" binding:"required,email,min=2,max=100"`
	Password   string   `form:"password" binding:"max=20"`
	IsLock     bool     `form:"is_lock"`
	RoleIDList []uint32 `form:"role_id_list[]"`
}

func (q *UserUpdateRequest) Valid() error {
	if q.Password != "" && len(q.Password) < 6 {
		return fmt.Errorf("password must 6 character long")
	}

	return nil
}

type UserUpdateSelfRequest struct {
	Name        string `form:"name" binding:"required,min=2,max=100"`
	Email       string `form:"email" binding:"required,email,min=2,max=100"`
	NewPassword string `form:"new_password" binding:"max=20"`
}

func (q *UserUpdateSelfRequest) Valid() error {
	if q.NewPassword != "" && len(q.NewPassword) < 6 {
		return fmt.Errorf("password must 6 character long")
	}

	return nil
}

func PhoneNumFormat(phoneNum string, canView bool) string {
	if phoneNum == "" {
		return ""
	}
	if canView {
		return phoneNum
	}

	return phoneNum[0:1] + "******" + phoneNum[len(phoneNum)-1-2:]
}

func PagePermissionFormat(permissions []*dao.Permission, userPermissions []string, isAdmin bool) map[string][]map[string]interface{} {
	result := make(map[string][]map[string]interface{})
	for ind := range permissions {
		if !isAdmin && !inUserPermissions(permissions[ind].Name, userPermissions) {
			continue
		}
		pages := strings.Split(permissions[ind].Pages, ",")
		for _, item := range pages {
			result[item] = append(result[item], map[string]interface{}{
				"name":          permissions[ind].Name,
				"id":            permissions[ind].ID,
				"display_name":  permissions[ind].DisplayName,
				"description":   permissions[ind].Description,
				"is_permission": true,
				"is_admin":      isAdmin,
			})
		}

	}

	return result
}

func inUserPermissions(permissions string, userPermissions []string) bool {
	for _, item := range userPermissions {
		if permissions == item {
			return true
		}
	}
	return false
}
