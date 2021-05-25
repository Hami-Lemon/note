# Netty

## NIO

NIO是一个非阻塞的,面向缓冲的处理
NIO主要有三个核心部分

- buffer缓冲区
- channel管道
- selector选择器

### buffer和channel

NIO使用buffer和channel一起对数据进行处理,channel是一个可以读写的双向通道

- Buffer缓冲区

```java
ByteBuffer buffer = ByteBuffer.allocate(1024);//初始化缓冲区,大小为1024
```

buffer四个属性:

- capacity:容量,代表缓冲区的最大容量
- limit:上界,表示缓冲区可以操作的数据大小(limit <= capacity)
- position:位置,表示下一个被操作元素的位置
- mark:标记,用于记录当前position的位置,可以通过reset()恢复

常用方法:

- file():limit设为position,将position设为0,在读取数据前应当调用此方法
- clear():"清空"缓冲区数据,并不是真正意义的清空,只是"遗忘"数据

### Channel通道

- FileChannel:操作文件
- SocketChannel:tcp连接,客户端
- ServerSocketChannel:tcp连接,服务端
- DatagramChannel:udp连接

channel只负责传输数据,不直接操作,操作数据都是通过buffer进行

```java
// 1. 通过本地IO的方式来获取通道
FileInputStream fileInputStream = new FileInputStream("F:\\3yBlog\\JavaEE常用框架\\Elasticsearch就是这么简单.md");
// 得到文件的输入通道
FileChannel inchannel = fileInputStream.getChannel();

// 2. jdk1.7后通过静态方法.open()获取通道
FileChannel.open(Paths.get("F:\\3yBlog\\JavaEE常用框架\\Elasticsearch就是这么简单2.md"), StandardOpenOption.WRITE);

```

### IO模型

根据UNIX网络编程对I/O模型的分类

- 阻塞I/O
- 非阻塞I/O
- 多路复用I/O
- 信号驱动I/O
- 异步I/O

#### 文件描述符

Linux 的内核将所有外部设备都看做一个文件来操作，对一个文件的读写操作会调用内核提供的系统命令(api)，返回一个file descriptor（fd，文件描述符）。而对一个socket的读写也会有响应的描述符，称为socket fd（socket文件描述符），描述符就是一个数字，指向内核中的一个结构体（文件路径，数据区等一些属性)  
所以说：在Linux下对文件的操作是利用文件描述符(file descriptor)来实现的。

#### I/O运行过程

![过程](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525112935.png)

#### 阻塞I/O

![阻塞](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525112938.png)

#### 非阻塞I/O

![非阻塞](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525112941.png)

#### I/O多路复用

在linux中,调用`select/poll/epoll/pselect`其中一个函数,会**传入多个文件描述符**,如果有一个文件描述符**就绪**,则返回,否则阻塞直到超时  
![多路复用](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525112945.png)

```txt
Java3y跟女朋友去麦当劳吃汉堡包，现在就厉害了可以使用微信小程序点餐了。
于是跟女朋友找了个地方坐下就用小程序点餐了。点餐了之后玩玩斗地主、聊聊天什么的。
时不时听到广播在复述XXX请取餐，反正我的单号还没到，就继续玩呗。
~~等听到广播的时候再取餐就是了。时间过得挺快的，此时传来：Java3y请过来取餐。
于是我就能拿到我的麦辣鸡翅汉堡了。听广播取餐，广播不是为我一个人服务。
广播喊到我了，我过去取就Ok了。
```

### NIO进行网络通信

在网络中使用的NIO往往时I/O模型的**多路复用模型**

- selector选择器可以比喻成麦当劳的广播
- 一个线程能够管理多个Channel的状态
![网络通信](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525112954.png)

#### 阻塞的NIO

客户端

```java
public class BlockClient {

    public static void main(String[] args) throws IOException {

        // 1. 获取通道
        SocketChannel socketChannel = SocketChannel.open(new InetSocketAddress("127.0.0.1", 6666));

        // 2. 发送一张图片给服务端吧
        FileChannel fileChannel = FileChannel.open(Paths.get("X:\\Users\\ozc\\Desktop\\新建文件夹\\1.png"), StandardOpenOption.READ);

        // 3.要使用NIO，有了Channel，就必然要有Buffer，Buffer是与数据打交道的呢
        ByteBuffer buffer = ByteBuffer.allocate(1024);

        // 4.读取本地文件(图片)，发送到服务器
        while (fileChannel.read(buffer) != -1) {

            // 在读之前都要切换成读模式
            buffer.flip();

            socketChannel.write(buffer);

            // 读完切换成写模式，能让管道继续读取文件的数据
            buffer.clear();
        }

        // 5. 关闭流
        fileChannel.close();
        socketChannel.close();
    }
}
```

服务端

```java
public class BlockServer {

    public static void main(String[] args) throws IOException {

        // 1.获取通道
        ServerSocketChannel server = ServerSocketChannel.open();

        // 2.得到文件通道，将客户端传递过来的图片写到本地项目下(写模式、没有则创建)
        FileChannel outChannel = FileChannel.open(Paths.get("2.png"), StandardOpenOption.WRITE, StandardOpenOption.CREATE);

        // 3. 绑定链接
        server.bind(new InetSocketAddress(6666));

        // 4. 获取客户端的连接(阻塞的)
        SocketChannel client = server.accept();

        // 5. 要使用NIO，有了Channel，就必然要有Buffer，Buffer是与数据打交道的呢
        ByteBuffer buffer = ByteBuffer.allocate(1024);

        // 6.将客户端传递过来的图片保存在本地中
        while (client.read(buffer) != -1) {

            // 在读之前都要切换成读模式
            buffer.flip();

            outChannel.write(buffer);

            // 读完切换成写模式，能让管道继续读取文件的数据
            buffer.clear();

        }

        // 7.关闭通道
        outChannel.close();
        client.close();
        server.close();
    }
}
```

#### 非阻塞NIO

客户端

```java
public class NoBlockClient {

    public static void main(String[] args) throws IOException {

        // 1. 获取通道
        SocketChannel socketChannel = SocketChannel.open(new InetSocketAddress("127.0.0.1", 6666));

        // 1.1切换成非阻塞模式
        socketChannel.configureBlocking(false);

        // 2. 发送一张图片给服务端吧
        FileChannel fileChannel = FileChannel.open(Paths.get("X:\\Users\\ozc\\Desktop\\新建文件夹\\1.png"), StandardOpenOption.READ);

        // 3.要使用NIO，有了Channel，就必然要有Buffer，Buffer是与数据打交道的呢
        ByteBuffer buffer = ByteBuffer.allocate(1024);

        // 4.读取本地文件(图片)，发送到服务器
        while (fileChannel.read(buffer) != -1) {

            // 在读之前都要切换成读模式
            buffer.flip();

            socketChannel.write(buffer);

            // 读完切换成写模式，能让管道继续读取文件的数据
            buffer.clear();
        }

        // 5. 关闭流
        fileChannel.close();
        socketChannel.close();
    }
}
```

服务端

```java
public class NoBlockServer {

    public static void main(String[] args) throws IOException {

        // 1.获取通道
        ServerSocketChannel server = ServerSocketChannel.open();

        // 2.切换成非阻塞模式
        server.configureBlocking(false);

        // 3. 绑定连接
        server.socket().bind(new InetSocketAddress(6666));

        // 4. 获取选择器
        Selector selector = Selector.open();

        // 4.1将通道注册到选择器上，指定接收“监听通道”事件
        server.register(selector, SelectionKey.OP_ACCEPT);

        // 5. 轮训地获取选择器上已“就绪”的事件--->只要select()>0，说明已就绪
        while (selector.select() > 0) {
            // 6. 获取当前选择器所有注册的“选择键”(已就绪的监听事件)
            Iterator<SelectionKey> iterator = selector.selectedKeys().iterator();

            // 7. 获取已“就绪”的事件，(不同的事件做不同的事)
            while (iterator.hasNext()) {

                SelectionKey selectionKey = iterator.next();

                // 接收事件就绪
                if (selectionKey.isAcceptable()) {

                    // 8. 获取客户端的链接
                    SocketChannel client = server.accept();

                    // 8.1 切换成非阻塞状态
                    client.configureBlocking(false);

                    // 8.2 注册到选择器上-->拿到客户端的连接为了读取通道的数据(监听读就绪事件)
                    client.register(selector, SelectionKey.OP_READ);

                } else if (selectionKey.isReadable()) { // 读事件就绪

                    // 9. 获取当前选择器读就绪状态的通道
                    SocketChannel client = (SocketChannel) selectionKey.channel();

                    // 9.1读取数据
                    ByteBuffer buffer = ByteBuffer.allocate(1024);

                    // 9.2得到文件通道，将客户端传递过来的图片写到本地项目下(写模式、没有则创建)
                    FileChannel outChannel = FileChannel.open(Paths.get("2.png"), StandardOpenOption.WRITE, StandardOpenOption.CREATE);

                    while (client.read(buffer) > 0) {
                        // 在读之前都要切换成读模式
                        buffer.flip();

                        outChannel.write(buffer);

                        // 读完切换成写模式，能让管道继续读取文件的数据
                        buffer.clear();
                    }
                }
                // 10. 取消选择键(已经处理过的事件，就应该取消掉了)
                iterator.remove();
            }
        }

    }
}
```

## NIO与零拷贝

在java中，常用的零拷贝有mmap(内存映射) 和sendFile

- 内存映射(mmap)
![mmap](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113002.png)
- sendFile
![sendFile](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113005.png)

- NIO中的零拷贝

  ```java
  //将当前channel的数据发送到target中,返回实际发送的数据
  long transferTo(long position, long count,WritableByteChannel target)
  //从src中获取数据到当前channel中,返回实际获取的数据
  long transferFrom(ReadableByteChannel src,long position, long count)
  ```

只有`FileChannel`的子类可以调用这两个方法
在windows下,`transferTo`一次只能发送8兆的数据,linux下无影响

## Netty框架

netty是一个异步，基于事件驱动的网络应用框架
![netty](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113009.png)

netty工作流程
![netty](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113013.png)
maven依赖

```xml
<dependency>
    <groupId>io.netty</groupId>
    <artifactId>netty-all</artifactId>
    <version>4.1.49.Final</version>
</dependency>
```

- Reactor模式
![reacotr](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113016.png)

### Netty创建TCP服务

![netty](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113019.png)

服务端

```java
public class Server {
    public static void main(String[] args) throws InterruptedException {
        //创建BossGroup和workerGroup
        //bossGroup只是处理连接请求
        //和客户端的业务处理交给workerGroup，两个都是无限循环
        //bossGroup 和workerGroup的子线程个数为cpu核心数 * 2
        NioEventLoopGroup bossGroup = new NioEventLoopGroup(1000);
        NioEventLoopGroup workerGroup = new NioEventLoopGroup();
        try{
            ServerBootstrap bootstrap = new ServerBootstrap();//创建服务器端的启动对象，配置 参数
            bootstrap.group(bossGroup, workerGroup) //设置两个线程组
                    .channel(NioServerSocketChannel.class)//使用NioSocketChannel作为服务器的通道
                    .option(ChannelOption.SO_BACKLOG, 128)//设置线程队列最大连接的个数
                    .childOption(ChannelOption.SO_KEEPALIVE, true)//设置保持活动连接状态
                    .childHandler(new ChannelInitializer<SocketChannel>() {
                        //给pipeline设置处理器
                        @Override
                        protected void initChannel(SocketChannel ch) throws Exception {
                            ch.pipeline().addLast(new ServerHandler());
                        }
                    });//给workerGroup的对应管道设置处理器
            System.out.println("服务器设置完成");
            //绑定端口，并同步，会生成一个ChannelFuture对象
            ChannelFuture cf = bootstrap.bind(8888).sync();
            //对关闭通道进行监听
            cf.channel().closeFuture().sync();
        }finally {
            //关闭
            bossGroup.shutdownGracefully();
            workerGroup.shutdownGracefully();
        }
    }
}
```

服务端处理器

```java
public class ServerHandler extends ChannelInboundHandlerAdapter {
    /**
     * 在这里读取客户端发送的消息
     * @param ctx 上下文对象，含有通道pipeline，管道
     * @param msg 客户端发送的数据
     * @throws Exception 异常
     */
    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        //将msg转成ByteBuf,netty提供的
        ByteBuf buf = (ByteBuf) msg;
        System.out.println("获取到消息" + buf.toString(CharsetUtil.UTF_8));
        System.out.println("客户端地址" + ctx.channel().remoteAddress());
    }
    /**
     * 数据读取完毕，回写数据
     * @param ctx
     * @throws Exception
     */
    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) throws Exception {
        //将数据写入到缓存，并刷新(向客户端会写数据)
        ctx.writeAndFlush
                (Unpooled.copiedBuffer("hello", CharsetUtil.UTF_8));
    }
    /**
     * 处理异常
     * @param ctx
     * @param cause
     * @throws Exception
     */
    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        //关闭通道
        ctx.close();
    }
}
```

客户端

```java
public class Client {
    public static void main(String[] args) throws InterruptedException {
        //客户端需要一个事件循环组
        EventLoopGroup group = new NioEventLoopGroup();
        try{
            //创建客户端启动对象
            Bootstrap bootstrap = new Bootstrap();
            //设置相关参数
            //设置线程组
            bootstrap.group(group)
                    //通道实现类
                    .channel(NioSocketChannel.class)
                    .handler(new ChannelInitializer<SocketChannel>() {
                        @Override
                        protected void initChannel(SocketChannel ch) throws Exception {
                            //加入自己的处理器
                            ch.pipeline().addLast(new ClientHandler());
                        }
                    });
            //客户端连接服务端
            ChannelFuture channelFuture =
                    bootstrap.connect("localhost", 8888).sync();
            //给关闭通道进行监听
            channelFuture.channel().closeFuture().sync();
        }finally {
            group.shutdownGracefully();
        }
    }
}
```

客户端处理器

```java
public class ClientHandler extends ChannelInboundHandlerAdapter {
    /**
     * 通道就绪时执行
     * @param ctx
     * @throws Exception
     */
    @Override
    public void channelActive(ChannelHandlerContext ctx) throws Exception {
        ctx.writeAndFlush
                (Unpooled.copiedBuffer("hello,server", CharsetUtil.UTF_8));
    }
    /**
     * 读取通道中的消息
     * @param ctx
     * @param msg
     * @throws Exception
     */
    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        ByteBuf buf = (ByteBuf) msg;
        System.out.println("服务器回复的消息");
        System.out.println(buf.toString(CharsetUtil.UTF_8));
        System.out.println("服务器的地址" +  ctx.channel().remoteAddress());
    }
    /**
     * 处理异常
     * @param ctx
     * @param cause
     * @throws Exception
     */
    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        ctx.channel();
    }
}
```

### 任务队列

当在处理请求时,有一个非常耗时的业务,可以采用异步执行  
将其提交到该`channel`对应的`NIOEventLoop`的`taskQueue`中

- 用户自定义的普通任务

```java
//将任务添加到eventLoop的taskQueue中,从而异步执行
ctx.channel().eventLoop().execcute(new Runable(){

    public void run(){
        //执行的操作
    }
});
```

- 用户自定义的定时任务

```java
//任务提交到scheduleTaskQueue中
ctx.channel.eventLoop().schedule(new Runable(){
    public void run(){
        //执行操作
    }
}, 5, TimeUnit.SECONDS);
```

- 非当前Reactor线程调用Channel的各种方法

### 异步模型

netty的IO操作时异步，通过`Future-Linstener`机制,方便用户可以主动获取或通过通知机制获得IO操作结果

![future](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113032.png)

```java
ChannelFuture cf = bootstrap.bind(8888).sync();
            //注册监听器
            cf.addListener(future -> {
                if(future.isSuccess()){
                    System.out.println("绑定端口成功");
                }
            });
```

### 核心模块

- `Bootstrap`
  `Bootstrap`主要用于配置整个netty程序,串联各个组件,netty中`Bootstrap`类是客户端的启动引导类,`ServerBootstrap`是服务端的
  ![bootstrap](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113036.png)
  - `handler()`方法对`bossGroup`生效
  
- `ChannelHandler`
![入站出站](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113041.png)
入站消息，顺序执行，出站消息倒序执行  
`ChannelInboundHandlerAdapter`中的方法
  - `void channelActive()`通道就绪事件
  - `void channelRead()`通道读取数据事件
  - `void channelReadComplete()`通道数据读取完毕事件  

- `pipeline`

`channelPipeline`是一个handler的集合，它负责处理和拦截inbound或者outbound的事件和操作  
![pipeline](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113045.png)  
入站事件和出站事件在一个双向链表中，入站事件会从链表head往后传递到最后一个入站的headler，出站事件则相反，两种类型的handler互不干扰

- `Unpooled`  
  Netty提供的一个用来操作缓冲区的工具类
  - `ByteBuf copiedBuf(CharSequence string, Charset charset)`

  ```java
  //创建一个包含byte[10]的ByteBuf
  ByteBuf buffer = Unpooled.buffer(10);
  //创建一个包含'hello'的字节ByteBuf
  ByteBuf buffer1 = Unpooled.copiedBuffer("hello", CharsetUtil.UTF_8);
  ```

## 心跳检测机制

```java
//加入netty提供的IdleStateHandler,用于处理空闲状态的处理器
//当3s没读操作时,发送一个心跳检测包,触发一个IdleStateEvent
//5s没写操作时,
//7s既没有读也没有写时,
//当IdelStateEvent事件触发后,就会传递给pipeline的下一个
//通过调用下一个handler的userEventTiggered()方法处理
pipeline.addLast(new IdleStateHandler(3,5,7,TimeUnit.SECONDS));
//加入一个自定义的handler,用于处理事件
pipeline.addLast(new xxxHandler())
```

自定义Handler

```java
public class HeartHandler extends ChannelInboundHandlerAdapter {
    /**
     * @param ctx
     * @param evt 事件
     * @throws Exception
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        if(evt instanceof IdleStateEvent){
            IdleStateEvent event = (IdleStateEvent)evt;
            switch (event.state()){
                case READER_IDLE:
                    System.out.println("读空闲");
                    break;
                case WRITER_IDLE:
                    System.out.println("写空闲");
                    break;
                case ALL_IDLE:
                    System.out.println("读写空闲");
                    break;
                default:
                    break;
            }
        }
    }
}
```

### Protobuf

数据在网络传输中都是二进制数据，在发送时就需要编码(encode)，接收时需要解码(decode)  

Protobuf是Google发布的开源项目,是一种轻便高效的结构化数据存储格式  
protobuf编译器能自动生成代码,protobuf是将类的定义使用.proto文件进行描述

1. 导入protobuf依赖

    ```java
    <dependency>
        <groupId>com.google.protobuf</groupId>
        <artifactId>protobuf-java</artifactId>
        <version>3.6.1</version>
    </dependency>
    ```

1. 编写`.proto`文件

    ```protobuf
    syntax = "proto3";//版本
    option java_outer_classname = "StudentPOJO";//java文件的外部类名
    //protobuf 使用message 管理数据
    message Student{
        //会在StudentPOJO中生成一个内部类，是一个真正发送的数据对象
        int32 id = 1;//1表示属性的一个序号，不是值
        string name = 2;
    }
    ```

1. 编译`.proto`文件
下载地址:`https://github.com/protocolbu ffers/protobuf/releases`

    ```bat
    //protoc.exe --java_out=${OUTPUT_DIR} path/to/your/proto/file
    protoc.exe --java_out=. Student.proto
    ```

1. 使用

    ```java
    //创建对象
    StudentPOJO.Student student = StudentPOJO.Student.newBuilder()
                .setId(4)
                .setName("demo")
                .build();
    ```

   - 编解码器

        ```java
        //加入proto编码器
        pipeline.addLast(new ProtobufEncoder());
        //加入解码器
        //指定对哪个对象解码
        pipeline.addLast(new ProtobufDecoder(StudentPOJO.Student.getDefaultInstance()));
        ```

1. 处理多种类型

    1. `.proto`文件

        ```protobuf
        syntax = "proto3";
        option optimize_for = SPEED;//加快解析
        option java_package = "com.hamilemon.codec2";//指定生成到哪个包
        option java_outer_classname = "DataInfo";//外部类名称

        //protobuf可以使用一个message管理其它message
        message MyMessage{
            //定义一个枚举
            enum DataType{
                //proto3中，enum要求编号从0开始
                StudentType = 0;
                WorkerType = 1;
            }
            //用DataType来标识传递的是哪一个枚举类型
            DataType data_type = 1;
            //表示每次枚举类型最多只能出现其中的一个
            oneof dataBody{
                Student student = 2;
                Worker worker = 3;
            }
        }
        message Student{
            int32 id = 1;
            string name = 2;
        }
        message Worker{
            string name = 1;
            int32 age = 2;
        }
        ```

    2. 编解码器

        ```java
        //加入proto编码器
        pipeline.addLast(new ProtobufEncoder());
        //加入解码器
        //指定对哪个对象解码
        pipeline.addLast(new ProtobufDecoder(DataInfo.MyMessage.getDefaultInstance()));
        ```

    3. 生成对象

        ```java
        DataInfo.MyMessage myMessage = null;
        if (random == 0) {
            //发送Student
            myMessage = DataInfo.MyMessage.newBuilder()
                    .setDataType(DataInfo.MyMessage.DataType.StudentType)
                    .setStudent(DataInfo.Student.newBuilder()
                            .setId(1)
                            .setName("我")
                            .build())
                    .build();
        } else {
            //发送Worker对象
            myMessage = DataInfo.MyMessage.newBuilder()
                    .setDataType(DataInfo.MyMessage.DataType.WorkerType)
                    .setWorker(DataInfo.Worker.newBuilder()
                            .setAge(20)
                            .setName("老王")
                            .build())
                    .build();
        }
        ```

    4. 接收对象

        ```java
        DataInfo.MyMessage.DataType dataType = msg.getDataType();
        if(dataType == DataInfo.MyMessage.DataType.StudentType){
            DataInfo.Student student = msg.getStudent();
            System.out.println(student.getId() + " " + student.getName());
        }else if(dataType == DataInfo.MyMessage.DataType.WorkerType){
            DataInfo.Worker worker = msg.getWorker();
            System.out.println(worker.getAge() + " " +worker.getName());
        }else{
            System.out.println("错误");
        }
        ```

### Netty编解码机制

ChannelPipeline提供了ChannelHandler链的容器,以客户端为例,事件运动方向是从客户端到服务端,这是出站,反之入站  
出站对应写,入站对应读,出站要编码,入站要解码  
入站:从head到tail,顺序执行  
出站:从tail到head,倒序执行

- 自定义解码器
每次入站从ByteBuf中读取数据,将其解码,然后将它添加到List中,当没有更多元素可以被添加是,它的内容将被发送到下一个`ChannelInboundHandler`

    ```java
    public class Byte2LongDecoder extends ByteToMessageDecoder{
        /**
        * @param in 入站的ByteBuf
        * @param out 将解码后的数据传给下一个handler
        * */
        protected void decode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) throws Exception{
            if(in.readableBytes() >= 8){
                out.add(in.readLong());
            }
        }
    }
    ```

- 自定义编码器

    ```java
     public class Long2ByteDecoder extends MessageToByteEncoder<Long>{

        protected void encode(ChannelHandlerContext ctx, Long msg, ByteBuf out) throws Exception{
            out.writeLong(msg);
        }
    }
    ```

- Netty常用编解码器
  - `ReplayingDecoder<S>`扩展了`ByteToMessageDecoder`,使用泛型指定了用户状态管理的类型,使用Void表示不需要管理
  - `LineBasedFrameDecoder`使用`\n`或者`\r\n`作为分割符来解析数据
  - `HttpObjectDecoder`http数据解码器
  - `LengthFieldBasdFrameDecoder`通过指定长度来标识整包消息

### TCP粘包和拆包

发送端为例将多个发给接收端的包,更有效的发给对方,使用优化方法,将多个较小的数据包合并为一个较大的数据包发送,称为粘包
![粘包](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113102.png)
