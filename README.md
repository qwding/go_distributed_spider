# go_sample_spider
the first time writing spider,this is a sample example

使用说明:\n -urls 和-targets输入要爬的url和关键词，目前只支持这个，整个过程可以随时ctrl+c掉，并且结果回存到当前目录下的has.txt和ret.txt，主要是保存符合条件的和所有爬多的url。

参数
* -help      ==>   see help detail.
* -file      ==>   save result into file.eg: -file=test.txt
* -timeout   ==>   dail tcp connect timeout.default is 2s. eg:-timeout=3
* -urls      ==>   give spider urls,use ',' split the urls. eg:-urls=www.baidu.com,www.xyz.com
* -targets   ==>   give spider targets,use ',' split the target. eg:-targets=golang,docker
