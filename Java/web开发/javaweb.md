

# JavaWeb

通过java语言编写,可以通过浏览器访问的程序(基于请求和响应开发)

## xml

用来传输和存储数据

```xml
<?xml version="1.0" encoding="utf-8" ?>
<books>
    <book sn="ss">
        <name>时间简史</name>
    </book>
</books>
```

### 约束

规定xml文档的书写规则

1. dtd:简单的约束技术
2. schema:复杂的约束技术

### 解析XML

[dom4j](https://dom4j.github.io/)

```java
//1. 创建一个saxReader，去读取xml文件
SAXReader saxReader = new SAXReader();
Document document = saxReader
        .read(Parse.class.getClassLoader().getResourceAsStream("books.xml"));
System.out.println(document);
//通过document获取根元素
Element rootElement = document.getRootElement();
System.out.println(rootElement);
//通过根元素获取标签对象
//遍历
for(Element book : rootElement.elements("book")){
    Element element = book.element("name");
    System.out.println(element.getText());
    System.out.println(element.asXML());
}
```

## Maven

1. 配置国内镜像`conf/setting.xml`中的`<mirrors>`

```xml
<mirror>
        <id>alimaven</id>
        <mirrorOf>central</mirrorOf>
        <name>aliyun maven</name>
        <url>http://maven.aliyun.com/nexus/content/repositories/central/</url>
</mirror>
```

2. 本地仓库

```xml
<localRepository>xxx</localRepository>
```

## Servlet

### 继承Servlet接口

- init
   初始化
- getServletConfig
  获取配置信息
- service
  提供服务
- getServletInfo
  获取Servlet基本信息
- destroy
  销毁
- 处理方式
  1. 第一次访问时，服务器会创建Servlet对象，调用init方法，再调用service方法
  2. 第二次访问时，Servlet对象已经存在，不再创建，直接执行service方法
  3. 当服务器停止时，会释放servlet，调用destroy方法，servlet对象会在堆中，由GC回收

### GenericServlet抽象类

提供生命周期方法init和destroy的简单实现，要编写一般的servlet，只需重写抽象service方法(与协议无关)

### HttpServlet类

继承GenericServlet的基础上进一步的扩展,提供将要被子类化创建使用于web站点的HTTP servlet的抽象类.
HttpServlet的子类必须重写一个方法

- doGet
- doPost
- doPut
- doDelete

### web.xml配置

1. `<url-patern>`
   1. 精确匹配
        `<url-pattern>/m</url-pattern>`
   2. 后缀匹配
        `<url-pattern>*.action</url-pattern>`
   3. 全匹配
        `<url-pattern>/*</url-pattern>`
        输入任何内容匹配的资源都是当前servlet,但是不会影响精确匹配
2. `<load-on-startup>`
放在`<servlet>`中

标记容器是否应该在web应用程序启动时加载这个servlet

值必须是一个整数,表示servlet被加载的先后顺序

如果该元素的值为负数或者没有设置,则容器会当servlet被请求时再加载

如果为正整数或者0时,表示容器在应用启动时就加载并初始化,值越小,优先级越高
3. `<welecom-file-list>`

配置默认显示的界面

```xml
<welcome-file-list>
    <welcome-file>xxx</welcome-file>
    <welcome-file>xxx</welcome-file>
</welcome-file-list>
```

1. `<error-page>`
配置错误页面

```xml
    <!-- 错误码 -->
    <error-code>404<error-code>
    <!-- 显示的页面 -->
    <location>/xxx</location>
```

### 使用注解配置

Servlet类上使用`@WebServlet`
参数:

1. name:名称
2. value:路径匹配
3. urlPatterns:同上
4. loadOnStartup:设置启动优先级

### 获取请求参数

通过`doGet`或`doPost`中的HttpServletRequest参数获取

使用`request`的`setCharacterEncoding("utf-8")`设置编码

- GET请求
  GET提交的数据会放在URL之后,以?分割,参数间以&连接

  提交的数据大小有限制

- POST请求
  post方法是把提交的数据放在HTTP包的body中

  提交的数据没有限制,相对安全

### 页面跳转

- 重定向
![重定向](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114949.png)
使用HttpServletResponse的`sendRedirect("/xxx")`(相对于整个容器,localhost:8080)

- 请求转发
![转发](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114954.png)

```java
// 相对于当前工程
request.getRequestDispatcher("/xxx").forward(request, response);
```

### Servlet生命周期

1. 实例化(调用构造方法)
2. 初始化(init方法) 只会被执行一次
3. 就绪/服务(service方法)
4. 销毁(destroy方法)

### 线程安全问题

多个线程访问同一个servlet时,共享成员变量

### Servlet初始化参数

- web.xml方式

  ```xml
  <init-param>
    <param-name>name</param-name>
    <param-value>张三</param-value>
  </init-param>

  <!-- init-param 用来定义Servlet启动的参数,可以有多个(一个init-param中只能有一对),放在相应的<servlet>中-->
  <!-- param-name表示参数名称 -->
  <!-- param-value表示参数值 -->
  ```

- 注解方式

  ```java
  initParams = {@WebInitParam(name="username",value="aaa")}
  ```

- 使用
  通过`getServletConfig`方法获取config对象

### 状态管理(Cookie)

web应用中的会话是指一个浏览器与web服务器之间连续发生的一系列请求和响应过程

web应用的会话状态是指web服务器于浏览器在会话过程中产生的状态信息,借助会话状态
web服务器能够把属于同意会话中的一系列的请求和响应过程关联起来

#### 客户端状态管理技术:将状态保存在客户端(Cookie)

  Cookie是在浏览器访问web服务器的某个资源时,由web服务器在响应消息头中附带传送给浏览器的一小段数据

```java
//创建一个Cookie,默认生命周期时浏览器关闭
Cookie cookie = new Cookie("key", "value");

//设置Cookie的生命周期
// 负数:浏览器关闭时销毁(默认)
// 0:失效,让已经存在的cookie失效
// 正数:有效时间(秒)
cookie.setMaxAge(0);

//设置共享范围
cookie.setPath("/");//整个服务器下

//发送给浏览器
response.addCookie(cookie);
```

- 查询cookie

```java
Cookie[] cookies = request.getCookies();
```

- 修改cookie

```java
//新建一个名和路径与原来相同的cookie,即可实现覆盖
Cookie cookie = new Cookie("key", "value");
cookie.setPath("/");
response.addCookie(cookie);
```

- cookie编码(解决中文乱码)  
  `URLEncoder.encode("", "UTF-8")`编码  
  `URLDecoder.decode("", "UTF-8")`解码

#### 服务器状态管理:将状态保存在服务器端(session)

session用于跟踪客户端的状态,session值的是在一段时间内,单个客户与web服务器的一连串交互过程  session被用于表示一个持续的连接状态,在网站访问中一般自带客户端浏览器的进程从开始到结束的过程  实现机制是当用户发起一个请求的时候,服务器会检查该请求中是否包含sessionid,如果未包含,则Tomcat会创造一个名为jsessionid的输出cookie返回给浏览器,当已经包含sessionid时,服务端会检查找到与该session相匹配的信息

```java
//获取session对象
HttpSession session = request.getSeesion();
session.getId();//sesssion唯一标记
session.getMaxInactiveInterval();//获取最大过期时间(秒)(默认30分钟)
//30分钟没有操作,就会过期

//删除session
session.invalidata();
```

- 修改session超时时间
  1. 使用`session.setMaxInactiveInterval(20*60)`(秒)
  1. web.xml设置(分钟)
  
  ```xml
  <session-config>
    <session-timeout>20</session-timeout>
  </session-config>
  ```

### ServletContext

servlet上下文,代表当前整个应用程序  当web服务器启动时,会为每个web应用程序创建一块共享存储区域 ServletContext在web服务器启动时创建,服务器关闭时销毁  

#### 获取ServletContext

1. GenericServlet提供了getServletContext()(推荐)  
   `this.getServletContext()`
2. ServletConfig提供了getServletContext()  
   `this.getServletConfig().getServletContext()`
3. HttpSession提供了getSrevletContext()  
   `request.getSession().getServletContext()`
4. HttpServletRequest提供了getServletContext()(推荐)  
   `request.getServletContext()`

#### ServletContext作用

- 获取当前项目的发布路径  
  `servletContext.getRealPath("/")`
- 获取容器附加信息
  
  ```java
  servletContext.getServerInfo();//Apache Tomcat/9.0.34
  servletContext.getContextPath();// /java
  ```

- 全局容器

  ```java
  servletContext.setAttribute("s", obj);//存放数据
  servletContext.getAttribute("s");//获取数据
  servletContext.removeAttribute("s");//移除数据
  ```

## 过滤器（Filter）

![过滤器](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115000.png)
可以通过filter对web服务器管理的所有资源进行拦截，servlet中提供了一个Filter接口,实现这个接口的类称为过滤器

1. 实现Filter接口
2. 实现doFilter方法
3. 设置拦截的URL

```java
public class HelloFilter implements Filter {
  /**
  * @param filterChain Filter链
  */
    @Override
    public void doFilter(ServletRequest servletRequest, ServletResponse servletResponse, FilterChain filterChain) throws IOException, ServletException {
      //调用下一个过滤器,如果没有,则传递到servlet
      filterChain.doFilter(servletRequest, servletResponse);

    }
}
```

#### 配置filter

1. web.xml

```xml
<filter>
<!-- 名称 -->
  <filter-name>sf</filter-name>
  <!-- 过滤器类 -->
  <filter-class>com.hamilemon.filters.HelloFilter</filter-class>
</filter>
<!-- 映射路径配置 -->
<filter-mapping>
  <filter-name>sf</filter-name>
  <!-- 匹配规则 同servlet一样-->
  <url-pattern>/*</url-pattern>
  <!-- 过滤指定的servlet -->
  <servlet-name>hello</servlet-name>
</filter-mapping>
```

2. 注解
   过滤器类上使用`@WebFilter`  

优先级:web.xml优先级最高,注解配置的按照类名顺序  

### filter链

所有注册的filter会构造一条链，按照优先级依次调用（过滤响应时顺序相反），每一个filter需要调用`chain`的`doFilter`方法向后传递，以便后续filter继续处理，如不调用，则该请求不会被继续传递也就不会到达servlet。

在`doFilter`方法前的代码在接受到请求时执行，而在其后的代码则在服务端响应时执行。

![image-20210724100747394](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/07/24/20210724100747.png)

## Listener

监听器用于监听WEB应用中某些对象的创建、销毁、增加、修改、删除等动作的发生，然后作出相应的处理。当监听范围对象的状态发生变化时，服务器自动调用监听器中的方法。常用于统计网站在线人数、系统加载是进行信息初始化、统计网站的访问量等等。

### 分类

- 按监听的对象分
  - ServletContext对象监听器
  - HttpSession对象监听器
  - ServertRequest对象监听器
- 按监听的事件分
  - 对象自身的创建与销毁监听器
  - 对象中属性的创建和消除监听器
  - session中的某个对象的状态变化监听器

### 统计网站在线人数

每当有一个访问连接到服务器时，都会创建一个session，所以可以统计session的数量来获得当前在线人数。使用`HttpSessionListener`

```java
//使用注解注册
@WebListener
public class HelloListener implements HttpSessionListener {
    private int count;

    @Override
    public void sessionCreated(HttpSessionEvent se) {
        count++;
        se.getSession().setAttribute("count", count);
    }

    @Override
    public void sessionDestroyed(HttpSessionEvent se) {
        count--;
        se.getSession().setAttribute("count", count);
    }
}
```

### 常用Listener

1. `HttpSessionAttributeListener` 监听session中属性的增加，移除以及改变
2. `ServletContextListener` 监听web上下文的初始化
3. `ServletContextAttributeListener` 监听web上下文中属性的变化
4. `ServletRequestListener`监听`request`的创建与销毁
5. `ServletRequestAttributeListener`监听request的属性的增加、删除、属性值变化
