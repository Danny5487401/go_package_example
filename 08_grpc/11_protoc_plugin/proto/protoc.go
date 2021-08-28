package proto

/*
 在go中使用google protobuf，有两个可选用的包: goprotobuf（go官方出品）和gogoprotobuf(gogo组织出品^_^)
 gogoprotobuf能够完全兼容google protobuf。而且经过我一些测试，它生成大代码质量确实要比goprotobuf高一些，主要是它在goprotobuf之上extend了一些option。这些option也是有级别区分的，有的option只能修饰field，有的可以修饰enum，有的可以修饰message，有的是修饰package（即对整个文件都有效

*/
