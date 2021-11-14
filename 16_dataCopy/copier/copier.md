#copier源码分析

##静态参数
```go
// These flags define options for tag handling
const (
	//结构体中标注了“must”标签的字段，必须被复制值，否则视为error
	tagMust uint8 = 1 << iota

	// 和tagMust配套使用，如果设置标签“nopanic”，则如果不满足must的条件，不直接报错，而是返回error代替
	tagNoPanic

	// 标签为“-”,设置该标签的字段直接忽略复制
	tagIgnore

	// 这个不是标签标识，而是字段复制的标识，在结构体复制结束之后，设置flag为已经复制.依据这个字段来判断是否复制成功
	hasCopied
)
```

##整体设计思路

    1.不可寻址和Invalid的数据直接报错或者返回

    2.判断两个数据结构是不是map，进行map的处理
    
    3.数组与结构体的处理，按照类型进行数据遍历
    
        a.循环所有字段，解析tag
        
        b.判断是否忽略不复制，不复制则跳过
        
        c.根据字段名进行赋值
        
        d.根据方法名（同名）进行赋值
        
        5.赋值成功之后设置本字段的赋值成功标识
        
        6.循环所有标签，判断不满足“must”标签的情况，进行相应处理

##辅助方法说明
###1.获取实际的Type和Value
在go中，如果一个参数是 *Struct类型的，也就是指针类型，当用这个方法去调用方法获取结构体属性的时候会报错，
比如我这边通过一个指针类型去调用FieldByName方法，就会出现下面的错误：
```go
--- FAIL: TestIndirect (0.00s)
panic: reflect: FieldByName of non-struct type [recovered]
	panic: reflect: FieldByName of non-struct type
```
在Go当中，如果是指针类型，可以通过【.Elem】方法获取对应的实际值。所以这边需要两个方法：

    1、通过type判断是否指针类型，返回具体的结构体类型
    2、通过value判断是否指针类型，发挥具体的结构体数据
```go

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}
 
func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType

```
###2.Tag处理
    1、解析tag字符串
```go
// parseTags Parses struct tags and returns uint8 bit flags.
func parseTags(tag string) (flags uint8) {
	for _, t := range strings.Split(tag, ",") {
		switch t {
		case "-":
			flags = tagIgnore
			return
		case "must":
			flags = flags | tagMust
		case "nopanic":
			flags = flags | tagNoPanic
		}
	}
	return
}
```
    2、解析Field对应的tag，获得每个字段对应的tag条件
```go
// getBitFlags Parses struct tags for bit flags.
func getBitFlags(toType reflect.Type) map[string]uint8 {
	//read note 存储的结构是  FieldName->tag对应的二进制数据(tag标签转换成程序标识)
	flags := map[string]uint8{}
	//read note 根据结构体的类型获取对应的Field切片
	toTypeFields := deepFields(toType)
 
	// Get a list dest of tags
	//read note 循环Field切片，获取切片对应的tag数据
	for _, field := range toTypeFields {
		tags := field.Tag.Get("copier") //tag标签是【copier】
		if tags != "" {
			//read note tag标签转换成程序处理标识（这边也是使用二进制的处理方式）
			flags[field.Name] = parseTags(tags)
		}
	}
	return flags
}
```
    3.获取结构体Field切片

```go
//read note 根据结构体类型，获取结构体对应的Field切片，注意这边的【Anonymous】表示的匿名变量，匿名变量的Field这边需要特殊处理
func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
 
	//read note 判断是不是结构体,只能对结构体进行处理
	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		//read note 循环处理对应的所有Field
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			//read note 对【嵌入（匿名）字段】结构体 的所有结构体进行添加
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}
 
	return fields
}
```
    4. 检查结构体复制结果
```go
// checkBitFlags Checks flags for error or panic conditions.
func checkBitFlags(flagsList map[string]uint8) (err error) {
	// Check flag conditions were met
	//read note 循环map（FieldName->tag对应的二进制数据）
	for name, flags := range flagsList {
		//read note 如果字段没有被复制
		if flags&hasCopied == 0 {
			switch {
			case flags&tagMust != 0 && flags&tagNoPanic != 0:
				//read note 处理1：返回错误信息
				err = fmt.Errorf("Field %s has must tag but was not copied", name)
				return
			case flags&(tagMust) != 0:
				//read note 处理2：直接报错
				panic(fmt.Sprintf("Field %s has must tag but was not copied", name))
			}
		}
	}
	return
}
```
    5.对结构体进行设值
处理步骤
![](.copier_images/struct_set_value_process.png)

    a.处理非零值的情况
        *如果Result字段是指针类型（前置处理）
            Origin字段是指针类型，并且为nil，直接设置Result字段为nil，返回
            Result字段是nil，赋值初值(包含结构体和指针的处理)
        *类型转换处理
            如果Origin类型可以转换为Result类型，进行结构体赋值
            如果Result是sql.Scanner类型，调用Scan方法
            如果Origin是指针类型，调用set(Result,Origin.Elem())方法，注意这边指针需要调用.Elem方法（递归）
            返回false
    b.零值的情况直接返回true
```go
func set(to, from reflect.Value) bool {
 
	//read note IsValid返回是否非零值的结果.所以这边的处理是针对非零值
	if from.IsValid() {
 
		//read note 前置条件：处理to的类型，如果from为空，则直接设置空值返回
 
		//	to是指针类型特殊处理
		if to.Kind() == reflect.Ptr {
			// set `to` to nil if from is nil
			//read note 如果from是空，则直接设置to为零值返回
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				//read note 如果to是nil且不满足上面的from为nil的条件，这个时候要给to设置默认值
				to.Set(reflect.New(to.Type().Elem()))
			}
			//read note 指针的转换处理
			to = to.Elem()
		}
 
		//read note from和to类型的转换处理，这边当from是ptr类型的时候，会调用set进行递归处理
 
		//read note 如果类型可以进行转换，则要设置对应的值（具体什么类型可以转换需要看一下源码，这里不多赘述）
		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			//read note sql.Scanner 这个不知道具体是干嘛的.
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			//read note from是指针类型，处理成结构体进行赋值（相当于递归再往下走）
			return set(to, from.Elem())
		} else {
			//read note 其他不能转换的直接返回false
			return false
		}
	}
 
	//read note 零值直接返回true，零值不处理
	return true
}
```

##copy主方法copier(toValue, fromValue, opt)说明
###参数说明
```go
	var (
		isSlice bool
		amount  = 1
		from    = indirect(reflect.ValueOf(fromValue))
		to      = indirect(reflect.ValueOf(toValue))
	)
```
###两个类型都是map的处理
流程

    1、判断key的类型是否可以转换，如果不能转换直接返回
    
    2、如果Result对应的Map为空，则初始化
    
    3、循环遍历Map的所有key，先复制key，如果复制失败则跳过该字段
    
    4、复制value，如果复制失败则跳过该字段（这边的复制key的复制是不一样的,因为key是不可变的，value需要再调用Copy进行复制）
    
    5、如果上面的复制操作都成功了，则对map进行key到value的映射
```go
//read note from和to都是map
	if fromType.Kind() == reflect.Map && toType.Kind() == reflect.Map {
		//read note 判断map的key的结构类型是否可以转换，因为这边已经判断是map，所以直接通过 .key来获取，不担心是否报错
		if !fromType.Key().ConvertibleTo(toType.Key()) {
			return
		}
		//read note 判断要转换的Map是否为空，为空则进行初始化
		if to.IsNil() {
			to.Set(reflect.MakeMapWithSize(toType, from.Len()))
		}
		//read note 遍历Map的所有key
		for _, k := range from.MapKeys() {
			//read note 根据to的key类型创建一个新的key
			toKey := indirect(reflect.New(toType.Key()))
			//read note 设置key的值
			if !set(toKey, k) {
				continue
			}
 
			//read note 设置value值
			toValue := indirect(reflect.New(toType.Elem()))
			if !set(toValue, from.MapIndex(k)) {
				//read note 对嵌套的结构体进行copy
				err = Copy(toValue.Addr().Interface(), from.MapIndex(k).Interface())
				if err != nil {
					continue
				}
			}
			to.SetMapIndex(toKey, toValue)
		}
	}
```
###只有一个类型是结构体的处理
    如果被复制的结构或者复制的结构有一个不是struct，则直接返回，不进行处理
```go
	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		// skip not supported type
		return
	}
```
###判断数组设置标识
    上面的操作已经把map和非结构体的情况处理掉，下面就是结构体\结构体数组的处理。
    
    这边要先设置对应的标识，isSlice是用来标识Result是不是数组的，所以这边会设置对应的标识，并且获取数组对应的长度
```go
	//read note 切片处理：设置切片的长度
	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}
```
根据amount进行循环处理

    1、通过isSlice标识，获取对应的Origin结构体数据，如果是数组则根据index获取，否则直接使用Origin
    
    2、获取字段对应的Field的tag，递归处理
    
    3、复制处理
    
    4、通过isSlice标识，对复制之后的结果进行处理
    
    5、通过map进行状态的校验
第一个步骤的代码如下，这边的处理就是根据index进行当前处理的结构体数据的获取：
```go
	//read note 循环被复制的切片的长度
	for i := 0; i < amount; i++ {
		var dest, source reflect.Value
 
		//read note Result的结果是数组
		if isSlice {
			// source
			// read note 如果Origin的类型是数组，需要根据index进行获取结构体
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				// read note 如果Origin的类型是结构体，直接获取该结构体
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			//read note Result的结果不是数组，直接获取结构体
			source = indirect(from)
			dest = indirect(to)
		}
```
第二个步骤的处理，是根据上面的 getBitFlags 方法来获取的，具体的代码如下：
```go
		// Get tag options
		//read note 获取tag的所有标签
		tagBitFlags := map[string]uint8{}
		if dest.IsValid() {
			//read note 根据结构体type获取所有的Field对应的tag标签
			tagBitFlags = getBitFlags(toType)
		}
```
第三个步骤的处理，这边也是根据类型获取对应的Field进行循环的处理，但是这边会比较特殊的地方就是copier包支持的两个功能

    1、复制的结果结构体，和方法同名的字段需要复制方法的返回值
    
    2、复制的结果结构体，和被复制结构体字段同名的方法，需要进行设置（没有返回值）
看一下 Origin的方法和Result字段同名的处理
```go
                //read note Origin的方法和Result字段同名的处理
 
				if fromField := source.FieldByName(name); fromField.IsValid() && !shouldIgnore(fromField, ignoreEmpty) {
					// has field
 
					//read note 根据名称获取Result结构体的字段.
					if toField := dest.FieldByName(name); toField.IsValid() {
						if toField.CanSet() {
							if !set(toField, fromField) {
								if err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil {
									return err
								}
							} else {
								//read note 赋值完成，设置对应的标识
								if fieldFlags != 0 {
									// Note that a copy was made
 
									//read note 设置复制标识
									tagBitFlags[name] = fieldFlags | hasCopied
								}
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						//read note 通过名称找到对应的Method
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(name)
						} else {
							toMethod = dest.MethodByName(name)
						}
						//read note 【被转换对象的方法调用】的校验还比较严格,这边可以看出来只能有一个字段，并且字段类型要对应上
						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							//read note 调用声明的方法
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
```

看一下 Result方法 与Origin字段同名的处理
```go
            // Copy from method to field
			//read note 处理目标结构体的方法，目标结构体的方法要和被复制结构体的字段名一致，就是这边控制的
			for _, field := range deepFields(toType) {
				name := field.Name
 
				//read note 根据Result的字段，获取Origin同名的方法
				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(name)
				} else {
					fromMethod = source.MethodByName(name)
				}
 
				//read note 如果方法符合规则，没有入参，有一个出参，则进行对应方法的调用处理
				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 && !shouldIgnore(fromMethod, ignoreEmpty) {
					if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							//read note 进行字段的设值
							set(toField, values[0])
						}
					}
				}
			}
```

第三步骤所有的代码如下
```go
        //read note 如果是非零值
		if source.IsValid() {
			//read note 获取结构体的所有Field的信息（数组）
			fromTypeFields := deepFields(fromType)
			// fmt.Printf("%#v", fromTypeFields)
			// Copy from field to field or method
 
			//todo 循环所有的Field处理
			for _, field := range fromTypeFields {
				name := field.Name
 
				// Get bit flags for field
				//read note 根据name获取tag数据
				fieldFlags, _ := tagBitFlags[name]
 
				// Check if we should ignore copying
				//read note ignore标签对结构体的影响处理
				if (fieldFlags & tagIgnore) != 0 {
					continue
				}
 
				//read note Origin的方法和Result字段同名的处理
 
				if fromField := source.FieldByName(name); fromField.IsValid() && !shouldIgnore(fromField, ignoreEmpty) {
					// has field
 
					//read note 根据名称获取Result结构体的字段.
					if toField := dest.FieldByName(name); toField.IsValid() {
						if toField.CanSet() {
							if !set(toField, fromField) {
								if err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil {
									return err
								}
							} else {
								//read note 赋值完成，设置对应的标识
								if fieldFlags != 0 {
									// Note that a copy was made
 
									//read note 设置复制标识
									tagBitFlags[name] = fieldFlags | hasCopied
								}
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						//read note 通过名称找到对应的Method
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(name)
						} else {
							toMethod = dest.MethodByName(name)
						}
						//read note 【被转换对象的方法调用】的校验还比较严格,这边可以看出来只能有一个字段，并且字段类型要对应上
						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							//read note 调用声明的方法
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}
 
			//read note Result方法 与Origin字段同名的处理
 
			// Copy from method to field
			//read note 处理目标结构体的方法，目标结构体的方法要和被复制结构体的字段名一致，就是这边控制的
			for _, field := range deepFields(toType) {
				name := field.Name
 
				//read note 根据Result的字段，获取Origin同名的方法
				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(name)
				} else {
					fromMethod = source.MethodByName(name)
				}
 
				//read note 如果方法符合规则，没有入参，有一个出参，则进行对应方法的调用处理
				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 && !shouldIgnore(fromMethod, ignoreEmpty) {
					if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							//read note 进行字段的设值
							set(toField, values[0])
						}
					}
				}
			}
		}
```
第四个步骤的处理，是对最终的复制结果进行处理
```go
    //read note 转换结果Result是切片的处理：分成两种情况，被复制的是 结构体指针 和 结构体
		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
```
第五个步骤是校验map，校验最终的复制结果。处理返回error.
```go
    //read note 这边是不是会有一个问题，就是err是不是会被覆盖，前面的字段有错误，最后一个没有错误则会覆盖之前的error
		err = checkBitFlags(tagBitFlags)
```