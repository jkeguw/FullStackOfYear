package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// 定义鼠标数据结构（与models中的结构一致）
type MouseSVGData struct {
	TopView  string `bson:"topView,omitempty" json:"topView,omitempty"`
	SideView string `bson:"sideView,omitempty" json:"sideView,omitempty"`
}

type MouseDimensions struct {
	Length float64 `bson:"length" json:"length"`
	Width  float64 `bson:"width" json:"width"`
	Height float64 `bson:"height" json:"height"`
	Weight float64 `bson:"weight" json:"weight"`
}

type MouseShape struct {
	Type              string `bson:"type" json:"type"`
	HumpPlacement     string `bson:"humpPlacement" json:"humpPlacement"`
	FrontFlare        string `bson:"frontFlare" json:"frontFlare"`
	SideCurvature     string `bson:"sideCurvature" json:"sideCurvature"`
	HandCompatibility string `bson:"handCompatibility" json:"handCompatibility"`
}

type MouseTechnical struct {
	Connectivity []string `bson:"connectivity" json:"connectivity"`
	Sensor       string   `bson:"sensor" json:"sensor"`
	MaxDPI       int      `bson:"maxDPI" json:"maxDPI"`
	PollingRate  int      `bson:"pollingRate" json:"pollingRate"`
	SideButtons  int      `bson:"sideButtons" json:"sideButtons"`
	Weight       float64  `bson:"weight,omitempty" json:"weight,omitempty"`
}

type MouseRecommended struct {
	GameTypes    []string `bson:"gameTypes" json:"gameTypes"`
	GripStyles   []string `bson:"gripStyles" json:"gripStyles"`
	HandSizes    []string `bson:"handSizes" json:"handSizes"`
	DailyUse     bool     `bson:"dailyUse" json:"dailyUse"`
	Professional bool     `bson:"professional" json:"professional"`
}

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

// 解析鼠标名称
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

func main() {
	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取MongoDB连接URI，默认值供本地测试
	uri := "mongodb://root:example@localhost:27017"
	if os.Getenv("MONGODB_URI") != "" {
		uri = os.Getenv("MONGODB_URI")
	}

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// 检查连接
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	
	log.Println("成功连接到MongoDB")

	// 获取数据库和集合
	dbName := "cpc"
	if os.Getenv("MONGODB_DATABASE") != "" {
		dbName = os.Getenv("MONGODB_DATABASE")
	}

	db := client.Database(dbName)
	collection := db.Collection("devices")

	// 定义SVG目录
	svgDir := "/mnt/e/Project/Go/FullStackOfYear/SVG"
	if os.Getenv("SVG_DIR") != "" {
		svgDir = os.Getenv("SVG_DIR")
	}

	// 读取目录中的SVG文件
	files, err := ioutil.ReadDir(svgDir)
	if err != nil {
		log.Fatalf("Failed to read SVG directory: %v", err)
	}

	// 按鼠标名称组织文件
	mouseMap := make(map[string]map[string]string)

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
				log.Printf("读取SVG文件失败 %s: %v", file.Name(), err)
				continue
			}

			// 创建唯一的鼠标标识符
			mouseKey := strings.ToLower(brand + "-" + name)

			// 确保mouseMap中有这个鼠标的条目
			if _, exists := mouseMap[mouseKey]; !exists {
				mouseMap[mouseKey] = make(map[string]string)
				mouseMap[mouseKey]["brand"] = brand
				mouseMap[mouseKey]["name"] = name
			}

			// 添加SVG数据
			mouseMap[mouseKey][viewType+"View"] = string(content)
		}
	}

	// 插入或更新数据库
	for _, mouseData := range mouseMap {
		brand := mouseData["brand"]
		name := mouseData["name"]

		// 设置过滤条件
		filter := map[string]interface{}{
			"name":  name,
			"brand": brand,
			"type":  "mouse",
		}

		now := time.Now()

		// 添加SVG数据
		svgData := &MouseSVGData{}
		if topView, ok := mouseData["topView"]; ok {
			svgData.TopView = topView
		}

		if sideView, ok := mouseData["sideView"]; ok {
			svgData.SideView = sideView
		}

		// 使用upsert操作插入或更新
		opts := options.Update().SetUpsert(true)
		update := map[string]interface{}{
			"$set": map[string]interface{}{
				"svgData":   svgData,
				"updatedAt": now,
			},
		}

		result, err := collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("更新鼠标失败 %s %s: %v", brand, name, err)
			continue
		}

		if result.UpsertedCount > 0 {
			log.Printf("插入新鼠标: %s %s", brand, name)
		} else if result.ModifiedCount > 0 {
			log.Printf("成功保存鼠标: %s %s", brand, name)
		} else {
			log.Printf("鼠标无变化: %s %s", brand, name)
		}
	}

	fmt.Println("SVG数据导入完成!")
}