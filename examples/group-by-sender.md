### Find all senders and sort by descending email count

Possibly you need modify labels:

```
labels: { $in: ['inbox'] }
```

You can see all your labels in [group-by-labels](https://github.com/Rpsl/mboximporter/blob/master/examples/group-by-labels.md)

```
db.mails.aggregate([
    { $match: { labels: { $in: ['inbox'] } } },
    { $unwind: "$sender" },
    { $group: {_id: "$sender", total: {$sum : 1} } },
    { $sort: {"total": -1 } }
]);
```