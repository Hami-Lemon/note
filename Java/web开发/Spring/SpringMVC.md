# Spring Web MVC

本笔记整理自[Spring Web MVC](https://docs.spring.io/spring-framework/docs/current/reference/html/web.html)

Spring Web MVC是一个基于Servlet而构建的web框架。

## DispatcherServlet

Spring MVC围绕前端控制器模式设计，其中`DispatcherServlet`作为中央`Servlet`，提供了公共的请求处理算法，而将实际工作交由具体的组件去完成。和其它`Servlet`一样，`DispatcherServlet`也需要根据Servlet规范，在`web.xml`中声明。反过来，`DispatcherServlt`则根据Spring的配置，去发现用于请求映射，页面处理的组件。

```xml
<web-app>

    <listener>
        <listener-class>org.springframework.web.context.ContextLoaderListener</listener-class>
    </listener>

    <context-param>
        <param-name>contextConfigLocation</param-name>
        <param-value>/WEB-INF/app-context.xml</param-value>
    </context-param>

    <servlet>
        <servlet-name>app</servlet-name>
        <servlet-class>org.springframework.web.servlet.DispatcherServlet</servlet-class>
        <init-param>
            <param-name>contextConfigLocation</param-name>
            <!--配置了context-param 所以这里可以留空-->
            <param-value></param-value>
        </init-param>
        <load-on-startup>1</load-on-startup>
    </servlet>

    <servlet-mapping>
        <servlet-name>app</servlet-name>
        <url-pattern>/</url-pattern>
    </servlet-mapping>

</web-app>
```

`DispatcherServlet`使用一个`WebApplicationContext`(Spring中`ApplicationContext`的子接口）作为自己的配置，它会通过这个`WebApplicationContext`去查找自己所需要的bean（例如视图解析器之类的）。

`DispatcherServlet`处理请求流程为：

- `WebApplicationContext`会被搜索到并将其作为一个属性绑定到请求上，以便于控制器和其它元素能够使用。
- 绑定相关的解析器到请求上，以便后续使用。
- 如果指定了multipart（http协议中的请求头）文件解析器 ，则会检查请求中的multipart。如果相关的multipart被找到，则将请求包装成`MultipartHttpServletRequest`，以便后续处理。
- 查找合适的处理程序，如果被找到则执行相关的处理链（预处理，后处理和控制器controllers）去准备model（包含视图中所需要的数据，也就是MVC模型中的M）用于渲染视图。或者，对于使用注解的controllers，可以直接返回响应而不返回视图（前后端分离）。
- 如果model被返回，则去渲染视图。如果没有model返回，则没有视图会被渲染，因为这个请求可能已经处理完成。

### 拦截器

所有的`HandlerMapping`的实现类都支持处理拦截器，拦截器可以在实际处理请求前后对其进行预处理和后处理。拦截器必须实现`HandlerInterceptor`接口并实现以下三个方法：

- `preHandle`:在实际处理程序运行前执行。
- `postHandle`:在实际处理程序运行后执行。
- `afterCompletion`:请求被处理完成后。

`preHandle`返回一个布尔值，如果为true，处理程序链将会继续执行，如果为false则会中断执行。

注意：`postHandle`对于`@ResponseBody`和`ResponseEntity`方法是无效的，因为这个响应在`postHandle`已经完成并提交。

### 语言环境

`DispatcherServlet`会自动使用本地的语言环境去解析消息，这些都交由`LocalResolver`对象来完成。

当接受到一个请求时，`DispatcherServlet`会去查找语言环境解析器（`LocalResolver`），如果成功找到，则使用它来处理语言环境。除了使用自动的语言环境解析，也可以通过定义拦截器进行处理。

### Multipart解析

HTTP协议中规定，在上传文件时，需要添加`multipart/form-data`请求头，在Spring MVC中则通过`MultipartResolver`对其进行处理。Spring MVC中提供了基于[commons-fileupload](https://commons.apache.org/proper/commons-fileupload/)和基于Servlet 3.0的解析实现。

在Spring中声明一个名为`multipartResolver`的bean即可使用其进行处理。当`DispatcherServlet`接受到一个带有`multipart/form-data`的POST请求时，会交由解析器解析并将请求包装为`MultipartHttpServletRequest`以提供对请求中文件的访问方法。

- 基于Apache Commons FileUpload

  通过使用Apache Commons FileUpload进行解析，需要创建一个类型为`CommonsMultipartResolver`并且名为`multipartResolver`的bean。（需要依赖`commons-fileupload`）

- 基于Servlet 3.0

  Servlet 3.0中提供了对于multipart请求的处理，在`web.xml`中添加`<multipart-config>`可配置此功能。

  然后，在Spring中添加一个类型为`StandardServletMultipartResolver`并且名为`multipartResolver`的bean。

## 控制器（Controllers）

在Spring MVC中提供了基于注解的编程模式，被`@Controller`或`@RestController`注解的组件表达了请求映射，请求输入，异常处理等。

通过Spring中的`<context:component-scan/>`和`<mvc:annotation-driven/>`则可自动查找并将它们注册为bean。

```java
@Controller
public class HelloController {

    @GetMapping("/hello")
    public String handle(Model model) {
        model.addAttribute("message", "Hello World!");
        return "index";
    }
}
```

上面的例子中，会使`handle`映射到`/hello`请求上，并且接受一个`Model`参数，然后返回一个视图的名称。

`@RestController`是一个复合注解，包含了`@Controller`和`@ResponseBody`，表明这个类中的所有方法都是直接返回响应，而不需要视图解析器去解析。

### 请求映射

通过使用`@RequestMapping`注解将一个请求映射到控制器中的一个方法上。它提供了多种属性用来匹配路径，HTTP请求方法，请求的参数，请求头等。

以下是指明了请求方法，对于`@RequestMapping`的简写注解：

- `@GetMapping`
- `@PostMapping`
- `@PutMapping`
- `@DeleteMapping`
- `@PatchMapping`

```java
@RestController
@RequestMapping("/persons") //使用在类上则会被类中的方法共有
class PersonController {

    @GetMapping("/{id}") //匹配 /persons/{id} 路径,方法为get
    public Person getPerson(@PathVariable Long id) {
        // ...
    }

    @PostMapping //匹配 /person 路径，方法为post
    @ResponseStatus(HttpStatus.CREATED)
    public void add(@RequestBody Person person) {
        // ...
    }
}
```

### 路径匹配

- `/resouces/ima?e.png`：?用来匹配任意一个字符。
- `/resources/*.png`：*用来匹配0个或多个字符。
- `/resources/**`：匹配多个路径段（如`/resources/a`，`/resources/a/b`）
- `/resources/{project}/version`：匹配一个路径段（用`/`隔开的称为一个路径段），并将其作为参数（这里的project则会作为参数）。
- `/resources/{project:[a-z]+}/version`：通过正则进行匹配。

```java
@GetMapping("/owners/{ownerId}/pets/{petId}") //变量会被自动转换成合适的类型。
public Pet findPet(@PathVariable Long ownerId, @PathVariable Long petId) {
    // ...
}
```

`{varName:regex}`语法表明路径正则来匹配，其中的`varName`作为变量名，`regex`则是具体的正则表达式，例如：匹配路径`/spring-web-3.0.5.jar`这种模式的路径。

```java
@GetMapping("/{name:[a-z-]+}-{version:\\d\\.\\d\\.\\d}{ext:\\.[a-z]+}")
public void handle(@PathVariable String name, 
                   @PathVariable String version, 
                   @PathVariable String ext) {
    // name : spring-web
    // version : 3.0.5
    // ext : .jar
}
```

### Media Type 匹配

- `Content-Type`匹配

  通过请求的`Content-Type`请求头来缩小匹配范围

  ```java
  @PostMapping(path = "/pets", consumes = "application/json") 
  public void addPet(@RequestBody Pet pet) {
      // ...
  }
  ```

- `Accept`匹配

  通过请求中`Accept`请求头中列出的`Content-Type`来缩小匹配范围

  ```java
  @GetMapping(path = "/pets/{petId}", produces = "application/json") 
  @ResponseBody
  public Pet getPet(@PathVariable String petId) {
      // ...
  }
  ```

注：在`MediaType`类中定义了常见的Content-Type

## 处理程序方法

### 方法参数和返回值

对于被`@RequestMapping`注解的方法，Spring可以自动根据其参数类型，传入适当的参数；同样的也能根据其返回值类型自动选择合适的处理。详见[Method Arguments](https://docs.spring.io/spring-framework/docs/current/reference/html/web.html#mvc-ann-methods)

### `@RequestParam`

使用`@RequestParam`注解可以将请求中的参数或者表单中的数据绑定到方法的参数上。

```java
@Controller
@RequestMapping("/pets")
public class EditPetForm {

    @GetMapping
    public String setupForm(@RequestParam("petId") int petId, Model model) { 
        Pet pet = this.clinic.loadPet(petId);
        model.addAttribute("pet", pet);
        return "petForm";
    }
}
```

### `@RequestBody`

 用在方法的参数上，表明将请求体中的数据反序列化成一个对象，并值作为参数传入。（需要结合一个`HttpMessageConverter`）

```java
@PostMapping("/accounts")
public void handle(@RequestBody Account account) {
    // ...
}
```

### HttpEntity

或多或少和使用`@RequestBody`相同，但它能够从中请求头和请求体（请求体为被反序列化后的对象）。

```java
@PostMapping("/accounts")
public void handle(HttpEntity<Account> entity) {
    // ...
}
```

### `@ResponseBody`

用在方法，表明直接将返回值序列化后作为响应体。（同样需要结合一个`HttpMessageConverter`）

```java
@GetMapping("/accounts/{id}")
@ResponseBody
public Account handle() {
    // ...
}
```

### ResponseEntity

和`@ResponseBody`相同，但可以设置状态码和响应头中的部分信息。

```java
@GetMapping("/something")
public ResponseEntity<String> handle() {
    String body = ... ;
    String etag = ... ;
    return ResponseEntity.ok().eTag(etag).build(body);
}
```

## URI

介绍了在Spring中处理URI的一些方式。

### UriComponents

提供了模板式的创建Uri的功能。

```java
UriComponents uriComponents = UriComponentsBuilder
        .fromUriString("https://example.com/hotels/{hotel}") //链接模板
        .queryParam("q", "{q}")  //参数模板
        .encode() //进行编码
        .build(); 
//构建的Uri为https://example.com/hotels/Westion?q=123
URI uri = uriComponents.expand("Westin", "123").toUri();
```

### Uri编码

- 使用`UriComponentsBuilder`中encode方法

  这会对Uri进行预编码，并且在Uri被展开时对变量进行严格编码。

- 使用`UriComponents`中encode方法

  在Uri被展开后进行编码，所以只能在`expand`方法后调用。

注：在路径中`;`是合法的，但在第一种方法时，它会将变量中的`;`进行编码，而第二种方法从不会对其编码。

### 相对于当前请求构建路径

`ServletUriComponentsBuilder`可以相对于当前的请求路径构建一个Uri。

```java
HttpServletRequest request = ...
// Re-uses host, scheme, port, path and query string...
ServletUriComponentsBuilder ucb = ServletUriComponentsBuilder.fromRequest(request)
        .replaceQueryParam("accountId", "{id}").build()
        .expand("123")
        .encode();
```

### 链接到Controllers

Spring MVC提供了一种机制用来构建一个链接到Controllers方法的Uri。

```java
@Controller
@RequestMapping("/hotels/{hotel}")
public class BookingController {

    @GetMapping("/bookings/{booking}")
    public ModelAndView getBooking(@PathVariable Long booking) {
        // ...
    }
}
```

```java
UriComponents uriComponents = MvcUriComponentsBuilder
    .fromMethodName(BookingController.class, "getBooking", 21).buildAndExpand(42);

URI uri = uriComponents.encode().toUri();
```

## 异步请求处理

在Spring 3.0之前，采用Thread-Per-Request的方式处理请求，即每一个请求由一个线程从头到尾处理，当过来一个请求时，从tomcat的线程池中获取一个线程，然后由其执行操作，处理完成后在归还给线程池。但线程池中的线程数量有限，当发生一些IO操作时，会导致线程被长时间占用，从而影响性能。在Spring 3.0中引入了异步请求的支持，即可以由另外一个线程去处理请求。

### 原生Servlet异步请求

```java
@WebServlet(value = "/async", asyncSupported = true) //设置asyncSupported开启异步请求支持
public class AsyncServlet extends HttpServlet {
    @Override
    protected void doGet(HttpServletRequest req,
                         HttpServletResponse resp) throws ServletException, IOException {
        //这个请求开启异步请求
        AsyncContext context = req.startAsync();
        context.start(() -> {
            //获取请求和响应
            ServletRequest request = context.getRequest();
            ServletResponse response = context.getResponse();
            try {
                response.getWriter().println("async");
            } catch (IOException e) {
                e.printStackTrace();
            } finally {
                //结束异步操作
                context.complete();
            }
        });
    }
}
```

### 配置

`web.xml`中，在`DispatcherServlet`下添加`<async-supported>true</async-supported>`即可开启异步请求支持。

在Spring配置中，可以进行详细配置（例如设置异步请求超时）

- 对于Java代码配置，使用`WebMvcConfigurer`中的`configureAsyncSupport`回调方法进行配置。

  ```java
  @Configuration
  public class AppConfig implements WebMvcConfigurer {
      @Override
      public void configureAsyncSupport(AsyncSupportConfigurer configurer) {
          configurer.setDefaultTimeout(100);//单位毫秒
      }
  }
  ```

- 对于XML配置，在`<mvc:annotation-driven/>`下的`<async-support/>`中配置。

  ```xml
  <mvc:annotation-driven>
      <!--如果用了这个标签就必需设置 default-timeout -->
      <mvc:async-support default-timeout="1000" task-executor="executor"/>
  </mvc:annotation-driven>
  <!--配置处理Callable返回值的线程池，默认会使用SimpleAsyncTaskExecutor-->
  <task:executor id="executor" pool-size="20" queue-capacity="25"/>
  ```

### DeferredResult

将方法的返回值使用`DeferredResult`进行包装，即可对其进行异步处理。

```java
@GetMapping("/quotes")
@ResponseBody
public DeferredResult<String> quotes() {
    DeferredResult<String> deferredResult = new DeferredResult<String>();
    // Save the deferredResult somewhere..
    return deferredResult;
}

// 在其它线程中设置值
deferredResult.setResult(result);
```

### Callable

返回值使用`Callable`进行包装，会使用预先配置的线程池来执行这个Callable操作，如果没有配置线程池则会默认使用一个`SimpleAsyncTaskExecutor`线程池，但不推荐使用这个默认线程池。

```java
@PostMapping
public Callable<String> processUpload(final MultipartFile file) {

    return new Callable<String>() {
        public String call() throws Exception {
            // ...
            return "someView";
        }
    };
}
```

### HTTP 流

上面的`DeferredResult`和`Callable`都只能异步地返回单个值，在Spring MVC中也提供了可以异步地返回多个数据的功能。

#### ResponseBodyEmitter

```java
@GetMapping("/events")
public ResponseBodyEmitter handle() {
    ResponseBodyEmitter emitter = new ResponseBodyEmitter();
    
    return emitter;
}

// 在其它线程中发送数据
emitter.send("Hello once");

// 等待一会儿后再次发送
emitter.send("Hello again");

// 数据发送完成
emitter.complete();
```

#### SSE

SSE（Server-Sent Events），是一种服务端到客户端的单向消息推送，基于HTTP协议，并且支持断线重连。

在Spring中使用`SseEmitter`来实现此操作，`SseEmitter`是`ResponseBodyEmitter`的子类，使用方法基本相同。

```java
//响应头里的Content-Type需设置为 text/event-stream
@GetMapping(value = "/sse", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
public SseEmitter sse() {
    SseEmitter emitter = new SseEmitter();
    new Thread(() -> {
        try {
            int number = 0;
            do {
                //每500毫秒发送一次消息，发送20次
                Thread.sleep(500);
                emitter.send("data" + number);
                number++;
            } while (number < 20);
            emitter.send("close");
            //操作完成
            emitter.complete();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }).start();
    return emitter;
}
```

前端接收数据：

```js
let source = new EventSource("http://localhost:8080/mvc/async/sse");
//建立链接时的监听
source.onopen = function (event) {
    console.log(event);
}
//接受到消息的监听
source.onmessage = function (message) {
    let data = message.data;
    console.log(data);
    if (data === "close") {
        source.close();
    }
}
//发生错误的监听
source.onerror = function (event) {
    console.log(event);
}
```

#### 原始数据

可以使用`StreamingResponseBody`发送原始二进行数据，这仍然采用的是流式传输，具体执行也会使用一个线程池，配置和使用`Callable`相同。

```java
@GetMapping("/download")
public StreamingResponseBody handle() {
    return new StreamingResponseBody() {
        @Override
        public void writeTo(OutputStream outputStream) throws IOException {
            // write...
        }
    };
}
```

## 跨域（CORS）

当一个请求URL的协议，域名，端口三者都与当前页面的URL不同时即为跨域。

![image-20210805190341816](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/08/05/20210805190349.png)

出于完全原因，浏览器禁止AJAX请求当前源范围外的资源。在W3C标准中定义了CORS（Cross-Origin Resource Sharing），用来解决跨域问题。

当浏览器发现发送了一个跨域请求时，会在请求头中添加`Origin`字段，用来表明这个请求来自哪个源。服务器则可以根据这个值判断是否允许本次请求。如果允许本次请求，则需要在响应头中设置`Access-Control-Allow-Origin`字段，其值要么是`Origin`中的值，要么是`*`，用来表示接受任意域名的请求。

### `@CrossOrigin`

可用于类或方法上，表示该链接允许跨域请求，默认允许所有源，可通过`origins`来指定特定的源。

```java
@RestController
@RequestMapping("/account")
public class AccountController {

    @CrossOrigin
    @GetMapping("/{id}")
    public Account retrieve(@PathVariable Long id) {
        // ...
    }

    @DeleteMapping("/{id}")
    public void remove(@PathVariable Long id) {
        // ...
    }
}
```

### 全局跨域设置

- 基于JAVA代码配置

  ```java
  @Configuration
  @EnableWebMvc
  public class WebConfig implements WebMvcConfigurer {
  
      @Override
      public void addCorsMappings(CorsRegistry registry) {
  
          registry.addMapping("/api/**")
              .allowedOrigins("https://domain2.com")
              .allowedMethods("PUT", "DELETE")
              .allowedHeaders("header1", "header2", "header3")
              .exposedHeaders("header1", "header2")
              .allowCredentials(true).maxAge(3600);
  
          // Add more mappings...
      }
  }
  ```

- 基于XML配置

  ```xml
  <mvc:cors>
  
      <mvc:mapping path="/api/**"
          allowed-origins="https://domain1.com, https://domain2.com"
          allowed-methods="GET, PUT"
          allowed-headers="header1, header2, header3"
          exposed-headers="header1, header2" allow-credentials="true"
          max-age="123" />
  
      <mvc:mapping path="/resources/**"
          allowed-origins="https://domain1.com" />
  
  </mvc:cors>
  ```

## HTTP 缓存

当客户端请求资源时，会先抵达浏览器缓存，如果在缓存中命中资源，则可以直接从浏览器中获取取而不用去请求服务器，从而减少服务器压力。详见[一文读懂HTTP缓存](https://www.jianshu.com/p/227cee9c8d15)

### CacheControl

`CacheControl`提供了对于`Cache-Control`请求头的相关设置。

```java
// 缓存一小时 - "Cache-Control: max-age=3600"
CacheControl ccCacheOneHour = CacheControl.maxAge(1, TimeUnit.HOURS);

// 不缓存 - "Cache-Control: no-store"
CacheControl ccNoStore = CacheControl.noStore();

// Cache for ten days in public and private caches,
// public caches should not transform the response
// "Cache-Control: max-age=864000, public, no-transform"
CacheControl ccCustom = CacheControl.maxAge(10,TimeUnit.DAYS)
    	.noTransform()
	    .cachePublic();
```

### Controllers

在Conttrollers中设置缓存。

```java
@GetMapping("/book/{id}")
public ResponseEntity<Book> showBook(@PathVariable Long id) {

    Book book = findBook(id);
    String version = book.getVersion();

    return ResponseEntity
            .ok()
            .cacheControl(CacheControl.maxAge(30, TimeUnit.DAYS))
            .eTag(version) // lastModified is also available
            .body(book);
}
```

## 视图处理

在Spring中可以使用多种模板引擎来进行视图处理，并且只需要简单的配置即可使用。

### Thymeleaf

Thymeleaf是一个现代的Java模板引擎，可以直接预览定义的模板并且支持HTML5。详见[thymeleaf](https://www.thymeleaf.org/)

在Spring中使用Thymeleaf详见[thymeleaf+spring](https://www.thymeleaf.org/doc/tutorials/3.0/thymeleafspring.html)

### PDF 和 Excel

Spring还提供展示PDF和Excel表格的相关支持。

注：如果展示Excel需要添加`Apache POI`库，生成PDF则需要添加`OpenPDF`库。

#### PDF

```java
public class PdfWordList extends AbstractPdfView {
    protected void buildPdfDocument(Map<String, Object> model, 
                                    Document doc, 
                                    PdfWriter writer,
                                    HttpServletRequest request, 
                                    HttpServletResponse response) throws Exception {
        List<String> words = (List<String>) model.get("wordList");
        for (String word : words) {
            doc.add(new Paragraph(word));
        }
    }
}
```

## MVC配置

和配置Spring相同，对于MVC的配置，Spring也提供了基于JAVA代码和基于XML的配置。

注：对于基于Java代码配置，需要tomcat也使用基于代码的配置。

### 开启配置

- 基于JAVA代码

  ```java
  @Configuration
  @EnableWebMvc
  public class WebConfig implements WebMvcConfigurer {
      //实现WebMvcConfigurer可以进行更多详细的配置，只需要去实现里面的方法即可
      // Implement configuration methods...
  }
  ```

- 基于XML

  ```xml
  <?xml version="1.0" encoding="UTF-8"?>
  <beans xmlns="http://www.springframework.org/schema/beans"
      xmlns:mvc="http://www.springframework.org/schema/mvc"
      xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
      xsi:schemaLocation="
          http://www.springframework.org/schema/beans
          https://www.springframework.org/schema/beans/spring-beans.xsd
          http://www.springframework.org/schema/mvc
          https://www.springframework.org/schema/mvc/spring-mvc.xsd">
  
      <mvc:annotation-driven/>
  
  </beans>
  ```

### 拦截器

通过注册一系列拦截器对请求进行预处理（关于拦截器在上面有作相关介绍）。

- 基于Java代码

  ```java
  @Configuration
  @EnableWebMvc
  public class WebConfig implements WebMvcConfigurer {
  
      @Override
      public void addInterceptors(InterceptorRegistry registry) {
          registry.addInterceptor(new LocaleChangeInterceptor());
          registry.addInterceptor(new ThemeChangeInterceptor())
              .addPathPatterns("/**")
              .excludePathPatterns("/admin/**");
          registry.addInterceptor(new SecurityInterceptor())
              .addPathPatterns("/secure/*");
      }
  }
  ```

- 基于XML

  ```xml
  <mvc:interceptors>
      <bean class="org.springframework.web.servlet.i18n.LocaleChangeInterceptor"/>
      <mvc:interceptor>
          <mvc:mapping path="/**"/>
          <mvc:exclude-mapping path="/admin/**"/>
          <bean class="org.springframework.web.servlet.theme.ThemeChangeInterceptor"/>
      </mvc:interceptor>
      <mvc:interceptor>
          <mvc:mapping path="/secure/*"/>
          <bean class="org.example.SecurityInterceptor"/>
      </mvc:interceptor>
  </mvc:interceptors>
  ```

### 内容协商

可以配置Spring如何从请求中来确定其Content-Type。Spring会去检查请求的`Accept`请求头，链接扩展名，请求参数等，然后根据配置去选择应该返回何种类型，最后由消息转换器去处理。例如对于同一个Controller，如果请求地址是`/test/account?format=xml`则返回xml格式的数据，如果地址是`/text/account?format=json`，则返回JSON格式的数据（注：这里format会被Spring去检查，而不需要自己处理）。详见[内容协商](https://cloud.tencent.com/developer/article/1497764)

- 基于Java代码

  ```java
  @Configuration
  @EnableWebMvc
  public class WebConfig implements WebMvcConfigurer {
  
      @Override
      public void configureContentNegotiation(ContentNegotiationConfigurer configurer) {
          configurer.mediaType("json", MediaType.APPLICATION_JSON);
          configurer.mediaType("xml", MediaType.APPLICATION_XML);
      }
  }
  ```

- 基于XML

  ```java
  <mvc:annotation-driven content-negotiation-manager="contentNegotiationManager"/>
  
  <bean id="contentNegotiationManager" class="org.springframework.web.accept.ContentNegotiationManagerFactoryBean">
      <property name="mediaTypes">
          <value>
              json=application/json
              xml=application/xml
          </value>
      </property>
  </bean>
  ```

### 消息转换器

消息转换器用来将一种数据格式转换成另一种数据格式，例如将表单数据转换成Java对象，或者将Java对象转换成JSON数据。Spring会根据请求或响应的`Content-Type`自动去选择合适的消息转换器。

在基于Java代码的配置中，可以通过实现方法`configureMessageConverters`去替换掉Spring中默认的消息转换器，通过实现方法`extendMessageConverters`添加额外的消息转换器。

下面的例子将注册Jaskson中处理XML和JSON的消息转换器。注：Spring默认会使用Jackson来处理JSON数据（JSON转Java对象或Java对象转JSON），所以如果只需要处理JSON，则只需要添加Jackson依赖即可，无需额外配置，但如果需要对Jackson进行配置的话则需要显式注册。

- 基于Java代码

  ```java
  @Configuration
  @EnableWebMvc
  public class WebConfiguration implements WebMvcConfigurer {
  
      @Override
      public void configureMessageConverters(List<HttpMessageConverter<?>> converters) {
          Jackson2ObjectMapperBuilder builder = new Jackson2ObjectMapperBuilder()
                  .indentOutput(true)
                  .dateFormat(new SimpleDateFormat("yyyy-MM-dd"))
                  .modulesToInstall(new ParameterNamesModule());
          converters.add(new MappingJackson2HttpMessageConverter(builder.build()));
          converters.add(
              new MappingJackson2XmlHttpMessageConverter(builder.createXmlMapper(true)
                                                                .build()));
      }
  }
  ```

- 基于XML

  ```xml
  <mvc:annotation-driven>
      <mvc:message-converters>
          <bean class="org.springframework.http.converter.json.MappingJackson2HttpMessageConverter">
              <property name="objectMapper" ref="objectMapper"/>
          </bean>
          <bean class="org.springframework.http.converter.xml.MappingJackson2XmlHttpMessageConverter">
              <property name="objectMapper" ref="xmlMapper"/>
          </bean>
      </mvc:message-converters>
  </mvc:annotation-driven>
  
  <bean id="objectMapper" class="org.springframework.http.converter.json.Jackson2ObjectMapperFactoryBean"
        p:indentOutput="true"
        p:simpleDateFormat="yyyy-MM-dd"
        p:modulesToInstall="com.fasterxml.jackson.module.paramnames.ParameterNamesModule"/>
  
  <bean id="xmlMapper" parent="objectMapper" p:createXmlMapper="true"/>
  ```

#### 补充1: 使用FastJson来处理JSON

当使用FastJson来处理JSON时，则需要注册FastJson中实现的消息转换器，以让Spring能够使用其来处理JSON数据。

- 基于Java代码

  ```java
  @Configuration
  @EnableWebMvc
  public class AppConfig implements WebMvcConfigurer {
      @Override
      public void configureMessageConverters(List<HttpMessageConverter<?>> converters) {
          FastJsonConfig config = new FastJsonConfig();
          config.setCharset(StandardCharsets.UTF_8);
          config.setDateFormat("yyyy-MM-dd");
          FastJsonHttpMessageConverter converter = new FastJsonHttpMessageConverter();
          converter.setSupportedMediaTypes(List.of(MediaType.APPLICATION_JSON));
          converter.setFastJsonConfig(config);
          converters.add(converter);
      }
  }
  ```

- 基于XML

  ```xml
  <mvc:annotation-driven>
      <mvc:message-converters>
          <bean 
               class="com.alibaba.fastjson.support.spring.FastJsonHttpMessageConverter">
              <!--此属性也可以不设置，但响应中的Content-Type会是text/html,
              而不是application/json，不过也可以通过在Controller中设置produces属性来解决
              -->
              <property name="supportedMediaTypes">
                  <list>
                      <value>application/json</value>
                  </list>
              </property>
              <property name="fastJsonConfig" ref="fastConfig"/>
          </bean>
      </mvc:message-converters>
  </mvc:annotation-driven>
  <bean id="fastConfig" class="com.alibaba.fastjson.support.config.FastJsonConfig">
      <!--设置字符集-->
      <property name="charset" value="UTF-8"/>
      <!--设置日期格式-->
      <property name="dateFormat" value="yyyy-MM-dd"/>
  </bean>
  ```

#### 补充2:自定义消息转换器

可以自定义消息转换器来支持自定义的`Content-Type`或者自己实现消息转换。

自定义消息转换器只需要创建一个类并实现`HttpMessageConverter`接口或者继承`AbsAbstractHttpMessageConverter`类。最后模仿上面的方法进行注册即可。

下面实现了一个简单的例子，对自定义的`Content-Type`--`application/demo`进行处理。

```java
//这里泛型指定处理哪些类，这里简单起见，只处理Pojo类
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

更多细节可以查看Fastjson中消息转换器的实现。

### 静态资源

用于处理中的静态资源，例如下面的例子中设置了对于以`/resources`开头的请求，在`/public`和类路径下的`/static`中查找相关资源，同时设置过期时间为1年。

- 基于Java代码

  ```java
  @Configuration
  @EnableWebMvc
  public class WebConfig implements WebMvcConfigurer {
  
      @Override
      public void addResourceHandlers(ResourceHandlerRegistry registry) {
          registry.addResourceHandler("/resources/**")
              .addResourceLocations("/public", "classpath:/static/")
              .setCacheControl(CacheControl.maxAge(Duration.ofDays(365)));
      }
  }
  ```

- 基于XML

  ```xml
  <mvc:resources mapping="/resources/**"
      location="/public, classpath:/static/"
      cache-period="31556926" />
  ```

  这里`mapping`的值使用的ant路径风格，`?`匹配任意单个字符，`*`匹配0到多个字符，`**`可以匹配多级目录。

### 默认Servlet

在Tomcat中有一个默认servlet，名称为`default`，主要用于处理静态资源请求。

但使用Spring MVC时，将`DispatcherServlet`映射到`/`路径下后，对于静态资源的请求也将会发送到`DispatcherServlet`从而导致无法正确获取到静态资源，不过可以由Spring创建一个`DefaultServletHttpRequestHandler`，它会检查所有请求，如果是请求静态资源则将其发送到Tomcat的默认Servlet中去处理。

- 基于Java代码

  ```java
  @Configuration
  @EnableWebMvc
  public class WebConfig implements WebMvcConfigurer {
  
      @Override
      public void configureDefaultServletHandling(DefaultServletHandlerConfigurer configurer) {
          configurer.enable();
      }
  }
  ```

- 基于XML

  ```xml
  <mvc:default-servlet-handler/>
  ```

#### 补充：解决静态资源无法访问的问题

在上面也有写到，对静态资源的请求也会被发送到`DispatchrServlet`中，所以静态资源无法被正确获取，主要有三种方法来解决。

- 使用Tomcat中默认Servlet

  ```xml
  <servlet-mapping>
      <servlet-name>default</servlet-name>
      <url-pattern>*.html</url-pattern>
  </servlet-mapping>
  <servlet-mapping>
      <servlet-name>default</servlet-name>
      <url-pattern>*.css</url-pattern>
  </servlet-mapping>
  ```

  这是tomcat中默认创建的Servlet，名称即为default，可以直接使用。

- 由Spring将静态资源转发到默认Servlet，参见上面默认Servlet部分。

  ```xml
  <mvc:default-servlet-handler/>
  ```

- 由Spring处理静态资源，参见上面静态资源部分。

  ```xml
  <mvc:resources mapping="/images/**" location="/images/"/>
  ```

三种方法可以混用，但从实现原理来看，第一种方式应单独使用。

## RestTemplate

Spring提供子一个更简洁的用于发送HTTP请求的`RestTemplate`，详见[RestTemplate](https://docs.spring.io/spring-framework/docs/current/reference/html/integration.html#rest-resttemplate)。其默认使用了JDK中的`HttpURLConnection`来进行发送HTTP请求，但也可以使用其它的库，包括：

- Apache HttpComponents
- Netty
- OkHttp

这里使用HttpComponents作为HTTP请求库。

添加maven依赖：

```xml
<dependency>
    <groupId>org.apache.httpcomponents</groupId>
    <artifactId>httpclient</artifactId>
    <version>4.5.13</version>
</dependency>
```

```java
RestTemplate template = new RestTemplate(new HttpComponentsClientHttpRequestFactory());
```

### URI

```java
String result = restTemplate.getForObject(
        "https://example.com/hotels/{hotel}/bookings/{booking}", 
    	String.class, "42", "21");
```

```java
Map<String, String> vars = Collections.singletonMap("hotel", "42");
//使用map保存参数
String result = restTemplate.getForObject(
        "https://example.com/hotels/{hotel}/rooms/{hotel}", String.class, vars);
```

### Headers

```java
String uriTemplate = "https://example.com/hotels/{hotel}";
URI uri = UriComponentsBuilder.fromUriString(uriTemplate).build(42);

RequestEntity<Void> requestEntity = RequestEntity.get(uri)
        .header("MyRequestHeader", "MyValue")
        .build();
//使用exchange使用指定的请求头
ResponseEntity<String> response = template.exchange(requestEntity, String.class);

String responseHeader = response.getHeaders().getFirst("MyResponseHeader");
String body = response.getBody();
```

## WebSocket

WebSocket是一种通信协议，提供了一种建立一个全双工，通过一个TCP链接双向交流的通信方式。HTTP协议中，只能客户端向服务端发起通信，而在WebSocket中，既可以由客户端发起通信，也可以由服务端发起。其具有以下特点：

- 建立在TCP协议之上
- 与HTTP有良好的兼容性，默认端口也是80和443
- 数据格式轻量，性能开销小。
- 没有同源限制。
- 协议标识符为`ws`（如果加密则是`wss`）。

### 何时使用WebSocket

虽然WebSocket可以让应用有更好的动态性和可交互性，但在多数场景下，使用轮询和长轮询方式也能达到同样的效果并且使用更简单。例如：新闻，邮件等信息需要动态的更新，但每隔几分钟更新一次或许更好，而社交信息，游戏等就需要更实时的更新。因此，相比之下，WebSocket更适合用在低延迟（这里的延迟不是指网络延迟），高频次，大信息量的场景下。

### WebSocketHandler

在Spring中，创建一个WebSocket服务中需要实现`WebSocketeHandler`接口即可（或者继承其己有实现类，例如：`TextWebSocketHandler`，`BinaryWebSocketHandler`）。

maven依赖：

```xml
<dependency>
    <groupId>org.springframework</groupId>
    <artifactId>spring-websocket</artifactId>
    <version>5.3.9</version>
</dependency>
```

```java
import org.springframework.web.socket.WebSocketHandler;
import org.springframework.web.socket.WebSocketSession;
import org.springframework.web.socket.TextMessage;

public class MyHandler extends TextWebSocketHandler {

    @Override
    public void handleTextMessage(WebSocketSession session, TextMessage message) {
        // ...
    }

}
```

然后在配置中注册这个Handler

- 基于java代码

  ```java
  import org.springframework.web.socket.config.annotation.EnableWebSocket;
  import org.springframework.web.socket.config.annotation.WebSocketConfigurer;
  import org.springframework.web.socket.config.annotation.WebSocketHandlerRegistry;
  
  @Configuration
  @EnableWebSocket
  public class WebSocketConfig implements WebSocketConfigurer {
  
      @Override
      public void registerWebSocketHandlers(WebSocketHandlerRegistry registry) {
          registry.addHandler(myHandler(), "/myHandler");
      }
  
      @Bean
      public WebSocketHandler myHandler() {
          return new MyHandler();
      }
  
  }
  ```

- 基于XML

  ```xml
  <beans xmlns="http://www.springframework.org/schema/beans"
      xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
      xmlns:websocket="http://www.springframework.org/schema/websocket"
      xsi:schemaLocation="
          http://www.springframework.org/schema/beans
          https://www.springframework.org/schema/beans/spring-beans.xsd
          http://www.springframework.org/schema/websocket
          https://www.springframework.org/schema/websocket/spring-websocket.xsd">
  
      <websocket:handlers>
          <websocket:mapping path="/myHandler" handler="myHandler"/>
      </websocket:handlers>
  
      <bean id="myHandler" class="org.springframework.samples.MyHandler"/>
  
  </beans>
  ```

前端代码：

```js
let client = new WebSocket("ws://localhost:8080/mvc/socket");
//建立连接事件监听
client.onopen = function (event) {
    console.log(event);
    //发送消息
    client.send("hello");
}
//接收到消息事件监听
client.onmessage = function (messageEvent) {
    console.log(messageEvent);
    let data = messageEvent.data;
    if (data === "close") {
        client.close();
    }
}
//关闭连接事件监听
client.onclose = function (closeEvent) {
    console.log(closeEvent);
}
//发生错误监听
client.onerror = function (event) {
    console.log(event);
    client.close();
}
```

### WebSocket握手

因为WebSocket只进行一次握手，并且之后一直保持连接，所以在某些场景下（例如验证是否登陆）则需要对握手过程进行拦截。

最简单地自定义建立HTTP进行WebSocket握手请求（WebSocket的握手是用HTTP协议进行）的方式是通过实现`HandshakeInterceptor`。

这里使用Spring中内置的`HttpSessionHandshakeInterceptor`，只需要进行配置即可。

```xml
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:websocket="http://www.springframework.org/schema/websocket"
    xsi:schemaLocation="
        http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd
        http://www.springframework.org/schema/websocket
        https://www.springframework.org/schema/websocket/spring-websocket.xsd">

    <websocket:handlers>
        <websocket:mapping path="/myHandler" handler="myHandler"/>
        <websocket:handshake-interceptors>
            <bean class="org.springframework.web.socket.server
                         .support.HttpSessionHandshakeInterceptor"/>
        </websocket:handshake-interceptors>
    </websocket:handlers>

    <bean id="myHandler" class="org.springframework.samples.MyHandler"/>

</beans>
```

### 服务端配置

可以通过`ServletServerContainerFactoryBean`对WebSocket进行一些配置。

```xml
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:websocket="http://www.springframework.org/schema/websocket"
    xsi:schemaLocation="
        http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd
        http://www.springframework.org/schema/websocket
        https://www.springframework.org/schema/websocket/spring-websocket.xsd">
	<!--只需要创建这个bean即可-->
    <bean class="org.springframework...ServletServerContainerFactoryBean">
        <property name="maxTextMessageBufferSize" value="8192"/>
        <property name="maxBinaryMessageBufferSize" value="8192"/>
    </bean>

</beans>
```

### 跨域允许

从Spring 4.1.5开始，对于WebSocket的默认行为是只接受同源请求，所以显式设置能够的接受其它源。

在这里有三种行为：

- 只允许同源请求（默认）：仅允许同源请求。
- 允许指定的源：每一个被指定允许的源都需要以`http://`或`https://`开头。
- 允许所有的源：使用`*`表示允许所有源。

```xml
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:websocket="http://www.springframework.org/schema/websocket"
    xsi:schemaLocation="
        http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd
        http://www.springframework.org/schema/websocket
        https://www.springframework.org/schema/websocket/spring-websocket.xsd">

    <websocket:handlers allowed-origins="https://mydomain.com">
        <websocket:mapping path="/myHandler" handler="myHandler" />
    </websocket:handlers>

    <bean id="myHandler" class="org.springframework.samples.MyHandler"/>

</beans>
```

## SockJS

在某些时候，可以因为各种原因而无法使用WebSocket（比如浏览器不支持，服务端的Web不支持），而最佳解决方案则是WebSocket模拟，在支持WebSocket时，直接使用WebSocket；在不支持时，则使用Http协议进行模拟。而SockJS的目标则是在可以使用WebSocket的情形下使用，在不可用时采用模拟的方法进行WebSocket通信。

### 在Spring中启用SockJS

```xml
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:websocket="http://www.springframework.org/schema/websocket"
    xsi:schemaLocation="
        http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd
        http://www.springframework.org/schema/websocket
        https://www.springframework.org/schema/websocket/spring-websocket.xsd">

    <websocket:handlers>
        <websocket:mapping path="/myHandler" handler="myHandler"/>
        <!--加上这一句就可以了-->
        <websocket:sockjs/>
    </websocket:handlers>

    <bean id="myHandler" class="org.springframework.samples.MyHandler"/>

</beans>
```

对于客户端而言，SockJS也提供了相应的实现，详见[SockJS Client](https://github.com/sockjs/sockjs-client/)

## 关于`ContextLoaderListener`

在`web.xml`中有这样一段配置

```xml
<listener>
    <listener-class>org.springframework.web.context.ContextLoaderListener</listener-class>
</listener>
<context-param>
    <param-name>contextConfigLocation</param-name>
    <param-value>classpath:spring.xml</param-value>
</context-param>
```

这里定义了一个监听器，其作用是在启动web容器时加载Spring中`ApplicationContext`（在web应用中应该是`WebApplicationContext`）的配置信息。这个监听器实现的是`ServletContextListener`，因此可以监听到web容器的启动事件，而对于配置文件的位置则在`<context-param>`中声明，如果没有声明则默认查找`/WEB-INF/applicationContext.xml`文件。









