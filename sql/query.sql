CREATE TABLE IF NOT EXISTS reserves(
    id INTEGER,
    item_id INTEGER,
    quantity INTEGER);

CREATE TABLE IF NOT EXISTS items(
    id BIGSERIAL PRIMARY KEY,
    item_name VARCHAR(100),
    main_shelf CHAR,
    other_shelfs CHAR[]);

CREATE TABLE IF NOT EXISTS shelfs(
    id BIGSERIAL PRIMARY KEY,
    shelf_type CHAR,
    items INTEGER[]);

CREATE TABLE IF NOT EXISTS orders(
    id BIGSERIAL PRIMARY KEY,
    main_shelf CHAR,
    items INTEGER[],
    order_number INTEGER,
    quantity INTEGER);

INSERT INTO reserves(id, item_id, quantity) VALUES(10, 1, 2);
INSERT INTO reserves(id, item_id, quantity) VALUES(11, 2, 3);
INSERT INTO reserves(id, item_id, quantity) VALUES(14, 1, 3);
INSERT INTO reserves(id, item_id ,quantity) VALUES(10, 3, 1);
INSERT INTO reserves(id, item_id, quantity) VALUES(14, 4, 4);
INSERT INTO reserves(id, item_id, quantity) VALUES(15, 5, 1);
INSERT INTO reserves(id, item_id, quantity) VALUES(10, 6, 1);

INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Laptop', 'А', NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('TV', 'А', NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Phone', 'Б', '{З,В}');
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('PC', 'Ж', NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Watch', 'Ж', '{А}');
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Microphone', 'Ж', NULL);

INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[1,1,8,1,1,1,8,8,2,2,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[8,8,3,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[6,4,8,4,4,5,8,8,8,8,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[1,1,8,1,1,1,8,8,2,2,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[8,8,3,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[6,4,8,4,4,5,8,8,8,8,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[1,1,8,1,1,1,8,8,2,2,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[8,8,3,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[6,4,8,4,4,5,8,8,8,8,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('А', ARRAY[1,1,8,1,1,1,8,8,2,2,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('Б', ARRAY[8,8,3,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8,8]);
INSERT INTO shelfs(shelf_type, items) VALUES('Ж', ARRAY[6,4,8,4,4,5,8,8,8,8,8,8,8,8,8,8,8,8,8]);