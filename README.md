# go_distributed_spider

之前写过单线程和多线程爬虫，放在了otherversion文件夹下

# 分布式爬虫

##主要思想是：

* 1个master，和多个slave，slave可以动态添加。
* master主要负责任务的分发和整理。slave负责将master分发的任务进行处理。
* master将任务分配给slave，slave将任务进行遍历查询，查询以后，将符合条件的MatchUrl和重新获得的NewUrl集合返回给master。
* slave的match工作主要为正则表达式。
* 初始url组 和 目标关键字组 都配置在config，运行程序需要加相关参数。
* 如果为master需要加-m参数。指定配置文件加-f，配置文件可以写多个方案，并用-e指定跑哪个方案。配置文件可以参考config目录下的config.json



## 改进
* slave处理任务为单线程，可以考虑改成多线程
* 之后会做成docker容器，写个Dockerfile，更容易跑
* 对MatchUrl没有做任何处理，单独保存起来。
