package device

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"project/backend/internal/errors"
	"project/backend/models"
	device "project/backend/types/device"
	deviceTypes "project/backend/types/device"
)

// ListPublicUserDevices 获取公开的用户设备配置
func (h *Handler) ListPublicUserDevices(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 调用服务
	result, err := h.deviceService.GetPublicUserDevices(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(errors.HTTPStatusFromError(err), gin.H{"error": err.Error()})
		return
	}

	// 响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    result,
	})
}

// GetMouseSVG 获取鼠标SVG数据
func (h *Handler) GetMouseSVG(c *gin.Context) {
	var req device.SVGRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 设置默认值
	if req.View == "" {
		req.View = "top"
	}

	// 解析设备ID
	deviceID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的设备ID",
		})
		return
	}

	// 获取鼠标设备
	mouse, err := h.deviceService.GetMouseDevice(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "未找到鼠标设备",
		})
		return
	}

	// 检查鼠标SVG数据是否存在
	var svgData []byte
	if mouse.SVGData != nil {
		// 根据请求视图返回对应的SVG数据
		var svgContent string
		if req.View == "top" && mouse.SVGData.TopView != "" {
			svgContent = mouse.SVGData.TopView
			fmt.Printf("使用数据库中的顶视图SVG: %s\n", mouse.Name)
		} else if req.View == "side" && mouse.SVGData.SideView != "" {
			svgContent = mouse.SVGData.SideView
			fmt.Printf("使用数据库中的侧视图SVG: %s\n", mouse.Name)
		}

		// 如果找到了对应视图的SVG
		if svgContent != "" {
			svgData = []byte(svgContent)
		}
	}

	// 如果没有找到SVG数据，返回占位图
	if svgData == nil || len(svgData) == 0 {
		svgData = []byte(`<svg viewBox="0 0 300 150" xmlns="http://www.w3.org/2000/svg">
  <rect width="300" height="150" fill="#eee"/>
  <text x="150" y="75" text-anchor="middle" fill="#999">SVG暂不可用</text>
  <text x="150" y="95" text-anchor="middle" fill="#999" font-size="12">` + mouse.Name + `</text>
</svg>`)
		// 前端调试用日志
		fmt.Printf("已生成替代SVG: %s\n", mouse.Name)
	}

	// 构建响应
	response := device.SVGResponse{
		DeviceID:   req.ID,
		DeviceName: mouse.Name,
		Brand:      mouse.Brand,
		View:       req.View,
		SVGData:    string(svgData),
		Scale:      1.0, // 默认比例
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取SVG数据成功",
		"data":    response,
	})
}

// CompareSVGs 比较多个鼠标的SVG
func (h *Handler) CompareSVGs(c *gin.Context) {
	var req device.SVGCompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	svgResponses := make([]device.SVGResponse, 0, len(req.DeviceIDs))

	// 对每个设备ID获取SVG数据
	for _, deviceID := range req.DeviceIDs {
		// 解析设备ID
		objID, err := primitive.ObjectIDFromHex(deviceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": fmt.Sprintf("无效的设备ID: %s", deviceID),
			})
			return
		}

		// 获取鼠标设备
		mouse, err := h.deviceService.GetMouseDevice(c.Request.Context(), objID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("未找到鼠标设备: %s", deviceID),
			})
			return
		}

		// 检查鼠标是否有SVG数据
		var svgData []byte
		if mouse.SVGData != nil {
			// 根据请求视图返回对应的SVG数据
			var svgContent string
			if req.View == "top" && mouse.SVGData.TopView != "" {
				svgContent = mouse.SVGData.TopView
				fmt.Printf("比较: 使用数据库中的顶视图SVG: %s\n", mouse.Name)
			} else if req.View == "side" && mouse.SVGData.SideView != "" {
				svgContent = mouse.SVGData.SideView
				fmt.Printf("比较: 使用数据库中的侧视图SVG: %s\n", mouse.Name)
			}

			// 如果找到了对应视图的SVG
			if svgContent != "" {
				svgData = []byte(svgContent)
			}
		}

		// 如果没有找到SVG数据，返回占位图
		if svgData == nil || len(svgData) == 0 {
			svgData = []byte(`<svg viewBox="0 0 300 150" xmlns="http://www.w3.org/2000/svg">
  <rect width="300" height="150" fill="#eee"/>
  <text x="150" y="75" text-anchor="middle" fill="#999">SVG暂不可用</text>
  <text x="150" y="95" text-anchor="middle" fill="#999" font-size="12">` + mouse.Name + `</text>
</svg>`)
			// 前端调试用日志
			fmt.Printf("已生成替代SVG用于比较: %s\n", mouse.Name)
		}

		// 添加到响应
		svgResponses = append(svgResponses, device.SVGResponse{
			DeviceID:   deviceID,
			DeviceName: mouse.Name,
			Brand:      mouse.Brand,
			View:       req.View,
			SVGData:    string(svgData),
			Scale:      1.0, // 默认比例
		})
	}

	// 构建最终响应
	response := device.SVGCompareResponse{
		Devices: svgResponses,
		Scale:   1.0, // 默认比例
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "比较SVG数据成功",
		"data":    response,
	})
}

// GetSVGMouseList 获取有SVG数据的鼠标列表
func (h *Handler) GetSVGMouseList(c *gin.Context) {
	var req device.SVGListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 构建查询条件
	filter := device.DeviceListFilter{
		Type:  string(models.DeviceTypeMouse),
		Brand: req.Brand,
	}

	// 获取所有鼠标设备
	result, err := h.deviceService.ListDevices(c.Request.Context(), filter)
	// 防止空指针错误
	if err != nil || result == nil {
		// 返回空列表而不是错误，以便前端仍然可以正常运行
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "获取SVG鼠标列表成功 (空列表)",
			"data": device.SVGListResponse{
				Devices: []device.DevicePreview{},
				Total:   0,
			},
		})
		return
	}

	// 使用结果中的设备列表和总数
	devices := result.Devices
	total := result.Total

	// 构建响应
	response := device.SVGListResponse{
		Devices: devices,
		Total:   total,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取SVG鼠标列表成功",
		"data":    response,
	})
}

// GetDevices 处理公开设备列表查询
func (h *Handler) GetDevices(c *gin.Context) {
	// 从查询参数获取分页信息
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 构建查询条件
	filter := deviceTypes.DeviceListFilter{
		Type:     c.Query("type"),
		Brand:    c.Query("brand"),
		Page:     page,
		PageSize: pageSize,
	}

	// 打印调试信息
	fmt.Printf("GetDevices API: 正在查询设备，类型=%s, 页码=%d\n", filter.Type, filter.Page)

	result, err := h.deviceService.ListDevices(c.Request.Context(), filter)
	if err != nil {
		fmt.Printf("GetDevices API错误: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.InternalError,
			"message": "获取设备列表失败",
			"data":    nil,
		})
		return
	}

	// 确保响应不为空
	if result == nil {
		fmt.Println("GetDevices API: 结果为nil，创建空响应")
		result = &deviceTypes.DeviceListResponse{
			Devices:  []deviceTypes.DevicePreview{},
			Total:    0,
			Page:     page,
			PageSize: pageSize,
		}
	}

	// 构建标准响应
	responseData := gin.H{
		"devices":  result.Devices,
		"total":    result.Total,
		"page":     page,
		"pageSize": pageSize,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功获取设备列表",
		"data":    responseData,
	})
}
