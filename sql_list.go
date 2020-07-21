package sqlpool

import (
	"container/list"
	"github.com/civet148/gotools/log"
	"sync"
)

type sqlList struct {
	elements list.List
	locker   sync.RWMutex
}

func init() {

}

func (s *sqlList) lock() {
	s.locker.Lock()
}

func (s *sqlList) unlock() {
	s.locker.Unlock()
}

func (s *sqlList) rlock() {
	s.locker.RLock()
}

func (s *sqlList) runlock() {
	s.locker.RUnlock()
}

/*
插入SQL事件节点
排序方式：优先级高的优先插入在链表前(保证每次插入后链表始终保持有序)
*/
func (s *sqlList) insert(node *sqlNode) {
	s.lock()
	defer s.unlock()
	if s.elements.Len() == 0 {
		s.elements.PushBack(node)
	} else {

		//从前往后找到第一个优先级低于node的链表节点并在其前面插入
		if ok := s.insertPriorityLess(node); ok {
			return
		}

		//从后往前找第一个优先级等于node优先级的链表节点并在其后面插入
		if ok := s.insertPriorityEqualReverse(node); ok {
			return
		}
		//只剩下比自己优先级高的则在最后插入
		s.elements.PushBack(node)
	}
}

/*
找到优先级低于当前节点的链表节点并插入(正向查找)
*/
func (s *sqlList) insertPriorityLess(node *sqlNode) (ok bool) {

	//正向查找
	for e := s.elements.Front(); e != nil; e = e.Next() {
		n := e.Value.(*sqlNode)
		if n.Priority < node.Priority { //找到链表中优先级低于要插入的节点
			s.elements.InsertBefore(node, e)
			return true
		}
	}
	return false
}

/*
找到优先级等于当前节点的链表节点并插入到最后(反向查找)
*/
func (s *sqlList) insertPriorityEqualReverse(node *sqlNode) (ok bool) {

	//反向查找
	for e := s.elements.Back(); e != nil; e = e.Prev() {
		n := e.Value.(*sqlNode)
		if n.Priority == node.Priority {
			s.elements.InsertAfter(node, e)
			return true
		}
	}
	return false
}

/*
获取链表最前面的节点
*/
func (s *sqlList) front() (node *sqlNode) {
	s.lock()
	defer s.unlock()
	e := s.elements.Front()
	if e == nil {
		return nil
	}
	node = e.Value.(*sqlNode)
	s.elements.Remove(e)
	return
}

/*
遍历输出链表所有节点
*/
func (s *sqlList) print(strQueueName string) (node *sqlNode) {
	s.rlock()
	defer s.runlock()
	log.Debugf("queue [%v] list element count [%v] ", strQueueName, s.elements.Len())
	return
}
