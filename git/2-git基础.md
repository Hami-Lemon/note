# Git 系列学习 -- Git基础

## 获取git仓库

1. 将本地项目转换为git仓库

   ```shell
   git init
   ```

2. 克隆现有仓库

   ```shell
   git clone <url>
   ```

### 记录更新到仓库

![文件状态变化](images/2-git%E5%9F%BA%E7%A1%80/image-20210511091520522.png)

### 检查当前文件状态

通过以下命令查看哪些文件处于什么状态

```shell
git status
```

- `changes to be committed`:已跟踪的文件，但还未提交到仓库
- `changes not staged for commit`:已跟踪文件中发生了修改的文件，但还未暂存

#### 状态简览

使用`git status -s`

![image-20210511095357695](images/2-git%E5%9F%BA%E7%A1%80/image-20210511095357695.png)

新添加的文件用`??`标记，新添加到暂存区中的文件前有`A`标记，修改过的文件用`M`标记。输出有两栏，左栏指明了暂存区的状态，右栏指明了工作区的状态。

例如，上面中显示，``README`已修改，但暂未暂存，`Rakefile`已经修改，且暂存之后又作了修改

### 跟踪新文件

创建一个新文件后，使用以下命令跟踪该文件

`git add`使用文件或目录的路径作为参数，如果参数是目录，则会递归地跟踪该目录下的所有文件

```shell
git add <file>
```

### 暂存已修改的文件

当已跟踪的文件内容发生变化时，使用`git add`可以把该文件放到暂存区，这样，在下次提交时就会一并记录到仓库中。

当暂存之后继续修改文件，如果此时提交，提交的版本是暂存区中的版本，而不是第二次修改后的版本，需要重新运行`git add`把最新版本暂存

### 忽略文件

可使用此仓库中提供的通用模版[gitignore](https://github.com/github/gitignore)

当某些文件不需要纳入git的管理时，可以创建一个名为`.gitignore`的文件，在里面编写忽略规则

格式规范如下：

- 所有空行或`#`开头会被忽略
- 可以使用标准的glob模式（shell使用的简化正则）匹配，它会递归地应用在整个工作区中
- 可以以`/`开关防止递归
- 以`/`结尾指定目录
- 在模式前加`!`取反，可忽略指定模式以外的文件或目录

glob模式：

- `*`匹配零个或多个任意字符
- `[abc]`匹配任何一个在括号中的字符
- `?`匹配一个任意字符
- `[0-9]`表示匹配0到9的数字
- `**`表示匹配任意中间目录，`a/**/z`

例:

```txt
# 忽略所有的 .a 文件
*.a
# 但跟踪所有的 lib.a，即便你在前面忽略了 .a 文件
!lib.a
# 只忽略当前目录下的 TODO 文件，而不忽略 subdir/TODO
/TODO
# 忽略任何目录下名为 build 的文件夹
build/
# 忽略 doc/notes.txt，但不忽略 doc/server/arch.txt
doc/*.txt
# 忽略 doc/ 目录及其所有子目录下的 .pdf 文件
doc/**/*.pdf
```

### 查看修改

`git diff`比较工作目录中当前文件和暂存区域之间的差异

`git diff --staged`查看已暂存文件与最后一次提交文件的差异

### 提交更新

使用`git commit`将暂存区中的文件提交到仓库

`git commit -m "提交信息"`

#### 跳过使用暂存区

在提交时使用`git commit -a`可跳过暂存区，直接把所有已经跟踪的文件暂存起来一并提交

### 移除文件

使用`git rm <file>`删除文件，会把暂存区和工作区中对应的文件一同删除，如果删除之前修改过文件，则需要加上`-f`参数

如果手动从工作区删除文件，则这此删除不会添加到暂存区，需要`git rm <file>`记录此次删除操作

使用`git rm --cached <file>`只删除暂存区的文件而不删除工作区的文件

### 移动文件

`git mv <old file> <new file>`修改文件名也用此方法

### 查看提交历史

`git log`查看所有的提交历史信息

`git log -p -2`显示每次提交的差异信息，并限制条数为`2`

`git log --stat`:列出每次提交的简略统计信息

`git log --pretty=oneline|short|full|fuller`:以不同的格式显示

`git log --pretty=format:"%h"`:自定义格式

![format常用选项](images/2-git%E5%9F%BA%E7%A1%80/image-20210511124124817.png)

`git log --pretty=format:"%h" --graph`:以图形方式显示

### 撤消操作

`git commit --amend`重新提交暂存区内容，用于提交后漏掉几个文件，或修改提交信息

#### 取消暂存

`git reset HEAD <file>`或`git restore --staged <file>（推荐）`撤消对应文件的暂存

#### 撤消修改

`git checkout -- <file>`或`git restore <file>`撤消未暂存文件的修改

## 远程仓库使用

### 查看远程仓库

`git remote`

添加远程仓库

`git remote ad <名称> <url>`