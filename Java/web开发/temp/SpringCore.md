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

beans-factory-nature

