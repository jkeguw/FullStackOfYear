package device

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"project/backend/internal/errors"
	"project/backend/models"
	deviceSvc "project/backend/services/device"
	deviceTypes "project/backend/types/device"
)

type Handler struct {
	deviceService deviceSvc.Service
}

func NewHandler(service deviceSvc.Service) *Handler {
	return &Handler{
		deviceService: service,
	}
}

// CreateMouseDevice 创建鼠标设备
func (h *Handler) CreateMouseDevice(c *gin.Context) {
	var request deviceTypes.CreateMouseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errObj := errors.NewBadRequestError("无效的请求: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 只有管理员可以创建设备
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		errObj := errors.NewForbiddenError("只有管理员可以创建设备")
		c.JSON(http.StatusForbidden, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.deviceService.CreateMouseDevice(c.Request.Context(), request)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 构建响应
	response := mapMouseDeviceToResponse(result)
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "成功",
		"data":    response,
	})
}

// GetMouseDevice 获取鼠标设备详情
func (h *Handler) GetMouseDevice(c *gin.Context) {
	deviceID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(deviceID)
	if err != nil {
		errObj := errors.NewBadRequestError("无效的设备ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.deviceService.GetMouseDevice(c.Request.Context(), objID)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 构建响应
	response := mapMouseDeviceToResponse(result)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    response,
	})
}

// UpdateMouseDevice 更新鼠标设备
func (h *Handler) UpdateMouseDevice(c *gin.Context) {
	deviceID := c.Param("id")

	var request deviceTypes.UpdateMouseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errObj := errors.NewBadRequestError("无效的请求: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 只有管理员可以更新设备
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		errObj := errors.NewForbiddenError("只有管理员可以更新设备")
		c.JSON(http.StatusForbidden, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.deviceService.UpdateMouseDevice(c.Request.Context(), deviceID, request)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 构建响应
	response := mapMouseDeviceToResponse(result)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    response,
	})
}

// DeleteDevice 删除设备
func (h *Handler) DeleteDevice(c *gin.Context) {
	deviceID := c.Param("id")

	// 只有管理员可以删除设备
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		errObj := errors.NewForbiddenError("只有管理员可以删除设备")
		c.JSON(http.StatusForbidden, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	err := h.deviceService.DeleteDevice(c.Request.Context(), deviceID)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    nil,
	})
}

// ListDevices 获取设备列表
func (h *Handler) ListDevices(c *gin.Context) {
	var request deviceTypes.DeviceListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的请求参数: "+err.Error()))
		return
	}

	// 创建DeviceListFilter
	filter := deviceTypes.DeviceListFilter{
		Page:     request.Page,
		PageSize: request.PageSize,
		Type:     request.Type,
		Brand:    request.Brand,
	}
	result, err := h.deviceService.ListDevices(c.Request.Context(), filter)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 确保结果不为空
	if result == nil {
		result = &deviceTypes.DeviceListResponse{
			Devices: []deviceTypes.DevicePreview{},
			Total:   0,
		}
	}

	// 打印更详细的调试信息
	log.Printf("Sending devices response: %+v, filter: %+v", result, filter)

	// 统一响应格式
	fmt.Printf("DEBUG: ListDevices result: %+v\n", result)

	if result == nil || (result.Devices == nil && result.Total == 0) {
		result = &deviceTypes.DeviceListResponse{
			Devices:  []deviceTypes.DevicePreview{},
			Total:    0,
			Page:     filter.Page,
			PageSize: filter.PageSize,
		}
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    result,
	})
}

// CreateDeviceReview 创建设备评测（仅限评测员和管理员）
func (h *Handler) CreateDeviceReview(c *gin.Context) {
	var request deviceTypes.CreateDeviceReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errObj := errors.NewBadRequestError("无效的请求: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errObj := errors.NewUnauthorizedError("用户未认证")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	role, roleExists := c.Get("userRole")
	if !roleExists || (role.(string) != string(models.RoleReviewer) && role.(string) != string(models.RoleAdmin)) {
		errObj := errors.NewForbiddenError("只有评测员才能创建评测")
		c.JSON(http.StatusForbidden, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 转换请求类型为CreateReviewRequest
	createReviewReq := deviceTypes.CreateReviewRequest(request)

	// 调用服务创建评测
	result, err := h.deviceService.CreateDeviceReview(c.Request.Context(), userID.(string), createReviewReq)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 构建响应
	response := mapDeviceReviewToResponse(result)
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "成功",
		"data":    response,
	})
}

// GetDeviceReview 获取评测详情
func (h *Handler) GetDeviceReview(c *gin.Context) {
	reviewID := c.Param("id")

	result, err := h.deviceService.GetDeviceReview(c.Request.Context(), reviewID)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 构建响应
	response := mapDeviceReviewToResponse(result)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    response,
	})
}

// UpdateDeviceReview 更新评测
func (h *Handler) UpdateDeviceReview(c *gin.Context) {
	reviewID := c.Param("id")

	var request deviceTypes.UpdateDeviceReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errObj := errors.NewBadRequestError("无效的请求: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errObj := errors.NewUnauthorizedError("用户未认证")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 转换请求类型为UpdateReviewRequest
	updateReviewReq := deviceTypes.UpdateReviewRequest(request)

	// 调用服务更新评测
	result, err := h.deviceService.UpdateDeviceReview(c.Request.Context(), userID.(string), reviewID, updateReviewReq)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 构建响应
	response := mapDeviceReviewToResponse(result)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    response,
	})
}

// DeleteDeviceReview 删除评测
func (h *Handler) DeleteDeviceReview(c *gin.Context) {
	reviewID := c.Param("id")

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errObj := errors.NewUnauthorizedError("用户未认证")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	err := h.deviceService.DeleteDeviceReview(c.Request.Context(), userID.(string), reviewID)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    nil,
	})
}

// ListDeviceReviews 获取评测列表
func (h *Handler) ListDeviceReviews(c *gin.Context) {
	var request deviceTypes.DeviceReviewListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		errObj := errors.NewBadRequestError("无效的请求参数: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.deviceService.ListDeviceReviews(c.Request.Context(), request)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    result,
	})
}

// CreateUserDevice 创建用户设备配置
func (h *Handler) CreateUserDevice(c *gin.Context) {
	var request deviceTypes.CreateUserDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errObj := errors.NewBadRequestError("无效的请求: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errObj := errors.NewUnauthorizedError("用户未认证")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.deviceService.CreateUserDevice(c.Request.Context(), userID.(string), request)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 获取设备详情需要通过service返回
	userDeviceResponse, err := h.mapUserDeviceToResponse(c, result)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "成功",
		"data":    userDeviceResponse,
	})
}

// GetUserDevice 获取用户设备配置详情
func (h *Handler) GetUserDevice(c *gin.Context) {
	configID := c.Param("id")

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errObj := errors.NewUnauthorizedError("用户未认证")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.deviceService.GetUserDevice(c.Request.Context(), userID.(string), configID)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 获取设备详情
	userDeviceResponse, err := h.mapUserDeviceToResponse(c, result)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    userDeviceResponse,
	})
}

// UpdateUserDevice 更新用户设备配置
func (h *Handler) UpdateUserDevice(c *gin.Context) {
	configID := c.Param("id")

	var request deviceTypes.CreateUserDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errObj := errors.NewBadRequestError("无效的请求: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errObj := errors.NewUnauthorizedError("用户未认证")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 转换请求类型为UpdateUserDeviceRequest
	updateRequest := deviceTypes.UpdateUserDeviceRequest{
		Name:        request.Name,
		Description: request.Description,
		Devices:     request.Devices,
		IsPublic:    request.IsPublic,
	}

	// 调用服务更新用户设备配置
	result, err := h.deviceService.UpdateUserDevice(c.Request.Context(), userID.(string), configID, updateRequest)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 获取设备详情
	userDeviceResponse, err := h.mapUserDeviceToResponse(c, result)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    userDeviceResponse,
	})
}

// DeleteUserDevice 删除用户设备配置
func (h *Handler) DeleteUserDevice(c *gin.Context) {
	configID := c.Param("id")

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		errObj := errors.NewUnauthorizedError("用户未认证")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	err := h.deviceService.DeleteUserDevice(c.Request.Context(), userID.(string), configID)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    nil,
	})
}

// ListUserDevices 获取用户设备配置列表
func (h *Handler) ListUserDevices(c *gin.Context) {
	var request deviceTypes.UserDeviceListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		errObj := errors.NewBadRequestError("无效的请求参数: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if exists && request.UserID == "" {
		// 如果没有指定用户ID，默认查询当前用户的配置
		request.UserID = userID.(string)
	}

	// 非管理员只能查看自己的或公开的配置
	role, roleExists := c.Get("userRole")
	if roleExists && role.(string) != "admin" && request.UserID != userID.(string) {
		// 非管理员要查看其他用户的配置，只能查看公开的
		isPublic := true
		request.IsPublic = &isPublic
	}

	result, err := h.deviceService.ListUserDevices(c.Request.Context(), request)
	if err != nil {
		status := errors.HTTPStatusFromError(err)
		var errObj *errors.AppError
		var ok bool
		if errObj, ok = err.(*errors.AppError); !ok {
			// 如果不是AppError类型，创建一个通用错误
			errObj = errors.NewInternalServerError("服务器内部错误")
		}
		c.JSON(status, gin.H{
			"code":    errObj.Code,
			"message": errObj.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    result,
	})
}

// 辅助方法

// mapMouseDeviceToResponse 将鼠标设备模型转换为响应
func mapMouseDeviceToResponse(device *models.MouseDevice) deviceTypes.MouseResponse {
	return deviceTypes.MouseResponse{
		ID:          device.ID.Hex(),
		Name:        device.Name,
		Brand:       device.Brand,
		Type:        string(device.Type),
		ImageURL:    device.ImageURL,
		Description: device.Description,
		Dimensions:  device.Dimensions,
		Shape:       device.Shape,
		Technical:   device.Technical,
		Recommended: device.Recommended,
		CreatedAt:   device.CreatedAt,
		UpdatedAt:   device.UpdatedAt,
	}
}

// mapDeviceReviewToResponse 将设备评测模型转换为响应
func mapDeviceReviewToResponse(review *models.DeviceReview) deviceTypes.DeviceReviewResponse {
	return deviceTypes.DeviceReviewResponse{
		ID:        review.ID.Hex(),
		DeviceID:  review.DeviceID.Hex(),
		UserID:    review.UserID.Hex(),
		Content:   review.Content,
		Pros:      review.Pros,
		Cons:      review.Cons,
		Score:     review.Score,
		Usage:     review.Usage,
		Status:    review.Status,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}

// mapUserDeviceToResponse 将用户设备配置模型转换为响应
func (h *Handler) mapUserDeviceToResponse(c *gin.Context, userDevice *models.UserDevice) (deviceTypes.UserDeviceResponse, error) {
	response := deviceTypes.UserDeviceResponse{
		ID:          userDevice.ID.Hex(),
		UserID:      userDevice.UserID.Hex(),
		Name:        userDevice.Name,
		Description: userDevice.Description,
		IsPublic:    userDevice.IsPublic,
		CreatedAt:   userDevice.CreatedAt,
		UpdatedAt:   userDevice.UpdatedAt,
		Devices:     make([]deviceTypes.UserDeviceSettingsResponse, len(userDevice.Devices)),
	}

	// 获取设备详情
	deviceIds := make([]string, len(userDevice.Devices))
	for i, d := range userDevice.Devices {
		deviceIds[i] = d.DeviceID.Hex()
	}

	// 这里需要调用service获取设备详情，但为简化，这里直接返回基本信息
	for i, d := range userDevice.Devices {
		response.Devices[i] = deviceTypes.UserDeviceSettingsResponse{
			DeviceID:    d.DeviceID.Hex(),
			DeviceType:  string(d.DeviceType),
			DeviceName:  "获取中...", // 实际应通过查询设备获取
			DeviceBrand: "获取中...", // 实际应通过查询设备获取
			Settings:    d.Settings,
		}
	}

	// 通过ListUserDevices方法获取完整设备信息
	listRequest := deviceTypes.UserDeviceListRequest{
		UserID: userDevice.UserID.Hex(),
	}
	userDevices, err := h.deviceService.ListUserDevices(c.Request.Context(), listRequest)
	if err != nil {
		return response, err
	}

	// 从结果中查找当前设备配置并获取详细信息
	for _, ud := range userDevices.UserDevices {
		if ud.ID == userDevice.ID.Hex() {
			response.Devices = ud.Devices
			break
		}
	}

	return response, nil
}
