# Spring Core

Spring框架的核心部分，主要涉及Spring的IOC容器以及AOP

## IOC容器

IOC(控制反转)也被称为依赖注入，它是在一个对象的构造方法，工厂方法的参数上或属性上声明相关依赖，然后由IOC容器在创建这个对象（Spring中称其为Bean）时注入相关的依赖。

> It is a process whereby objects define their dependencies (that is, the other objects they work with) only through constructor arguments, arguments to a factory method, or properties that are set on the object instance after it is constructed or returned from a factory method. The container then injects those dependencies when it creates the bean.

### 使用IOC

一个基本的XML配置文件

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd">
	<!--使用bean标签表示向容器中注入一个对象，id为这个对象的唯一标识，class则是指定对象类型-->
    <bean id="..." class="...">  
        <!-- collaborators and configuration for this bean go here -->
    </bean>
    <bean id="..." class="...">
        <!-- collaborators and configuration for this bean go here -->
    </bean>
</beans>
```

加载容器：

```java
ApplicationContext context = new ClassPathXmlApplicationContext("services.xml", "daos.xml");
//获取容器中的对象
PetStoreService service = context.getBean("petStore", PetStoreService.class);
```

#### 构建Bean

- 使用默认无参构造方法

  ```xml
  <bean id="exampleBean" class="examples.ExampleBean"/>
  <bean name="anotherExample" class="examples.ExampleBeanTwo"/>
  ```
  
- 使用静态工厂方法构建

  ```xml
  <bean id="clientService"
      class="examples.ClientService"
      factory-method="createInstance"/>
  ```

  ```java
  public class ClientService {
      private static ClientService clientService = new ClientService();
      private ClientService() {}
      public static ClientService createInstance() {
          return clientService;
      }
  }
  ```
  
- 使用实例工厂方法构建

  ```xml
  <!--工厂bean-->
  <bean id="serviceLocator" class="examples.DefaultServiceLocator">
      <!-- inject any dependencies required by this locator bean -->
  </bean>
  <!--会调用工厂中的createClientServiceInstance方法来创建这个bean-->
  <bean id="clientService"
      factory-bean="serviceLocator"
      factory-method="createClientServiceInstance"/>
  ```
  

#### 依赖注入

- 构造方法注入

  在调用构造方法可能会需要一系列的参数，而这些参数则可看作是这个bean需要的依赖，Spring中可以使用`<constructor-arg/>`标签来设置构造方法中参数的值。

  如果使用静态工厂方法来创建也可以使用这个标签来设置工厂方法中的参数（实例工厂也相同），此时Spring则只负责传参，而不管对象的创建。（因为创建由工厂来完成）

  ```java
  package x.y;
  public class ThingOne {
      public ThingOne(ThingTwo thingTwo, ThingThree thingThree) {
          // ...
      }
  }
  ```

  ```xml
  <beans>
      <bean id="beanOne" class="x.y.ThingOne">
          <!--index指定在参数列表中对应的顺序，从0开始-->
          <constructor-arg index="0" ref="beanTwo"/>
          <constructor-arg index="1" ref="beanThree"/>
          <!--ref用于引用已经存在的bean-->
          <!--可以使用value属性设置基本类型的值-->
      </bean>
      <bean id="beanTwo" class="x.y.ThingTwo"/>
      <bean id="beanThree" class="x.y.ThingThree"/>
  </beans>
  ```

  注：Spring还可以使用name来对应参数名，`<constructor-arg name="thingThree" ref="beanThree"/>`但由于参数名在运行时不可见，所以需要在构造方法上使用注解`@ConstructorProperties({"thingTow", "thingThree"})`来显式指定参数名。

- Setter方法注入

  通过调用无参构造创建出实例后，调用Setter方法去注入依赖。当然Setter方法注入也可以和构造方法注入一起使用。

  > 在官方文档中有提到，最佳的做法是使用构造方法来注入那些必需的属性（不可变，不可空，对象创建时就应该被确定的属性），使用settter方法来注入可选的属性（即使 使用默认值也不会影响这个类的功能。

  ```xml
  <bean id="exampleBean" class="examples.ExampleBean">
      <property name="beanOne">
          <!--可以在property内嵌一个bean-->
          <ref bean="anotherExampleBean"/>
          <!--<bean class="examples.AnotherBean"/>-->
      </property>
      <!--ref属性引用bean,value给基本类型赋值-->
      <property name="beanTwo" ref="yetAnotherBean"/>
      <property name="integerProperty" value="1"/>
  </bean>
  <bean id="anotherExampleBean" class="examples.AnotherBean"/>
  <bean id="yetAnotherBean" class="examples.YetAnotherBean"/>
  ```
  
  ```java
  public class ExampleBean {
      private AnotherBean beanOne;
      private YetAnotherBean beanTwo;
      private int i;
      public void setBeanOne(AnotherBean beanOne) {
          this.beanOne = beanOne;
      }
      public void setBeanTwo(YetAnotherBean beanTwo) {
          this.beanTwo = beanTwo;
      }
      public void setIntegerProperty(int i) {
          this.i = i;
      }
  }
  ```

Spring解决依赖的过程，详见[beans-dependency-resolution](https://docs.spring.io/spring-framework/docs/current/reference/html/core.html#beans-dependency-resolution)：

1. 首先创建出所有配置的bean对象
2. 对每一个bean，其依赖关系都由属性，构造方法的参数或者静态工厂方法中的参数来定义，在创建对象后，罗列出这些依赖关系。
3. 对每一个对象所依赖属性或构造方法中的参数，都会使用定义好的值或者引用的其它bean去注入。
4. 在配置文件中，value中的内容虽然以字符串形式表示但Spring会在注入时将其转换为对应的类型。

Spring在创建容器后会初始化所有的bean为预实例化状态（创建出了对象但不设置属性），只有在这个bean被请求时才去完整的创建出来（属性值被注入），同时，这个bean所依赖的其它bean也会被完整创建出来。

注：环形依赖。例如，对象A的构造方法需要对象B，而对象B的构造方法又需要对象A（类似于先有鸡还是先有蛋的问题），对于这种环形依赖会导致`BeanCurrentlyInCreationException`，解决方法是将其修改为settter方法注入。

##### `idref`标签

如果对象中的某一个属性需要的是一个bean的id(只是id，一个字符串而不是一个bean)，则可以使用idref，当然也可以直接使用value属性赋值，但使用idref的好处是能够确保这个id一定指向了一个bean。

```xml
<bean id="theTargetBean" class="..."/>
<bean id="theClientBean" class="...">
    <property name="targetName">
        <idref bean="theTargetBean"/>
    </property>
</bean>
```

##### 内嵌Bean

可以在`<property/>`和 ` <constructor-arg/>`中内嵌`<bean/>`

```xml
<bean id="outer" class="...">
    <property name="target">
        <bean class="com.example.Person"> <!-- this is the inner bean -->
            <property name="name" value="Fiona Apple"/>
        </bean>
    </property>
</bean>
```

##### 集合注入

```xml
<bean id="moreComplexObject" class="example.ComplexObject">
    <!-- results in a setAdminEmails(java.util.Properties) call -->
    <property name="adminEmails">
        <props>
            <prop key="administrator">administrator@example.org</prop>
            <prop key="development">development@example.org</prop>
        </props>
    </property>
    <!-- results in a setSomeList(java.util.List) call -->
    <property name="someList">
        <list>
            <value>a list element followed by a reference</value>
            <ref bean="myDataSource" />
        </list>
    </property>
    <!-- results in a setSomeMap(java.util.Map) call -->
    <property name="someMap">
        <map>
            <entry key="an entry" value="just some string"/>
            <entry key ="a ref" value-ref="myDataSource"/>
        </map>
    </property>
    <!-- results in a setSomeSet(java.util.Set) call -->
    <property name="someSet">
        <set>
            <value>just some string</value>
            <ref bean="myDataSource" />
        </set>
    </property>
</bean>
```

##### 空值和空字符串

空值（null）可以使用`<null/>`设置，空字符串则`value=""`即可。

```xml
<bean class="ExampleBean">
    <!--这是空字符串-->
    <property name="email" value=""/>
    <!--这是空值null-->
    <property name="name">
        <null/>
    </property>
</bean>
```

##### 简写

可以使用p-命名空间来简化`<property/>`的书写。

```xml
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:p="http://www.springframework.org/schema/p"
    xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd">
    <bean name="jane" class="com.example.Person">
        <property name="name" value="Jane Doe"/>
    </bean>
	<!--传统写法-->
    <bean name="john-classic" class="com.example.Person">
        <property name="name" value="John Doe"/>
        <property name="spouse" ref="jane"/>
    </bean>
	<!--p命名空间写法-->
    <bean name="john-modern"
        class="com.example.Person"
        p:name="John Doe"
        p:spouse-ref="jane"/>
    <!--使用 -ref 后缀来引用其它bean-->
</beans>
```

类似的还有c-命名空间来简化`<constructor-arg/>`

```xml
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:c="http://www.springframework.org/schema/c"
    xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd">

    <bean id="beanTwo" class="x.y.ThingTwo"/>
    <bean id="beanThree" class="x.y.ThingThree"/>
    <!-- 传统方式 -->
    <bean id="beanOne" class="x.y.ThingOne">
        <constructor-arg name="thingTwo" ref="beanTwo"/>
        <constructor-arg name="thingThree" ref="beanThree"/>
        <constructor-arg name="email" value="something@somewhere.com"/>
    </bean>
    <!-- c-namespace 使用name-->
    <bean id="beanOne" class="x.y.ThingOne" c:thingTwo-ref="beanTwo"
        c:thingThree-ref="beanThree" c:email="something@somewhere.com"/>
    <!--使用索引, 数字前必须加 _  -->
	<bean id="beanOne" class="x.y.ThingOne" c:_0-ref="beanTwo" c:_1-ref="beanThree"
    	c:_2="something@somewhere.com"/>
</beans>
```

##### 复合属性名

```xml
<bean id="something" class="things.ThingOne">
    <property name="fred.bob.sammy" value="123" />
</bean>
```

在ThingOne类中有一个属性为fred，fred中有一个属性为bob，bob有一个属性为sammy。需要在ThingOne对象被构造出来时，fred和bob都不为null。

#### `depends-on`标签

通常表明两个bean之间的依赖可以使用`ref`属性，但有时候两个bean的依赖并不那么直接，比如一个类中一静态代码块需要先执行，例如注册JDBC驱动。需要先加载驱动类，然后再使用`DriverManager`。在Spring中可以使用`depends-on`强制一个或多个bean在当前bean（使用`depends-on`标签的bean）被初始化前初始化。

```xml
<!--manager,accountDao会比beanOne先初始化-->
<bean id="beanOne" class="ExampleBean" depends-on="manager,accountDao">
    <property name="manager" ref="manager" />
</bean>
<bean id="manager" class="ManagerBean" />
<bean id="accountDao" class="x.y.jdbc.JdbcAccountDao" />
```

在销毁时，如果所有的bean都是单例的（只有单例的bean才会由spring管理销毁），则`depends-on`中的bean销毁顺序会在当前bean之后。

#### 懒加载

Spring默认在创建容器时预实例化所有的bean，当指定其需要懒加载时，只有在这个bean被请求使用时才会去实例化（不会在容器创建时被预实例化）。

```xml
<!--懒加载-->
<bean id="lazy" class="com.something.ExpensiveToCreateBean" lazy-init="true"/>
<bean name="not.lazy" class="com.something.AnotherBean"/>
```

但是当一个非懒加载的bean依赖一个懒加载的bean时，为了满足第一个bean的依赖，这个懒加载的bean并不会被懒加载。

可以在`<beans/>`上指定里面的所有bean都被懒加载。

```xml
<beans default-lazy-init="true">
    <!-- no beans will be pre-instantiated... -->
</beans>
```

#### 自动注入

Spring可以自动注入一个所需要的依赖，而不用显式的描述出来。

在`<bean/>`标签上设置`autowire`属性即可使用此功能，其有以下几种取值

| `no`          | 默认值，不使用自动注入                                       |
| ------------- | ------------------------------------------------------------ |
| `byName`      | 通过名称注入属性，例如在一个bean中有一个master属性且有一个setMaster方法，那么spring会在容器中寻找名称为master的bean并自动注入。 |
| `byType`      | 通过类型注入，类似`byName`只是通过类型去查找，当有不止一个bean满足时，会出现异常；当一个也没有找到时，则不会注入。 |
| `constructor` | 类似于`byType`只是作用于构造方法的参数，并且当一个也没有找到时，会出现异常。 |

##### 排除被自动注入

可以设置`<bean/>`标签上的`autowire-candidate`为`false`，在自动注入查找依赖时，会排除当前这个bean。（包括使用注解进行自动注入）

注：这个属性主要设计给通过类型自动注入，如果在自动注入时，有多个bean因为类型而满足，那么通过设置`autowire-candidata`为`false`可以将当前bean排除。但如果是通过名称自动注入，且当前bean正好满足，则仍然会被自动注入到目标bean中。

### Bean作用域

可以通过设置`<bean/>`标签的`scope`属性来指定当前bean的作用域（姑且叫做作用域，Spring文档中为Bean Scope）。

```xml
<bean id="accountService" class="com.something.DefaultAccountService" scope="singleton"/>
```

| singleton   | 默认值，单例模式，在IOC容器中只有一个实例。                 |
| ----------- | ----------------------------------------------------------- |
| prototype   | 中IOC容器中有多个实例。                                     |
| request     | 生命周期和一个http request相同，即一个request拥有一个实例。 |
| session     | 生命周期和一个Session相同。                                 |
| application | 生命周期和一个`ServletContext`相同。                        |
| websocket   | 生命周期和一个`WebSocket`相同。                             |

#### Singleton

 在IOC容器中只存在一个实例，所有引用都指向同一个对象。

注意：Spring中的单例和设计模式中的单例并不完全相同，设计模式中的单例主要通过编码方式来实现每个ClassLoader加载对象时都只会获取到同一个对象，而在Spring中则是一个容器中只有一个对象，即在Spring的一个容器中，一个单例对象只会创建一个且只有一个对象（或许可以理解为有多个容器时会有多个对象？）。

![singleton](https://docs.spring.io/spring-framework/docs/current/reference/html/images/singleton.png)

#### prototype

每次获取时都会得到一个新的对象，官方文档中建议，对于有状态的bean最好定义为prototype，而无状态的bean则可以定义为singleton。

相比于其它几种作用域，对于prototype，Spring并不会管理它的整个生命周期，Spring只负责创建对象，而对于释放资源这类操作需自己去处理。

![prototype](https://docs.spring.io/spring-framework/docs/current/reference/html/images/prototype.png)

#### Singleton对象依赖于Prototype

如果一个Singleton的对象（称作A对象）依赖于一个Prototype（称作B对象），Spring会在创建对象时解决这个依赖，但当A对象每次使用B对象，都想重新获取一个新的B对象时，则需要使用方法注入。

- 方式一

  实现`ApplicationContextAware`接口，通过`applicationContext`来手动获取bean。但这种方法会使这个类对Spring产生依赖，耦合度会增加。

  ```java
  public class CommandManager implements ApplicationContextAware {
  
      private ApplicationContext applicationContext;
  
      public Object process(Map commandState) {
          Command command = createCommand();
          command.setState(commandState);
          return command.execute();
      }
      //通过applicationContext来获取bean,由于这个bean是prototype所以每次都是获取到一个新的对象
      protected Command createCommand() {
          return this.applicationContext.getBean("command", Command.class);
      }
      //Spring会在实例化这个类时自动通过setter方法注入applicationContext，无需配置
      public void setApplicationContext(
              ApplicationContext applicationContext) throws BeansException {
          this.applicationContext = applicationContext;
      }
  }
  ```

- 方式二

  使用的动态代理的方式，由Spring动态生成一个子类来实现或重载获取对象的方法。这需要这个类不能是`final`且被覆盖的方法也不能是`final`。

  ```java
  //并不强制要求这个类是抽象类，也可以是非抽象类
  public abstract class CommandManager {
      public void process() {
          createCommand().run();
      }
  
      protected abstract Command createCommand();
  }
  ```

  ```xml
  <bean id="command" class="pojo.Command" scope="prototype"/>
  <bean id="commandManager" class="manager.CommandManager">
      <!--name为待覆盖的方法名，bean为其返回的bean-->
      <lookup-method name="createCommand" bean="command"/>
  </bean>
  ```

#### 网络相关几个的作用域

`request`、`session`、`application`和`websocket`四个作用域只有在web应用中才可用，并且需要使用具有网络感知的`ApplicationContext`，例如`XmlWebApplicationContext`。

如是在项目有使用Spring Web MVC 那么将不需要特别的设置，因为这些都将由`DispatcherServlet`来完成。

#### 不同作用域的bean相互依赖

注：此部分涉及到动态代理相关知识。

主要针对一个长生命周期的bean（称作对象A）依赖于一个短生命周期的bean（称作对象B），如果按照传统方式注入，在对象A使用对象B时，可能这时的对象B已经过期，但在对象A中引用的实例仍然是那个己经过期的对象B。(等价于前面提到的将一个prototype对象注入到一个singleton对象中)

在Spring中除了前面提到的使用方法注入外，还可以使用AOP代理，注入一个代理对象。

```xml
<!--一个session scope的bean，在实际注入时会注入一个代理对象，但是客户端无感知，正常使用就行-->
<bean id="userPreferences" class="com.something.UserPreferences" scope="session">
    <aop:scoped-proxy/>
</bean>
<bean id="userManager" class="com.something.UserManager">
    <property name="userPreferences" ref="userPreferences"/>
</bean>
```

注：Spring会使用CGLIB来生成代理类，但代理类只会拦截public方法，这意味着，不能在其它对象中调用它的非public方法，并且CGLIB是通过生成一个子类来进行代理，因此目标类（需要被代理的那个类）不能为`final`。

##### 使用接口方式生成代理类

Spring也支持使用JDK中基于接口的动态代理，因此，被代理的类（短生命周期的那个类）至少需要实现一个接口，并且在引用它的地方需通过接口来引用（类型定义为接口而不是那个类的具体类型）。

```xml
<!-- DefaultUserPreferences implements the UserPreferences interface -->
<bean id="userPreferences" class="com.stuff.DefaultUserPreferences" scope="session">
    <!--设置proxy-target-class为false即可-->
    <aop:scoped-proxy proxy-target-class="false"/>
</bean>
<bean id="userManager" class="com.stuff.UserManager">
    <property name="userPreferences" ref="userPreferences"/>
</bean>
```

### 优雅地关闭IOC容器

主要针对非web应用，在web应用中，Spring提供了自动关闭的功能。

在非web应用中，需要注册一个shutdown hook来确保正确关闭IOC容器，并且相关的销毁方法也会被执行。

调用`ConfigurableApplicationContext`中的`registerShutdownHook`方法即可

```java
public final class Boot {

    public static void main(final String[] args) throws Exception {
        ConfigurableApplicationContext ctx = new ClassPathXmlApplicationContext("beans.xml");
        // add a shutdown hook for the above context...
        ctx.registerShutdownHook();
        // app runs here...
        // main method exits, hook is called prior to the app shutting down...
    }
}
```

### 基于注解配置

对于Bean的配置可以通过java的注解来完成。要启用此功能需要注册注解支持。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:context="http://www.springframework.org/schema/context"
    xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd
        http://www.springframework.org/schema/context
        https://www.springframework.org/schema/context/spring-context.xsd">
    <!--可以简单地理解为启用注解支持-->
    <context:annotation-config/>
</beans>
```

#### `@Autowired`

自动注入相关属性

- 通过构造方法注入

  在Spring 4.3 中，当bean中只有一个构造方法时，可以不显式声明`Autowired`注解，但当有多个构造方法时，至少有一个构造方法被声明有`Autowired`。

  注：当有多个构造方法上都声明有`Autowired`时，需要设置`required=false`，然后Spring会根据其参数选择哪一个构造方法能被满足（IOC容器中有那个构造方法需要的参数），如果都没有满足，就会去调用默认构造方法（如果有的话）。

  ```java
  public class MovieRecommender {
      private final CustomerPreferenceDao customerPreferenceDao;
      @Autowired
      public MovieRecommender(CustomerPreferenceDao customerPreferenceDao) {
          this.customerPreferenceDao = customerPreferenceDao;
      }
  }
  ```

- 通过setter方法注入

  也可以声明在其它任意名称的方法上

  ```java
  public class SimpleMovieLister {
      private MovieFinder movieFinder;
      @Autowired
      public void setMovieFinder(MovieFinder movieFinder) {
          this.movieFinder = movieFinder;
      }
  }
  ```

- 直接写在属性上

  还可以和构造方法注入一同使用（setter方法也可以）

  ```java
  public class MovieRecommender {
      private final CustomerPreferenceDao customerPreferenceDao;
      @Autowired
      private MovieCatalog movieCatalog;
      @Autowired
      public MovieRecommender(CustomerPreferenceDao customerPreferenceDao) {
          this.customerPreferenceDao = customerPreferenceDao;
      }
  }
  ```

注：

1. 在Spring 5.0中，可以使用`@Nullable`表明允许这个属性为`null`

   ```java
   @Autowired
   public void setMovieFinder(@Nullable MovieFinder movieFinder) {
   }
   ```

2. 当某个属性是Spring中的`BeanFactory`, `ApplicationContext`, `Environment`, `ResourceLoader`, `ApplicationEventPublisher`, and `MessageSource`及其子类时，可以直接属性上声明`@Autowired`。

   ```java
   public class MovieRecommender {
       @Autowired
       private ApplicationContext context;
       public MovieRecommender() {
       }
   }
   ```

##### 结合`@Qualifier`使用

可以在`@Qualifier`指定bean的名称（默认的名称匹配方式是将bean的名称和属性名进行匹配），从而在`@Autowired`匹配多个bean时，从中选择需要的bean。

```java
@Autowired
@Qualifier("main")
private MovieCatalog movieCatalog;
}
```

当`@Qualifier`用在集合上时，也可以具有元素过滤的作用，来将名称满足的bean作为集合进行注入。

#### 通过`@Resource`注入

Spring也支持通过JSR-250中的`@Resource`注解来注入，该注解只能用于属性和settter方法上，并且主要通过名称进行匹配（默认使用属性名）。

```java
public class SimpleMovieLister {

    private MovieFinder movieFinder;

    @Resource(name="myMovieFinder") 
    public void setMovieFinder(MovieFinder movieFinder) {
        this.movieFinder = movieFinder;
    }
}
```

#### `@Value`

主要用来注入外部的配置文件中定义的值（如一个properties文件中的值）以及基本类型，并且支持SpEL表达式。

```java
public class MovieRecommender {
    private final String catalog;
    public MovieRecommender(@Value("${catalog.name}") String catalog) {
        this.catalog = catalog;
    }
}
```

在properties文件中有如下定义

```
catalog.name=MovieCatalog
```

#### `PostConstruct`和`PreDestroy`

同样来自JSC-250中的注解，可用于生命周期回调。

```java
public class CachingMovieLister {

    @PostConstruct
    public void populateMovieCache() {
        // populates the movie cache upon initialization...
    }

    @PreDestroy
    public void clearMovieCache() {
        // clears the movie cache upon destruction...
    }
}
```

#### 扫描Bean

Spring提供了`@Component`、`@Service`、`Controller`和`Repository`四个注解用来声明当前类是一个需要由IOC容器管理的bean，其中`@Component`是一个通用注解，而`@Service`、`@Controller`和`@Repository`分别表示这个bean属于业务层，表现层，持久层（详见[三层架构](https://baike.baidu.com/item/%E4%B8%89%E5%B1%82%E6%9E%B6%E6%9E%84/11031448?fr=aladdin)）。

然后，指定需要扫描的包

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:context="http://www.springframework.org/schema/context"
    xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd
        http://www.springframework.org/schema/context
        https://www.springframework.org/schema/context/spring-context.xsd">
    <context:component-scan base-package="org.example"/>
</beans>
```

注：`<context:component-scan/>`同时隐含着启用`<context:annotation-config/>`，因此，当使用了前者，就可以不用写后者。

##### 扫描包时过滤

在扫描包时，可以在`<context:component-scan/>`中添加子标签`<context:include-filter/>`（被匹配的bean才会保留）或者`<context:exclude-filter/>`（被匹配的bean会被排除）（注解则是`@ComponentScan`中的`includeFilters`和`excludeFilters`）来进行过滤，每一个filter都需要type和expression。type的可选值如下：

| type       | 示例expression               | 描述                             |
| ---------- | ---------------------------- | -------------------------------- |
| annotation | `org.example.SomeAnnotation` | 在类上存在指定的注解时，被匹配。 |
| assignable | `org.example.SomeClass`      | 按照指定的类型进行过滤。         |
| aspectj    | `org.example..*Service+`     | 通过aspectj表达式匹配。          |
| regex      | `org\.example\.Default.*`    | 通过正则表达式匹配类名。         |

xml表示：

```xml
<beans>
    <context:component-scan base-package="org.example">
        <context:include-filter type="regex"
                expression=".*Stub.*Repository"/>
        <context:exclude-filter type="annotation"
                expression="org.springframework.stereotype.Repository"/>
    </context:component-scan>
</beans>
```

注解表示：

```java
@Configuration
@ComponentScan(basePackages = "org.example",
        includeFilters = @Filter(type = FilterType.REGEX, pattern = ".*Stub.*Repository"),
        excludeFilters = @Filter(Repository.class))
public class AppConfig {
    ...
}
```

#### `@Bean`

用在方法上，表明该方法的返回值是一个需要被IOC容器管理的bean，当方法需要参数时，Spring会尝试自动解决依赖。

```java
@Component
public class FactoryMethodComponent {

    @Bean//也可是写在一个被@Configuration注解的类中
    @Qualifier("public")
    public TestBean publicInstance() {
        return new TestBean("publicInstance");
    }
    public void doWork() {
        // Component method implementation omitted
    }
}
```

#### 生成候选bean的索引

虽然类路径扫描非常快，但是Spring内部存在大量的类，添加此依赖，可以通过在编译时创建候选对象的静态列表来提高大型应用程序的启动性能。在此模式下，作为组件扫描目标的所有模块都必须使用此机制。

##### Maven依赖

```xml
<dependencies>
    <dependency>
        <groupId>org.springframework</groupId>
        <artifactId>spring-context-indexer</artifactId>
        <version>5.3.9</version>
        <optional>true</optional>
    </dependency>
</dependencies>
```

这会在应用编译之后生成`META-INF/spring.components`文件，并在运行时自动加载。

#### 使用JSR 330中的标准注解

##### Maven依赖

```xml
<dependency>
    <groupId>javax.inject</groupId>
    <artifactId>javax.inject</artifactId>
    <version>1</version>
</dependency>
```

##### 使用`@Inject`和`@Named`进行依赖注入

等同于`@Autowired`，可使用`@Named`指定bean的名称。

```java
public class SimpleMovieLister {
    private MovieFinder movieFinder;
    @Inject
    public void setMovieFinder(MovieFinder movieFinder) {
        this.movieFinder = movieFinder;
    }
    public void listMovies() {
        this.movieFinder.findMovies(...);
        // ...
    }
}
```

##### `@Named`用在类上

用在类上时等价于`@Component`

```java
@Named("movieListener")
public class SimpleMovieLister {
    private MovieFinder movieFinder;
    @Inject
    public void setMovieFinder(MovieFinder movieFinder) {
        this.movieFinder = movieFinder;
    }
}
```

### 纯注解配置

在上面的基于注解配置中，只是对于bean的配置改为注解方式，但仍需要一个XML文件用于对IOC容器进行一些配置（比如配置组件扫描）。

在Spring中也支持使用纯注解的方式来配置IOC容器，这主要使用到`@Configuration`和`@Bean`注解，前者表明当前类是一个配置类，而后者则用于方法上，表明方法的返回值是一个由IOC容器管理的bean。

```java
@Configuration
public class AppConfig {

    @Bean
    public MyService myService() {
        return new MyServiceImpl();
    }
}
```

#### 实例化IOC容器

使用XML作为配置文件时，实例化时使用`ClassPathXmlApplicationContext`，当使用纯注解进行配置时，则需要使用`AnnotationConfigApplicationContext`来加载配置类，并实例化IOC容器。

```java
public static void main(String[] args) {
    ApplicationContext ctx = new AnnotationConfigApplicationContext(AppConfig.class);
    MyService myService = ctx.getBean(MyService.class);
    myService.doStuff();
}
```

也可以使用`register`方法来注册配置类：

```java
public static void main(String[] args) {
    AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext();
    ctx.register(AppConfig.class, OtherConfig.class);
    ctx.register(AdditionalConfig.class);
    ctx.refresh();
    MyService myService = ctx.getBean(MyService.class);
    myService.doStuff();
}
```

##### 启用组件扫描

```java
@Configuration
@ComponentScan(basePackages = "com.acme") 
public class AppConfig  {
    ...
}
```

##### web应用中纯注解配置

在web应用中，需要使用`AnnotationConfigWebApplicationContext`来代替`AnnotationConfigApplicationContext`。

web.xml配置为：

```xml
<web-app>
    <!-- 配置 ContextLoaderListener 使用 AnnotationConfigWebApplicationContext
        去代替默认的 XmlWebApplicationContext -->
    <context-param>
        <param-name>contextClass</param-name>
        <param-value>
            org.springframework.web.context.support.AnnotationConfigWebApplicationContext
        </param-value>
    </context-param>

    <!-- 指定配置类，有多个时可以使用逗号或空格进行分隔 -->
    <context-param>
        <param-name>contextConfigLocation</param-name>
        <param-value>com.acme.AppConfig</param-value>
    </context-param>

    <!-- Bootstrap the root application context as usual using ContextLoaderListener -->
    <listener>
        <listener-class>org.springframework.web.context.ContextLoaderListener</listener-class>
    </listener>

    <!-- Declare a Spring MVC DispatcherServlet as usual -->
    <servlet>
        <servlet-name>dispatcher</servlet-name>
        <servlet-class>org.springframework.web.servlet.DispatcherServlet</servlet-class>
        <!-- Configure DispatcherServlet to use AnnotationConfigWebApplicationContext
            instead of the default XmlWebApplicationContext -->
        <init-param>
            <param-name>contextClass</param-name>
            <param-value>
             org.springframework.web.context.support.AnnotationConfigWebApplicationContext
            </param-value>
        </init-param>
        <init-param>
            <param-name>contextConfigLocation</param-name>
            <param-value>com.acme.web.MvcConfig</param-value>
        </init-param>
    </servlet>
    <servlet-mapping>
        <servlet-name>dispatcher</servlet-name>
        <url-pattern>/app/*</url-pattern>
    </servlet-mapping>
</web-app>
```

#### 注解方式的Lookup Method注入

在前面有写到，当一个Singleton对象需要一个Prototype对象作为依赖时，可以使用Lookup 方法注入，在基于注解时，配置如下：

```java
public abstract class CommandManager {
    public Object process(Object commandState) {
        // grab a new instance of the appropriate Command interface
        Command command = createCommand();
        // set the state on the (hopefully brand new) Command instance
        command.setState(commandState);
        return command.execute();
    }
    protected abstract Command createCommand();
}
```

```java
@Bean
@Scope("prototype")
public AsyncCommand asyncCommand() {
    AsyncCommand command = new AsyncCommand();
    return command;
}

@Bean
public CommandManager commandManager() {
    return new CommandManager() {
        protected Command createCommand() {
            return asyncCommand();
        }
    }
}
```

## SpEL(Spring EL 表达式)

Spring中的expression语句，可以在运行时查询和操纵对象。

### 定义bean时使用SpEL表达式

无论是基于注解还是XML的配置都可以使用SpEL表达式，语法为`#{ <expression string> }`

```xml
<bean id="numberGuess" class="org.spring.samples.NumberGuess">
    <property name="randomNumber" value="#{ T(java.lang.Math).random() * 100.0 }"/>
</bean>
```

```java
public class FieldValueTestBean {

    @Value("#{ systemProperties['user.region'] }")
    private String defaultLocale;
    public void setDefaultLocale(String defaultLocale) {
        this.defaultLocale = defaultLocale;
    }
    public String getDefaultLocale() {
        return this.defaultLocale;
    }
}
```

### SpEL语法

#### 字面量

支持的字面量类型为字符串、数字、布尔和null值。字符串由单引号包裹，如果需要在将单引号本身包含在字符串中，则需要两个单引号。

#### 属性和列表

使用`.`来获取对象中的属性，`[]`来获取列表中对应索引的值，对应map集合，可以在`[]`传入键来取值。

#### 内联列表

使用`{}`来创建一个列表，如`{1,2,3,4}`

#### 内联map

使用`{key:value}`来创建一个map，如`{dob:{day:10,month:'July',year:1856}}`

#### 创建数组

和java语法相同，如`new int[4]`，`new int[]{1,2,3}`，`new int[4][5]`。

#### 调用方法

和java语法相同，如`'abc'.substring(1, 3)`，`isMember('Mihajlo Pupin')`。

#### 运算符

##### 关系运算符

和java相同，在比较时，如果有`null`，那么它会被当作 ‘无’ 来对待（不是0），也就是任何值都比`null`大。

同时Spring还提供了类型匹配`instanceof`和正则匹配`matches`。

```
'xyz' instanceof T(Integer)  --- false
'5.00' matches '^-?\\d+(\\.\\d{2})?$'  ---true
```

注：作为基本类型的数据，都会使用包装类型来表示，例如：`1 instanceof T(int)`为false，而`1 instanceof T(Integer)`为`true`。

每一个关系运算符都对应一个等价的字母表示形式：

- `lt` (`<`)
- `gt` (`>`)
- `le` (`<=`)
- `ge` (`>=`)
- `eq` (`==`)
- `ne` (`!=`)
- `div` (`/`)
- `mod` (`%`)
- `not` (`!`)

##### 逻辑运算符

- `and` (`&&`)
- `or` (`||`)
- `not` (`!`)

#### 类型

通过`T()`表明这是一个类型数据（Class类型），可用于类型匹配或者获取其中的static成员，默认会按照给定的名称去`java.lang`包查找，其它包则需要使用全限定类名。如`T(String)`，`T(java.util.Date)`

#### 构造器

使用`new`操作符来创建对象，除`java.lang`包外的类，都需要使用全限定类名。

#### 三元运算符

和java中相同，如`false ? 'trueExp' : 'falseExp'`

#### 猫王运算符

`name?:'Unknown'`，主要用于判空，等价于`name != null ? name : 'Unknown'`，和kotlin和的`?:`相同。

## AOP(面向切面编程)

面向切面编程（AOP）是对面向对象编程（OOP）地一种补充，在OOP中，最小单位是一个类，而在AOP中则是一个切面。

### AOP相关术语

- 切面（aspect）：一种跨越多个类的模块化概念，由切点和增强组成。
- 连接点（join point）：程序运行中的某个特定位置。
- 通知（advice）：切面在一个特定连接点的行为，分为环绕通知，后置通知，前置通知等。
- 切入点（pointcut）：定义通知在哪些连接点上执行。通知会和一个切入点表达式相关联并且运行在被切入点匹配的连接点上。
- 目标对象（target object）：包含一个或多个切面的对象。
- 代理对象（AOP proxy）：由AOP框架创建，用于实现切面操作的代理对象，在Spring中由JDK的动态代理或CGLIB来实现。
- 织入（weaving）：将切面应用到目标对象，并创建代理对象的过程。

在Spring AOP中包含以下几种通知（advice）：

- 前置通知（before advice）：运行在连接点前，但不会阻止运行到连接点的通知。
- 后置返回通知（after returning advice）：在连接点运行完成后运行的通知（例如一个方法正常返回且没有抛出异常）
- 后置异常通知（after throwing advice）：在方法因异常而结束时运行的通知。
- 后置最终通知（after advice）：连接点运行结束时运行的通知。（不管是正常返回还是抛出异常）
- 环绕通知（around advice）：围绕在连接点前后的通知，可以在方法调用前后执行自定义的行为，并且需要决定是继续运行还是中断。

### 启用@AspectJ 支持

在AspectJ项目中使用注解的风格来定义一个切面，在Spring中也采用了这种风格，并且使用AspectJ提供的库来解析和匹配切入点。但在运行时仍然是Spring的原生AOP，并不依赖于AspectJ的解释器和织入器。

@AspectJ支持则是基于@AspectJ 切面来配置AOP，并且当检测到一个bean被一个或多个切面通知时，会自动生成其代理类并且拦截方法以执行通知。

注：因为切入点需要使用AspectJ提供的库来解析，所以需要添加`aspectjweaver`依赖。

- 使用注解开启@AspectJ支持

  ```java
  @Configuration
  @EnableAspectJAutoProxy
  public class AppConfig {
  
  }
  ```

- xml方式

  ```xml
  <aop:aspectj-autoproxy/>
  ```

### 定义切面

当开启@AspectJ支持后，任何被定义为@AspectJ切面的bean都会被Spring检测到并用于配置Spring AOP。

- XML方式

  ```xml
  <bean id="myAspect" class="org.xyz.NotVeryUsefulAspect">
      <!-- configure properties of the aspect here -->
  </bean>
  ```

- 注解方式

  ```java
  @Component
  @Aspect
  public class NotVeryUsefulAspect {
  
  }
  ```

### 定义切入点

由于Spring只支持对方法执行时的连接点，所以定义一个切入点也可理解为匹配一个运行的方法。一个切入点的声明包括两个部分：包含了名称和任意参数的签名以及用来确定方法的切入点表达式。在Spring中签名可以通过一个常规的方法来定义（这个方法返回值必须为void，因为签名并不包含返回值，并且这个方法主要用来定义签名，所以应为空方法），而切入点表达式则使用`@Pointcut`注解。

例如，以下例子定义了一个名为`anyOldTransfer`并且匹配任意方法名为`transfer`的切入点。

```java
@Pointcut("execution(* transfer(..))") // the pointcut expression
private void anyOldTransfer() {} // the pointcut signature
```

在切入点表达式（关于切入点表达式的详细语法见[AspectJ](https://www.eclipse.org/aspectj/doc/released/progguide/index.html)）中，Spring支持使用以下AspectJ切入点指示符（pointcut designators，PCD）:

- `execution`：用于匹配方法运行的连接点。
- `within`：通过指定的类型来限制匹配的连接点。
- `this`：限制匹配的连接点，其中bean的引用是一个指定类型的实例。
- `target`：限制匹配的连接点，其中目标对象是一个指定类型的实例。
- `args`：限制匹配的连接点，其中参数是指定类型的实例。

注：Spring AOP是基于代理实现，所以对于代理（对应this）和代理对象后的目标对象（对应target）是两个不同的对象。

#### 组合切入点表达式

可以使用`&&`，`||`，`!`对切入点表达式进行组合。

```java
//匹配所有public方法
@Pointcut("execution(public * *(..))")
private void anyPublicOperation() {} 
//匹配所有trading包中的方法
@Pointcut("within(com.xyz.myapp.trading..*)")
private void inTrading() {} 
//上面两个的组合
@Pointcut("anyPublicOperation() && inTrading()")
private void tradingOperation() {} 
```

#### 示例

在Spring AOP中多使用`execution`切入点指示符，其格式如下：

```
execution(modifiers-pattern? ret-type-pattern declaring-type-pattern?name-pattern(param-pattern)throws-pattern?)
```

其中，除了`ret-type-pattern`，`name-pattern`和`param-pattern`以外，其余部分均为可选的。`ret-type-pattern`指定方法的返回值类型（需要写全限定名称，可使用`*`通配）；`name-pattern`则是匹配方法的名称，可以使用`*`通配全部或一部分名称；对于`param-pattern`，`()`匹配无参方法，`(..)`匹配任意个数（0个或多个）的参数，`(*)`匹配只有一个参数，类型为任意类型，`(*,String)`则匹配有两个参数且第一个为任意类型，第二个为String类型。

- 任何public方法

  ```
  execution(public * *(..))
  ```

- 任何方法名是以set开头的方法

  ```
  execution(* set*(..))
  ```

- 任何定义在`AccountService`中的方法

  ```
  execution(* com.xyz.service.AccountService.*(..))
  ```

- 任何定义在service包中的方法

  ```
  execution(* com.xyz.service.*.*(..))
  ```

- 任何定义在service包及其子包的方法

  ```
  execution(* com.xyz.service..*.*(..))
  ```

- 任何service包中的连接点

  ```
  within(com.xyz.service.*)
  ```

- 任何service包及其子包中的连接点

  ```
  within(com.xyz.service..*)
  ```

- 任何bean的名称以Service结尾的bean中的连接点

  ```
  bean(*Service)
  ```

### 定义通知

#### 前置通知（Before Advice）

在切入点匹配的方法执行前运行。

```java
@Aspect
public class BeforeExample {
	//也可以引用定义好的切入点@Before("com.xyz.myapp.CommonPointcuts.dataAccessOperation()")
    @Before("execution(* com.xyz.myapp.dao.*.*(..))")
    public void doAccessCheck() {
        // ...
    }
}
```

#### 后置返回通知（After Returning Advice)

在方法return时运行

```java
@Aspect
public class AfterReturningExample {

    @AfterReturning("com.xyz.myapp.CommonPointcuts.dataAccessOperation()")
    public void doAccessCheck() {
        // ...
    }
}
```

```java
@Aspect
public class AfterReturningExample {
	//获取方法的返回值，并作为参数绑定到retVal上
    @AfterReturning(
        pointcut="com.xyz.myapp.CommonPointcuts.dataAccessOperation()",
        returning="retVal")
    public void doAccessCheck(Object retVal) {
        // ...
    }
}
```

#### 后置异常通知（After Throwing Advice）

在方法因抛出异常而退出时运行。注：如果方法在内部将处理（即不抛出异常），那么此通知不会运行。

```java
@Aspect
public class AfterThrowingExample {

    @AfterThrowing("com.xyz.myapp.CommonPointcuts.dataAccessOperation()")
    public void doRecoveryActions() {
        // ...
    }
}
```

可以指定`throwing`属性来限制异常的类型并将其绑定到参数上。

```java
@Aspect
public class AfterThrowingExample {

    @AfterThrowing(
        pointcut="com.xyz.myapp.CommonPointcuts.dataAccessOperation()",
        throwing="ex")
    public void doRecoveryActions(DataAccessException ex) {
        // ...
    }
}
```

#### 后置通知（After Finally Advice）

在方法退出时运行，类似于try-catch语句中finally，一定会被执行，而后置返回通知只在方法成功返回（无异常）时才会运行。

```java
@Aspect
public class AfterFinallyExample {

    @After("com.xyz.myapp.CommonPointcuts.dataAccessOperation()")
    public void doReleaseLock() {
        // ...
    }
}
```

#### 环绕通知

既可以在方法执行前运行，也可以在方法执行后运行，并且能够决定方法是否继续运行。

在环绕通知中，第一个参数必须为`ProceedingJoinPoint`类型，可以通过调用其`proceed`方法来运行切入点匹配的方法，该方法可以调用一次，多次或者一次也不调用。

```java
@Aspect
public class AroundExample {

    @Around("com.xyz.myapp.CommonPointcuts.businessService()")
    public Object doBasicProfiling(ProceedingJoinPoint pjp) throws Throwable {
        // start stopwatch
        Object retVal = pjp.proceed();
        // stop stopwatch
        return retVal;
    }
}
```

#### 执行顺序

![image-20210731175116841](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/07/31/20210731175124.png)

#### 通知优先级

当在同一个连接点有多个通知需要运行时，可以指定每个切面的优先级来确定执行顺序，在确定顺序时，对目标方法执行前运行的通知，优先级越大，越先执行（如对于两个前置通知，优先级大的先运行）；而在目标方法执行后运行的通知，优先级越大，越后执行（如对于两个后置通知，优先级大的后运行）。

在Spring中可以给切面加上`@Order`注解，并设置值。其中，值越小，优先级越大。

