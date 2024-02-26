CREATE TABLE IF NOT EXISTS order_description(
    id INTEGER,
    item_id INTEGER,
    quantity INTEGER);
INSERT INTO order_description(id, item_id, quantity) VALUES(10, 1, 2);
INSERT INTO order_description(id, item_id, quantity) VALUES(10, 3, 1);
INSERT INTO order_description(id, item_id, quantity) VALUES(10, 6, 1);

INSERT INTO order_description(id, item_id , quantity) VALUES(11, 2, 3);

INSERT INTO order_description(id, item_id, quantity) VALUES(14, 1, 3);
INSERT INTO order_description(id, item_id, quantity) VALUES(14, 2, 4);

INSERT INTO order_description(id, item_id, quantity) VALUES(15, 5, 1);
