# MongoDB数据库设计文档

## 1. 数据处理哲学

### 1.1 优先实时获取
对于可从第三方API实时获取的数据（酒店价格与库存、航班价格与状态、实时天气等），一律不进行本地数据库存储，保证数据的实时性和准确性。

### 1.2 本地存储范围
MongoDB仅负责存储以下数据：
- 用户生成的核心数据（UGC）：行程结构、个人笔记、评论、订单信息等
- 无法从外部获取或需要与用户绑定的静态关联信息：收藏列表等
- 系统运行必需的相对静态基础数据：目的地地理坐标、简介等

### 1.3 缓存策略

为高效处理高频第三方API调用，采用Redis作为分布式缓存解决方案，具体实现如下：

### 1.3.1 缓存架构
- **多级缓存**：应用内存缓存(LRU) → Redis分布式缓存 → 第三方API
- **缓存集群**：主从复制+哨兵模式确保高可用
- **数据分片**：按业务域(用户/目的地/行程)进行数据分片

### 1.3.2 缓存键设计规范
采用统一命名格式：`{业务域}:{数据类型}:{唯一标识}:{版本}`

示例：
```
# 目的地天气缓存
weather:destination:10001:v1

# 酒店价格缓存
hotel:price:hotel12345:v2

# 用户行程缓存
itinerary:user:60001:active:v1
```

### 1.3.3 过期策略
根据数据特性设置差异化TTL：
- **高频变动数据**（酒店价格/航班状态）：5-15分钟
- **中等变动数据**（天气/景点开放状态）：1-3小时
- **低频变动数据**（目的地基础信息）：24-72小时
- **用户私有数据**（草稿行程）：7天（滑动窗口）

### 1.3.4 缓存更新机制
- **主动更新**：数据变更时通过发布订阅模式主动更新缓存
- **被动失效**：设置合理TTL+定期后台刷新
- **缓存预热**：系统启动时预加载热门目的地数据

### 1.3.5 缓存一致性保障
- **缓存穿透防护**：布隆过滤器过滤无效Key
- **缓存击穿防护**：热点数据永不过期+定时更新
- **缓存雪崩防护**：TTL随机偏移(±5%)+熔断降级机制
- **数据一致性**：采用"更新数据库→删除缓存"策略，避免脏数据

### 1.3.6 缓存监控
- 实时监控缓存命中率(目标>90%)
- 缓存容量预警(阈值85%)
- 热点Key自动识别与保护
对于高频调用的第三方API，在应用层或使用Redis设计合理的缓存策略，提升性能并降低API调用成本。

## 2. 核心Collection设计

### 2.1 users (用户信息表)
```json
{
  "_id": ObjectId,
  "openid": String,          // 微信用户唯一标识
  "nickname": String,        // 用户昵称
  "avatar_url": String,      // 头像URL
  "gender": Number,          // 性别：0-未知，1-男，2-女
  "phone": String,           // 手机号码（可选）
  "email": String,           // 邮箱（可选）
  "created_at": Date,        // 创建时间
  "updated_at": Date,        // 更新时间
  "last_login_at": Date,     // 最后登录时间
  "status": Number,          // 状态：0-正常，1-禁用
  "preferences": {
    "favorite_destinations": [ObjectId],  // 收藏的目的地
    "travel_styles": [String],            // 旅行风格偏好
    "notification_settings": Object       // 通知设置
  },
  "published_itineraries": [ObjectId],    // 已发布的行程
  "purchased_itineraries": [ObjectId]     // 已购买的行程
}
```

### 2.2 destinations (目的地POI信息表)
```json
{
  "_id": ObjectId,
  "name": String,                // 名称
  "address": String,             // 地址
  "city": String,                // 城市
  "province": String,            // 省份
  "country": String,             // 国家
  "coordinates": {
    "latitude": Number,          // 纬度
    "longitude": Number          // 经度
  },
  "category": String,            // 类别：景点、城市、地区等
  "tags": [String],              // 标签
  "description": String,         // 官方简介
  "images": [String],            // 图片URL列表
  "popularity": Number,          // 热度指数
  "created_at": Date,            // 创建时间
  "updated_at": Date             // 更新时间
}
```

### 2.3 itineraries (行程单信息表)
```json
{
  "_id": ObjectId,
  "author_id": ObjectId,         // 创建者ID
  "title": String,               // 行程标题
  "description": String,         // 行程描述
  "cover_image": String,         // 封面图片URL
  "start_date": Date,            // 开始日期
  "end_date": Date,              // 结束日期
  "destination_ids": [ObjectId], // 目的地IDs
  "status": String,              // 状态：draft-草稿，active-活跃，archived-已归档
  "visibility": String,          // 可见性：private-私有，public-公开，shared-共享
  "collaborators": [{
    "user_id": ObjectId,         // 协作者ID
    "role": String,              // 角色：editor-编辑者，viewer-查看者
    "joined_at": Date            // 加入时间
  }],
  "days": [{
    "day_number": Number,        // 第几天
    "title": String,             // 当日标题
    "description": String,       // 当日描述
    "nodes": [{
      "_id": ObjectId,           // 节点ID
      "type": String,            // 类型：transport-交通，accommodation-住宿，attraction-景点，food-餐饮，shopping-购物，activity-活动
      "title": String,           // 标题
      "description": String,     // 描述
      "start_time": Date,        // 开始时间
      "end_time": Date,          // 结束时间
      "duration": Number,        // 持续时间(分钟)
      "location": {
        "name": String,          // 地点名称
        "address": String,       // 详细地址
        "coordinates": {
          "latitude": Number,    // 纬度
          "longitude": Number    // 经度
        },
        "destination_id": ObjectId // 关联的目的地ID
      },
      "third_party_info": {
        "platform": String,      // 平台名称
        "id": String,            // 第三方平台ID
        "url": String            // 第三方平台URL
      },
      "user_budgeted_price": Number, // 用户预算价格
      "notes": String,           // 用户笔记
      "images": [String],        // 图片URL列表
      "order": Number            // 排序序号
    }]
  }],
  "budget": {
    "currency": String,          // 货币类型
    "total_estimated": Number,   // 总预算估计
    "daily_breakdown": [Number]  // 每日预算
  },
  "tags": [String],              // 标签
  "created_at": Date,            // 创建时间
  "updated_at": Date,            // 更新时间
  "share_count": Number,         // 分享次数
  "view_count": Number           // 查看次数
}
```

### 2.4 itinerary_products (行程商品表)
```json
{
  "_id": ObjectId,
  "author_id": ObjectId,         // 作者ID
  "original_itinerary_id": ObjectId, // 原始行程ID
  "title": String,               // 商品标题
  "description": String,         // 商品描述
  "cover_image": String,         // 封面图片URL
  "price": {
    "amount": Number,            // 价格金额
    "currency": String           // 货币类型
  },
  "status": String,              // 状态：draft-草稿，published-已发布，archived-已归档，sold_out-售罄
  "tags": [String],              // 标签
  "categories": [String],        // 分类
  "snapshot_data": {
    "estimated_cost": Number,    // 发布时估算的总成本
    "days": Number,              // 天数
    "destinations": [String],    // 目的地列表
    "nodes_count": Number        // 节点数量
  },
  "sales_count": Number,         // 销售数量
  "average_rating": Number,      // 平均评分
  "review_count": Number,        // 评论数量
  "created_at": Date,            // 创建时间
  "updated_at": Date,            // 更新时间
  "published_at": Date           // 发布时间
}
```

### 2.5 orders (订单表)
```json
{
  "_id": ObjectId,
  "order_number": String,        // 订单编号
  "buyer_id": ObjectId,          // 购买者ID
  "product_id": ObjectId,        // 行程商品ID
  "original_itinerary_id": ObjectId, // 原始行程ID
  "amount": {
    "total": Number,             // 总金额
    "currency": String           // 货币类型
  },
  "payment_method": String,      // 支付方式
  "payment_status": String,      // 支付状态：pending-待支付，completed-已完成，failed-失败，refunded-已退款
  "payment_time": Date,          // 支付时间
  "transaction_id": String,      // 交易ID
  "status": String,              // 订单状态：created-已创建，paid-已支付，delivered-已交付，completed-已完成，cancelled-已取消
  "created_at": Date,            // 创建时间
  "updated_at": Date             // 更新时间
}
```

### 2.6 reviews (评论表)
```json
{
  "_id": ObjectId,
  "itinerary_product_id": ObjectId, // 行程商品ID
  "author_id": ObjectId,         // 评论作者ID
  "order_id": ObjectId,          // 关联订单ID
  "rating": Number,              // 评分(1-5)
  "content": String,             // 评论内容
  "images": [String],            // 图片URL列表
  "status": String,              // 状态：pending-待审核，approved-已审核，rejected-已拒绝
  "created_at": Date,            // 创建时间
  "updated_at": Date,            // 更新时间
  "reply": {
    "content": String,           // 回复内容
    "created_at": Date           // 回复时间
  }
}
```

## 2.6 数据验证规则

为确保数据一致性和完整性，为各集合添加以下验证规则：

### 2.6.1 users集合验证规则
```json
{
  "$jsonSchema": {
    "bsonType": "object",
    "required": ["openid", "created_at", "status"],
    "properties": {
      "openid": {"bsonType": "string", "description": "微信用户唯一标识，不能为空"},
      "nickname": {"bsonType": "string", "maxLength": 50, "description": "用户昵称，最多50个字符"},
      "avatar_url": {"bsonType": "string", "pattern": "^https?://", "description": "头像URL必须是有效的HTTP/HTTPS URL"},
      "gender": {"enum": [0, 1, 2], "description": "性别只能是0(未知)、1(男)或2(女)"},
      "status": {"enum": [0, 1], "description": "状态只能是0(正常)或1(禁用)"},
      "phone": {"bsonType": "string", "pattern": "^1[3-9]\\d{9}$", "description": "手机号格式必须正确"},
      "email": {"bsonType": "string", "pattern": "^[\\w-]+(\\.[\\w-]+)*@[\\w-]+(\\.[\\w-]+)+$", "description": "邮箱格式必须正确"}
    }
  }
}
```

### 2.6.2 itineraries集合验证规则
```json
{
  "$jsonSchema": {
    "bsonType": "object",
    "required": ["author_id", "title", "start_date", "end_date", "status", "visibility"],
    "properties": {
      "title": {"bsonType": "string", "minLength": 1, "maxLength": 100, "description": "行程标题不能为空且最多100个字符"},
      "status": {"enum": ["draft", "active", "archived"], "description": "状态只能是draft、active或archived"},
      "visibility": {"enum": ["private", "public", "shared"], "description": "可见性只能是private、public或shared"},
      "start_date": {"bsonType": "date", "description": "开始日期必须是有效的日期类型"},
      "end_date": {
        "bsonType": "date",
        "description": "结束日期必须晚于开始日期",
        "$expr": {"$gt": ["$end_date", "$start_date"]}
      },
      "days": {
        "bsonType": "array",
        "minItems": 1,
        "description": "行程至少包含一天",
        "items": {
          "bsonType": "object",
          "properties": {
            "nodes": {
              "bsonType": "array",
              "items": {
                "bsonType": "object",
                "required": ["type", "title", "start_time", "end_time", "order"],
                "properties": {
                  "type": {"enum": ["transport", "accommodation", "attraction", "food", "shopping", "activity"], "description": "节点类型必须是预定义值之一"},
                  "user_budgeted_price": {"bsonType": "number", "minimum": 0, "description": "预算价格不能为负数"}
                }
              }
            }
          }
        }
      }
    }
  }
}
```

## 3. 索引设计

### 3.1 users集合索引
```
{
  "openid": 1
}
{
  "email": 1
}
{
  "phone": 1
}
{
  "created_at": -1
}
```

### 3.2 destinations集合索引
```
{
  "name": "text", "description": "text"
}
{
  "city": 1, "category": 1
}
{
  "coordinates": "2dsphere"
}
{
  "tags": 1
}
```

### 3.3 itineraries集合索引
```
{
  "author_id": 1, "status": 1
}
{
  "destination_ids": 1
}
{
  "start_date": 1, "end_date": 1
}
{
  "status": 1, "visibility": 1
}
{
  "title": "text", "description": "text"
}
```

### 3.4 itinerary_products集合索引
```
{
  "author_id": 1
}
{
  "original_itinerary_id": 1
}
{
  "status": 1, "price.amount": 1
}
{
  "tags": 1, "categories": 1
}
{
  "sales_count": -1
}
{
  "average_rating": -1
}
{
  "title": "text", "description": "text"
}
```

### 3.5 orders集合索引
```
{
  "order_number": 1
}
{
  "buyer_id": 1, "created_at": -1
}
{
  "product_id": 1
}
{
  "payment_status": 1, "status": 1
}
{
  "transaction_id": 1
}
```

### 3.6 reviews集合索引
```
{
  "itinerary_product_id": 1, "status": 1
}
{
  "author_id": 1
}
{
  "rating": -1
}
{
  "created_at": -1
}
```