package cms

type CustomEntityStructure struct {
	// Type of the entity
	Type string
	// Label to display referencing the entity
	TypeLabel string
	// Name of the entity
	Name string
	// AttributeList list of attributes
	AttributeList []CustomAttributeStructure
	// Group to which this entity belongs (i.e. Shop, Users, etc)
	Group string
}

type CustomAttributeStructure struct {
	// Name the name of the attribute
	Name string
	// Type of the attribute - string, float, int
	Type string
	// FormControlLabel label to display for the control
	FormControlLabel string
	// FormControlType the type of form control - input, textarea. etc
	FormControlType string
	// FormControlHelp help message to display for the control
	FormControlHelp string
	// BelongsToType describes a Belong To relationsip
	BelongsToType string
	// HasOneType describes a Has One relationsip
	HasOneType string
	// HasManyType describes a Has Many relationsip
	HasManyType string
}
