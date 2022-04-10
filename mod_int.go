type ModInt struct {
        mod, value int
}

func NewModInt(mod int) *ModInt {
        return &ModInt{mod: mod}
}

func (m *ModInt) Get() int {
        return m.value
}

func (m *ModInt) Add(x int) {
        m.value = (m.value + x % m.mod) % m.mod
}

func (m *ModInt) Mul(x int) {
        m.value = (m.value * (x % m.mod)) % m.mod
}

func (m *ModInt) Sub(x int) {
        m.value = m.value - x % m.mod
        if m.value < 0 {
                m.value += m.mod
        }
}

