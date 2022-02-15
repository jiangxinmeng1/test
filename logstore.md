# Improve logstore

## 概要设计

1. 增加分区层
   
   在vinfo和syncbase中，记录commt和checkpoint信息的同时记录group名。

   ```golang
   type vInfo struct {
	    commits     map [string]common.ClosedInterval//[groupName]commits
	    checkpoints map [string][]common.ClosedInterval//[groupName]checkpoint
   }
   type syncBase struct {
	    synced, syncing             map [string]uint64//[groupName]commit id
	    checkpointed, checkpointing map[string]uint64//[groupName]checkpoint id
   }   
   ```

   在entry的info中记录group name。

2. 支持分区级别的checkpoint
   
   checkpoint中logstore只用提供group对应的commit id，其他的都在外部生成。
   
   syncbase提供接口，返回全局或指定组的元数据。

3. 获取分区和全局的元数据
   
   全局元数据存在syncBase中。

   分区元数据存在vinfo中，持久化在vfile的file末尾，在vfile.Commit()时持久化。

   全局和分区的元数据都在打开rotatefile的时候恢复到内存中。

   vinfo提供接口，返回当前vfile的元数据。供replay handler使用，以及syncbase同步元数据时读取。

   syncbase提供接口，返回logstore的元数据，在生成checkpoint entry的时候使用。

   对外，元数据从syncbase中取，例如：

   ```golang
   func (base *syncBase) GetPenddings(groupName string) uint64 {
   	ckp := base.GetCheckpointed(groupName)
   	commit := base.GetSynced(groupName)
   	return commit - ckp
   }
   ```


## 任务拆解

1. 更改vinfo和syncbase格式，和相关函数（e.g. syncBase.SetCheckpointed）

2. vinfo提供接口，从文件中读取元数据

3. rotatefile打开时恢复history和全局元数据（syncbase）
