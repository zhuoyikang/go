package bomb

/*
  炸弹爆炸是个单独的事件，不同的事件由不同的数据包表示
*/

type EvtBombExplode struct {
	x int
	y int
	r int  //范围十字.
}
