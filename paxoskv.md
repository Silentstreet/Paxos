<!--
 * @Author: your name
 * @Date: 2021-05-01 11:11:29
 * @LastEditTime: 2021-05-01 11:57:29
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /Codes/paxoskv.md
-->
#### paxos的解释
1. paxos的工作：就是把一堆运行的机器协同起来，让多个机器成为一个整体系统。
    在这个系统中，每个机器都必须让系统中的状态保持一致。
        例如：3副本集群如果一个机器上传了一张图片，那么另外2台机器也必须复制这张图片过来，整个系统处于一个一致的状态