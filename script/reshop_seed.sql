USE reshop;

DELETE FROM seckill_products
WHERE activity_id IN (1001, 1002);

DELETE FROM seckill_activities
WHERE id IN (1001, 1002);

DELETE FROM product_skus
WHERE product_id IN (2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010);

DELETE FROM product_stock
WHERE product_id IN (2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010);

DELETE FROM products
WHERE id IN (2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010);

DELETE FROM brands
WHERE id IN (301, 302, 303, 304);

DELETE FROM categories
WHERE id IN (101, 102, 103, 104);

DELETE FROM users
WHERE id IN (1, 2, 3, 4);

INSERT INTO users (id, username, password, nickname, avatar, phone, role, status)
VALUES
    (1, 'admin', '$2a$10$PIpzO.5dNhQg7zvcynpnTuRWu2n7.UP43GULmvbN8FjQpEacMLAwO', 'shop_admin', '/static/img/fj1.jpg', '13800000001', 1, 1),
    (2, 'alice', '$2a$10$PIpzO.5dNhQg7zvcynpnTuRWu2n7.UP43GULmvbN8FjQpEacMLAwO', 'alice', '/static/img/krm.jpg', '13800000002', 0, 1),
    (3, 'bob', '$2a$10$PIpzO.5dNhQg7zvcynpnTuRWu2n7.UP43GULmvbN8FjQpEacMLAwO', 'bob', '/static/img/krm2.jpg', '13800000003', 0, 1),
    (4, 'cathy', '$2a$10$PIpzO.5dNhQg7zvcynpnTuRWu2n7.UP43GULmvbN8FjQpEacMLAwO', 'cathy', '/static/img/krm3.jpg', '13800000004', 0, 1);

INSERT INTO categories (id, name, parent_id, sort)
VALUES
    (101, 'beauty', 0, 1),
    (102, 'digital', 0, 2),
    (103, 'snacks', 0, 3),
    (104, 'home', 0, 4);

INSERT INTO brands (id, name, logo, description, sort, status)
VALUES
    (301, 'PureLab', '/static/img/krm.jpg', 'skin care and personal care', 1, 1),
    (302, 'VoltCore', '/static/img/krm3.jpg', 'digital accessories and gadgets', 2, 1),
    (303, 'SnackTime', '/static/img/krm5.jpg', 'snacks and instant drinks', 3, 1),
    (304, 'WarmNest', '/static/img/krm7.jpg', 'home and daily lifestyle', 4, 1);

INSERT INTO products (id, name, title, spu_code, category_id, brand_id, price, original_price, cover_image, detail, status)
VALUES
    (2001, 'Hydrating Mask Set', 'night care and deep hydration', 'SPU-BEA-001', 101, 301, 79.00, 99.00, '/static/img/krm.jpg', 'A daily hydration mask set for skin care.', 1),
    (2002, 'Amino Facial Cleanser', 'gentle clean and soft finish', 'SPU-BEA-002', 101, 301, 59.00, 79.00, '/static/img/krm2.jpg', 'A gentle cleanser for morning and evening use.', 1),
    (2003, '65W GaN Charger', 'multi port fast charging', 'SPU-DIG-001', 102, 302, 129.00, 169.00, '/static/img/krm3.jpg', 'A compact charger for phones and tablets.', 1),
    (2004, 'Wireless Earbuds', 'daily commute and long battery', 'SPU-DIG-002', 102, 302, 299.00, 399.00, '/static/img/krm4.jpg', 'Wireless earbuds for commute and workout use.', 1),
    (2005, 'Freeze Dried Coffee', 'easy brew and energy boost', 'SPU-SNK-001', 103, 303, 39.90, 49.90, '/static/img/krm5.jpg', 'Portable black coffee for office and home.', 1),
    (2006, 'Nut Energy Gift Box', 'daily nuts and mixed nutrition', 'SPU-SNK-002', 103, 303, 89.00, 119.00, '/static/img/krm6.jpg', 'A gift box of mixed nuts for snacks and sharing.', 1),
    (2007, 'Aroma Diffuser', 'quiet mist and warm night light', 'SPU-HOM-001', 104, 304, 149.00, 199.00, '/static/img/krm7.jpg', 'Aroma diffuser for bedroom and office use.', 1),
    (2008, 'Cotton Towel Set', 'soft absorbent and quick dry', 'SPU-HOM-002', 104, 304, 69.00, 89.00, '/static/img/krm8.jpg', 'A practical cotton towel set for home use.', 1),
    (2009, 'Vitamin Essence', 'brightening and light texture', 'SPU-BEA-003', 101, 301, 139.00, 179.00, '/static/img/krmm.jpg', 'Lightweight essence for daily care.', 0),
    (2010, 'Mini Power Bank', 'small size and travel ready', 'SPU-DIG-003', 102, 302, 99.00, 129.00, '/static/img/fj1.jpg', 'Portable mini power bank for short trips.', 1);

INSERT INTO product_skus (id, product_id, sku_code, name, specs, price, original_price, stock, status)
VALUES
    (4001, 2001, 'SKU-BEA-001-A', 'Hydrating Mask Set 5pcs', '{"size":"5pcs","type":"basic"}', 79.00, 99.00, 120, 1),
    (4002, 2001, 'SKU-BEA-001-B', 'Hydrating Mask Set 10pcs', '{"size":"10pcs","type":"gift"}', 149.00, 189.00, 80, 1),
    (4003, 2002, 'SKU-BEA-002-A', 'Amino Facial Cleanser 120g', '{"weight":"120g"}', 59.00, 79.00, 100, 1),
    (4004, 2002, 'SKU-BEA-002-B', 'Amino Facial Cleanser 200g', '{"weight":"200g"}', 89.00, 109.00, 80, 1),
    (4005, 2003, 'SKU-DIG-001-A', '65W GaN Charger CN Plug', '{"plug":"CN","ports":"2C1A"}', 129.00, 169.00, 70, 1),
    (4006, 2003, 'SKU-DIG-001-B', '65W GaN Charger EU Plug', '{"plug":"EU","ports":"2C1A"}', 139.00, 179.00, 50, 1),
    (4007, 2004, 'SKU-DIG-002-A', 'Wireless Earbuds White', '{"color":"white"}', 299.00, 399.00, 45, 1),
    (4008, 2004, 'SKU-DIG-002-B', 'Wireless Earbuds Black', '{"color":"black"}', 299.00, 399.00, 45, 1),
    (4009, 2005, 'SKU-SNK-001-A', 'Freeze Dried Coffee 12 cups', '{"pack":"12"}', 39.90, 49.90, 180, 1),
    (4010, 2005, 'SKU-SNK-001-B', 'Freeze Dried Coffee 30 cups', '{"pack":"30"}', 89.90, 109.90, 120, 1),
    (4011, 2006, 'SKU-SNK-002-A', 'Nut Energy Gift Box Standard', '{"weight":"900g"}', 89.00, 119.00, 90, 1),
    (4012, 2006, 'SKU-SNK-002-B', 'Nut Energy Gift Box Deluxe', '{"weight":"1500g"}', 139.00, 169.00, 60, 1),
    (4013, 2007, 'SKU-HOM-001-A', 'Aroma Diffuser White', '{"color":"white","tank":"300ml"}', 149.00, 199.00, 65, 1),
    (4014, 2007, 'SKU-HOM-001-B', 'Aroma Diffuser Green', '{"color":"green","tank":"300ml"}', 159.00, 209.00, 55, 1),
    (4015, 2008, 'SKU-HOM-002-A', 'Cotton Towel Set 3pcs', '{"pack":"3pcs","color":"grey"}', 69.00, 89.00, 110, 1),
    (4016, 2008, 'SKU-HOM-002-B', 'Cotton Towel Set 5pcs', '{"pack":"5pcs","color":"beige"}', 99.00, 129.00, 90, 1),
    (4017, 2009, 'SKU-BEA-003-A', 'Vitamin Essence 30ml', '{"volume":"30ml"}', 139.00, 179.00, 50, 0),
    (4018, 2009, 'SKU-BEA-003-B', 'Vitamin Essence 50ml', '{"volume":"50ml"}', 199.00, 239.00, 30, 0),
    (4019, 2010, 'SKU-DIG-003-A', 'Mini Power Bank White', '{"color":"white","capacity":"10000mAh"}', 99.00, 129.00, 95, 1),
    (4020, 2010, 'SKU-DIG-003-B', 'Mini Power Bank Blue', '{"color":"blue","capacity":"10000mAh"}', 109.00, 139.00, 85, 1);

INSERT INTO product_stock (product_id, stock, locked_stock)
VALUES
    (2001, 200, 0),
    (2002, 180, 0),
    (2003, 120, 0),
    (2004, 90, 0),
    (2005, 300, 0),
    (2006, 150, 0),
    (2007, 120, 0),
    (2008, 200, 0),
    (2009, 80, 0),
    (2010, 180, 0);

INSERT INTO seckill_activities (id, name, start_time, end_time, status)
VALUES
    (1001, 'night_flash_sale', '2026-04-27 18:00:00', '2026-04-27 23:59:59', 1),
    (1002, 'weekend_preview_sale', '2026-05-01 10:00:00', '2026-05-01 22:00:00', 0);

INSERT INTO seckill_products (activity_id, product_id, seckill_price, seckill_stock, locked_stock, limit_per_user)
VALUES
    (1001, 2001, 49.00, 80, 0, 2),
    (1001, 2003, 99.00, 60, 0, 1),
    (1001, 2005, 19.90, 120, 0, 3),
    (1002, 2010, 79.00, 70, 0, 1);
