# go_sample_spider

更适合说这是一个初学者学习写爬虫的过程

### 单线程实现

主要是用一个线程爬网页，把整个流程实现了一遍，具体的匹配模块和下载都没有，只是简单的记录url（匹配成功和全部的url两个文件）。想细化这部分应该并不难，写个接口，可以支持更全面。

实现心得：

第一次写爬虫，先简单写一下，目前只支持爬关键字， 屏蔽了图片。

使用说明: -urls 和-targets输入要爬的url和关键词，用逗号隔开输入多个，整个过程可以随时ctrl+c掉，并且结果回存到当前目录下的has.txt和ret.txt，主要是保存符合条件的和所有爬多的url。

为了减少IO，所以只在最后存储到文件里。

参数
* -help      ==>   see help detail.
* -file      ==>   save result into file.eg: -file=test.txt
* -timeout   ==>   dail tcp connect timeout.default is 2s. eg:-timeout=3
* -urls      ==>   give spider urls,use ',' split the urls. eg:-urls=www.baidu.com,www.xyz.com
* -targets   ==>   give spider targets,use ',' split the target. eg:-targets=golang,docker


### 多线程实现 

单线程跑的真心慢， 并且很多url请求很慢，当然必须呀换成多线程了

具体实现：

* 并且将单线程里故意用递归实现的爬取逻辑改成了循环。
* 制造了一个线程池，保持线程数和cpu核数相同， 尽量高可用cpu。
* 每个线程申请时候为未申请url的 1/cpuNum.超过100个就按100个算。
* 记录的url数组和为线程申请url都设置成了全局变量。并且用各种锁来帮助写操作。
* 每个线程结束后将遍历的url保存到文件里，所以ctrl+c后的正在跑的线程结果是得不到的。
* 程序出口为所有spider线程退出，并且没有未分配的url。
 
实现心得：

写多线程这个学到不少

* 锁的使用过程，被各种死锁，各种return前一定把锁解开，不然就锁，锁一多还容易乱，认为一个锁只管一个变量比较好管理。
* 程序性能，看了pprof包和输出gc信息。想好好调调关于性能方面， 看看能不能优化性能。


程序代码写的挺渣的，没有具体分层，主要也是对爬虫不了解，写前没有考虑那么多，逻辑都罗列在了一起。


### 分布式爬虫

当然对分布式这个概念还是模糊的，要怎么实现还不太懂，做成分布式以后再好好弄弄

### 后记：主要还是因为没有写过，也算是对写爬虫的一个过程吧~
