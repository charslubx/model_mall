# 智能服饰商城系统 - 后端接口需求文档

## 文档说明

**项目名称**: 智能服饰商城系统  
**文档版本**: v1.0  
**更新日期**: 2025-01-05  
**接口协议**: RESTful API  
**数据格式**: JSON  
**字符编码**: UTF-8

---

## 目录

- [1. 认证模块](#1-认证模块)
- [2. 商品模块](#2-商品模块)
- [3. 购物车模块](#3-购物车模块)
- [4. 订单模块](#4-订单模块)
- [5. 用户个人中心模块](#5-用户个人中心模块)
- [6. 商户管理模块](#6-商户管理模块)
- [7. 卖家店铺模块](#7-卖家店铺模块)
- [8. 系统管理模块](#8-系统管理模块)
- [9. 文件上传模块](#9-文件上传模块)
- [10. 搜索与推荐模块](#10-搜索与推荐模块)
- [通用规范](#通用规范)

---

## 1. 认证模块

### 1.1 用户登录

**接口地址**: `POST /api/auth/login`  
**接口描述**: 客户和商户登录接口  
**是否需要认证**: 否

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 是 | 用户邮箱 |
| password | string | 是 | 用户密码（明文传输需使用HTTPS） |
| userType | string | 是 | 用户类型：customer(客户) / merchant(商户) |

**请求示例**:
```json
{
  "email": "zhangsan@example.com",
  "password": "Password123!",
  "userType": "customer"
}
```

**响应参数**:

| 参数名 | 类型 | 说明 |
|--------|------|------|
| token | string | JWT访问令牌，有效期2小时 |
| refreshToken | string | 刷新令牌，有效期7天 |
| user | object | 用户信息对象 |
| user.id | string | 用户ID |
| user.name | string | 用户姓名 |
| user.email | string | 用户邮箱 |
| user.avatar | string | 用户头像URL |
| user.userType | string | 用户类型 |
| user.createdAt | string | 注册时间（ISO 8601格式） |

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "u123456",
      "name": "张三",
      "email": "zhangsan@example.com",
      "avatar": "https://cdn.example.com/avatars/u123456.jpg",
      "userType": "customer",
      "createdAt": "2023-01-15T08:30:00Z"
    }
  }
}
```

---

### 1.2 用户注册

**接口地址**: `POST /api/auth/register`  
**接口描述**: 新用户注册接口  
**是否需要认证**: 否

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 用户姓名 |
| email | string | 是 | 用户邮箱 |
| password | string | 是 | 密码（至少8位，包含字母和数字） |
| userType | string | 是 | 用户类型：customer / merchant |
| phone | string | 否 | 手机号码 |
| merchantName | string | 条件必填 | 商户名称（userType为merchant时必填） |
| businessLicense | string | 否 | 营业执照号码（商户注册时选填） |

**请求示例**:
```json
{
  "name": "张三",
  "email": "zhangsan@example.com",
  "password": "Password123!",
  "userType": "customer",
  "phone": "13800138000"
}
```

**响应参数**: 同登录接口

**响应示例**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "u123456",
      "name": "张三",
      "email": "zhangsan@example.com",
      "avatar": null,
      "userType": "customer",
      "createdAt": "2025-01-05T10:30:00Z"
    }
  }
}
```

---

### 1.3 用户登出

**接口地址**: `POST /api/auth/logout`  
**接口描述**: 用户登出，使当前token失效  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**: 无

**响应示例**:
```json
{
  "code": 200,
  "message": "登出成功",
  "data": {
    "success": true
  }
}
```

---

### 1.4 获取当前用户信息

**接口地址**: `GET /api/auth/user`  
**接口描述**: 获取当前登录用户的详细信息  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**响应参数**:

| 参数名 | 类型 | 说明 |
|--------|------|------|
| id | string | 用户ID |
| name | string | 用户姓名 |
| email | string | 用户邮箱 |
| phone | string | 手机号码 |
| avatar | string | 头像URL |
| userType | string | 用户类型 |
| status | string | 账户状态：active(正常) / disabled(禁用) |
| createdAt | string | 注册时间 |
| lastLoginAt | string | 最后登录时间 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "u123456",
    "name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "avatar": "https://cdn.example.com/avatars/u123456.jpg",
    "userType": "customer",
    "status": "active",
    "createdAt": "2023-01-15T08:30:00Z",
    "lastLoginAt": "2025-01-05T10:30:00Z"
  }
}
```

---

### 1.5 刷新访问令牌

**接口地址**: `POST /api/auth/refresh`  
**接口描述**: 使用刷新令牌获取新的访问令牌  
**是否需要认证**: 否

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| refreshToken | string | 是 | 刷新令牌 |

**请求示例**:
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### 1.6 管理员登录

**接口地址**: `POST /api/auth/admin/login`  
**接口描述**: 系统管理员登录接口  
**是否需要认证**: 否

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 是 | 管理员账号 |
| password | string | 是 | 管理员密码 |

**请求示例**:
```json
{
  "email": "admin@example.com",
  "password": "AdminPassword123!"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin": {
      "id": "admin001",
      "name": "超级管理员",
      "email": "admin@example.com",
      "role": "admin",
      "permissions": ["user_manage", "order_manage", "product_manage", "system_config"]
    }
  }
}
```

---

## 2. 商品模块

### 2.1 获取商品列表

**接口地址**: `GET /api/products`  
**接口描述**: 获取商品列表，支持分页、筛选和排序  
**是否需要认证**: 否

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认20，最大100 |
| category | string | 否 | 商品分类 |
| keyword | string | 否 | 搜索关键词 |
| sortBy | string | 否 | 排序方式：price-asc(价格升序) / price-desc(价格降序) / name-asc / name-desc / sales(销量) / newest(最新) |
| minPrice | number | 否 | 最低价格 |
| maxPrice | number | 否 | 最高价格 |

**请求示例**:
```
GET /api/products?page=1&pageSize=20&category=上衣&sortBy=price-asc
```

**响应参数**:

| 参数名 | 类型 | 说明 |
|--------|------|------|
| products | array | 商品列表数组 |
| products[].id | string | 商品ID |
| products[].name | string | 商品名称 |
| products[].category | string | 商品分类 |
| products[].price | number | 商品价格 |
| products[].image | string | 商品主图URL |
| products[].rating | number | 商品评分(0-5) |
| products[].sales | number | 销量 |
| products[].stock | number | 库存数量 |
| products[].tags | array | 商品标签数组 |
| products[].seller | object | 卖家信息 |
| products[].seller.id | string | 卖家ID |
| products[].seller.name | string | 卖家名称 |
| products[].seller.avatar | string | 卖家头像 |
| total | number | 总记录数 |
| page | number | 当前页码 |
| pageSize | number | 每页数量 |
| totalPages | number | 总页数 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "products": [
      {
        "id": "p001",
        "name": "简约纯棉T恤",
        "category": "上衣",
        "price": 99.00,
        "image": "https://cdn.example.com/products/p001.jpg",
        "rating": 4.5,
        "sales": 1250,
        "stock": 120,
        "tags": ["舒适", "百搭", "新品"],
        "seller": {
          "id": "m001",
          "name": "时尚优选",
          "avatar": "https://cdn.example.com/merchants/m001.jpg"
        }
      }
    ],
    "total": 156,
    "page": 1,
    "pageSize": 20,
    "totalPages": 8
  }
}
```

---

### 2.2 获取商品详情

**接口地址**: `GET /api/products/:id`  
**接口描述**: 获取指定商品的详细信息  
**是否需要认证**: 否

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 商品ID |

**请求示例**:
```
GET /api/products/p001
```

**响应参数**:

| 参数名 | 类型 | 说明 |
|--------|------|------|
| id | string | 商品ID |
| name | string | 商品名称 |
| description | string | 商品描述 |
| category | string | 商品分类 |
| price | number | 商品价格 |
| stock | number | 库存数量 |
| rating | number | 商品评分 |
| reviews | number | 评价数量 |
| images | array | 商品图片URL数组 |
| colors | array | 可选颜色数组 |
| colors[].name | string | 颜色名称 |
| colors[].value | string | 颜色值 |
| colors[].hex | string | 颜色十六进制代码 |
| sizes | array | 可选尺码数组 |
| tags | array | 商品标签数组 |
| features | array | 商品特点数组 |
| specifications | object | 规格参数对象 |
| specifications.material | string | 材质 |
| specifications.care | string | 洗涤说明 |
| seller | object | 卖家信息 |
| seller.id | string | 卖家ID |
| seller.name | string | 卖家名称 |
| seller.avatar | string | 卖家头像 |
| seller.rating | number | 卖家评分 |
| relatedProducts | array | 相关推荐商品数组 |
| createdAt | string | 创建时间 |
| updatedAt | string | 更新时间 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "p001",
    "name": "简约纯棉T恤",
    "description": "采用优质纯棉面料，柔软舒适，透气性好，简约百搭设计，适合日常休闲穿着。",
    "category": "上衣",
    "price": 99.00,
    "stock": 120,
    "rating": 4.5,
    "reviews": 128,
    "images": [
      "https://cdn.example.com/products/p001-1.jpg",
      "https://cdn.example.com/products/p001-2.jpg",
      "https://cdn.example.com/products/p001-3.jpg"
    ],
    "colors": [
      { "name": "白色", "value": "white", "hex": "#FFFFFF" },
      { "name": "黑色", "value": "black", "hex": "#000000" },
      { "name": "灰色", "value": "gray", "hex": "#808080" }
    ],
    "sizes": ["S", "M", "L", "XL", "XXL"],
    "tags": ["舒适", "百搭", "新品", "热卖"],
    "features": ["100%纯棉面料", "圆领设计", "简约百搭", "四季可穿"],
    "specifications": {
      "material": "100%纯棉",
      "care": "30°C机洗，不可漂白，不可烘干"
    },
    "seller": {
      "id": "m001",
      "name": "时尚优选",
      "avatar": "https://cdn.example.com/merchants/m001.jpg",
      "rating": 4.8
    },
    "relatedProducts": [
      {
        "id": "p002",
        "name": "轻薄防晒衬衫",
        "price": 129.00,
        "image": "https://cdn.example.com/products/p002.jpg"
      }
    ],
    "createdAt": "2024-12-01T10:00:00Z",
    "updatedAt": "2025-01-03T15:30:00Z"
  }
}
```

---

### 2.3 获取商品分类列表

**接口地址**: `GET /api/categories`  
**接口描述**: 获取所有商品分类  
**是否需要认证**: 否

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "categories": [
      { "id": "c001", "name": "上衣", "slug": "tops", "count": 1250, "icon": null },
      { "id": "c002", "name": "裤装", "slug": "pants", "count": 856, "icon": null },
      { "id": "c003", "name": "裙装", "slug": "skirts", "count": 432, "icon": null },
      { "id": "c004", "name": "外套", "slug": "outerwear", "count": 678, "icon": null },
      { "id": "c005", "name": "鞋履", "slug": "shoes", "count": 543, "icon": null },
      { "id": "c006", "name": "配饰", "slug": "accessories", "count": 234, "icon": null }
    ]
  }
}
```

---

### 2.4 创建商品（商户）

**接口地址**: `POST /api/products`  
**接口描述**: 商户创建新商品  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 商品名称 |
| category | string | 是 | 商品分类 |
| price | number | 是 | 商品价格（>0） |
| stock | number | 是 | 库存数量（>=0） |
| description | string | 否 | 商品描述 |
| images | array | 是 | 商品图片URL数组（至少1张） |
| tags | array | 否 | 商品标签数组 |
| colors | array | 否 | 可选颜色数组 |
| sizes | array | 否 | 可选尺码数组 |
| features | array | 否 | 商品特点数组 |
| specifications | object | 否 | 规格参数对象 |

**请求示例**:
```json
{
  "name": "简约纯棉T恤",
  "category": "上衣",
  "price": 99.00,
  "stock": 120,
  "description": "采用优质纯棉面料，柔软舒适",
  "images": [
    "https://cdn.example.com/products/p001-1.jpg",
    "https://cdn.example.com/products/p001-2.jpg"
  ],
  "tags": ["舒适", "百搭", "新品"],
  "colors": [
    { "name": "白色", "value": "white", "hex": "#FFFFFF" },
    { "name": "黑色", "value": "black", "hex": "#000000" }
  ],
  "sizes": ["S", "M", "L", "XL"],
  "features": ["100%纯棉面料", "圆领设计"],
  "specifications": {
    "material": "100%纯棉",
    "care": "30°C机洗"
  }
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": "p001",
    "name": "简约纯棉T恤",
    "status": "active",
    "createdAt": "2025-01-05T10:30:00Z"
  }
}
```

---

### 2.5 更新商品（商户）

**接口地址**: `PUT /api/products/:id`  
**接口描述**: 商户更新商品信息  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 商品ID |

**请求参数**: 同创建商品，所有字段均为可选

**请求示例**:
```json
{
  "price": 89.00,
  "stock": 150,
  "tags": ["舒适", "百搭", "新品", "促销"]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": "p001",
    "updated": true,
    "updatedAt": "2025-01-05T11:00:00Z"
  }
}
```

---

### 2.6 删除商品（商户）

**接口地址**: `DELETE /api/products/:id`  
**接口描述**: 商户删除商品  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 商品ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": {
    "success": true
  }
}
```

---

### 2.7 获取商户商品列表

**接口地址**: `GET /api/merchant/products`  
**接口描述**: 商户查询自己的商品列表  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认20 |
| category | string | 否 | 商品分类 |
| status | string | 否 | 商品状态：active(在售) / inactive(下架) |
| keyword | string | 否 | 搜索关键词 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "products": [
      {
        "id": "p001",
        "name": "简约纯棉T恤",
        "category": "上衣",
        "price": 99.00,
        "stock": 120,
        "status": "active",
        "image": "https://cdn.example.com/products/p001.jpg",
        "tags": ["舒适", "百搭", "新品"],
        "sales": 1250,
        "createdAt": "2024-12-01T10:00:00Z"
      }
    ],
    "total": 45,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 2.8 AI生成商品标签

**接口地址**: `POST /api/products/generate-tags`  
**接口描述**: 使用AI分析商品图片自动生成标签  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| imageUrl | string | 是 | 商品图片URL |
| category | string | 是 | 商品分类 |

**请求示例**:
```json
{
  "imageUrl": "https://cdn.example.com/products/p001.jpg",
  "category": "上衣"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "生成成功",
  "data": {
    "tags": ["舒适", "时尚", "百搭", "休闲", "简约", "新品"]
  }
}
```

---

## 3. 购物车模块

### 3.1 获取购物车列表

**接口地址**: `GET /api/cart`  
**接口描述**: 获取当前用户的购物车商品列表  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "items": [
      {
        "id": "cart001",
        "productId": "p001",
        "name": "简约纯棉T恤",
        "price": 99.00,
        "quantity": 2,
        "color": "白色",
        "size": "M",
        "image": "https://cdn.example.com/products/p001.jpg",
        "stock": 120,
        "selected": true
      }
    ],
    "subtotal": 198.00,
    "total": 198.00
  }
}
```

---

### 3.2 添加商品到购物车

**接口地址**: `POST /api/cart`  
**接口描述**: 添加商品到购物车  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| productId | string | 是 | 商品ID |
| color | string | 否 | 选择的颜色 |
| size | string | 否 | 选择的尺码 |
| quantity | number | 是 | 数量（>0） |

**请求示例**:
```json
{
  "productId": "p001",
  "color": "白色",
  "size": "M",
  "quantity": 2
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "添加成功",
  "data": {
    "id": "cart001",
    "productId": "p001",
    "quantity": 2
  }
}
```

---

### 3.3 更新购物车商品数量

**接口地址**: `PUT /api/cart/:itemId`  
**接口描述**: 更新购物车中商品的数量  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| itemId | string | 是 | 购物车项ID |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| quantity | number | 是 | 新的数量（>0） |

**请求示例**:
```json
{
  "quantity": 3
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": "cart001",
    "quantity": 3
  }
}
```

---

### 3.4 删除购物车商品

**接口地址**: `DELETE /api/cart/:itemId`  
**接口描述**: 从购物车中删除指定商品  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| itemId | string | 是 | 购物车项ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": {
    "success": true
  }
}
```

---

### 3.5 清空购物车

**接口地址**: `DELETE /api/cart`  
**接口描述**: 清空购物车中的所有商品  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "清空成功",
  "data": {
    "success": true
  }
}
```

---

## 4. 订单模块

### 4.1 创建订单

**接口地址**: `POST /api/orders`  
**接口描述**: 创建新订单  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| items | array | 是 | 订单商品数组 |
| items[].productId | string | 是 | 商品ID |
| items[].quantity | number | 是 | 数量 |
| items[].color | string | 否 | 颜色 |
| items[].size | string | 否 | 尺码 |
| address | object | 是 | 收货地址 |
| address.name | string | 是 | 收货人姓名 |
| address.phone | string | 是 | 收货人电话 |
| address.province | string | 是 | 省份 |
| address.city | string | 是 | 城市 |
| address.district | string | 是 | 区县 |
| address.detail | string | 是 | 详细地址 |
| paymentMethod | string | 是 | 支付方式：alipay / wechat / union |
| note | string | 否 | 订单备注 |

**请求示例**:
```json
{
  "items": [
    {
      "productId": "p001",
      "quantity": 2,
      "color": "白色",
      "size": "M"
    }
  ],
  "address": {
    "name": "张三",
    "phone": "13800138000",
    "province": "北京市",
    "city": "北京市",
    "district": "朝阳区",
    "detail": "某某街道某某小区1号楼1单元101"
  },
  "paymentMethod": "alipay",
  "note": "请尽快发货"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "orderId": "ord20250105001",
    "orderNo": "ORD20250105123456",
    "total": 208.00,
    "paymentUrl": "https://payment.example.com/pay?orderNo=ORD20250105123456"
  }
}
```

---

### 4.2 获取订单列表

**接口地址**: `GET /api/orders`  
**接口描述**: 获取用户的订单列表  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认20 |
| status | string | 否 | 订单状态：pending(待付款) / paid(已付款) / shipped(已发货) / completed(已完成) / cancelled(已取消) |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "orders": [
      {
        "id": "ord20250105001",
        "orderNo": "ORD20250105123456",
        "date": "2025-01-05T10:30:00Z",
        "status": "paid",
        "statusText": "已付款",
        "total": 208.00,
        "itemCount": 2,
        "items": [
          {
            "productId": "p001",
            "name": "简约纯棉T恤",
            "image": "https://cdn.example.com/products/p001.jpg",
            "quantity": 2
          }
        ]
      }
    ],
    "total": 15,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 4.3 获取订单详情

**接口地址**: `GET /api/orders/:id`  
**接口描述**: 获取指定订单的详细信息  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 订单ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "ord20250105001",
    "orderNo": "ORD20250105123456",
    "date": "2025-01-05T10:30:00Z",
    "status": "shipped",
    "statusText": "已发货",
    "items": [
      {
        "id": "item001",
        "productId": "p001",
        "name": "简约纯棉T恤",
        "price": 99.00,
        "quantity": 2,
        "color": "白色",
        "size": "M",
        "image": "https://cdn.example.com/products/p001.jpg"
      }
    ],
    "shipping": {
      "method": "快递配送",
      "address": "北京市朝阳区某某街道某某小区1号楼1单元101",
      "recipient": "张三",
      "phone": "13800138000",
      "trackingNumber": "SF1234567890",
      "shippingCompany": "顺丰速运"
    },
    "payment": {
      "method": "alipay",
      "methodText": "支付宝",
      "subtotal": 198.00,
      "shipping": 10.00,
      "tax": 0.00,
      "total": 208.00,
      "paidAt": "2025-01-05T10:35:00Z"
    },
    "timeline": [
      {
        "date": "2025-01-05T10:30:00Z",
        "status": "订单已创建",
        "description": "您的订单已提交"
      },
      {
        "date": "2025-01-05T10:35:00Z",
        "status": "支付成功",
        "description": "订单已支付"
      },
      {
        "date": "2025-01-05T14:20:00Z",
        "status": "商品出库",
        "description": "商品已从仓库发出"
      },
      {
        "date": "2025-01-05T16:45:00Z",
        "status": "已发货",
        "description": "快递已揽收，运单号：SF1234567890"
      }
    ]
  }
}
```

---

### 4.4 支付订单

**接口地址**: `POST /api/orders/:id/pay`  
**接口描述**: 支付订单  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 订单ID |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| paymentMethod | string | 是 | 支付方式：alipay / wechat / union |

**请求示例**:
```json
{
  "paymentMethod": "alipay"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "paymentUrl": "https://payment.example.com/pay?orderNo=ORD20250105123456",
    "orderNo": "ORD20250105123456",
    "qrCode": "https://payment.example.com/qr/ORD20250105123456.png"
  }
}
```

---

### 4.5 取消订单

**接口地址**: `PUT /api/orders/:id/cancel`  
**接口描述**: 取消订单  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 订单ID |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| reason | string | 否 | 取消原因 |

**请求示例**:
```json
{
  "reason": "不想要了"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "取消成功",
  "data": {
    "success": true,
    "refundStatus": "processing",
    "refundAmount": 208.00
  }
}
```

---

### 4.6 获取商户订单列表

**接口地址**: `GET /api/merchant/orders`  
**接口描述**: 商户查询自己商品的订单列表  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认20 |
| status | string | 否 | 订单状态 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "orders": [
      {
        "id": "ord20250105001",
        "orderNo": "ORD20250105123456",
        "customer": {
          "name": "张三",
          "avatar": "https://cdn.example.com/avatars/u123456.jpg"
        },
        "date": "2025-01-05T10:30:00Z",
        "status": "paid",
        "statusText": "待发货",
        "total": 208.00,
        "items": [
          {
            "id": "item001",
            "name": "简约纯棉T恤",
            "quantity": 2,
            "price": 99.00
          }
        ],
        "payment": "支付宝",
        "address": "北京市朝阳区某某街道某某小区1号楼1单元101"
      }
    ],
    "total": 156,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 4.7 商户发货

**接口地址**: `PUT /api/merchant/orders/:id/ship`  
**接口描述**: 商户为订单发货  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 订单ID |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| trackingNumber | string | 是 | 物流单号 |
| shippingCompany | string | 是 | 物流公司 |

**请求示例**:
```json
{
  "trackingNumber": "SF1234567890",
  "shippingCompany": "顺丰速运"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "发货成功",
  "data": {
    "success": true,
    "status": "shipped",
    "shippedAt": "2025-01-05T16:45:00Z"
  }
}
```

---

## 5. 用户个人中心模块

### 5.1 获取用户资料

**接口地址**: `GET /api/user/profile`  
**接口描述**: 获取当前用户的个人资料  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "u123456",
    "name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "avatar": "https://cdn.example.com/avatars/u123456.jpg",
    "gender": "male",
    "birthday": "1990-01-01",
    "createdAt": "2023-01-15T08:30:00Z",
    "addresses": [
      {
        "id": "addr001",
        "name": "张三",
        "phone": "13800138000",
        "province": "北京市",
        "city": "北京市",
        "district": "朝阳区",
        "detail": "某某街道某某小区1号楼1单元101",
        "isDefault": true
      }
    ]
  }
}
```

---

### 5.2 更新用户资料

**接口地址**: `PUT /api/user/profile`  
**接口描述**: 更新用户个人资料  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 用户姓名 |
| email | string | 否 | 邮箱 |
| phone | string | 否 | 手机号 |
| gender | string | 否 | 性别：male / female / other |
| birthday | string | 否 | 生日（YYYY-MM-DD） |

**请求示例**:
```json
{
  "name": "张三",
  "phone": "13800138000",
  "gender": "male"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "success": true,
    "user": {
      "id": "u123456",
      "name": "张三",
      "phone": "13800138000"
    }
  }
}
```

---

### 5.3 上传用户头像

**接口地址**: `POST /api/user/avatar`  
**接口描述**: 上传用户头像  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| avatar | File | 是 | 头像文件（支持jpg/png，最大5MB） |

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "avatarUrl": "https://cdn.example.com/avatars/u123456.jpg"
  }
}
```

---

### 5.4 修改密码

**接口地址**: `PUT /api/user/password`  
**接口描述**: 修改用户密码  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| oldPassword | string | 是 | 旧密码 |
| newPassword | string | 是 | 新密码（至少8位） |

**请求示例**:
```json
{
  "oldPassword": "OldPassword123!",
  "newPassword": "NewPassword456!"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "修改成功",
  "data": {
    "success": true
  }
}
```

---

### 5.5 获取用户订单历史

**接口地址**: `GET /api/user/orders`  
**接口描述**: 获取用户的订单历史记录  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认20 |
| status | string | 否 | 订单状态 |

**响应示例**: 同订单列表接口

---

## 6. 商户管理模块

### 6.1 获取商户数据统计

**接口地址**: `GET /api/merchant/dashboard`  
**接口描述**: 获取商户的数据统计概览  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "productsCount": 128,
    "salesCount": 1024,
    "revenue": 98765.43,
    "ordersCount": 856,
    "pendingOrders": 23,
    "todaySales": 5432.10,
    "todayOrders": 45
  }
}
```

---

### 6.2 获取商户数据分析

**接口地址**: `GET /api/merchant/analytics`  
**接口描述**: 获取商户的数据分析报表  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| timeRange | string | 否 | 时间范围：7days / 30days / 90days / year，默认7days |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "salesData": [
      { "date": "2025-01-01", "sales": 1200.00, "orders": 12 },
      { "date": "2025-01-02", "sales": 1800.00, "orders": 18 },
      { "date": "2025-01-03", "sales": 1500.00, "orders": 15 }
    ],
    "categorySales": [
      { "category": "上衣", "sales": 42, "percentage": 35, "revenue": 4158.00 },
      { "category": "裤装", "sales": 28, "percentage": 23, "revenue": 4452.00 }
    ],
    "topProducts": [
      {
        "id": "p001",
        "name": "简约纯棉T恤",
        "sales": 120,
        "revenue": 11880.00,
        "image": "https://cdn.example.com/products/p001.jpg"
      }
    ]
  }
}
```

---

### 6.3 获取商户信息

**接口地址**: `GET /api/merchant/profile`  
**接口描述**: 获取商户的店铺信息  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "m001",
    "name": "时尚优选",
    "avatar": "https://cdn.example.com/merchants/m001.jpg",
    "description": "专注时尚服饰，品质保证",
    "productsCount": 128,
    "salesCount": 1024,
    "rating": 4.8,
    "reviewsCount": 456,
    "createdAt": "2022-12-05T10:00:00Z"
  }
}
```

---

### 6.4 更新商户信息

**接口地址**: `PUT /api/merchant/profile`  
**接口描述**: 更新商户店铺信息  
**是否需要认证**: 是（商户）

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 商户名称 |
| description | string | 否 | 店铺描述 |
| avatar | string | 否 | 店铺头像URL |

**请求示例**:
```json
{
  "name": "时尚优选旗舰店",
  "description": "专注时尚服饰，品质保证，正品直销"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "success": true
  }
}
```

---

## 7. 卖家店铺模块

### 7.1 获取卖家信息

**接口地址**: `GET /api/sellers/:id`  
**接口描述**: 获取指定卖家的店铺信息  
**是否需要认证**: 否

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 卖家ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "m001",
    "name": "时尚优选",
    "avatar": "https://cdn.example.com/merchants/m001.jpg",
    "description": "专注时尚服饰，品质保证",
    "totalSales": 1024,
    "rating": 4.8,
    "reviewsCount": 456,
    "joinDate": "2022-12-05",
    "productsCount": 128
  }
}
```

---

### 7.2 获取卖家商品列表

**接口地址**: `GET /api/sellers/:id/products`  
**接口描述**: 获取指定卖家的商品列表  
**是否需要认证**: 否

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 卖家ID |

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认20 |
| category | string | 否 | 商品分类 |
| sortBy | string | 否 | 排序方式 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "products": [
      {
        "id": "p001",
        "name": "简约纯棉T恤",
        "price": 99.00,
        "image": "https://cdn.example.com/products/p001.jpg",
        "rating": 4.5,
        "sales": 1250
      }
    ],
    "total": 128,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## 8. 系统管理模块

### 8.1 获取系统数据统计

**接口地址**: `GET /api/admin/dashboard`  
**接口描述**: 获取系统整体数据统计  
**是否需要认证**: 是（管理员）

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "totalUsers": 2568,
    "totalMerchants": 156,
    "totalCustomers": 2412,
    "totalOrders": 15689,
    "totalRevenue": 1256789.45,
    "totalProducts": 8754,
    "todayOrders": 234,
    "todayRevenue": 45678.90,
    "activeUsers": 1856
  }
}
```

---

### 8.2 获取用户列表

**接口地址**: `GET /api/admin/users`  
**接口描述**: 管理员获取所有用户列表  
**是否需要认证**: 是（管理员）

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认20 |
| type | string | 否 | 用户类型：customer / merchant |
| status | string | 否 | 账户状态：active / disabled |
| keyword | string | 否 | 搜索关键词（姓名/邮箱） |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "users": [
      {
        "id": "u123456",
        "name": "张三",
        "email": "zhangsan@example.com",
        "type": "customer",
        "status": "active",
        "registerDate": "2023-01-15",
        "lastLogin": "2025-01-05T10:30:00Z",
        "avatar": "https://cdn.example.com/avatars/u123456.jpg",
        "orderCount": 15,
        "totalSpent": 2345.67
      }
    ],
    "total": 2568,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 8.3 获取用户详情

**接口地址**: `GET /api/admin/users/:id`  
**接口描述**: 管理员获取指定用户的详细信息  
**是否需要认证**: 是（管理员）

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 用户ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "u123456",
    "name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "type": "customer",
    "status": "active",
    "registerDate": "2023-01-15T08:30:00Z",
    "lastLogin": "2025-01-05T10:30:00Z",
    "avatar": "https://cdn.example.com/avatars/u123456.jpg",
    "orderCount": 15,
    "totalSpent": 2345.67,
    "recentOrders": [
      {
        "orderNo": "ORD20250105123456",
        "date": "2025-01-05T10:30:00Z",
        "total": 208.00,
        "status": "completed"
      }
    ]
  }
}
```

---

### 8.4 更新用户状态

**接口地址**: `PUT /api/admin/users/:id/status`  
**接口描述**: 管理员更新用户账户状态（启用/禁用）  
**是否需要认证**: 是（管理员）

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 用户ID |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | string | 是 | 账户状态：active / disabled |
| reason | string | 否 | 操作原因 |

**请求示例**:
```json
{
  "status": "disabled",
  "reason": "违规操作"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "success": true,
    "status": "disabled"
  }
}
```

---

### 8.5 获取标签分析数据

**接口地址**: `GET /api/admin/tags/analytics`  
**接口描述**: 获取系统标签使用分析数据  
**是否需要认证**: 是（管理员）

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "topTags": [
      { "name": "舒适", "count": 1245, "percentage": 18 },
      { "name": "时尚", "count": 1120, "percentage": 16 },
      { "name": "百搭", "count": 980, "percentage": 14 }
    ],
    "tagsByCategory": [
      {
        "category": "上衣",
        "topTags": ["舒适", "透气", "百搭", "简约"]
      },
      {
        "category": "裤装",
        "topTags": ["修身", "显瘦", "百搭", "时尚"]
      }
    ],
    "tagGrowth": [
      { "name": "新品", "growth": 25 },
      { "name": "时尚", "growth": 15 }
    ]
  }
}
```

---

### 8.6 获取订单分析数据

**接口地址**: `GET /api/admin/orders/analytics`  
**接口描述**: 获取系统订单分析数据  
**是否需要认证**: 是（管理员）

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| timeRange | string | 否 | 时间范围：7days / 30days / 90days / year |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "ordersByStatus": [
      { "status": "已完成", "count": 9845, "percentage": 62 },
      { "status": "已发货", "count": 2568, "percentage": 16 },
      { "status": "待发货", "count": 1256, "percentage": 8 }
    ],
    "ordersByPayment": [
      { "method": "支付宝", "count": 7845, "percentage": 50 },
      { "method": "微信支付", "count": 6258, "percentage": 40 }
    ],
    "recentSales": [
      { "date": "2025-01-01", "sales": 125689.00, "orders": 456 },
      { "date": "2025-01-02", "sales": 138456.00, "orders": 512 }
    ],
    "topCategories": [
      { "category": "上衣", "sales": 456789.00, "percentage": 36 },
      { "category": "裤装", "sales": 325678.00, "percentage": 26 }
    ]
  }
}
```

---

## 9. 文件上传模块

### 9.1 上传图片

**接口地址**: `POST /api/upload/image`  
**接口描述**: 上传图片文件  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | File | 是 | 图片文件（支持jpg/png/gif，最大10MB） |

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "url": "https://cdn.example.com/uploads/2025/01/05/abc123.jpg",
    "thumbnailUrl": "https://cdn.example.com/uploads/2025/01/05/abc123_thumb.jpg",
    "width": 1920,
    "height": 1080,
    "size": 245678
  }
}
```

---

### 9.2 上传文件

**接口地址**: `POST /api/upload/file`  
**接口描述**: 上传通用文件  
**是否需要认证**: 是

**请求头**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | File | 是 | 文件（最大20MB） |

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "url": "https://cdn.example.com/uploads/2025/01/05/document.pdf",
    "filename": "document.pdf",
    "size": 1245678,
    "mimeType": "application/pdf"
  }
}
```

---

## 10. 搜索与推荐模块

### 10.1 全局搜索商品

**接口地址**: `GET /api/search`  
**接口描述**: 全局搜索商品  
**是否需要认证**: 否

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| category | string | 否 | 商品分类 |
| minPrice | number | 否 | 最低价格 |
| maxPrice | number | 否 | 最高价格 |
| sortBy | string | 否 | 排序方式 |
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认20 |

**请求示例**:
```
GET /api/search?keyword=T恤&category=上衣&minPrice=50&maxPrice=200&sortBy=price-asc
```

**响应示例**:
```json
{
  "code": 200,
  "message": "搜索成功",
  "data": {
    "products": [
      {
        "id": "p001",
        "name": "简约纯棉T恤",
        "price": 99.00,
        "image": "https://cdn.example.com/products/p001.jpg",
        "category": "上衣",
        "rating": 4.5,
        "sales": 1250,
        "seller": {
          "id": "m001",
          "name": "时尚优选"
        }
      }
    ],
    "total": 45,
    "page": 1,
    "pageSize": 20,
    "keyword": "T恤"
  }
}
```

---

### 10.2 获取推荐商品

**接口地址**: `GET /api/recommendations`  
**接口描述**: 获取推荐商品列表  
**是否需要认证**: 否

**请求参数** (Query String):

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| productId | string | 否 | 基于此商品推荐 |
| userId | string | 否 | 基于用户推荐 |
| limit | number | 否 | 返回数量，默认10 |

**请求示例**:
```
GET /api/recommendations?productId=p001&limit=10
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "products": [
      {
        "id": "p002",
        "name": "轻薄防晒衬衫",
        "price": 129.00,
        "image": "https://cdn.example.com/products/p002.jpg",
        "rating": 4.6,
        "reason": "相似商品"
      }
    ]
  }
}
```

---

### 10.3 获取轮播图数据

**接口地址**: `GET /api/carousel`  
**接口描述**: 获取首页轮播图数据  
**是否需要认证**: 否

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "banners": [
      {
        "id": "banner001",
        "image": "https://cdn.example.com/banners/banner1.jpg",
        "title": "新春大促",
        "link": "/products?category=上衣",
        "order": 1
      },
      {
        "id": "banner002",
        "image": "https://cdn.example.com/banners/banner2.jpg",
        "title": "冬季新品",
        "link": "/products?tag=新品",
        "order": 2
      }
    ]
  }
}
```

---

## 通用规范

### 请求规范

#### 请求头

所有需要认证的接口必须在请求头中携带JWT令牌：

```
Authorization: Bearer {token}
Content-Type: application/json
```

文件上传接口使用：

```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

#### 请求方法

- **GET**: 查询数据
- **POST**: 创建数据
- **PUT**: 更新数据
- **DELETE**: 删除数据

### 响应规范

#### 统一响应格式

所有接口统一使用以下响应格式：

**成功响应**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 具体数据
  }
}
```

**错误响应**:
```json
{
  "code": 400,
  "message": "错误描述",
  "error": "详细错误信息"
}
```

#### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（未登录或token过期） |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如邮箱已存在） |
| 500 | 服务器内部错误 |

#### 业务错误码

| 错误码 | 说明 |
|--------|------|
| 1001 | 用户不存在 |
| 1002 | 密码错误 |
| 1003 | 邮箱已被注册 |
| 1004 | 账户已被禁用 |
| 2001 | 商品不存在 |
| 2002 | 库存不足 |
| 2003 | 商品已下架 |
| 3001 | 订单不存在 |
| 3002 | 订单状态不允许此操作 |
| 3003 | 支付失败 |
| 4001 | 权限不足 |
| 4002 | Token已过期 |
| 4003 | Token无效 |

### 分页规范

所有列表接口统一使用以下分页参数和响应格式：

**请求参数**:
```
page: 页码（从1开始）
pageSize: 每页数量
```

**响应格式**:
```json
{
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "pageSize": 20,
    "totalPages": 5
  }
}
```

### 时间格式

所有时间字段统一使用ISO 8601格式：

```
2025-01-05T10:30:00Z
```

### 金额格式

所有金额字段统一使用数字类型，保留两位小数：

```json
{
  "price": 99.00,
  "total": 208.00
}
```

### 图片URL规范

所有图片URL使用完整的HTTPS地址：

```
https://cdn.example.com/products/p001.jpg
```

### 安全规范

1. **HTTPS**: 所有接口必须使用HTTPS协议
2. **密码加密**: 密码在传输前应使用前端加密
3. **Token过期**: 访问令牌有效期2小时，刷新令牌有效期7天
4. **请求限流**: 同一IP每分钟最多请求100次
5. **SQL注入防护**: 所有输入参数需进行SQL注入检测
6. **XSS防护**: 所有用户输入需进行XSS过滤

### 接口版本控制

当前版本：v1

接口路径格式：`/api/v1/{resource}`

示例：`/api/v1/products`

---

## 附录

### 接口清单汇总

| 序号 | 接口地址 | 方法 | 说明 | 认证 |
|------|----------|------|------|------|
| 1 | /api/auth/login | POST | 用户登录 | 否 |
| 2 | /api/auth/register | POST | 用户注册 | 否 |
| 3 | /api/auth/logout | POST | 用户登出 | 是 |
| 4 | /api/auth/user | GET | 获取当前用户信息 | 是 |
| 5 | /api/auth/refresh | POST | 刷新令牌 | 否 |
| 6 | /api/auth/admin/login | POST | 管理员登录 | 否 |
| 7 | /api/products | GET | 获取商品列表 | 否 |
| 8 | /api/products/:id | GET | 获取商品详情 | 否 |
| 9 | /api/categories | GET | 获取商品分类 | 否 |
| 10 | /api/products | POST | 创建商品 | 商户 |
| 11 | /api/products/:id | PUT | 更新商品 | 商户 |
| 12 | /api/products/:id | DELETE | 删除商品 | 商户 |
| 13 | /api/merchant/products | GET | 获取商户商品列表 | 商户 |
| 14 | /api/products/generate-tags | POST | AI生成标签 | 商户 |
| 15 | /api/cart | GET | 获取购物车 | 是 |
| 16 | /api/cart | POST | 添加到购物车 | 是 |
| 17 | /api/cart/:itemId | PUT | 更新购物车商品 | 是 |
| 18 | /api/cart/:itemId | DELETE | 删除购物车商品 | 是 |
| 19 | /api/cart | DELETE | 清空购物车 | 是 |
| 20 | /api/orders | POST | 创建订单 | 是 |
| 21 | /api/orders | GET | 获取订单列表 | 是 |
| 22 | /api/orders/:id | GET | 获取订单详情 | 是 |
| 23 | /api/orders/:id/pay | POST | 支付订单 | 是 |
| 24 | /api/orders/:id/cancel | PUT | 取消订单 | 是 |
| 25 | /api/merchant/orders | GET | 获取商户订单 | 商户 |
| 26 | /api/merchant/orders/:id/ship | PUT | 商户发货 | 商户 |
| 27 | /api/user/profile | GET | 获取用户资料 | 是 |
| 28 | /api/user/profile | PUT | 更新用户资料 | 是 |
| 29 | /api/user/avatar | POST | 上传用户头像 | 是 |
| 30 | /api/user/password | PUT | 修改密码 | 是 |
| 31 | /api/user/orders | GET | 获取用户订单 | 是 |
| 32 | /api/merchant/dashboard | GET | 获取商户统计 | 商户 |
| 33 | /api/merchant/analytics | GET | 获取商户分析 | 商户 |
| 34 | /api/merchant/profile | GET | 获取商户信息 | 商户 |
| 35 | /api/merchant/profile | PUT | 更新商户信息 | 商户 |
| 36 | /api/sellers/:id | GET | 获取卖家信息 | 否 |
| 37 | /api/sellers/:id/products | GET | 获取卖家商品 | 否 |
| 38 | /api/admin/dashboard | GET | 获取系统统计 | 管理员 |
| 39 | /api/admin/users | GET | 获取用户列表 | 管理员 |
| 40 | /api/admin/users/:id | GET | 获取用户详情 | 管理员 |
| 41 | /api/admin/users/:id/status | PUT | 更新用户状态 | 管理员 |
| 42 | /api/admin/tags/analytics | GET | 获取标签分析 | 管理员 |
| 43 | /api/admin/orders/analytics | GET | 获取订单分析 | 管理员 |
| 44 | /api/upload/image | POST | 上传图片 | 是 |
| 45 | /api/upload/file | POST | 上传文件 | 是 |
| 46 | /api/search | GET | 全局搜索 | 否 |
| 47 | /api/recommendations | GET | 获取推荐商品 | 否 |
| 48 | /api/carousel | GET | 获取轮播图 | 否 |

---

**文档结束**

