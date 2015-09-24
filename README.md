# saseat: simulated annealing seating charts

This is a one-off program I wrote to assign seats at my wedding reception.  I
don't intend to maintain it, but it may be useful or interesting to you.

Making a wedding seating chart is notoriously difficult and time-consuming,
especially if you fully assign seats instead of only assigning guests to a
table.  But it's really a combinatorial optimization problem.  I decided I would
enjoy spending a few hours coding more than optimizing by hand.

The algorithm used here is simulated annealing, mostly because it works
surprisingly well and is very easy to implement.  Because it is so simple, it is
also quite fast, which means you can take advantage of the power of brute force.
After a few rounds of tuning, I let it run overnight.  In the morning, I had a
seating chart that was the result of evaluating 3.05 billion candidate charts.
It was very good.

An interesting twist on the standard algorithm is that it does not swap only
single seats; up to an entire side will be swapped.  Also, it spends 90% of its
time optimizing within a table, only occasionally moving groups between tables.

## Utility function

For any candidate seating chart, the room's score is the sum of the guest
preferences and table scores.

Guest preferences are a bilateral score of how much a pair of guests will enjoy
being seated near each other.

Table balance is a measure of how well-balanced the seating arrangement is.  A
seating chart will be penalized if it has everyone on one side of a table, too
many people crowded around a table, or too few people at a table.

### Guest preferences

For each guest, all the other members at the table are considered according to
their bilateral preferences.

* Guests seated next to each other are weighted at 100%.
* Guests seated directly across are weighted at 60%.
* Guests seated diagonally are weighted at 30%.
* Guests seated at the same table in any other position are weighted at 10%.

When guests are adjacent and have the same gender, we also apply a small
penalty.

### Table balance

Table score is zero or negative, and is the sum of two scores.

Table capacity score is based on the ideal of a table set for two less than full
capacity.  A table with more or fewer guests will be penalized.

Left/right balance per table is an exponential penalty if the number of guests
on each side of the table is unequal.

## Limitations

There is no attempt to generalize beyond what I needed for my wedding.

Only rectangular tables are supported.  Many weddings use round banquet-style
tables.

All tables are considered equal.  You could decide later where tables should be
placed in the room.

For simplicity, guest preferences are bilateral, even though in reality people
like or dislike each other asymmetrically.

Weights and algorithms are hardcoded.

## Input

### Guest list

A CSV file with one or two fields.

1. Name; any unique Unicode string.  If it starts with "Mr.", "Mrs.", "Ms.", or
   "Miss", that will be used to default the gender field to "male" or "female".
2. Gender; also any Unicode string, but not unique.

### Neighbor preferences

A CSV file with three fields.  The first two are names, which must exactly match
either a name or nickname from the guest list.  The third field is a real number
given these guests' preference to be near each other; this may be negative.

### Table layout

A command-line flag, for example `--tables=12,12,10,12,10` for five tables
seating 56 in total.

## Output

The final output is a guest list with table assignments fully populated.

`saseat` prints its current state once a second.  It runs indefinitely.
