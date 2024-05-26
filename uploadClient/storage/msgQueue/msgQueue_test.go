package msgQueue

import (
	"context"
	"fmt"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const testNum = 1000
const msgGroupID = "0"

func TestInit(t *testing.T) {
	queue := MessageQueue{}
	md, _ := clientModel.NewPreUploadFileData(commonModel.LocalScraperType, "1", msgGroupID, "test.jpg")
	msgData := &MsgQueueData{
		MsgMetaData: MsgMetaData{
			UploadType: commonModel.PreUploadType,
			MsgMetaID: MsgMetaID{
				ScraperID:   "1",
				MsgGroupID:  strconv.Itoa(1),
				ScraperType: commonModel.APIScraperType,
			},
		},
		FileMetaData: &clientModel.AnyFileMetaDataModel{
			PreUploadFileMetaDataModel: &md.PreUploadFileMetaDataModel,
		},
	}
	if err := msgData.Commit(); err == nil {
		t.Errorf("Commit() should fail before initialization")
	}
	if err := msgData.GoDead(); err == nil {
		t.Errorf("GoDead() should fail before initialization")
	}
	if err := queue.AddElement(msgData); err == nil {
		t.Errorf("AddElement() should fail before initialization")
	}
	if err := queue.AddElements([]*MsgQueueData{msgData}); err == nil {
		t.Errorf("AddElements() should fail before initialization")
	}
	if _, err := queue.ListenUploadType(context.Background(), commonModel.PreUploadType); err == nil {
		t.Errorf("ListenScpID() should fail before initialization")
	}
	if _, err := queue.ListenMsgMetaData(msgData.MsgMetaData); err == nil {
		t.Errorf("ListenMsgID() should fail before initialization")
	}
	if _, err := queue.PopData(msgData.MsgMetaData); err == nil {
		t.Errorf("PopID() should fail before initialization")
	}
	if _, err := queue.PopAll(ActivateQueue); err == nil {
		t.Errorf("PopAll() should fail before initialization")
	}
}

func TestAddElementAndPop(t *testing.T) {
	queue := NewMessageQueue()
	md, _ := clientModel.NewPreUploadFileData(commonModel.LocalScraperType, "1", msgGroupID, "test.jpg")
	msgData := &MsgQueueData{
		MsgMetaData: MsgMetaData{
			UploadType: commonModel.PostUploadType,
			MsgMetaID: MsgMetaID{
				ScraperType: commonModel.LocalScraperType,
				ScraperID:   "1",
				MsgGroupID:  strconv.Itoa(1),
			},
		},
		FileMetaData: &clientModel.AnyFileMetaDataModel{
			PreUploadFileMetaDataModel: &md.PreUploadFileMetaDataModel,
		},
	}
	if err := queue.AddElement(msgData); err != nil {
		t.Fatalf("AddElement() failed: %v", err)
	}
	if val, ok := queue.activateQueue.Load(msgData.MsgMetaData); !ok {
		t.Errorf("Element was not added to the activateQueue")
	} else {
		set := val.(*sync.Map)
		if _, exist := set.Load(getMsgPureData(msgData)); !exist {
			t.Errorf("Added element is not in the set")
		}
	}
	poppedData, err := queue.PopData(msgData.MsgMetaData)
	if err != nil {
		t.Fatalf("PopID() failed: %v", err)
	}
	if poppedData[0] != msgData {
		t.Errorf("Popped data does not match the sent data")
	} else {
		t.Logf("Popped data matches the sent data")
	}
}

func TestAddElementsAndPopAll(t *testing.T) {
	queue := NewMessageQueue()
	msgQueueDataSlice := make([]*MsgQueueData, 0)
	// 3-1003
	scpStartNo := 3
	// Same MsgMetaData + Same FileMetaData
	for i := scpStartNo; i < testNum/2+scpStartNo; i++ {
		md, _ := clientModel.NewPreUploadFileData(commonModel.LocalScraperType, "1", msgGroupID, "test.jpg")
		msgQueueDataSlice = append(msgQueueDataSlice, &MsgQueueData{
			MsgMetaData: MsgMetaData{
				UploadType: commonModel.UploadType(strconv.Itoa(scpStartNo)),
				MsgMetaID: MsgMetaID{
					ScraperID:   "1",
					MsgGroupID:  strconv.Itoa(1),
					ScraperType: commonModel.LocalScraperType,
				},
			},
			FileMetaData: &clientModel.AnyFileMetaDataModel{
				PreUploadFileMetaDataModel: &md.PreUploadFileMetaDataModel,
			},
		})
	}
	// Same MsgMetaData + Different FileMetaData
	for i := scpStartNo + testNum/2; i < scpStartNo+testNum; i++ {
		md, _ := clientModel.NewPreUploadFileData(
			commonModel.LocalScraperType,
			"1",
			msgGroupID,
			fmt.Sprintf("test%d.jpg", i),
		)
		msgQueueDataSlice = append(msgQueueDataSlice, &MsgQueueData{
			MsgMetaData: MsgMetaData{
				UploadType: commonModel.UploadType(strconv.Itoa(scpStartNo)),
				MsgMetaID: MsgMetaID{
					ScraperID:   "1",
					MsgGroupID:  strconv.Itoa(1),
					ScraperType: commonModel.LocalScraperType,
				},
			},
			FileMetaData: &clientModel.AnyFileMetaDataModel{
				PreUploadFileMetaDataModel: &md.PreUploadFileMetaDataModel,
				UploadFileMetaDataModel:    &clientModel.UploadFileMetaDataModel{},
			},
		})
	}
	err := queue.AddElements(msgQueueDataSlice)
	if err != nil {
		t.Errorf("AddElements failed: %v", err)
	}
	for _, data := range msgQueueDataSlice {
		if val, ok := queue.activateQueue.Load(data.MsgMetaData); !ok {
			t.Errorf("Element with ScraperID %s was not added to the activateQueue", data.MsgMetaData.ScraperID)
		} else {
			set := val.(*sync.Map)
			_, exist := set.Load(getMsgPureData(data))
			if !exist {
				t.Errorf("Added element with ScraperID %s is not in the set", data.MsgMetaData.ScraperID)
			}
		}
	}
	poppedDataSlice, err := queue.PopAll(ActivateQueue)
	if err != nil {
		t.Fatalf("PopAll() failed: %v", err)
	}
	// len(poppedDataSlice) must be testNum/2+1 cuz they have same MsgMetaData(1) and different MsgMetaData(testNum/2)
	if len(poppedDataSlice) != testNum/2+1 {
		t.Errorf("Popped data slice length does not match the sent data slice length")
	} else {
		t.Logf("Popped data slice length matches the sent data slice length")
	}
	foundTime := 0
	for _, poppedData := range poppedDataSlice {
		for _, originalData := range msgQueueDataSlice {
			if originalData.MsgMetaData == poppedData.MsgMetaData {
				foundTime++
				break
			}
		}
	}
	if foundTime != testNum/2+1 {
		t.Errorf("Popped data slice does not match the sent data slice")
	} else {
		t.Logf("Popped data slice matches the sent data slice")
	}
	poppedDataSlice2, err := queue.PopAll(DeadQueue)
	if err != nil {
		t.Fatalf("PopAll() failed: %v", err)
	}
	if len(poppedDataSlice2) != 0 {
		t.Errorf("Popped data slice length does not match the sent data slice length")
	} else {
		t.Logf("Popped data slice length matches the sent data slice length")
	}
	if _, _err := queue.PopAll(msgQueueType(114514)); _err == nil {
		t.Errorf("PopAll() should fail with invalid msgQueue type")
	}
}

func TestListenChan(t *testing.T) {
	queue := NewMessageQueue()
	scpStartNo := 2000
	wg := sync.WaitGroup{}
	// element
	md, _ := clientModel.NewPreUploadFileData(
		commonModel.LocalScraperType,
		strconv.Itoa(scpStartNo),
		msgGroupID,
		"test.jpg",
	)
	msgData1 := &MsgQueueData{
		MsgMetaData: MsgMetaData{
			UploadType: commonModel.UploadType(strconv.Itoa(scpStartNo)),
			MsgMetaID: MsgMetaID{
				ScraperID:   "1",
				MsgGroupID:  strconv.Itoa(1),
				ScraperType: commonModel.LocalScraperType,
			},
		},
		FileMetaData: &clientModel.AnyFileMetaDataModel{
			PreUploadFileMetaDataModel: &md.PreUploadFileMetaDataModel,
		},
	}
	msgData2 := &MsgQueueData{
		MsgMetaData: MsgMetaData{
			UploadType: commonModel.UploadType(strconv.Itoa(scpStartNo)),
			MsgMetaID: MsgMetaID{
				ScraperID:   "3",
				MsgGroupID:  strconv.Itoa(2),
				ScraperType: commonModel.APIScraperType,
			},
		},
		FileMetaData: &clientModel.AnyFileMetaDataModel{
			PreUploadFileMetaDataModel: &md.PreUploadFileMetaDataModel,
		},
	}
	// NOTICE: all element in msg1Slice || msg2Slice have same MsgMetaData,
	// which means there are only one element after storage in set
	msg1Slice := make([]*MsgQueueData, 0)
	msg2Slice := make([]*MsgQueueData, 0)
	for i := 0; i < testNum; i++ {
		msg1Slice = append(msg1Slice, msgData1)
		msg2Slice = append(msg2Slice, msgData2)
	}
	if err := queue.AddElements(msg1Slice); err != nil {
		t.Fatalf("AddElement() failed after ListenScpID: %v", err)
	}
	if err := queue.AddElements(msg2Slice); err != nil {
		t.Fatalf("AddElement() failed after ListenScpID: %v", err)
	}
	wg.Add(3)
	// global scpID chan
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		var counter int64
		defer wg.Done()
		defer cancel()
		ch, err := queue.ListenUploadType(ctx, commonModel.UploadType(strconv.Itoa(scpStartNo)))
		if err != nil {
			t.Errorf("ListenScpID() failed: %v", err)
		}
		for {
			select {
			case receivedMsgData := <-ch:
				atomic.AddInt64(&counter, 1)
				if receivedMsgData != msgData1 && receivedMsgData != msgData2 {
					t.Errorf("Global scpID chan Received data does not match the sent data")
				} else {
					t.Logf("Global scpID chan Received data matches the sent data")
				}
			// BUG:
			case <-time.After(2 * time.Second):
				logrus.Debug("Received ", atomic.LoadInt64(&counter), " messages")
				if atomic.LoadInt64(&counter) != 2 {
					t.Errorf("[0] Expected %d message, but received %d", 2, atomic.LoadInt64(&counter))
				}
				return
			}
		}
	}()
	// another only MsgQueueID chan
	// TODO: timeout
	go func() {
		defer wg.Done()
		var counter int64
		ch1, err := queue.ListenMsgMetaData(msgData1.MsgMetaData)
		if err != nil {
			t.Errorf("ListenScpID() failed: %v", err)
			return
		}
		for ele := range ch1 {
			if ele != msgData1 {
				t.Errorf("Only MsgQueueID chan Received data does not match the sent data")
			} else {
				atomic.AddInt64(&counter, 1)
				t.Logf("Only MsgQueueID chan Received data matches the sent data")
			}
		}
	}()
	go func() {
		defer wg.Done()
		var counter int64
		ch2, err := queue.ListenMsgMetaData(msgData2.MsgMetaData)
		if err != nil {
			t.Errorf("ListenScpID() failed: %v", err)
			return
		}
		for ele := range ch2 {
			if ele != msgData2 {
				t.Errorf("Only MsgQueueID chan Received data does not match the sent data")
			} else {
				atomic.AddInt64(&counter, 1)
				t.Logf("Only MsgQueueID chan Received data matches the sent data")
			}
		}
	}()
	wg.Wait()
	// check repeat chan
	chan1, err := queue.ListenUploadType(context.Background(), "114514")
	if err != nil {
		t.Errorf("ListenMsgID() failed: %v", err)
	}
	chan2, err := queue.ListenUploadType(context.Background(), "114514")
	if err != nil {
		t.Errorf("ListenMsgID() failed: %v", err)
	}
	if chan1 != chan2 {
		t.Errorf("chan1 and chan2 should be the same")
	}
}

// TestCommitAndGoDead TODO: test chan listen after eme.Commit() and eme.GoDead()
func TestCommitAndGoDead(t *testing.T) {
	queue := NewMessageQueue()
	md, _ := clientModel.NewPreUploadFileData(commonModel.LocalScraperType, "1", msgGroupID, "test.jpg")
	scpStartNo1 := 3000
	scpStartNo2 := 4000
	msgData1 := &MsgQueueData{
		MsgMetaData: MsgMetaData{
			UploadType: commonModel.UploadType(strconv.Itoa(scpStartNo1)),
			MsgMetaID: MsgMetaID{
				ScraperID:   "1",
				MsgGroupID:  strconv.Itoa(2),
				ScraperType: commonModel.LocalScraperType,
			},
		},
		FileMetaData: &clientModel.AnyFileMetaDataModel{
			PreUploadFileMetaDataModel: &md.PreUploadFileMetaDataModel,
		},
	}
	msgData2 := &MsgQueueData{
		MsgMetaData: MsgMetaData{
			UploadType: commonModel.UploadType(strconv.Itoa(scpStartNo2)),
			MsgMetaID: MsgMetaID{
				ScraperID:   "3",
				MsgGroupID:  strconv.Itoa(4),
				ScraperType: commonModel.LocalScraperType,
			},
		},
		FileMetaData: &clientModel.AnyFileMetaDataModel{
			PreUploadFileMetaDataModel: &md.PreUploadFileMetaDataModel,
		},
	}
	if err := queue.AddElement(msgData1); err != nil {
		t.Fatalf("AddElement() failed: %v", err)
	}
	err := msgData1.Commit()
	if err != nil {
		t.Fatalf("Commit() failed: %v", err)
	}
	if ele, ok := queue.activateQueue.Load(msgData1.MsgMetaData); ok {
		if _, _ok := ele.(*sync.Map).Load(msgData1.MsgMetaData); _ok {
			t.Errorf("MsgQueueData was not removed from the activateQueue after Commit()")
		} else {
			t.Logf("MsgQueueData was removed from the activateQueue after Commit()")
		}
	} else {
		t.Errorf("MsgQueueData was not removed from the activateQueue after Commit()")
	}
	// Test GoDead
	if err = queue.AddElement(msgData2); err != nil {
		t.Fatalf("AddElement() failed: %v", err)
	}
	err = msgData2.GoDead()
	if err != nil {
		t.Fatalf("GoDead() failed: %v", err)
	}
	if ele, ok := queue.activateQueue.Load(msgData2.MsgMetaData); ok {
		if _, _ok := ele.(*sync.Map).Load(msgData2.MsgMetaData); _ok {
			t.Errorf("MsgQueueData was not removed from the activateQueue after Commit()")
		} else {
			t.Logf("MsgQueueData was removed from the activateQueue after Commit()")
		}
	} else {
		t.Errorf("MsgQueueData was not removed from the activateQueue after Commit()")
	}
	if ele, ok := queue.deadQueue.Load(msgData2.MsgMetaData); ok {
		if _, _ok := ele.(*sync.Map).Load(getMsgPureData(msgData2)); _ok {
			t.Logf("MsgQueueData was added to the DeadQueue after GoDead()")
		} else {
			t.Errorf("MsgQueueData was not added to the DeadQueue after GoDead()")
		}
	} else {
		t.Errorf("MsgQueueData was not removed from the activateQueue after Commit()")
	}
}

func TestConcurrentAddAndRemove(t *testing.T) {
	queue := NewMessageQueue()
	wg := sync.WaitGroup{}
	scpStartNo := 5000
	numGoroutines := 50
	numElementsPerRoutine := 100
	dataAdded := make(chan *MsgQueueData, numGoroutines*numElementsPerRoutine)
	for i := scpStartNo; i < numGoroutines+scpStartNo; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numElementsPerRoutine; j++ {
				msgData := &MsgQueueData{
					MsgMetaData: MsgMetaData{
						UploadType: commonModel.UploadType(strconv.Itoa(id)),
						MsgMetaID: MsgMetaID{
							ScraperID:   strconv.Itoa(j),
							MsgGroupID:  strconv.Itoa(j + 1),
							ScraperType: commonModel.LocalScraperType,
						},
					},
					FileMetaData: &clientModel.AnyFileMetaDataModel{
						PreUploadFileMetaDataModel: &clientModel.PreUploadFileMetaDataModel{},
					},
				}
				if err := queue.AddElement(msgData); err != nil {
					t.Errorf("AddElement() failed: %v", err)
				} else {
					dataAdded <- msgData
				}
				if poppedData, err := queue.PopData(msgData.MsgMetaData); err != nil {
					t.Errorf("PopID() failed: %v", err)
				} else if poppedData[0] != msgData {
					t.Errorf("Popped data does not match added data: added %v, popped %v", msgData, poppedData)
				}
			}
		}(i)
	}
	wg.Wait()
	close(dataAdded)
	addedDataMap := make(map[MsgMetaData]*MsgQueueData)
	for data := range dataAdded {
		addedDataMap[data.MsgMetaData] = data
	}
	if len(addedDataMap) != numGoroutines*numElementsPerRoutine {
		t.Errorf("Mismatch in the number of added and processed data: added %d, processed %d",
			numGoroutines*numElementsPerRoutine, len(addedDataMap))
	}
}

// TestConcurrentListening1 test different UploadType & one listenChan
func TestConcurrentListening1(t *testing.T) {
	queue := NewMessageQueue()
	wg := sync.WaitGroup{}
	numListeners := 200
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	scpStartNo := 6000
	wg.Add(numListeners)
	var all int32
	go func() {
		ch, err := queue.ListenUploadType(ctx, commonModel.UploadType(strconv.Itoa(scpStartNo)))
		if err != nil {
			t.Errorf("ListenUploadType() failed for UploadType %d: %v", scpStartNo, err)
			return
		}
		for itm := range ch {
			if itm.MsgMetaData.UploadType != commonModel.UploadType(strconv.Itoa(scpStartNo)) {
				t.Errorf("Received wrong UploadType data: got %s, want %s",
					itm.MsgMetaData.UploadType, commonModel.UploadType(strconv.Itoa(scpStartNo)))
			} else {
				t.Logf("Received correct UploadType %s", commonModel.UploadType(strconv.Itoa(scpStartNo)))
				atomic.AddInt32(&all, 1)
			}
		}
	}()
	for i := scpStartNo; i < numListeners+scpStartNo; i++ {
		msgData := &MsgQueueData{
			MsgMetaData: MsgMetaData{
				UploadType: commonModel.UploadType(strconv.Itoa(scpStartNo)),
				MsgMetaID: MsgMetaID{
					ScraperID:   "1",
					MsgGroupID:  strconv.Itoa(1),
					ScraperType: commonModel.LocalScraperType,
				},
			},
		}
		if err := queue.AddElement(msgData); err != nil {
			t.Fatalf("AddElement() failed: %v", err)
		}
		logrus.Debugf("Added element with ScraperType %d", i)
		wg.Done()
	}
	wg.Wait()
	if atomic.LoadInt32(&all) != 1 {
		t.Errorf("Mismatch in the number of listeners and received data: listeners %d, received %d",
			numListeners, atomic.LoadInt32(&all))
	} else {
		t.Logf("All listeners received data!")
	}
}

// TestConcurrentListening2 test same UploadType & one listenChan
func TestConcurrentListening2(t *testing.T) {
	queue := NewMessageQueue()
	numListeners := 200
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	scpStartNo := 7000
	var all int32
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch, err := queue.ListenUploadType(ctx, commonModel.PreUploadType)
		if err != nil {
			t.Errorf("ListenUploadType() failed for UploadType %d: %v", scpStartNo, err)
			return
		}
		for itm := range ch {
			if itm.MsgMetaData.UploadType != commonModel.PreUploadType {
				t.Errorf("Received wrong UploadType data: got %s, want %s",
					itm.MsgMetaData.UploadType, commonModel.PreUploadType)
			} else {
				t.Logf("Received correct UploadType %s", commonModel.PreUploadType)
				atomic.AddInt32(&all, 1)
			}
			if atomic.LoadInt32(&all) == int32(numListeners) {
				return
			}
		}
	}()
	for i := scpStartNo; i < numListeners+scpStartNo; i++ {
		msgData := &MsgQueueData{
			MsgMetaData: MsgMetaData{
				UploadType: commonModel.PreUploadType,
				MsgMetaID: MsgMetaID{
					ScraperID:   strconv.Itoa(rand.Intn(100) + 1),
					MsgGroupID:  strconv.Itoa(rand.Intn(100) + 1),
					ScraperType: commonModel.ScraperType(strconv.Itoa(i)),
				},
			},
			FileMetaData: &clientModel.AnyFileMetaDataModel{
				PreUploadFileMetaDataModel: &clientModel.PreUploadFileMetaDataModel{
					ScraperType: commonModel.ScraperType(strconv.Itoa(i)),
					ScraperID:   strconv.Itoa(rand.Intn(100) + 1),
					ResourceUri: fmt.Sprintf("test%d.jpg", i),
				},
			},
		}
		if err := queue.AddElement(msgData); err != nil {
			t.Fatalf("AddElement() failed: %v", err)
		}
		logrus.Debugf("Added element with ScraperType %d", i)
	}
	wg.Wait()
	if atomic.LoadInt32(&all) != int32(numListeners) {
		t.Errorf("Mismatch in the number of listeners and received data: listeners %d, received %d",
			numListeners, atomic.LoadInt32(&all))
	} else {
		t.Logf("All listeners received data!")
	}
}

// TestConcurrentListening3 test different UploadType & one listenChan
func TestConcurrentListening3(t *testing.T) {
	queue := NewMessageQueue()
	wg := sync.WaitGroup{}
	numListeners := 200
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	scpStartNo := 8000
	wg.Add(1)
	var all int32
	go func() {
		defer wg.Done()
		ch, err := queue.ListenUploadType(ctx, commonModel.UploadType(strconv.Itoa(scpStartNo)))
		if err != nil {
			t.Errorf("ListenUploadType() failed for UploadType %d: %v", scpStartNo, err)
			return
		}
		for itm := range ch {
			if itm.MsgMetaData.UploadType != commonModel.UploadType(strconv.Itoa(scpStartNo)) {
				t.Errorf("Received wrong UploadType data: got %s, want %s",
					itm.MsgMetaData.UploadType, commonModel.UploadType(strconv.Itoa(scpStartNo)))
			} else {
				t.Logf("Received correct UploadType %s", commonModel.UploadType(strconv.Itoa(scpStartNo)))
				atomic.AddInt32(&all, 1)
			}
			if atomic.LoadInt32(&all) == int32(numListeners) {
				return
			}
		}
	}()
	for i := scpStartNo; i < numListeners+scpStartNo; i++ {
		msgData := &MsgQueueData{
			MsgMetaData: MsgMetaData{
				UploadType: commonModel.UploadType(strconv.Itoa(scpStartNo)),
				MsgMetaID: MsgMetaID{
					ScraperID:   strconv.Itoa(1),
					MsgGroupID:  strconv.Itoa(1),
					ScraperType: commonModel.LocalScraperType,
				},
			},
			FileMetaData: &clientModel.AnyFileMetaDataModel{
				PreUploadFileMetaDataModel: &clientModel.PreUploadFileMetaDataModel{
					ScraperType: commonModel.LocalScraperType,
					ScraperID:   "1",
					ResourceUri: fmt.Sprintf("test%d.jpg", i),
				},
			},
		}
		if err := queue.AddElement(msgData); err != nil {
			t.Fatalf("AddElement() failed: %v", err)
		}
		logrus.Debugf("Added element with ScraperType %d", i)
	}
	wg.Wait()
	if atomic.LoadInt32(&all) != int32(numListeners) {
		t.Errorf("Mismatch in the number of listeners and received data: listeners %d, received %d",
			1, atomic.LoadInt32(&all))
	} else {
		t.Logf("All listeners received data!")
	}
}

// TestConcurrentListening4 test different UploadType & different listenChan
func TestConcurrentListening4(t *testing.T) {
	queue := NewMessageQueue()
	wg := sync.WaitGroup{}
	numListeners := 200
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	scpStartNo := 9000
	wg.Add(numListeners * 2)
	var all int32
	for i := scpStartNo; i < numListeners+scpStartNo; i++ {
		go func(id int) {
			defer wg.Done()
			ch, err := queue.ListenUploadType(ctx, commonModel.UploadType(strconv.Itoa(id)))
			if err != nil {
				t.Errorf("ListenUploadType() failed for UploadType %d: %v", id, err)
				return
			}
			select {
			case msgData := <-ch:
				if string(msgData.MsgMetaData.ScraperType) != strconv.Itoa(id) {
					t.Errorf("Received wrong ScraperType data: got %s, want %d",
						msgData.MsgMetaData.ScraperType, id)
				} else {
					t.Logf("Received correct ScraperType %d", id)
					atomic.AddInt32(&all, 1)
				}
			case <-time.After(5 * time.Second):
				t.Errorf("Timeout for UploadType %d", id)
			}
		}(i)
	}
	for i := scpStartNo; i < numListeners+scpStartNo; i++ {
		msgData := &MsgQueueData{
			MsgMetaData: MsgMetaData{
				UploadType: commonModel.UploadType(strconv.Itoa(i)),
				MsgMetaID: MsgMetaID{
					ScraperID:   strconv.Itoa(rand.Intn(100) + 1),
					MsgGroupID:  strconv.Itoa(rand.Intn(100) + 1),
					ScraperType: commonModel.ScraperType(strconv.Itoa(i)),
				},
			},
		}
		if err := queue.AddElement(msgData); err != nil {
			t.Fatalf("AddElement() failed: %v", err)
		}
		logrus.Debugf("Added element with ScraperType %d", i)
		wg.Done()
	}
	wg.Wait()
	if atomic.LoadInt32(&all) != int32(numListeners) {
		t.Errorf("Mismatch in the number of listeners and received data: listeners %d, received %d",
			numListeners, atomic.LoadInt32(&all))
	} else {
		t.Logf("All listeners received data!")
	}
}

func TestCloseChanAndReopen(t *testing.T) {
	queue := NewMessageQueue()
	numListeners := 200
	ctx, cancel := context.WithCancel(context.Background())
	scpStartNo := 10000
	var all int32
	var wg sync.WaitGroup
	var uploadType1 = commonModel.UploadType(strconv.Itoa(10001))
	var uploadType2 = commonModel.UploadType(strconv.Itoa(10002))
	wg.Add(1)
	go func() {
		ch, err := queue.ListenUploadType(ctx, uploadType1)
		if err != nil {
			t.Errorf("ListenUploadType() failed for UploadType %d: %v", scpStartNo, err)
			return
		}
		for itm := range ch {
			if itm.MsgMetaData.UploadType != uploadType1 {
				t.Errorf("Received wrong UploadType data: got %s, want %s",
					itm.MsgMetaData.UploadType, uploadType1)
			} else {
				t.Logf("Received correct UploadType %s", uploadType1)
				atomic.AddInt32(&all, 1)
			}
			if atomic.LoadInt32(&all) == int32(numListeners/2) {
				wg.Done()
			}
		}
		// after chan close, will exit to here
		if atomic.LoadInt32(&all) == int32(numListeners/2) {
			return
		}
	}()
	for i := scpStartNo; i < numListeners/2+scpStartNo; i++ {
		msgData := &MsgQueueData{
			MsgMetaData: MsgMetaData{
				UploadType: uploadType1,
				MsgMetaID: MsgMetaID{
					ScraperID:   strconv.Itoa(rand.Intn(100) + 1),
					MsgGroupID:  strconv.Itoa(rand.Intn(100) + 1),
					ScraperType: commonModel.ScraperType(strconv.Itoa(i)),
				},
			},
			FileMetaData: &clientModel.AnyFileMetaDataModel{
				PreUploadFileMetaDataModel: &clientModel.PreUploadFileMetaDataModel{
					ScraperType: commonModel.ScraperType(strconv.Itoa(i)),
					ScraperID:   strconv.Itoa(rand.Intn(100) + 1),
					ResourceUri: fmt.Sprintf("test%d.jpg", i),
				},
			},
		}
		if err := queue.AddElement(msgData); err != nil {
			t.Fatalf("AddElement() failed: %v", err)
		}
		logrus.Debugf("Added element with ScraperType %d", i)
	}
	wg.Wait()
	cancel()
	if atomic.LoadInt32(&all) != int32(numListeners/2) {
		t.Errorf("Mismatch in the number of listeners and received data: listeners %d, received %d",
			numListeners, atomic.LoadInt32(&all))
	} else {
		t.Logf("All listeners received data!")
	}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		ch, err := queue.ListenUploadType(ctx, uploadType2)
		if err != nil {
			t.Errorf("ListenUploadType() failed for UploadType %d: %v", scpStartNo, err)
			return
		}
		for itm := range ch {
			if itm.MsgMetaData.UploadType != uploadType2 {
				t.Errorf("Received wrong UploadType data: got %s, want %s",
					itm.MsgMetaData.UploadType, uploadType2)
			} else {
				t.Logf("2-Received correct UploadType %s", uploadType2)
				atomic.AddInt32(&all, 1)
			}
			if atomic.LoadInt32(&all) == int32(numListeners) {
				wg.Done()
			}
		}
		// after chan close, will exit to here
		if atomic.LoadInt32(&all) == int32(numListeners) {
			return
		}
	}()
	for i := scpStartNo; i < numListeners/2+scpStartNo; i++ {
		msgData := &MsgQueueData{
			MsgMetaData: MsgMetaData{
				UploadType: uploadType2,
				MsgMetaID: MsgMetaID{
					ScraperID:   strconv.Itoa(rand.Intn(100) + 1),
					MsgGroupID:  strconv.Itoa(rand.Intn(100) + 1),
					ScraperType: commonModel.ScraperType(strconv.Itoa(i)),
				},
			},
			FileMetaData: &clientModel.AnyFileMetaDataModel{
				PreUploadFileMetaDataModel: &clientModel.PreUploadFileMetaDataModel{
					ScraperType: commonModel.ScraperType(strconv.Itoa(i)),
					ScraperID:   strconv.Itoa(rand.Intn(100) + 1),
					ResourceUri: fmt.Sprintf("test%d.jpg", i),
				},
			},
		}
		if err := queue.AddElement(msgData); err != nil {
			t.Fatalf("AddElement() failed: %v", err)
		}
		logrus.Debugf("Added element with ScraperType %d", i)
	}
	wg.Wait()
	if atomic.LoadInt32(&all) != int32(numListeners) {
		t.Errorf("Mismatch in the number of listeners and received data: listeners %d, received %d",
			numListeners, atomic.LoadInt32(&all))
	} else {
		t.Logf("All listeners received data!")
	}
}
