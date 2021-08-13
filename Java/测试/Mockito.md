# Mockito

## Mock测试

Mock测试是指在测试过程中，对于某些不易构建的对象通过一个虚拟对象来创建，以便于进行测试。同时还能降降低测试时代码的耦合度，例如：A类依赖于某些类，当测试A类时，为了能使A类能正常工作，还需要构建出A类所需要的整个依赖树，当使用Mock测试时，相关依赖可以完全通过Mock对象提供而不依赖于具体的类。

## Maven依赖

```xml
<dependency>
    <groupId>org.mockito</groupId>
    <artifactId>mockito-core</artifactId>
    <version>3.6.28</version>
    <scope>test</scope>
</dependency>
```

## 简单使用

```java
//静态导入所有静态方法
import static org.mockito.Mockito.*;
//创建Mock对象
List mockedList = mock(List.class);
//使用Mock对象
mockedList.add("one");
mockedList.clear();
//验证Mock对对象的add方法和clear方法是否被调用
verify(mockedList).add("one");
verify(mockedList).clear();
```

## 打桩

通过对Mock对象中的方法打桩，使其在被调用时返回指定的值。

```java
 //创建Mock对象
 LinkedList mockedList = mock(LinkedList.class);

 //打桩
 when(mockedList.get(0)).thenReturn("first");
 when(mockedList.get(1)).thenThrow(new RuntimeException());

 //这会打印 first
 System.out.println(mockedList.get(0));

 //这会抛出 RuntimeException
 System.out.println(mockedList.get(1));

 //这会打印 null ，因为并没有对 get(999)方法进行打桩
 System.out.println(mockedList.get(999));
//验证 get(0)是否被调用
 verify(mockedList).get(0);
```

## 参数匹配

在Mockito中，同一个方法，如果传入参数不同（默认通过`equals`方法来比较是否相同），也会被认为是不同的操作，例如：`get(0)`方法和`get(999)`并不会被认为是相同的操作，但Mockito也提供了相关方法来对参数进行“通配”。

```java
 //匹配所有参数为int类型
 when(mockedList.get(anyInt())).thenReturn("element");

 //自定义匹配方法
 when(mockedList.contains(argThat(isValid()))).thenReturn(true);

 //following prints "element"
 System.out.println(mockedList.get(999));

 //同样可以将参数匹配也于验证
 verify(mockedList).get(anyInt());

 //argument matchers can also be written as Java 8 Lambdas
 verify(mockedList).add(argThat(someString -> someString.length() > 5));
```

## 验证方法调用次数

```java
 mockedList.add("once");
 mockedList.add("twice");
 mockedList.add("twice");
 mockedList.add("three times");
 mockedList.add("three times");
 mockedList.add("three times");

 //默认验证是是否被调用一次
 verify(mockedList).add("once");
 verify(mockedList, times(1)).add("once");

 //通过 times方法验证方法被调用的次数
 verify(mockedList, times(2)).add("twice");
 verify(mockedList, times(3)).add("three times");

 //验证从未被调用
 verify(mockedList, never()).add("never happened");

 //至多，至少调用几次
 verify(mockedList, atMostOnce()).add("once");
 verify(mockedList, atLeastOnce()).add("three times");
 verify(mockedList, atLeast(2)).add("three times");
 verify(mockedList, atMost(5)).add("three times");
```

## 验证调用顺序

```java
 List singleMock = mock(List.class);
 singleMock.add("was added first");
 singleMock.add("was added second");

 //创建顺序验证器
 InOrder inOrder = inOrder(singleMock);
//验证方法的调用顺序
 inOrder.verify(singleMock).add("was added first");
 inOrder.verify(singleMock).add("was added second");

 // 多个不同的对象也可以验证
 List firstMock = mock(List.class);
 List secondMock = mock(List.class);
 firstMock.add("was called first");
 secondMock.add("was called second");

 InOrder inOrder = inOrder(firstMock, secondMock);
 inOrder.verify(firstMock).add("was called first");
 inOrder.verify(secondMock).add("was called second");
```

## 验证Mock对象是否被使用

```java
 //只有mockOne被使用
 mockOne.add("one");

 verify(mockOne).add("one");

 verify(mockOne, never()).add("two");

 //验证mockTwo和mockThree从未被使用
 verifyZeroInteractions(mockTwo, mockThree);
```

## 连续打桩

```java

 when(mock.someMethod("some arg"))
   .thenThrow(new RuntimeException())
   .thenReturn("foo");
 //第一次调用时会抛出 RuntimeException异常
 mock.someMethod("some arg");

 //第二次调用时则返回 foo
 System.out.println(mock.someMethod("some arg"));

 //此后都会返回 foo
 System.out.println(mock.someMethod("some arg"));
```

## 通过回调方法打桩

```java
 when(mock.someMethod(anyString())).thenAnswer(
     new Answer() {
         public Object answer(InvocationOnMock invocation) {
             Object[] args = invocation.getArguments();
             Object mock = invocation.getMock();
             return "called with arguments: " + Arrays.toString(args);
         }
 });

 //Following prints "called with arguments: [foo]"
 System.out.println(mock.someMethod("foo"));
```

## 对Void方法打桩

对于Void方法，对其打桩的需要在调用`when`方法前调用`doReturn`，`doThrow`，`doAnswer`,`doNothing`或`doCallRealMethod`方法来指定其行为。

```java
   doThrow(new RuntimeException()).when(mockedList).clear();
   //抛出RuntimeException异常
   mockedList.clear();
```

## 使用Spy

通过Spy可以不仅可以对Mock对象打桩，还可以直接调用Mock对象的真实方法

```java
   List list = new LinkedList();
   List spy = spy(list);
   //optionally, you can stub out some methods:
   when(spy.size()).thenReturn(100);
   //using the spy calls *real* methods
   spy.add("one");
   spy.add("two");
   //打印 one
   System.out.println(spy.get(0));
   //打印 100 ，因为这个方法被打桩
   System.out.println(spy.size());
   //optionally, you can verify
   verify(spy).add("one");
   verify(spy).add("two");
```

