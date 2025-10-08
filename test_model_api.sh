#!/bin/bash

# 模型服务API测试脚本

# 配置
BASE_URL="http://localhost:8888"
TOKEN="YOUR_ACCESS_TOKEN"  # 替换为实际的token
TEST_IMAGE="test.jpg"      # 替换为实际的测试图片

echo "========================================="
echo "模型服务API测试脚本"
echo "========================================="
echo ""

# 检查是否提供了token
if [ "$TOKEN" == "YOUR_ACCESS_TOKEN" ]; then
    echo "❌ 错误: 请先设置有效的access token"
    echo "   修改脚本中的 TOKEN 变量"
    exit 1
fi

# 检查测试图片是否存在
if [ ! -f "$TEST_IMAGE" ]; then
    echo "❌ 错误: 测试图片 $TEST_IMAGE 不存在"
    echo "   请准备一张测试图片并更新脚本中的 TEST_IMAGE 变量"
    exit 1
fi

echo "✅ 配置检查通过"
echo ""

# 1. 上传图片
echo "📤 步骤1: 上传图片..."
echo "----------------------------------------"
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/images/upload" \
  -H "Authorization: Bearer $TOKEN" \
  -F "image=@$TEST_IMAGE" \
  -F "model_name=resnet50")

echo "响应: $UPLOAD_RESPONSE"
echo ""

# 解析响应（简单的方式，实际项目建议使用jq）
IMAGE_ID=$(echo $UPLOAD_RESPONSE | grep -o '"image_id":[0-9]*' | grep -o '[0-9]*')
TASK_ID=$(echo $UPLOAD_RESPONSE | grep -o '"task_id":"[^"]*"' | sed 's/"task_id":"//;s/"//')

if [ -z "$IMAGE_ID" ] || [ -z "$TASK_ID" ]; then
    echo "❌ 上传失败，请检查响应"
    exit 1
fi

echo "✅ 上传成功!"
echo "   Image ID: $IMAGE_ID"
echo "   Task ID: $TASK_ID"
echo ""

# 2. 查询任务状态
echo "🔍 步骤2: 查询任务状态..."
echo "----------------------------------------"
sleep 1  # 等待一秒

TASK_RESPONSE=$(curl -s -X GET "$BASE_URL/api/tasks/$TASK_ID/status" \
  -H "Authorization: Bearer $TOKEN")

echo "响应: $TASK_RESPONSE"
echo ""

# 3. 获取任务列表
echo "📋 步骤3: 获取任务列表..."
echo "----------------------------------------"
TASK_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/tasks?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

echo "响应: $TASK_LIST_RESPONSE"
echo ""

# 4. 模拟模型服务回调（这一步通常由模型服务执行）
echo "🤖 步骤4: 模拟模型服务回调..."
echo "----------------------------------------"
CALLBACK_RESPONSE=$(curl -s -X POST "$BASE_URL/api/model/callback" \
  -H "Content-Type: application/json" \
  -d "{
    \"task_id\": \"$TASK_ID\",
    \"status\": \"completed\",
    \"progress\": 100,
    \"results\": [
      {
        \"name\": \"测试标签1\",
        \"code\": \"test_label_1\",
        \"confidence\": 0.95
      },
      {
        \"name\": \"测试标签2\",
        \"code\": \"test_label_2\",
        \"confidence\": 0.87,
        \"bbox\": {
          \"x\": 100,
          \"y\": 150,
          \"width\": 200,
          \"height\": 180
        }
      }
    ]
  }")

echo "响应: $CALLBACK_RESPONSE"
echo ""

# 5. 再次查询任务状态（应该已完成）
echo "🔍 步骤5: 再次查询任务状态（应该已完成）..."
echo "----------------------------------------"
sleep 1

TASK_RESPONSE_2=$(curl -s -X GET "$BASE_URL/api/tasks/$TASK_ID/status" \
  -H "Authorization: Bearer $TOKEN")

echo "响应: $TASK_RESPONSE_2"
echo ""

# 6. 获取图片标签
echo "🏷️  步骤6: 获取图片标签..."
echo "----------------------------------------"
LABELS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/images/$IMAGE_ID/labels" \
  -H "Authorization: Bearer $TOKEN")

echo "响应: $LABELS_RESPONSE"
echo ""

# 7. 获取图片列表
echo "🖼️  步骤7: 获取图片列表..."
echo "----------------------------------------"
IMAGE_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/images?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

echo "响应: $IMAGE_LIST_RESPONSE"
echo ""

echo "========================================="
echo "✅ 所有测试完成!"
echo "========================================="
echo ""
echo "总结:"
echo "  - Image ID: $IMAGE_ID"
echo "  - Task ID: $TASK_ID"
echo ""
echo "你可以在浏览器中访问上传的图片:"
echo "  $BASE_URL/uploads/[文件路径]"
echo ""
