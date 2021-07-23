# CMAKE

CMake是一个用于管理项目源代码构建的工具，可以用简单的语句来描述所有平台的安装

## 基本项目

创建一个`CMakeLists.txt`文件配置CMAKE项目，一个最基础的项目需要包含以下信息

1. `cmake_minimum_required(VERSION 3.10)`指定cmake的最低版本
2. `projet(Tutorial)`设置项目的名称，会生成一个`${PROJECT_NAME}`变量，可以在其它地方使用
3. `add_exectuable(Tutorial tutorial.cpp)`设置生成的可执行文件，第一个参数为生成文件的名称，后面的参数为源代码文件

## 设置版本号并配置头文件

`project(Tutorial VERSION 1.0)`指定版本号

配置头文件并获取到版本号

```cmake
# 会在${PROJECT_BINARY_DIR}(编译文件的目录)中生成TutorialConfig.h
configure_file(TutorialConfig.h.in TutorialConfig.h)
# 将生成的头文件引入到项目中,需要写在最后
target_include_directories(Tutorial PUBLIC "${PEOJECT_BINARY_DIR}")
```

编写`TutorialConfig.h.in`文件

```cpp
#define Tutorial_VERSION_MAJOR @Tutorial_VERSION_MAJOR@
#define Tutorial_VERSION_MINOR @Tutorial_VERSION_MINOR@
```

在使用CMAKE配置项目后，会其中的`@Tutorial_VERSION_MAJOR`替换成相应的数字

## 指定C++标准

`set(CMAKE_CXX_STANDARD 11)`需要放在`add_executable`前面

## 添加Library

将一个源文件生成为一个库

```cmake
add_library(Name SHARED xxx.cpp)
# 还项目中引用上面生成的库文件
target_link_libraries(PROJECT PUBLIC Name)
```

