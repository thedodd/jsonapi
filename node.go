package jsonapi

type OnePayload struct {
	Data     *Node              `json:"data"`
	Included []*Node            `json:"included,omitempty"`
	Links    *map[string]string `json:"links,omitempty"`
}

type ManyPayload struct {
	Data     []*Node            `json:"data"`
	Included []*Node            `json:"included,omitempty"`
	Links    *map[string]string `json:"links,omitempty"`
}

type Node struct {
	Type          string                 `json:"type"`
	ID            string                 `json:"id"`
	ClientID      string                 `json:"client-id,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

// Extend applies the attirbutes on the rhs to the callee
func (n *Node) Extend(node *Node) {
	for attr, val := range node.Attributes {
		n.AddAttriute(attr, val)
	}

	for rel, val := range node.Relationships {
		n.AddRelationship(rel, val)
	}
}

// AddRelationship adds a relationship to the Node
func (n *Node) AddRelationship(name string, val interface{}) {
	if n.Relationships == nil {
		n.Relationships = make(map[string]interface{})
	}

	n.Relationships[name] = val
}

// AddAttriute adds an attribute to the Node
func (n *Node) AddAttriute(name string, val interface{}) {
	if n.Attributes == nil {
		n.Attributes = make(map[string]interface{})
	}

	n.Attributes[name] = val
}

type RelationshipOneNode struct {
	Data  *Node              `json:"data"`
	Links *map[string]string `json:"links,omitempty"`
}

type RelationshipManyNode struct {
	Data  []*Node            `json:"data"`
	Links *map[string]string `json:"links,omitempty"`
}
