package device

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	
	"project/backend/models"
	"project/backend/types/device"
	"project/backend/internal/errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	mouse, err := h.deviceService.GetMouseDevice(c, deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "未找到鼠标设备",
		})
		return
	}

	// 构建SVG文件名
	svgFileName := fmt.Sprintf("%s-%s.svg", mouse.Name, req.View)
	svgFileName = strings.ReplaceAll(svgFileName, " ", "")
	
	// 尝试读取SVG文件
	// 注意: 这里假设SVG文件存放在项目根目录，实际生产环境应使用配置或环境变量
	filePath := path.Join(".", svgFileName)
	svgData, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 如果找不到文件，尝试使用品牌名称
		filePath = path.Join(".", strings.ReplaceAll(mouse.Brand, " ", "") + "-" + req.View + ".svg")
		svgData, err = ioutil.ReadFile(filePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "SVG数据不可用",
			})
			return
		}
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
		mouse, err := h.deviceService.GetMouseDevice(c, objID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": fmt.Sprintf("未找到鼠标设备: %s", deviceID),
			})
			return
		}

		// 构建SVG文件名
		svgFileName := fmt.Sprintf("%s-%s.svg", mouse.Name, req.View)
		svgFileName = strings.ReplaceAll(svgFileName, " ", "")
		
		// 尝试读取SVG文件
		filePath := path.Join(".", svgFileName)
		svgData, err := ioutil.ReadFile(filePath)
		if err != nil {
			// 如果找不到文件，尝试使用品牌名称
			filePath = path.Join(".", strings.ReplaceAll(mouse.Brand, " ", "") + "-" + req.View + ".svg")
			svgData, err = ioutil.ReadFile(filePath)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    http.StatusNotFound,
					"message": fmt.Sprintf("设备 %s 的SVG数据不可用", deviceID),
				})
				return
			}
		}

		// 添加到响应
		svgResponses = append(svgResponses, device.SVGResponse{
			DeviceID:    deviceID,
			DeviceName:  mouse.Name,
			Brand:       mouse.Brand,
			View:        req.View,
			SVGData:     string(svgData),
			Scale:       1.0, // 默认比例
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
	devices, total, err := h.deviceService.ListDevices(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	// 转换为预览格式
	devicePreviews := make([]device.DevicePreview, 0, len(devices))
	for _, d := range devices {
		// 检查是否有对应的SVG文件
		// 这里简化处理，实际应该基于Views参数检查特定视图
		hasSVG := true
		
		if hasSVG {
			devicePreviews = append(devicePreviews, device.DevicePreview{
				ID:          d.ID.Hex(),
				Name:        d.Name,
				Brand:       d.Brand,
				Type:        string(d.Type),
				ImageURL:    d.ImageURL,
				Description: d.Description,
				CreatedAt:   d.CreatedAt,
			})
		}
	}

	// 构建响应
	response := device.SVGListResponse{
		Devices: devicePreviews,
		Total:   total,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取SVG鼠标列表成功",
		"data":    response,
	})
}