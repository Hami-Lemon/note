JAVA集合，多线程，网络编程，放射，注解
<!--more-->

## Collections

`static <T> void sort(List list)`对list进行排序，对于自定义的的类型，需要实现`Compareable<T>`接口并重写`compareTo()`（将`this`与参数进行比较）
`static <T> void sort(List list, Comparator com)`
`Comparable`:需要自己实现comparalble接口，并重写compareTo()方法
`Comparator`:比较器

```JAVA
Collections.sort(list, new Comparator<Integer>(){
    @Override
    public int compare(Integer o1, Integer o2){
        return o1 - o2;
    }
});
```

## Map集合

`Map<K,V>`
键值对

- HashMap<K, V>

1. 基于哈希表实现的map集合，查询速度快
2. 是一个无序的集合

- linkedHashMap<k, V>

1. 基于哈希表+链表实现
2. 是一个有序的集合，存储元素和去除元素的顺序一致

- 常用方法
  `V put(K key, V value)`添加，如果key重复，用value替换，如果不重复，返回null
  `V get(Object key)`获取
  `V remove(Object key)`删除
  `boolean containsKey(Object key)`判断键是否存在
  `Set<K> keySet()`获取所有的key，存在set集合中
  `Set<Map.Entry<K, V> entrySet()`把map集合的多个entry对象取出来，存储到一个set集合中
  ![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113916.png)
  - Map.Entry<k, v>
    Map中的内部接口
    当map集合一创建，就会创建一个entry对象，用于记录键与值
    ![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113921.png)
    `getKey()`获取键
    `getValue()`获取值
- Hashtable<k,v>
不能存储null，线程安全（同步）
只要使用其子类`Properties`
- 新特性
set， list，map:增加了一个静态方法of(),可以给集合一次性添加多个元素

1. 只适用于list，set，map，不适用于其实现类
2. 返回值是一个不可变集合，不能使用add，put方法
3. set，map接口不能有重复元素

## 异常

java.lang.throwable,所有异常的超类

- 处理异常
  - throws
    抛出异常，由虚拟机处理（程序会中断）
  - try-catch
    手动处理异常

    ```JAVA
    try{
    
    }catch(Exception e/*捕获的异常类*/){
    
    }
    ```

- 抛出异常
`throw new xxxException("内容")`
- Throwable异常处理方法
`String getMessage()`异常的简短描述
`String toString()`详细信息
- finally
无论是否有异常，一定会执行finally中的代码
- 注意

1. 一个try多个catch时，如果异常有子父类关系，子类必须写在父类上面
2. 如果finally里中有return，永远是返回finally中的return
3. 父类方法没有抛出异常时，子类重写时不能抛出，只能捕获
4. 一般，子父类异常相同

- 自定义异常

1. 一般以Exception结尾
2. 必须继承Exception(编译器异常)或者RuntimeException(运行期异常,一般不需要处理)

```JAVA
public class xxxException extends Exception{
    public xxxException(){

    }
    public xxxException(String msg){
        super(msg);
    }
}
```

## 多线程

- 主线程
执行主（main）方法的线程，主线程的结束意味着程序的结束
- 开辟子线程

```JAVA
//方式一
public class Demo extends Threads{
    //重写run()方法
}
new Demo().start();//启动线程
//方式二
class Demo implements Runable{
    //重写run()方法
}
new Thread(new Demo()).start;//开启线程
//匿名形式

//方式一
new Thread(){
    public void run(){

    }
}.start();
//方式二
new Thread(new Runable(){
    public void run(){

    }
}).start;
```

- 内存图
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113928.png)
- thread常用方法
`String getName()`获取线程名称
`static Thead currentThread()`获取当前执行的线程对象
`void setName(String name)`设置线程名称
`static void sleep(long millis)`使当前线程睡觉
- 线程安全
多个线程访问共享的数据时，会出现数据更新错误的问题

```JAVA
//同步代码块
synchronized(锁对象){
    //访问共享数据的代码
}
/*1.锁对象可以是任意东
2. 但是必须保证多个线程的锁对象是同一个
*/

//同步方法
public synchronized void method(){
    //访问共享数据的代码
    //锁对象是this
}

//静态同步方法
public static synchronized void method(){
    //同上
    //锁对象是本类的class文件对象
}

//锁机制
//java.util.concurrent.locsk.Lock
//使用其实现类ReentrantLock
Look l = new ReentrantLock();

l.look();//获取锁
//访问共享数据的代码
l.unlook();//释放锁，有异常时放在finally中
```

- 生命周期
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113933.png)
- 线程等待与唤醒
`void wait()`线程等待，无限等待（Object中）
`void notify()`唤醒线程
使用锁对象调用上面两个方法
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113938.png)
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113946.png)
- 线程池
  实现线程复用
  `java.util.Concurrent.Executors`生产线程池的工厂类
  `static ExecutorService newFixedThreadPool(int nThreads)`参数：nThread 线程池中的线程数

  - 使用

    ```JAVA
    ExecutorService es = Executors.newFixedThreadPool(3);
    es.submit(new Runnable()J{
        public void run(){
            //
        }
    });//调用线程池中的一个线程
    es.shutdown();//结束线程池
    ```

## 函数式编程

Lambda表达式，用于在重写接口函数时简化代码

```JAVA
public class Demo{
    public static void main(Strings args[]){
        //Lambda表达式，实现多线程
        new Thread(() -> {
            //run（）函数的方法体
        }).start();
    }
}
```

- 格式
`(参数列表) -> {一些重写方法的代码}`
()：接口中抽象方法的参数列表
-> :
{...}: 重写抽象方法的接口
- 无参无返回值
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113952.png)
- 有参数有返回值
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113957.png)
![](/images/JAVA/l3.png)
lanbda中凡是根据上下文推到出来的即可省略
例：参数列表中，类型可以省略
如果参数只有一个，类型和小括号（只是小括号）都可以省略
{}：中的代码只有一行，无论是否有返回值，可以省略{}，return，分号(三者需一起省略)
- 使用前提

1. 必须是接口，并且接口中有且只有一个抽象方法，这种接口被称为函数式接口
2. 必须具有上下文推断，也就是方法的参数混局部变量类型必须为lambd对应的接口类型
3. 一个函数的参数或者一个局部变量的类型是一个接口，并且该接口满足条件1，可以使用lambda，简化使用匿名内部类的语法

## File类

- 路径
    1. 绝对路径
    以盘符开始的路径
    2. 相对路径
    相对于当前项目的根目录
  - 注意
    路径不区分大小写
- 常用方法
  - 获取功能方法
    `String getAbsolutePath()`返回文件的绝对路径
    `String getPath()`将file对象转化为路径名称支付穿
    `String getName()`放回当前文件或目录的名称
    `long length()`放回当前文件的字节大小，如果文件不存在则返回0
  - 判断功能方法
    `boolean exists()`判断文件是否存在
    `boolean isDirectory()`判断是否是目录
    `boolean isFile()`判断是否是文件
  - 创建删除功能方法
    `boolean createNewFile()`创建新的文件，只能是文件，且路径必须存在
    `boolean mkdir()`创建单级目录
    `boolean delete()`删除File对象表示的文件或目录，目录中有内容时，不会删除
    `boolean mkdirs()`创建多级目录
  - 目录遍历
    `String[] list()`
    `File[] listFiles()`
- 文件过滤器
  `java.io.FilFilter`接口：用于File对象的过滤器
  抽象方法：`boolean accept(File pathname)`测试指定File对象是否应该包好在某个路径名列表中
  `java.io.FilenameFilter`接口：实现文件名称的过滤器
  抽象方法：`boolean accept(File dir, String name)`测试指定文件是否应该包含在某一文件列表中
  - 使用
    listFiles方法会将遍历每一个文件或目录，并调用accept()方法，将每一个文件或目录作为参数传递给accept()方法,如果accept()返回true，则将文件对象存在返回的数组中

    ````JAVA
    File f = new File("...");
    File[] files = f.listFiles(pathname -> {
        //过滤规则
    })
    ```

## IO

- 字节输出流
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114003.png)

```JAVA
//创建Fileoutputstream对象
FileOutputStream fos = new FileOutputStream("...");
//调用write()写入数据
fos.write(...);
//释放资源
fos.close();
```

- 字节输入流

```JAVA
//创建FileInputStream对象
FileInputStream fis = new FileInputStream("...");
//使用read()读取
fis.read();//读取一个字节
//释放资源
fis.close();
```

- 关闭与刷新
`fluse()`刷新缓冲区，流对象可以继续使用
`close()`先刷新缓冲区，然后释放资源，流对象不能继续使用
- 异常处理
try前面可以定义流对象，后面可以加()，并且可以映入流对象的名称，在try代码执行完后，会自动释放流对象，无需finally

```JAVA
//方式一
FileInputStream fis = new FileInputStream("...");
FileOutputStream fos = new FileOutputStream("...");

try(fis;fos){
    //有异常的代码
}catch(IOException e){
    //处理代码
}

//方式二
try(
    FileInputStream fis = new FileInputStream("...");
    FileOutputStream fos = new FileOutputStream("...");){
    //有异常的代码
}catch(IOException e){
    //处理代码
}
```

## Properties(属性集)

可以使用properties集合中的方法store，把集合中的临时数据，持久化写入到硬盘中
使用load，把保存的文件(键值对)读取到集合中

```JAVA
Proerties prop = new Proerties();
prop.setProperty(key, value);//添加数据
Set<String> set = prop.stringPropertyNames();//获取所有键
prop.getProperty(key);//通过键获取值

FileWriter fw = new FileWriter("...");
prop.store(fw, "注释，最好不要用中文");//保存集合中的数据到文件

FileReader fr = new FileReader("...");
prop.load(fr);
```

- 缓冲流
`BufferedInputStream`
`BufferedOutputStream`
`BufferedReader`
`BufferedWriter`
- 转换流
可以指定编码表
`InputStreamReader`字节转字符
`OutputStreamWriter`字符转字节
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114009.png)

```JAVA
OutputStreamWriter osw = new OutputStreamWriter(new FileOutputStream("path"), "utf-8");//使用utf-8编码表
osw.write("...");
osw.flush();//刷新缓冲区，写入内容到文件
osw.close();//释放资源
```

```JAVA
InputStreamReader isr = new InputStreamReader(new FileInputStream("path"), "utf-8");
isr.read();//读取
isr.close();
```

- 序列化与反序列化
序列化：把一个对象的信息写入到文件中
反序列化：将一个对象的信息从文件中读取出来
需要被序列化的类实现Serializable接口，只是一个标记，不需要时实现方法
被static修饰的成员不能被序列化
transient(瞬态)：被transient修饰的成员变量不能被序列化，仅仅是让成员不被序列化

```JAVA
//序列化
ObjectOutputStream oos = new ObjectOutputStream(new FileOutputStream("path"));
oos.writeObject(obj);
oos.close();

//反序列化
ObjectInputStream ois = new ObjectInputStream(new FileInputStream(""));
Object obj = ois.readObject();
ois.close();
```

    - 自定义序列化编号serialVsersionUIOD
    `private static finall long serialVersionUID = ...;`

- 打印流
`java.io.PrintStream`
可以使用System.setOut(PrintStream ps)改变输出语句的流向

## 网络编程

在网络中，连接和通信的规则

- TCP/IP协议(传输控制协议/因特网互联协议)
最基本最广泛的协议
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114014.png)
- 通信协议分类
  - UDP:无连接通信协议，效率高，耗资小，不能保证数据完整型
  - TCP:传输控制协议，面向连接的通信协议，在传输前先三次握手
- 三要素
    1. 协议
    2. IP地址
    3. 端口号
- TCP通信
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114018.png)

```JAVA
//Client
class Client{
    //绑定服务器socket
   Socket socket = new Socket("127.0.0.1", 8888);
   //通过socket获取网络输出流
   OutputStream outputStream = socket.getOutpuStream();
   //通过流传输数据
    os.write("hello".getBytes());
    //读取服务端发送的数据
    InputStream is = socket.getInputStream();
    byte[] bytes = new byte[1024];
    int len = is.read(bytes);
    //释放资源
    socket.close();
}

//Server
class Server{
    //创建服务器socket
    ServerSocket server = new ServerSocket(8888);
    //监听
    Socket socket = server.accept();
    //获取输入流
    InputStream is = socket.getInputStream();
    //读取数据
    byte[] bytes = new byte[1024];
    int len = is.read(bytes);
    //向客户端发送数据
    OutputStream os = socket.getOutputStream();
    os.write("over".getBytes());
    //释放资源
    socket.close();
    server.close();
}
```

- 阻塞问题
结束标记不会被上传，所以服务器端永远不会读取到结束标记，被read()方法阻塞
解决：客户端主动上传一个结束标记，调用Socket中的shutdownOutput()方法亦或者自己定义一个结束标记，由服务端判断

- B/S 版TCP

```JAVA
class Server{
    ServerSocket server = new ServerSocket(8080);
    Socket socket = server.accept();
    InputStream is = Socket.getInputStream();
    byte[] bytes = new byte[1024];
    int len = 0;
    while((len = is.read(bytes)) != -1){
        //...
    }
}
```

- UDP通信

```JAVA
//server
DatagramSocket socket = new DatagramSocket(8888);
DatagramPacket packet = new DatagramPacket(new byte[1024], 1024);

socket.receive(packet);//接收数据
byte[] bytes = packet.getData();//获取数据
int len = packet.getLength();
System.out.println(new String(bytes, 0, len));
socket.close();

//client
String info = "info";
DatagramSocket socket = new DatagramSocket();
DatagramPacket packet = new DatagramPacket(info.getBytes(),info.length,InetAddress.getByName("127.0.0.1"), 8888);

socket.send(packet);//发送数据
socket.close();
```

## 函数式接口

有且只有一个抽象方法的接口，称为函数式接口，当然也可以有其它的非抽象方法,可以使用lambda表达式
lambda可以作为方法的参数或者返回值

```JAVA
//demo为函数式接口
//作为放回值
public static void main(String[] args) {
        System.out.println(d().method());
    }
private static demo d(){
        return () -> {return 10 + 10;};
    }
}
```

- 常用的函数式接口
`Supplier<T>{ T get();}`生产型接口
`Consumer<T>{void accept(T t);}`消费型接口
    默认方法，andThen(Consumer c):需要两个Consumer接口，可以把两个接口组合到一起
    `con1.accept(s);con2.accept(s);`
    等价于`con1.andThen(con2).accept(s);`
`Predicate<T> {boolean test(T t);}`对数据进行某种判断

##  Stream流

(不同于IO流)，对集合和数组的元素进行操作

```JAVA
//获取流
//方式一
//通过Collection集合中的默认方法stream();
List<String> list = new ArrayList<>();
Stream<String> stream = list.stream();
//方式二
//通过Stream接口中静态方法of(),参数是可变参数
String[] arr  = {...};
Stream<String> stream = Stream.of(arr);
```

- 常用方法
`void forEach(Consumer)`
用来遍历中的数据，是一个终结方法，遍历之后就不能继续使用流
`stream.forEach(() -> {/*...*/})`
`Stream filter(Predicate)`
筛选出流中符合条件的数据
`stream.filter(() -> {...});`
`Stream map(Function)`
映射，将流中的数据按照指定的规则进行转换

```JAVA
Stream<Integer> stream2 = stream1.map((String s) -> {
    return Integer.parseInt(s);//将字符串转为整数
});
```

`long count();`统计流的元素个数，终结方法
`Stream limit(long maxSize)`对流进行截取，取流中前maxSize个元素
`Stream skip(long n)`跳过流中前n个元素，如果总个数小于n，则返回一个空的流
`Stream concat(Stream stream1, Stream stream2)`将两个流组合到一起
stream属于管道流，只能被消费一次，第一个stream调用一个方法后，数据流向下一个Stream流，前一个就已经关闭

## 方法引用

(对lambda表达式的优化，当方法体只是对其它方法的调用时使用)
lamdba写法：`s -> System.out.println(s);`
方法引用写法：`System.out::println`因为其中syste.out对象已经存在，并且println方法也存在
调用函数的参数列表与lambda参数列表相同

- 通过对象名引用成员方法

```JAVA
MethodObject obj = new MethodObject();
obj::method;
```

- 通过类名引用静态成员方法
`MethodObject::method;`method是MethodObject中的静态方法
- 通过super引用成员方法
存在继承关系时，用super引用父类方法
`super::method`method为父类的方法
- 通过this引用成员方法
`this::method`method为当前类中的方法
- 类构造器的引用
`MethodObject::new`创建一个MethodObject对象，用于构造一个对象，自动根据参数选择构造方法
- 数组构造器的引用
`int[]::new`

## Junit

单元测试
步骤：

1. 定义一个测试类(测试用例)(类名建议为原类名+Test，包名：xxx.xxx.xxx.test)
2. 定义测试方法，可以独立运行(方法名test+测试的方法名，放回值建议使用void，参数列表建议空参)
3. 给方法加注解@Test
4. 导入junit依赖

- 断言`Assert`
使用断言判断结果是否是期望的值
`@Before`在所有测试方法前自动执行的方法
`@After`在所有方法结束时自动执行

## 反射：框架设计的灵魂

框架：半成品软件，可以在框架的基础上进行软件开发，简化编码
反射：将类的各个组成部分封装成其它对象
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114026.png)

- 优点

1. 可以在程序运行中操作这些对象
2. 可以解耦，提高程序的可扩展性

- 获取Classs类对象的方式

1. `Class.forName("全类名");`将字节码文件加载进内存，返回class对象
2. `类名.class;`通过类名的class属性
3. `对象.getClass();`Object类中的方法

```JAVA
//方式一
Class cls = Class.forname("....");//需要穿全类名(带上包名)
//方式二
Class cls = Person.class;
//方式三
Person p = new Person();
Class cls = p.getClass();
```

同一个字节码(*.class)文件在程序运行过程中，只会被加载一次，所以使用三种方式获取的class对象是同一个

- 使用class对象
  * 获取成员变量
    + `Field[] getFields();`获取所有public修饰的成员变量
    + `Field getField(String name);`获取指定名称的public变量
    + `Field[] getDeclaredFields();`获取所有成员变量，不考虑修饰符
    + `Field getDeclaredField(String name);`
  * 获取构造方法
    + `Constructor<?>[] getConstructors();`获取构造方法
    + `Constructor<T> getConstructor(Class<?>... parameterTypes)`获取指定参数列表的构造方法
            `Constructor construcotr = personClass.getConstructor(String.class, int.class);`
    + `Constructor<T> getDeclaredConstructor(Class<?>...parameterTypes);`
    + `Constructor<?>[] getDeclaredConstructors();`
  * 获取成员方法
    + `Method[] getMethods();`
    + `Method getMethod(String name, Class<?>... parameterTypes);`
    + `Method[] getDeclaredMethods();`
    + `Method getDeclaredMethod(String name, Class<?>... parameterTypes);`
  * 获取类名
    + `String getName();`
  * Field对象操作
    + 设置值
        `void set(Object obj, Object value);`给obj设置值
    + 获取值
        `Object get(Object obj);`获取obj的值
  * Constructor对象操作
    + 创建对象
        `Object newInstance(Object... initargs);`

  * Method对象操作
    + 执行方法
        `Object invoke(Object obj, Class<?>... args)`需要传递当前方法的所属对象
    + 获取方法名称
        `String getName()`

## 注解

一种代码级别的说明，与类，接口三，枚举同一层次，用来对包，类，字段，方法，局部变量，方法参数进行说明

+ 分类
  - 编写文档
  - 代码分析
  - 编译检查
+　预定义注解

```JAVA
@Override：该方法是否是继承父类
@Deprecated：该注解标注的内容，表示为已过时
@SuppressWarnings：压制警告(忽略警告)`@SuppressWarnings("all")`
@SafeVarargs:在使用可变数量的参数是,而参数类型是泛型,使用此注解去掉警告,必须是static或final的方法或构造方法
@FunctionalInterfacea:函数式接口,只有一个抽象方法的接口
```

+ 自定义注解
格式：

```JAVA

    @Target({METHOD, TYPE})//作用在什么(类,接口,方法...)上面
    @Retention(RetentionPolicy.RUNTIME)//生命周期
    @Inherited//可被子类继承
    @Documented//会生成相关文档
    public @interface JDBCConfig{
        String ip();
        String database();
        String encoding();
        int port() default 3306;
    }


//使用
@JDBCConfig(ip='',database=''...)
public class DBUtil{
    ...
}
```

本质上是一个继承了java.lang.annotation.Annotation的接口

+ 解析注解
使用反射获取

```JAVA
//从注解使用对象上获取
JDBCConfig config = DBUtil.class.getAnnotation(JDBCConfig.class);
String ip = config.ip();
int port = config.port();
...

```

+ 元注解
@Target:放在什么位置
  - ElementType.Type:类,接口或枚举
  - ElementType.FIELD:成员变量
  - ElementType.METHOD:方法
  - ElementType.PARAMETER:修饰参数
  - ElementType.CONSTRUCTOR:构造器
  - ElementType.LOCAL_VARIABLE:局部变量
  - ElementType.ANNOTATION_TYPE:修饰注解
  - ElementType.PACKAGE:修饰包
@Retention:生命周期
  - RetentionPolicy.SOURCE:只在源代码中存在
  - RetentionPolicy.CLASS:存在于.class中,运行中不存在
  - RetentionPolicy.RUNTIME:运行中仍然存在
@Inherited:具有继承性,子类会自动继承父类上的此注解
@Documented:文档中会出现此注解
@Repeatable:可以重复出现(1.8新增)

