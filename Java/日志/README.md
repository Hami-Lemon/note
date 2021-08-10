# 日志

## Log4j2使用

### Maven依赖

详见:[log4j2 maven](http://logging.apache.org/log4j/2.x/maven-artifacts.html)

```xml
<dependencies>
  <dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-api</artifactId>
    <version>2.14.1</version>
  </dependency>
  <dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-core</artifactId>
    <version>2.14.1</version>
  </dependency>
</dependencies>
```

### Log4j2 API

```java
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
 
public class HelloWorld {
    private static final Logger logger = LogManager.getLogger("HelloWorld");
    public static void main(String[] args) {
        logger.info("Hello, World!");
        //使用占位符
        logger.debug("Logging in user {} with birthday {}", user.getName(), 
                     user.getBirthdayCalendar());
        //参数格式化
        logger.debug("Logging in user %s with birthday %s", user.getName(), 
                     user.getBirthdayCalendar());
        //labbda表达式支持
        logger.debug("lambda test {}", () -> "message");
    }
}
```

### 配置

Log4j2会自动加载配置文件，其查找顺序为（哪一个文件先被查找到就使用哪一个）：

1. 类路径下的log4j2-test.properties
2. 类路径下的log4j2-test.yaml或log4j-text.yml
3. 类路径下的log4j2-test.json或log4j2-test.jsn
4. 类路径下的log4j2-text.xml
5. 然后查找路径下无`-text`后缀的文件，顺序同上。

如果上述文件都没有找到则会使用默认配置，其等价的配置为：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="WARN">
  <Appenders>
    <Console name="Console" target="SYSTEM_OUT">
      <PatternLayout pattern="%d{HH:mm:ss.SSS} [%t] %-5level %logger{36} - %msg%n"/>
    </Console>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="Console"/>
    </Root>
  </Loggers>
</Configuration>
```

Log4j2 可以通过JSON,XML，YAML格式的配置文件进行配置，其配置文件的结构为：

- Configuration
  - properties
  - Appenders
    - Console
      - PatternLayout
    - File
    - RollingRandomAccessFile
    - Async
  - Loggers
    - Logger
    - Root
      - AppenderRef

#### Configuration

Configuration作为根节点，有status和monitorInterval等多个属性

- status的值有trace，debug，info，warn，error，和fatal用于控制Log4j2本身的日志输出信息。
- monitorInterval表示每隔多少秒重新读取配置文件，这样可以不用重新启动而更新日志配置。
- nama：配置的名称

#### properties

用于定义全局变量，减少重复编码。

```xml
  <Properties>
    <Property name="filename">target/rolling1/rollingtest-$${sd:type}.log</Property>
  </Properties>
<!--使用 ${filename} -->
```

#### Appenders

Appenders负责将日志信息输出到指定的地方，可以是控制台，也可以是文件等等。

##### AsyncAppenders

AsyncAppenders引用其它的Appenders，然后使用多线程输出日志信息。

其属性有：

| 属性        | 类型    | 描述                                               |
| ----------- | ------- | -------------------------------------------------- |
| AppenderRef | String  | 多线程调用时使用的Appdenders和名称，可以配置多个。 |
| name        | String  | 设置这个Appenders的名称。                          |
| bufferSize  | Integer | 指定存放日志信息的缓冲区在最大容量，默认为1024。   |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <File name="MyFile" fileName="logs/app.log">
      <PatternLayout>
        <Pattern>%d %p %c{1.} [%t] %m%n</Pattern>
      </PatternLayout>
    </File>
    <Async name="Async">
      <AppenderRef ref="MyFile"/>
    </Async>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="Async"/>
    </Root>
  </Loggers>
</Configuration>
```

##### ConsoleAppender

将日志信息输出到标准输出流（控制台）。

其属性有：

| 属性   | 类型   | 描述                                                      |
| ------ | ------ | --------------------------------------------------------- |
| filter | Filter | 指定过滤器用于判断该日志信息是否需要被这个Appenders处理。 |
| layout | Layout | 指定日志信息的输出格式。                                  |
| target | String | 可选值为`SYSTEM_OUT`和`SYSTEM_ERR`，默认为第一个。        |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <Console name="STDOUT" target="SYSTEM_OUT">
      <PatternLayout pattern="%m%n"/>
    </Console>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="STDOUT"/>
    </Root>
  </Loggers>
</Configuration>
```

##### FailOverAppender

引用多个Appenders，当定义的主Appender失败时，则依次使用failovers中定义的次Appender。

其属性有：

| 属性      | 类型     | 描述                                                    |
| --------- | -------- | ------------------------------------------------------- |
| primary   | String   | 使用的主Appender的名称。                                |
| failovers | String[] | 使用的次Appenders的名称。                               |
| target    | String   | 可选值为`SYSTEM_OUT`或`STYSTEM_ERR`，默认为`SYSTEM_ERR` |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <RollingFile name="RollingFile" fileName="logs/app.log" filePattern="logs/app-%d{MM-dd-yyyy}.log.gz"
                 ignoreExceptions="false">
      <PatternLayout>
        <Pattern>%d %p %c{1.} [%t] %m%n</Pattern>
      </PatternLayout>
      <TimeBasedTriggeringPolicy />
    </RollingFile>
    <Console name="STDOUT" target="SYSTEM_OUT" ignoreExceptions="false">
      <PatternLayout pattern="%m%n"/>
    </Console>
    <Failover name="Failover" primary="RollingFile">
      <Failovers>
        <AppenderRef ref="Console"/>
      </Failovers>
    </Failover>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="Failover"/>
    </Root>
  </Loggers>
</Configuration>
```

##### FileAppender

FileAppender将日志信息输出到文件中。

其属性有：

| 属性           | 类型    | 描述                                                         |
| -------------- | ------- | ------------------------------------------------------------ |
| append         | boolean | 当为true时，日志信息会被追加到文件中，为false时则会覆盖文件中已有的内容，默认为true。 |
| bufferedIO     | boolean | 默认为true，日志信息会被先写入缓冲区中，当缓冲区满了后再写入文件中， |
| bufferSize     | int     | 缓冲区的大小，默认为8192字节。                               |
| fileName       | String  | 保存日志的文件，如果不存在会被自动创建。                     |
| immediateFlush | boolean | 每一次写入日志信息后都刷新缓冲区，这能保证日志信息成功写入文件但可能会影响性能。默认为true。 |
| layout         | Layout  | 指定日志信息的输出格式。                                     |
| name           | String  | 指定这个Appender的名称。                                     |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <File name="MyFile" fileName="logs/app.log">
      <PatternLayout>
        <Pattern>%d %p %c{1.} [%t] %m%n</Pattern>
      </PatternLayout>
    </File>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="MyFile"/>
    </Root>
  </Loggers>
</Configuration>
```

##### JDBCAppender

通过标准JDBC样品将日志信息写入关系性数据库中。

其属性有：

| 属性             | 类型             | 描述                                                         |
| ---------------- | ---------------- | ------------------------------------------------------------ |
| name             | String           | 设置这个Appender的名称。                                     |
| bufferSize       | int              | 如果值大于0，则日志信息会被先存放在缓冲区中，然后缓冲区满了后刷新缓冲区。 |
| connectionSource | ConnectionSource | 从connectionSource中获取数据库连接。                         |
| tableName        | String           | 保存日志信息的表名。                                         |
| columnConfig     | ColumnConfig[]   | 设置日志信息如何与列名对应。                                 |

在配置`connectionSource`属性时，需要使用以下标签来设置，详见[log4j2](http://logging.apache.org/log4j/2.x/manual/appenders.html#JDBCDataSource)：

- `<DataSource>`：使用JNDI
- `<ConnectionFactory>`：指定一个提供JDBC Connection的方法。
- `<DriverManager>`：通过DriverManager获取，没有池化技术。
- `<PoolingDriver>`：通过Apache DBCP数据库连接池来提供。

在配置`columnConfig`属性时，使用`<Column>`标签去指明日志信息应该如何写入表中。

其属性有：

| 属性             | 类型    | 描述                                                         |
| ---------------- | ------- | ------------------------------------------------------------ |
| name             | String  | 对应的字段名。                                               |
| pattern          | String  | 通过PatternLayout获取被插入的具体值。如：`%level`则插入的值为对应日志级别。 |
| isEventTimestamp | boolean | 设为true时，表明这个字段插入的值为当前的时间戳。             |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="error">
  <Appenders>
    <JDBC name="databaseAppender" tableName="LOGGING.APPLICATION_LOG">
      <ConnectionFactory class="net.example.db.ConnectionFactory" 
                         method="getDatabaseConnection" />
      <Column name="EVENT_DATE" isEventTimestamp="true" />
      <Column name="LEVEL" pattern="%level" />
      <Column name="LOGGER" pattern="%logger" />
      <Column name="MESSAGE" pattern="%message" />
      <Column name="THROWABLE" pattern="%ex{full}" />
    </JDBC>
  </Appenders>
  <Loggers>
    <Root level="warn">
      <AppenderRef ref="databaseAppender"/>
    </Root>
  </Loggers>
</Configuration>
```

其中ConnectionFactory对应的类如下：

```java
public class ConnectionFactory {
    private static interface Singleton {
        final ConnectionFactory INSTANCE = new ConnectionFactory();
    }
 
    private final DataSource dataSource;
 
    private ConnectionFactory() {
        Properties properties = new Properties();
        properties.setProperty("user", "logging");
        properties.setProperty("password", "abc123"); // or get properties from some configuration file
 
        GenericObjectPool<PoolableConnection> pool = 
            new GenericObjectPool<PoolableConnection>();
        DriverManagerConnectionFactory connectionFactory = 
            new DriverManagerConnectionFactory(
                "jdbc:mysql://example.org:3306/exampleDb", properties);
        new PoolableConnectionFactory(connectionFactory, pool, null,
                                      "SELECT 1", 3, false, false, 
                                      Connection.TRANSACTION_READ_COMMITTED);
 		//这里使用DBCP作为数据源
        this.dataSource = new PoolingDataSource(pool);
    }
 
    public static Connection getDatabaseConnection() throws SQLException {
        return Singleton.INSTANCE.dataSource.getConnection();
    }
}
```

##### RandomAccessFileAppender

和FileAppender相似，但它会一直启用缓冲区（无法关闭）并且内部采用ByteBuffer和RandomAccessFile实现读写，而FileAppender内部使用BufferedOutputStream实现读写，在性能上，RandomAccessFileAppender能提升20-200%。

其属性和FileAppender大致相同，只是没有`bufferedIO`属性。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <RandomAccessFile name="MyFile" fileName="logs/app.log">
      <PatternLayout>
        <Pattern>%d %p %c{1.} [%t] %m%n</Pattern>
      </PatternLayout>
    </RandomAccessFile>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="MyFile"/>
    </Root>
  </Loggers>
</Configuration>
```

##### RollingFileAppender

RollingFileAppender使用OutputStream实现将日志输出到文件中，并且可能通过`TriggeringPolicy`和`RolloverPolicy`进行滚动写入。例如：设置按照时间进行滚动，则每隔一段时间将产生一个新的日志文件，而不是继续在原来的文件中写入信息，详见[RollingFileAppender](http://logging.apache.org/log4j/2.x/manual/appenders.html#RollingFileAppender)。

其属性有（FileAppender包含的属性，这里也包含）：

| 属性        | 类型             | 描述                                                         |
| ----------- | ---------------- | ------------------------------------------------------------ |
| filePattern | String           | 保存文件名的表达式，当满足滚动条件时，新创建的文件使用此规则创建。 |
| policy      | TriggeringPolicy | 触发滚动的条件。                                             |
| strategy    | RolloverStrategy | 确定新文件名和位置的滚动策略。                               |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <RollingFile name="RollingFile" fileName="logs/app.log"
                 filePattern="logs/$${date:yyyy-MM}/app-%d{MM-dd-yyyy}-%i.log.gz">
      <PatternLayout>
        <Pattern>%d %p %c{1.} [%t] %m%n</Pattern>
      </PatternLayout>
      <Policies>
          <!--按照时间触发 这里按照filePattern中定义的将每一天产生一个日志文件-->
        <TimeBasedTriggeringPolicy />
          <!--按照大小触发，当文件超过250MB时产生一个新的日志文件-->
        <SizeBasedTriggeringPolicy size="250 MB"/>
      </Policies>
    </RollingFile>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="RollingFile"/>
    </Root>
  </Loggers>
</Configuration>
```

##### RollingRandomAccessFileAppender

由RandomAccessFile实现的RollingAppender，详见[RollingRandomAccessFileAppender](https://logging.apache.org/log4j/2.x/manual/appenders.html#RollingRandomAccessFileAppender)。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <RollingRandomAccessFile name="RollingRandomAccessFile" fileName="logs/app.log"
                 filePattern="logs/$${date:yyyy-MM}/app-%d{MM-dd-yyyy}-%i.log.gz">
      <PatternLayout>
        <Pattern>%d %p %c{1.} [%t] %m%n</Pattern>
      </PatternLayout>
      <Policies>
        <TimeBasedTriggeringPolicy />
        <SizeBasedTriggeringPolicy size="250 MB"/>
      </Policies>
    </RollingRandomAccessFile>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="RollingRandomAccessFile"/>
    </Root>
  </Loggers>
</Configuration>
```

#### Layouts

Layout用于定义日志信息的格式。

##### HtmlLayout

输出Html页面保存日志信息。

其属性有：

| 属性         | 类型    | 描述                                                        |
| ------------ | ------- | ----------------------------------------------------------- |
| charset      | String  | 设置编码，默认为UTF-8                                       |
| contentType  | String  | 设置content-type头，默认为“text/html”                       |
| locationInfo | boolean | 如果为true，将会在Html文件中输出文件名和行号，默认为false。 |
| title        | String  | Html文件的标题                                              |
| datePattern  | String  | 日期格式                                                    |
| timezone     | String  | 设置时区。                                                  |

```xml
<Appenders>
  <Console name="console">
      <!--GMT+8 表示东八区-->
    <HtmlLayout datePattern="ISO8601" timezone="GMT+8"/>
  </Console>
</Appenders>
```

##### PatternLayout

通过设置模板样式来输出日志信息。

其属性有：

| 属性    | 类型   | 描述                                                         |
| ------- | ------ | ------------------------------------------------------------ |
| charset | String | 指定编码                                                     |
| pattern | String | 定义模板，详见[Pattern](https://logging.apache.org/log4j/2.x/manual/layouts.html#Patterns) |
| header  | String | 在每个日志文件第一行显示的信息                               |
| footer  | String | 在最后一行显示的信息                                         |

```xml
<Appenders>
        <Console name="Console" target="SYSTEM_OUT">
            <PatternLayout header="header" 
                           pattern="%d{HH:mm:ss.SSS} [%t] %-5level %logger{36} - %msg %n" 
                           footer="footer"/>
        </Console>
    </Appenders>
```

#### Filters

Filter允许按照一定规则对日志信息进行过滤，选择丢弃或保留某些日志。

##### BurstFilter

BurstFilter可以通过在达到最大限制后丢弃日志来控制日志的处理速率。

其属性有：

| 属性       | 类型    | 描述                                                         |
| ---------- | ------- | ------------------------------------------------------------ |
| level      | String  | 被过滤的日志的最高级别，当达到最大限制时，在这个级别及之下的日志将会被过滤掉而不被记录。 |
| rate       | float   | 每秒允许处理日志的平均数。                                   |
| maxBurst   | integer | 最大日志处理速率，当超过这个值时，会触发过滤。默认为10倍rate。 |
| onMatch    | String  | 对于被过滤器匹配的日志所采用的处理方式，可选值有`ACCEPT`（保留日志），`DENY`（丢弃日志）和`NEUTRAL`（由其它过滤器处理）。默认为`NEUTRAL`。 |
| onMismatch | String  | 同上，只是作用于未被过滤器匹配的日志。默认为`DENTY`。        |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <RollingFile name="RollingFile" fileName="logs/app.log"
                 filePattern="logs/app-%d{MM-dd-yyyy}.log.gz">
      <BurstFilter level="INFO" rate="16" maxBurst="100"/>
      <PatternLayout>
        <pattern>%d %p %c{1.} [%t] %m%n</pattern>
      </PatternLayout>
      <TimeBasedTriggeringPolicy />
    </RollingFile>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="RollingFile"/>
    </Root>
  </Loggers>
</Configuration>
```

##### RegexFilter

使用正则表达式对日志信息进行匹配。

其属性有：

| 属性       | 类型    | 描述                                                         |
| ---------- | ------- | ------------------------------------------------------------ |
| regex      | String  | 正则表达式                                                   |
| useRawMsg  | boolean | 如果为true，则对未被格式化的日志进行匹配，则false则匹配格式化后的日志，默认为false。 |
| onMatch    | String  | 默认为NEUTRAL                                                |
| onMismatch | String  | 默认为DENY                                                   |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <RollingFile name="RollingFile" fileName="logs/app.log"
                 filePattern="logs/app-%d{MM-dd-yyyy}.log.gz">
      <RegexFilter regex=".* test .*" onMatch="ACCEPT" onMismatch="DENY"/>
      <PatternLayout>
        <pattern>%d %p %c{1.} [%t] %m%n</pattern>
      </PatternLayout>
      <TimeBasedTriggeringPolicy />
    </RollingFile>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="RollingFile"/>
    </Root>
  </Loggers>
</Configuration>
```

##### ThresholdFilter

通过指定日志级别来过滤。

其属性有：

| 属性       | 类型   | 描述                                           |
| ---------- | ------ | ---------------------------------------------- |
| level      | String | 匹配的日志级别，只有这个级别的日志才会被匹配。 |
| onMatch    | String | 默认为NEUTRAL                                  |
| onMismatch | String | 默认为DENY                                     |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <RollingFile name="RollingFile" fileName="logs/app.log"
                 filePattern="logs/app-%d{MM-dd-yyyy}.log.gz">
      <ThresholdFilter level="TRACE" onMatch="ACCEPT" onMismatch="DENY"/>
      <PatternLayout>
        <pattern>%d %p %c{1.} [%t] %m%n</pattern>
      </PatternLayout>
      <TimeBasedTriggeringPolicy />
    </RollingFile>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="RollingFile"/>
    </Root>
  </Loggers>
</Configuration>
```

##### TimeFilter

通过设置一天中的一个时间段来过滤日志。

其属性有（以下都省略onMatch和onMismatch）：

| 属性     | 类型   | 描述                     |
| -------- | ------ | ------------------------ |
| start    | String | 开始时间，格式为HH:mm:ss |
| end      | String | 结束时间，格式同上。     |
| timezone | String | 指定时区。               |

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Configuration status="warn" name="MyApp" packages="">
  <Appenders>
    <RollingFile name="RollingFile" fileName="logs/app.log"
                 filePattern="logs/app-%d{MM-dd-yyyy}.log.gz">
        <!--只记录每天5:00到5:30之间产生的日志-->
      <TimeFilter start="05:00:00" end="05:30:00" onMatch="ACCEPT" onMismatch="DENY"/>
      <PatternLayout>
        <pattern>%d %p %c{1.} [%t] %m%n</pattern>
      </PatternLayout>
      <TimeBasedTriggeringPolicy />
    </RollingFile>
  </Appenders>
  <Loggers>
    <Root level="error">
      <AppenderRef ref="RollingFile"/>
    </Root>
  </Loggers>
</Configuration>
```

#### Logger

日志器分为根日志器（必须有）和自定义日志器，类中的`getLogger`则是通过`name`属性来获取日志器，未获取到指定名称的日志器时，则使用根日志器。日志器可以指定对应的日志级别并选择appender来处理日志。

其属性有：

| 属性        | 类型        | 描述                                                         |
| ----------- | ----------- | ------------------------------------------------------------ |
| name        | String      | 日志器的名称。                                               |
| level       | String      | 日志器的日志级别，只有在这个级别及之上的日志才会被appenders处理 |
| additivity  | boolean     | 当为true时会同时将日志输出到根日志器中。                     |
| AppenderRef | AppenderRef | 定义用来处理日志的Appenders，可以有多个。                    |

```xml
<Loggers>
    <Root level="DEBUG">
        <AppenderRef ref="Console"/>
    </Root>
    <Logger name="demo.log.Test" additivity="false" level="INFO">
        <AppenderRef ref="File"/>
    </Logger>
</Loggers>
```

#### 补充1：Idea中设置提示

1. 在[github](https://github.com/apache/logging-log4j2/blob/master/log4j-core/src/main/resources/Log4j-config.xsd)上下载配置文件的XSD文件，并放在类路径下。

   在Log4j2的jar包中也包含对应的XSD文件，但似乎并没有更新，所以需要自己下载最新文件。

2. `Configuration`中添加如下信息

   其中`Log4j-config.xsd`为第一步下载的文件的文件名。

   ```xml
   <Configuration xmlns="http://logging.apache.org/log4j/2.0/config"
                  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                  xsi:schemaLocation="http://logging.apache.org/log4j/2.0/config 
                                      Log4j-config.xsd">
   ```

#### 补充2: 使用SLF4j作为日志门面

SLF4j作为日志门面，Log4j2作为日志实现。

maven依赖：

```xml
<dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-slf4j-impl</artifactId>
    <version>2.14.1</version>
</dependency>
```

## 日志门面（SLF4j）

### 日志门面

为什么要使用日志门面？

在java中有各种各样实现日志功能的库，包括Log4j2，logback，jul等等。如果直接使用这些日志库，则在会有大量包含相关日志框架的代码，如果某一天需要修改使用的日志库，则会巨大的工作量。同时，如果项目中使用了一些第三库，而这些第三库中又使用另一个日志框架，那么为了能够对日志进行配置，则不得不也使用其对应的日志框架，可如果有多个库并且每个都使用了不同的日志框架以如何解决呢？

> **计算机科学领域的任何问题都可以通过增加一个间接的中间层来解决。**

而日志门面则是采用这样的方式，在日志框架与应用之间建立一座桥梁，对于应用来说，无论底层使用的日志框架是什么，都只需要和日志门面打交道即可，因为日志门面提供了一系列统一的接口，而具体实现都交由日志框架来完成。

常见的日志门面有[SLF4j](http://www.slf4j.org)和[commons-logging](http://commons.apache.org/proper/commons-logging/)

### SLF4j

#### Maven依赖

因为SLF4j只是一个日志门面，不提供日志实现，所以实际使用时还需要添加日志实现的依赖，例如`slf4j-simple`或`logback-classic`（这两个库都直接提供了对SLF4j的实现，如果是其它日志框架可能还需要添加相关的实现库，例如上面的`log4j-slf4j-impl`）。

- slf4j-simple作日志实现

  ```xml
  <dependency>
      <groupId>org.slf4j</groupId>
      <artifactId>slf4j-simple</artifactId>
      <version>1.7.32</version>
   </dependency>
  ```

- logback作日志实现

  ```xml
  <dependency> 
    <groupId>ch.qos.logback</groupId>
    <artifactId>logback-classic</artifactId>
    <version>1.2.3</version>
  </dependency>
  ```

- JUL（JDK 1.4开始提供的日志框架`java.util.logging`）作日志实现

  ```xml
  <dependency> 
    <groupId>org.slf4j</groupId>
    <artifactId>slf4j-jdk14</artifactId>
    <version>1.7.31</version>
  </dependency>
  ```

#### 使用

```java
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class HelloWorld {
  public static void main(String[] args) {
    Logger logger = LoggerFactory.getLogger(HelloWorld.class);
    logger.info("Hello World");
  }
}
```

#### 参数占位符

```java
public class Wombat {
  final Logger logger = LoggerFactory.getLogger(Wombat.class);
  Integer t;
  Integer oldT;
  public void setTemperature(Integer temperature) {
   
    oldT = t;        
    t = temperature;
    logger.debug("Temperature set to {}. Old temperature was {}.", t, oldT);
    if(temperature.intValue() > 50) {
      logger.info("Temperature has risen above 50 degrees.");
    }
  }
} 
```

#### 流式API

```java
int newT = 15;
int oldT = 16;

// using traditional API
logger.debug("Temperature set to {}. Old temperature was {}.", newT, oldT);

// using fluent API, add arguments one by one and then log message
logger.atDebug().addArgument(newT)
    .addArgument(oldT)
    .log("Temperature set to {}. Old temperature was {}.");

// using fluent API, log message with arguments
logger.atDebug().log("Temperature set to {}. Old temperature was {}.", newT, oldT);

// using fluent API, add one argument and then log message providing one more argument
logger.atDebug().addArgument(newT)
    .log("Temperature set to {}. Old temperature was {}.", oldT);

// using fluent API, add one argument with a Supplier and then log message with one more argument.
// Assume the method t16() returns 16.
logger.atDebug().addArgument(() -> t16())
    .log(msg, "Temperature set to {}. Old temperature was {}.", oldT);
```

