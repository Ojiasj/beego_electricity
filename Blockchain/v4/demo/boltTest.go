package main

import (
	"fmt"
	"go一期/bolt"
	"log"
)

func main() {
	//1. 打开数据库
	db, err := bolt.Open("test.db", 0600, nil)
	defer db.Close()

	if err != nil {
		log.Panic("打开数据库失败！")
	}

	//将要操作数据库（改写）
	db.Update(func(tx *bolt.Tx) error {
		//2. 找到抽屉bucket(如果没有，就创建）
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			//没有抽屉，我们需要创建
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				log.Panic("创建bucket(b1)失败")
			}
		}

		//3. 写数据
		bucket.Put([]byte("11111"), []byte("hello"))
		bucket.Put([]byte("22222"), []byte("world"))

		return nil
	})

	//4. 读数据
	db.View(func(tx *bolt.Tx) error {
		//1. 找到抽屉，没有的化直接报错退出
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			log.Panic("bucket b1 不应该为空，请检查!!!!")
		}

		//2. 直接读取数据
		v1 := bucket.Get([]byte("11111"))
		v2 := bucket.Get([]byte("22222"))

		fmt.Printf("v1 : %s\n", v1)
		fmt.Printf("v2 : %s\n", v2)

		return nil
	})

}
