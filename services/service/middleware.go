package service

import (
	"errors"
	"fmt"
	"go-blog/services/store"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func authorization(ctx *gin.Context) {
	// 认证用户

	// 获取认证信息
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing."})
		return
	}

	// Authorization 是否合法
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format is not valid."})
		return
	}

	// Authorization 是否是Bearer类型
	if headerParts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing bearer part."})
		return
	}

	// 从token中提取id
	userID, err := verifyJWT(headerParts[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 用id查询用户
	user, err := store.FetchUser(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 设置上下文属性
	ctx.Set("user", user)
	// 结束中间件
	ctx.Next()
}

func currentUser(ctx *gin.Context) (*store.User, error) {
	// 从上下文中获取用户

	var err error
	// 从上下文获取用户变量
	_user, exists := ctx.Get("user")
	if !exists {
		err = errors.New("Current context user not set")
		log.Error().Err(err).Msg("")
		return nil, err
	}

	// 把interface格式的user转化为User格式的user
	user, ok := _user.(*store.User)
	if !ok {
		err = errors.New("Context user is not valid type")
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return user, nil
}

func customErrors(ctx *gin.Context) {
	// 修改gin内部的报错信息

	// 先正常运行程序
	ctx.Next()

	// 捕捉error
	if len(ctx.Errors) > 0 {
		for _, err := range ctx.Errors {
			// 检查错误类型
			switch err.Type {
			case gin.ErrorTypePublic:
				// Show public errors only if nothing has been written yet
				if !ctx.Writer.Written() {
					ctx.AbortWithStatusJSON(ctx.Writer.Status(), gin.H{"error": err.Error()})
				}
			case gin.ErrorTypeBind:
				errMap := make(map[string]string)
				if errs, ok := err.Err.(validator.ValidationErrors); ok {
					for _, fieldErr := range []validator.FieldError(errs) {
						errMap[fieldErr.Field()] = customValidationError(fieldErr)
					}
				}
				status := http.StatusBadGateway
				// Preserve current status
				if ctx.Writer.Status() != http.StatusOK {
					status = ctx.Writer.Status()
				}
				ctx.AbortWithStatusJSON(status, gin.H{"error": errMap})
			default:
				log.Panic().Err(err).Msg("TypeBind: other error")
			}
		}
		// If there was no public or bind error, display default 500 message
		if !ctx.Writer.Written() {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		}
	}
}

func customValidationError(err validator.FieldError) string {
	// 自定义的验证错误
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required.", err.Field())
	case "min":
		return fmt.Sprintf("%s must be longer than or equal %s characters.", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters.", err.Field(), err.Param())
	default:
		return err.Error()
	}

}
