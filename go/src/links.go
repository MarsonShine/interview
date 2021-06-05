package interview

type ListNode struct {
	Val  int
	Next *ListNode
}

func DivideMerge(lists []*ListNode, l, r int) *ListNode {
	if l == r {
		return lists[l]
	}
	if l > r {
		return nil
	}
	mid := (l + r) / 2

	return MergeTowList(DivideMerge(lists, l, mid), DivideMerge(lists, mid+1, r))
}

func MergeTowList(one *ListNode, two *ListNode) *ListNode {
	if one == nil && two == nil {
		return &ListNode{}
	}
	if (one == nil || one == &ListNode{}) {
		return two
	}
	if (two == nil || two == &ListNode{}) {
		return one
	}
	head := ListNode{}
	tail := &head
	for one != nil && two != nil {
		if one.Val < two.Val {
			tail.Next = one
			one = one.Next
		} else {
			tail.Next = two
			two = two.Next
		}
		tail = tail.Next
	}
	// 剩余节点
	if one != nil {
		tail.Next = one
	} else {
		tail.Next = two
	}
	return head.Next
}

func MergeKList(lists []*ListNode) *ListNode {
	return DivideMerge(lists, 0, len(lists)-1)
	// var r *ListNode = nil
	// for i := 0; i < len(lists); i++ {
	// 	r = MergeTowList(r, lists[i])
	// }
	// return r
}

func HasCycle(head *ListNode) bool {
	if head == nil {
		return false
	}
	if head.Next == nil {
		return false
	}
	if head.Next.Next == nil {
		return false
	}
	slow, fast := head, head.Next
	for slow != fast {
		if fast == nil || fast.Next == nil {
			return false
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	return true
	// step1
	// root := head
	// for root.Next != nil {
	// 	if root.Next == root.Next.Next {
	// 		return true
	// 	}
	// 	root = root.Next
	// }
	// return false
	// step2
}

func CreateListNode() *ListNode {
	root := ListNode{
		Val: 3,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 0,
				Next: &ListNode{
					Val:  -4,
					Next: nil,
				},
			},
		},
	}
	return &root
}
