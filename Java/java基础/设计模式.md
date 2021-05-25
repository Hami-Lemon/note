# 设计模式

## UML

UML(统一建模语言)是一种由一整套图组成的标准化建模风格

### 泛化关系

一种继承关系，表示一般与特色的关系，指定了子类如何特化父类的属性和方法
![泛化](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113421.png)

### 实现关系

一种类与接口的关系,表示类是接口所有属性和方法的实现
![实现](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113423.png)

### 聚合关系

是整体和部分的关系,部分可以离开整体而单独存在

eg:一个公司有多个员工,但员工可以离开公司单独存在  
![聚合](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113425.png)

### 组合关系

表示整体和部分的关系,版部分离开整体后无法单独存在,是一种比聚合更强的关系

![组合](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113427.png)

### 关联关系

一种拥有关系,是一个类知道另一个类的属性和方法,关联可以是双向的,也可以是单向的

![关联](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113429.png)

### 依赖关系

一种使用的关系,即一个类的实现需要另一个类的协

![依赖](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113431.png)

## 监听模式(观察者模式)

是一种一对多的关系,可以有任意个观察者对象同时监听一个对象  
监听的对象叫观察者,被监听的的对象叫被观察者  
被观察者对象在状态或内容发生变化是,会通知所有观察者对象,使它们能够做出相应的变化  

```java
/**
 * 被观察者基类
 * @author Hami Lemon
*/
public abstract class BaseObservable {
    List<Observer> observers;

    public BaseObservable() {
        observers = new ArrayList<>();
    }
    public  void addObserver(Observer observer){
        observers.add(observer);
    }
    public void removeObserver(Observer observer){
        observers.remove(observer);
    }
    public void notifyObservers(){
        for (Observer observer : observers) {
            observer.update(this);
        }
    }
}

```

```java
/**
 * 观察者接口
 * @author Hami Lemon
 */
public interface Observer {
    void update(Observable obj);
}

```

## 状态模式

允许一个对象载器内部状态发生改变时改变其行为,使这个对象看上去就像改变了它的类型一样,状态即事物所处的某一形态,一个对象在其内部状态发生改变时,其表现的行为和外在属性不一样,状态模式又称为对象的行为模式

```java
/**
 * 上下环境类
 * @author Hami Lemon
 */
public abstract class BaseContext {
    private final Set<BaseState> states;
    private BaseState nowState;

    public BaseContext() {
        states = new HashSet<>();
        nowState = null;
    }

    public void addState(BaseState state){
        states.add(state);
    }

    public boolean changeState(BaseState state){
        if(state == null) {
            return false;
        }
        if(nowState == null){
            //初始化
        }else{
            //变为xxx
        }
        nowState = state;
        addState(state);
        return true;
    }

    public BaseState getNowState() {
        return nowState;
    }
}

```

```java
/**
 * 状态的基类
 * @author Hami Lemon
 */
public abstract class BaseState {
    @Override
    public int hashCode() {
        return super.hashCode();
    }

    @Override
    public boolean equals(Object obj) {
        return super.equals(obj);
    }

    /**
     * 不同状态的行为,不同的状态应该时单例的,或者设置属性标识身份
     * @return 不确定的返回值
     */
    abstract Object behavior();
}
```

## 中介模式

用一个中介对象来封装一系列的对象交互,中介者使各对象不需要显示地相互利用,从而使其耦合松散,而且可以独立地改变它们之间的交互
![中介模式](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113436.png)

## 装饰模式

动态地给一个对象增加一些额外的职责,扩展对象的功能

### 增强功能的装饰模式(透明装饰模式)

```java
//颜值接口
public interface IBeauty{
    int getBeautyValue();
}
```

```java
//新建Me,实现颜值接口
public class Me implements IBeauty{
    @Override
    public int gtBeautyValue(){
        return 100;
    }
}
```

```java
//戒指装饰类
public class RingDecorator implements IBeauty {
    private final IBeauty me;

    public RingDecorator(IBeauty me) {
        this.me = me;
    }

    @Override
    public int getBeautyValue() {
        return me.getBeautyValue() + 20;
    }
}
```

添加更多装饰类

```java
//耳环装饰类
public class EarringDecorator implements IBeauty {
    private final IBeauty me;

    public EarringDecorator(IBeauty me) {
        this.me = me;
    }

    @Override
    public int getBeautyValue() {
        return me.getBeautyValue() + 50;
    }
}
```

```java
//项链装饰类
public class NecklaceDecorator implements IBeauty {
    private final IBeauty me;

    public NecklaceDecorator(IBeauty me) {
        this.me = me;
    }

    @Override
    public int getBeautyValue() {
        return me.getBeautyValue() + 80;
    }
}
```

测试

```java
public class Client {
    @Test
    public void show() {
        IBeauty me = new Me();
        System.out.println("我原本的颜值：" + me.getBeautyValue());

        // 随意挑选装饰
        IBeauty meWithNecklace = new NecklaceDecorator(me);
        System.out.println("戴上了项链后，我的颜值：" + meWithNecklace.getBeautyValue());

        // 多次装饰
        IBeauty meWithManyDecorators = new NecklaceDecorator(new RingDecorator(new EarringDecorator(me)));
        System.out.println("戴上耳环、戒指、项链后，我的颜值：" + meWithManyDecorators.getBeautyValue());

        // 任意搭配装饰
        IBeauty meWithNecklaceAndRing = new NecklaceDecorator(new RingDecorator(me));
        System.out.println("戴上戒指、项链后，我的颜值：" + meWithNecklaceAndRing.getBeautyValue());
    }
}
```

### 添加功能的装饰模式(半透明装饰模式)

```java
//房屋接口
public interface IHouse{
    void live();
}
```

```java
//房屋类
public class House implements IHouse{

    @Override
    public void live() {
        System.out.println("房屋原有的功能：居住功能");
    }
}
```

```java
//粘钩装饰器接口
public interface IStickyHookHouse extends IHouse{
    void hangThings();
}
```

```java
//粘钩装饰类
public class StickyHookDecorator implements IStickyHookHouse {
    private final IHouse house;

    public StickyHookDecorator(IHouse house) {
        this.house = house;
    }

    @Override
    public void live() {
        house.live();
    }

    @Override
    public void hangThings() {
        System.out.println("有了粘钩后，新增了挂东西功能");
    }
}
```

测试

```java
public class Client {
    @Test
    public void show() {
        IHouse house = new House();
        house.live();

        IStickyHookHouse stickyHookHouse = new StickyHookDecorator(house);
        stickyHookHouse.live();
        stickyHookHouse.hangThings();
    }
}
```

在半透明装饰模式中,无法对同一个对象多次装饰

## 单例模式

确保一个类只有一个实例,并且提供一个获取它的全局方法

### 饿汉式

```java
public class Singleton {
  
    private static Singleton instance = new Singleton();

    private Singleton() {
    }

    public static Singleton getInstance() {
        return instance;
    }
}
```

### 懒汉式

#### 线程不安全

```java
public class Singleton {
  
    private static Singleton instance = null;
  
    private Singleton() {
    }

    public static Singleton getInstance(){
        if (instance == null) {
            instance = new Singleton();
        }
        return instance;
    }
}
```

#### 线程安全

- 单检锁

```java
public class Singleton {

    private static Singleton instance = null;

    private Singleton() {
    }

    public static Singleton getInstance() {
        synchronized (Singleton.class) {
            if (instance == null) {
                instance = new Singleton();
            }
        }
        return instance;
    }
}
```

- 双检锁  
如果 instance 已经被实例化，则不会执行同步化操作，大大提升了程序效率

```java
public class Singleton {

    private static Singleton instance = null;

    private Singleton() {
    }

    public static Singleton getInstance() {
        if (instance == null) {
            synchronized (Singleton.class) {
                if (instance == null) {
                    instance = new Singleton();
                }
            }
        }
        return instance;
    }
}
```

- 静态内部类方式

java的内部类采用懒加载方式,当内部类在使用使才会加载,同时,当访问一个类的静态字段使,如果该类未初始化,则立即初始化此类  
jvm运行时,会自动为初始化方法加锁,同步,确保一次只有一个线程执行初始化方法

```java
public class Singleton {

    private static class SingletonHolder {
        public static Singleton instance = new Singleton();
    }

    private Singleton() {
    }

    public static Singleton getInstance() {
        return SingletonHolder.instance;
    }
}
```

## 克隆模式(原型模式)

用原型实例指定要创建对象的种类,并通过拷贝这些原型的属性来创建新的对象

```java
public class MilkTea{
    public String type;
    public boolean ice;

    public MilkTea clone(){
        MilkTea milkTea = new MilkTea();
        milkTea.type = this.type;
        milkTea.ice = this.ice;
        return milkTea;
    }
}
```

可使用java中的语法糖,让需要拷贝的类实现Cloneable接口,从而不需要手写clone()方法  
此方法是浅拷贝,只有基本类型的参数会被拷贝一份,引用类型的对象不会被拷贝,而是继续使用传递引用的方式  
如果需要实现深拷贝,需要之间修改clone()方法

```java
public class MilkTea implements Cloneable{
    public String type;
    public boolean ice;

    @NonNull
    @Override
    protected MilkTea clone() throws CloneNotSupportedException {
        return (MilkTea) super.clone();
    }
}
```

## 职责模式(责任链模式)

为避免请求发送者与接收者耦合在一起,让多个对象都有可能接收请求.将这些接收的对象连接成一条链,并且沿着这条链传递请求,直到有对象处理为之

```java
/**
 * 责任链模式
 * 责任人的基类
 * @author Hami Lemon
 */
public abstract class Responsible {
    /**
     * 下一个责任人
     */
    private Responsible nextHandler;

    public Responsible getNextHandler() {
        return nextHandler;
    }

    public void setNextHandler(Responsible nextHandler) {
        this.nextHandler = nextHandler;
    }

    /**
     * 处理业务
     * @param obj 需要处理的请求
     * @return 视情况而定的返回值
     */
    public Object handle(Object obj){
        //判断能否处理请求
        //如果不能，则传递给下一个

        return new Object();
    }
}
```

## 代理模式

为其它对象提供一种代理以控制对这个对象的访问

## 外观模式

为子系统中的一组接口提供一个一致的界面,外观模式定义了一个高层接口,这个接口使得这一子系统更易使用
![外观模式](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113445.png)

## 迭代模式

提供一种方法顺序地访问一组聚合对象(一个容器)中的各个元素,而又不需要对外暴露该对象的细节  
![迭代模式](https://gitee.com/Hami-Lemon/image-repo/raw/master/images/2021/05/25/20210525113447.png)

## 组合模式

又叫部分整体模式，是用于把一组相似的对象当作一个单一的对象。组合模式依据树形结构来组合对象，用来表示部分以及整体层次。这种类型的设计模式属于结构型模式，它创建了对象组的树形结构  
组合模式用于整体与部分的结构,当整体与部分有相似的结构,在操作时可以被一致对待时,就可以使用组合模式