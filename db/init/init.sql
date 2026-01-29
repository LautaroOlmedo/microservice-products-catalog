CREATE DATABASE IF NOT EXISTS products_catalog_db;
USE products_catalog_db;

-- PRODUCTS
CREATE TABLE products (
                          id CHAR(36) PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
                          stock INT NOT NULL CHECK (stock >= 0),
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                          UNIQUE KEY uq_product_name (name)
) ENGINE=InnoDB;


-- ORDERS
CREATE TABLE orders (
                        id CHAR(36) PRIMARY KEY,
                        product_id CHAR(36) NOT NULL,
                        quantity INT NOT NULL CHECK (quantity > 0),
                        total DECIMAL(10,2) NOT NULL CHECK (total >= 0),
                        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        CONSTRAINT fk_orders_product
                            FOREIGN KEY (product_id)
                                REFERENCES products(id)
) ENGINE=InnoDB;


CREATE INDEX idx_orders_product_id ON orders(product_id);