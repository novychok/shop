# How program works

1) Immitate Orders in database.
2) Retrieve Orders from DB (reserved items and their count).
3) Compair Order Id to given Args in command line.

4) Shelfs have array with given numbers, the numbers immitate Items in Shelf array,
   each number match the Item Id.

5) Create indexes on shelf types to retrieve needed shelfs faster.

6) At first program start to search Items in Item main shelf type, with limit 1
   (Maybe we have enough Items on main shelf, don't need to retrieve more data).
7) If we don't get enough data, program start to search in Other Item Shelfs ('З', 'В') with some limit
   (They are not main shelf types, so limit is for see in several places at once).
8) Still don't get enough Items ? Start a loop that will retrieve data from main Item shelfs with offset and limit,
   (offset and limit increments on each iteration while finding items, immitate passing the shelfs,
   to not have them all at once, maybe on next iteration we will have enough Items,
   if not Program will loop forever, sorry:).