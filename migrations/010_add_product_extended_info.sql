-- 创建商品颜色表
CREATE TABLE IF NOT EXISTS product_colors (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    value VARCHAR(50) NOT NULL,
    hex VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_colors_product_id ON product_colors(product_id);

-- 创建商品尺码表
CREATE TABLE IF NOT EXISTS product_sizes (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    size VARCHAR(50) NOT NULL,
    stock INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_sizes_product_id ON product_sizes(product_id);

-- 创建商品特点表
CREATE TABLE IF NOT EXISTS product_features (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    feature TEXT NOT NULL,
    sort INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_features_product_id ON product_features(product_id);

-- 创建商品规格参数表
CREATE TABLE IF NOT EXISTS product_specifications (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    spec_key VARCHAR(100) NOT NULL,
    spec_value TEXT NOT NULL,
    sort INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_specifications_product_id ON product_specifications(product_id);

-- 创建商户扩展信息表
CREATE TABLE IF NOT EXISTS merchant_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    shop_name VARCHAR(255) NOT NULL,
    shop_description TEXT,
    shop_avatar VARCHAR(255),
    business_license VARCHAR(255),
    rating DECIMAL(3,2) DEFAULT 5.0,
    total_sales INTEGER DEFAULT 0,
    products_count INTEGER DEFAULT 0,
    reviews_count INTEGER DEFAULT 0,
    join_date DATE DEFAULT CURRENT_DATE,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_merchant_profiles_user_id ON merchant_profiles(user_id);
CREATE INDEX idx_merchant_profiles_status ON merchant_profiles(status);

-- 添加订单项颜色和尺码字段
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS color VARCHAR(50);
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS size VARCHAR(50);

-- 添加商品评分字段
ALTER TABLE products ADD COLUMN IF NOT EXISTS rating DECIMAL(3,2) DEFAULT 0;
ALTER TABLE products ADD COLUMN IF NOT EXISTS reviews_count INTEGER DEFAULT 0;
