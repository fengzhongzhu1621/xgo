package opentelemetry

import (
	"testing"
)

func TestNewRandomIDGenerator(t *testing.T) {
	// 测试使用固定种子创建生成器
	seed := int64(42)
	generator := newRandomIDGenerator(seed)

	if generator == nil {
		t.Fatal("newRandomIDGenerator should not return nil")
	}

	if generator.randSource == nil {
		t.Fatal("randSource should not be nil")
	}
}

func TestRandomIDGenerator_NewSpanID(t *testing.T) {
	// 使用固定种子以获得可预测的结果
	seed := int64(42)
	generator := newRandomIDGenerator(seed)

	// 生成第一个ID
	id1 := generator.newSpanID()

	// 确保ID不为负数
	if id1 < 0 {
		t.Errorf("Generated SpanID should not be negative, got: %d", id1)
	}

	// 生成第二个ID
	id2 := generator.newSpanID()

	// 确保连续生成的ID不同
	if id1 == id2 {
		t.Errorf("Two consecutively generated IDs should be different, got: %d and %d", id1, id2)
	}
}

func TestRandomIDGenerator_Deterministic(t *testing.T) {
	// 使用相同的种子创建两个生成器
	seed := int64(42)
	generator1 := newRandomIDGenerator(seed)
	generator2 := newRandomIDGenerator(seed)

	// 从两个生成器分别生成ID
	id1 := generator1.newSpanID()
	id2 := generator2.newSpanID()

	// 使用相同种子的生成器应该产生相同的第一个ID
	if id1 != id2 {
		t.Errorf("Generators with same seed should produce same first ID, got: %d and %d", id1, id2)
	}
}

func TestRandomIDGenerator_DifferentSeeds(t *testing.T) {
	// 使用不同的种子创建两个生成器
	generator1 := newRandomIDGenerator(42)
	generator2 := newRandomIDGenerator(43)

	// 从两个生成器分别生成ID
	id1 := generator1.newSpanID()
	id2 := generator2.newSpanID()

	// 使用不同种子的生成器很可能产生不同的ID
	// 注意：理论上有极小概率产生相同的ID，但实际上几乎不可能
	if id1 == id2 {
		t.Errorf("Generators with different seeds should produce different IDs, got: %d for both", id1)
	}
}
