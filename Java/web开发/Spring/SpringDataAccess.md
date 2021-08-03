# Spring Data Access

本笔记整理自[Spring Data Access](https://docs.spring.io/spring-framework/docs/current/reference/html/data-access.html)

主要涉及Spring中的事务管理。

## 事务管理

### 声明式事务管理

Spring提供一种声明式的事务管理，这种方式能够最小程度的影响代码。

声明式事务通过AOP代理来实现其事务通知由XML中的配置或注解来驱动。生成的AOP代理会使用`TransactionInterceptor`结合一个合适的`TransactionManager`实现类来驱动事务。Spring中的`TransactionInterceptor`提供了命令式编程和响应式编程的事务管理，它会自动根据方法的返回值选择合适的方式，如果方法返回值是一个响应式类型（如`Publisher`）则会使用响应式事务管理，其它任何类型（包括void）都会使用命令式事务管理。

以下为通过事务代理调用方法的概念图：

![tx](https://docs.spring.io/spring-framework/docs/current/reference/html/images/tx.png)

### 声明式事务示例

- 一个service接口

  ```java
  public interface FooService {
  
      Foo getFoo(String fooName);
  
      Foo getFoo(String fooName, String barName);
  
      void insertFoo(Foo foo);
  
      void updateFoo(Foo foo);
  
  }
  ```

- 其实现类

  ```java
  public class DefaultFooService implements FooService {
  
      @Override
      public Foo getFoo(String fooName) {
          // ...
      }
  
      @Override
      public Foo getFoo(String fooName, String barName) {
          // ...
      }
  
      @Override
      public void insertFoo(Foo foo) {
          // ...
      }
  
      @Override
      public void updateFoo(Foo foo) {
          // ...
      }
  }
  ```

  假设`getFoo`方法需要运行在只读（read-only）事务中，而其它方法则运行在读写（read-write）事务中。
  
- Spring 配置文件

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
      <!-- 需要事务管理的service bean -->
      <bean id="fooService" class="x.y.service.DefaultFooService"/>
      <!-- 配置事务通知 -->
      <tx:advice id="txAdvice" transaction-manager="txManager">
          <!-- the transactional semantics... -->
          <tx:attributes>
              <!-- 所有以get开头的方法，且为只读 -->
              <tx:method name="get*" read-only="true"/>
              <!-- other methods use the default transaction settings (see below) -->
              <tx:method name="*"/>
          </tx:attributes>
      </tx:advice>
  
      <!-- 使上面定义的事务通知运行在所有FooService的方法中 -->
      <aop:config>
          <aop:pointcut id="fooServiceOperation" expression="execution(* x.y.service.FooService.*(..))"/>
          <aop:advisor advice-ref="txAdvice" pointcut-ref="fooServiceOperation"/>
      </aop:config>
  
      <!-- 数据源 ，这里使用DBCP数据源-->
      <bean id="dataSource" class="org.apache.commons.dbcp.BasicDataSource" destroy-method="close">
          <property name="driverClassName" value="oracle.jdbc.driver.OracleDriver"/>
          <property name="url" value="jdbc:oracle:thin:@rj-t42:1521:elvis"/>
          <property name="username" value="scott"/>
          <property name="password" value="tiger"/>
      </bean>
  
      <!-- 事务管理器 -->
      <bean id="txManager" class="org.springframework.jdbc.datasource.DataSourceTransactionManager">
          <property name="dataSource" ref="dataSource"/>
      </bean>
      <!-- other <bean/> definitions here -->
  </beans>
  ```

### 事务回滚

在Spring中最佳的做法是通过抛出异常来触发事务回滚。Spring会捕获所有未被处理的异常，并确定是否回滚事务。

对于默认配置，Spring只有当抛出运行时异常时（`RuntimeException`的子类），才会触发事务回滚，而其它检查异常则不会触发事务回滚。可以显式地配置异常的类型，从而让其也会触发事务回滚。

```xml
<tx:advice id="txAdvice" transaction-manager="txManager">
    <tx:attributes>
    <tx:method name="get*" read-only="true" rollback-for="NoProductInStockException"/>
    <tx:method name="*"/>
    </tx:attributes>
</tx:advice>
```

相反地，也可以显式地指定某些异常不触发事务回滚：

```xml
<tx:advice id="txAdvice">
    <tx:attributes>
    <tx:method name="updateStock" no-rollback-for="InstrumentNotFoundException"/>
    <tx:method name="*"/>
    </tx:attributes>
</tx:advice>
```

### 配置`<tx:advice/>`

对于`<tx:advice/>`默认的设置为：

- 事务设置为`REQUIRED`
- 隔离级别为`DEFAULT`
- 事务为读写事务（read-write）
- 任意的`RuntimeException`都会触发事务回滚

下表展示了内嵌在`<tx:advice/>`和`<tx:attributes/>`中的`<tx:method/>`标签中的属性：

| 属性              | 是否必需 | 默认值     | 描述                                                         |
| ----------------- | -------- | ---------- | ------------------------------------------------------------ |
| `name`            | Y        |            | 事务关联的方法名，可以使用`*`进行通配                        |
| `propagation`     | N        | `REQUIRED` | 事务的传播行为                                               |
| `isolation`       | N        | `DEFAULT`  | 事务的隔离级别，只有`propagation`为`REQUIRED` 或`REQUIRES_NEW`时可用。 |
| `timeout`         | N        | -1         | 事务超时时间（单位秒），同样只在`propagation`为`REQUIRED` 或`REQUIRES_NEW`时可用。 |
| `read-only`       | N        | false      | 是否是只读，默认为`read-write`，只在`propagation`为`REQUIRED` 或`REQUIRES_NEW`时可用。 |
| `rollback-for`    | N        |            | 指定触发事务回滚的异常类型，有多个时使用逗号分隔。           |
| `no-rollback-for` | N        |            | 和上面相反，指定不触发回滚的异常。                           |

### 使用`@Transactional`

除了使用XML配置文件来声明事务，还可以使用注解方式来声明事务。

```java
@Transactional
public class DefaultFooService implements FooService {

    @Override
    public Foo getFoo(String fooName) {
        // ...
    }

    @Override
    public Foo getFoo(String fooName, String barName) {
        // ...
    }

    @Override
    public void insertFoo(Foo foo) {
        // ...
    }

    @Override
    public void updateFoo(Foo foo) {
        // ...
    }
}
```

在类上声明`@Transactional`可使其中（及其子类）所有方法应用默认的事务配置，如果把上面的类作为bean由IOC容器进行管理，则可以使用`@EnableTransactionManagement`注解在一个`@Configuration`类上，来使其配置的事务生效。

在XML文件中则是使用`<tx:annotation-driven/>`：

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
    <!-- this is the service object that we want to make transactional -->
    <bean id="fooService" class="x.y.service.DefaultFooService"/>
    <!-- enable the configuration of transactional behavior based on annotations -->
    <!-- a TransactionManager is still required -->
    <tx:annotation-driven transaction-manager="txManager"/> 
    <bean id="txManager" class="org.springframework.jdbc.datasource.DataSourceTransactionManager">
        <!-- (this dependency is defined somewhere else) -->
        <property name="dataSource" ref="dataSource"/>
    </bean>
    <!-- other <bean/> definitions here -->
</beans>
```

注：在`<tx:annotation-driven/>`中，如果事务管理器`TransactionManager`的名称为`transactionManager`，则可以不显式的指定`transaction-manager`。

`@Transactional`可以用在接口，接口中的方法，类，类中的方法上，但建议用在一个具体的类及其中的方法上。

#### 配置基于注解方式的事务设置

注：此配置可用于`<tx:annotation-driven/>`和`@EnableTransactionManagement`中

| XML中属性             | 注解中属性                                                   | 默认值               | 描述                        |
| ---------- | ------------------------------------------------------------ | -------- | --------------------------- |
| `transaction-manager` | 不适用，见[`TransactionManagementConfigurer`](https://docs.spring.io/spring-framework/docs/5.3.9/javadoc-api/org/springframework/transaction/annotation/TransactionManagementConfigurer.html) | `transactionManager` | 指定事务管理器的名称。      |
| `mode`                | `mode`                                                       | `proxy`              | 默认使用Spring的AOP进行代理，另一种则是使用AspectJ来代理。 |
| `proxy-target-class` | `proxyTargetClass` | `false` | 控制代理方式，默认使用JDK的基于接口代理，为`true`时则是基于子类代理。 |

#### `@Transacatonal`的设置

主要用于对事务进行配置

| 属性    | 类型     | 描述                             |
| ------- | -------- | -------------------------------- |
| `value` | `String` | 指定所使用的事务管理器。（可选） |

其余配置可参考`<tx:advice/>`的配置。

```java
@Transactional(readOnly = true)
public class DefaultFooService implements FooService {

    public Foo getFoo(String fooName) {
        // ...
    }

    // these settings have precedence for this method
    @Transactional(readOnly = false, propagation = Propagation.REQUIRES_NEW)
    public void updateFoo(Foo foo) {
        // ...
    }
}
```

### 事务传播行为（Transaction Propagation）

指当一个事务方法调用另一个事务方法时，另一个事务方法应该如何执行。例如：事务方法methodA中调用了方法事务方法methodB，methodB执行时继续在methodA中的事务中运行，还是新建一个事务来执行，由methodB的传播行为决定。在Spring中定义了七种传播行为：

- `REQUIRED`

  使用当前事务，如果当前不存在事务就新建一个。默认行为。

- `SUPPORT`

  使用当前事务，如果当前不存在事务就以非事务方式运行。

- `MANDATORY`

  使用当前事务，如果当前不存在事务就抛出异常。

- `REQUIRES_NEW`

  创建一个新事务，如果当前存在事务就将其暂停。需要使用`JtaTransactionManager`。

- `NOT_SUPPORTED`

  以非事务方式运行，如果当前存在事务就将其暂停，需要使用`JtaTransactionManager`。

- `NEVER`

  以非事务方式运行，如果当前存在事务就抛出异常。

- `NESTED`

  如果当前存在事务就创建一个内嵌事务，否则和`REQUIRED`相同。只在`DataSourceTransactionManager`支持。

### 事务隔离级别

#### 由事务产生的问题

- 脏读

  读取到其它事务未提交的数据，因为未提交的数据可能会被回滚，所以读取的数据是“脏数据”。

- 不可重复读

  指在一个事务中多次读取数据，而在这个事务尚未完成时，另一个事务也访问该数据并对其进行修改，从而导致在第一个事务中，每次读取到的数据不一样。

- 幻读

  于不可重复读类似，在第一个事务读取了几行数据后，另一个事务添加了一些数据，从而导致第一个事务在再次读取时会多出一些原本不存在的数据，如同发生了幻觉。

#### Spring中的事务隔离级别

- DEFAULT

  使用数据库默认的隔离级别。

- READ_UNCOMMITED

  允许读取未提交的数据，会发生脏读，不可重复读，幻读。

- READ_COMMITED

  只允许读取已提交的数据，可以阻止脏读，但仍会发生不可重复读，幻读。

- REPEATABLE_READ

  禁止读取未提交的数据，也禁止在第一个事务数据某一行数据时，另一个事务对其进行修改，可以阻止脏读，不可重复读，但仍会发生幻读。

- SERIALIZABLE

  包含了`REPEATABLE_READ`中禁止的内容，同时也禁止第一个事务在读取出一系列满足条件的数据时，另一个事务添加满足条件的数据，可以阻止脏读，不可重复读，幻读。

### 程序式事务管理

Spring也支持程度式的事务管理，通过使用以下两个类：

- `TransactionTemplate`或`TransactionOperator`
- 直接使用`TransactionManager`的实现类

#### 使用`TransactionTemplate`

```java
public class SimpleService implements Service {

    private final TransactionTemplate transactionTemplate;

    //通过构造方法注入TransactionManager实例化TransactionTemplate
    public SimpleService(PlatformTransactionManager transactionManager) {
        this.transactionTemplate = new TransactionTemplate(transactionManager); 
        //配置事务
        this.transactionTemplate
            .setIsolationLevel(TransactionDefinition.ISOLATION_READ_UNCOMMITTED);
        this.transactionTemplate.setTimeout(30); 
    }

    public Object someServiceMethod() {
        return transactionTemplate.execute(new TransactionCallback() {
            //此方法将在具有事务的上下文中执行
            public Object doInTransaction(TransactionStatus status) {
                updateOperation1();
                return resultOfUpdateOperation2();
            }
        });
    }
    
    public void otherServiceMethod(){
        //没有返回值
        transactionTemplate.execute(new TransactionCallbackWithoutResult() {
    		protected void doInTransactionWithoutResult(TransactionStatus status) {
        		try {
            		updateOperation1();
            		updateOperation2();
       			 } catch (SomeBusinessException ex) {
                    //回滚事务
           			 status.setRollbackOnly();
        		}
    		}
		});
    }
}
```

#### 使用`TransactionOperator`

采用了和响应式操作类似的设计，主要也是用于响应式事务操作。

```java
public class SimpleService implements Service {

    // single TransactionOperator shared amongst all methods in this instance
    private final TransactionalOperator transactionalOperator;

    // use constructor-injection to supply the ReactiveTransactionManager
    public SimpleService(ReactiveTransactionManager transactionManager) {
        this.transactionOperator = TransactionalOperator.create(transactionManager);
    }

    public Mono<Object> someServiceMethod() {
        // the code in this method runs in a transactional context
        Mono<Object> update = updateOperation1();

        return update.then(resultOfUpdateOperation2)
            		 .as(transactionalOperator::transactional);
    }
}
```

#### 直接使用`TransactionManager`的实现类

可以直接使用`TransactionManager`的实现类，例如：`PlatformTransactionManager`进行事务管理。

```java
//对事务进行配置
DefaultTransactionDefinition def = new DefaultTransactionDefinition();
def.setName("SomeTxName");
def.setPropagationBehavior(TransactionDefinition.PROPAGATION_REQUIRED);
//通过TransactionManager获取一个事务
TransactionStatus status = txManager.getTransaction(def);
try {
    // put your business logic here
}
catch (MyException ex) {
    //事务回滚
    txManager.rollback(status);
    throw ex;
}
//事务提交
txManager.commit(status);
```

## JdbcTemplate

Spring中提供了一个JdbcTemplate可以简化传统的JDBC操作（诸如获取资源，关闭链接等操作）。

```java
@Repository
public class JdbcMovieFinder implements MovieFinder {
    private JdbcTemplate jdbcTemplate;
    @Autowired
    public void init(DataSource dataSource) {
        this.jdbcTemplate = new JdbcTemplate(dataSource);
    }
}
```

只需要传入一个数据源进行初始化，然后就可以直接使用它进行操作。

### 查询（SELECT）

查询共有多少数据：

```java
int rowCount = jdbcTemplate.queryForObject("select count(*) from t_actor", Integer.class);
```

带有一个参数进行查询：

```java
int countOfActorsNamedJoe = jdbcTemplate.queryForObject(
        "select count(*) from t_actor where first_name = ?", Integer.class, "Joe");
```

查询结果是一个`String`:

```java
String lastName = jdbcTemplate.queryForObject(
        "select last_name from t_actor where id = ?",
        String.class, 1212L);
```

查询结果一一个复杂对象:

```java
Actor actor = jdbcTemplate.queryForObject(
        "select first_name, last_name from t_actor where id = ?",
        (resultSet, rowNum) -> {
            Actor newActor = new Actor();
            newActor.setFirstName(resultSet.getString("first_name"));
            newActor.setLastName(resultSet.getString("last_name"));
            return newActor;
        },
        1212L);
```

查询结果一个集合:

```java
List<Actor> actors = this.jdbcTemplate.query(
        "select first_name, last_name from t_actor",
        (resultSet, rowNum) -> {
            Actor actor = new Actor();
            actor.setFirstName(resultSet.getString("first_name"));
            actor.setLastName(resultSet.getString("last_name"));
            return actor;
        });
```

### 更新（INSERT,UPDATE,DELETE）

插入一条数据：

```java
this.jdbcTemplate.update(
        "insert into t_actor (first_name, last_name) values (?, ?)",
        "Leonor", "Watling");
```

修改数据:

```java
this.jdbcTemplate.update(
        "update t_actor set last_name = ? where id = ?",
        "Banjo", 5276L);
```

删除数据：

```java
this.jdbcTemplate.update(
        "delete from t_actor where id = ?",
        Long.valueOf(actorId));
```

## 内嵌式数据库支持

Spring中提供了对于内存级数据库引擎的支持（包括HSQL，H2，Dery）。可以创建一个内嵌的数据库并由IOC容器管理。

### 通过XML方式创建

```xml
<jdbc:embedded-database id="dataSource" generate-name="true">
    <jdbc:script location="classpath:schema.sql"/>
    <jdbc:script location="classpath:test-data.sql"/>
</jdbc:embedded-database>
```

以上方式会创建一个HSQL的内嵌数据库（可通过type属性指定创建哪种数据库），并且执行`sechema.sql`和`test-data.sql`中的语句，并且获取此数据库链接的数据源对象会被加入到IOC容器中。注：还需要相关数据库的依赖才能成功创建。

### 通过代码方式创建

```java
@Configuration
public class DataSourceConfig {

    @Bean
    public DataSource dataSource() {
        return new EmbeddedDatabaseBuilder()
                .generateUniqueName(true)
                .setType(H2)
                .setScriptEncoding("UTF-8")
                .ignoreFailedDrops(true)
                .addScript("schema.sql")
                .build();
    }
}
```
