package device

// 鼠标比较相关类型

// ComparisonRequest 鼠标比较请求
type ComparisonRequest struct {
	IDs []string `form:"ids" binding:"required,min=2,max=3"`
}

// SimilarityRequest 相似度请求
type SimilarityRequest struct {
	ID    string `uri:"id" binding:"required"`
	Limit int    `form:"limit" binding:"omitempty,min=1,max=20"`
}

// PropertyDiff 属性差异
type PropertyDiff struct {
	Property         string      `json:"property"`
	Values           []any       `json:"values"`
	DifferencePercent float64    `json:"differencePercent"`
}

// ComparisonResponse 比较响应
type ComparisonResponse struct {
	Mice            []MouseResponse          `json:"mice"`
	Differences     map[string]PropertyDiff  `json:"differences"`
	SimilarityScore float64                  `json:"similarityScore"`
}

// SimilarityResponse 相似度响应
type SimilarityResponse struct {
	Reference       MouseResponse            `json:"reference"`
	SimilarMice     []SimilarMouse           `json:"similarMice"`
}

// SimilarMouse 相似鼠标
type SimilarMouse struct {
	Mouse           MouseResponse            `json:"mouse"`
	SimilarityScore float64                  `json:"similarityScore"`
	KeyDifferences  []PropertyDiff           `json:"keyDifferences"`
}

// SVG相关类型

// SVGResponse SVG响应
type SVGResponse struct {
	DeviceID    string    `json:"deviceId"`
	DeviceName  string    `json:"deviceName"`
	Brand       string    `json:"brand"`
	View        string    `json:"view"`
	SVGData     string    `json:"svgData"`
	Scale       float64   `json:"scale,omitempty"`
}

// SVGRequest SVG请求
type SVGRequest struct {
	ID   string `uri:"id" binding:"required"`
	View string `form:"view" binding:"omitempty,oneof=top side"` 
}

// SVGCompareRequest SVG比较请求
type SVGCompareRequest struct {
	DeviceIDs []string `json:"deviceIds" binding:"required,min=2,max=3"`
	View      string   `json:"view" binding:"required,oneof=top side"`
}

// SVGCompareResponse SVG比较响应
type SVGCompareResponse struct {
	Devices []SVGResponse `json:"devices"`
	Scale   float64       `json:"scale"`
}

// SVGListRequest SVG列表请求
type SVGListRequest struct {
	Type   string    `form:"type" binding:"omitempty"`
	Brand  string    `form:"brand" binding:"omitempty"`
	Views  []string  `form:"views" binding:"omitempty"`
}

// SVGListResponse SVG列表响应
type SVGListResponse struct {
	Devices []DevicePreview `json:"devices"`
	Total   int             `json:"total"`
}