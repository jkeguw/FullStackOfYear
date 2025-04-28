db = db.getSiblingDB('cpc');

// 创建用户并分配权限
db.createUser({
    user: "appuser",
    pwd: "apppassword",
    roles: [
        {
            role: "readWrite",
            db: "cpc"
        }
    ]
});

// 创建必要的集合
db.createCollection("users");
db.createCollection("devices");
db.createCollection("reviews");
db.createCollection("orders");
db.createCollection("carts");
db.createCollection("measurements");

// 插入一些示例鼠标设备
// 清除现有设备数据
db.devices.drop();

// 插入一些示例鼠标设备
db.devices.insertMany([
  {
    name: "Logitech G Pro X Superlight 2",
    brand: "Logitech",
    type: "mouse",
    imageUrl: "/images/devices/gpx2.png",
    description: "Ultra-lightweight professional gaming mouse",
    dimensions: {
      length: 125,
      width: 63.5,
      height: 40,
      weight: 60
    },
    shape: {
      type: "ambi",
      humpPlacement: "mid",
      frontFlare: "minimal",
      sideCurvature: "moderate",
      handCompatibility: "universal"
    },
    technical: {
      connectivity: ["wireless"],
      sensor: "HERO 2",
      maxDPI: 32000,
      pollingRate: 8000,
      sideButtons: 2,
      battery: {
        type: "lithium-ion",
        capacity: 900,
        life: 95
      }
    },
    recommended: {
      gameTypes: ["FPS", "MOBA", "Battle Royale"],
      gripStyles: ["claw", "fingertip", "palm"],
      handSizes: ["medium", "large"],
      dailyUse: true,
      professional: true
    },
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Razer Viper V3 Pro",
    brand: "Razer",
    type: "mouse",
    imageUrl: "/images/devices/viperv3.png",
    description: "Ultra-lightweight wireless esports mouse",
    dimensions: {
      length: 128,
      width: 67,
      height: 38,
      weight: 54
    },
    shape: {
      type: "ambi",
      humpPlacement: "mid",
      frontFlare: "minimal",
      sideCurvature: "moderate",
      handCompatibility: "universal"
    },
    technical: {
      connectivity: ["wireless"],
      sensor: "Focus Pro 30K",
      maxDPI: 30000,
      pollingRate: 8000,
      sideButtons: 2,
      battery: {
        type: "lithium-ion",
        capacity: 950,
        life: 90
      }
    },
    recommended: {
      gameTypes: ["FPS", "MOBA"],
      gripStyles: ["claw", "fingertip"],
      handSizes: ["medium", "large"],
      dailyUse: true,
      professional: true
    },
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Pulsar X2 Mini",
    brand: "Pulsar",
    type: "mouse",
    imageUrl: "/images/devices/x2mini.png",
    description: "Lightweight symmetrical mini gaming mouse",
    dimensions: {
      length: 114,
      width: 59,
      height: 39,
      weight: 51
    },
    shape: {
      type: "ambi",
      humpPlacement: "mid",
      frontFlare: "minimal",
      sideCurvature: "moderate",
      handCompatibility: "small-medium"
    },
    technical: {
      connectivity: ["wireless"],
      sensor: "PAW3395",
      maxDPI: 26000,
      pollingRate: 1000,
      sideButtons: 2,
      battery: {
        type: "lithium-ion",
        capacity: 550,
        life: 60
      }
    },
    recommended: {
      gameTypes: ["FPS", "Battle Royale"],
      gripStyles: ["claw", "fingertip"],
      handSizes: ["small", "medium"],
      dailyUse: true,
      professional: true
    },
    createdAt: new Date(),
    updatedAt: new Date()
  }
]);

console.log("Added " + db.devices.find({type: "mouse"}).count() + " mouse devices");

// 创建索引
db.users.createIndex({ "email": 1 }, { unique: true });
db.devices.createIndex({ "name": 1 });
db.devices.createIndex({ "type": 1 });
db.devices.createIndex({ "brand": 1 });
db.reviews.createIndex({ "deviceId": 1 });
db.orders.createIndex({ "userId": 1 });
db.carts.createIndex({ "userId": 1 }, { unique: true });
db.measurements.createIndex({ "deviceId": 1 });

console.log("MongoDB初始化完成");