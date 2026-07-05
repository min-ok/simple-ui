package ui


type anyOption interface {
	applyToNode(Node)
}

type BoxOption interface {
	anyOption
	applyToBox(*BoxNode)
}

type SpriteOption interface {
	anyOption
	applyToSprite(*SpriteNode)
}

type ButtonOption interface {
	anyOption
	applyToButton(*ButtonNode)
}

type ScrollButtonOption interface {
	anyOption
	applyToScrollButton(*ScrollButtonNode)
}

type LabelOption interface {
	anyOption
	applyToLabel(*LabelNode)
}


type SpriteElementOpt func(*spriteElement)

func (f SpriteElementOpt) applyToNode(n Node) {}
func (f SpriteElementOpt) applyToSprite(sn *SpriteNode) { f(&sn.spriteElement) }
func (f SpriteElementOpt) applyToButton(bn *ButtonNode) { f(&bn.spriteElement) }



type CommonOpt func(Node)
func (f CommonOpt) applyToNode(n Node) { f(n) }
func (f CommonOpt) applyToBox(be *BoxNode) {}
func (f CommonOpt) applyToButton(bn *ButtonNode) {}
func (f CommonOpt) applyToScrollButton(sbn *ScrollButtonNode) {}
func (f CommonOpt) applyToSprite(sn *SpriteNode) {}
func (f CommonOpt) applyToLabel(sn *LabelNode) {}


type BoxOpt func(*BoxNode)
func (f BoxOpt) applyToNode(n Node) {}
func (f BoxOpt) applyToBox(bn *BoxNode) { f(bn) }

type ButtonOpt func(*ButtonNode)
func (f ButtonOpt) applyToNode(n Node) {}
func (f ButtonOpt) applyToButton(bn *ButtonNode) { f(bn) }

type ScrollButtonOpt func(*ScrollButtonNode)
func (f ScrollButtonOpt) applyToNode(n Node) {}
func (f ScrollButtonOpt) applyToScrollButton(sbn *ScrollButtonNode) { f(sbn) }

type SpriteOpt func(*SpriteNode)
func (f SpriteOpt) applyToNode(n Node) {}
func (f SpriteOpt) applyToSprite(sn *SpriteNode) { f(sn) }

type LabelOpt func(*LabelNode)
func (f LabelOpt) applyToNode(n Node) {}
func (f LabelOpt) applyToLabel(ln *LabelNode) { f(ln) }


func (node *BoxNode) applyToNode(n Node) { n.base().children = append(n.base().children, node) }
func (node *BoxNode) applyToBox(be *BoxNode) {}
func (node *BoxNode) applyToButton(be *ButtonNode) {}
func (node *BoxNode) applyToScrollButton(sbn *ScrollButtonNode) {}
func (node *BoxNode) applyToSprite(sn *SpriteNode) {}
func (node *BoxNode) applyToLabel(sn *LabelNode) {}

func (node *ButtonNode) applyToNode(n Node) { n.base().children = append(n.base().children, node) }
func (node *ButtonNode) applyToBox(be *BoxNode) {}
func (node *ButtonNode) applyToButton(be *ButtonNode) {}
func (node *ButtonNode) applyToScrollButton(sbn *ScrollButtonNode) {}
func (node *ButtonNode) applyToSprite(sn *SpriteNode) {}
func (node *ButtonNode) applyToLabel(sn *LabelNode) {}

func (node *ScrollButtonNode) applyToNode(n Node) { n.base().children = append(n.base().children, node) }
func (node *ScrollButtonNode) applyToBox(be *BoxNode) {}
func (node *ScrollButtonNode) applyToButton(be *ButtonNode) {}
func (node *ScrollButtonNode) applyToScrollButton(sbn *ScrollButtonNode) {}
func (node *ScrollButtonNode) applyToSprite(sn *SpriteNode) {}
func (node *ScrollButtonNode) applyToLabel(sn *LabelNode) {}

func (node *SpriteNode) applyToNode(n Node) { n.base().children = append(n.base().children, node) }
func (node *SpriteNode) applyToBox(be *BoxNode) {}
func (node *SpriteNode) applyToButton(be *ButtonNode) {}
func (node *SpriteNode) applyToScrollButton(sbn *ScrollButtonNode) {}
func (node *SpriteNode) applyToSprite(sn *SpriteNode) {}
func (node *SpriteNode) applyToLabel(sn *LabelNode) {}

func (node *LabelNode) applyToNode(n Node) { n.base().children = append(n.base().children, node) }
func (node *LabelNode) applyToBox(be *BoxNode) {}
func (node *LabelNode) applyToButton(be *ButtonNode) {}
func (node *LabelNode) applyToScrollButton(sbn *ScrollButtonNode) {}
func (node *LabelNode) applyToSprite(sn *SpriteNode) {}
func (node *LabelNode) applyToLabel(sn *LabelNode) {}
