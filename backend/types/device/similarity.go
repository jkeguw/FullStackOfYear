package device

// ComparisonRequest 鼠标比较请求
type ComparisonRequest struct {
	IDs []string `json:"ids" form:"ids" binding:"required,min=1"`
}

// PropertyDiff 属性差异
type PropertyDiff struct {
	Property          string `json:"property"`
	Values            []any  `json:"values"`
	DifferencePercent float64 `json:"differencePercent"`
}

// ComparisonResult 鼠标比较结果
type ComparisonResult struct {
	Mice            []any                  `json:"mice"`
	Differences     map[string]PropertyDiff `json:"differences"`
	SimilarityScore int                    `json:"similarityScore"`
}

// ComparisonResponse 鼠标比较响应
type ComparisonResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    ComparisonResult `json:"data"`
}

// SimilarityRequest 相似度请求
type SimilarityRequest struct {
	MouseID string `json:"mouseId" form:"mouseId" binding:"required"`
	Limit   int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=20"`
}

// SimilarityResponse 相似度响应
type SimilarityResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}