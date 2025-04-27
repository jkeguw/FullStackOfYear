package sensitivity

// 灵敏度相关的类型定义

// ThreeStageCalibrationRequest 三阶校准法请求
type ThreeStageCalibrationRequest struct {
	InitialValue float64 `json:"initialValue"` // 初始输入值 (cm/360°)
	DPI          int     `json:"dpi"`          // 鼠标DPI
	Stage        int     `json:"stage"`        // 校准阶段 (1-3)
	Direction    string  `json:"direction"`    // 选择方向 (left/right)
	CurrentBase  float64 `json:"currentBase"`  // 当前基准值
}

// ThreeStageCalibrationResponse 三阶校准法响应
type ThreeStageCalibrationResponse struct {
	CurrentBase float64                 `json:"currentBase"` // 当前基准值
	LeftValue   float64                 `json:"leftValue"`   // 左值
	RightValue  float64                 `json:"rightValue"`  // 右值
	GameSens    map[string]GameSensData `json:"gameSens"`   // 各游戏灵敏度数据
	History     []CalibrationStep       `json:"history"`     // 校准历史
}

// CalibrationStep 校准步骤
type CalibrationStep struct {
	Stage       int     `json:"stage"`       // 阶段
	Direction   string  `json:"direction"`   // 选择方向
	BaseValue   float64 `json:"baseValue"`   // 基准值
	LeftValue   float64 `json:"leftValue"`   // 左值
	RightValue  float64 `json:"rightValue"`  // 右值
	NewBase     float64 `json:"newBase"`     // 新基准值
	Timestamp   string  `json:"timestamp"`   // 时间戳
}

// BinaryMethodRequest 二分法请求
type BinaryMethodRequest struct {
	InitialValue float64 `json:"initialValue"` // 初始输入值
	CurrentBase  float64 `json:"currentBase"`  // 当前基准值
	CurrentStep  int     `json:"currentStep"`  // 当前步骤
	Choice       string  `json:"choice"`       // 选择 (low/high)
}

// BinaryMethodResponse 二分法响应
type BinaryMethodResponse struct {
	CurrentBase     float64           `json:"currentBase"`     // 当前基准值
	LowValue        float64           `json:"lowValue"`        // 低值
	HighValue       float64           `json:"highValue"`       // 高值
	CurrentStep     int               `json:"currentStep"`     // 当前步骤
	IsComplete      bool              `json:"isComplete"`      // 是否完成
	FinalValue      float64           `json:"finalValue"`      // 最终值(如果完成)
	History         []BinaryStep      `json:"history"`         // 历史记录
}

// BinaryStep 二分法步骤
type BinaryStep struct {
	Step        int     `json:"step"`        // 步骤编号
	BaseValue   float64 `json:"baseValue"`   // 基准值
	LowValue    float64 `json:"lowValue"`    // 低值
	HighValue   float64 `json:"highValue"`   // 高值
	Choice      string  `json:"choice"`      // 选择
	NewBase     float64 `json:"newBase"`     // 新基准值
	Timestamp   string  `json:"timestamp"`   // 时间戳
}

// InterpolationMethodRequest 极敏内推法请求
type InterpolationMethodRequest struct {
	FastValue  float64 `json:"fastValue"`  // 最快可承受灵敏度(k)
	SlowValue  float64 `json:"slowValue"`  // 最慢可承受灵敏度(m)
	TargetType string  `json:"targetType"` // 目标类型 (aiming/gaming)
}

// InterpolationMethodResponse 极敏内推法响应
type InterpolationMethodResponse struct {
	Result      float64 `json:"result"`      // 计算结果
	Formula     string  `json:"formula"`     // 使用的公式
	Calculation string  `json:"calculation"` // 计算过程
}

// SensitivityConversionRequest 灵敏度转换请求
type SensitivityConversionRequest struct {
	DPI         int     `json:"dpi"`         // 鼠标DPI
	Sensitivity float64 `json:"sensitivity"` // 灵敏度数值
	SourceGame  string  `json:"sourceGame"`  // 源游戏
	TargetGame  string  `json:"targetGame"`  // 目标游戏
}

// SensitivityConversionResponse 灵敏度转换响应
type SensitivityConversionResponse struct {
	SourceGame     string       `json:"sourceGame"`     // 源游戏
	SourceValue    float64      `json:"sourceValue"`    // 源灵敏度值
	TargetGame     string       `json:"targetGame"`     // 目标游戏
	TargetValue    float64      `json:"targetValue"`    // 目标灵敏度值
	CM360Value     float64      `json:"cm360Value"`     // cm/360°值
	OtherGames     []GameSens   `json:"otherGames"`     // 其他游戏的灵敏度
}

// GameSens 游戏灵敏度
type GameSens struct {
	Game  string  `json:"game"`  // 游戏名称
	Value float64 `json:"value"` // 灵敏度值
}

// GameSensData 游戏灵敏度数据
type GameSensData struct {
	Game  string  `json:"game"`  // 游戏名称
	Value float64 `json:"value"` // 灵敏度值
}

// SensitivityConstants 灵敏度常量
type SensitivityConstants struct {
	CSGOFactor      float64   `json:"csgoFactor"`      // CS/Apex系数
	ValorantFactor  float64   `json:"valorantFactor"`  // Valorant系数
	OverwatchFactor float64   `json:"overwatchFactor"` // Overwatch系数
	ApexFactor      float64   `json:"apexFactor"`      // Apex系数
	R6Factor        float64   `json:"r6Factor"`        // Rainbow Six系数
	ThreeStageAdjustments []ThreeStageAdjustment `json:"threeStageAdjustments"` // 三阶校准法调整值
	BinaryRatios    []float64 `json:"binaryRatios"`    // 二分法比例系数
}

// ThreeStageAdjustment 三阶校准法调整值
type ThreeStageAdjustment struct {
	Stage        int     `json:"stage"`       // 阶段
	LeftValue    float64 `json:"leftValue"`   // 左值调整
	RightValue   float64 `json:"rightValue"`  // 右值调整
}