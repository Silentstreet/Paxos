<!--
 * @Author: your name
 * @Date: 2021-05-01 11:11:29
 * @LastEditTime: 2021-05-01 12:12:18
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /Codes/paxoskv.md
-->
#### paxos的解释
1. paxos的工作：就是把一堆运行的机器协同起来，让多个机器成为一个整体系统。
    在这个系统中，每个机器都必须让系统中的状态保持一致。
        例如：3副本集群如果一个机器上传了一张图片，那么另外2台机器也必须复制这张图片过来，整个系统处于一个一致的状态

##### 技术演变：
早些年有各种各样的复制策略都被提出来解决各种场景下的需要，除了复制的份数之外，各种各样的算法实际上都是在尝试解决一致问题，
每种复制策略都有其优缺点， 最后引出paxos如何解决副本一致性的问题
！[image]()