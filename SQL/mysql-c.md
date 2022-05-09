c语言连接MySQL数据库，vistual studio环境下
<!--more-->

1. 将mysql server 中的 include 和 lib 目录复制到项目目录下

2. 创建空项目，在上方项目中选择属性c

![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115031.png)
3. 包含目录中添加 include目录，库目录中添加 lib目录
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115037.png)

4. 附加依赖项中添加libmysql.lib, 并将lib中的libmysql.lib和libmysql.dll，复制到项目目录下
![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115042.png)

![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115048.png)


```c
#include <winsock.h>
#include <mysql.h>
#include <stdio.h>
int main() {
	MYSQL* mysql = mysql_init(NULL);//初始化，获得mysql实例
	mysql_real_connect(mysql, "39.107.39.204", "root", "root", "test", 3306, NULL, 0);//链接数据库
	mysql_query(mysql, "show databases;");//执行sql语句
	MYSQL_RES *result;//结果集
	MYSQL_ROW row;//每一行的结果
	result = mysql_store_result(mysql);//获取结果集
	while (row = mysql_fetch_row(result)) {//获取每一行的数据，row[0],为第一列的数据，row[1]则为第二列，...

	}
	mysql_free_result(result);//释放结果集
	mysql_close(mysql);//关闭连接
	return 0;
}

```