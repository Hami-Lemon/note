# 五笔（98版）

98版五笔入门学习
> 本文参考于：https://www.zhihu.com/question/19816777/answer/1054269608

![](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115520.jpg)

<!--more-->

![image-20210107115215489](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115534.png)

## 单字输入

将汉字分为键面汉字和键外汉字

### 键面汉字

在字根表里已经存在的汉字

键面汉字又可分为键名汉字和成字汉字

#### 键名汉字

字根表中的第一个字根

输入规则：**把所在键连打四下**，例如：“土 =FFFF”

#### 成字汉字

字根表里有的，但出去键名汉字外的汉字

输入规则：**报户口+第一单笔画+第二单笔画+最末单笔画**，报户口是指按下该字根所做的键，之后的笔画就是按区域划分的第一个键，即：G H T Y N，不足四笔就用空格键。例如：“士（shi）”可分为“一” “|” “一” 编码为：FGHG

#### 补码码元

98五笔中有三个比较特殊的码元（字根），犭，礻，衤，一个码元需要用两个键输入（**报户口+最后一笔**) 

输入“犭”时，先打犭所在的键Ｑ，再补最后一撇所在的键Ｔ。

输入“礻”时，先打礻所在的键Ｐ，再补最后一点所在的键Ｙ。

输入“衤”时，先打衤所在的键Ｐ，再补最后两点所在的键Ｕ。

### 键外汉字

字根表里不存在的汉字，可以理解为用字根凑成的字

- 只能拆分为两个字根，**第一个字根+第二字根+空格**,如“明”字，拆分为“日”和“月”两个字根打“J” “E”再打空格
- 只能拆分为三个字根，**第一字根+第二字根+第三字根+空格**，如“些”字，拆分为“止”和“匕”“二”三个字根打“H” “X” “F”再打空格
- 只能拆分为四个字根，**第一字根+第二字根+第三字根+第四字根**，如“都”，拆分为“土” “丿”“日”和“阝”四个字根打“F” “T” “J” “B”
- 超过四个字根，**第一字根+第二字根+第三字根+最末字根**，如“幅”， 拆分为“冂”“丨”“一”“口”和“田”五个字根，所以打“M” “H” “G” “L”

#### 识别码

“只”和“叭”的编码都是“KW”，为了更快地识别想打的是哪个，就加入了识别码。

识别码是由**末根末笔笔画**和**汉字的字形**组成的一个附加码，只在不足够四个码的时候才用，四码没有识别码的。字形组成是字的结构，比如左右结构、上下结构、半包围这些，**在五笔中只分为左右型、上下型和除此以外的称为杂合型**。

识别码的规则：汉字的**最后一个字根的未笔笔画**作为区号（即前面所说的按笔画的类型分的区），将汉字的字形作为位号（每个区里面的号码数），左右型、上下型和杂合型分别是1、2、3号位，两者结合就能确定识别码了。![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115540.jpg)

#### 一级简码

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115547.jpg)

#### 拆字规则

按结构的区别将汉字分为四类，“单结构”，“散结构”，“连结构”，“交结构”

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115552.jpg)



![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115557.jpg)

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115601.jpg)

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115605.jpg)

拆字是按照书写书写顺序、取大优先、能连不交、能散不连的原则

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115609.jpg)

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115613.jpg)

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115616.jpg)

![img](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115620.jpg)

## 词组输入

- 二字词，**由每字全码的前两码组成**

  例如，“后悔”，后=RGK，悔=NTX，那么各取前两个编码，即后悔=RGNT。

- 三字词输入规则：**前两字各取全码的第一码，最后一字取前两码，共四码。**

  例如，“合格证”，合=WGK，格=STK，证=YGH，那么合格各取第一个，证取前两个，即合格证=WSYG。

- 四字词输入规则：**由每字全码的第一码组成，共四码。**

  例如，“心想事成”，心=NY，想=SHN，事=GK，成=DN，那么各取第一个编码，即心想事成=NSGD。

- 多于四字的词组输入规则：**取第一、二、三和最后一个字的第一个码，共四码。**多字词包括**特定的词组和诗句**等。

  例如，“新疆维吾尔自治区”，新=USR，疆=XFG，维=XWY，区=AQ，那么分别取这几个字的第一个码，即新疆维吾尔自治区=UXXA。
  
   **五个单笔画的编码硬性规定为： “一”是GGLL, “|” 是HHLL “丿”TTLL, “丶”为YYLL, “乙”为 NNLL**

## 特殊字(个人认为)

- 州

![image-20210205165911342](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115625.png)

- 永

![image-20210205165935818](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115630.png)

- 假

![image-20210205170010102](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115634.png)

- 兰 

![image-20210205170029628](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115638.png)

- 承

![image-20210205170045054](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115643.png)

- 养

![image-20210205170102036](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115647.png)

- 藏

![image-20210205170136762](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115650.png)

- 臧

  ![image-20210205170402047](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115654.png)

-   曲

  ![image-20210205170449416](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115657.png)

-  班
  
  ![image-20210205170557518](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115701.png)

-  盖

  ![image-20210205170629495](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115704.png)

-  睡

  ![image-20210205170646762](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115708.png)

- 久

  ![image-20210205170701276](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115712.png)

- 段

  ![image-20210205170721307](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115716.png)

-  七

  ![image-20210205170742596](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115720.png)

- 忍

  ![image-20210205170805134](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115725.png)

- 尴

  ![image-20210205170915420](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115728.png)

- 尬

  ![image-20210205170948757](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115732.png)
  
- 丈

  ![image-20210205220238877](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115738.png)   
  
- 延

  ![image-20210205221128143](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115744.png)
  
- 垂

  ![image-20210206171930438](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115746.png)

- 燕

  ![image-20210206175103191](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525115750.png)

## 重码字

`fcu` ：去 云 支(`fc`)