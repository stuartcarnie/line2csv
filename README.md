line2csv
========

change 1
change 2

Transform InfluxDB line protocol to CSV.

Test Edit.

Usage
-----

Specifying both input and output files:

```bash
$ ./line2csv input.txt -o output.csv
```

Format
------

For each line, every field will generate a separate line in the CSV.

**Input:**

```text
cpu,host=foo,region=us-west-1 user=0.50,system=0.25 1000
cpu,host=foo,region=us-west-1 user=0.10,system=0.15 1001
```

**Output:**

```csv
timestamp,key,field_name,field_value
1000,"cpu,host=foo,region=us-west-1",user,0.50
1000,"cpu,host=foo,region=us-west-1",system,0.25
1001,"cpu,host=foo,region=us-west-1",user,0.10
1001,"cpu,host=foo,region=us-west-1",system,0.15
```

