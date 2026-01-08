-- 创建分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_id BIGINT DEFAULT 0,
    level INTEGER DEFAULT 1,
    sort INTEGER DEFAULT 0,
    icon VARCHAR(255),
    image VARCHAR(255),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_status ON categories(status);

-- 创建商品表
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    original_price DECIMAL(10,2),
    stock INTEGER DEFAULT 0,
    sales INTEGER DEFAULT 0,
    image VARCHAR(255),
    images TEXT,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (seller_id) REFERENCES users(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE INDEX idx_products_seller_id ON products(seller_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_status ON products(status);

-- 创建购物车表
CREATE TABLE IF NOT EXISTS cart_items (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INTEGER DEFAULT 1,
    selected BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX idx_cart_items_user_id ON cart_items(user_id);
CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);

-- 创建订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    seller_id BIGINT NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(50),
    payment_status SMALLINT DEFAULT 0,
    order_status SMALLINT DEFAULT 0,
    shipping_address TEXT,
    shipping_name VARCHAR(100),
    shipping_phone VARCHAR(20),
    tracking_number VARCHAR(100),
    remark TEXT,
    paid_at TIMESTAMP,
    shipped_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (seller_id) REFERENCES users(id)
);

CREATE INDEX idx_orders_order_no ON orders(order_no);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_seller_id ON orders(seller_id);
CREATE INDEX idx_orders_status ON orders(order_status);

-- 创建订单项表
CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    image VARCHAR(255),
    price DECIMAL(10,2) NOT NULL,
    quantity INTEGER NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- 创建地址表
CREATE TABLE IF NOT EXISTS addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    province VARCHAR(50),
    city VARCHAR(50),
    district VARCHAR(50),
    address TEXT NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_addresses_user_id ON addresses(user_id);

-- 创建轮播图表
CREATE TABLE IF NOT EXISTS banners (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255),
    image VARCHAR(255) NOT NULL,
    link VARCHAR(255),
    sort INTEGER DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_banners_status ON banners(status);
CREATE INDEX idx_banners_sort ON banners(sort);

-- 创建标签表
CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50),
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tags_type ON tags(type);

-- 创建商品标签关联表
CREATE TABLE IF NOT EXISTS product_tags (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

CREATE INDEX idx_product_tags_product_id ON product_tags(product_id);
CREATE INDEX idx_product_tags_tag_id ON product_tags(tag_id);

-- 插入示例分类数据
INSERT INTO categories (name, parent_id, level, sort) VALUES
('服装', 0, 1, 1),
('男装', 1, 2, 10),
('女装', 1, 2, 20),
('童装', 1, 2, 30),
('鞋靴', 0, 1, 2),
('男鞋', 5, 2, 10),
('女鞋', 5, 2, 20),
('配饰', 0, 1, 3);

-- 插入示例轮播图数据
INSERT INTO banners (title, image, link, sort, status) VALUES
('新春特惠', 'https://example.com/banner1.jpg', '/products?tag=sale', 1, 1),
('新品上市', 'https://example.com/banner2.jpg', '/products?tag=new', 2, 1);

