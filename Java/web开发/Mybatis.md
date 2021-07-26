# Mybatis

*本笔记整理自[mybatis](https://mybatis.org/mybatis-3/zh/index.html)*

Mybatis是一款优秀的持久层框架，支持自定义SQL、存储过程以及高级映射。免除了几乎所有的JDBC代码以及设置参数和获取结果集的工作。

## Maven依赖

```xml
<dependency>
  <groupId>org.mybatis</groupId>
  <artifactId>mybatis</artifactId>
  <version>x.x.x</version>
</dependency>
```

## 入门

### 从XML中构建SqlSessionFactory

每个Mybatis应用都是以一个SqlSessionFactory实例为核心，可以通过SqlSessionFactoryBuilder获得其实例。而SqlSessionFactoryBuilder则可以从XML配置文件或一个Configuration实例来构建出SqlSessionFactory

```java
String resource = "SqlSession.xml";
InputStream is = Resources.getResourceAsStream(resource);
SqlSessionFactory factory =
            new SqlSessionFactoryBuilder().build(is);
```

XML配置文件中包含了对Mybatis的核心设置，包括获取数据源，事务作用域和事务管理器。

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE configuration
  PUBLIC "-//mybatis.org//DTD Config 3.0//EN"
  "http://mybatis.org/dtd/mybatis-3-config.dtd">
<configuration>
  <environments default="development">
    <environment id="development">
      <transactionManager type="JDBC"/>
        <!--数据源-->
      <dataSource type="POOLED">
        <property name="driver" value="${driver}"/>
        <property name="url" value="${url}"/>
        <property name="username" value="${username}"/>
        <property name="password" value="${password}"/>
      </dataSource>
    </environment>
  </environments>
    <!--设置映射器-->
  <mappers>
    <mapper resource="pojo-mapper.xml"/>
  </mappers>
</configuration>
```

### 纯Java代码构建SqlSessionFactory

```java
//设置链接数据库的参数
Properties properties = new Properties();
properties.setProperty("driver", "com.mysql.cj.jdbc.Driver");
properties.setProperty("url", "jdbc:mysql://localhost:3306");
properties.setProperty("username", "root");
properties.setProperty("password", "root");
//数据源
PooledDataSourceFactory dataSourceFactory = new PooledDataSourceFactory();
dataSourceFactory.setProperties(properties);
DataSource dataSource = dataSourceFactory
                .getDataSource();
JdbcTransactionFactory transactionFactory = new JdbcTransactionFactory();
Environment mysql = new Environment("mysql", transactionFactory, dataSource);
Configuration configuration = new Configuration(mysql);
configuration.addMapper(IPojoMapper.class);
SqlSessionFactory sessionFactory = new SqlSessionFactoryBuilder()
                .build(configuration);
```

### 获取SqlSession

SqlSession提供了在数据库执行Sql语句的所有方法，可以通过SqlSession来获取一个Dao接口的具体实现类（这个实现类由Mybatis生成），然后去完全各种数据库操作。

- xml方式

  1. IPojoMapper接口

     ```java
     public interface IPojoMapper {
         List<Pojo> getAllPojo();
         
         Pojo getPojo(Integer id);
     }
     ```

  2. pojo-mapper.xml

     ```xml
     <?xml version="1.0" encoding="UTF-8" ?>
     <!DOCTYPE mapper
             PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
             "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
     <!--指定所映射类-->
     <mapper namespace="mapper.IPojoMapper">
         <!--查询方法，id为方法名,注意这里返回指定为对应的pojo类型而不是list-->
         <select id="getAllPojo" resultType="pojo.Pojo">
             select *
             from pojo;
         </select>
         <!--#{id} 等同于占位，运行时使用方法中对应的参数来填充-->
         <select id="getPojo" resultType="pojo.Pojo">
             select *
             from pojo
             where id = #{id}
         </select>
     </mapper>
     ```

- 注解方式

  1. IPojoMapper接口

     ```java
     public interface IPojoMapper {
         @Select("select * from pojo")
         List<Pojo> getAllPojo();
         
         @Select("select * from pojo where id = #{id}")
         Pojo getPojo(Integer id);
     }
     ```

  2. Mybatis主配置文件

     ```xml
     <?xml version="1.0" encoding="UTF-8" ?>
     <!DOCTYPE configuration
             PUBLIC "-//mybatis.org//DTD Config 3.0//EN"
             "http://mybatis.org/dtd/mybatis-3-config.dtd">
     <configuration>
         <environments default="mysql">
             <!--配置mysql环境-->
     ...
             </environment>
         </environments>
     	<!--configuration.addMapper(IPojoMapper.class)-->
     	<!--指定映射器所在的包-->
         <mappers>
             <package name="mapper"/>
         </mappers>
     </configuration>
     ```

- 获取SqlSession并进行操作

  ```java
  try(SqlSession sqlsession = sessionFactory.openSession()){
      //由于mybatis生成相应接口的实现类
  	final IPojoMapper mapper = sqlsession.getMapper(IPojoMapper.class);
          System.out.println(mapper.getAllPojo());
          System.out.println(mapper.getPojo(1));
  }
  ```

  ### 生命周期

  注：如果使用依赖注入框架获取相应实例，可以不用考虑生命周期。

  - SqlSessionFactoryBuilder

    这个实例主要用于创建SqlSessionFactory，创建之后就可以丢弃了，因此可以使其作用域为方法作用域，这样就可以在方法结束后被GC所回收。

  - SqlSessionFactory

    该实例在被创建后就应该一直存在，即和应用的生命周期相同，并且应该每次获取是都是同一个对象，所以应该采用单例模式。

  - SqlSession

    每个线程都应该有一个自己的SqlSession，这个实例不是线程安全的，因此不能被共享。绝对不能将SqlSession的实例引用放在一个类的静态域，甚至一个类的实例变量，也绝不能将SqlSession实例的引用放在任何类型的托管作用域中，比如Servlet的HttpSession。应该是在每次收到一个Http请求时，打开一个SqlSession，返回响应后就关闭它（这个关闭很重要！）。

## 主配置文件

### 结构：

- configuration
  - properties
  - settings
  - typeAliases
  - typeHandlers
  - objectFactory
  - plugins
  - environments
    - environment
      - transactionManager
      - dataSource
  - databaseIdProvider
  - mappers

### 属性（Properties）

预先定义一系列属性，然后在需要的地方使用，类似定义变量。

```xml
<!--resource指定外部引用的内容-->
<properties resource="data.properties">
    <property name="driver" value="com.mysql.cj.jdbc.Driver"/>
    <!--使用时${driver}-->
</properties>
```

如果一个属性出现多次，加载顺序为

1. 首先读取properties标签下指定的内容
2. 根据resource指定的外部文件去读取对应内容，并覆盖同名属性
3. 最后读取作为方法参数传递的属性，并覆盖之前读取过的同名属性。（代码中传递的参数）

#### 指定默认值

```xml
<properties resource="org/mybatis/example/config.properties">
  <property name="org.apache.ibatis.parsing.PropertyParser.enable-default-value" value="true"/> <!-- 启用默认值特性 -->
</properties>

<dataSource type="POOLED">
  <property name="username" value="${username:ut_user}"/> <!-- 如果属性 'username' 没有被配置，'username' 属性的值将为 'ut_user' -->
</dataSource>
```

### 设置（Settings）

对mybatis的默认行为进行设置，详见[mybatis设置](https://mybatis.org/mybatis-3/zh/configuration.html#settings)

```xml
<settings>
  <setting name="cacheEnabled" value="true"/>
  <setting name="lazyLoadingEnabled" value="true"/>
  <setting name="multipleResultSetsEnabled" value="true"/>
  <setting name="useColumnLabel" value="true"/>
  <setting name="useGeneratedKeys" value="false"/>
  <setting name="autoMappingBehavior" value="PARTIAL"/>
  <setting name="autoMappingUnknownColumnBehavior" value="WARNING"/>
  <setting name="defaultExecutorType" value="SIMPLE"/>
  <setting name="defaultStatementTimeout" value="25"/>
</settings>
```

### 类型别名（typeAliases）

起一个别名，仅用于XML配置。

```xml
<typeAliases>
  <typeAlias alias="Author" type="domain.blog.Author"/>
  <typeAlias alias="Blog" type="domain.blog.Blog"/>
  <typeAlias alias="Comment" type="domain.blog.Comment"/>
  <typeAlias alias="Post" type="domain.blog.Post"/>
  <typeAlias alias="Section" type="domain.blog.Section"/>
  <typeAlias alias="Tag" type="domain.blog.Tag"/>
</typeAliases>
```

也可以指定一个包名，这样会搜索需要的java Bean对象，在没有指定注解时，会类名作为别名（首字母小写）。

```xml
<typeAliases>
  <package name="domain.blog"/>
</typeAliases>
```

```java
@Alias("author")//将domain.blog.Author别名为author
public class Author {
    ...
}
```

### 类型处理器（typeHandlers）

在设置PreparedStatement中的参数或从结果集中取出值时，mybatis会用类型处理器将获取到的值转换成合适的java对象，详见[类型处理器](https://mybatis.org/mybatis-3/zh/configuration.html#typeHandlers)

可以重写已有的类型处理器或创建自己的类型处理器来处理不支持的类型。实现`TypeHandler`接口或继承`BaseTypeHandler`

```java
// 指定其对应的JDBC类型（可选）
@MappedJdbcTypes(JdbcType.VARCHAR)
public class ExampleTypeHandler extends BaseTypeHandler<String> {

  @Override
  public void setNonNullParameter(PreparedStatement ps, int i, String parameter, JdbcType jdbcType) throws SQLException {
    ps.setString(i, parameter);
  }

  @Override
  public String getNullableResult(ResultSet rs, String columnName) throws SQLException {
    return rs.getString(columnName);
  }

  @Override
  public String getNullableResult(ResultSet rs, int columnIndex) throws SQLException {
    return rs.getString(columnIndex);
  }

  @Override
  public String getNullableResult(CallableStatement cs, int columnIndex) throws SQLException {
    return cs.getString(columnIndex);
  }
}
```

```xml
<typeHandlers>
  <typeHandler handler="org.mybatis.example.ExampleTypeHandler"/>
  <!--也可以指定包名，mybatis会自动查找-->
  <package name="org.mybatis.example"/>
</typeHandlers>
```

### 环境配置（environments）

Mybatis中可以配置多种环境，一种环境即为一种数据库，例如开发时开发、测试和生产环境下都对应不同的配置。不过每个SqlSessionFactory只能选择一种环境。但可以有多个SqlSessionFactory。

```java
//第二个参数即为指定环境，不指定时会使用默认环境
SqlSessionFactory factory =
                new SqlSessionFactoryBuilder().build(is,"mysql");
```

```xml
<environments default="development">
  <environment id="development">
    <transactionManager type="JDBC">
      <property name="..." value="..."/>
    </transactionManager>
    <dataSource type="POOLED">
      <property name="driver" value="${driver}"/>
      <property name="url" value="${url}"/>
      <property name="username" value="${username}"/>
      <property name="password" value="${password}"/>
    </dataSource>
  </environment>
</environments>
```

#### 事务管理器（transactionManager）

mybatis中有两种事务管理器

- JDBC 直接使用了JDBC的提交和回滚功能
- MANAGED 从不提交或回滚，而是让容器来管理事务。

注：如果项目中有使用Spring，则可以不配置事务管理器，因为Spring会使用自带的管理器。

#### 数据源（datasource）

有三种内建的数据源类型

- UNPOOLED 每次请求是都会打开一个新的连接。
- POOLED 利用池的概念，会有一个数据库链接池。
- JNDI 为了能在如EJB或应用服务器这类容器中使用，容器可以集中或在外部配置数据源，然后放置一个JNDI上下文的数据源引用，如使用Tomcat的JNDI。

### 数据库厂商标识（砌筑baseIdProvider）

mybatis可以根据不同的数据库执行不同的SQL语句，这种实现是基于映射语句中的`databaseId`属性进行区分。

```xml
<!--开启多厂商支持-->
<databaseIdProvider type="DB_VENDOR">
  <!--设置别名-->
  <property name="SQL Server" value="sqlserver"/>
  <property name="DB2" value="db2"/>
  <property name="Oracle" value="oracle" />
</databaseIdProvider>
```

### 映射器（mappers）

定义映射的SQL语句，告诉Mybatis到哪里去找到这些语句。

```xml
<!-- 使用相对于类路径的资源引用 -->
<mappers>
  <!--类路径下查找对应的配置文件-->
  <mapper resource="org/mybatis/builder/AuthorMapper.xml"/>
  <!--指定的类-->
  <mapper class="org.mybatis.builder.PostMapper"/>
  <!--包内的所有接口-->
  <package name="org.mybatis.builder"/>
</mappers>
```

## 映射文件

包含的元素

- cache 该命名空间的缓存配置
- cache-ref 引用其它命名空间的缓存配置
- resultMap 描述如何从结果集中加载对象
- sql 可被重用的SQL语句
- insert
- update
- delete
- select

### select

查询语句，详见[select](https://mybatis.org/mybatis-3/zh/sqlmap-xml.html#select)

```xml
<select id="selectPerson" parameterType="int" resultType="hashmap">
  SELECT * FROM PERSON WHERE ID = #{id}
</select>
```

### sql

定义可重用的SQL语句，以便在其它语句中使用。

```xml
<sql id="userColumns"> ${alias}.id,${alias}.username,${alias}.password </sql>
<select id="selectUsers" resultType="map">
  select
    <include refid="userColumns"><property name="alias" value="t1"/></include>  	from some_table t1
</select>
```

### 字符串替换

默认情况下，使用`#{}`等同于PreparedStatement中的占位符`?`然后通过占位符设置参数，不过有时候仅仅动态的替换SQL语句中的一个字符串（不将其转为占位符，而是直接替换），则可以使用`${}`，不过这种方式存在SQL注入隐患。

```java
@Select("select * from user where ${column} = #{value}")
User findByColumn(@Param("column") String column, @Param("value") String value);
```

### 结果映射（resultMap）

将结果集中的字段映射到一个java对象中。

例如，有一个POJO

```java
public class User{
    int id;
    String name;
}
```

```xml
<!--当名称不同时可以手动设置别名-->
<select id="selectUsers" resultType="com.someapp.model.User">
  select id, username as 'name'
  from some_table
  where id = #{id}
</select>
```

在这种情况下，Mybatis会自动进行映射，将读取出的数据映射到一个Java对象上。

当然，Mybatis也支持手动定义映射规则，这时则需要使用`resultMap`。

```xml
<!--property指定类中属性的名称，column则数据库中字段的名称-->
<resultMap id="userMap" type="User">
    <!--id表示这个属性是主键-->
    <id property="id" column="id"/>
    <result property="name" column="username"/>
    <!--<result javaType="对应的java类型" jdbcType="对应的JDBC类型" typeHandler="类型处理器"/> -->
</resultMap>

<select id="selectUsers" resultMap="userMap">
  select id, username
  from some_table
  where id = #{id}
</select>
```

#### 高级映射

由于表与表之间存在一对一，一对多和多对多这样的复杂关系，因此映射也不只是简简单单的字段名与属性对应。在resultMap中包含以下标签，以应对表与表之间的复杂关系。

- constructor 在实例化类时，注入结果到构造方法中
  - idArg 主键参数
  - arg 其它结果参数
- id 这个属性是主键
- result 字段和java对象属性这间的关联
- association 一个复杂类型的关联，外键
- collection 一个复杂类型的集合
- discriminator 使用结果值来决定使用哪个`resultMap`

##### constructor

一个java对象

```java
public class User{
    User(Integer id, String name, int age){}
}
```

```xml
<constructor>
	<idArg column="id" javaType="int" name="id"/>
    <arg column="age" javaType="_int" name="age"/>
    <arg column="name" javaType="String" name="name"/>
</constructor>
```

##### association（关联）

处理一对一关系，可以将一系列字段映射成一个对象（连接查询），或者进行Select嵌套查询。

例如：一个Blog对应一个Author。

- 嵌套查询

  这种方式虽然简单但对于大型数据集上表现不佳，这被“N+1”问题，首先单独执行一个SQL语句来获取一个结果的列表（就是”+1“，如这里的Blog），然后对这个列表中的每一条数据又去执行另一条SQL语句来加载详细信息（就是”N“，如这里去获取Blog中的Author）。这将会导致大量SQL语句被执行，不过Mybatis对这种查询会采用延迟加载的方式。如这里会先加载出Blog，而对于Blog中的Author则会在使用时才会去加载。

  ```xml
  <resultMap id="blogResult" type="Blog">
      <!--Blog和Author是一对一关系，会将author_id字段的内容用于去查询Author对象，并将结果设置给author属性-->
      <!--可以理解为将column中指定的列作为参数，传递给selectAuthor-->
    <association property="author" column="author_id" javaType="Author" select="selectAuthor"/>
  </resultMap>
  <select id="selectBlog" resultMap="blogResult">
    SELECT * FROM BLOG WHERE ID = #{id}
  </select>
  <select id="selectAuthor" resultType="Author">
    SELECT * FROM AUTHOR WHERE ID = #{id}
  </select>
  ```

- 连接查询

  ```xml
  <resultMap id="authorMap" type="author">
          <id property="id" column="author_id"/>
          <result property="name" column="author_name"/>
  </resultMap>
  <resultMap id="blogMap" type="blog">
          <id property="id" column="blog_id"/>
          <result property="name" column="blog_name"/>
          <!--属性author根据authorMap的映射规则 解析-->
          <association property="author" column="blog_author_id"
                       javaType="author" resultMap="authorMap"/>
      <!--也可以直接在内部定义映射规则-->
      <!--<association property="author" column="blog_author_id"
                       javaType="author">
              <id property="id" column="author_id"/>
              <result property="name" column="author_name"/>
          </association>-->
  </resultMap>
  <!--由于存在重名字段，所以每个字段因该设置别名，且不能使用表名.列名的格式-->
  <select id="getBlog" resultMap="blogMap">
          select b.id        as blog_id,
                 b.name      as blog_name,
                 b.author_id as blog_author_id,
                 a.id        as author_id,
                 a.name      as author_name
          from blog b
                   left outer join author a on a.id = b.author_id
          where b.id = #{id}
  </select>
  ```

  `columnPrefix`使用，在连接多个表时，可以指定别名来标识不同表中重复的列名，通过指定`columnPrefix`列名的前缀，可以将带有指定前缀的列映射到一个resultMap中。

  例如：每一个Blog除了有一个author外，还有一个共同作者（co-author）

  ```xml
  <resultMap id="authorMap" type="author">
      <id property="id" column="author_id" javaType="_int"/>
      <result property="name" column="author_name" javaType="String"/>
  </resultMap>
  <resultMap id="blogMap" type="blog">
      <id property="id" column="blog_id" javaType="int"/>
      <result property="name" column="blog_name" javaType="String"/>
      <association property="author" column="blog_author_id"
                   javaType="author" resultMap="authorMap"/>
      <!--将带有co_前缀的列也由authorMap映射，映射时不用再考虑前缀-->
      <association property="coAuthor" column="blog_co_author_id"
                   javaType="author" resultMap="authorMap"
                   columnPrefix="co_"/>
  </resultMap>
  <select id="getBlog" resultMap="blogMap">
          select b.id           as blog_id,
                 b.name         as blog_name,
                 b.author_id    as blog_author_id,
                 b.co_author_id as blog_co_author_id,
                 a.id           as author_id,
                 a.name         as author_name,
                 ca.id          as co_author_id,
                 ca.name        as co_author_name
          from blog b
                   left join author a on a.id = b.author_id
                   left join author ca on ca.id = b.co_author_id
          where b.id = #{id}
  </select>
  ```

##### collection（集合）

用于处理一对多关系，而对于多对多关系，等价于一对多关系。

例如：一个Blog只有一个Author，但可以有多个Post(文章)。所以在Blog中会有一个`List<Post> posts`属性。不过不同于数据库中表的设计，对于一对多关系，会在”多“的一方中添加一个外键，去引用”少“的一方的主键。

和association相同，collection也可以使用连接查询或select嵌套查询。

- select嵌套查询

  ```xml
  <resultMap id="blogMap" type="blog">
      <id property="id" column="id"/>
      <association property="author" column="author_id"
                       javaType="author" select="getAuthor"/>
      <association property="coAuthor" column="co_author_id"
                       javaType="author" select="getAuthor"/>
       <collection property="posts" column="id"
                      javaType="ArrayList" select="getPost"
                      ofType="post"/>
  </resultMap>
  <select id="getBlog" resultMap="blogMap">
          select *
          from blog
          where id = #{id}
  </select>
  <select id="getAuthor" resultType="author">
          select *
          from author
          where id = #{id}
  </select>
  <select id="getPost" resultType="post">
          select *
          from post
          where blog_id = #{id}
  </select>
  ```

  使用方式和association类似，只是多了一个`ofType`属性，用来指定集合中元素的类型。

- 连接查询

  ```xml
  <resultMap id="authorMap" type="author">
      <id property="id" column="author_id" />
      <result property="name" column="author_name" />
  </resultMap>
  <resultMap id="blogMap" type="blog">
      <id property="id" column="blog_id" />
      <result property="name" column="blog_name" />
      <association property="author" column="blog_author_id" 
                   javaType="author" resultMap="authorMap" />
      <association property="coAuthor" column="blog_co_author_id" 
                   javaType="author" resultMap="authorMap"
                   columnPrefix="co_" />
      <collection property="posts" ofType="post">
          <id property="id" column="post_id" />
          <result property="title" column="post_title" />
          <result property="content" column="post_content" />
      </collection>
  </resultMap>
  <select id="getBlog" resultMap="blogMap">
      	select b.id           as blog_id,
                 b.name         as blog_name,
                 b.author_id    as blog_author_id,
                 b.co_author_id as blog_co_author_id,
                 a.id           as author_id,
                 a.name         as author_name,
                 ca.id          as co_author_id,
                 ca.name        as co_author_name,
                 p.id           as post_id,
                 p.title        as post_title,
                 p.content      as post_content,
                 p.blog_id      as post_blog_id
          from blog b
                   left join author a on a.id = b.author_id
                   left join author ca on ca.id = b.co_author_id
                   left join post p on p.blog_id = b.id
          where b.id = #{id}
  </select>
  ```

### 鉴别器（discriminator）

类似于java的switch语句，根据结果集返回不同的值来选择不同的resultMap去映射。

```xml
<resultMap id="vehicleResult" type="Vehicle">
  <id property="id" column="id" />
  <result property="vin" column="vin"/>
  <result property="year" column="year"/>
  <result property="make" column="make"/>
  <result property="model" column="model"/>
  <result property="color" column="color"/>
  <!--根据vechicle_type具体的值去选择不同的映射器-->
  <discriminator javaType="int" column="vehicle_type">
    <case value="1" resultMap="carResult"/>
    <case value="2" resultMap="truckResult"/>
    <case value="3" resultMap="vanResult"/>
    <case value="4" resultMap="suvResult"/>
  </discriminator>
</resultMap>
<resultMap id="carResult" type="Car">
  <result property="doorCount" column="door_count" />
</resultMap>
```

### 缓存

Mybatis中内置了一级缓存和二级缓存，默认只启用一级缓存，仅对一个会话中的数据进行缓存。可在映射文件中添加以下标签启用二级缓存，详见[cache](https://mybatis.org/mybatis-3/zh/sqlmap-xml.html#cache)

```xml
<cache/>
```

## 动态SQL

根据参数的值动态的拼接SQL语句，并且可以使用OGNL表达式。

### 条件（if）

```xml
<select id="findActiveBlogLike" resultType="Blog">
  SELECT * FROM BLOG WHERE state = ‘ACTIVE’
  <!--如果test中的内容为真，则会拼接标签内的SQL语句-->
  <if test="title != null">
    AND title like #{title}
  </if>
  <if test="author != null and author.name != null">
    AND author_name like #{author.name}
  </if>
</select>
```

### 选择（choose，when，otherwise）

类似于switch语句，从多个条件中选择一个来使用，从上到下依次匹配，匹配成功则不再继续。

```xml
<select id="findActiveBlogLike" resultType="Blog">
  SELECT * FROM BLOG WHERE state = ‘ACTIVE’
  <choose>
    <when test="title != null">
      AND title like #{title}
    </when>
    <when test="author != null and author.name != null">
      AND author_name like #{author.name}
    </when>
    <!--上面都匹配失败时采用otheriwse-->
    <otherwise>
      AND featured = 1
    </otherwise>
  </choose>
</select>
```

### where

有如下例子

```xml
<select id="findActiveBlogLike"
     resultType="Blog">
  SELECT * FROM BLOG
  WHERE
  <if test="state != null">
    state = #{state}
  </if>
  <if test="title != null">
    AND title like #{title}
  </if>
  <if test="author != null and author.name != null">
    AND author_name like #{author.name}
  </if>
</select>
```

如果if标签中没有一个被匹配，则最终的SQL语句会变成

```sql
SELECT * FROM BLOG WHERE
```

而当`state !=null`不成立时，又会变成

```SQL
SELECT * FROM BLOG WHERE AND xxx
```

在这样场景中，mybatis提供了一个where标签，它只在子元素有内容时去拼接WHERE子句，并且若子句开头是”and”或“or”，也会将它们去除。

```xml
<select id="findActiveBlogLike" resultType="Blog">
  SELECT * FROM BLOG
  <where>
    <if test="state != null">
         state = #{state}
    </if>
    <if test="title != null">
        AND title like #{title}
    </if>
    <if test="author != null and author.name != null">
        AND author_name like #{author.name}
    </if>
  </where>
</select>
```

### 动态更新语句（set）

```xml
<update id="updateAuthorIfNecessary">
  update Author
    <set>
      <if test="username != null">username=#{username},</if>
      <if test="password != null">password=#{password},</if>
      <if test="email != null">email=#{email},</if>
      <if test="bio != null">bio=#{bio}</if>
    </set>
  where id=#{id}
</update>
```

### 迭代（foreach）

```xml
<select id="selectPostIn" resultType="domain.blog.Post">
  SELECT *
  FROM POST P
  WHERE ID in
  <!--item为值，index为键或索引-->
  <!--可以迭代list,set,map,数组-->
  <foreach item="item" index="index" collection="list"
      open="(" separator="," close=")">
        #{item}
  </foreach>
</select>
```

## 使用注解

详见[mybatis文档](https://mybatis.org/mybatis-3/zh/java-api.html#sqlSessions)

## 整合Spring（mybatis-spring）

### Maven依赖

```xml
<dependency>
  <groupId>org.mybatis</groupId>
  <artifactId>mybatis-spring</artifactId>
  <version>2.0.6</version>
</dependency>
```

### 配置SqlSessionFactory

使用SqlSessionFactoryBean来创建SqlSessionFactory，并且必须指定一个数据源，这个数据源也可以使用其它框架来获取（如c3p0）。

注意：这里使用的SqlSessionFactoryBean是一个Spring的`FactoryBean`接口实现类（详见[自定义工厂bean](https://docs.spring.io/spring-framework/docs/current/reference/html/core.html#beans-factory-extension-factorybean)），最终创建的对象并不是一个SqlSessionFactoryBean，而是其实现的方法`getObject`的返回值，在这里则是一个`SqlSessionFactory`

```xml
<!--使用mybatis内置的数据源-->
<bean id="dataSource" class="org.apache.ibatis.datasource.pooled.PooledDataSource">
    <property name="driver" value="com.mysql.cj.jdbc.Driver"/>
    <property name="url" value="jdbc:mysql://localhost:3306/mybatis?serverTimezone=Asia/Shanghai"/>
    <property name="username" value="root"/>
    <property name="password" value="root"/>
</bean>
<bean id="sqlSessionFactory" class="org.mybatis.spring.SqlSessionFactoryBean">
    <property name="dataSource" ref="dataSource"/>
    <!--指定Mybatis的xml配置文件，其中和环境有关的配置会被忽略-->
    <!--<property name="configLocation" value="classpath:"/>-->
    <!--指定xml映射器文件的位置，使用注解时可以不配置,可以使用Ant风格路径-->
    <property name="mapperLocations" value="classpath:blog-mapper.xml"/>
</bean>
```

### 注入Mapper

注意：这里同SqlSessionFactoryBean一样，MapprFactoryBean同样是`FactoryBean`的实现类。

```xml
<bean id="blogMapper" class="org.mybatis.spring.mapper.MapperFactoryBean">
     <property name="mapperInterface" value="mapper.IBlogMapper"/>
     <property name="sqlSessionFactory" ref="sqlSessionFactory"/>
</bean>
```

### 事务配置

可以借助Spring中的`DataSourceTransactionManager`来管理事务。

```xml
<bean id="transactionManager"
          class="org.springframework.jdbc.datasource.DataSourceTransactionManager">
    <!--需要SqlSessionFactoryBean中传入dataSource和这个是同一个-->
    <constructor-arg name="dataSource" ref="dataSource"/>
</bean>
```

### 使用SqlSession

整合Spring后，将不使用SqlSessionFatory，因为已经直接将需要的mapper注入到Spring的IOC容器中。但可以通过`SqlSessionTemplate`来获取`SqlSession`。

```xml
<bean id="sqlSession" class="org.mybatis.spring.SqlSessionTemplate">
    <constructor-arg name="sqlSessionFactory" ref="sqlSessionFactory"/>
</bean>
```

### 发现映射器

可以不用一个一个的去手动注册映射器，而是让mybatis-spring对类路径扫描并发现它们。

- XML配置方式

  在Spring配置文件中使用`<mybatis:scan/>`，类似Spring的`<context:component-scan/>`，base-package指定要扫描的包。

  ```xml
  <beans xmlns="http://www.springframework.org/schema/beans"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xmlns:mybatis="http://mybatis.org/schema/mybatis-spring"
         xsi:schemaLocation="http://www.springframework.org/schema/beans
          http://www.springframework.org/schema/beans/spring-beans.xsd
          http://mybatis.org/schema/mybatis-spring
          https://mybatis.org/schema/mybatis-spring.xsd">
  <mybatis:scan base-package="mapper"/>
  ```

  被发现的映射器会使用Spring的默认命名规则（类名首字母小写），可以使用注解`@Component`自定义名称。

  在这里并没有指定SqlSessionFactory，因为它会自动注入一个MapperFactoryBean（同样也会为这个对象自动注入一个SqlSessionFactory），但当有多个数据源时（也意味着存在多个SqlSessionFactory），则需要指定所使用的SqlSessionFactory。`<mybatis:scan base-package="mapper" factory-ref="sqlSessionFactory"/>`。

- 注解方式

  在Spring的配置类上使用`@MapperScan`，等价于上一种方式，只是以注解方式
  
  ```java
  @Configuration
  @MapperScan(basePackages = "mapper")
  public class Application {
  
  }
  ```
  
  
  
  



