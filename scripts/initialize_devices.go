package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MouseSVGData 鼠标SVG数据
type MouseSVGData struct {
	TopView  string `bson:"topView,omitempty" json:"topView,omitempty"`
	SideView string `bson:"sideView,omitempty" json:"sideView,omitempty"`
}

// MouseDimensions 鼠标尺寸信息
type MouseDimensions struct {
	Length float64 `bson:"length" json:"length"`
	Width  float64 `bson:"width" json:"width"`
	Height float64 `bson:"height" json:"height"`
	Weight float64 `bson:"weight" json:"weight"`
}

// MouseShape 鼠标形状信息
type MouseShape struct {
	Type              string `bson:"type" json:"type"`
	HumpPlacement     string `bson:"humpPlacement" json:"humpPlacement"`
	FrontFlare        string `bson:"frontFlare" json:"frontFlare"`
	SideCurvature     string `bson:"sideCurvature" json:"sideCurvature"`
	HandCompatibility string `bson:"handCompatibility" json:"handCompatibility"`
}

// Battery 电池信息
type Battery struct {
	Type     string `bson:"type" json:"type"`
	Capacity int    `bson:"capacity" json:"capacity"`
	Life     int    `bson:"life" json:"life"`
}

// MouseTechnical 鼠标技术信息
type MouseTechnical struct {
	Connectivity []string `bson:"connectivity" json:"connectivity"`
	Sensor       string   `bson:"sensor" json:"sensor"`
	MaxDPI       int      `bson:"maxDPI" json:"maxDPI"`
	PollingRate  int      `bson:"pollingRate" json:"pollingRate"`
	SideButtons  int      `bson:"sideButtons" json:"sideButtons"`
	Battery      *Battery `bson:"battery,omitempty" json:"battery,omitempty"`
}

// MouseRecommended 鼠标推荐信息
type MouseRecommended struct {
	GameTypes    []string `bson:"gameTypes" json:"gameTypes"`
	GripStyles   []string `bson:"gripStyles" json:"gripStyles"`
	HandSizes    []string `bson:"handSizes" json:"handSizes"`
	DailyUse     bool     `bson:"dailyUse" json:"dailyUse"`
	Professional bool     `bson:"professional" json:"professional"`
}

// MouseDevice 鼠标设备模型
type MouseDevice struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Brand       string             `bson:"brand" json:"brand"`
	Type        string             `bson:"type" json:"type"`
	ImageURL    string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Dimensions  MouseDimensions    `bson:"dimensions" json:"dimensions"`
	Shape       MouseShape         `bson:"shape" json:"shape"`
	Technical   MouseTechnical     `bson:"technical" json:"technical"`
	Recommended MouseRecommended   `bson:"recommended" json:"recommended"`
	SVGData     *MouseSVGData      `bson:"svgData,omitempty" json:"svgData,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// MouseData CSV解析后的鼠标数据
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
	Connectivity      string
	Sensor            string
	DPI               int
	PollingRate       int
	SideButtons       int
}

func main() {
	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB连接信息
	uri := "mongodb://root:example@localhost:27017/?authSource=admin"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("无法连接到MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// 检查连接
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("无法ping MongoDB: %v", err)
	}
	log.Println("成功连接到MongoDB!")

	// 获取数据库和集合
	db := client.Database("cpc")
	collection := db.Collection("devices")

	// 清空现有设备集合
	collection.Drop(ctx)
	
	// 处理CSV数据
	miceFromCSV := processCSVData(ctx)
	
	// 处理SVG数据
	addSVGData(ctx, miceFromCSV, collection)
	
	log.Println("数据初始化完成!")
}

func processCSVData(ctx context.Context) map[string]*MouseDevice {
	miceMap := make(map[string]*MouseDevice)
	
	// 定义CSV文件路径
	csvPaths := []string{
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_.csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (1).csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (2).csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (3).csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (4).csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (5).csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (6).csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (7).csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (8).csv",
		"/mnt/e/Project/Go/FullStackOfYear/data/101_x_25_ (9).csv",
	}
	
	for _, csvPath := range csvPaths {
		log.Printf("处理CSV文件: %s\n", csvPath)
		
		file, err := os.Open(csvPath)
		if err != nil {
			log.Printf("无法打开CSV文件 %s: %v", csvPath, err)
			continue
		}
		defer file.Close()
		
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			log.Printf("无法读取CSV文件 %s: %v", csvPath, err)
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
				log.Printf("解析鼠标数据错误(行 %d): %v", i+1, err)
				continue
			}
			
			// 创建鼠标设备
			mouseDevice := createMouseDevice(mouseData)
			
			// 使用品牌-名称组合作为键
			key := strings.ToLower(mouseDevice.Brand + "-" + mouseDevice.Name)
			miceMap[key] = &mouseDevice
		}
	}
	
	log.Printf("从CSV文件中解析了 %d 个鼠标设备\n", len(miceMap))
	return miceMap
}

func parseMouseData(record []string) (MouseData, error) {
	var mouse MouseData
	var err error
	
	mouse.Name = strings.TrimSpace(record[2])
	if mouse.Name == "" {
		return mouse, fmt.Errorf("鼠标名称为空")
	}
	
	// 解析尺寸数据
	mouse.Length, err = parseFloat(record[3])
	if err != nil {
		return mouse, fmt.Errorf("无效的长度: %v", err)
	}
	
	mouse.Width, err = parseFloat(record[4])
	if err != nil {
		return mouse, fmt.Errorf("无效的宽度: %v", err)
	}
	
	mouse.Height, err = parseFloat(record[5])
	if err != nil {
		return mouse, fmt.Errorf("无效的高度: %v", err)
	}
	
	mouse.Weight, err = parseFloat(record[6])
	if err != nil {
		mouse.Weight = 0
	}
	
	// 形状信息
	mouse.Shape = strings.TrimSpace(record[7])
	mouse.HumpPlacement = strings.TrimSpace(record[8])
	mouse.FrontFlare = strings.TrimSpace(record[9])
	mouse.SideCurvature = strings.TrimSpace(record[10])
	mouse.HandCompatibility = strings.TrimSpace(record[11])
	mouse.Connectivity = strings.TrimSpace(record[15])
	mouse.Sensor = strings.TrimSpace(record[16])
	
	// 技术数据
	mouse.DPI, err = parseInt(record[19])
	if err != nil {
		mouse.DPI = 0
	}
	
	mouse.PollingRate, err = parseInt(record[20])
	if err != nil {
		mouse.PollingRate = 0
	}
	
	if len(record) > 23 {
		mouse.SideButtons, err = parseInt(record[23])
		if err != nil {
			mouse.SideButtons = 0
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
		Description: fmt.Sprintf("%s gaming mouse with %d DPI sensor", brandName, data.DPI),
	}
	
	// 设置尺寸信息
	mouse.Dimensions = MouseDimensions{
		Length: data.Length,
		Width:  data.Width,
		Height: data.Height,
		Weight: data.Weight,
	}
	
	// 设置形状信息
	shapeType := "ambi"
	if strings.Contains(strings.ToLower(data.Shape), "ergonomic") {
		shapeType = "ergo"
	}
	
	mouse.Shape = MouseShape{
		Type:              shapeType,
		HumpPlacement:     strings.ToLower(strings.Split(data.HumpPlacement, " - ")[0]),
		FrontFlare:        strings.ToLower(strings.Split(data.FrontFlare, " - ")[0]),
		SideCurvature:     strings.ToLower(strings.Split(data.SideCurvature, " - ")[0]),
		HandCompatibility: strings.ToLower(data.HandCompatibility),
	}
	
	// 设置技术参数
	connectivitySlice := []string{"wired"}
	if strings.Contains(strings.ToLower(data.Connectivity), "wireless") {
		connectivitySlice = []string{"wireless"}
	}
	
	mouse.Technical = MouseTechnical{
		Connectivity: connectivitySlice,
		Sensor:       data.Sensor,
		MaxDPI:       data.DPI,
		PollingRate:  data.PollingRate,
		SideButtons:  data.SideButtons,
		Battery: &Battery{
			Type:     "lithium-ion",
			Capacity: 500,
			Life:     70,
		},
	}
	
	// 设置推荐信息
	mouse.Recommended = MouseRecommended{
		GameTypes:    inferGameTypes(data),
		GripStyles:   inferGripStyles(data),
		HandSizes:    inferHandSizes(data),
		DailyUse:     true,
		Professional: data.DPI >= 16000,
	}
	
	return mouse
}

func extractBrandName(fullName string) string {
	parts := strings.SplitN(fullName, " ", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return "Unknown"
}

func inferGameTypes(data MouseData) []string {
	gameTypes := []string{}
	
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
		gameTypes = append(gameTypes, "Battle Royale")
	}
	
	return gameTypes
}

func inferGripStyles(data MouseData) []string {
	gripStyles := []string{}
	
	if strings.Contains(strings.ToLower(data.HumpPlacement), "back") {
		gripStyles = append(gripStyles, "palm")
	}
	
	if strings.Contains(strings.ToLower(data.HumpPlacement), "center") {
		gripStyles = append(gripStyles, "claw")
	}
	
	if data.Weight < 80 && data.Length < 120 {
		gripStyles = append(gripStyles, "fingertip")
	}
	
	// 确保至少有一个握持方式
	if len(gripStyles) == 0 {
		gripStyles = append(gripStyles, "palm")
	}
	
	return gripStyles
}

func inferHandSizes(data MouseData) []string {
	handSizes := []string{}
	
	if data.Length <= 115 && data.Width <= 60 {
		handSizes = append(handSizes, "small")
	}
	
	if (data.Length > 115 && data.Length <= 125) || (data.Width > 60 && data.Width <= 68) {
		handSizes = append(handSizes, "medium")
	}
	
	if data.Length > 125 || data.Width > 68 {
		handSizes = append(handSizes, "large")
	}
	
	// 确保至少有一个手型大小
	if len(handSizes) == 0 {
		handSizes = append(handSizes, "medium")
	}
	
	return handSizes
}

func addSVGData(ctx context.Context, miceMap map[string]*MouseDevice, collection *mongo.Collection) {
	// 定义SVG目录
	svgDir := "/mnt/e/Project/Go/FullStackOfYear/SVG"
	
	// 读取目录中的SVG文件
	files, err := ioutil.ReadDir(svgDir)
	if err != nil {
		log.Printf("无法读取SVG目录: %v", err)
		return
	}
	
	// 处理每个SVG文件
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".svg") {
			brand, name, viewType := parseMouseInfo(file.Name())
			if viewType == "unknown" {
				log.Printf("跳过未知视图类型的文件: %s", file.Name())
				continue
			}
			
			// 读取SVG内容
			content, err := ioutil.ReadFile(filepath.Join(svgDir, file.Name()))
			if err != nil {
				log.Printf("无法读取SVG文件 %s: %v", file.Name(), err)
				continue
			}
			
			// 创建唯一的鼠标标识符
			mouseKey := strings.ToLower(brand + "-" + name)
			
			// 检查鼠标是否在我们的地图中
			mouse, exists := miceMap[mouseKey]
			if !exists {
				// 如果没有在CSV中找到，创建一个新的鼠标对象
				now := time.Now()
				mouse = &MouseDevice{
					ID:        primitive.NewObjectID(),
					Name:      name,
					Brand:     brand,
					Type:      "mouse",
					CreatedAt: now,
					UpdatedAt: now,
					Dimensions: MouseDimensions{
						Length: 120,
						Width:  65,
						Height: 40,
						Weight: 70,
					},
					Shape: MouseShape{
						Type:              "ambi",
						HumpPlacement:     "mid",
						FrontFlare:        "minimal",
						SideCurvature:     "moderate",
						HandCompatibility: "universal",
					},
					Technical: MouseTechnical{
						Connectivity: []string{"wireless"},
						Sensor:       "PixArt PAW3950",
						MaxDPI:       26000,
						PollingRate:  8000,
						SideButtons:  2,
						Battery: &Battery{
							Type:     "lithium-ion",
							Capacity: 500,
							Life:     70,
						},
					},
					Recommended: MouseRecommended{
						GameTypes:    []string{"FPS", "MOBA"},
						GripStyles:   []string{"claw", "fingertip"},
						HandSizes:    []string{"medium", "large"},
						DailyUse:     true,
						Professional: true,
					},
					SVGData: &MouseSVGData{},
				}
				miceMap[mouseKey] = mouse
			}
			
			// 确保有SVG数据结构
			if mouse.SVGData == nil {
				mouse.SVGData = &MouseSVGData{}
			}
			
			// 添加SVG数据
			if viewType == "top" {
				mouse.SVGData.TopView = string(content)
			} else if viewType == "side" {
				mouse.SVGData.SideView = string(content)
			}
		}
	}
	
	// 批量插入/更新鼠标数据到数据库
	for _, mouse := range miceMap {
		_, err := collection.UpdateOne(
			ctx,
			map[string]interface{}{
				"name":  mouse.Name,
				"brand": mouse.Brand,
			},
			map[string]interface{}{
				"$set": mouse,
			},
			options.Update().SetUpsert(true),
		)
		
		if err != nil {
			log.Printf("更新鼠标失败 %s %s: %v", mouse.Brand, mouse.Name, err)
		} else {
			log.Printf("成功保存鼠标: %s %s", mouse.Brand, mouse.Name)
		}
	}
}

// 解析鼠标名称和视图类型
func parseMouseInfo(filename string) (string, string, string) {
	// 移除文件扩展名
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	
	// 判断是否包含视图类型
	isTop := strings.Contains(strings.ToLower(filename), "top")
	isSide := strings.Contains(strings.ToLower(filename), "side")
	
	var viewType string
	if isTop {
		viewType = "top"
		filename = strings.Replace(strings.ToLower(filename), "top", "", -1)
	} else if isSide {
		viewType = "side"
		filename = strings.Replace(strings.ToLower(filename), "side", "", -1)
	} else {
		viewType = "unknown"
	}
	
	// 尝试分析品牌和名称
	parts := strings.Split(filename, " ")
	var brand, name string
	
	if len(parts) > 1 {
		// 假设第一个部分是品牌
		brand = strings.TrimSpace(parts[0])
		// 剩余部分是名称
		name = strings.TrimSpace(strings.Join(parts[1:], " "))
	} else {
		brand = "Unknown"
		name = strings.TrimSpace(filename)
	}
	
	return brand, name, viewType
}