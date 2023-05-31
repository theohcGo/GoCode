/*
	Go语言中的反射（reflection）是指在程序运行期对程序本身进行访问和修改的能力。
	通过反射，我们可以在运行期间动态地获取一个变量的类型信息和值信息，并且可以通过反射来修改变量的值。

	学习Go语言中的反射，可以先了解reflect包中提供的相关函数和方法，
	例如reflect.TypeOf()、reflect.ValueOf()、reflect.New()等。同时，可以通过实践编写一些使用反射的小程序来加深对反射的理解。

	在一些开源项目中，反射通常用于实现一些通用的功能，
	例如序列化和反序列化、ORM框架等。下面是一个使用反射实现JSON序列化的示例代码：

	复制
*/
package main

import (
    "encoding/json"
    "fmt"
    "reflect"
)

type Person struct {
    Name string
    Age  int
}


func main() {
    p := Person{Name: "张三", Age: 20}
    // 序列化为json字符串
    b, err := json.Marshal(p)
    if err != nil {
        fmt.Println("json.Marshal error:", err)
        return
    }
    fmt.Println(string(b))

    // 使用反射修改Person的Name字段
    v := reflect.ValueOf(&p).Elem()
    v.FieldByName("Name").SetString("李四")
    fmt.Println(p)
}
	在上面的代码中，我们首先定义了一个Person结构体，并且使用json.Marshal()函数将其序列化为JSON字符串。
	接着，我们使用反射获取了Person结构体的Value对象，
	并且通过Value对象的FieldByName()方法获取了Name字段的Value对象，
	并且使用SetString()方法修改了Name字段的值。最后，我们再次打印Person结构体的值，发现Name字段已经被修改为了"李四"。


