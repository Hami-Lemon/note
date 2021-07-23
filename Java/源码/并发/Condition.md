# Condition

将Object中的监测方法（wait，notify，notifyAll）分解出来，	结合一个`Lock`接口的实现类来替代原来的写法。即用Lock对象代替同步方法，用Condition中的方法进行wait和signal（等价于object中的notify）。

Condition可以使一个线程阻塞，直到被其它线程唤醒。在多个线程访问同一个资源时，为保证线程安全，所以Condition需要和一个Lock对象相关联。在调用wait方法时，会阻塞当前线程，并且释放锁。

例如：生产者消费者问题

```java
   class BoundedBuffer<E> {
     final Lock lock = new ReentrantLock();
     final Condition notFull  = lock.newCondition(); 
     final Condition notEmpty = lock.newCondition(); 
  
     final Object[] items = new Object[100];
     int putptr, takeptr, count;
  
     public void put(E x) throws InterruptedException {
       lock.lock();
       try {
         while (count == items.length)
           notFull.await();
         items[putptr] = x;
         if (++putptr == items.length) putptr = 0;
         ++count;
         notEmpty.signal();
       } finally {
         lock.unlock();
       }
     }
  
     public E take() throws InterruptedException {
       lock.lock();
       try {
         while (count == 0)
           notEmpty.await();
         E x = (E) items[takeptr];
         if (++takeptr == items.length) takeptr = 0;
         --count;
         notFull.signal();
         return x;
       } finally {
         lock.unlock();
       }
     }
   }
```

Condition相比Objecct提供了更好的行为和语义性，可以保证唤醒线程的有序性，或者在执行通知时不持有锁。

## AQS中的Condition实现（ConditionObject）

在AQS中维护了两个队列，一个阻塞队列，一个条件队列。

条件队列由Condition形成，当一个线程被Condition挂起（await）时，则会被加入到条件队列中，直到被唤醒（signal）或中断，然后将其加入到阻塞队列中，等待获取锁后被执行。条件队列是一个单向队列，而阻塞队列是一个双向队列。在AQS中的Condition要求每个线程独占资源。

### await

```java
public final void await() throws InterruptedException {
    if (Thread.interrupted())
        throw new InterruptedException();
    //加入结点到条件队列尾部
    Node node = addConditionWaiter();
    //释放当前持有的锁，返回释放之前锁的状态
    int savedState = fullyRelease(node);
    int interruptMode = 0;
    //判断是否转移到阻塞队列中，没有则阻塞
    while (!isOnSyncQueue(node)) {
        LockSupport.park(this);
        //判断在阻塞过程中是否被中断
        //被中断过则返回THROW_IE
        //在唤醒后被则返回REINTERRUPT
        //0则是未被中断
        if ((interruptMode = checkInterruptWhileWaiting(node)) != 0)
            break;
    }
    //此时已进入阻塞队列，则 
    if (acquireQueued(node, savedState) && interruptMode != THROW_IE)
        interruptMode = REINTERRUPT;
    if (node.nextWaiter != null) // clean up if cancelled
        unlinkCancelledWaiters();
    if (interruptMode != 0)
        reportInterruptAfterWait(interruptMode);
}
```

- addConditionWaiter

  加入结点到条件队列，并返回新加入的结点

  ```java
  private Node addConditionWaiter() {
      //判断当前线程是否正在独占资源，该方法须自己实现
              if (!isHeldExclusively())
                  throw new IllegalMonitorStateException();
              Node t = lastWaiter;
              // If lastWaiter is cancelled, clean out.
              if (t != null && t.waitStatus != Node.CONDITION) {
                  //删除掉队列中status不是CONDITION的结点
                  unlinkCancelledWaiters();
                  t = lastWaiter;
              }
              Node node = new Node(Node.CONDITION);
              if (t == null)
                  firstWaiter = node;
              else
                  t.nextWaiter = node;
              lastWaiter = node;
              return node;
  }
  ```

- fullyRelease

  释放当前线程持有的全部资源（即释放锁），并返回释放锁之前的锁状态。

  ```java
  final int fullyRelease(Node node) {
          try {
              int savedState = getState();
              if (release(savedState))
                  return savedState;
              throw new IllegalMonitorStateException();
          } catch (Throwable t) {
              node.waitStatus = Node.CANCELLED;
              throw t;
          }
  }
  ```

- isOnSyncQueue

  判断某个结点是否已经从条件队列转移到阻塞队列中

  ```java
  final boolean isOnSyncQueue(Node node) {
          if (node.waitStatus == Node.CONDITION || node.prev == null)
              return false;
          if (node.next != null) // If has successor, it must be on queue
              return true;
          //由于在加入阻塞队列的过程中，可能CAS会失败，从而出现结点正在加入队列但并没有成功加入（加入了但没完全加入），因此需要从尾部开始遍历，确定它真的被加入了。新加入阻塞队列的结点因该离尾部更近
          return findNodeFromTail(node);
  }
  private boolean findNodeFromTail(Node node) {
          for (Node p = tail;;) {
              if (p == node)
                  return true;
              if (p == null)
                  return false;
              p = p.prev;
          }
  }
  ```

### signal

```java
public final void signal() {
    //判断当前线程是否独占资源，需自己实现
            if (!isHeldExclusively())
                throw new IllegalMonitorStateException();
            Node first = firstWaiter;
            if (first != null)
                doSignal(first);
}
```

- doSignal

  将条件队列中的结点转移到阻塞队列中，一次最多只转移一个。
  
  ```java
  private void doSignal(Node first) {
              do {
                  //从头开始移出结点
                  if ( (firstWaiter = first.nextWaiter) == null)
                      //队列空了时，设置尾指针为null
                      lastWaiter = null;
                  first.nextWaiter = null;
                  //如果该结点转移失败，则尝试转移下一个
              } while (!transferForSignal(first) &&
                       (first = firstWaiter) != null);
         }
  ```
  
- transferForSignal

  将条件队列中的结点转移到阻塞队列中，返回是否成功。

  ```java
  final boolean transferForSignal(Node node) {
          //如果设置失败，则表明该结点被取消或者被其它线程转移
          if (!node.compareAndSetWaitStatus(Node.CONDITION, 0))
              return false;
      //将结点加入到阻塞队列队尾，p为其前驱结点
          Node p = enq(node);
          int ws = p.waitStatus;
          if (ws > 0 || !p.compareAndSetWaitStatus(ws, Node.SIGNAL))
              LockSupport.unpark(node.thread);
          return true;
  }
  ```

## 被唤醒后的中断检测

线程在被park方法阻塞后，可能会由以下几种情况被唤醒

1. 其它线程调用signal方法，将结点加入到阻塞队列中，当获取到锁时，线程则被唤醒。
2. 该线程被其它线程中断
3. 在进行唤醒时（signal），结点的waitStatus为CANCELLED，或者cas操作失败

- checkInterruptWhileWaiting

  检查线程在阻塞期间是否被中断

  - THROW_IE：线程被中断
  - REINTERRUPT：线程是被转移到阻塞队列后被中断
  - 0：没有被中断

  ```java
  private int checkInterruptWhileWaiting(Node node) {
              return Thread.interrupted() ?
                  (transferAfterCancelledWait(node) ? THROW_IE : REINTERRUPT) :
                  0;
  }
  final boolean transferAfterCancelledWait(Node node) {
      //如果设置成功则结点还在条件队列中，说明是在条件队列中等待期间被中断
          if (node.compareAndSetWaitStatus(Node.CONDITION, 0)) {
              //加入到阻塞队列中
              enq(node);
              return true;
          }
          //如果cast失败，说明线程是在signal后发生的中断，此时自旋等待将其加入到阻塞队列中
          while (!isOnSyncQueue(node))
              Thread.yield();
          return false;
  }
  ```

- reportInterruptAfterWait

  根据线程是在signal前或后做出不同的操作

  ```java
  private void reportInterruptAfterWait(int interruptMode)
              throws InterruptedException {
              if (interruptMode == THROW_IE)
                  //在signal前被中断，直接抛出异常
                  throw new InterruptedException();
              else if (interruptMode == REINTERRUPT)
                  //在signal后中断，则中断
                  selfInterrupt();
  }
  ```

  