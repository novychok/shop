CREATE TABLE IF NOT EXISTS items(
		id BIGSERIAL PRIMARY KEY NOT NULL,
		item_name VARCHAR(100),
		main_shelf CHAR,
		other_shelfs CHAR[]);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Laptop', 'А', NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Moitor', 'А', NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Phone', 'Б', '{З,В}');
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('PC', 'Ж', NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Watch', 'Ж', NULL);
INSERT INTO items(item_name, main_shelf, other_shelfs) VALUES ('Microphone', 'Ж', NULL);

