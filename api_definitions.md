# API与RPC定义文档

## 1. API设计原则

- 遵循RESTful设计规范
- 使用JSON作为数据交换格式
- 统一响应格式
- 版本控制：在URL中包含版本号，如`/api/v1/resource`
- 身份验证：使用JWT令牌
- 错误处理：统一的错误码和错误信息

## 2. 核心API定义

### 2.1 用户服务 (User API)

#### 2.1.1 认证相关
- **登录**
  - URL: `/api/v1/users/login`
  - Method: POST
  - Request: `{"code": string}`
  - Response: `{"token": string, "user": User}`

- **获取当前用户信息**
  - URL: `/api/v1/users/me`
  - Method: GET
  - Headers: `Authorization: Bearer {token}`
  - Response: `{"user": User}`

#### 2.1.2 用户信息管理
- **更新用户资料**
  - URL: `/api/v1/users/profile`
  - Method: PUT
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"nickname": string, "avatar_url": string, ...}`
  - Response: `{"user": User}`

- **获取用户资产**
  - URL: `/api/v1/users/assets`
  - Method: GET
  - Headers: `Authorization: Bearer {token}`
  - Response: `{"published_itineraries": [ItineraryProduct], "purchased_itineraries": [ItineraryProduct]}`

### 2.2 目的地服务 (Destination API)

- **搜索目的地**
  - URL: `/api/v1/destinations/search`
  - Method: GET
  - Query: `{"q": string, "page": number, "limit": number, "tags": string[]}`
  - Response: `{"destinations": [Destination], "total": number, "page": number, "limit": number}`

- **获取目的地详情**
  - URL: `/api/v1/destinations/{id}`
  - Method: GET
  - Response: `{"destination": Destination, "real_time_info": RealTimeInfo}`

- **获取目的地推荐**
  - URL: `/api/v1/destinations/recommendations`
  - Method: GET
  - Query: `{"preferences": string[], "season": string, "budget": number}`
  - Response: `{"destinations": [Destination]}`

### 2.3 行程服务 (Itinerary API)

#### 2.3.1 行程管理
- **创建行程**
  - URL: `/api/v1/itineraries`
  - Method: POST
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"title": string, "description": string, "start_date": string, "end_date": string, "destination_ids": string[]}`
  - Response: `{"itinerary": Itinerary}`

- **获取行程列表**
  - URL: `/api/v1/itineraries`
  - Method: GET
  - Headers: `Authorization: Bearer {token}`
  - Query: `{"status": string, "page": number, "limit": number}`
  - Response: `{"itineraries": [Itinerary], "total": number, "page": number, "limit": number}`

- **获取行程详情**
  - URL: `/api/v1/itineraries/{id}`
  - Method: GET
  - Headers: `Authorization: Bearer {token}`
  - Response: `{"itinerary": Itinerary, "real_time_info": RealTimeInfo[]}`

- **更新行程**
  - URL: `/api/v1/itineraries/{id}`
  - Method: PUT
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"title": string, "description": string, ...}`
  - Response: `{"itinerary": Itinerary}`

- **删除行程**
  - URL: `/api/v1/itineraries/{id}`
  - Method: DELETE
  - Headers: `Authorization: Bearer {token}`
  - Response: `{"success": boolean}`

#### 2.3.2 行程节点管理
- **添加行程节点**
  - URL: `/api/v1/itineraries/{id}/days/{day_number}/nodes`
  - Method: POST
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"type": string, "title": string, "start_time": string, "end_time": string, ...}`
  - Response: `{"node": ItineraryNode}`

- **更新行程节点**
  - URL: `/api/v1/itineraries/{id}/days/{day_number}/nodes/{node_id}`
  - Method: PUT
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"title": string, "description": string, ...}`
  - Response: `{"node": ItineraryNode}`

- **删除行程节点**
  - URL: `/api/v1/itineraries/{id}/days/{day_number}/nodes/{node_id}`
  - Method: DELETE
  - Headers: `Authorization: Bearer {token}`
  - Response: `{"success": boolean}`

### 2.4 社交服务 (Social API)

- **分享行程**
  - URL: `/api/v1/social/share`
  - Method: POST
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"itinerary_id": string, "visibility": string, "message": string}`
  - Response: `{"share_url": string, "success": boolean}`

- **添加协作者**
  - URL: `/api/v1/itineraries/{id}/collaborators`
  - Method: POST
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"user_id": string, "role": string}`
  - Response: `{"collaborator": Collaborator}`

### 2.5 行程市场服务 (Marketplace API)

#### 2.5.1 行程商品管理
- **创建行程商品**
  - URL: `/api/v1/marketplace/products`
  - Method: POST
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"itinerary_id": string, "title": string, "description": string, "price": number, "tags": string[]}`
  - Response: `{"product": ItineraryProduct}`

- **获取行程商品列表**
  - URL: `/api/v1/marketplace/products`
  - Method: GET
  - Query: `{"category": string, "tags": string[], "price_min": number, "price_max": number, "sort": string, "page": number, "limit": number}`
  - Response: `{"products": [ItineraryProduct], "total": number, "page": number, "limit": number}`

- **获取行程商品详情**
  - URL: `/api/v1/marketplace/products/{id}`
  - Method: GET
  - Response: `{"product": ItineraryProduct, "author_info": UserInfo, "reviews": [Review]}`

- **更新行程商品**
  - URL: `/api/v1/marketplace/products/{id}`
  - Method: PUT
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"title": string, "description": string, "price": number, "status": string, ...}`
  - Response: `{"product": ItineraryProduct}`

#### 2.5.2 订单与支付
- **创建订单**
  - URL: `/api/v1/marketplace/orders`
  - Method: POST
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"product_id": string, "payment_method": string}`
  - Response: `{"order": Order, "payment_url": string}`

- **获取订单列表**
  - URL: `/api/v1/marketplace/orders`
  - Method: GET
  - Headers: `Authorization: Bearer {token}`
  - Query: `{"status": string, "page": number, "limit": number}`
  - Response: `{"orders": [Order], "total": number, "page": number, "limit": number}`

#### 2.5.3 评论与评分
- **添加评论**
  - URL: `/api/v1/marketplace/products/{id}/reviews`
  - Method: POST
  - Headers: `Authorization: Bearer {token}`
  - Request: `{"rating": number, "content": string, "images": string[]}`
  - Response: `{"review": Review}`

## 3. RPC服务定义

### 3.1 通用RPC定义格式
所有RPC服务均使用Protocol Buffers定义，遵循以下格式：

```protobuf
syntax = "proto3";

package service_name;

option go_package = "./service_name_rpc/pb"; // Go包路径

// 消息定义
message RequestMessage {
  // 字段定义
}

message ResponseMessage {
  // 字段定义
}

// 服务定义
service ServiceName {
  rpc MethodName(RequestMessage) returns (ResponseMessage);
}
```

### 3.2 Marketplace RPC服务

```protobuf
syntax = "proto3";

package marketplace;

option go_package = "./marketplace_rpc/pb";

import "google/protobuf/empty.proto";
import "user_rpc/pb/user.proto";
import "itinerary_rpc/pb/itinerary.proto";

// 行程商品消息
message ItineraryProduct {
  string id = 1;                // 商品ID
  string author_id = 2;         // 作者ID
  string original_itinerary_id = 3; // 原始行程ID
  string title = 4;             // 标题
  string description = 5;       // 描述
  string cover_image = 6;       // 封面图片
  Price price = 7;              // 价格
  string status = 8;            // 状态
  repeated string tags = 9;     // 标签
  repeated string categories = 10; // 分类
  SnapshotData snapshot_data = 11; // 快照数据
  int32 sales_count = 12;       // 销售数量
  float average_rating = 13;    // 平均评分
  int32 review_count = 14;      // 评论数量
  string created_at = 15;       // 创建时间
  string updated_at = 16;       // 更新时间
  string published_at = 17;     // 发布时间
}

// 价格消息
message Price {
  float amount = 1;             // 金额
  string currency = 2;          // 货币类型
}

// 快照数据消息
message SnapshotData {
  float estimated_cost = 1;     // 估算总成本
  int32 days = 2;               // 天数
  repeated string destinations = 3; // 目的地
  int32 nodes_count = 4;        // 节点数量
}

// 创建行程商品请求
message CreateProductRequest {
  string author_id = 1;         // 作者ID
  string itinerary_id = 2;      // 行程ID
  string title = 3;             // 标题
  string description = 4;       // 描述
  string cover_image = 5;       // 封面图片
  Price price = 6;              // 价格
  repeated string tags = 7;     // 标签
  repeated string categories = 8; // 分类
}

// 创建行程商品响应
message CreateProductResponse {
  ItineraryProduct product = 1; // 行程商品
  bool success = 2;             // 是否成功
  string message = 3;           // 消息
}

// 获取行程商品请求
message GetProductRequest {
  string id = 1;                // 商品ID
}

// 获取行程商品响应
message GetProductResponse {
  ItineraryProduct product = 1; // 行程商品
  user.UserInfo author_info = 2; // 作者信息
}

// 列出行程商品请求
message ListProductsRequest {
  string category = 1;          // 分类
  repeated string tags = 2;     // 标签
  float price_min = 3;          // 最低价格
  float price_max = 4;          // 最高价格
  string sort = 5;              // 排序方式
  int32 page = 6;               // 页码
  int32 limit = 7;              // 每页数量
}

// 列出行程商品响应
message ListProductsResponse {
  repeated ItineraryProduct products = 1; // 行程商品列表
  int32 total = 2;               // 总数
  int32 page = 3;                // 页码
  int32 limit = 4;               // 每页数量
}

// 更新行程商品请求
message UpdateProductRequest {
  string id = 1;                // 商品ID
  string author_id = 2;         // 作者ID
  string title = 3;             // 标题
  string description = 4;       // 描述
  string cover_image = 5;       // 封面图片
  Price price = 6;              // 价格
  string status = 7;            // 状态
  repeated string tags = 8;     // 标签
  repeated string categories = 9; // 分类
}

// 更新行程商品响应
message UpdateProductResponse {
  ItineraryProduct product = 1; // 行程商品
  bool success = 2;             // 是否成功
  string message = 3;           // 消息
}

// 行程市场RPC服务
service Marketplace {
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
  rpc DeleteProduct(GetProductRequest) returns (google.protobuf.Empty);
}
```

### 3.3 其他RPC服务
其他服务（用户、行程、社交等）的RPC定义将遵循类似模式，包含各自领域的消息和方法定义。