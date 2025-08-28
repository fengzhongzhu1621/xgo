package randutils

// Option 用于指定概率的百分比
type Option struct {
	Percent int
}

// Probability 表示一个概率值
type Probability struct {
	Percent int // 当前的运行概率，单位是%，默认为100
}

// Should 判断当前概率是否应该执行
func (p *Probability) Should() bool {
	return RandomInt(1, 100) <= p.Percent
}

// SetPercent 设置概率的百分比
func (p *Probability) SetPercent(percent int) {
	p.Percent = (&Probability{Percent: percent}).valid().Percent
}

// valid 确保设置的百分比在有效范围内[0 到 100]
func (p *Probability) valid() *Probability {
	p.Percent = If(p.Percent > 100, 100, p.Percent)
	p.Percent = If(p.Percent < 0, 0, p.Percent)
	return p
}

func NewProbability(o *Option) *Probability {
	return (&Probability{
		Percent: o.Percent,
	}).valid()
}

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
