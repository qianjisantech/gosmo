### 覆盖率统计回放接入

<br>

覆盖率统计回放，即流量回放的同时，支持被测代码覆盖率统计，并生成覆盖率报告。

<br>

#### 一、接入方式

回放前提：已经完成流量录制。[录制接入文档](../recorder/README.md)

对于服务启动阶段有TCP请求的SUT，如初始化连接池等，推荐 服务启动顺序：
* 先启动Replayer-Agent 
* 再启动SUT

##### 1. 配置并启动Replayer-Agent

同 [回放接入-1. 配置并启动Replayer-Agent](./README.md#1-配置并启动replayer-agent)

##### 2. 配置并启动SUT

> 同样需要使用定制的golang。

首先，配置定制版golang环境。同 [回放接入-2. 配置并启动SUT](./README.md#2-配置并启动sut) 内定制版golang的配置操作。
```shell script
curl https://raw.githubusercontent.com/didi/sharingan-go/recorder/install/go1.13 | sh
&& export GOROOT=/tmp/recorder-go1.13
&& export PATH=$GOROOT/bin:$PATH
```
然后，新增 [main_test.go](../../replayer-agent/install/codeCov/main_test.go) 文件到SUT的根目录。(对于[支持flag的SUT](#2-支持sut使用flag)，仅此步不同)

最后，使用go test 命令编译代码 并 启动SUT服务。
> 温馨提示：
> 默认统计全模块覆盖率，如需缩小统计范围，请修改-coverpkg指定特定目录，可以提高代码覆盖率
```shell script
go test -gcflags="all=-N -l" -tags="replayer" -v -c -covermode=count -coverpkg ./...
```
> 相比普通回放，启动命令 多了一个参数 -test.coverprofile

##### 可选环境变量**BAN_SYSTEM_TEST**, 用以标识是否关闭集成测试覆盖率支持，设置为1时为原始test流程(如跑单测等流程) 
>示例：export BAN_SYSTEM_TEST=1 #关闭集成测试

```shell script
#linux下启动：
nohup ./$binName.test -test.coverprofile=/tmp/ShaRinGan/coverage.$binName.cov >> run.log 2>&1 &
#mac下启动(务必 绝对路径 启动):
nohup /xx/$binName.test -test.coverprofile=/tmp/ShaRinGan/coverage.$binName.cov >> run.log 2>&1 &
```
> SUT一键接入和启动[脚本](../../example/replayer/sut_replayer.sh) 及其 [使用方法](./replayer-sut.md)

至此，浏览器打开 [http://127.0.0.1:8998](http://127.0.0.1:8998) 或 local_ip 即可开始回放啦~

<br>

#### 二、使用

##### 1. 覆盖率报告

覆盖率报告支持 多次回放后的 **累计覆盖率统计** ；即 在SUT启动后，可以进行多次回放测试，最后统一查看覆盖率报告。

覆盖率报告的查看点有两处：

A. 单流量回放结果页

![codeCover_report_singal](http://img-hxy021.didistatic.com/static/sharingan/codeCover_report_singal_v2.png)

B. [批量回放](./replayer-parallel.md)结果页

![codeCover_report_parallel](http://img-hxy021.didistatic.com/static/sharingan/codeCover_report_parallel_v2.png)

> 注意：
> 
> a.点击 "覆盖率报告" 链接后，Replayer-Agent会 **自动重启SUT服务**。所以，覆盖率数据会重新统计!!!
>
> b.上面两处 "覆盖率报告" 链接 点击效果一样。覆盖率统计不区分单流量回放和批量回放，会一起统计。


覆盖率报告 不仅会给出 整体覆盖结果，还会给出 每个函数的覆盖率结果和覆盖详情。

![codeCover_report_sum](http://img-hxy021.didistatic.com/static/sharingan/codeCover_report_sum_v2.png)

下面是每个函数具体覆盖详情，其中 红色为未覆盖代码。

![codeCover_report_detail](http://img-hxy021.didistatic.com/static/sharingan/codeCover_report_detail_v2.png)

<br>

##### 2. 支持SUT使用flag

对于使用flag的SUT,需要保证SUT在启动后10s内完成flag.Parse()

若10s仍然不能完成，请首先考虑SUT的代码是否合理。确认必须后，可修改replayer-agent/install/codeCov/main_test.go中waitFlagParseTime常量的值

其他接入操作不变，同 [配置并启动SUT](#2-配置并启动sut)。

<br>

##### 3. 历史覆盖率报告

历史覆盖率报告存放在本机 /tmp/ShaRinGan/ 目录下。

a. 覆盖率原始报告, 命名格式 coverage.$binName.$timestamp
> $timestamp 一般包括起始和终止两个时间戳。对于初始手动启动SUT的情况，只有终止时间戳。对于Replayer-Agent自动启动SUT的情况，都有起始和终止两个时间戳。

b. 格式化后的 *.html 报告，命名格式 coverage.$binName.$timestamp.html
> 意外惊喜：历史 *.html 报告，可通过接口 http://127.0.0.1:8998/coverage/report/coverage.$binName.$timestamp.html 查看。

由于历史报告保存在/tmp临时目录下，所以对于想长期保存的历史报告，请及时备份。

