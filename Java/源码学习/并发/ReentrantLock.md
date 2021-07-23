# ReentrantLock(重入锁)

一个线程获取锁后可以反复的加锁，不会出现自己阻塞自己的情况。但在释放锁时也得多次释放，获取了几次锁，就得释放几次锁。

获取锁时会判断当前线程是否为获取锁的线程，如果是则同步状态+1，释放锁时同步状态-1。只有同步状态置为0时才会最终释放锁。

ReentrantLock分为公平锁和非公平锁（默认创建非公平锁），公平锁表示在锁可用时，让等待时间最长的线程获取锁，而非公平则是随机分配。相比而言，非公平锁的性能表现更好一点，但公平锁能够避免饥饿问题。

## 使用

```java
//创建锁，默认为非公平锁
ReentrantLock lock = new ReentrantLock();
//ReentrantLock lock = new ReentrantLock(true); 创建公平锁
//获取锁，获取失败则会被阻塞
lock.lock();
//释放锁
lock.unlock();
```

## 实现分析

ReentrantLock中有一个Sync的静态内部类，该类继承AQS，并实现相应的方法

### 获取锁

获取锁是直接调用 AQS的aqcuqire方法，因此主要内容在对tryAcquire方法的实现。

#### 非公平锁

```java
static final class NonfairSync extends Sync {
        private static final long serialVersionUID = 7316153563782823691L;
        protected final boolean tryAcquire(int acquires) {
            return nonfairTryAcquire(acquires);
        }
}
final boolean nonfairTryAcquire(int acquires) {
            final Thread current = Thread.currentThread();
    //获取AQS中的state，此处可理解为锁
            int c = getState();
    //为0时表示还没有线程获取锁
            if (c == 0) {
                //如果cas失败，可能是被其它线程获取锁，则当前线程获取失败
                //这里也可以看出非公平锁获取锁是一个竞争过程
                if (compareAndSetState(0, acquires)) {
                    setExclusiveOwnerThread(current);
                    return true;
                }
            }
    //如果锁已经被获取过，则判断是否是同一个线程再次加锁（重入）
            else if (current == getExclusiveOwnerThread()) {
                int nextc = c + acquires;
                if (nextc < 0) // overflow
                    throw new Error("Maximum lock count exceeded");
                setState(nextc);
                return true;
            }
            return false;
}
```

#### 公平锁

```java
static final class FairSync extends Sync {
        private static final long serialVersionUID = -3000897897090466540L;
        
        @ReservedStackAccess
        protected final boolean tryAcquire(int acquires) {
            final Thread current = Thread.currentThread();
            int c = getState();
            if (c == 0) {
                //和非公平锁相同，但多调用了一个hasQueuedPredecessors方法
                //当不存在比当前等待时间更长的线程时，则让当前线程获取锁
                if (!hasQueuedPredecessors() &&
                    compareAndSetState(0, acquires)) {
                    setExclusiveOwnerThread(current);
                    return true;
                }
            }
            //锁重入
            else if (current == getExclusiveOwnerThread()) {
                int nextc = c + acquires;
                if (nextc < 0)
                    throw new Error("Maximum lock count exceeded");
                setState(nextc);
                return true;
            }
            return false;
        }
}
//判断是否存在比当前线程等待时间更长的线程
public final boolean hasQueuedPredecessors() {
        Node h, s;
        if ((h = head) != null) {
            if ((s = h.next) == null || s.waitStatus > 0) {
                s = null; // traverse in case of concurrent cancellation
                for (Node p = tail; p != h && p != null; p = p.prev) {
                    if (p.waitStatus <= 0)
                        s = p;
                }
            }
            if (s != null && s.thread != Thread.currentThread())
                return true;
        }
        return false;
}
```

### 释放锁

非公平锁和公平锁的差别主要在获取锁上，释放锁的操作相同，unlock方法也是直接调用AQS的release方法

```java
//返回锁是否被完全释放
protected final boolean tryRelease(int releases) {
    // releases通常为1
            int c = getState() - releases;
    //判断释放当前线程是否是获取到锁的线程
            if (Thread.currentThread() != getExclusiveOwnerThread())
                throw new IllegalMonitorStateException();
            boolean free = false;
    //为0时表示锁被完全释放，其它线程可 
            if (c == 0) {
                free = true;
                setExclusiveOwnerThread(null);
            }
            setState(c);
            return free;
}
```



