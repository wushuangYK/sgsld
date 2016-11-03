package battle

import (
	"fmt"
	"errors"
)

type Queue struct {
	Element []Event //Element
}

func NewQueue() *Queue {
	return &Queue{}
}

func (queue *Queue)Enqueue(value ...Event){
	queue.Element = append(queue.Element,value...)
}

//返回下一个元素
func (queue *Queue)Next()(Event){
	if queue.Size() > 0 {
		return queue.Element[0]
	}
	return Event{} //read empty queue
}

//返回下一个元素,并从Queue移除元素
func (queue *Queue)Dequeue()(Event, error){
	next_value := queue.Next()
	if queue.Size()> 0 {
		queue.Element = queue.Element[1:queue.Size()]
		return next_value,nil
	}
	return next_value,errors.New("Queue为空.") //read empty queue
}

//交换值
func (queue *Queue)Swap(other *Queue){
	switch{
	case queue.Size() == 0 && other.Size() == 0:
		return
	case other.Size() == 0 :
		other.Element = queue.Element[:queue.Size()]
		queue.Element = nil
	case queue.Size()== 0 :
		queue.Element = other.Element
		other.Element = nil
	default:
		queue.Element,other.Element = other.Element,queue.Element
	}
	return
}

//修改指定索引的元素
func (queue *Queue)Set(idx int,value Event)(error){
	if idx >= 0 && queue.Size() > 0 && queue.Size() > idx{
		queue.Element[idx] = value
		return nil
	}
	return errors.New("Set失败!")
}

//返回指定索引的元素
func (queue *Queue)Get(idx int)(Event){
	if idx >= 0 && queue.Size() > 0 && queue.Size() > idx {
		return queue.Element[idx]
	}
	return Event{} //read empty queue
}

//Queue的size
func (queue *Queue)Size()(int){
	return len(queue.Element)
}

//是否为空
func (queue *Queue)Empty()(bool){
	if queue.Element == nil || queue.Size() == 0 {
		return true
	}
	return false
}

//打印
func (queue *Queue)Print(){
	for i := len(queue.Element) - 1; i >= 0; i--{
		fmt.Println(i,"=>",queue.Element[i])
	}
}
