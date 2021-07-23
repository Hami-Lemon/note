# Spring Core

Spring是分层的java全栈是的轻量级开源框架，以IOC(控制反转)和AOP(面向切面编程)为内核，提供了展现层SpringMVC和持久层Sprng JDBC以及业务层事务管理等众多的企业级应用技术。

## 体系结构

![1](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115142.png)

## IOC(控制反转)

把对象创建的权力交给框架,可以削减程序的耦合。

在传统的开以方式中，创建一个对象最直接的方式是使用`new`关键字，并且在创建这个对象时，可能这个对象又依赖于其它对象，为了能够将其创建出来还需要解决其中的依赖关系，而使用Spring容器，则是将对象的创建交给容器去解决，同样的依赖关系也将由容器去解决，从而实现控制反转。

但Ioc（Inversion of Control）并不能让人更加直观和清晰地理解其背后所代表的含义，于Martin Fowler创造了一个新词 —— 依赖注入（Dependency Injection, DI）。

- 使用new创建对象
  ![2](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115147.png)

- 使用工厂创建对象
  ![3](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115151.png)

### spring中的IOC

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd">
    <!-- 把对象创建交给spring管理-->
    <bean id="userService" class="com.hamilemon.service.impl.IUserServiceImpl"/>
    <bean id="userDao" class="com.hamilemon.dao.impl.IUserDaoImpl"/>
</beans>
```

使用配置文件,让spring的ioc获取对象

```java
//获取spring的ioc核心容器，并根据id获取对象
//1.获取核心容器对象
ApplicationContext ac = new ClassPathXmlApplicationContext("bean.xml");
//2.根据id获取bean对象
IUserService us = (IUserService) ac.getBean("userService");
IUserDao ud = ac.getBean("userDao",IUserDao.class);

```

- AplicationContext三个常用实现类:
    1. ClassPathXmlApplicationContext:可以加载类路径下的配置文件
    2. FileSysteamXmlApplicationContext: 可以加载磁盘任意路径下的配置文件
    3. AnnotationConfigApplicationContext: 用于读取注解创建容器的

- 核心容器的两个接口:
    1. ApplicationContext:适用于单例对象(常用)
    在构造核心容器时,创建对象采取立即加载的方式,也就是,一读完配置文件就马上创建对象
    2. BeanFactory:适用于多例对象
    采用延迟加载方式,什么时候获取对象,什么时候创建对象
- 创建bean的三种方式
    1. 使用默认构造函数,配置文件中使用bean标签,只有id和class属性

       ```xml
       <bean id="userService" class="com.hamilemon.service.impl.IUserServiceImpl"/>
       ```

    2. 使用工厂类中的方法创建对象,并存入spring容器

       ```xml
       <!--使用instanceFactory工厂创建userService对象-->
       <bean id="instanceFactory" class="xxx"/>
       <bean id="userService" factory-bean="instanceFactory" factory-method="getUserService/>
       ```

    3. 使用静态工厂中的静态方法创建对象,并存入spring容器

       ```xml
       <bean id="userService" class="com.hamilemon.factory.StaticFactory"
       factory-method="getUserService">
       ```

- bean的作用范围
spring创建的bean对象默认为单例，使用scope指定bean的作用范围

1. singleton:单例(默认值)

2. prototype:多例，每次获取时都会创建一个新对象

3. request:作用于web应用的请求返回

4. session:作用web应用的作用返回

5. global-session: 作用集群环境的会话范围

   ```xml
   <bean id="xxx" class="xxx" scope="singleton"/>
   ```

- bean对象生命周期
    1. 单例对象
    创建容器时出生,销毁容器时销毁,和容器相同

    ```xml
    <bean id="xxx" class="xxx" scope="singleton"
    init-method="xxx" destroy-method="xxx"/>
    ```

    2. 多例对象
       使用对象时才创建,由GC回收

### 依赖注入

依赖关系都交给spring来维护,在当前类需要用到其它类的对象,由spring提供,只需要在配置文件中说明即可，依赖关系的维护称为依赖注入。

#### 能注入的数据

- 基本类型和String

- 其它bean类型(配置文件中或者注解配置过的bean)

- 复杂类型/集合类型

#### 使用构造函数注入

  bean中使用constructor-arg

- type：指定要注入的数据的数据类型

- index:指定要注入的数据给构造函数中指定索引位置的参数赋值（0开始）

- name: 指定给构造函数中指定名称的参数赋值

- ref:用于指定其它的bean类型

- value:提供基本类型和String

```xml
<bean id="accountService" class="com.hamilemon.service.impl.AccountServiceImpl">
    <constructor-arg name="name" value="test"/>
    <constructor-arg name="age" value="15"/>
    <constructor-arg name="birthday" ref="now"/>
</bean>
<bean id="now" class="java.util.Date"/>
```

#### set方法注入

bean中使用property标签

- name: 注入时所调用的set方法名
- ref:用于指定其它的bean类型
- value:提供基本类型和String

```xml
<bean id="now" class="java.util.Date"/>

<bean id="accountService" class="com.hamilemon.service.implAccountServiceImpl">
    <property name="age" value="20"/>
    <property name="name" value="塔斯特"/>
    <property name="birthday" ref="now"/>
</bean>
```

#### 复杂类型注入

用于给list注入的标签:`list` `array` `set`
用于给map注入的标签: `map` `props`

```xml
<bean id="accountService" class="xxx">
    <property name="xxx">
        <arraly>
            <value>AAA</value>
            <value>BBB</value>
        </array>
    </property>

    <property name="xxx">
        <list>
            <value>AAA</value>
            <value>BBB</value>
        </list>
    </property>

    <property name="xxx">
        <map>
            <entry key="xxx" value="xxx"/>
            <entry key="xxx">
                <value>xxx</value>
            </entry>
        </map>
    </property>
</bean>
```

### 基于注解的IOC

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:context="http://www.springframework.org/schema/context"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
       http://www.springframework.org/schema/beans/spring-beans.xsd
       http://www.springframework.org/schema/context
       https://www.springframework.org/schema/context/spring-context.xsd">
    <!--告知spring在创建容器时要扫描的包-->
    <context:component-scan base-package="com.hamilemon"/>
</beans>
```

- 用于创建对象的
    1. `@Compoent`:用于把当前类对象存入spring容器中,作用在类上
    属性:value:用于指定bean的id,默认为当前类名的首字母小写
    2. `@Controller`:表现层
    3. `@Service`:业务层
    4. `@Repository`:持久层
    作用和属性于`@Compoent`相同,是spring提供明确的三层使用的注解
- 用于注入数据的
    (集合类型的注入只能通过xml实现)
    1. `Autowired`:可以作用在成员和方法上,自动按照类型注入,只要容器中有唯一的一个bean对象类型和要注入的变量类型匹配,就可以注入成功;
    如果容器中,没有任何bean类型和要注入的变量类型匹配,则报错;如果有多个,先匹配类型,再匹配名称,如果两个都不匹配,则报错
    2. `@Qualifier`:在按照类中注入的基础上,在按照名称注入.在给类成员注入时不能单独使用需要配合`@Autowired`使用,在给方法参数注入时可以
    属性: value:用于指定注入bean的id
    3. `@Resource`:直接按照bean的id注入,可以独立使用(需要导入javax.annotation)
    属性:name:用于指定bean的id
    4. `@Value`:用于注入基本类型和String类型
       属性:value:用于指定数据的值,可以使用spring中的SpEL(Spring的EL表达式)
        SpEL写法:`${表达式}`
- 用于改变作用范围的

    1. `@Scope`:用于指定bean的作用范围
        属性:value:指定范围的取值(singleton,prototype)
- 和生命周期相关

    1. `@PreDestroy`:用于指定销毁方法
    2. `@PostConstruct`:用于指定初始化方法

- 指定配置类
创建一个类加上`@Configuration`
`@ComponentScan`:通过注解指定spring在创建容器时要扫描的包
属性:value:和basePackages的作用是一样的,都是指定创建容器时要扫描的包
`@ComponentScans`:配置多个`@ComponentScan`

`@Bean`:用于把当前方法的返回值作文bean对象存入spring的IOC容器中
属性:name:用于指定bean的id,默认为当前方法的名称
使用注解配置方法时,如果方法有参数,spring会去容器中查找,查找方式和`@Autowired`相同

```java
/**一个配置类，用来替代bean.xml
 * @author Hami Lemon
 */

@Configuration
@ComponentScan("com.hamilemon")
public class SpringConfiguration {

    /**
     * 创建QueryRunner对象,并存入IOC容器
     * @param dataSource
     * @return
     */
    @Bean(name="runner")
    public QueryRunner createQueryRunner(DataSource dataSource){
        return new QueryRunner(dataSource);
    }
    @Bean(name="dataSource")
    public DataSource createDataSource(){
        ComboPooledDataSource ds = new ComboPooledDataSource();
        try {
            ds.setDriverClass("com.mysql.cj.jdbc.Driver");
        } catch (PropertyVetoException e) {
            e.printStackTrace();
        }
        ds.setJdbcUrl("jdbc:mysql://localhost:3306/spring?serverTimezone=Asia/Shanghai");
        ds.setUser("root");
        ds.setPassword("root");
        return ds;
    }
}
```

使用

```java
ApplicationContext ac =
new AnnotationConfigApplicationContext(SpringConfiguration.class);
```

- `@Import`:用于导入其它配置类
  属性:value:用于指定其它配置类的字节码
  `@Import(JdbcConfig.class)`

- `@PropertySource`:指定配置文件的位置
  在使用的地方用`@Value(${key})`的方式注入
  属性:value:指定文件的名称和路径,关键字`classpath`表示在类路径下
  `@PropertySource("classpath:jdbcConfig.properties")`
  `@PropertySources`:指定多个

### spring整合junit

1. 添加spring整合junit的依赖

   ```xml
   <dependency>
       <groupId>org.springframework</groupId>
       <artifactId>spring-test</artifactId>
       <version>5.2.5.RELEASE</version>
   </dependency>
   ```

2. 使用`@Runwith(SpringJUnit4ClassRunner.class)`把原本的main方法替换成sprig提供的

3. 告知spring运行器,spring和ioc创建是xml还是注解
   `@ContextConfiguration`
   属性:locations:指定xml的位置,加上classpath,表示在类路径下
    classes:注解类的位置

4. 需要spring提供的对象上添加`@Autowired`自动注入

## AOP面向切面编程

### 动态代理

字节码随用随创建，随用随加载，不修改源码的基础上对方法增强

- 基于接口的动态代理
    使用`Proxy`对象中的`newProxyInstance`方法。
    
    - 方法参数:
          1. `ClassLoader`:类加载器,用于在家代理对象的字节码,传入被代理对象的类加载器 `被代理类.getClass().getClassLoader()`
             2. `Class[]`:字节码数组,用于让代理对象和被代理对象有相同方法, 传入`被代理类.getClass().getInterfaces()`
             3. `InvocationHandler`:用于提供增强的代码,让我们写如何代理， 一般传入一个该接口的实现类
    
    ```java
    Proxy.newProxyInstance(producer.getClass().getClassLoader(),
        producer.getClass().getInterfaces(),
        new InvocationHandler(){
    
            /**执行被代理对象的任何接口方法都会经过该方法
            * @param proxy 代理对象的引用
            * @param method 当前执行的方法
            * @param args 当前执行方法所需的参数
            * @return 和被代理对象方法有相同返回值
            */
            @Override
            public Object invoke(Object proxy, Method method, Object[] args) throws Throwable{
                //提供增强的代码
                Object returnValue = null;
                1. 获取方法执行的参数
                Float money = (Float)args[0];
                2. 判断当前方法是不是销售方法
                if("saleProduct" .equals(method.getName())){
                    //producer:被代理的对象
                    returnValue = method.invoke(producer,money * 0.8f);
                }
    
                return returnValue;
            }
        });
    ```
    
- 基于子类的动态代理

  使用[cglib](https://blog.csdn.net/danchu/article/details/70238002)实现

  ```xml
  <dependency>
      <groupId>cglib</groupId>
      <artifactId>cglib</artifactId>
      <version></version>
  </dependency>
  ```

  使用`cglib`库的`Enhancer`中的`create`方法，要求:被代理类不能时最终类（非final）

  方法参数:

  	1. Class:字节码,指定被代理对象的字节码
   	2. Callback:用于提供增强的代码,一般传入MethodInterceptor的实现类

```java
Enhancer.create(producer.getClass, new MethodInterceptor(){
    /**
    执行被代理对象任何方法都会执行该方法,前三个参数和invoke方法参数含义相同
    methodProxy:当前执行方法的代理对象
    */
    @Override
    public Object intercept(Object proxy, Method method, Object[] args, MethodProxy methodProxy) throws Throwable{
            //提供增强的代码
            Object returnValue = null;
            1. 获取方法执行的参数
            Float money = (Float)args[0];
            2. 判断当前方法是不是销售方法
            if("saleProduct" .equals(method.getName())){
                //producer:被代理的对象
                returnValue = method.invoke(producer,money * 0.8f);
            }
            return returnValue;
        }
    });
```

### AOP(面向切面编程)

通过预编译方式和运行期动态代理实现程序功能的统一维护的一种技术,spring中通过配置实现动态代理

- 相关术语
  - Joinpoint(连接点):指那些被拦截到的点,在spring中,这些点指的是方法,因为spring只支持方法类型的连接点

  - Pointcout(切入点):指我们要对哪些Joinpoint进行拦截的定义

  - Advice(通知/增强):拦截到joinpoint之后要做的事情
        通知的类型:前置通知,后置通知,异常通知,最终通知,环绕通知
    ![4](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115206.png)

  - Introduction(引介):一种特殊的通知,在不修改类代码的前提下,Introduction可以在运行期为类动态的添加一些方法或Field
  - Target(目标对象):代理的目标对象

  - Weaving(织入):是指把增强应用到目标对象来创建新的代理对象的过程

  - Proxy(代理):一个类被AOP织入增强后,就产生一个结果代理类

  - Aspect(切面):是切入点和通知的结合

### 基于xml的aop

xml约束

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:aop="http://www.springframework.org/schema/aop"
    xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd
        http://www.springframework.org/schema/aop
        https://www.springframework.org/schema/aop/spring-aop.xsd">
```

配置spring的IOC

```xml
<bean id="accountService" class="com.hamilemon.service.impl.AccountServiceImpl"/>
```

配置AOP

1. 把通知的bean交给spring管理
2. 使用`aop:config`标签表明开始AOP的配置r
3. 使用`aop:aspect`表明开始配置切面
    id属性:提供一个唯一标识
    ref属性:指定通知类的id
4. 在`aop:aspect`中配置通知类型
    `aop:before`:前置通知
        `method`属性:用于指定哪个方法是前置通知
        `pointcut`属性:用于指定切入点表达式,表示对业务层中哪些方法增强
    `aop:after-returning`:后置通知
    `aop:after-throwing`:异常通知
    `aop:after`:最终通知
    `aop:around`:环绕通知
    ![5](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115216.png)
    ![6](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115221.png)
    - 切入点表达式写法
        写法:execution(表达式)
        表达式:访问修饰符 返回值 包名.类名.方法名(参数列表)
        `public void com.hamilemon.service.impl.AccountServiceImpl.saveAccount()`
        全通配写法:
        `* *..*.*(..)`
        访问修饰符可以省略,返回值、包名、类名、方法名可使用通配符`*`;
        包名可以使用`..`表示当前包及其子包
        参数列表:基本类型直接写,引用类型`包名.类名`,可使用`*`表示任意类型
        一般写成:`* com.hamilemon.service.impl.*.*(..)`
    - 配置切入点表达式
        id用于指定表达式的唯一标识,expression用于指定表达式内容
        此标签写在`aop:aspect`内部只能当前切面可用
        写在`aop:aspect`外面则所有切面都可用(只能放在`aop:aspect`前面)

        ```xml
        <aop:pointcut id="pt1" expression="execution(* com.hamilemon.service.impl.*.*(..))"/>
        ```

```xml
    <!--配置spring的IOC,把service对象配置进来-->
<bean id="accountService" class="com.hamilemon.service.impl.AccountServiceImpl"/>
    <!--    基于xml的aop配置步骤-->
    <!--    配置Logger类-->
<bean id="logger" class="com.hamilemon.utils.Logger"/>
    <!--    配置AOP-->
<aop:config>
    <!--配置切面-->
    <aop:aspect id="logAdvice" ref="logger">
        <!--前置通知-->
        <aop:before method="printLog"
                    pointcut="execution(* com.hamilemon.service.impl.*.*(..))"/>
        <!--后置通知-->
        <aop:after-returning method="after"
                    pointcut="execution(* com.hamilemon.service.impl.*.*(..))"/>
        <!--异常通知-->
        <aop:after-throwing method="exception"
                            pointcut="execution(* com.hamilemon.service.impl.*.*(..))"/>
        <!--最终通知-->
        <aop:after method="finallyAd"
                    pointcut="execution(* com.hamilemon.service.impl.*.*(..))"/>
        <!--环绕通知-->
        <aop:around method="around()"
                    pointcut-ref="pt">

        <aop:pointcut id="pt" expression="execution(* com.hamilemon.service.impl.*.*(..))"/>
    </aop:aspect>
</aop:config>
```

### 基于注解的AOP

调用顺序会有一定的问题,建议使用环绕通知

```java
@Component("logger")
@Aspect//当前类是一个切面
public class Logger {
    @Pointcut("execution(* com.hamilemon.service.impl.*.*(..))")
    private void pt(){}
    /**
     * 前置通知
     */
    @Before("pt()")
    public void printLog(){
        System.out.println("开始记录");
    }
    /**
     * 后置通知
     */
    @AfterReturning("pt()")
    public void after(){
        System.out.println("后置");
    }
    /**
     * 异常通知
     */
    @AfterThrowing("pt()")
    public void exception(){
        System.out.println("异常");
    }
    /**
     * 最终通知
     */
    @After("pt()")
    public void finallyAd(){
        System.out.println("最终");
    }
    /**
    * 环绕通知
    */
//    @Around("pt()")
    public void around(){

    }
}
```

## spring中的jdbcTemplate

用于和数据库交互，实现对表的crud操作

- 依赖关系

```xml
<dependency>
    <groupId>org.springframework</groupId>
    <artifactId>spring-context</artifactId>
    <version>5.2.0.RELEASE</version>
</dependency>

<dependency>
    <groupId>org.springframework</groupId>
    <artifactId>spring-jdbc</artifactId>
    <version>5.2.5.RELEASE</version>
</dependency>

<dependency>
    <groupId>org.springframework</groupId>
    <artifactId>spring-tx</artifactId>
    <version>5.2.5.RELEASE</version>
</dependency>
```

- 声明式事务控制
spring提供了一组事务控制的接口，在`spring-tx`依赖中

spring的事务控制都是基于AOP的,既可以使用编程的方式实现,也可以使用配置的方式实现

1. 配置事务管理器

   ```xml
   <bean id="transactionManager" class="org.springframework.jdbc.datasource.DataSourceTransactionManager">
           <property name="dataSource" ref="dataSource"/>
   </bean>
   ```

2. 导入事务的约束,配置事务通知

   ```xml
   <?xml version="1.0" encoding="UTF-8"?>
   <beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:aop="http://www.springframework.org/schema/aop"
       xmlns:tx="http://www.springframework.org/schema/tx"
       xsi:schemaLocation="
           http://www.springframework.org/schema/beans
           https://www.springframework.org/schema/beans/spring-beans.xsd
           http://www.springframework.org/schema/tx
           https://www.springframework.org/schema/tx/spring-tx.xsd
           http://www.springframework.org/schema/aop
           https://www.springframework.org/schema/aop/spring-aop.xsd">
   <tx:advice id="txAdvice" transaction-manager="transactionManager"/>
   ```

3. 配置AOP的切入点表达式

4. 建立切入点表达式和事务通知的对应关系

   ```xml
   <!--    配置aop-->
   <!--    配置切入点表达式-->
       <aop:config>
           <aop:pointcut id="pt" expression="execution(* com.hamilemon.service.impl.*.*(..))"/>
   <!--        建立切入点表达式和事务通知的对应关系-->
           <aop:advisor advice-ref="txAdvice" pointcut-ref="pt"/>
       </aop:config>
   ```

5. 配置事务的属性
   在`<tx:advice>`标签中
   `<tx:method>`中的属性:

 - isolation:用于指定事务的隔离级别
 - propagation:用于指定事务的传播行为,默认为REQUITRED,表示一定会有事务.
 只有查询时可以选择SUPPORTS
 - read-only:用于指定事务是否只读.只有查询方法时才能设置为true
 - timeout:用于指定事务的超时时间.默认为-1,表示永不超时(以秒为单位)
 - rollback-for:用于指定一个异常,当产生该异常时,事务回滚;产生其它异常,不会滚
 没有默认值,表示任何异常都回滚
 - no-rollback-for:用于指定一个异常,当产生时,事务不回滚;
 没有默认值,表示任何异常都回滚

```xml
<!--配置事务属性-->
        <tx:attributes>
            <tx:method name="*" propagation="REQUIRED" read-only="false"/>
            <tx:method name="find*" propagation="SUPPORTS" read-only="true"/>
        </tx:attributes>
```

- 基于注解的事务控制

    1. 配置事务管理器(同上)
    1. 开启spring对注解事务的支持
`<tx:annotation-driven transation-manager="transactionManager/>`
    1. 在Service类上添加`@Transactional`注解,可以作用在方法上
`@Transactional(propagation=Propagation.SUPORTS)`
