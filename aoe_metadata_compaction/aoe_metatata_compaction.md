# AOE Metadata Compaction

## 需求背景

当前的元数据日志没有压缩。压缩后，可以减少日志的存储空间，并且节省replay的时间。

## 概要设计

1. checkpoint

   1. 触发时机
   
        通过catalog.checkpointer触发。每过一分钟检查当前的commit id和上次checkpoint id的间隔。如果间隔大于DefaultCheckpointDelta，则触发checkpoint。

   2. 生成checkpoint entry

        * logentrys

        database, table, segment, block

        记录的范围是上次checkpoint后发生变化的database, table, segment, block 。(i.e. commitId > previousCheckpointId)

        记录的内容与对应的logEntry相似（e.g. 对 table 会记录类似tableLogEntry的内容），会记录基本信息和commit信息。只会记录最新的commit信息。
   
        * databaseset & tableset
   
        会额外记录所有database的名字，和每个database中所有table的名字。

        用来检查上次checkpoint与本次之间被删除的database和table。

2. replay
   
  遍历catalog, databaseSet和tableSet
   
  删除databaseSet和tableSet中不存在的database和table。

  更新logEntry中提到的database, table, segment, block, 类似onReplay函数。

## 任务拆解

1. 添加logstore, metadata包下的注释

   checkpoint和replay过程相关的代码 
   
   3day

2. 自测 2day
