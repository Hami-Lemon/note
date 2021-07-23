# SpringMVC

Spring的web框架围绕DispatcherServlet设计，DispatcherServlet的作用是将请求分发的看不同的处理器  

![springMvc](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115252.png)

## 中心控制器

spring mvc以请求为驱动,围绕一个中心Servlet分派请求及提供其它功能,DispatcherServlet是一个实际的Servlet  
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
        <param-value>classpath:spring-config.xml</param-value>
    </init-param>
        <load-on-startup>1</load-on-startup>
</servlet>

<servlet-mapping>
    <servlet-name>DispatcherServlet</servlet-name>
    <url-pattern>/</url-pattern>
</servlet-mapping>
```

- spring配置
```xml
<!-- 处理器映射器 -->
<bean class="org.springframework.web.servlet.handler.BeanNameUrlHandlerMapping"/>
<!--处理器设配器 -->
<bean class="org.springframework.web.servlet.mvc.SimpleControllerHandlerAdapter"/>
<!-- 视图解析器:模板引擎 thymeleaf Freemarker -->
<bean class="org.springframework.web.servlet.view.InternalResourceViewResolver"
        id="internalResourceViewResolver">
        <!-- 前缀 -->
    <property name="prefix" value="/WEB-INF/jsp/"/>
    <!-- 后缀 -->
    <property name="suffix" value=".jsp"/>
</bean>
<!-- 控制器 -->
<bean id="/hello" class="com.hamilemon.controller.HelloController"/>
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

- spring配置文件

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:context="http://www.springframework.org/schema/context"
       xmlns:mvc="http://www.springframework.org/schema/mvc"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
       http://www.springframework.org/schema/beans/spring-beans.xsd
       http://www.springframework.org/schema/context
       https://www.springframework.org/schema/context/spring-context.xsd
       http://www.springframework.org/schema/mvc
       https://www.springframework.org/schema/mvc/spring-mvc.xsd">

    <context:component-scan base-package="com.hamilemon.controller"/>
<!--资源过滤，让spring不处理静态资源-->
    <mvc:default-servlet-handler/>
<!--    开启注解主持-->
    <mvc:annotation-driven/>

<!--    识图解析器-->
    <bean class="org.springframework.web.servlet.view.InternalResourceViewResolver"
          id="internalResourceViewResolver">
        <property name="prefix" value="/WEB-INF/jsp/"/>
        <property name="suffix" value=".jsp"/>
    </bean>
</beans>
```

- 创建Controller类  
使用`@Controller`

```java
/**
 * @author Hami Lemon
 */
@Controller
@RequestMapping("/hello")//可不设置
public class HelloController {
    /**访问地址：ip:端口/项目名/hello/h1
     *
     * @param model 添加数据
     * @return 需要跳转的jsp页面名称
     */
    @RequestMapping("/h1")
    public String hello(Model model){
        //封装数据
        model.addAttribute("msg", "hellofasf");
        //回被视图解析器处理,返回需要跳转的jsp页面名称
        //WEB-INF/jsp/hello.jsp
        return "hello";
    }
}
```

## RestFul风格

Restful就是一个资源定位及资源操作的风格，不是标准也不是协议，只是一种风格  
基于这个风格设计的软件更简洁，更有层次,更易于实现缓存等机制  

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

使用`http://localhost:8080/spring-anno/restful/1/2`访问

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

  1. 方式二

    ```java
    @GetMapping("/restful/{a}/{b}")
    @PostMapping("/restful/{a}/{b}")
    @PutMapping("/restful/{a}/{b}")
    @DeleteMapping("/restful/{a}/{b}")
    ```

## 重定向和转发

- 没有视图解析器实现
需要注释掉视图解析器

```java
return "/index.jsp";//转发

return "forward:/index.jsp"; //转发

return "redirect:/index.jsp"; //重定向
```

- 有视图解析器

```java
return "test"; //需要转发到的另一个页面

return "redirect:/index.jsp"; //重定向
```

## 处理数据

### 接收数据

- 前端传递的参数名与后端的变量名相同  

  ```java
  @xxxMapping("xxx")
  public String func(String name){
      //ip:端口?name=xxx
      //name会自动接收到
      return "xxx";
  }
  ```

- 传递的参数名和后端变量名不同

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
<!-- 处理乱码-->
    <filter>
        <filter-name>encoding</filter-name>
        <filter-class>org.springframework.web.filter.CharacterEncodingFilter</filter-class>
    </filter>
    <filter-mapping>
        <filter-name>encoding</filter-name>
        <url-pattern>/*</url-pattern>
    </filter-mapping>
```

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

json解析工具:jackson, fastson  

jackson maven依赖

```xml
<!-- https://mvnrepository.com/artifact/com.fasterxml.jackson.core/jackson-databind -->
<dependency>
    <groupId>com.fasterxml.jackson.core</groupId>
    <artifactId>jackson-databind</artifactId>
    <version>2.11.0</version>
</dependency>

```

- 使用`@ResponseBody`作用在方法上,则该方法不会走视图解析器,会直接返回字符串
- 使用`@RestController`替换`@Controller`表示当前类的所有方法都不走视图解析器

### jackson使用

```java
@Controller
public class JsonController {

    @RequestMapping("/j1")
    @ResponseBody
    public String json() throws JsonProcessingException {
        User user = new User("李华",20, "男");
        ObjectMapper mapper = new ObjectMapper();
        //将user对象转为json字符串
        String str = mapper.writeValueAsString(user);
        return str;
    }
}
```

- 解决json乱码

```xml
<!--    解决json乱码-->
    <mvc:annotation-driven>
        <mvc:message-converters register-defaults="true">
            <bean class="org.springframework.http.converter.StringHttpMessageConverter">
                <constructor-arg value="UTF-8"/>
            </bean>

            <bean class="org.springframework.http.converter.json.MappingJackson2HttpMessageConverter">
                <property name="objectMapper">
                    <bean class="org.springframework.http.converter.json.Jackson2ObjectMapperFactoryBean">
                        <property name="failOnEmptyBeans" value="false"/>
                    </bean>
                </property>
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

```java
@RequestMapping(value = "/j2",produces = "application/json;charset=utf-8")
    @ResponseBody
    public String json2(){
        List<User> users = new ArrayList<>();
        users.add(new User("李华", 20, "女"));
        users.add(new User("李华", 20, "女"));
        users.add(new User("李华", 20, "女"));
        users.add(new User("李华", 20, "女"));
        users.add(new User("李华", 20, "女"));
        users.add(new User("李华", 20, "女"));
        users.add(new User("李华", 20, "女"));
        return JSON.toJSONString(users);
    }
```
