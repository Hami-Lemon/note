# AQS(抽象队列同步器)

提供了一种实现阻塞锁和一系列依赖FIFO等待队列的同步器的框架。在框架中维护了一个代表资源状态的整数值（state）和一个FIFO的阻塞队列（双向队列）。阻塞队列中则保存着等待获取资源的线程。

子类应该被定义为一个非公有的内部类对象，并实现其中的方法。这个类支持*独占模式(exclusive)*和*共享模式(shared)*

。在独占模式下，其它线程获取资源时将会失败；在共享模式下，多个线程同时获取资源可能（但不一定）会成功。

使用这个类的基本功能只需要重写以下方法即可：

- `tryAcquire`独占模式获取资源
- `tryRelease`独占模式释放资源
- `tryAcquireShared`共享模式获取资源
- `tryReleaseShared`共享模式释放资源
- `isHeldExclusively`该线程是否独占资源

## 补充1: LockSupport.park

详见：[csdn](https://blog.csdn.net/weixin_39687783/article/details/85058686)

挂起当前线程，在以下几种情况下恢复

- 其它线程调用uppark，唤醒当前线程
- 当前线程被中断

## 属性

- Node

  一个静态内部类，作为阻塞队列中的每一个结点。

  - waitStatus

    当前结点的阻塞状态

  - prev

    上一个结点，第一个结点会指向头结点，并且当队列初始化后，头结点始终存在。因此，该属性始终不为null。

  - next

    指向下一个结点。

  - thread

    保存被阻塞的线程。

  - nextWaiter

    指向下一个在条件（condition）队列中的结点，或者为`SHARED`，因为只有在独占模式下才会使用。

- head

  队列的头结点

- tail

  队列的尾结点

## 结点状态（waitStatus)

- SIGNAL(-1) : 当前结点的后继结点处于阻塞状态，当前结点必须在释放资源或被取消后唤醒后继结点。
- CANCELLED(1)：结点因超时或中断而被取消，进入该状态的结点不会再转换为其它状态。
- CONDITION(-2)：表明该结点处于一个条件（condition)队列中，在被转换之前它将不会被作为一个同步队列中的结点被使用。
- PROPAGATE(-3)：在共享模式下，当该结点被释放后，不仅会唤醒后继结点，还会将唤醒操作传播下去。
- 0：新结点入队时的默认值。

## 获取资源（独占模式）

- acquire

  独占模式下获取资源的顶层入口，如果获取到资源则直接返回，否则将线程加入阻塞队列。

  ```java
  public final void acquire(int arg) {
      //tryAcquire表示尝试获取资源，具体方法由子类实现。
          if (!tryAcquire(arg) &&
              acquireQueued(addWaiter(Node.EXCLUSIVE), arg))
              selfInterrupt();
  }
  ```

- addWaiter

  将当前线程加入阻塞队列队尾，并返回对应的Node

  ```java
  private Node addWaiter(Node mode) {
      //创建一个结点，并保存当前线程。
          Node node = new Node(mode);
          for (;;) {
              Node oldTail = tail;
              if (oldTail != null) {
                  //设置结点的前驱结点
                  node.setPrevRelaxed(oldTail);
                  //cas自旋设置尾结点
                  if (compareAndSetTail(oldTail, node)) {
                      oldTail.next = node;
                      return node;
                  }
              } else {
                  // 初始化队列，以cas的方式创建头结点，并让tail也指向头结点
                  initializeSyncQueue();
              }}}
  ```

- acquireQueued

  线程尝试去获取资源，返回当前线程是否被中断

  ```java
  final boolean acquireQueued(final Node node, int arg) {
          boolean interrupted = false;
          try {
              for (;;) {
                  final Node p = node.predecessor();
                  //当前结点的前驱结点为头结点，则去尝试获取资源
                  if (p == head && tryAcquire(arg)) {
                      //成功获取资源，则当前结点作为头结点
                      setHead(node);
                      p.next = null; // help GC
                      return interrupted;
                  }
                  //资源获取失败，检查是否应当阻塞当前线程
                  if (shouldParkAfterFailedAcquire(p, node))
                      //阻塞线程，并检查是否被中断
                      interrupted |= parkAndCheckInterrupt();
              }
          } catch (Throwable t) {
              ...
          }
      }
  ```

- shouldParkAfterFailedAcquire

  检查并更新结点的状态（waitStatus），返回是否需要阻塞该结点。

  ```java
  private static boolean shouldParkAfterFailedAcquire(Node pred, Node node) {
      //前驱结点的状态
          int ws = pred.waitStatus;
          if (ws == Node.SIGNAL)
             //当前结点需要被前驱结点唤醒，则当前结点可被阻塞
              return true;
          if (ws > 0) {
              //前驱结点被取消调度，一直向前找到一个未被取消的结点
              //中间被跳过的结点将会被GC回收
              do {
                  node.prev = pred = pred.prev;
              } while (pred.waitStatus > 0);
              pred.next = node;
          } else {
              //此时的waitStatus应为 0 或 PROPAGATE. 则将其改为SIGNAL
              //并在下一次获取资源失败时，阻塞当前线程。
              pred.compareAndSetWaitStatus(ws, Node.SIGNAL);
          }
          return false;
      }
  ```

整个过程如同在医院排队取号，当排在第一个时（head的后面）则可以直接拿号走人，而在其它情况下，先查看前面的是否已经放弃（CANCELLED），如果已经放弃则可以占掉它的位置，否则就告诉前一个人唤醒自己，然后去休息（park)。

## 释放资源（独占模式）

- release

  独占模式下线程释放资源的顶层入口，会释放资源，并在资源完全释放后唤醒等待队列中的后继结点

  ```java
  public final boolean release(int arg) {
      //tryRelease同样由子类实现
          if (tryRelease(arg)) {
              //资源被完全释放
              Node h = head;
              if (h != null && h.waitStatus != 0)
                  //唤醒后继结点
                  unparkSuccessor(h);
              return true;
          }
          return false;
  }
  ```

- unparkSuccessor

  唤醒当前结点后第一个应该被唤醒的结点，在这里需要和`acquireQueued`联系起来，当s线程被唤醒后，会去判断`p == head && tryAcquire(arg)`，由于此时资源已被完全释放，所以一定能获取资源，而当该结点前有`CANCELLED`的线程时（`p == head`为false），在调整位置之后，必会成为head的后继结点，则在下一个循环中去执行对应的操作。

  ```java
  private void unparkSuccessor(Node node) {
          //尝试清空waitStatus，但失败或者被其它线程修改也没关系
      	//因为node多数情况下都为head(个人理解)
          int ws = node.waitStatus;
          if (ws < 0)
              node.compareAndSetWaitStatus(ws, 0);
          //唤醒等待队列中的下一个结点，当为空或被取消时（CANCELLED），则从后向前遍历
          Node s = node.next;
          if (s == null || s.waitStatus > 0) {
              s = null;
              //从后向前遍历，找到在当前结点后，且未被取消的结点
              for (Node p = tail; p != node && p != null; p = p.prev)
                  if (p.waitStatus <= 0)
                      s = p;
          }
          if (s != null)
              //唤醒对应的线程
              LockSupport.unpark(s.thread);
      }
  ```

  为什么要从后向前遍历？
  
  首先线程Node入队列的时候，先把自己的prev节点连到队列的尾节点，然后CAS设置自己为尾节点成功后，才将原来的尾结点的next指针指向自己，如果在这一步的时候，当前节点要被唤醒的话，由于next指针还没有指向当前节点，因此从头向尾遍历可能找不到当前节点，从尾部开始遍历就不存在这个问题了。

## 获取资源（共享模式）

共享模式下可以理解为资源有多个，当一个线程获取资源后，其它线程仍然有机会继续获取到资源

- acquireShared

  共享模式下获取资源的顶层入口，获取成功则直接返回，否则加入等待队列。

  ```java
  public final void acquireShared(int arg) {
      //仍然是需要自己实现的获取资源的方法
      //当获取失败时返回负数
      //当获取成功，但其它处于共享模式下的线程不会成功时返回0（没有资源了）
      //当获取成功，且其它牌共享模式下的线程可能会成功时返回正数（还有资源）
          if (tryAcquireShared(arg) < 0)
              doAcquireShared(arg);
  }
  ```

- doAcquireShared

  资源获取失败时，将线程加入阻塞队列队尾，和独占模式下有相似之处

  ```java
  private void doAcquireShared(int arg) {
      //加入阻塞队列的队尾
          final Node node = addWaiter(Node.SHARED);
          boolean interrupted = false;
          try {
              for (;;) {
                  final Node p = node.predecessor();
                  if (p == head) {
                      int r = tryAcquireShared(arg);
                      //大于等于0时代表获取资源成功
                      if (r >= 0) {
                          setHeadAndPropagate(node, r);
                          p.next = null; // help GC
                          return;
                      }
                  }
                  if (shouldParkAfterFailedAcquire(p, node))
                      interrupted |= parkAndCheckInterrupt();
              }
          } catch (Throwable t) {
              cancelAcquire(node);
              throw t;
          } finally {
              //try里面的代码return之后，会继续执行这里
              if (interrupted)
                  selfInterrupt();
          }
  }
  ```

- setHeadAndPropagate

  改变头结点为当前成功获取到资源的结点，并判断是否唤醒阻塞队列中的其它结点，唤醒操作是一个传播过程，一次只唤醒一个结点，但结点唤醒后会继续唤醒后继结点，直到某一个结点无法获取到资源。

  ```java
  private void setHeadAndPropagate(Node node, int propagate) {
          Node h = head; // Record old head for check below
          setHead(node);
          if (propagate > 0 || h == null || h.waitStatus < 0 ||
              (h = head) == null || h.waitStatus < 0) {
              Node s = node.next;
              if (s == null || s.isShared())
                  doReleaseShared();
          }
      }
  ```

- doReleaseShared

  ```java
  private void doReleaseShared() {
          for (;;) {
              Node h = head;
              if (h != null && h != tail) {
                  int ws = h.waitStatus;
                  if (ws == Node.SIGNAL) {
                      if (!h.compareAndSetWaitStatus(Node.SIGNAL, 0))
                          continue;            // loop to recheck cases
                      unparkSuccessor(h);
                  }
                  else if (ws == 0 &&
                           !h.compareAndSetWaitStatus(0, Node.PROPAGATE))
                      continue;                // loop on failed CAS
              }
              if (h == head)                   // loop if head changed
                  break;
          }
      }
  ```

## 释放资源（共享模式）

- releaseShared

  共享模式下释放资源的顶层入口，如果资源释放后（可以不完全释放）就去唤醒阻塞队列中等待的结点

  ```java
  public final boolean releaseShared(int arg) {
          if (tryReleaseShared(arg)) {
              //资源释放后，唤醒阻塞队列中等待的结点
              doReleaseShared();
              return true;
          }
          return false;
  }
  ```

    

