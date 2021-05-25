java持久层操作，JDBC以及Mybatis框架的使用
<!--more-->

## JDBC

1. 导入jar包（以mysql为例)

2. 初始化驱动类
`Class.forname("com.mysql.jdbc.Driver")`
3. 建立连接

```JAVA
Connection c = DriverManager
                .getConnection("jdbc:mysql://ip:3306/库名,"root","password");
```

4. 创建Statement
使用Steatement执行sql语句（使用java.sql.Statement)
`Statement s = c.createStatement();`
5. 执行sql语句
```JAVA
String sql = "sql";
s.execute(sql);
```

6. 关闭连接
```JAVA
//先关闭Steatement，再关闭Connection
//可使用try-with-resource自动关闭
s.close();
c.close();
```

7. CRUD增删改查操作
```JAVA
//增删改与上面操作相同，只是sql语句不同
//查
Class.forname("com.mysql.jdbc.Driver");
Connection c = DriverManager.getConnection("jdbc:mysql://ip:3306/库","root","password");
Statement s = c.createStatement();
String sql = "SELECT * from ";
ResultSet rs = s.executeQuery(sql);//获取查询的结果集
while(rs.next()){
    int id = rs.getInt("id");//使用字段名获取
    String name = rs.getString(2);//使用字段顺序，从1开始
}
//Steatement关闭时，ResultSet会自动关闭
```

8. 使用PreparedStatement
```JAVA
Class.forname("xxx");
String sql = "insert into xx values(null,?,?,?)";//?为占位符
Connection c = DriverManager.getConnection("xx","root","password");
PreparedStatement ps = c.prepareStatement(sql);
//设置参数
ps.setString(1,"xxx");
ps.setFloat(2,xxx);
ps.setInt(3,xxx);
ps.execute();//执行
```
9. execute()与executeUpdate()
    - 不同点
        1. execute 可以查询，使用getResultSet获取结果集,executeUpdate不能查询
        2. execute返回boolean，true为执行的是查询语句
        3. executeUpdate返回int，表示有多少条数据受影响

10. 事务
```JAVA
//关闭自动提交
c.setAutoCommit(false);
//手动提交
//要么所有的操作都执行，要么都不执行
s.commit();
```

## Mybatis框架（持久层框架）
内部封装了jdbc，只需要关注sql语句本身
ORM思想：对象关系映射，把数据库表和实体类的属性对应起来
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114219.png)

- 创建一个javabean作为实体类，并实现`Serializable`接口
- 创建`Dao`接口
- 在resources目录下建立配置文件的xml文件
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114228.png)
```xml
<?xml version="1.0" encoding="UTF-8" ?>

<!DOCTYPE configuration
        PUBLIC "-//mybatis.org//DTD Config 3.0//EN"
        "http://mybatis.org/dtd/mybatis-3-config.dtd">
```
- 建立`Dao`的xml文件
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114241.png)
```xml
<?xml version="1.0" encoding="UTF-8" ?>

<!DOCTYPE mapper
        PUBLIC "-//mybatis.org//DTD Config 3.0//EN"
        "http://mybatis.org/dtd/mybatis-3-mapper.dtd">

<mapper namespace="com.hamilemon.dao.IUserDao">
```

### 入门案例
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114245.png)
1. 读取配置文件
2. 创建SqlSessionFactory工厂
3. 创建SqlSession
4. 创建Dao接口的代理对象
5. 执行dao中的方法
6. 释放资源
注意: 配置中要指定封装的实体类`resultType`

- 用注解配置
1. 主配置文件中,mapper中改为`class="类名"`
2. Dao方法上添加注解`@Select("sql语句")`;
- 总结
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114251.png)
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114255.png)

### Mybatis中的CRUD
见示例

### Mybatis中的参数
1. 传递简单数据
2. 传递pojo对象
 使用ognl表达式解析对象字段的值,#{}或${}括号中的值为pojo属性名称
 - ognl表达式
    Object Graphic Navigation Language(对象图导航语言)
    通过对象的的取值方法获取数据,在写法上把get给省略了
    例:类中写法:user.getName()
        OGNL写法:user.name
 3. 传递pojo包装对象

 ### Mybatis结果类型的封装
 pojo对象的属性名与数据库字段保持一致
 1. 输出基本类型
 2. 输出pojo对象
 3. 输出pojo列表

 ### Mybatis的配置文件
 1. 可在主配置文件中定义一个<properties>标签中,将jdbc的配置写入配置文件中
 2. 配置别名（typeAliases）
 ```xml
 <typeAliases>
    <!-- 指定别名，使用时可直接写别名，放在主配置文件中 -->
    <typeAlias type="com.hamilemon.xxx" alias="user"/>
    <!-- 指定要配置别名的包，指定后，该包下的实体类都会注册别名，并且类名就是别名，不分大小写 -->
    <package name=""/>
 </typeAliases>
 ```

### Mybatis中的连接池以及事务控制
1. 连接池以及事务
    - mybatis连接池提供了3种配置方式
        + 主配置的dataSource标签，type属性就是采用何种连接池
            - type的取值：
                - POOLED 传统的javax.sql.DataSource规范
                - UNPOOLED 传统的获取连接的方式，没有池的思想
                - JNDI 采用服务器提供的JNDI技术实现，来获取DataSource对象
    - 事务
    
2. 配置的动态sql语句使用
    - `<if>`
        做条件判断
    - `<where>`
    - `<foreach>`
        实现遍历功能
        ![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114301.png)
    - `<sql>`
        抽取重复sql语句
3. 多表操作
    - 一对多
        `<collection>`添加集合映射,对于相同的字段名需要设置别名
    - 一对一
        `<resultMap>`中添加`<assocation>`
    - 多对一
        看成一对一
    - 多对多
        建立中间表
### Mybatis延迟加载
一对多中
在查询时,用户下的账户信息应该是什么时候使用,什么时候查询
在查询账户时,账户的所属用户信息应该一起查询出来
- 延迟加载
    在真正使用数据时才发起查询,不用的时候不查询,叫延迟加载(按需加载,懒加载)
- 立即加载    
    不管用不用,调用方法就发起查询
- 四种对应关系中
    - 一对多,多对多:采用延迟加载
    - 多对一,一对一:采用立即加载
- 一对一中的延迟加载
    主配置文件中
    ```xml
    <settings>
        <setting name = "lazyLoadingEnabled" value="true"/>
        <setting name="aggressiveLazyLoading" value="false"/>
    </settings>
    ```
    映射配置文件中
    ```xml
    <resultMap>
    <!--一对一的关系映射 延迟加载-->
        <!--select属性的内容 查询用户的唯一标识-->
        <!--column 用户根据id查询时，所需要的参数的值-->
        <association property="user" column="uid" javaType="user"
                     select="com.hamilemon.dao.IUserDao.findById"/>
    </resultMap>
    ```
- 一对多的延迟加载
    ```xml
    <resultMap>
    <!-- 一对多 -->
    <collection property="accounts" column="id" ofType="account" 
                    select="com.hamilemon.dao.IAccountDao.findByUid"/>
    </resultMap>
    ```

### Mybatis缓存
存在于内存中的临时数据,减少和数据库的交互次数

适用于缓存的数据:经常查询并且不经常改变的,数据的正确与否对最终结果影响不大

不适用于缓存的数据:经常改变得数据,数据的正确与否对最终结果影响很大的

- 一级缓存
Mybatis中SqlSession对象的缓存,自动开启

执行查询后,查询的结果会同时存入到SqlSession为我们提供的一块区域中,

再次查询同样的数据,mybatis会先去SqlSession中查询是否有

当调用修改,删除,commit(),close(),clearCache(),都会清空一级缓存

- 二级缓存
  Mybatis中的SqlSessionFactory对象的缓存,由同一个SqlSessionFactory对象建立的SqlSession共享其缓存
  ![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525114312.png)
    - 使用
        1. 主配置文件中
        ```
        <setting name="cacheEnabled" value="true"/>
        ```
        2. 映射配置文件中
        ```
        <cache/>
        3.
        在查询的标签上加上'useCache="true"'属性
        ```
### Mybatis注解开发
针对crud有四个注解
`@Select,@Insert,@Update,@Delete`
- 配置属性名,与字段名的对应关系
```JAVA
@Select("xxxx")
@Results(id="userMap",
    value={@Result(id=true,column="id",property="userId"),
    @Result(column="name",property="userName")
})
List<User> findAll();

// 引用上面的对应关系
@ResultMap("userMap")
User findById();
```

- 多表查询
    - 一对一,多对一
    ```JAVA
    /**
     * 查询所有账户，并且获取所属的用户信息
     * @return
     */
    @Select("select * from account")
    @Results(id = "accountMap",value={
            @Result(id=true,column = "id",property = "id"),
            @Result(column = "uid",property = "uid"),
            @Result(column = "money", property = "money"),
            @Result(property = "user",column = "uid",
                    one = @One(select = "com.hamilemon.dao.IUserDao.findById",fetchType = FetchType.EAGER))
    })
    //通过指定fetchType的值,确定是否使用懒加载
    List<Account> findAll();
    
    ```
    - 一对多
    ```java
    /**
     * 查询所有用户并获取旗下的所有账户
     * @return 所有用户
     */
    @Select("select * from user")
    @Results(id="userMap",value = {
            @Result(id=true,column = "id",property = "id"),
            @Result(column = "name",property = "name"),
            @Result(column = "age",property = "age"),
            @Result(column = "id",property = "accounts",
            many = @Many(select = "com.hamilemon.dao.IAccountDao.findByUid",
            fetchType = FetchType.LAZY))
    })
    List<User> findAll();
    ```
