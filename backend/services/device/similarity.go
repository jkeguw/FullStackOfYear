package device

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"project/backend/models"
	"project/backend/types/device"
)

// CompareMice 比较鼠标形状和尺寸
func (s *ServiceImpl) CompareMice(ctx context.Context, ids []string) (*device.ComparisonResponse, error) {
	// 先获取鼠标设备
	mice := make([]*models.MouseDevice, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("无效的鼠标ID: %s", id)
		}

		mouse, err := s.GetMouseDevice(ctx, objID)
		if err != nil {
			return nil, err
		}
		mice = append(mice, mouse)
	}

	// 计算差异
	differences := make(map[string]device.PropertyDiff)

	// 处理尺寸差异
	handleDimensionsDiff(mice, differences)

	// 处理形状差异
	handleShapeDiff(mice, differences)

	// 处理技术参数差异
	handleTechnicalDiff(mice, differences)

	// 计算总体相似度分数
	similarityScore := calculateSimilarityScore(differences)

	// 构建响应
	// 确保mice不为空，避免返回null
	if len(mice) == 0 {
		// 返回一个有效的空响应，而不是null
		return &device.ComparisonResponse{
			Mice:            []device.MouseResponse{},
			Differences:     map[string]device.PropertyDiff{},
			SimilarityScore: 0,
		}, nil
	}

	response := &device.ComparisonResponse{
		Mice:            make([]device.MouseResponse, len(mice)),
		Differences:     differences,
		SimilarityScore: similarityScore,
	}

	// 转换鼠标数据格式
	for i, mouse := range mice {
		response.Mice[i] = mapMouseToResponse(mouse)
	}

	return response, nil
}

// FindSimilarMice 根据给定的鼠标ID查找相似的鼠标
func (s *ServiceImpl) FindSimilarMice(ctx context.Context, id string, limit int) (*device.SimilarityResponse, error) {
	// 解析ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的鼠标ID: %s", id)
	}

	// 获取参考鼠标
	reference, err := s.GetMouseDevice(ctx, objID)
	if err != nil {
		return nil, err
	}

	// 获取所有鼠标
	filter := device.DeviceListFilter{
		Type: string(models.DeviceTypeMouse),
	}
	deviceList, err := s.ListDevices(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 转换回鼠标模型
	allMice := make([]*models.MouseDevice, 0, len(deviceList.Devices))
	for _, dev := range deviceList.Devices {
		mouseID, err := primitive.ObjectIDFromHex(dev.ID)
		if err != nil {
			continue
		}
		if mouse, err := s.GetMouseDevice(ctx, mouseID); err == nil {
			allMice = append(allMice, mouse)
		}
	}

	// 计算与每个鼠标的相似度
	similarities := make([]struct {
		Mouse           *models.MouseDevice
		SimilarityScore float64
		KeyDifferences  map[string]device.PropertyDiff
	}, 0, len(allMice))

	for _, mouse := range allMice {
		// 跳过自身
		if mouse.ID == reference.ID {
			continue
		}

		// 计算相似度
		differences := make(map[string]device.PropertyDiff)
		handleMouseCompare(reference, mouse, differences)
		score := calculateSimilarityScore(differences)

		// 提取关键差异
		keyDiffs := extractKeyDifferences(differences, 5)

		similarities = append(similarities, struct {
			Mouse           *models.MouseDevice
			SimilarityScore float64
			KeyDifferences  map[string]device.PropertyDiff
		}{
			Mouse:           mouse,
			SimilarityScore: score,
			KeyDifferences:  keyDiffs,
		})
	}

	// 按相似度排序
	sortSimilarities(similarities)

	// 截取前N个结果
	if limit > len(similarities) {
		limit = len(similarities)
	}

	// 确保有相似的鼠标
	if len(similarities) == 0 {
		// 返回只有参考鼠标但没有相似鼠标的响应
		return &device.SimilarityResponse{
			Reference:   mapMouseToResponse(reference),
			SimilarMice: []device.SimilarMouse{},
		}, nil
	}

	similarities = similarities[:limit]

	// 构建响应
	response := &device.SimilarityResponse{
		Reference:   mapMouseToResponse(reference),
		SimilarMice: make([]device.SimilarMouse, len(similarities)),
	}

	for i, sim := range similarities {
		keyDiffs := make([]device.PropertyDiff, 0, len(sim.KeyDifferences))
		for _, diff := range sim.KeyDifferences {
			keyDiffs = append(keyDiffs, diff)
		}

		response.SimilarMice[i] = device.SimilarMouse{
			Mouse:           mapMouseToResponse(sim.Mouse),
			SimilarityScore: sim.SimilarityScore,
			KeyDifferences:  keyDiffs,
		}
	}

	return response, nil
}

// 将鼠标设备模型转换为API响应类型
func mapMouseToResponse(mouse *models.MouseDevice) device.MouseResponse {
	return device.MouseResponse{
		ID:          mouse.ID.Hex(),
		Name:        mouse.Name,
		Brand:       mouse.Brand,
		Type:        string(mouse.Type),
		ImageURL:    mouse.ImageURL,
		Description: mouse.Description,
		Dimensions:  mouse.Dimensions,
		Shape:       mouse.Shape,
		Technical:   mouse.Technical,
		Recommended: mouse.Recommended,
		CreatedAt:   mouse.CreatedAt,
		UpdatedAt:   mouse.UpdatedAt,
	}
}

// 分别处理各部分差异
func handleDimensionsDiff(mice []*models.MouseDevice, differences map[string]device.PropertyDiff) {
	// 处理尺寸参数: 长度、宽度、高度、重量
	lengthValues := make([]any, len(mice))
	widthValues := make([]any, len(mice))
	heightValues := make([]any, len(mice))
	weightValues := make([]any, len(mice))

	for i, mouse := range mice {
		lengthValues[i] = mouse.Dimensions.Length
		widthValues[i] = mouse.Dimensions.Width
		heightValues[i] = mouse.Dimensions.Height
		weightValues[i] = mouse.Technical.Weight // 使用技术参数中的重量
	}

	// 计算百分比差异
	lengthDiff := calculateNumericDiff(lengthValues)
	widthDiff := calculateNumericDiff(widthValues)
	heightDiff := calculateNumericDiff(heightValues)
	weightDiff := calculateNumericDiff(weightValues)

	// 添加到差异map
	differences["length"] = device.PropertyDiff{
		Property:          "长度 (mm)",
		Values:            lengthValues,
		DifferencePercent: lengthDiff,
	}
	differences["width"] = device.PropertyDiff{
		Property:          "宽度 (mm)",
		Values:            widthValues,
		DifferencePercent: widthDiff,
	}
	differences["height"] = device.PropertyDiff{
		Property:          "高度 (mm)",
		Values:            heightValues,
		DifferencePercent: heightDiff,
	}
	differences["weight"] = device.PropertyDiff{
		Property:          "重量 (g)",
		Values:            weightValues,
		DifferencePercent: weightDiff,
	}
}

func handleShapeDiff(mice []*models.MouseDevice, differences map[string]device.PropertyDiff) {
	// 处理形状参数
	typeValues := make([]any, len(mice))
	humpValues := make([]any, len(mice))
	flareValues := make([]any, len(mice))
	curvatureValues := make([]any, len(mice))
	handCompValues := make([]any, len(mice))

	for i, mouse := range mice {
		typeValues[i] = mouse.Shape.Type
		humpValues[i] = mouse.Shape.HumpPlacement
		flareValues[i] = mouse.Shape.FrontFlare
		curvatureValues[i] = mouse.Shape.SideCurvature
		handCompValues[i] = mouse.Shape.HandCompatibility
	}

	// 计算差异 - 对于字符串值使用相等性比较
	typeDiff := calculateCategoryDiff(typeValues)
	humpDiff := calculateCategoryDiff(humpValues)
	flareDiff := calculateCategoryDiff(flareValues)
	curvatureDiff := calculateCategoryDiff(curvatureValues)
	handCompDiff := calculateCategoryDiff(handCompValues)

	// 添加到差异map
	differences["shape_type"] = device.PropertyDiff{
		Property:          "形状类型",
		Values:            typeValues,
		DifferencePercent: typeDiff,
	}
	differences["hump_placement"] = device.PropertyDiff{
		Property:          "坑位位置",
		Values:            humpValues,
		DifferencePercent: humpDiff,
	}
	differences["front_flare"] = device.PropertyDiff{
		Property:          "前端开叉",
		Values:            flareValues,
		DifferencePercent: flareDiff,
	}
	differences["side_curvature"] = device.PropertyDiff{
		Property:          "侧面曲线",
		Values:            curvatureValues,
		DifferencePercent: curvatureDiff,
	}
	differences["hand_compatibility"] = device.PropertyDiff{
		Property:          "手型适配",
		Values:            handCompValues,
		DifferencePercent: handCompDiff,
	}
}

func handleTechnicalDiff(mice []*models.MouseDevice, differences map[string]device.PropertyDiff) {
	// 处理技术参数
	dpiValues := make([]any, len(mice))
	pollingRateValues := make([]any, len(mice))
	sideButtonsValues := make([]any, len(mice))

	for i, mouse := range mice {
		dpiValues[i] = mouse.Technical.MaxDPI
		pollingRateValues[i] = mouse.Technical.PollingRate
		sideButtonsValues[i] = mouse.Technical.SideButtons
	}

	// 计算差异
	dpiDiff := calculateNumericDiff(dpiValues)
	pollingRateDiff := calculateNumericDiff(pollingRateValues)
	sideButtonsDiff := calculateNumericDiff(sideButtonsValues)

	// 添加到差异map
	differences["max_dpi"] = device.PropertyDiff{
		Property:          "最大DPI",
		Values:            dpiValues,
		DifferencePercent: dpiDiff,
	}
	differences["polling_rate"] = device.PropertyDiff{
		Property:          "轮询率 (Hz)",
		Values:            pollingRateValues,
		DifferencePercent: pollingRateDiff,
	}
	differences["side_buttons"] = device.PropertyDiff{
		Property:          "侧键数量",
		Values:            sideButtonsValues,
		DifferencePercent: sideButtonsDiff,
	}
}

// 单独比较两只鼠标
func handleMouseCompare(mouse1, mouse2 *models.MouseDevice, differences map[string]device.PropertyDiff) {
	mice := []*models.MouseDevice{mouse1, mouse2}
	handleDimensionsDiff(mice, differences)
	handleShapeDiff(mice, differences)
	handleTechnicalDiff(mice, differences)
}

// 计算数值型参数的差异百分比
func calculateNumericDiff(values []any) float64 {
	if len(values) < 2 {
		return 0
	}

	// 找出最大值和最小值
	var min, max float64
	min = math.MaxFloat64
	max = -math.MaxFloat64

	for _, val := range values {
		var num float64
		switch v := val.(type) {
		case int:
			num = float64(v)
		case int32:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		case float64:
			num = v
		default:
			continue
		}

		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}

	// 防止除以零
	if min == 0 {
		if max == 0 {
			return 0 // 所有值都是0，差异为0
		}
		// 最小值为0，使用最大值作为基准
		return 100.0
	}

	// 计算差异百分比
	return ((max - min) / min) * 100.0
}

// 计算分类型参数的差异百分比
func calculateCategoryDiff(values []any) float64 {
	if len(values) < 2 {
		return 0
	}

	// 检查所有值是否相同
	firstVal := values[0]
	for _, val := range values[1:] {
		if val != firstVal {
			return 100.0 // 不同类别，差异100%
		}
	}

	return 0.0 // 所有值相同，差异0%
}

// 计算总体相似度分数
func calculateSimilarityScore(differences map[string]device.PropertyDiff) float64 {
	if len(differences) == 0 {
		return 100.0
	}

	// 定义权重
	weights := map[string]float64{
		// 尺寸参数 (总权重: 0.5)
		"length": 0.15,
		"width":  0.15,
		"height": 0.10,
		"weight": 0.10,

		// 形状参数 (总权重: 0.35)
		"shape_type":         0.07,
		"hump_placement":     0.07,
		"front_flare":        0.07,
		"side_curvature":     0.07,
		"hand_compatibility": 0.07,

		// 技术参数 (总权重: 0.15)
		"max_dpi":      0.05,
		"polling_rate": 0.05,
		"side_buttons": 0.05,
	}

	// 计算加权平均差异
	totalWeight := 0.0
	weightedDiffSum := 0.0

	for prop, diff := range differences {
		weight := weights[prop]
		if weight == 0 {
			weight = 0.05 // 默认权重
		}
		weightedDiffSum += diff.DifferencePercent * weight
		totalWeight += weight
	}

	// 防止除以零
	if totalWeight == 0 {
		return 100.0
	}

	// 计算平均差异百分比
	avgDiffPercent := weightedDiffSum / totalWeight

	// 将差异百分比转换为相似度分数 (100 - 差异百分比)
	// 限制在0-100范围内
	similarityScore := 100.0 - avgDiffPercent
	if similarityScore < 0 {
		similarityScore = 0
	}
	if similarityScore > 100 {
		similarityScore = 100
	}

	return math.Round(similarityScore)
}

// 提取最显著的几个差异
func extractKeyDifferences(differences map[string]device.PropertyDiff, limit int) map[string]device.PropertyDiff {
	// 按差异大小排序
	type diffPair struct {
		key  string
		diff device.PropertyDiff
	}

	pairs := make([]diffPair, 0, len(differences))
	for k, v := range differences {
		pairs = append(pairs, diffPair{k, v})
	}

	// 冒泡排序 (简单实现，数据量小)
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].diff.DifferencePercent < pairs[j].diff.DifferencePercent {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	// 提取前N个
	if limit > len(pairs) {
		limit = len(pairs)
	}

	result := make(map[string]device.PropertyDiff, limit)
	for i := 0; i < limit; i++ {
		result[pairs[i].key] = pairs[i].diff
	}

	return result
}

// 按相似度排序
func sortSimilarities(similarities []struct {
	Mouse           *models.MouseDevice
	SimilarityScore float64
	KeyDifferences  map[string]device.PropertyDiff
}) {
	// 冒泡排序 (简单实现，数据量小)
	for i := 0; i < len(similarities); i++ {
		for j := i + 1; j < len(similarities); j++ {
			if similarities[i].SimilarityScore < similarities[j].SimilarityScore {
				similarities[i], similarities[j] = similarities[j], similarities[i]
			}
		}
	}
}
