# API与RPC定义文档

## 1. API设计原则

- 遵循RESTful设计规范
- 使用JSON作为数据交换格式
- 统一响应格式
- 版本控制：在URL中包含版本号，如`/api/v1/resource`
- 身份验证：使用JWT令牌
- 错误处理：统一的错误码和错误信息

### 1.1 API版本控制策略

采用URI路径版本控制，格式为：`/api/v{version}/{resource}`

#### 1.1.1 版本号规则
- 主版本号：`v1`, `v2`...，表示不兼容的API变更
- 次版本更新：通过API文档和变更日志传达，不修改版本号

#### 1.1.2 版本管理原则
- 每个主版本API独立部署，可并行运行
- 新功能优先添加到最新版本
- 旧版本API提供至少6个月的过渡期支持
- 版本弃用前3个月发布公告

#### 1.1.3 版本迁移策略
- 提供版本迁移指南和自动化工具
- 保留关键业务数据的向下兼容性
- 提供版本切换的灰度发布能力

#### 1.1.4 示例
```
# 当前版本
GET /api/v1/itineraries/60001

# 旧版本（过渡期内可用）
GET /api/v0/itineraries/60001
```

### 1.2 错误码规范

采用分层错误码设计，格式为：`XXX-YYYY`，其中：
- `XXX`：服务域标识（100-999）
- `YYYY`：具体错误编号（0001-9999）

#### 1.1.1 服务域划分
- **100**：用户服务（User Service）
- **200**：目的地服务（Destination Service）
- **300**：行程服务（Itinerary Service）
- **400**：社交服务（Social Service）
- **500**：市场服务（Marketplace Service）
- **900**：系统级错误

#### 1.1.2 通用错误码
| 错误码 | 描述 | HTTP状态码 | 处理建议 |
|--------|------|------------|----------|
| 900-0001 | 系统内部错误 | 500 | 查看服务日志 |
| 900-0002 | 参数验证失败 | 400 | 检查请求参数格式 |
| 900-0003 | 资源不存在 | 404 | 确认资源ID有效性 |
| 900-0004 | 权限不足 | 403 | 检查用户权限 |
| 900-0005 | 请求频率限制 | 429 | 减少请求频率 |
| 900-0006 | 第三方服务错误 | 503 | 稍后重试 |

#### 1.1.3 错误响应格式
所有API错误响应遵循统一格式：
```json
{
  "code": "300-0002",
  "message": "行程不存在",
  "details": {
    "itinerary_id": "60001",
    "timestamp": "2023-11-15T14:30:00Z"
  },
  "request_id": "req-123e4567-e89b-12d3-a456-426614174000"
}
```

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
// 通用错误消息
message Error {
  int32 code = 1;        // 错误码
  string message = 2;    // 错误信息
  string details = 3;    // 详细错误描述
}

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
### 3.4 用户RPC服务

```protobuf
syntax = "proto3";

package user;

option go_package = "./user_rpc/pb";

import "google/protobuf/empty.proto";

// 用户信息消息
message User {
  string id = 1;                // 用户ID
  string openid = 2;            // 微信openid
  string nickname = 3;          // 昵称
  string avatar_url = 4;        // 头像URL
  int32 gender = 5;             // 性别
  string phone = 6;             // 手机号
  string email = 7;             // 邮箱
  string created_at = 8;        // 创建时间
  string updated_at = 9;        // 更新时间
  string last_login_at = 10;    // 最后登录时间
  int32 status = 11;            // 状态
}

// 用户信息请求
message GetUserRequest {
  string id = 1;                // 用户ID
}

// 用户资产响应
message UserAssetsResponse {
  repeated string published_itineraries = 1; // 已发布行程
  repeated string purchased_itineraries = 2; // 已购买行程
  double balance = 3;             // 账户余额
  int32 points = 4;               // 积分
}

// 用户信息响应
message GetUserResponse {
  User user = 1;                // 用户信息
}

// 更新用户信息请求
message UpdateUserRequest {
  string id = 1;                // 用户ID
  string nickname = 2;          // 昵称
  string avatar_url = 3;        // 头像URL
  int32 gender = 4;             // 性别
  string phone = 5;             // 手机号
  string email = 6;             // 邮箱
}

// 更新用户信息响应
message UpdateUserResponse {
  User user = 1;                // 更新后的用户信息
  bool success = 2;             // 是否成功
  string message = 3;           // 消息
}

// 用户RPC服务
service User {
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {
  // 错误码:
  // 400: 请求参数无效
  // 404: 用户不存在
  // 500: 服务器内部错误
}

rpc GetUserAssets(GetUserAssetsRequest) returns (UserAssetsResponse) {
  // 错误码:
  // 400: 请求参数无效
  // 401: 未授权访问
  // 404: 用户不存在
}
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc GetUserAssets(GetUserRequest) returns (UserAssetsResponse);
}
```

### 3.5 行程RPC服务

```protobuf
syntax = "proto3";

package itinerary;

option go_package = "./itinerary_rpc/pb";

import "google/protobuf/empty.proto";
import "user_rpc/pb/user.proto";

// 行程节点消息
message ItineraryNode {
  string id = 1;                // 节点ID
  string type = 2;              // 类型
  string title = 3;             // 标题
  string description = 4;       // 描述
  string start_time = 5;        // 开始时间
  string end_time = 6;          // 结束时间
  int32 duration = 7;           // 持续时间(分钟)
  Location location = 8;        // 位置信息
  ThirdPartyInfo third_party_info = 9; // 第三方信息
  double user_budgeted_price = 10; // 用户预算价格
  string notes = 11;            // 用户笔记
  repeated string images = 12;  // 图片URL列表
  int32 order = 13;             // 排序序号
}

// 行程消息
message Itinerary {
  string id = 1;                // 行程ID
  string author_id = 2;         // 创建者ID
  string title = 3;             // 行程标题
  string description = 4;       // 行程描述
  string cover_image = 5;       // 封面图片URL
  string start_date = 6;        // 开始日期
  string end_date = 7;          // 结束日期
  repeated string destination_ids = 8; // 目的地IDs
  string status = 9;            // 状态
  string visibility = 10;       // 可见性
  repeated Collaborator collaborators = 11; // 协作者
  repeated ItineraryDay days = 12; // 每日行程
  Budget budget = 13;           // 预算信息
  repeated string tags = 14;    // 标签
  string created_at = 15;       // 创建时间
  string updated_at = 16;       // 更新时间
}

// 行程RPC服务
service Itinerary {
  rpc CreateItinerary(CreateItineraryRequest) returns (CreateItineraryResponse);
  rpc GetItinerary(GetItineraryRequest) returns (GetItineraryResponse);
  rpc UpdateItinerary(UpdateItineraryRequest) returns (UpdateItineraryResponse);
  rpc DeleteItinerary(DeleteItineraryRequest) returns (google.protobuf.Empty);
  rpc AddItineraryNode(AddItineraryNodeRequest) returns (AddItineraryNodeResponse);
}
```

### 3.6 社交RPC服务

```protobuf
syntax = "proto3";

package social;

option go_package = "./social_rpc/pb";

import "google/protobuf/empty.proto";
import "user_rpc/pb/user.proto";
import "itinerary_rpc/pb/itinerary.proto";

// 分享消息
message Share {
  string id = 1;                // 分享ID
  string itinerary_id = 2;      // 行程ID
  string user_id = 3;           // 用户ID
  string title = 4;             // 标题
  string description = 5;       // 描述
  string share_url = 6;         // 分享URL
  string visibility = 7;        // 可见性
  string created_at = 8;        // 创建时间
  int32 view_count = 9;         // 查看次数
}

// 评论消息
message Comment {
  string id = 1;                // 评论ID
  string itinerary_id = 2;      // 行程ID
  string user_id = 3;           // 用户ID
  string content = 4;           // 内容
  repeated string images = 5;   // 图片URL列表
  string created_at = 6;        // 创建时间
  string updated_at = 7;        // 更新时间
  int32 status = 8;             // 状态
}

// 社交RPC服务
service Social {
  rpc ShareItinerary(ShareItineraryRequest) returns (ShareItineraryResponse);
  rpc AddCollaborator(AddCollaboratorRequest) returns (AddCollaboratorResponse);
  rpc AddComment(AddCommentRequest) returns (AddCommentResponse);
  rpc GetComments(GetCommentsRequest) returns (GetCommentsResponse);
}
```

其他服务的RPC定义遵循上述类似模式，包含各自领域的消息和方法定义。