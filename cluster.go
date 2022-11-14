package main

type MsgBody struct {
	blockA M
	blockB M
	step   int // which iteration
}
type MsgChan chan MsgBody

type Channles struct {
	msgChanMat    [][]MsgChan
	resultChanMat [][]chan M
	dimByBlock    int
}

func CreateChannels(dimByBlock int) Channles {
	channels := Channles{}
	channels.dimByBlock = dimByBlock
	channels.msgChanMat = make([][]MsgChan, dimByBlock)
	channels.resultChanMat = make([][]chan M, dimByBlock)
	for i := 0; i < dimByBlock; i++ {
		channels.msgChanMat[i] = make([]MsgChan, dimByBlock)
		channels.resultChanMat[i] = make([]chan M, dimByBlock)
		for j := 0; j < dimByBlock; j++ {
			channels.msgChanMat[i][j] = make(MsgChan)
			channels.resultChanMat[i][j] = make(chan M)
		}
	}
	return channels
}

type Cluster struct {
	procMat    [][]Proc
	dimByBlock int
}

func CreatCluster(
	dimByBlock int, dimOfBlock int, bmA BlockMat, bmB BlockMat, channels Channles) Cluster {

	cluster := Cluster{}
	cluster.dimByBlock = dimByBlock
	cluster.procMat = make([][]Proc, dimByBlock)
	for i := 0; i < dimByBlock; i++ {
		cluster.procMat[i] = make([]Proc, dimByBlock)
		for j := 0; j < dimByBlock; j++ {
			cluster.procMat[i][j] = CreateProcessor(
				i, j, dimOfBlock, channels.msgChanMat[i][j], channels.resultChanMat[i][j])
		}
	}
	return cluster
}

func (o *Cluster) joinBlocks(dimOfBlock int, mC *M) {
	for ir := 0; ir < o.dimByBlock; ir++ {
		for ic := 0; ic < o.dimByBlock; ic++ {
			blockC := <-o.procMat[ir][ic].resultChan
			joinMat(*mC, blockC,
				ic*dimOfBlock, ic*dimOfBlock+dimOfBlock,
				ir*dimOfBlock, ir*dimOfBlock+dimOfBlock)
		}
	}
}

func (o *Cluster) goProcs() {
	for ir := 0; ir < o.dimByBlock; ir++ {
		for ic := 0; ic < o.dimByBlock; ic++ {
			proc := o.procMat[ir][ic]
			go proc.multiply()
		}
	}
}

func (o *Cluster) closeProcs() {
	for ir := 0; ir < o.dimByBlock; ir++ {
		for ic := 0; ic < o.dimByBlock; ic++ {
			close(o.procMat[ir][ic].msgChan)
		}
	}
}

func (o *Cluster) broadcastNextBlocksToProcs(bmA BlockMat, bmB BlockMat, step int) {
	for ir := 0; ir < o.dimByBlock; ir++ {
		for ic := 0; ic < o.dimByBlock; ic++ {
			//fmt.Println("brd: ", ir, ic, step, " mat ", bmA.dim, bmA.dimOfBlock, bmB.dim, bmB.dimOfBlock)
			msg := MsgBody{bmA.bm[ir][step], bmB.bm[step][ic], step}
			o.procMat[ir][ic].msgChan <- msg
		}
	}
}

type Proc struct {
	ir, ic     int
	blockA     M
	blockB     M
	blockC     M
	msgChan    MsgChan
	resultChan chan M
}

func CreateProcessor(
	ir int, ic int, dimOfBlock int,
	msgChan MsgChan, resultChan chan M) Proc {

	proc := Proc{}
	proc.ir = ir
	proc.ic = ic
	proc.blockC = CreateZeroMatrix(dimOfBlock, dimOfBlock)
	proc.msgChan = msgChan
	proc.resultChan = resultChan
	return proc
}

func (o *Proc) multiply() {
	for {
		msgBody, isOpen := <-o.msgChan
		if isOpen {
			o.blockA = msgBody.blockA
			o.blockB = msgBody.blockB
			o.blockC = addMat(o.blockC, MultiplyStandardParallel(o.blockA, o.blockB))
		} else {
			o.resultChan <- o.blockC
			return
		}
	}
}
