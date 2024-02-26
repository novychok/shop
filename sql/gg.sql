EXPLAIN ANALYSE SELECT * FROM gg WHERE shelf_type = 'А';

CREATE INDEX gg_shelf_type ON gg USING HASH (shelf_type);
DROP INDEX IF EXISTS gg_shelf_type;

INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);
INSERT INTO gg(shelf_type, items) VALUES('А', ARRAY[8, 8, 8, 8, 8, 8, 8, 8, 8, 1, 2, 8, 8, 1]);

-----------------
 Seq Scan on gg  (cost=0.00..333.60 rows=12436 width=88) (actual time=0.005..0.753 rows=12436 loops=1)
   Filter: (shelf_type = 'А'::bpchar)
   Rows Removed by Filter: 12
 Planning Time: 0.087 ms
 Execution Time: 0.951 ms
(5 rows)

 Seq Scan on gg  (cost=0.00..333.60 rows=12436 width=88) (actual time=0.005..0.757 rows=12436 loops=1)
   Filter: (shelf_type = 'А'::bpchar)
   Rows Removed by Filter: 12
 Planning Time: 0.087 ms
 Execution Time: 0.955 ms
(5 rows)


--------------------------------------

 Seq Scan on gg  (cost=0.00..1.20 rows=1 width=48) (actual time=0.009..0.010 rows=4 loops=1)
   Filter: (shelf_type = 'А'::bpchar)
   Rows Removed by Filter: 12
 Planning Time: 0.042 ms
 Execution Time: 0.019 ms
(5 rows)

-------
 Seq Scan on gg  (cost=0.00..1.20 rows=1 width=48) (actual time=0.006..0.007 rows=4 loops=1)
   Filter: (shelf_type = 'А'::bpchar)
   Rows Removed by Filter: 12
 Planning Time: 0.083 ms
 Execution Time: 0.015 ms
(5 rows)
