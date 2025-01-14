package tree

import (
	"fmt"
	"testing"

	tree "github.com/duke-git/lancet/v2/datastructure/tree"
)

// BSTree是一种二叉搜索树数据结构，其中每个节点最多有两个子节点，分别称为左子节点和右子节点。
// 在BSTree中：leftNode < rootNode < rightNode。类型T应该在约束.Comparator接口中实现比较函数。

//func NewBSTree[T any](rootData T, comparator constraints.Comparator) *BSTree[T]
//
//type BSTree[T any] struct {
//	root       *datastructure.TreeNode[T]
//	comparator constraints.Comparator
//}
//
//type TreeNode[T any] struct {
//	Value T
//	Left  *TreeNode[T]
//	Right *TreeNode[T]
//}

type intComparator struct{}

func (c *intComparator) Compare(v1, v2 any) int {
	val1, _ := v1.(int)
	val2, _ := v2.(int)

	if val1 < val2 {
		return -1
	} else if val1 > val2 {
		return 1
	}
	return 0
}

// func (t *BSTree[T]) Insert(data T)
func TestBSTreeInsert(t *testing.T) {
	bstree := tree.NewBSTree(6, &intComparator{})
	bstree.Insert(7)
	bstree.Insert(5)
	bstree.Insert(2)
	bstree.Insert(4)

	fmt.Println(bstree.PreOrderTraverse()) //6, 5, 2, 4, 7
}

// func (t *BSTree[T]) Delete(data T)
func TestBSTreeDelete(t *testing.T) {
	bstree := tree.NewBSTree(6, &intComparator{})
	bstree.Insert(7)
	bstree.Insert(5)
	bstree.Insert(2)
	bstree.Insert(4)

	bstree.Delete(4)

	fmt.Println(bstree.PreOrderTraverse()) //2, 5, 6, 7
}

// func (t *BSTree[T]) PreOrderTraverse() []T
func TestBSTreePreOrderTraverse(t *testing.T) {
	bstree := tree.NewBSTree(6, &intComparator{})
	bstree.Insert(7)
	bstree.Insert(5)
	bstree.Insert(2)
	bstree.Insert(4)

	fmt.Println(bstree.PreOrderTraverse()) //6, 5, 2, 4, 7
}

// func (t *BSTree[T]) InOrderTraverse() []T
func TestBSTreeInOrderTraverse(t *testing.T) {
	bstree := tree.NewBSTree(6, &intComparator{})
	bstree.Insert(7)
	bstree.Insert(5)
	bstree.Insert(2)
	bstree.Insert(4)

	fmt.Println(bstree.InOrderTraverse()) //2, 4, 5, 6, 7
}

// func (t *BSTree[T]) PostOrderTraverse() []T
func TestBSTreePostOrderTraverse(t *testing.T) {
	bstree := tree.NewBSTree(6, &intComparator{})
	bstree.Insert(7)
	bstree.Insert(5)
	bstree.Insert(2)
	bstree.Insert(4)

	fmt.Println(bstree.PostOrderTraverse()) //5, 2, 4, 7, 6
}

// func (t *BSTree[T]) LevelOrderTraverse() []T
func TestBSTreeLevelOrderTraverse(t *testing.T) {
	bstree := tree.NewBSTree(6, &intComparator{})
	bstree.Insert(7)
	bstree.Insert(5)
	bstree.Insert(2)
	bstree.Insert(4)

	fmt.Println(bstree.LevelOrderTraverse()) //6, 5, 7, 2, 4
}

// func (t *BSTree[T]) Depth() int
func TestBSTreeDepth(t *testing.T) {
	bstree := tree.NewBSTree(6, &intComparator{})
	bstree.Insert(7)
	bstree.Insert(5)
	bstree.Insert(2)
	bstree.Insert(4)

	fmt.Println(bstree.Depth()) //4
}

// func (t *BSTree[T]) HasSubTree(subTree *BSTree[T]) bool
func TestBSTreeHasSubTree(t *testing.T) {
	superTree := tree.NewBSTree(8, &intComparator{})
	superTree.Insert(4)
	superTree.Insert(5)
	superTree.Insert(6)
	superTree.Insert(9)
	superTree.Insert(4)

	subTree := tree.NewBSTree(5, &intComparator{})
	subTree.Insert(4)
	subTree.Insert(6)

	fmt.Println(superTree.HasSubTree(subTree)) //true
	fmt.Println(subTree.HasSubTree(superTree)) //false
}

// func (t *BSTree[T]) Print()
func TestBSTreePrint(t *testing.T) {
	bstree := tree.NewBSTree(6, &intComparator{})
	bstree.Insert(7)
	bstree.Insert(5)
	bstree.Insert(2)
	bstree.Insert(4)

	bstree.Print()
	//        6
	//       / \
	//      /   \
	//     /     \
	//    /       \
	//    5       7
	//   /
	//  /
	//  2
	//   \
	//    4
}
