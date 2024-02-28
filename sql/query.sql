CREATE TABLE IF NOT EXISTS reserves(
    id INTEGER,
    item_id INTEGER,
    quantity INTEGER);

CREATE TABLE IF NOT EXISTS items(
    id BIGSERIAL PRIMARY KEY,
    item_name VARCHAR(100),
    main_shelf INTEGER,
    other_shelfs INTEGER[]);

CREATE TABLE IF NOT EXISTS shelfs(
    id BIGSERIAL PRIMARY KEY,
    shelf_type INTEGER,
    items INTEGER[]);

CREATE INDEX IF NOT EXISTS shelf_type_index ON shelfs (shelf_type);

INSERT INTO reserves(id, item_id, quantity) VALUES(10, 1, 2);
INSERT INTO reserves(id, item_id, quantity) VALUES(11, 2, 3);
INSERT INTO reserves(id, item_id, quantity) VALUES(14, 1, 3);
INSERT INTO reserves(id, item_id ,quantity) VALUES(10, 3, 1);
INSERT INTO reserves(id, item_id, quantity) VALUES(14, 4, 4);
INSERT INTO reserves(id, item_id, quantity) VALUES(15, 5, 1);
INSERT INTO reserves(id, item_id, quantity) VALUES(10, 6, 1);

INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Laptop', 1040, NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('TV', 1040, NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Phone', 1041, ARRAY[1047, 1042]);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('PC', 1046, NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Watch', 1046, ARRAY[1040]);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Microphone', 1046, NULL);

INSERT INTO shelfs(shelf_type, items) VALUES(1040, ARRAY[1,1,1,1,2,5]);
INSERT INTO shelfs(shelf_type, items) VALUES(1040, ARRAY[1,2,2,8,8,8]);

INSERT INTO shelfs(shelf_type, items) VALUES(1041, ARRAY[1,1,1,1,2,8]);
INSERT INTO shelfs(shelf_type, items) VALUES(1041, ARRAY[1,2,2,8,8,8]);

INSERT INTO shelfs(shelf_type, items) VALUES(1046, ARRAY[6,1,1,1,2,5]);
INSERT INTO shelfs(shelf_type, items) VALUES(1046, ARRAY[1,2,4,4,4,4]);

INSERT INTO shelfs(shelf_type, items) VALUES(1047, ARRAY[1,1,1,1,2,3]);
INSERT INTO shelfs(shelf_type, items) VALUES(1042, ARRAY[1,2,2,8,8,3]);