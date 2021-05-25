# <center>整数的分解（递归）</center>

### 题目描述：

给定一个整数n和m，将整数n划分成几个数相加的形式，其中最大的数不超过m（m可能大于n）,求能分解成几种结果
<!--more-->
### 例：给定`n = 3, m = 3`
输出：

    3 = 3
    3 = 2 + 1
    3 = 1 + 1 + 1

### 分析：
![ ](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113652.png)

### 代码实现：

```C
int q(n, m){
    if(m == 1)
        return 1;
    if(n < m)
        return q(n, n);
    if(n == m)
        return q(n, n - 1) + 1;
    if(n > m)
        return q(n, m - 1) + q(n - m, m);
}
```