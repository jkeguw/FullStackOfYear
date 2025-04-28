package order

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"project/backend/internal/errors"
	"project/backend/models"
	orderService "project/backend/services/order"
	orderTypes "project/backend/types/order"
)

// Handler 订单处理程序
type Handler struct {
	orderService *orderService.Service
}

// NewHandler 创建订单处理程序
func NewHandler(orderService *orderService.Service) *Handler {
	return &Handler{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (h *Handler) CreateOrder(c *gin.Context) {
	// 获取用户ID
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAppError(errors.Unauthorized, "未授权访问"))
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的用户ID"))
		return
	}

	// 解析请求
	var req orderTypes.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的请求数据"))
		return
	}

	// 创建订单
	orderObj, err := h.orderService.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, err.Error()))
		}
		return
	}

	// 返回响应
	c.JSON(http.StatusCreated, convertToOrderResponse(orderObj))
}

// GetOrder 获取订单详情
func (h *Handler) GetOrder(c *gin.Context) {
	// 获取用户ID
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAppError(errors.Unauthorized, "未授权访问"))
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的用户ID"))
		return
	}

	// 获取订单ID
	orderIDStr := c.Param("id")
	orderID, err := primitive.ObjectIDFromHex(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的订单ID"))
		return
	}

	// 获取订单
	orderObj, err := h.orderService.GetOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, err.Error()))
		}
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, convertToOrderResponse(orderObj))
}

// GetOrderByNumber 通过订单号获取订单
func (h *Handler) GetOrderByNumber(c *gin.Context) {
	// 获取用户ID
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAppError(errors.Unauthorized, "未授权访问"))
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的用户ID"))
		return
	}

	// 获取订单号
	orderNumber := c.Param("number")
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "订单号不能为空"))
		return
	}

	// 获取订单
	orderObj, err := h.orderService.GetOrderByNumber(c.Request.Context(), userID, orderNumber)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, err.Error()))
		}
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, convertToOrderResponse(orderObj))
}

// ListUserOrders 获取用户订单列表
func (h *Handler) ListUserOrders(c *gin.Context) {
	// 获取用户ID
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAppError(errors.Unauthorized, "未授权访问"))
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的用户ID"))
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 10

	pageStr := c.Query("page")
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	pageSizeStr := c.Query("page_size") // 修改为 page_size 以匹配前端API
	if pageSizeStr == "" {
		pageSizeStr = c.Query("pageSize") // 后备兼容
	}
	if pageSizeStr != "" {
		pageSizeInt, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSizeInt > 0 && pageSizeInt <= 100 {
			pageSize = pageSizeInt
		}
	}

	// 获取订单列表
	orderList, err := h.orderService.ListUserOrders(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, err.Error()))
		}
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, orderList)
}

// UpdateOrderStatus 更新订单状态
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	// 获取用户ID
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAppError(errors.Unauthorized, "未授权访问"))
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的用户ID"))
		return
	}

	// 获取订单ID
	orderIDStr := c.Param("id")
	orderID, err := primitive.ObjectIDFromHex(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的订单ID"))
		return
	}

	// 解析请求
	var req orderTypes.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的请求数据"))
		return
	}

	// 验证状态
	status := models.OrderStatusEnum(req.Status)
	if !isValidOrderStatus(status) {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的订单状态"))
		return
	}

	// 更新状态
	orderObj, err := h.orderService.UpdateOrderStatus(c.Request.Context(), userID, orderID, status, req.CancelReason)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, err.Error()))
		}
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, convertToOrderResponse(orderObj))
}

// ProcessPayment 处理支付
func (h *Handler) ProcessPayment(c *gin.Context) {
	// 获取订单ID
	orderIDStr := c.Param("id")
	orderID, err := primitive.ObjectIDFromHex(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的订单ID"))
		return
	}

	// 解析请求
	var req orderTypes.PaymentCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "无效的请求数据"))
		return
	}

	// 处理支付
	orderObj, err := h.orderService.ProcessPayment(c.Request.Context(), orderID, req.TransactionID, req.PaymentStatus)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			c.JSON(appErr.HTTPStatus(), appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewAppError(errors.InternalError, err.Error()))
		}
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, convertToOrderResponse(orderObj))
}

// 验证订单状态是否有效
func isValidOrderStatus(status models.OrderStatusEnum) bool {
	validStatus := map[models.OrderStatusEnum]bool{
		models.OrderStatusPending:   true,
		models.OrderStatusPaid:      true,
		models.OrderStatusShipped:   true,
		models.OrderStatusDelivered: true,
		models.OrderStatusCancelled: true,
		models.OrderStatusRefunded:  true,
	}
	return validStatus[status]
}

// 转换订单为响应格式
func convertToOrderResponse(order *models.Order) orderTypes.OrderResponse {
	// 转换订单项
	items := make([]orderTypes.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = orderTypes.OrderItemResponse{
			ProductID:   item.ProductID.Hex(),
			ProductType: item.ProductType,
			Name:        item.Name,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
			ImageURL:    item.ImageURL,
		}
	}

	// 转换配送信息
	shippingInfo := orderTypes.ShippingInfoResponse{
		Name:           order.ShippingInfo.Name,
		Phone:          order.ShippingInfo.Phone,
		Email:          order.ShippingInfo.Email,
		Address:        order.ShippingInfo.Address,
		City:           order.ShippingInfo.City,
		State:          order.ShippingInfo.State,
		ZipCode:        order.ShippingInfo.ZipCode,
		Country:        order.ShippingInfo.Country,
		ShippingMethod: order.ShippingInfo.ShippingMethod,
	}

	// 转换支付信息
	paymentInfo := orderTypes.PaymentInfoResponse{
		Method:          string(order.PaymentInfo.Method),
		TransactionID:   order.PaymentInfo.TransactionID,
		LastFourDigits:  order.PaymentInfo.LastFourDigits,
		PaymentStatus:   order.PaymentInfo.PaymentStatus,
		PaymentProvider: order.PaymentInfo.PaymentProvider,
	}

	return orderTypes.OrderResponse{
		ID:           order.ID.Hex(),
		UserID:       order.UserID.Hex(),
		OrderNumber:  order.OrderNumber,
		Status:       string(order.Status),
		Items:        items,
		ShippingInfo: shippingInfo,
		PaymentInfo:  paymentInfo,
		Subtotal:     order.Subtotal,
		ShippingFee:  order.ShippingFee,
		Tax:          order.Tax,
		Discount:     order.Discount,
		Total:        order.Total,
		Notes:        order.Notes,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		PaidAt:       order.PaidAt,
		ShippedAt:    order.ShippedAt,
		DeliveredAt:  order.DeliveredAt,
		CancelledAt:  order.CancelledAt,
		CancelReason: order.CancelReason,
	}
}
