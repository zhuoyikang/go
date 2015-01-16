2叉树

# 普通2叉树

```
var btree BinTreeI
btree = &BinTree{}
btree.Insert(BinInt(1), 1)
btree.Insert(BinInt(2), 2)
btree.Delete(BinInt(v))
```

Key类型:BinInt BinString BinFloat，自定义key类型

```
type BinInt int64

func (x BinInt) Cmp(y BinKeyI) int {
	y1 := y.(BinInt)
	return int(x - y1)
}
```


# 红黑树


```
var btree BinTreeI
btree = &RBTree{}
btree.Insert(BinInt(1), 1)
btree.Insert(BinInt(2), 2)
btree.Delete(BinInt(v))
```


# AVL树



