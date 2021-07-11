# 变量句柄（VarHandle)

java9新增内容，在并发操作时，需要对类中的某个字段进行原子操作。

VarHandle提供了一系列标准的内存屏障操作，用于更加细粒度的控制内存排序。在安全性、可用性、性能上都要优于现有的API。VarHandle 可以与任何字段、数组元素或静态变量关联，支持在不同访问模型下对这些类型变量的访问，包括简单的 read/write 访问，volatile 类型的 read/write 访问，和 CAS(compare-and-swap)等。

## 创建VarHandle

```java
	private int plainStr;    //普通变量
    private static int staticStr;    //静态变量
    private int reflectStr;    //通过反射生成句柄的变量
    private int[] arrayStr = new int[]{100, 200, 300};    //数组变量

  	private static final VarHandle plainVar;    //普通变量句柄
    private static final VarHandle staticVar;    //静态变量句柄
    private static final VarHandle reflectVar;    //反射字段句柄
    private static final VarHandle arrayVar;    //数组句柄
    static
    {
        try
        {
            MethodHandles.Lookup lookup = MethodHandles.lookup();
            plainVar = lookup.findVarHandle(VarHandleTest.class, "plainStr", int.class);
            staticVar = lookup.findStaticVarHandle(VarHandleTest.class, "staticStr", int.class);
            reflectVar = lookup.unreflectVarHandle(VarHandleTest.class.getDeclaredField("reflectStr"));
            arrayVar = MethodHandles.arrayElementVarHandle(int[].class);
        } catch (ReflectiveOperationException e)
        {
            throw new ExceptionInInitializerError(e);
        }
    }
```

