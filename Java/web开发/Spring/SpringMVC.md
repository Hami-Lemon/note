# SpringMVC

Spring的web框架围绕DispatcherServlet设计，DispatcherServlet的作用是将请求分发的看不同的处理器  

![springMvc](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115252.png)

## 中心控制器

spring mvc以请求为驱动,围绕一个中心Servlet分派请求及提供其它功能,DispatcherServlet实际上是一个Servlet 
![继承关系](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115247.png)
![过程图](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115240.png)

### 配置

- web.xml

```xml
<!-- 配置DispatcherServlet -->
<servlet>
    <servlet-name>DispatcherServlet</servlet-name>
    <servlet-class>org.springframework.web.servlet.DispatcherServlet</servlet-class>
    <init-param>
        <param-name>contextConfigLocation</param-name>
        <!-- 指定SpringMVC的配置文件路径 -->
        <param-value>classpath:spring-config.xml</param-value>
    </init-param>
        <load-on-startup>1</load-on-startup>
</servlet>

<servlet-mapping>
    <servlet-name>DispatcherServlet</servlet-name>
    <url-pattern>/</url-pattern>
</servlet-mapping>
```

## 使用XML配置

- SpringMVC配置文件

  ```xml
  <!-- 处理器映射器 -->
  <bean class="org.springframework.web.servlet.handler.BeanNameUrlHandlerMapping"/>
  <!--处理器设配器 -->
  <bean class="org.springframework.web.servlet.mvc.SimpleControllerHandlerAdapter"/>
  <!-- 视图解析器-->
  <bean class="org.springframework.web.servlet.view.InternalResourceViewResolver"
          id="internalResourceViewResolver">
          <!-- 前缀 -->
      <property name="prefix" value="/WEB-INF/jsp/"/>
      <!-- 后缀 -->
      <property name="suffix" value=".jsp"/>
  </bean>
  <!-- 控制器 ,id为其所响应的链接-->
  <bean id="/hello" class="com.hamilemon.controller.HelloController"/>
  ```

- 控制器代码

  ```java
  public class HelloController implements Controller {
      @Override
      public ModelAndView handleRequest(HttpServletRequest request, HttpServletResponse response) throws Exception {
          ModelAndView mv = new ModelAndView();
          //向视图中添加数据，以用于模板进行解析
          mv.addObject("msg", "hello");
          //name为该视图对应的文件名，完全名为前缀+name+后缀
          //如：/WEB-INF/jsp/hello.jsp
          mv.setViewName("hello");
          return mv;
      }
  }
  ```

### 使用Thymeleaf作为模板引擎

1. 添加依赖

   ```xml
   <dependency>
           <groupId>org.thymeleaf</groupId>
           <artifactId>thymeleaf-spring5</artifactId>
           <version>3.0.12.RELEASE</version>
   </dependency>
   ```

2. 配置thymeleaf

   - XML配置

     ```xml
     <!--设置模板解析器-->
     <bean id="templateResolver"        class="org.thymeleaf.spring5.templateresolver.SpringResourceTemplateResolver">
             <!--设置前缀-->
             <property name="prefix" value="/WEB-INF/html/"/>
             <!--设置后缀-->
             <property name="suffix" value=".html"/>
             <!--解析模式，默认即为HTML也能处理js和css-->
             <property name="templateMode" value="HTML"/>
             <!--是否开启缓存-->
             <property name="cacheable" value="true"/>
         </bean>
         <!--设置模板引擎-->
         <bean id="templateEngine"
               class="org.thymeleaf.spring5.SpringTemplateEngine">
             <property name="templateResolver" ref="templateResolver"/>
             <!--是否使用SpringEL表达式-->
             <property name="enableSpringELCompiler" value="false"/>
         </bean>
     ```

   - 注解配置

     ```java
     @Bean
     public SpringResourceTemplateResolver templateResolver(){
         SpringResourceTemplateResolver templateResolver = new SpringResourceTemplateResolver();
         templateResolver.setApplicationContext(this.applicationContext);
         //前缀
         templateResolver.setPrefix("/WEB-INF/html/");
         //后缀
         templateResolver.setSuffix(".html");
         templateResolver.setTemplateMode(TemplateMode.HTML);
         templateResolver.setCacheable(true);
         return templateResolver;
     }
     
     @Bean
     public SpringTemplateEngine templateEngine(){
         SpringTemplateEngine templateEngine = new SpringTemplateEngine();
         templateEngine.setTemplateResolver(templateResolver());
         templateEngine.setEnableSpringELCompiler(true);
         return templateEngine;
     }
     ```

3. 设置thymeleaf作为模板引擎（视图解析器）

   - xml配置

     ```xml
     <!--使用thymeleaf作为视图解析器-->
         <bean class="org.thymeleaf.spring5.view.ThymeleafViewResolver">
             <property name="templateEngine" ref="templateEngine"/>
         </bean>
     ```

   - 注解配置

     ```java
     @Bean
     public ThymeleafViewResolver viewResolver(){
         ThymeleafViewResolver viewResolver = new ThymeleafViewResolver();
         viewResolver.setTemplateEngine(templateEngine());
         return viewResolver;
     }
     ```

4. 使用

   ```html
   <!doctype html>
   <html lang="zh-CN" xmlns:th="http://www.thymeleaf.org">
   <head>
       <meta charset="UTF-8">
       <title>Document</title>
   </head>
   <body>
   <p th:text="${msg}"></p>
   </body>
   </html>
   ```

## 使用注解开发

- web.xml

```xml
<?xml version="1.0" encoding="UTF-8"?>
<web-app xmlns="http://xmlns.jcp.org/xml/ns/javaee"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://xmlns.jcp.org/xml/ns/javaee http://xmlns.jcp.org/xml/ns/javaee/web-app_4_0.xsd"
         version="4.0">
    <!-- 配置DispatcherServlet -->
    <servlet>
        <servlet-name>DispatcherServlet</servlet-name>
        <servlet-class>org.springframework.web.servlet.DispatcherServlet</servlet-class>
        <init-param>
            <param-name>contextConfigLocation</param-name>
            <param-value>classpath:spring-config.xml</param-value>
        </init-param>
        <load-on-startup>1</load-on-startup>
    </servlet>

    <servlet-mapping>
        <servlet-name>DispatcherServlet</servlet-name>
        <url-pattern>/</url-pattern>
    </servlet-mapping>
</web-app>
```

- SpringMVC配置文件

```xml
<!--扫描controller包-->
    <context:component-scan base-package="com.hamilemon.controller"/>
    <!--资源过滤，让spring不处理静态资源-->
    <mvc:default-servlet-handler/>
    <!--开启MVC注解主持-->
    <mvc:annotation-driven/>
    <!--视图解析器-->
    <bean class="org.springframework.web.servlet.view.InternalResourceViewResolver"
          id="internalResourceViewResolver">
        <property name="prefix" value="/WEB-INF/jsp/"/>
        <property name="suffix" value=".jsp"/>
    </bean>
```

- 创建Controller类

```java
@Controller
@RequestMapping("/hello")//设置该controller处理的路径，可不设置
public class HelloController {
    /**
     * 访问地址：ip:端口/项目名/hello/h1
     * @param model 添加数据
     * @return 需要跳转的jsp页面名称
     */
    @RequestMapping("/h1")
    public String hello(Model model){
        //封装数据
        model.addAttribute("msg", "hello");
        //回被视图解析器处理,返回需要跳转的jsp页面名称
        //WEB-INF/jsp/hello.jsp
        return "hello";
    }
}
```

## RestFul风格

Restful就是一个资源定位及资源操作的风格，不是标准也不是协议，只是一种风格，基于这个风格设计的软件更简洁，更有层次,更易于实现缓存等机制。

- 功能
  - 资源:互联网所有的事物都可以被抽象为资源
  - 资源操作:使用POST,DELETE,PUT,GET,对资源进行操作
- 使用Restful风格操作资源

```md
http://localhost/item/1     (查询,get)
http://localhost/item       (新增,post)
http://localhost/item       (更新,put)
htto://localhost/item/1     (删除,delete)
```

- 使用

例如：访问`http://localhost:8080/spring-anno/restful/1/2`

```java
@RequestMapping("/restful/{a}/{b}")
    public String restful(@PathVariable int a,@PathVariable int b, Model model){
        model.addAttribute("msg", a + b);
        return "hello";
    }
```

- 限定请求方式
  1. 方式一

     ```java
     @RequestMapping(value = "/restful/{a}/{b}",method = RequestMethod.GET)
     ```

  2. 方式二

     ```java
     @GetMapping("/restful/{a}/{b}")
     @PostMapping("/restful/{a}/{b}")
     @PutMapping("/restful/{a}/{b}")
     @DeleteMapping("/restful/{a}/{b}")
     ```


## 重定向和转发

- 没有视图解析器(需要注释掉视图解析器)

```java
return "/index.jsp";//转发到某一个页面
return "/hello"; //转发到一个链接上
return "forward:/index.jsp"; //转发
return "redirect:/index.jsp"; //重定向
```

- 有视图解析器

```java
return "forward:/hello/h1"; //转发到/hello/h1上
return "redirect:/index.jsp"; //重定向
```

## 处理数据

### 接收数据

- 前端传递的参数名与后端的变量名相同  

  自动根据名称进行匹配

  ```java
  @xxxMapping("xxx")
  public String func(String name){
      //ip:端口?name=xxx
      return "xxx";
  }
  ```

- 传递的参数名和后端变量名不同

    通过`RequestParam`注解设置对应前端参数的名称

    ```java
    @xxxMapping("xxx")
    public String func(@RequestParam("username") String name){
          //ip:端口?username=xxx
          //name会自动接收到
        return "xxx";
    }
    ```

- 传递的是一个对象 
  会自动匹配User对象中的属性名和表单名,如果相同,就会自动封装

    ```java
    @xxxMapping("xxx")
    public String func(User user){
          //ip:端口?id=xxx&name=xxx
          //name会自动接收到
        return "xxx";
    }
    ```

### 数据回显

数据会被写入到视图中，由视图解析器读取并解析

- 使用ModelMap
`ModelMap`继承自`LinkedHashMap`

```java
@xxxMapping("xxx")
public String func(Model map){
    map.addAttribute("msg", "xxx");
    return "xxx";
}
```

- 使用Model
`Model`继承了`ModelMap`

```java
@xxxMapping("xxx")
public String user(Model model){
    model.addAttribute("msg", "xxx");
    return "hello";
}
```

## 中文乱码问题

配置过滤器(Spring已经实现了一个)

```xml
<filter>
        <filter-name>encoding</filter-name>
        <filter-class>org.springframework.web.filter.CharacterEncodingFilter</filter-class>
	<init-param>
     	<param-name>encoding</param-name>
          <param-value>UTF-8</param-value>
	</init-param>
</filter>
<filter-mapping>
    <filter-name>encoding</filter-name>
    <url-pattern>/*</url-pattern>
</filter-mapping>
```

## 静态资源处理

由于SpringMVC的DispatcherServlet在配置时会处理所有的路径，所以会把对静态资源的路径也给拦截，从而导致静态资源无法访问。

### 处理方式一

使用Tomcat中的defaultServlet来处理静态资源

```xml
<servlet-mapping>
    <servlet-name>default</servlet-name>
    <url-pattern>*.html</url-pattern>
</servlet-mapping>
<servlet-mapping>
    <servlet-name>default</servlet-name>
    <url-pattern>*.jsp</url-pattern>
</servlet-mapping>
```

defaultServlet为Tomcat默认创建的Servlet，可直接使用。

### 处理方式二

使用SpringMVC的`mvc:default-servlet-handler`，会直接将所有静态资源转发到defaultServlet中。

在SpringMVC的配置文件中添加`<mvc:default-servlet-handler/>`即可。

### 处理方式三

使用SpringMVC的`<mvc:resources>`，自定义静态资源的路径

```xml
<!--location 表示资源所在的路径，相对于webapp目录，mapping则是路径的匹配规则-->
<mvc:resources mapping="/images/**" location="/images/"/>
```

通常会将方式二和方式三结合使用。

## JSON

JSON(JavaScript Object Notation, JS 对象简谱) 是一种轻量级的数据交换格式。它基于 ECMAScript (欧洲计算机协会制定的js规范)的一个子集，采用完全独立于编程语言的文本格式来存储和表示数据。简洁和清晰的层次结构使得 JSON 成为理想的数据交换语言。 易于人阅读和编写，同时也易于机器解析和生成，并有效地提升网络传输效率。

```json
{
    "name":"xxx",
    "age":3,
    "set":"女"
}
```

### Controller 返回JSON

jackson maven依赖

```xml
<!-- https://mvnrepository.com/artifact/com.fasterxml.jackson.core/jackson-databind -->
<dependency>
    <groupId>com.fasterxml.jackson.core</groupId>
    <artifactId>jackson-databind</artifactId>
    <version>2.11.0</version>
</dependency>
```

- 使用`@ResponseBody`作用在方法上,则该方法不会走视图解析器,会直接将返回的java对象转换成json。
- 使用`@RestController`替换`@Controller`表示当前类的所有方法都不走视图解析器

### Jackson使用

默认会使用jackson作处理，将java对象转换成json，无需配置，添加依赖即可。

```java
@Controller
public class JsonController {

    @RequestMapping("/j1")
    @ResponseBody
    public User json() throws JsonProcessingException {
        User user = new User("李华",20, "男");
        return user;
    }
}
```

- 解决json乱码

```xml
    <mvc:annotation-driven>
        <!--register-defaults指定是否使用默认的转换器-->
        <mvc:message-converters register-defaults="true">
            <bean class="org.springframework.http.converter.StringHttpMessageConverter">
                <constructor-arg value="UTF-8"/>
            </bean>
        </mvc:message-converters>
    </mvc:annotation-driven>
```

### Fastjson使用

依赖

```xml
<dependency>
       <groupId>com.alibaba</groupId>
       <artifactId>fastjson</artifactId>
       <version>1.2.68</version>
</dependency>
```

SpringMVC配置

```xml
    <mvc:annotation-driven>
        <mvc:message-converters>
            <bean class="com.alibaba.fastjson.support.spring.FastJsonHttpMessageConverter">
                <property name="supportedMediaTypes">
                    <list>
                        <value>application/json</value>
                    </list>
                </property>
            </bean>
        </mvc:message-converters>
    </mvc:annotation-driven>
```

## 自定义数据转换

在json处理中使用了一个`XXXHttpMessageConverter`，这个类由不同的json处理框架实现，用于将java对象转成json，我们也可以自定义一个数据转换器，用于支持自定义的数据格式。

- 新建一个类，并继承`AbstractHttpMessageConverter<T>`，重写其中的方法。

  这里仅作演示，只对Pojo一个类作处理，且处理方式为直接获取前端发送的数据而不作解析。

  ```java
  public class MyConverters extends AbstractHttpMessageConverter<Pojo> {
      private Charset charset;
  
      public MyConverters() {
          //设置支持处理的contentType为application/demo和text/html
          super(StandardCharsets.UTF_8,
                  new MediaType("application", "demo"),
                  MediaType.TEXT_HTML);
          charset = getDefaultCharset();
      }
  
      @Override
      protected boolean supports(Class<?> clazz) {
          //判断数据是否为支持的类型
          return Pojo.class == clazz;
      }
  
      @Override
      protected Pojo readInternal(Class<? extends Pojo> clazz, HttpInputMessage inputMessage) throws IOException, HttpMessageNotReadableException {
          //读取前端发送的数据
          final InputStream in = inputMessage.getBody();
          byte[] buffer = new byte[1024];
          int len;
          StringBuilder str = new StringBuilder();
          while ((len = in.read(buffer)) != -1) {
              str.append(new String(buffer, 0, len, charset));
          }
          //Pojo中只有一个String属性，这里直接将读取到的数据转成String，并作为参数传入
          //开发中应根据实际情况作相应的解析
          return new Pojo(str.toString());
      }
  
      @Override
      protected void writeInternal(Pojo pojo, HttpOutputMessage outputMessage) throws IOException, HttpMessageNotWritableException {
          //向前端写数据
          final OutputStream body = outputMessage.getBody();
          //这里header只能读，不能设置值，contentType和contentLength 都会由父类进行设置
          final HttpHeaders headers = outputMessage.getHeaders();
          //这里直接将前端传过来的数据再写回去
          byte[] buffer = pojo.name.getBytes(charset);
          body.write(buffer);
      }
  }
  ```

- 在SpringMVC中注册自定义的数据转换器

  ```xml
  <mvc:annotation-driven>
  	<mvc:message-converters>
          <bean class="com.hamilemon.controller.MyConverters">
              <!--这里的属性可以不设置，因为自定义的类中默认已设置-->
              <!--<property name="supportedMediaTypes">
                  <list>
                      <value>application/demo</value>
       				<value>text/html</value>
                  </list>
               </property>
               <property name="defaultCharset" value="UTF-8"/>-->
          </bean>
      </mvc:message-converters>
  </mvc:annotation-driven>
  ```

- Controller类

  ```java
  @RequestMapping(value = "/j3", 
              consumes = "application/demo", //指定请求中的contentType
              produces = "text/html") //指定响应的contentType
      @ResponseBody
      public Pojo json3(@RequestBody Pojo json) {
          //直接返回前端发送的请求体
          return json;
      }
  ```

- 前端请求

  ```js
  $.ajax({
  	type: "post",
  	url: "/json/j3",
  	contentType: "application/demo;charset=UTF-8",
  	data: {
          data: "data",
  		demo: "demo"
  	},
  	success: function (response) {
  		console.log(response);
  	},
  	error: function (error) {
  		console.log(error);
  	}
  	});
  });
  ```

- 前端获取到的响应结果

  ```
  data=data&demo=demo
  ```

  

