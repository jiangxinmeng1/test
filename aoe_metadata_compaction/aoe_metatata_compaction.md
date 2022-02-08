# AOE Metadata Compaction

## 需求背景

当前的元数据日志没有压缩。压缩后，可以减少日志的存储空间，并且节省replay的时间。

## 概要设计

1. checkpoint

   1. 触发时机
   
        通过catalog.checkpointer触发。每间隔DefaultCheckpointInterval时间检查当前的commit id和上次checkpoint id之间的commit数量。如果commit数量大于DefaultCheckpointDelta，则触发checkpoint。
        ```golang
        //用HeartBeater定时检查
          catalog.checkpointer = worker.NewHeartBeater(DefaultCheckpointInterval, &catalogCheckpointer{
	          	catalog:   catalog,
	          })

          func (c *catalogCheckpointer)OnExec(){
          	previousCheckpointId := c.catalog.GetCheckpointId()
          	commitId := c.catalog.Store.GetSyncedId()
               //检查commit数量
          	if commitId < previousCheckpointId+DefaultCheckpointDelta{
          		return
          	}
          	c.catalog.Checkpoint()
          }
          
        ```

   2. 生成checkpoint entry

      记录的范围是上次checkpoint后发生变化的database, table, segment, block和所有database的safeid。只会记录最新的commit信息。

      会额外记录所有database的名字，和每个database中所有table的名字。用来检查上次checkpoint与本次之间被删除的database和table。
    ```golang
    
     
    type segmentCheckpoint struct {
	      Blocks     []*blockLogEntry
	     NeedReplay bool
	     LogEntry   segmentLogEntry
     }

     type tableCheckpoint struct {
	     Segments   []*segmentCheckpoint
	     NeedReplay bool
	     LogEntry   tableLogEntry
     }

     type databaseCheckpoint struct {
	     Tables     map[string]*tableCheckpoint
	     NeedReplay bool
	     LogEntry   databaseLogEntry
     }

     type catalogLogEntry struct {
	     Databases map[string]*databaseCheckpoint
         SafeId    map[uint64]uint64
	     Range     *common.Range
     }

    ```

2. replay
   
  遍历catalog, catalogLogEntry。
   
  删除databaseSet和tableSet中不存在的database和table。

  更新logEntry中提到的database, table, segment, block, 类似onReplay函数。

  覆盖恢复wal里每个shard的safeId。

## 任务拆解

1. 添加logstore, metadata包下的注释

   checkpoint和replay过程相关的代码 
   
   2day

2. 自测 2day
