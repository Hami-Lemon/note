# Spring Boot（未完）

Spring Boot用来简化Spring应用开发，约定大于配置

## 微服务

一种架构风格，一个应用应该是一组小型服务，可以通过HTTP的方式进行互通

## Spring Initializer快速创建Spring Boot项目

![1](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113138.png)
![2](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113142.png)

- 主程序已经自动生成，只需要实现自己的逻辑
- resources文件夹中目录
  - static：保存所有的静态资源
  - templates: 保存所有的模板页面
  - application：spring boot 的配置文件

## Spring Boot配置

springboot使用一个全局的配置文件

- `application.properties`
- `application.yml`

配置文件的作用:修改spring Boot的默认配置

### YAML

`k: v`表示一对键值对(空格必须有)  
以空格的缩进控制层级关系  
属性和值大小写敏感

```yml
server:
 port: 8888
 path: /hello
```

#### 值的写法

- 字面量:普通的值(数字,字符串,布尔)  
  `k: v`:直接写  
  字符串默认不用加单引号或者双引号  
  `""`:双引号,会转义字符串里的特殊字符
    `name: "vvv\nsas"`  输出:vvv换行nsas
  `''`:单引号,不会转义特殊字符
- 对象,Map:

  ```yml
  friends:
    lastName: zhangsan
    age: 20
  ```

  行内写法:

  ```yml
  friends: {lastName: zhangsan, age: 20}
  ```

- 数组
用`-`表示数组中的一个元素

    ```yml
    pets:
    - cat
    - dog
    - pig
    ```

    行内写法:

    ```yml
    pets: [cat, dog, pig]
    ```

- 实例

java对象

```java
@ConfigurationProperties(prefix = "person")
public class Person{
    private String lastName;
    private Integer age;
    private Boolean boss;
    private Date birth;

    private Map<String, Object> maps;
    private List<Object> lists;
    private Dog dog;
}

class Dog{
    private String name;
    private Integer age;

}
```

yml表示

```yml
person:
  lastName: zhangsan
  age: 20
  boss: false
  birth: 2018/12/12
  maps: { k1: v1, k2: v2}
  lists:
    - lisi
    - zhaoliu
  dog:
    name: 狗
    age: 2
```

将yml中配置的属性的值,映射到对象中
`@ConfigurationProperties`将本类中所有属性和配置文件中的相关配置进行绑定

### `@value`和`@ConfigurationProperties`比较

|                    | @ConfigurationProperties | @Value                          |
| ------------------ | ------------------------ | ------------------------------- |
| 功能               | 批量注入配置文件中的属性 | 一个一个指定                    |
| 松散绑定(松散语法) | 支持                     | 不支持                          |
| SpEL               | 不支持                   | 支持                            |
| JSR303数据校验     | 支持                     | 不支持                          |
| 复杂类型封装       | 支持                     | 不支持,只能注入普通类型和字符串 |

### `@PropertySource`和`@importResource`

- `PropertySource`

导入配置文件

```java
@PropertySource(value={"classpath:xxx.properties"})//加载指定的配置文件
@ConfigurationProperties(prefix="person")
public class Person{
    
}
```

- `importResource`

  导入Spring的配置文件

```xml
<bean id = "xxx" class = "xxx"/>
```

Spring Boot不会自动加载Spring的配置文件

使用`@ImportResource`注解导入Spring的配置文件@ImportResource(locations = {"classpath:beans.xml"})

```java
@SpringBootApplication
public class xxxxx{  
}
```

- spring boot推荐的方式

  给容器中添加组件的方式--全注解方式

  ```java
  //声明当前类是一个配置类
  @Configuration
  public class MyAppConfig{
      //会将该方法的返回值添加到IOC容器中,默认id为方法名,也可以自己指定
      @Bean
      public HelloService helloService(){
          return new HelloService();
      }
  }
  ```

  

### 配置文件占位符

```properties
person.name=张三${random.uuid}
person.age=${random.int}
person.dog.name=${person.name}_dog
```

