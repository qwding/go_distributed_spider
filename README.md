# go_sample_spider
第一次写爬虫，先简单写一下，目前只支持爬关键字， 屏蔽了图片。

使用说明: -urls 和-targets输入要爬的url和关键词，用逗号隔开输入多个，整个过程可以随时ctrl+c掉，并且结果回存到当前目录下的has.txt和ret.txt，主要是保存符合条件的和所有爬多的url。

为了减少IO，所以只在最后存储到文件里。

参数
* -help      ==>   see help detail.
* -file      ==>   save result into file.eg: -file=test.txt
* -timeout   ==>   dail tcp connect timeout.default is 2s. eg:-timeout=3
* -urls      ==>   give spider urls,use ',' split the urls. eg:-urls=www.baidu.com,www.xyz.com
* -targets   ==>   give spider targets,use ',' split the target. eg:-targets=golang,docker
