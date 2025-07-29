package game

type Cell struct {
	IsMine     bool
	IsRevealed bool
	IsFlagged  bool
	Adjacent   int
}

func NewCell() *Cell {
	return &Cell{
		IsMine:     false,
		IsRevealed: false,
		IsFlagged:  false,
		Adjacent:   0,
	}
}

func (c *Cell) Reveal() {
	if !c.IsFlagged {
		c.IsRevealed = true
	}
}

func (c *Cell) ToggleFlag() {
	if !c.IsRevealed {
		c.IsFlagged = !c.IsFlagged
	}
}

func (c *Cell) SetMine() {
	c.IsMine = true
}

func (c *Cell) SetAdjacent(count int) {
	c.Adjacent = count
}