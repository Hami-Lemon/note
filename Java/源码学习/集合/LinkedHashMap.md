# LinkedHashMap

有序的hashmap，基于一个hashmap和一个双向链表实现有序，默认的排序方式是插入顺序，先插入的元素排在前面，通过构造函数设置以访问顺序进行排序，最近访问的元素会排在后面（等同与LRU算法中的哈希链表）。

## Example

```java
Map<Integer, Integer> map = new LinkedHashMap<>(10, 0.75f, true);//以访问顺序排序
map.put(4, 1);
map.put(3, 2);
map.put(2, 3);
map.put(1, 4);
map.forEach((k, v) -> {
   System.out.println(k + "  " + v); //4=1,3=2,2=3,1=4
});
map.get(4);
System.out.println(" ");
map.forEach((k, v) -> {
    System.out.println(k + "  " + v);//3=2,2=3,1=4,4=1
});
```

## 源码

`Entry`中维护一个`before,after`作为该节点的前驱节点和后继节点。

```java
static class Entry<K,V> extends HashMap.Node<K,V> {
        Entry<K,V> before, after;
        Entry(int hash, K key, V value, Node<K,V> next) {
            super(hash, key, value, next);
        }
    }
```

同时在`LinkedHashMap`中的`head,tail`则作为双向链表的头节点和尾节点。

```java
   transient LinkedHashMap.Entry<K,V> head;
   transient LinkedHashMap.Entry<K,V> tail;
```

`put,` `get`均是调用`HashMap`中的函数，在`LinkedHashMap`中重载了`afterNodeInsertion`，`afterNodeAccess`

### afterNodeAccess

```java
void afterNodeAccess(Node<K,V> e) { // move node to last
        LinkedHashMap.Entry<K,V> last;
        if (accessOrder && (last = tail) != e) {
            LinkedHashMap.Entry<K,V> p =
                (LinkedHashMap.Entry<K,V>)e, b = p.before, a = p.after;
            p.after = null;
            if (b == null)
                head = a;
            else
                b.after = a;
            if (a != null)
                a.before = b;
            else
                last = b;
            if (last == null)
                head = p;
            else {
                p.before = last;
                last.after = p;
            }
            tail = p;
            ++modCount;
        }
    }
```

