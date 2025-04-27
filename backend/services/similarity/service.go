package similarity

import (
	"math"
	"sort"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"project/backend/models"
)

// ComparisonResult 鼠标比较结果
type ComparisonResult struct {
	Mice            []models.MouseDevice     `json:"mice"`
	Differences     map[string]PropertyDiff  `json:"differences"`
	SimilarityScore int                      `json:"similarityScore"`
}

// PropertyDiff 属性差异
type PropertyDiff struct {
	Property          string  `json:"property"`
	Values            []any   `json:"values"`
	DifferencePercent float64 `json:"differencePercent"`
}

// SimilarityService 相似度计算服务接口
type SimilarityService interface {
	// CompareMice 比较多个鼠标并生成比较结果
	CompareMice(mice []models.MouseDevice) (*ComparisonResult, error)
	
	// FindSimilarMice 寻找与目标鼠标相似的鼠标列表
	FindSimilarMice(targetMouseID primitive.ObjectID, allMice []models.MouseDevice, limit int) ([]models.MouseDevice, error)
}

type similarityService struct {
	// 如果需要依赖项，可以在这里添加
}

// NewSimilarityService 创建新的相似度服务
func NewSimilarityService() SimilarityService {
	return &similarityService{}
}

// CompareMice 比较多个鼠标并生成比较结果
func (s *similarityService) CompareMice(mice []models.MouseDevice) (*ComparisonResult, error) {
	if len(mice) < 1 {
		return nil, errors.New("至少需要一个鼠标进行比较")
	}

	if len(mice) == 1 {
		return &ComparisonResult{
			Mice:            mice,
			Differences:     make(map[string]PropertyDiff),
			SimilarityScore: 100,
		}, nil
	}

	// 计算差异
	differences := make(map[string]PropertyDiff)

	// 尺寸差异 - 长度
	lengths := make([]float64, len(mice))
	for i, mouse := range mice {
		lengths[i] = mouse.Dimensions.Length
	}
	differences["dimensions.length"] = PropertyDiff{
		Property:          "长度",
		Values:            anySlice(lengths),
		DifferencePercent: calculatePercentageDifference(maxFloat64(lengths), minFloat64(lengths)),
	}

	// 尺寸差异 - 宽度
	widths := make([]float64, len(mice))
	for i, mouse := range mice {
		widths[i] = mouse.Dimensions.Width
	}
	differences["dimensions.width"] = PropertyDiff{
		Property:          "宽度",
		Values:            anySlice(widths),
		DifferencePercent: calculatePercentageDifference(maxFloat64(widths), minFloat64(widths)),
	}

	// 尺寸差异 - 高度
	heights := make([]float64, len(mice))
	for i, mouse := range mice {
		heights[i] = mouse.Dimensions.Height
	}
	differences["dimensions.height"] = PropertyDiff{
		Property:          "高度",
		Values:            anySlice(heights),
		DifferencePercent: calculatePercentageDifference(maxFloat64(heights), minFloat64(heights)),
	}

	// 尺寸差异 - 重量
	weights := make([]float64, len(mice))
	for i, mouse := range mice {
		weights[i] = mouse.Dimensions.Weight
	}
	differences["dimensions.weight"] = PropertyDiff{
		Property:          "重量",
		Values:            anySlice(weights),
		DifferencePercent: calculatePercentageDifference(maxFloat64(weights), minFloat64(weights)),
	}

	// 技术参数差异 - 最大DPI
	dpis := make([]int, len(mice))
	for i, mouse := range mice {
		dpis[i] = mouse.Technical.MaxDPI
	}
	differences["technical.maxDPI"] = PropertyDiff{
		Property:          "最大DPI",
		Values:            anySlice(dpis),
		DifferencePercent: calculatePercentageDifference(float64(maxInt(dpis)), float64(minInt(dpis))),
	}

	// 技术参数差异 - 轮询率
	pollingRates := make([]int, len(mice))
	for i, mouse := range mice {
		pollingRates[i] = mouse.Technical.PollingRate
	}
	differences["technical.pollingRate"] = PropertyDiff{
		Property:          "轮询率",
		Values:            anySlice(pollingRates),
		DifferencePercent: calculatePercentageDifference(float64(maxInt(pollingRates)), float64(minInt(pollingRates))),
	}

	// 技术参数差异 - 侧键数量
	sideButtons := make([]int, len(mice))
	for i, mouse := range mice {
		sideButtons[i] = mouse.Technical.SideButtons
	}
	differences["technical.sideButtons"] = PropertyDiff{
		Property:          "侧键数量",
		Values:            anySlice(sideButtons),
		DifferencePercent: calculatePercentageDifference(float64(maxInt(sideButtons)), float64(minInt(sideButtons))),
	}

	// 形状类型差异
	shapeTypes := make([]string, len(mice))
	for i, mouse := range mice {
		shapeTypes[i] = mouse.Shape.Type
	}
	allSame := true
	for i := 1; i < len(shapeTypes); i++ {
		if shapeTypes[i] != shapeTypes[0] {
			allSame = false
			break
		}
	}
	differences["shape.type"] = PropertyDiff{
		Property:          "形状类型",
		Values:            anySlice(shapeTypes),
		DifferencePercent: 0,
	}
	if !allSame {
		differences["shape.type"] = PropertyDiff{
			Property:          "形状类型",
			Values:            anySlice(shapeTypes),
			DifferencePercent: 100,
		}
	}

	// 推荐握持方式
	gripStyles := make([]string, len(mice))
	for i, mouse := range mice {
		gripStyles[i] = joinStrings(mouse.Recommended.GripStyles, ", ")
	}
	differences["recommended.gripStyles"] = PropertyDiff{
		Property:          "推荐握持方式",
		Values:            anySlice(gripStyles),
		DifferencePercent: 0, // 不计算差异百分比
	}

	// 计算总体相似度得分
	dimensionsSimilarity := calculateDimensionsSimilarity(mice)
	shapeSimilarity := calculateShapeSimilarity(mice)
	technicalSimilarity := calculateTechnicalSimilarity(mice)

	// 加权计算总相似度
	similarityScore := int(math.Round(
		(dimensionsSimilarity * 0.4) + (shapeSimilarity * 0.4) + (technicalSimilarity * 0.2),
	))

	return &ComparisonResult{
		Mice:            mice,
		Differences:     differences,
		SimilarityScore: similarityScore,
	}, nil
}

// FindSimilarMice 寻找与目标鼠标相似的鼠标列表
func (s *similarityService) FindSimilarMice(targetMouseID primitive.ObjectID, allMice []models.MouseDevice, limit int) ([]models.MouseDevice, error) {
	if limit <= 0 {
		limit = 5 // 默认限制
	}

	// 找到目标鼠标
	var targetMouse models.MouseDevice
	var otherMice []models.MouseDevice

	for _, mouse := range allMice {
		if mouse.ID == targetMouseID {
			targetMouse = mouse
		} else {
			otherMice = append(otherMice, mouse)
		}
	}

	if targetMouse.ID.IsZero() {
		return nil, errors.New("未找到目标鼠标")
	}

	type mouseSimilarity struct {
		mouse          models.MouseDevice
		similarityScore int
	}

	// 计算每个鼠标与目标鼠标的相似度
	similarities := make([]mouseSimilarity, 0, len(otherMice))
	for _, mouse := range otherMice {
		result, err := s.CompareMice([]models.MouseDevice{targetMouse, mouse})
		if err != nil {
			return nil, errors.Wrap(err, "比较鼠标时出错")
		}

		similarities = append(similarities, mouseSimilarity{
			mouse:          mouse,
			similarityScore: result.SimilarityScore,
		})
	}

	// 按相似度降序排序
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].similarityScore > similarities[j].similarityScore
	})

	// 返回限制数量的结果
	resultCount := limit
	if resultCount > len(similarities) {
		resultCount = len(similarities)
	}

	result := make([]models.MouseDevice, resultCount)
	for i := 0; i < resultCount; i++ {
		result[i] = similarities[i].mouse
	}

	return result, nil
}

// 计算两个数值之间的百分比差异
func calculatePercentageDifference(a, b float64) float64 {
	if a == 0 && b == 0 {
		return 0
	}
	if a == 0 || b == 0 {
		return 100
	}

	max := math.Max(a, b)
	min := math.Min(a, b)
	return ((max - min) / min) * 100
}

// 计算鼠标物理尺寸的相似度得分
func calculateDimensionsSimilarity(mice []models.MouseDevice) float64 {
	if len(mice) < 2 {
		return 100
	}

	// 获取所有维度的平均值
	var avgLength, avgWidth, avgHeight, avgWeight float64
	for _, mouse := range mice {
		avgLength += mouse.Dimensions.Length
		avgWidth += mouse.Dimensions.Width
		avgHeight += mouse.Dimensions.Height
		avgWeight += mouse.Dimensions.Weight
	}
	avgLength /= float64(len(mice))
	avgWidth /= float64(len(mice))
	avgHeight /= float64(len(mice))
	avgWeight /= float64(len(mice))

	// 计算每个鼠标与平均值的偏差
	deviations := make([]float64, len(mice))
	for i, mouse := range mice {
		lengthDev := math.Abs(mouse.Dimensions.Length-avgLength) / avgLength
		widthDev := math.Abs(mouse.Dimensions.Width-avgWidth) / avgWidth
		heightDev := math.Abs(mouse.Dimensions.Height-avgHeight) / avgHeight
		weightDev := math.Abs(mouse.Dimensions.Weight-avgWeight) / avgWeight

		// 返回总偏差
		deviations[i] = (lengthDev + widthDev + heightDev + weightDev) / 4
	}

	// 平均偏差
	var avgDeviation float64
	for _, dev := range deviations {
		avgDeviation += dev
	}
	avgDeviation /= float64(len(deviations))

	// 将偏差转换为相似度得分(0-100)
	return math.Max(0, 100-(avgDeviation*100))
}

// 计算形状相似度得分
func calculateShapeSimilarity(mice []models.MouseDevice) float64 {
	if len(mice) < 2 {
		return 100
	}

	// 计算鼠标形状特征的相似性 (简化版)
	// 检查形状类型、隆起位置等是否一致
	totalMatches := 0.0
	for i := 1; i < len(mice); i++ {
		baseShape := mice[0].Shape
		currentShape := mice[i].Shape

		score := 0.0
		// 形状类型匹配
		if baseShape.Type == currentShape.Type {
			score += 1.0
		}
		// 隆起位置匹配
		if baseShape.HumpPlacement == currentShape.HumpPlacement {
			score += 1.0
		}
		// 前部曲线匹配
		if baseShape.FrontFlare == currentShape.FrontFlare {
			score += 1.0
		}
		// 侧面曲率匹配
		if baseShape.SideCurvature == currentShape.SideCurvature {
			score += 1.0
		}
		// 手型兼容性匹配
		if baseShape.HandCompatibility == currentShape.HandCompatibility {
			score += 1.0
		}

		// 一个鼠标最多5分(所有属性都匹配)
		totalMatches += (score / 5.0)
	}

	// 计算总分数 (满分是每个鼠标都与第一个完全匹配)
	return (totalMatches / float64(len(mice)-1)) * 100
}

// 计算技术参数相似度
func calculateTechnicalSimilarity(mice []models.MouseDevice) float64 {
	if len(mice) < 2 {
		return 100
	}

	// 获取平均DPI和轮询率
	var avgDPI, avgPollingRate float64
	for _, mouse := range mice {
		avgDPI += float64(mouse.Technical.MaxDPI)
		avgPollingRate += float64(mouse.Technical.PollingRate)
	}
	avgDPI /= float64(len(mice))
	avgPollingRate /= float64(len(mice))

	// 计算每个鼠标与平均值的偏差
	deviations := make([]float64, len(mice))
	for i, mouse := range mice {
		dpiDev := math.Abs(float64(mouse.Technical.MaxDPI)-avgDPI) / avgDPI
		pollingDev := math.Abs(float64(mouse.Technical.PollingRate)-avgPollingRate) / avgPollingRate

		// 加权平均(DPI差异权重较低，因为它往往有较大的数值范围)
		deviations[i] = (dpiDev * 0.3) + (pollingDev * 0.7)
	}

	// 平均偏差
	var avgDeviation float64
	for _, dev := range deviations {
		avgDeviation += dev
	}
	avgDeviation /= float64(len(deviations))

	// 将偏差转换为相似度得分(0-100)
	return math.Max(0, 100-(avgDeviation*100))
}

// 工具函数 - 将任意类型切片转换为any切片
func anySlice[T any](slice []T) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// 工具函数 - 找出浮点数切片中的最大值
func maxFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// 工具函数 - 找出浮点数切片中的最小值
func minFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// 工具函数 - 找出整数切片中的最大值
func maxInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// 工具函数 - 找出整数切片中的最小值
func minInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// 工具函数 - 将字符串切片用分隔符连接
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}