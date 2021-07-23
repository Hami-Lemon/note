# <center>数据库-MySQL</center>

学习MySQL的一些笔记
<!---more--->
## 数据库管理
### 查看所有数据库  

    SHOW DATABASES;
### 创建数据库

    CREATE DATABASE 库名;

### 使用指定数据库

    USE 库名;

### 查看当前库的所有表

    SHOW TABLES;

### 创建表

    CREATE TABLE 表名(
        名称 类型 [约束],
        ...
        [表约束]
    );
### 查看表结构

    DESC 表名;

## 约束
- 主键:`PRYMARY KEY` 或 `PRIMARY KEY(id)`
例:

```
CREATE TABLE name(
    id INT(11)  PRIMARY KEY,

);
或
CREATE TABLE name(
    id INT(11),
    PRIMARY KEY(id)
);
```
- 外键：`CONSTRAINT 外键名 FOREIGN KEY(字段) REFERENCES 主表(字段)`
例：
```
CREATE TABLE  name(
id INT(11),
CONSTRAINT fk_emp_dept1 FOREIGN KEY(deptId) REFERENCES tb_dept1(id)
);
```
- 非空约束
`字段名 类型 NOT NULL`
- 唯一性约束
`字段名 类型 UNIQUE`
- 默认约束
`字段名 类型 DEFAULT 默认值`
- 设置标点属性值自动增加（只能用于主键，且一个表最多只能有一个）
`字段 类型 AUTO_INCREMENT`
 ## 修改数据表
 `ALTER TABLE`

- 修改表名
`ALTER TABLE 旧表名 RENAME 新表名;`
- 修改字段的数据类型
`ALTER TABLE 表名 MODIFY 字段名 类型;
- 修改字段名
`ALTER TABLE 表名 CHANGE 旧字段名 新字段名 类型;`(类型不可缺少，可以与原来不一样)
- 添加字段
```
ALTER TABLE 表名 ADD  字段 类型
    [约束]
    [FIRST || AFTER 已存在的字段名]
```
(FIRST指将字段设为第一个，after指放在指定字段的后面，都是可选参数)
- 删除字段
` ALTER TABLE 表名 DROP 字段;`
- 修改字段的排序位置
`ALTER TABLE 表名 MODIFY 字段 类型 FIRST | AFTER 字段名;`
- 更改表的存储引擎
`ALTER TABLE 表名 ENGINE=引擎名;`
- 删除外键约束
`ALTER TABLE 表名 DROP FOREIGN KEY 外键名;`
- 删除没有被关联的表
`DROP TABLE IF EXISTS 表名1，表名2,...;`(使用if exists 判是否存在)
- 删除被关联的表
    1. 解除外键约束
    2. 删除
## 8.0新特性
- 默认字符集为utf8mb4
- 自增变量的值会持久化
## 数据类型和运算符
### 数据类型
- 整数

类型名称|说明|储存需求
:-:|:-:|:-:
TINYINT(INTEGER)|很小的整数|1字节
SMALLINT(INTEGER)|小的整数|2字节
MEDIUMINT(INTEGER)|中等大小的整数|3字节
INT(INTEGER)|普通大小的整数|4字节
BIGING(INTEGER)|大整数|8字节
注：INT(INTEGER)中的INTEGER表示该数据类型指定的显示宽度，表示最大能够显示的数字个数，与范围无关，数值位小于指定的宽度时会由空格填充，大于时仍完整显示出来，不指定时会使用默认值
- 小数

类型名称|说明|储存需求
:-:|:-:|:-:
FLOAT(M,D)|单精度|4字节
DOUBLE(M,D)|双精度|8字节
DECIMAL(M,D), DEC|压缩的“严格”定点数|M+2字节  
注：M表示数据总长度（不含小数点），D表示小数点后保留几位，高精时使用DECIMAL

- 日期与时间类型

类型名称|日期格式|日期范围|存储需求
:-:|:-:|:-:|:-:
YEAR|YYYY|1901-2155|1字节
TIME|HH:MM:SS|-838:59:59-838:59:59|3字节
DATE|YYYY-MM-DD|1000-01-01~9999-12-3|3字节
DATETIME|YYYY-MMM-DD HH:MM:SS|1000-01-01 00:00:00~9999-12-31 23:59:59|8字节
TIMESTAMP|YYYY-MM-DD HH:MM:SS|1970-01-01 00:00:01 UTC ~ 2038-01-19 03:14:07 UTC|4字节
注：UTC 世界标准时间，TIMESTAMP的储存是以世界标准时间格式保存，会根据当前时区，返回不同结果

- 字符串类型

类型名称|说明|储存需求
:-:|:-:|:-:
CHAR(M)|固定长度非二进制字符串|M字节，1 <=M <=255
VARCHAR(M)|变长非二进制字符串|L+1字节，L<=M和1 <=M <=255
TINYTEXT|非常小的非二进制字符串|L+1 字节
TEXT|小的非二进制字符串|
MEDIUMTEXT|中等大小的非二进制字符串|
LONGTEXT|大的非二进制字符串|
ENUM|枚举类型，只能有一个枚举字符串值|
SET|一个设置，字符串对象可以有零个或多个SET成员|
注：char的长度固定，为M,varchar长度可变，最大为M,text类型主要用于保存文字内容，评论等，enum以编号存储

- 二进制字符串类型

类型名称|说明|储存需求
:-:|:-:|:-:
BIT(M)|位字段|
BINARY(M)|固定长度的二进制字符串|
VARBINARY(M)|可变长度的二进制字符串|
TINYBLOB(M)|非常小的二进制字符串|
BLOB(M)|小BLOB|
MEDIUNBLOB(M)|中等大小的BLOB|
LONGBLOB(M)|非常大的BLOB|

注：M表示每个值的位数（bit位）

## 运算符
- 算术运算符

    `+  -   *   /   %`
- 比较运算符
  `>  <   =   >=  <=  !=(<>)  <=>`
  以及`IN、NOT IN、BETWEEN AND、IS NULL、IS NOT NULL、GREATEST、LEAST、LIKE、REGEXP`
    - `<>`：可以用来判断null
    - `BETWEEN A AND B`:等价与大于等于A且小于等于B
    - `LEAST(值1，值2，...)`:返回最小值，若有一个为null,则值为null
    - `GREATEST（...）`：同上，取最大值
    - `LIKE '...'`:匹配字符串，
        - `%` 匹配任何数目的字符
        - `_` 只能匹配一个字符
    -  `REGEXP` 正则匹配
        - `^` 匹配以该字符后面的字符开头的字符串 
        - `$` 匹配以该字符后面的字符结尾的字符串
        - `.` 匹配任何一个单字符
        - `[...]` 匹配方括号呢的任何字符
            例如:`[abc]`匹配'a' 'b' 'c'，可使用‘-’ 选定范围
        - `*` 匹配零个或多个在它前面的字符
- 逻辑运算符
`NOT AND OR XOR（异或）` 
- 位运算符
`&  |   ~   ^   <<  >>`
## 插入、更新、删除数据
- 插入
    - 为表的所有字段插入数据
    `INSERT INTO tablename (字段列表) VALUES (值列表);`
    可同时插入多条记录
    ```
    INSERT INTO tablename (字段列表) 
    VALUES (值列表),
    VALUES (值列表),
    VALUES (值列表),
    VALUES (值列表),
    ...
    VALUES (值列表);
    ```

    - 将查询结果插入表中
    ```
    INSERT INTO  tablename (字段列表)
    SELECT ... FROM ...;
    ```
- 更新
`UPDATE 表名 SET 字段=值, ... WHERE 条件;`(不加条件时则全部更新！！！)
- 删除
`DELETE FROM 表名 WEHERE 条件;`(不加条件时则全部删除！！！)
- 为表增加计算列
指定某一列的值是由其它列计算而来，不需要手动插入
例：
```
CREATE TABLE name(
    id INT(9),
    a INT(9),
    b INT(9),
    c INT(9) GENERATED ALWAYS AS ((a+b)) VIRTUAL

);
# c的值将由a+b决定，不需要手动插入
```
