package export_model

type ColumnReplace struct {
	oldField []OldField
	newField []string
}

type OldField struct {
	oldField string
	isDelete string
}

func NewColumnReplace() (m *ColumnReplace) {
	return new(ColumnReplace)
}

func (m *ColumnReplace) AddNewField(newField ...string) {
	m.newField = append(m.newField, newField...)
	return
}

func (m *ColumnReplace) AddOldField(field, isDelete string) {
	m.oldField = append(m.oldField, OldField{
		oldField: field,
		isDelete: isDelete,
	})

	return
}

func (m *ColumnReplace) GetNewField() []string {
	return m.newField
}

func (m *ColumnReplace) GetOldField() []OldField {
	return m.oldField
}

func (item OldField) GetOldField() string {
	return item.oldField
}

func (item OldField) GetIsDelete() string {
	return item.isDelete
}