# <center>全排列问题（递归解决）</center>

输出一个数列的全排列  
例：数列{1，2，3}  
它的全排列为：

    {1，2，3}、{2，1，3}、{1，3，2}
    {2，3，1}、{3，2，1}、{3，1，2}
<!--more-->

#### 代码实现：
```c
void perm(int *, int, int);
int main(){
    int list[5] = [1,2,3,4,5];
    perm(list, 0, 4);
    return 0;
}
void perm(int *list, int k, int m){
    if(k == m){
        int i;
        for(i = 0; i <= m; i++){
            printf("%d \n", list[i]);
        }
    }else{
        int t;
        for(t = k; t < m; t++){
            int temp = list[t];
            list[t] = list[k];
            list[k] = temp;

            perm(list, k+1, m);

            temp = list[t];
            list[t] = list[k];
            list[k] = temp;
        }
    }
}

```