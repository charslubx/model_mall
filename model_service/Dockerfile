# 使用Python官方镜像作为基础镜像
FROM python:3.10-slim

# 设置工作目录
WORKDIR /app

# 安装系统依赖
RUN apt-get update && apt-get install -y \
    libgl1-mesa-glx \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

# 复制依赖文件
COPY requirements.txt .

# 安装Python依赖
RUN pip install --no-cache-dir -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 复制应用代码
COPY app/ ./app/

# 创建模型目录
RUN mkdir -p /app/models

# 设置环境变量
ENV MODEL_PATH=/app/models/model.h5
ENV MODEL_NAME=image-classifier
ENV PORT=5000
ENV HOST=0.0.0.0
ENV PYTHONUNBUFFERED=1

# 暴露端口
EXPOSE 5000

# 使用gunicorn启动服务（生产环境）
CMD ["gunicorn", "--bind", "0.0.0.0:5000", "--workers", "4", "--threads", "2", "--timeout", "120", "--worker-class", "gevent", "--access-logfile", "-", "--error-logfile", "-", "app.api:app"]

# 开发环境可以使用Flask自带服务器
# CMD ["python", "-m", "app.api"]
