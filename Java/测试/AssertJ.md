# AssertJ断言库

AssertJ提供了一种流式API的断言操作。

例如：

```java
// 引用所有的静态方法
import static org.assertj.core.api.Assertions.*;

// 基础断言
assertThat(frodo.getName()).isEqualTo("Frodo");
assertThat(frodo).isNotEqualTo(sauron);

// 对String断言操作
assertThat(frodo.getName()).startsWith("Fro")
                           .endsWith("do")
                           .isEqualToIgnoringCase("frodo");

// 对集合断言操作
// fellowshipOfTheRing 是 List<TolkienCharacter>
assertThat(fellowshipOfTheRing).hasSize(9)
                               .contains(frodo, sam)
                               .doesNotContain(sauron);

// as()方法设置测试断言描述，会在失败时显示。
assertThat(frodo.getAge()).as("check %s's age", frodo.getName()).isEqualTo(33);

// 断言异常，标准风格
assertThatThrownBy(() -> { throw new Exception("boom!"); }).hasMessage("boom!");
// BDD风格
Throwable thrown = catchThrowable(() -> { throw new Exception("boom!"); });
assertThat(thrown).hasMessageContaining("boom");

// using the 'extracting' feature to check fellowshipOfTheRing character's names
assertThat(fellowshipOfTheRing).extracting(TolkienCharacter::getName)
                               .doesNotContain("Sauron", "Elrond");

// extracting multiple values at once grouped in tuples
assertThat(fellowshipOfTheRing).extracting("name", "age", "race.name")
                               .contains(tuple("Boromir", 37, "Man"),
                                         tuple("Sam", 38, "Hobbit"),
                                         tuple("Legolas", 1000, "Elf"));

// 断言前对集合进行过滤
assertThat(fellowshipOfTheRing).filteredOn(character -> character.getName().contains("o"))
                               .containsOnly(aragorn, frodo, legolas, boromir);

// combining filtering and extraction (yes we can)
assertThat(fellowshipOfTheRing).filteredOn(character -> character.getName().contains("o"))
                               .containsOnly(aragorn, frodo, legolas, boromir)
                               .extracting(character -> character.getRace().getName())
                               .contains("Hobbit", "Elf", "Man");
```

## Maven依赖

```xml
<dependency>
  <groupId>org.assertj</groupId>
  <artifactId>assertj-core</artifactId>
  <version>3.20.2</version>
  <scope>test</scope>
</dependency>
```

## 断言描述

通过`as()`方法来设置断言的描述信息，这将会在断言失败时显示出来，注：`as`方法应在开始断言前调用 ，这样才不会因为前面的断言失败而导致`as`方法未被调用。

```java
TolkienCharacter frodo = new TolkienCharacter("Frodo", 33, Race.HOBBIT);

// 一个会失败的断言，会显示：[check Frodo's age] expected:<100> but was:<33>
assertThat(frodo.getAge()).as("check %s's age", frodo.getName())
                          .isEqualTo(100);
```

## 重写错误信息

通过`withFailMessage`方法来自定义错误信息的显示，同样应该在断言开始前调用。

```java
TolkienCharacter frodo = new TolkienCharacter("Frodo", 33, Race.HOBBIT);
TolkienCharacter sam = new TolkienCharacter("Sam", 38, Race.HOBBIT);
// java.lang.AssertionError: should be TolkienCharacter [name=Frodo, age=33, race=HOBBIT]
assertThat(frodo.getAge()).withFailMessage("should be %s", frodo)
                          .isEqualTo(sam);
```

## 集合和数组断言

### 内容检查

| 方法                        | 描述                                                         |
| --------------------------- | ------------------------------------------------------------ |
| `contains`                  | 验证集合中是否包含给定的元素，不考虑顺序。                   |
| `containsOnly`              | 验证集合中是否**只**包含给定的元素，不考虑顺序，会忽略重复的元素。 |
| `containsExactly`           | 同上，但会考虑顺序，不忽略重复的元素。                       |
| `containsExactlyInAnyOrder` | 同上，但不考虑顺序。                                         |
| `containsSequence`          | 验证集合中是否包含指定的序列，且序列之间不能有其它元素。     |
| `containsSubsequence`       | 同上，但允许序列之间有其它元素。                             |
| `containsOnlyOnce`          | 验证集合中是否只包含给定的元素，且不重复。                   |
| `containsAnyOf`             | 验证集合中是否包含给定元素中的一个或多个。                   |

### 检验元素是否满足条件

```java
List<TolkienCharacter> hobbits = list(frodo, sam, pippin);

// 所有的元素都满足条件
assertThat(hobbits).allSatisfy(character -> {
  assertThat(character.getRace()).isEqualTo(HOBBIT);
  assertThat(character.getName()).isNotEqualTo("Sauron");
});

// 至少有一个元素满足条件
assertThat(hobbits).anySatisfy(character -> {
  assertThat(character.getRace()).isEqualTo(HOBBIT);
  assertThat(character.getName()).isEqualTo("Sam");
});

// 没有任何元素满足条件
assertThat(hobbits).noneSatisfy(character -> 
                                assertThat(character.getRace()).isEqualTo(ELF));
```

### 对指定元素断言

```java
// 只能使用对象断言方法
Iterable<TolkienCharacter> hobbits = list(frodo, sam, pippin);
assertThat(hobbits).first().isEqualTo(frodo);//断言第一个元素
assertThat(hobbits).element(1).isEqualTo(sam);//断言索引为1的元素
assertThat(hobbits).last().isEqualTo(pippin);//断言最后一个元素

// 获取到指定元素后，根据类型切换的对应的断言器，从而可用更多的断言方法
Iterable<String> hobbitsName = list("frodo", "sam", "pippin");
// STRING is an InstanceOfAssertFactory from 
//org.assertj.core.api.InstanceOfAssertFactories.STRING
// as() is just synthetic sugar for readability
assertThat(hobbitsName).first(as(STRING))
                       .startsWith("fro")
                       .endsWith("do");
assertThat(hobbitsName).element(1, as(STRING))
                       .startsWith("sa")
                       .endsWith("am");
assertThat(hobbitsName).last(as(STRING))
                       .startsWith("pip")
                       .endsWith("pin");

// 通过第二个参数指定断言器的类型
assertThat(hobbitsName, StringAssert.class).first()
                                           .startsWith("fro")
                                           .endsWith("do");
```

#### 对单元素集合断言

```java
Iterable<String> babySimpsons = list("Maggie");
// 集合中只包含一个元素，获取到这个元素并对其断言
assertThat(babySimpsons).singleElement()
                        .isEqualTo("Maggie");
assertThat(babySimpsons).singleElement(as(STRING))
                        .endsWith("gie");
assertThat(babySimpsons, StringAssert.class).singleElement()
                                            .startsWith("Mag");
```

### 过滤元素

注：以下所有过滤都是指过滤出，即只有满足条件的元素才会被保留。

```java
//传入lambda方法来判断，元素是否需要被过滤，返回true的元素会被保留
assertThat(fellowshipOfTheRing).filteredOn( 
    character -> character.getName().contains("o") )                 
    .containsOnly(aragorn, frodo, legolas, boromir);
```

#### 根据对象的属性进行过滤

```java
// 第一个参数为属性名称，第三个参数为期望的值，私有属性也支持
assertThat(fellowshipOfTheRing).filteredOn("race", HOBBIT)
                               .containsOnly(sam, frodo, pippin, merry);
// 通过属性的属性进行过滤
assertThat(fellowshipOfTheRing).filteredOn("race.name", "Man")
                               .containsOnly(aragorn, boromir);

assertThat(fellowshipOfTheRing).filteredOn("race", notIn(HOBBIT, MAN))
                               .containsOnly(gandalf, gimli, legolas);
assertThat(fellowshipOfTheRing).filteredOn("race", in(MAIA, MAN))
                               .containsOnly(gandalf, boromir, aragorn);
assertThat(fellowshipOfTheRing).filteredOn("race", not(HOBBIT))
                               .containsOnly(gandalf, boromir, aragorn, gimli, legolas);
//可以多次过滤
assertThat(fellowshipOfTheRing).filteredOn("race", MAN)
                               .filteredOn("name", not("Boromir"))
                               .containsOnly(aragorn);
```

#### 通过方法的返回值过滤

```java
assertThat(fellowshipOfTheRing).filteredOn(TolkienCharacter::getRace, HOBBIT)
                               .containsOnly(sam, frodo, pippin, merry);
```

#### 过滤出属性值的`null`的元素

```java
//age属性为null的元素会被保留
assertThat(hobbits).filteredOnNull("age"))
                   .singleElement()
                   .isEqualTo(mysteriousHobbit);
```

#### 通过断言进行过滤

```java
//age小于34的元素会被保留
assertThat(hobbits).filteredOnAssertions(hobbit -> assertThat(hobbit.age).isLessThan(34))
                   .containsOnly(frodo, pippin);
```

### 提取元素的属性

#### 提取单个属性

```java
// 提取出集合中每个元素的 name 属性的值，然后将其作为集合进行断言
assertThat(fellowshipOfTheRing).extracting("name")
                               .contains("Boromir", "Gandalf", "Frodo", "Legolas")
                               .doesNotContain("Sauron", "Elrond");
// 属性的属性也可以提取
assertThat(fellowshipOfTheRing).extracting("race.name")
                               .contains("Man", "Maia", "Hobbit", "Elf");
// 通过lambda方法来提取出属性的值，这样类型安全
assertThat(fellowshipOfTheRing).extracting(TolkienCharacter::getName)
                               .contains("Boromir", "Gandalf", "Frodo", "Legolas");
// 也可以调用map方法，和extracting方法相同，但不接收String参数
assertThat(fellowshipOfTheRing).map(TolkienCharacter::getName)
                               .contains("Boromir", "Gandalf", "Frodo", "Legolas");
//第二个参数指定提取出来的值的类型
assertThat(fellowshipOfTheRing).extracting("name", String.class)
                               .contains("Boromir", "Gandalf", "Frodo", "Legolas")
                               .doesNotContain("Sauron", "Elrond");
```

#### 提取多个属性

```java
import static org.assertj.core.api.Assertions.tuple;
// 提取出 name, age 和 race.name
assertThat(fellowshipOfTheRing).extracting("name", "age", "race.name")
                               .contains(tuple("Boromir", 37, "Man"),
                                         tuple("Sam", 38, "Hobbit"),
                                         tuple("Legolas", 1000, "Elf"));

// 同样可以使用lambda来提取
assertThat(fellowshipOfTheRing).extracting(TolkienCharacter::getName,
                                            tolkienCharacter -> tolkienCharacter.age,
                                            tolkienCharacter -> tolkienCharacter.getRace()
                                           									.getName())
                                .contains(tuple("Boromir", 37, "Man"),
                                          tuple("Sam", 38, "Hobbit"),
                                          tuple("Legolas", 1000, "Elf"));
```

#### 提取并聚合属性

```java
//teamMates的类型是一个集合
//flatExtracting方法会提取出每个元素的teamMates属性，然后将其合并成一个集合
assertThat(reallyGoodPlayers).flatExtracting("teamMates")
                             .contains(pippen, kukoc, jabbar, worthy);
//同样可以使用lambda来提取
assertThat(reallyGoodPlayers).flatExtracting(BasketBallPlayer::getTeamMates)
                             .contains(pippen, kukoc, jabbar, worthy);

// flatMap也具有相同的功能，但不接收String参数
assertThat(reallyGoodPlayers).flatMap(BasketBallPlayer::getTeamMates)
                             .contains(pippen, kukoc, jabbar, worthy);

// 如果用extracting方法来处理，则每个元素都是一个list
assertThat(reallyGoodPlayers).extracting("teamMates")
                             .contains(list(pippen, kukoc), list(jabbar, worthy));

// 也可以用flatExtracting方法来提取多个参数
assertThat(fellowshipOfTheRing).flatExtracting("name", "race.name")
                               .contains("Frodo", "Hobbit", "Legolas", "Elf");
assertThat(fellowshipOfTheRing).flatExtracting(TolkienCharacter::getName,
                                               tc -> tc.getRace().getName())
                               .contains("Frodo", "Hobbit", "Legolas", "Elf");
```

### 自定义比较器比较元素

```java
// 对于对象而言，默认会使用其equals方法来比较两个对象是否相同，从而确定其是否包含在集合中，
//也可以通过 usingElementComparator方法 来自定义比较方法
assertThat(fellowshipOfTheRing).usingElementComparator(
    (t1, t2) -> t1.getRace().compareTo(t2.getRace()))
                               .contains(sauron);
```

## 断言异常

当程序抛出异常时，断言其是否是期望的异常。

### 断言异常消息

```java
Throwable throwable = new IllegalArgumentException("wrong amount 123");

assertThat(throwableWithMessage).hasMessage("wrong amount 123")
                                .hasMessage("%s amount %d", "wrong", 123)
                                // 检查开始的内容
                                .hasMessageStartingWith("wrong")
                                .hasMessageStartingWith("%s a", "wrong")
                                // 检查包含的内容
                                .hasMessageContaining("wrong amount")
                                .hasMessageContaining("wrong %s", "amount")
                                .hasMessageContainingAll("wrong", "amount")
                                // 检查结尾的内容
                                .hasMessageEndingWith("123")
                                .hasMessageEndingWith("amount %s", "123")
                                // 通过正则检查
                                .hasMessageMatching("wrong amount .*")
                                // 检查是否不包含某些内容
                                .hasMessageNotContaining("right")
                                .hasMessageNotContainingAny("right", "pri
```

### 检查异常的原因和根本原因

```java
NullPointerException cause = new NullPointerException("boom!");
Throwable throwable = new Throwable(cause);

assertThat(throwable).hasCause(cause)
                     // hasCauseInstanceOf will match inheritance.
                     .hasCauseInstanceOf(NullPointerException.class)
                     .hasCauseInstanceOf(RuntimeException.class)
                     // hasCauseExactlyInstanceOf will match only exact same type
                     .hasCauseExactlyInstanceOf(NullPointerException.class);
// navigate before checking
assertThat(throwable).getCause()
                     .hasMessage("boom!")
                     .hasMessage("%s!", "boom")
                     .hasMessageStartingWith("bo")
                     .hasMessageEndingWith("!")
                     .hasMessageContaining("boo")
                     .hasMessageContainingAll("bo", "oom", "!")
                     .hasMessageMatching("b...!")
                     .hasMessageNotContaining("bam")
                     .hasMessageNotContainingAny("bam", "bim")
                     // isInstanceOf will match inheritance.
                     .isInstanceOf(NullPointerException.class)
                     .isInstanceOf(RuntimeException.class)
                     // isExactlyInstanceOf will match only exact same type
                     .isExactlyInstanceOf(NullPointerException.class);
```

```java
NullPointerException rootCause = new NullPointerException("null!");
Throwable throwable = new Throwable(new IllegalStateException(rootCause));

// direct root cause check
assertThat(throwable).hasRootCause(rootCause)
                     .hasRootCauseMessage("null!")
                     .hasRootCauseMessage("%s!", "null")
                     // hasRootCauseInstanceOf will match inheritance
                     .hasRootCauseInstanceOf(NullPointerException.class)
                     .hasRootCauseInstanceOf(RuntimeException.class)
                     // hasRootCauseExactlyInstanceOf will match only exact same type
                     .hasRootCauseExactlyInstanceOf(NullPointerException.class);

// navigate to root cause and check
assertThat(throwable).getRootCause()
                     .hasMessage("null!")
                     .hasMessage("%s!", "null")
                     .hasMessageStartingWith("nu")
                     .hasMessageEndingWith("!")
                     .hasMessageContaining("ul")
                     .hasMessageContainingAll("nu", "ull", "l!")
                     .hasMessageMatching("n...!")
                     .hasMessageNotContaining("NULL")
                     .hasMessageNotContainingAny("Null", "NULL")
                     // isInstanceOf will match inheritance.
                     .isInstanceOf(NullPointerException.class)
                     .isInstanceOf(RuntimeException.class)
                     // isExactlyInstanceOf will match only exact same type
                     .isExactlyInstanceOf(NullPointerException.class);
```

### 捕获异常

```java
//通过catchThrowable来捕获方法抛出的异常，没有异常时会返回null
Throwable thrown = catchThrowable(() -> System.out.println(names[9]));

//通过catchThrowableOfType来捕获指定类型的异常
class TextException extends Exception {
   int line;
   int column;

   public TextException(String msg, int line, int column) {
     super(msg);
     this.line = line;
     this.column = column;
   }
 }

 TextException textException = catchThrowableOfType(
     () -> { throw new TextException("boom!", 1, 5); },TextException.class);

 // 这将会出错，因为类型不匹配
 catchThrowableOfType(() -> { throw new TextException("boom!", 1, 5); }, 
                      RuntimeException.class);
```

### assertThatThrownBy

自动捕获异常然后进行断言，如果没有异常抛出，则断言失败。

```java
assertThatThrownBy(() -> { throw new Exception("boom!"); }).isInstanceOf(Exception.class)
                                                           .hasMessageContaining("boom");
```

### assertThatExceptionOfType

自动捕获指定类型的异常然后进行断言。

```java
assertThatExceptionOfType(IOException.class).isThrownBy(
    								() -> { throw new IOException("boom!"); })
                                            .withMessage("%s!", "boom")
                                            .withMessageContaining("boom")
                                            .withNoCause();
```

这里也提供了一些用于获取常见异常的断言方法：

- `assertThatNullPointerException`
- `assertThatIllegalArgumentException`
- `assertThatIllegalStateException`
- `assertThatIOException`

```java
assertThatIOException().isThrownBy(() -> { throw new IOException("boom!"); })
                       .withMessage("%s!", "boom")
                       .withMessageContaining("boom")
                       .withNoCause();
```

### 断言无异常抛出

```java
// standard style
assertThatNoException().isThrownBy(() -> System.out.println("OK"));
// BDD style
thenNoException().isThrownBy(() -> System.out.println("OK"));

// standard style
assertThatCode(() -> System.out.println("OK")).doesNotThrowAnyException();
// BDD style
thenCode(() -> System.out.println("OK")).doesNotThrowAnyException();
```

### BDD风格

```java
// GIVEN
String[] names = { "Pier ", "Pol", "Jak" };
// WHEN
Throwable thrown = catchThrowable(() -> System.out.println(names[9]));
// THEN
then(thrown).isInstanceOf(ArrayIndexOutOfBoundsException.class)
            .hasMessageContaining("9");
```

## 逐字段比较

```java
Person sherlock = new Person("Sherlock", 1.80);
sherlock.home.ownedSince = new Date(123);
sherlock.home.address.street = "Baker Street";
sherlock.home.address.number = 221;

Person sherlock2 = new Person("Sherlock", 1.80);
sherlock2.home.ownedSince = new Date(123);
sherlock2.home.address.street = "Baker Street";
sherlock2.home.address.number = 221;

// 断言为真，因为这会逐个逐个比较字段的内容
assertThat(sherlock).usingRecursiveComparison()
                    .isEqualTo(sherlock2);

// 断言为假，这只会比较对象的引用
assertThat(sherlock).isEqualTo(sherlock2);
```

### isNotEqualTo

```java
// equals not overridden in TolkienCharacter
TolkienCharacter frodo = new TolkienCharacter("Frodo", 33, HOBBIT);
TolkienCharacter frodoClone = new TolkienCharacter("Frodo", 33, HOBBIT);
TolkienCharacter youngFrodo = new TolkienCharacter("Frodo", 22, HOBBIT);

// 只会比较对象的引用
assertThat(frodo).isNotEqualTo(frodoClone);

// 断言失败，因为这会逐个比较对象属性的值
assertThat(frodo).usingRecursiveComparison()
                 .isNotEqualTo(frodoClone);

// 断言成功，因为有一个属性不一样
assertThat(frodo).usingRecursiveComparison()
                 .isNotEqualTo(youngFrodo);
```

### 比较时忽略属性

```java
// 通过属性名忽略
assertThat(sherlock).usingRecursiveComparison()
                    .ignoringFields("name", "home.address.street")
                    .isEqualTo(moriarty);
assertThat(sherlock).usingRecursiveComparison()
                    .ignoringFields("name", "home")
                    .isEqualTo(moriarty);

// 通过正则表达式匹配属性名来忽略
assertThat(sherlock).usingRecursiveComparison()
                    .ignoringFieldsMatchingRegexes(".*me")
                    .isEqualTo(moriarty);

// 忽略值为null的属性
assertThat(sherlock).usingRecursiveComparison()
                    .ignoringActualNullFields()
                    .isEqualTo(moriarty);

// 通过类型来忽略
assertThat(sherlock).usingRecursiveComparison()
                    .ignoringFieldsOfTypes(double.class, Address.class)
                    .isEqualTo(tallSherlock);
```

### 忽略重写的equals方法

当对象中有一个属性重写的equals方法时，默认会使用equals方法来比较时否相同，而不再逐字段比较。

可以通过以下方法来忽略某一个或多个属性重写的equals方法：

- `ignoringOverriddenEqualsForTypes(Class...)`:通过类型指定对应的属性
- `ignoringOverriddenEqualsForFields(String...)`:通过属性名来指定对应的属性
- `ignoringOverriddenEqualsForFieldMatchRegexes(String...)`:通过正则表达式来指定对应的属性
- `ignoringAllOverriddenEquals()`：除了java中定义的类型，所有重写equals方法的属性都会逐字段比较。

```java
public class Person {
  String name;
  double height;
  Home home = new Home();
}

public class Home {
  Address address = new Address();
}

public static class Address {
  int number;
  String street;

  // 只比较number
  @Override
  public boolean equals(final Object other) {
    if (!(other instanceof Address)) return false;
    Address castOther = (Address) other;
    return Objects.equals(number, castOther.number);
  }
}

Person sherlock = new Person("Sherlock", 1.80);
sherlock.home.address.street = "Baker Street";
sherlock.home.address.number = 221;

Person sherlock2 = new Person("Sherlock", 1.80);
sherlock2.home.address.street = "Butcher Street";
sherlock2.home.address.number = 221;

//断言为真，因为home.address.street在address重写的equals中并没有比较
assertThat(sherlock).usingRecursiveComparison()
                    .isEqualTo(sherlock2);

//断言为假，这会忽略Address中重写的equals方法，而逐字段比较
assertThat(sherlock).usingRecursiveComparison()
                    .ignoringOverriddenEqualsForTypes(Address.class)
                    .isEqualTo(sherlock2);
```

### 忽略被期望对象中为null的属性

```java
Person sherlock = new Person("Sherlock", 1.80);
sherlock.home.address.street = "Baker Street";
sherlock.home.address.number = 221;

Person noName = new Person(null, 1.80);
noName.home.address.street = null;
noName.home.address.number = 221;

// 断言为真，noName中的name和street为null,比较时会被忽略
assertThat(sherlock).usingRecursiveComparison()
                    .ignoringExpectedNullFields()
                    .isEqualTo(noName);

// 断言为假，sherlock中的name的street不为null,比较时不会被忽略
assertThat(noName).usingRecursiveComparison()
                  .ignoringExpectedNullFields()
                  .isEqualTo(sherlock);
```

## 软断言

软断言可以收集所有的断言错误而不会在一条断言失败后就停止执行。

```java
SoftAssertions softly = new SoftAssertions(); 

softly.assertThat("George Martin").as("great authors").isEqualTo("JRR Tolkien");  
softly.assertThat(42).as("response to Everything").isGreaterThan(100); 
softly.assertThat("Gandalf").isEqualTo("Sauron"); 

//调用assertAll方法报告断言错误！！！
softly.assertAll(); 
```

### BDD风格软断言

```java
BDDSoftAssertions softly = new BDDSoftAssertions();

softly.then("George Martin").as("great authors").isEqualTo("JRR Tolkien");
softly.then(42).as("response to Everything").isGreaterThan(100);
softly.then("Gandalf").isEqualTo("Sauron");

// Don't forget to call assertAll() otherwise no assertion errors are reported!
softly.assertAll();
```

### Junit5扩展

AssertJ针对Junit5提供了`SoftAssertionsExtension`扩展：

- 自动在每个测试完成后调用`assertAll()`方法
- 通过`@InjectSoftAssertions`注解来初始化`softAssertions`属性
- 如果测试方法的参数需要`softAssertions`，则会自动传入。

```java
import org.assertj.core.api.SoftAssertions;
import org.assertj.core.api.junit.jupiter.SoftAssertionsExtension;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

@ExtendWith(SoftAssertionsExtension.class)
public class JUnit5SoftAssertionsExtensionAssertionsExamples {

  @InjectSoftAssertions
  private SoftAssertions softly;

  @Test
  public void chained_soft_assertions_example() {
    String name = "Michael Jordan - Bulls";
    softly.assertThat(name).startsWith("Mi")
                           .contains("Bulls");
    // no need to call softly.assertAll(), this is done by the extension
  }
}
```

### 可自动关闭的软断言

```java
@Test
void auto_closeable_soft_assertions_example() {
  try (AutoCloseableSoftAssertions softly = new AutoCloseableSoftAssertions()) {
  softly.assertThat("George Martin").as("great authors").isEqualTo("JRR Tolkien");  
  softly.assertThat(42).as("response to Everything").isGreaterThan(100); 
  softly.assertThat("Gandalf").isEqualTo("Sauron"); 
    // no need to call assertAll, this is done when softly is closed.
  }
}
```

### 使用assertSoftly静态方法

```java
@Test
void assertSoftly_example() {
  SoftAssertions.assertSoftly(softly -> {
    softly.assertThat("George Martin").as("great authors").isEqualTo("JRR Tolkien");
    softly.assertThat(42).as("response to Everything").isGreaterThan(100);
    softly.assertThat("Gandalf").isEqualTo("Sauron");
    // no need to call assertAll(), assertSoftly does it for us.
  });
}
```

## AssertJ-DB

