/*
很多资源需要管理
1.在线玩家{id => 数据}
2.网络连接{id => 网络连接}

此功能可大量复用，提供以下操作

1.通过Key查询data.
2.增加{k, v}
3.删除k
4.换新key{k1,v}  => {k2,v}
5.遍历所有的k,v 有锁或者无锁
*/

package agent

import (
	//"fmt"
)

//
func test() {
}
