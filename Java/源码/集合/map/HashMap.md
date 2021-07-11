# HashMap

## 补充1：位移运算

位移运算都是对补码进行操作，操作后的结果也为补码，最终值需要转换为原码，当移动大于等于32时，实际移动位数为位数%32

- 左移：`<<`，每一个二进制位向左移动，丢弃移出的位，最低位用0补齐，结果等价于乘2
- 右移：`>>`，每一个进行位向右移动，丢弃移出的位，最高用符号位补齐，结果等价于除2
- 无符号右移：`>>>`，操作同右移，只是最高位用0补齐，无论正负。

## 补充2：补码转原码

- 正数：补码与原码相同
- 负数：
  - 将从右向左找到的第一个1与符号位之间的所有数字按位取反。例：`10010110`是补码，符号位与最后一个1之间的所有数字按位取反，得`11101010`
  - 同求补码的操作，除符号位外，各位取反然后+1

## Map接口

map是一个包含键值对的容器，一个map不能有两个重复的键，每一个只对应一个值。在JDK中，每一个Map的实现都会有两个“标准”构造函数：1. 一个无参的构造函数。用于创建一个空的map。2. 一个参数为`Map`类型的构造函数，用于创建一个键值对与其相同的map，当然这个函数也可能用来复制map

### 不可变map

通过`Map.of`,`Map.ofEntries`和`Map.copyOf`创建的map将是一个不可变map，它具有以下特征：

- 不能添加，删除或者修改里面的键值对，否则会抛出`UnsupportedOperationException`。
- 不允许空键和空值，否则会抛出`NullPointerException`
- 如果包含的键值都是可序列化的，则map也可序列化

## HashMap

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/06/16/20210616165038.jpg)

map的实现类之一，允许存入空键或空值，是线程不安全，不保证元素顺序的。

影响HashMap性能的有两个重要参数：初始容量（`initial capacity`)和负载因子(`load factory`)。初始容量默认为16，负载因子默认为0.75。容量是指在哈希表中的元素(官方文档中叫`bucket`)的数量，负载因子则是指在被扩容前，可以存多“满”，当元素个数大于$capacity \ * load\ factory$时，则会进行扩容。如需存储大量数据，最好指定一个合适的初始容量。

索引计算，使用位运算的效率更高，采用此公式：$index = hash(key) \& (capacity - 1)$计算索引，但需要容量为2的幂（即2的几次方），当容量为2的幂时，则进行与运算时，结果为$hash(key)$的后几位。例：

![image-20210616220732575](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/06/16/20210616220732.png)

如果capacity不为2的幂时，部分索引则会不可能被计算出，而使哈希表的分布不均匀。

当出现哈希碰撞（两个不同的key计算出同一个索引）时，则在对应位置建立链表进行存储（采用尾插法），当链表的结点个数大于一定值（称作树化临界值，默认为8）并且哈希表容量达到最小树化容量（未达到此条件时，会直接扩容）时，为提升查询效率，会将链表转变为红黑树。

### 属性

```java
//存储元素的哈希表
transient Node<K,V>[] table;
//每一个键值对结点的集合
transient Set<Map.Entry<K,V>> entrySet;
//存储的键值对个数
transient int size;
//表结构被修改的次数
transient int modCount;
//数组大小的临界值 = 当前大小 * 负载因子
int threshold;
//负载因子
final float loadFactor;
```

### hash表中存储的结点（哈希桶）

```java
static class Node<K,V> implements Map.Entry<K,V> {
        final int hash;
        final K key;
        V value;
        Node<K,V> next;
    }
```

### get

```java
    public V get(Object key) {
        Node<K,V> e;
        return (e = getNode(hash(key), key)) == null ? null : e.value;
    }
    final Node<K,V> getNode(int hash, Object key) {
        Node<K,V>[] tab; Node<K,V> first, e; int n; K k;
        if ((tab = table) != null && (n = tab.length) > 0 &&
            (first = tab[(n - 1) & hash]) != null) {
            if (first.hash == hash && // always check first node
                ((k = first.key) == key || (key != null && key.equals(k))))
                return first;
            if ((e = first.next) != null) {
                if (first instanceof TreeNode)
                    //链表被转化成了红黑树
                    return ((TreeNode<K,V>)first).getTreeNode(hash, key);
                do {
                    //遍历链表
                    if (e.hash == hash &&
                        ((k = e.key) == key || (key != null && key.equals(k))))
                        return e;
                } while ((e = e.next) != null);
            }
        }
        return null;
    }

```

### put

```java
public V put(K key, V value) {
    return putVal(hash(key), key, value, false, true);
}
/**
 * @param onlyIfAbsent if true, don't change existing value
 * @param evict if false, the table is in creation mode.
 * @return previous value, or null if none
 */
final V putVal(int hash, K key, V value, boolean onlyIfAbsent, boolean evict) {
    Node<K,V>[] tab; Node<K,V> p; int n, i;
    if ((tab = table) == null || (n = tab.length) == 0)
	    //哈希表暂为被初始化
        n = (tab = resize()).length;
    if ((p = tab[i = (n - 1) & hash]) == null)
        //这个位置还没被占用，直接插入新结点
        tab[i] = newNode(hash, key, value, null);
    else {
        Node<K,V> e; K k;
        if (p.hash == hash &&
            ((k = p.key) == key || (key != null && key.equals(k))))
            //当前结点即为key对应的结点
            e = p;
        else if (p instanceof TreeNode)
            //链表已被转化为红黑树
            e = ((TreeNode<K,V>)p).putTreeVal(this, tab, hash, key, value);
        else {
            for (int binCount = 0; ; ++binCount) {
                //遍历链表
                if ((e = p.next) == null) {
                    //到达链表尾，插入新结点
                    p.next = newNode(hash, key, value, null);
                    if (binCount >= TREEIFY_THRESHOLD - 1) // -1 for 1st
                        //链表的结点数达到定义的值，转化为红黑树
                        treeifyBin(tab, hash);
                    break;
                }
                //找到对应的结点
                if (e.hash == hash &&
                    ((k = e.key) == key || (key != null && key.equals(k))))
                    break;
                p = e;
            }
        }
        if (e != null) { // existing mapping for key
            V oldValue = e.value;
            if (!onlyIfAbsent || oldValue == null)
                e.value = value;
            //回调函数，在LinkedHashMap有具体实现，此类中为空实现
            afterNodeAccess(e);
            return oldValue;
        }
    }
    ++modCount;
    if (++size > threshold)
        //达到扩容条件
        resize();
    afterNodeInsertion(evict);
    return null;
}
```

### 扩容（resize)

```java
final Node<K,V>[] resize() {
        Node<K,V>[] oldTab = table;
        int oldCap = (oldTab == null) ? 0 : oldTab.length;
        int oldThr = threshold;
        int newCap, newThr = 0;
        if (oldCap > 0) {
            if (oldCap >= MAXIMUM_CAPACITY) {
                //无法继续扩容，只能增大临界值
                threshold = Integer.MAX_VALUE;
                return oldTab;
            }
            else if ((newCap = oldCap << 1) < MAXIMUM_CAPACITY &&
                     oldCap >= DEFAULT_INITIAL_CAPACITY)
                //容量扩大两倍
                newThr = oldThr << 1; // double threshold
        }
    	//哈希表未被初始化
        else if (oldThr > 0) // initial capacity was placed in threshold
            newCap = oldThr;
        else {               // zero initial threshold signifies using defaults
            newCap = DEFAULT_INITIAL_CAPACITY;
            newThr = (int)(DEFAULT_LOAD_FACTOR * DEFAULT_INITIAL_CAPACITY);
        }
        if (newThr == 0) {
            float ft = (float)newCap * loadFactor;
            newThr = (newCap < MAXIMUM_CAPACITY && ft < (float)MAXIMUM_CAPACITY ?
                      (int)ft : Integer.MAX_VALUE);
        }
        threshold = newThr;
        @SuppressWarnings({"rawtypes","unchecked"})
        Node<K,V>[] newTab = (Node<K,V>[])new Node[newCap];
        table = newTab;
    	//将旧表的数据转移到新中，会重新进行映射（rehash)，并不是单纯的复制
        if (oldTab != null) {
            for (int j = 0; j < oldCap; ++j) {
                Node<K,V> e;
                if ((e = oldTab[j]) != null) {
                    oldTab[j] = null;
                    if (e.next == null)
                        newTab[e.hash & (newCap - 1)] = e;
                    else if (e instanceof TreeNode)
                        ((TreeNode<K,V>)e).split(this, newTab, j, oldCap);
                    else { // preserve order
                        //低位索引处形成的新链表
                        Node<K,V> loHead = null, loTail = null;
                        //高位索引处形成的新链表
                        Node<K,V> hiHead = null, hiTail = null;
                        Node<K,V> next;
                        //遍历链表，一段精妙的代码
                        //根据索引的计算公式可以得出，新的索引只会被“新增的一位”所影响
                        //且只有两种结果：1. 索引不变（新增位为0），此类称为低位索引
                        //2. 索引变为 原来的值+oldCap（新增位为1），此类称为高位索引
                        do {
                            next = e.next;
                            if ((e.hash & oldCap) == 0) {
                                //新增位为0
                                if (loTail == null)
                                    loHead = e;
                                else
                                    loTail.next = e;
                                loTail = e;
                            }
                            else {
                                if (hiTail == null)
                                    hiHead = e;
                                else
                                    hiTail.next = e;
                                hiTail = e;
                            }
                        } while ((e = next) != null);
                        if (loTail != null) {
                            loTail.next = null;
                            newTab[j] = loHead;
                        }
                        if (hiTail != null) {
                            hiTail.next = null;
                            newTab[j + oldCap] = hiHead;
                        }
                    }
                }
            }
        }
        return newTab;
    }
```

### 树化（treeifyBin)

为提升查询效率，当某一个索引值下的链表结点达到树化临界值（默认为8），并且哈希表的容量大于等于最小树化容量（此值最小应为4 * 树化临界值，默认为64)时，会将链表转化为红黑树

```java
final void treeifyBin(Node<K,V>[] tab, int hash) {
        int n, index; Node<K,V> e;
        if (tab == null || (n = tab.length) < MIN_TREEIFY_CAPACITY)
            resize();
        else if ((e = tab[index = (n - 1) & hash]) != null) {
            TreeNode<K,V> hd = null, tl = null;
            do {
                TreeNode<K,V> p = replacementTreeNode(e, null);
                //将链表转化为双向链表
                if (tl == null)
                    hd = p;
                else {
                    p.prev = tl;
                    tl.next = p;
                }
                tl = p;
            } while ((e = e.next) != null);
            if ((tab[index] = hd) != null)
                //转化为红黑树
                hd.treeify(tab);
        }
    }
```

### 移除元素（remove)

```java
public V remove(Object key) {
        Node<K,V> e;
        return (e = removeNode(hash(key), key, null, false, true)) == null ?
            null : e.value;
    }
    /**
     * @param matchValue if true only remove if value is equal
     * @param movable if false do not move other nodes while removing
     * @return the node, or null if none
     */
    final Node<K,V> removeNode(int hash, Object key, Object value,
                               boolean matchValue, boolean movable) {
        Node<K,V>[] tab; Node<K,V> p; int n, index;
        if ((tab = table) != null && (n = tab.length) > 0 &&
            (p = tab[index = (n - 1) & hash]) != null) {
            Node<K,V> node = null, e; K k; V v;
            if (p.hash == hash &&
                ((k = p.key) == key || (key != null && key.equals(k))))
                node = p;
            else if ((e = p.next) != null) {
                if (p instanceof TreeNode)
                    node = ((TreeNode<K,V>)p).getTreeNode(hash, key);
                else {
                    do {
                        if (e.hash == hash &&
                            ((k = e.key) == key ||
                             (key != null && key.equals(k)))) {
                            node = e;
                            break;
                        }
                        p = e;
                    } while ((e = e.next) != null);
                }
            }
            if (node != null && (!matchValue || (v = node.value) == value ||
                                 (value != null && value.equals(v)))) {
                if (node instanceof TreeNode)
                    ((TreeNode<K,V>)node).removeTreeNode(this, tab, movable);
                else if (node == p)
                    tab[index] = node.next;
                else
                    p.next = node.next;
                ++modCount;
                --size;
                afterNodeRemoval(node);
                return node;
            }
        }
        return null;
    }
```

