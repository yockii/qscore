package util

import snowflake "github.com/yockii/snowflake_ext"

var snowflakeWorker *snowflake.Worker

func InitNode(workerId uint64) (err error) {
	w, err := snowflake.NewSnowflake(workerId)
	if err != nil {
		return err
	}
	snowflakeWorker = w
	return nil
}

func SnowflakeId() uint64 {
	return snowflakeWorker.NextId()
}

//
//var snowflakeNode *snowflake.Node
//
//func InitNode(node int64) (err error) {
//	snowflakeNode, err = snowflake.NewNode(1)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func SnowflakeId() int64 {
//	return snowflakeNode.Generate().Int64()
//}
