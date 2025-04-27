package similarity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"project/backend/models"
)

func TestSimilarityService_CompareMice(t *testing.T) {
	service := NewSimilarityService()

	// 创建测试数据
	mouse1 := models.MouseDevice{
		HardwareDevice: models.HardwareDevice{
			ID:        primitive.NewObjectID(),
			Name:      "Mouse 1",
			Brand:     "Brand A",
			Type:      models.DeviceTypeMouse,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Dimensions: models.MouseDimensions{
			Length: 120,
			Width:  65,
			Height: 40,
			Weight: 80,
		},
		Shape: models.MouseShape{
			Type:              "ergo",
			HumpPlacement:     "center",
			FrontFlare:        "medium",
			SideCurvature:     "curved",
			HandCompatibility: "medium",
		},
		Technical: models.MouseTechnical{
			Connectivity: []string{"wired"},
			Sensor:       "PAW3370",
			MaxDPI:       16000,
			PollingRate:  1000,
			SideButtons:  2,
		},
		Recommended: models.MouseRecommended{
			GameTypes:    []string{"FPS", "MOBA"},
			GripStyles:   []string{"palm", "claw"},
			HandSizes:    []string{"medium", "large"},
			DailyUse:     true,
			Professional: true,
		},
	}

	mouse2 := models.MouseDevice{
		HardwareDevice: models.HardwareDevice{
			ID:        primitive.NewObjectID(),
			Name:      "Mouse 2",
			Brand:     "Brand B",
			Type:      models.DeviceTypeMouse,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Dimensions: models.MouseDimensions{
			Length: 125,
			Width:  68,
			Height: 42,
			Weight: 85,
		},
		Shape: models.MouseShape{
			Type:              "ergo",
			HumpPlacement:     "center",
			FrontFlare:        "medium",
			SideCurvature:     "curved",
			HandCompatibility: "medium",
		},
		Technical: models.MouseTechnical{
			Connectivity: []string{"wireless"},
			Sensor:       "PAW3370",
			MaxDPI:       18000,
			PollingRate:  1000,
			SideButtons:  2,
		},
		Recommended: models.MouseRecommended{
			GameTypes:    []string{"FPS", "MOBA"},
			GripStyles:   []string{"palm", "claw"},
			HandSizes:    []string{"medium", "large"},
			DailyUse:     true,
			Professional: true,
		},
	}

	// 测试比较两个非常相似的鼠标
	t.Run("两个相似的鼠标", func(t *testing.T) {
		result, err := service.CompareMice([]models.MouseDevice{mouse1, mouse2})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Mice))
		assert.Greater(t, result.SimilarityScore, 80, "相似度得分应该很高")
		assert.LessOrEqual(t, result.SimilarityScore, 100, "相似度得分不应超过100")
	})

	// 测试比较一个鼠标
	t.Run("单个鼠标", func(t *testing.T) {
		result, err := service.CompareMice([]models.MouseDevice{mouse1})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result.Mice))
		assert.Equal(t, 100, result.SimilarityScore, "单个鼠标的相似度得分应为100")
	})

	// 测试空切片
	t.Run("空切片", func(t *testing.T) {
		result, err := service.CompareMice([]models.MouseDevice{})
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	// 创建一个差异很大的鼠标进行测试
	mouse3 := models.MouseDevice{
		HardwareDevice: models.HardwareDevice{
			ID:        primitive.NewObjectID(),
			Name:      "Mouse 3",
			Brand:     "Brand C",
			Type:      models.DeviceTypeMouse,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Dimensions: models.MouseDimensions{
			Length: 90, // 明显更小
			Width:  50,
			Height: 30,
			Weight: 50, // 明显更轻
		},
		Shape: models.MouseShape{
			Type:              "ambi", // 不同的形状类型
			HumpPlacement:     "none",
			FrontFlare:        "narrow",
			SideCurvature:     "straight",
			HandCompatibility: "small",
		},
		Technical: models.MouseTechnical{
			Connectivity: []string{"bluetooth"},
			Sensor:       "PMW3389",
			MaxDPI:       8000, // 明显更低
			PollingRate:  500, // 更低
			SideButtons:  0, // 没有侧键
		},
		Recommended: models.MouseRecommended{
			GameTypes:    []string{"casual"},
			GripStyles:   []string{"fingertip"},
			HandSizes:    []string{"small"},
			DailyUse:     true,
			Professional: false,
		},
	}

	// 测试比较差异很大的鼠标
	t.Run("差异很大的鼠标", func(t *testing.T) {
		result, err := service.CompareMice([]models.MouseDevice{mouse1, mouse3})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Mice))
		assert.Less(t, result.SimilarityScore, 50, "相似度得分应该很低")
	})
}

func TestSimilarityService_FindSimilarMice(t *testing.T) {
	service := NewSimilarityService()

	// 创建一组测试鼠标
	mouse1 := models.MouseDevice{
		HardwareDevice: models.HardwareDevice{
			ID:        primitive.NewObjectID(),
			Name:      "Target Mouse",
			Brand:     "Brand A",
			Type:      models.DeviceTypeMouse,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Dimensions: models.MouseDimensions{
			Length: 120,
			Width:  65,
			Height: 40,
			Weight: 80,
		},
		Shape: models.MouseShape{
			Type:              "ergo",
			HumpPlacement:     "center",
			FrontFlare:        "medium",
			SideCurvature:     "curved",
			HandCompatibility: "medium",
		},
		Technical: models.MouseTechnical{
			Connectivity: []string{"wired"},
			Sensor:       "PAW3370",
			MaxDPI:       16000,
			PollingRate:  1000,
			SideButtons:  2,
		},
		Recommended: models.MouseRecommended{
			GameTypes:    []string{"FPS", "MOBA"},
			GripStyles:   []string{"palm", "claw"},
			HandSizes:    []string{"medium", "large"},
			DailyUse:     true,
			Professional: true,
		},
	}

	// 相似鼠标
	mouse2 := models.MouseDevice{
		HardwareDevice: models.HardwareDevice{
			ID:        primitive.NewObjectID(),
			Name:      "Similar Mouse 1",
			Brand:     "Brand B",
			Type:      models.DeviceTypeMouse,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Dimensions: models.MouseDimensions{
			Length: 125,
			Width:  68,
			Height: 42,
			Weight: 85,
		},
		Shape: models.MouseShape{
			Type:              "ergo",
			HumpPlacement:     "center",
			FrontFlare:        "medium",
			SideCurvature:     "curved",
			HandCompatibility: "medium",
		},
		Technical: models.MouseTechnical{
			Connectivity: []string{"wireless"},
			Sensor:       "PAW3370",
			MaxDPI:       18000,
			PollingRate:  1000,
			SideButtons:  2,
		},
		Recommended: models.MouseRecommended{
			GameTypes:    []string{"FPS", "MOBA"},
			GripStyles:   []string{"palm", "claw"},
			HandSizes:    []string{"medium", "large"},
			DailyUse:     true,
			Professional: true,
		},
	}

	// 较少相似的鼠标
	mouse3 := models.MouseDevice{
		HardwareDevice: models.HardwareDevice{
			ID:        primitive.NewObjectID(),
			Name:      "Less Similar Mouse",
			Brand:     "Brand C",
			Type:      models.DeviceTypeMouse,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Dimensions: models.MouseDimensions{
			Length: 110,
			Width:  60,
			Height: 38,
			Weight: 70,
		},
		Shape: models.MouseShape{
			Type:              "ergo",
			HumpPlacement:     "back", // 不同位置
			FrontFlare:        "medium",
			SideCurvature:     "curved",
			HandCompatibility: "medium",
		},
		Technical: models.MouseTechnical{
			Connectivity: []string{"wired"},
			Sensor:       "PMW3389", // 不同传感器
			MaxDPI:       12000, // 更低
			PollingRate:  1000,
			SideButtons:  1, // 更少侧键
		},
		Recommended: models.MouseRecommended{
			GameTypes:    []string{"FPS"},
			GripStyles:   []string{"claw"},
			HandSizes:    []string{"medium"},
			DailyUse:     true,
			Professional: false, // 非专业级
		},
	}

	// 完全不同的鼠标
	mouse4 := models.MouseDevice{
		HardwareDevice: models.HardwareDevice{
			ID:        primitive.NewObjectID(),
			Name:      "Very Different Mouse",
			Brand:     "Brand D",
			Type:      models.DeviceTypeMouse,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Dimensions: models.MouseDimensions{
			Length: 90, // 明显更小
			Width:  50,
			Height: 30,
			Weight: 50, // 明显更轻
		},
		Shape: models.MouseShape{
			Type:              "ambi", // 不同形状
			HumpPlacement:     "none",
			FrontFlare:        "narrow",
			SideCurvature:     "straight",
			HandCompatibility: "small",
		},
		Technical: models.MouseTechnical{
			Connectivity: []string{"bluetooth"},
			Sensor:       "PMW3325",
			MaxDPI:       6000, // 明显更低
			PollingRate:  500, // 更低
			SideButtons:  0, // 没有侧键
		},
		Recommended: models.MouseRecommended{
			GameTypes:    []string{"casual"},
			GripStyles:   []string{"fingertip"},
			HandSizes:    []string{"small"},
			DailyUse:     true,
			Professional: false,
		},
	}

	// 所有鼠标
	allMice := []models.MouseDevice{mouse1, mouse2, mouse3, mouse4}

	t.Run("找到相似鼠标", func(t *testing.T) {
		similarMice, err := service.FindSimilarMice(mouse1.ID, allMice, 3)
		assert.NoError(t, err)
		assert.NotNil(t, similarMice)
		assert.Equal(t, 3, len(similarMice))

		// 验证排序 - 最相似的应该排在前面
		// mouse2 应该比 mouse3 更相似，而 mouse3 应该比 mouse4 更相似
		assert.Equal(t, "Similar Mouse 1", similarMice[0].Name)
		assert.Equal(t, "Less Similar Mouse", similarMice[1].Name)
		assert.Equal(t, "Very Different Mouse", similarMice[2].Name)
	})

	t.Run("限制结果数量", func(t *testing.T) {
		similarMice, err := service.FindSimilarMice(mouse1.ID, allMice, 1)
		assert.NoError(t, err)
		assert.NotNil(t, similarMice)
		assert.Equal(t, 1, len(similarMice))
		assert.Equal(t, "Similar Mouse 1", similarMice[0].Name)
	})

	t.Run("未找到目标鼠标", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		similarMice, err := service.FindSimilarMice(nonExistentID, allMice, 3)
		assert.Error(t, err)
		assert.Nil(t, similarMice)
	})
}