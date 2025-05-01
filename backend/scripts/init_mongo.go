package scripts

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

// MouseData 代表CSV文件中的一行数据
type MouseData struct {
	Name              string
	Length            float64
	Width             float64
	Height            float64
	Weight            float64
	Shape             string
	HumpPlacement     string
	FrontFlare        string
	SideCurvature     string
	HandCompatibility string
	ThumbRest         string
	RingFingerRest    string
	Material          string
	Connectivity      string
	Sensor            string
	SensorTechnology  string
	SensorPosition    string
	DPI               int
	PollingRate       int
	TrackingSpeed     int
	Acceleration      int
	SideButtons       int
	MiddleButtons     int
}

// MouseDevice 数据库中的鼠标设备模型
type MouseDevice struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Brand       string             `bson:"brand"`
	Type        string             `bson:"type"`
	ImageURL    string             `bson:"imageUrl,omitempty"`
	Description string             `bson:"description,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
	DeletedAt   *time.Time         `bson:"deletedAt,omitempty"`
	Dimensions  struct {
		Length float64 `bson:"length"`
		Width  float64 `bson:"width"`
		Height float64 `bson:"height"`
		Weight float64 `bson:"weight"`
	} `bson:"dimensions"`
	Shape struct {
		Type              string `bson:"type"`
		HumpPlacement     string `bson:"humpPlacement"`
		FrontFlare        string `bson:"frontFlare"`
		SideCurvature     string `bson:"sideCurvature"`
		HandCompatibility string `bson:"handCompatibility"`
	} `bson:"shape"`
	Technical struct {
		Connectivity []string `bson:"connectivity"`
		Sensor       string   `bson:"sensor"`
		MaxDPI       int      `bson:"maxDPI"`
		PollingRate  int      `bson:"pollingRate"`
		SideButtons  int      `bson:"sideButtons"`
	} `bson:"technical"`
	Recommended struct {
		GameTypes    []string `bson:"gameTypes"`
		GripStyles   []string `bson:"gripStyles"`
		HandSizes    []string `bson:"handSizes"`
		DailyUse     bool     `bson:"dailyUse"`
		Professional bool     `bson:"professional"`
	} `bson:"recommended"`
}

func InitMongo() {
	// 加载环境变量
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// 获取MongoDB连接URI，如果环境变量中没有，则使用默认值
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		// 尝试使用docker容器的默认连接
		mongoURI = "mongodb://root:example@localhost:27017"
		log.Printf("Using default MongoDB URI: %s", mongoURI)
	} else {
		log.Printf("Using MongoDB URI from environment: %s", mongoURI)
	}

	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// 确认连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB successfully!")

	// 获取数据库和集合
	db := client.Database("cpc")
	collection := db.Collection("devices")

	// 读取CSV文件
	csvPaths := []string{
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_.csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (1).csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (2).csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (3).csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (4).csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (5).csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (6).csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (7).csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (8).csv",
		"/mnt/e/Project/Go/FullStackOfYear/101_x_25_ (9).csv",
	}

	processedMice := make(map[string]bool)
	var successCount int
	var errorCount int

	for _, csvPath := range csvPaths {
		fmt.Printf("Processing file: %s\n", csvPath)

		file, err := os.Open(csvPath)
		if err != nil {
			log.Printf("Error opening CSV file %s: %v", csvPath, err)
			continue
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			log.Printf("Error reading CSV file %s: %v", csvPath, err)
			continue
		}

		// 跳过标题行
		for i, record := range records {
			if i == 0 || len(record) < 23 {
				continue // 跳过标题行或数据不完整的行
			}

			// 解析鼠标数据
			mouseData, err := parseMouseData(record)
			if err != nil {
				log.Printf("Error parsing mouse data at row %d: %v", i+1, err)
				errorCount++
				continue
			}

			// 检查是否已处理过这个鼠标
			if _, exists := processedMice[mouseData.Name]; exists {
				continue
			}
			processedMice[mouseData.Name] = true

			// 创建鼠标设备文档
			mouseDevice := createMouseDevice(mouseData)

			// 插入到数据库
			_, err = collection.InsertOne(ctx, mouseDevice)
			if err != nil {
				log.Printf("Error inserting mouse device into database: %v", err)
				errorCount++
				continue
			}
			successCount++
		}
	}

	fmt.Printf("Successfully imported %d mice into the database.\n", successCount)
	if errorCount > 0 {
		fmt.Printf("Encountered %d errors during import.\n", errorCount)
	}
}

func parseMouseData(record []string) (MouseData, error) {
	var mouse MouseData
	var err error

	mouse.Name = strings.TrimSpace(record[2])
	if mouse.Name == "" {
		return mouse, fmt.Errorf("mouse name is empty")
	}

	// 解析尺寸数据
	mouse.Length, err = parseFloat(record[3])
	if err != nil {
		return mouse, fmt.Errorf("invalid length: %v", err)
	}

	mouse.Width, err = parseFloat(record[4])
	if err != nil {
		return mouse, fmt.Errorf("invalid width: %v", err)
	}

	mouse.Height, err = parseFloat(record[5])
	if err != nil {
		return mouse, fmt.Errorf("invalid height: %v", err)
	}

	mouse.Weight, err = parseFloat(record[6])
	if err != nil {
		// 有些鼠标可能没有重量数据，设为0
		mouse.Weight = 0
	}

	// 形状信息
	mouse.Shape = strings.TrimSpace(record[7])
	mouse.HumpPlacement = strings.TrimSpace(record[8])
	mouse.FrontFlare = strings.TrimSpace(record[9])
	mouse.SideCurvature = strings.TrimSpace(record[10])
	mouse.HandCompatibility = strings.TrimSpace(record[11])
	mouse.ThumbRest = strings.TrimSpace(record[12])
	mouse.RingFingerRest = strings.TrimSpace(record[13])
	mouse.Material = strings.TrimSpace(record[14])
	mouse.Connectivity = strings.TrimSpace(record[15])
	mouse.Sensor = strings.TrimSpace(record[16])
	mouse.SensorTechnology = strings.TrimSpace(record[17])
	mouse.SensorPosition = strings.TrimSpace(record[18])

	// 技术数据
	mouse.DPI, err = parseInt(record[19])
	if err != nil {
		mouse.DPI = 0
	}

	mouse.PollingRate, err = parseInt(record[20])
	if err != nil {
		mouse.PollingRate = 0
	}

	mouse.TrackingSpeed, err = parseInt(record[21])
	if err != nil {
		mouse.TrackingSpeed = 0
	}

	mouse.Acceleration, err = parseInt(record[22])
	if err != nil {
		mouse.Acceleration = 0
	}

	if len(record) > 23 {
		mouse.SideButtons, err = parseInt(record[23])
		if err != nil {
			mouse.SideButtons = 0
		}
	}

	if len(record) > 24 {
		mouse.MiddleButtons, err = parseInt(record[24])
		if err != nil {
			mouse.MiddleButtons = 0
		}
	}

	return mouse, nil
}

func parseFloat(value string) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "-" {
		return 0, nil
	}
	return strconv.ParseFloat(value, 64)
}

func parseInt(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "-" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func createMouseDevice(data MouseData) MouseDevice {
	now := time.Now()
	brandName := extractBrandName(data.Name)

	mouse := MouseDevice{
		ID:        primitive.NewObjectID(),
		Name:      data.Name,
		Brand:     brandName,
		Type:      "mouse",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 设置尺寸信息
	mouse.Dimensions.Length = data.Length
	mouse.Dimensions.Width = data.Width
	mouse.Dimensions.Height = data.Height
	mouse.Dimensions.Weight = data.Weight

	// 设置形状信息
	mouse.Shape.Type = data.Shape
	mouse.Shape.HumpPlacement = data.HumpPlacement
	mouse.Shape.FrontFlare = data.FrontFlare
	mouse.Shape.SideCurvature = data.SideCurvature
	mouse.Shape.HandCompatibility = data.HandCompatibility

	// 设置技术参数
	connectivityList := []string{}
	if data.Connectivity != "" {
		connectivityList = append(connectivityList, data.Connectivity)
	}
	mouse.Technical.Connectivity = connectivityList
	mouse.Technical.Sensor = data.Sensor
	mouse.Technical.MaxDPI = data.DPI
	mouse.Technical.PollingRate = data.PollingRate
	mouse.Technical.SideButtons = data.SideButtons

	// 设置推荐信息
	// 根据形状和性能特征推断适合的游戏类型和握持方式
	mouse.Recommended.GameTypes = inferGameTypes(data)
	mouse.Recommended.GripStyles = inferGripStyles(data)
	mouse.Recommended.HandSizes = inferHandSizes(data)
	mouse.Recommended.DailyUse = inferDailyUse(data)
	mouse.Recommended.Professional = inferProfessional(data)

	return mouse
}

func extractBrandName(fullName string) string {
	// 从完整名称中提取品牌名称
	parts := strings.SplitN(fullName, " ", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return "Unknown"
}

func inferGameTypes(data MouseData) []string {
	gameTypes := []string{}

	// 根据DPI和轮询率推断游戏类型
	if data.DPI >= 16000 && data.PollingRate >= 1000 {
		gameTypes = append(gameTypes, "FPS")
	}

	if data.DPI >= 12000 {
		gameTypes = append(gameTypes, "MOBA")
	}

	if data.SideButtons >= 3 {
		gameTypes = append(gameTypes, "MMO")
	}

	// 确保至少有一个游戏类型
	if len(gameTypes) == 0 {
		gameTypes = append(gameTypes, "General Gaming")
	}

	return gameTypes
}

func inferGripStyles(data MouseData) []string {
	gripStyles := []string{}

	// 根据鼠标形状和尺寸推断握持方式
	if strings.Contains(data.HumpPlacement, "Back") {
		gripStyles = append(gripStyles, "Palm")
	}

	if strings.Contains(data.HumpPlacement, "Center") {
		gripStyles = append(gripStyles, "Claw")
	}

	if data.Weight < 80 && data.Length < 120 {
		gripStyles = append(gripStyles, "Fingertip")
	}

	// 确保至少有一个握持方式
	if len(gripStyles) == 0 {
		gripStyles = append(gripStyles, "Universal")
	}

	return gripStyles
}

func inferHandSizes(data MouseData) []string {
	handSizes := []string{}

	// 根据鼠标尺寸推断适合的手型大小
	if data.Length <= 115 && data.Width <= 60 {
		handSizes = append(handSizes, "Small")
	}

	if (data.Length > 115 && data.Length <= 125) || (data.Width > 60 && data.Width <= 68) {
		handSizes = append(handSizes, "Medium")
	}

	if data.Length > 125 || data.Width > 68 {
		handSizes = append(handSizes, "Large")
	}

	// 确保至少有一个手型大小
	if len(handSizes) == 0 {
		handSizes = append(handSizes, "Medium")
	}

	return handSizes
}

func inferDailyUse(data MouseData) bool {
	// 推断是否适合日常使用
	if strings.Contains(strings.ToLower(data.Connectivity), "wireless") ||
		data.Weight < 100 ||
		strings.Contains(strings.ToLower(data.Shape), "ergonomic") {
		return true
	}
	return false
}

func inferProfessional(data MouseData) bool {
	// 推断是否为专业级
	if data.DPI >= 20000 || data.PollingRate >= 4000 ||
		strings.Contains(strings.ToLower(data.Sensor), "paw3950") ||
		strings.Contains(strings.ToLower(data.Sensor), "paw3395") {
		return true
	}
	return false
}
