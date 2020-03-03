### Find all labels and sort by descending email count

```
db.mails.aggregate([
    { $unwind : "$labels" },
	{ $group: {_id: "$labels", total: {$sum : 1} } },
	{ $sort : {"total": -1 } }
])
```
