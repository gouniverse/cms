package cms

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	uid "github.com/lesichkovm/gouid"
	"gorm.io/gorm"
)

const (
	// EntityStatusActive entity "active" status
	EntityStatusActive = "active"
	// EntityStatusInactive entity "inactive" status
	EntityStatusInactive = "inactive"
)

// Entity type
type Entity struct {
	ID          string     `gorm:"type:varchar(40);column:id;primary_key;"`
	Status      string     `gorm:"type:varchar(10);column:status;"`
	Type        string     `gorm:"type:varchar(40);column:type;"`
	Name        string     `gorm:"type:varchar(255);column:name;DEFAULT NULL;"`
	Description string     `gorm:"type:longtext;column:description;"`
	CreatedAt   time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt   time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt   *time.Time `gorm:"type:datetime;olumn:deleted_at;DEFAULT NULL;"`

	Attributes []EntityAttribute
}

// TableName teh name of the User table
func (Entity) TableName() string {
	return prefix + "entities_entity"
}

// BeforeCreate adds UID to model
func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.NanoUid()
	e.ID = uuid
	return nil
}

// GetAttribute the name of the User table
func (e *Entity) GetAttribute(attributeKey string) *EntityAttribute {
	entityAttribute := &EntityAttribute{}

	result := GetDb().First(&entityAttribute, "entity_id=? AND attribute_key = ?", e.ID, attributeKey)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return entityAttribute
}

// EntityAttribute type
type EntityAttribute struct {
	ID             string     `gorm:"type:varchar(40);column:id;primary_key;"`
	EntityID       string     `gorm:"type:varchar(40);column:entity_id;"`
	AttributeKey   string     `gorm:"type:varchar(255);column:attribute_key;DEFAULT NULL;"`
	AttributeValue string     `gorm:"type:longtext;column:attribute_value;"`
	CreatedAt      time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt      time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt      *time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`
}

// TableName teh name of the User table
func (EntityAttribute) TableName() string {
	return prefix + "entities_attribute"
}

// BeforeCreate adds UID to model
func (e *EntityAttribute) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.NanoUid()
	e.ID = uuid
	return nil
}

// SetValue serializes the values
func (e *EntityAttribute) SetValue(value interface{}) bool {
	// interfaceType := reflect.ValueOf(value).Kind() //value.(type) //reflect.TypeOf(value)

	// if interfaceType.String() == "bool" {
	// 	e.AttributeValue = fmt.Sprintf("%v", value)
	// 	return true
	// }
	// if interfaceType.String() == "string" {
	// 	e.AttributeValue = fmt.Sprintf("%v", value)
	// 	return true
	// }

	bytes, err := json.Marshal(value)

	if err != nil {
		return false
	}

	e.AttributeValue = string(bytes)

	return true
}

// GetValue serializes the values
func (e *EntityAttribute) GetValue() interface{} {
	var value interface{}
	err := json.Unmarshal([]byte(e.AttributeValue), &value)

	if err != nil {
		panic("JSOB error unmarshaliibg attribute" + err.Error())
	}

	return value
}

// EntityCreate creates a new entity
func EntityCreate(entityType string) *Entity {
	entity := &Entity{Type: entityType, Status: EntityStatusActive}

	dbResult := GetDb().Create(&entity)

	if dbResult.Error != nil {
		return nil
	}

	return entity
}

// EntityCreateWithAttributes func
func EntityCreateWithAttributes(entityType string, attributes map[string]interface{}) *Entity {
	// Note the use of tx as the database handle once you are within a transaction
	tx := GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil
	}

	//return tx.Commit().Error

	entity := &Entity{Type: entityType, Status: EntityStatusActive}

	dbResult := tx.Create(&entity)

	if dbResult.Error != nil {
		tx.Rollback()
		return nil
	}

	//entityAttributes := make([]EntityAttribute, 0)
	for k, v := range attributes {
		ea := EntityAttribute{EntityID: entity.ID, AttributeKey: k} //, AttributeValue: value}
		ea.SetValue(v)

		dbResult2 := tx.Create(&ea)
		if dbResult2.Error != nil {
			tx.Rollback()
			return nil
		}
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return nil
	}

	return entity

}

// EntityAttributeFind finds an entity by ID
func EntityAttributeFind(entityID string, attributeKey string) *EntityAttribute {
	entityAttribute := &EntityAttribute{}

	result := GetDb().First(&entityAttribute, "entity_id = ? AND attribute_key = ?", entityID, attributeKey)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	return entityAttribute
}

// EntityAttributesUpsert upserts and entity attribute
func EntityAttributesUpsert(entityID string, attributes map[string]interface{}) bool {
	tx := GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false
	}

	for k, v := range attributes {
		entityAttribute := EntityAttributeFind(entityID, k)

		if entityAttribute == nil {
			entityAttribute = &EntityAttribute{EntityID: entityID, AttributeKey: k}
			entityAttribute.SetValue(v)

			dbResult := tx.Create(&entityAttribute)
			if dbResult.Error != nil {
				tx.Rollback()
				return false
			}

		}

		entityAttribute.SetValue(v)
		dbResult := tx.Save(entityAttribute)
		if dbResult.Error != nil {
			return false
		}
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return false
	}

	return true

}

// EntityAttributeUpsert upserts and entity attribute
func EntityAttributeUpsert(entityID string, attributeKey string, attributeValue interface{}) bool {
	entityAttribute := EntityAttributeFind(entityID, attributeKey)

	if entityAttribute == nil {
		entityAttribute = &EntityAttribute{EntityID: entityID, AttributeKey: attributeKey}
		entityAttribute.SetValue(attributeValue)

		dbResult := GetDb().Create(&entityAttribute)
		if dbResult.Error != nil {
			return false
		}

		return true
	}

	entityAttribute.SetValue(attributeValue)
	dbResult := GetDb().Save(entityAttribute)
	if dbResult.Error != nil {
		return false
	}

	return true

}

// EntityFindByID finds an entity by ID
func EntityFindByID(entityID string) *Entity {
	entity := &Entity{}

	resultEntity := GetDb().First(&entity, "id = ?", entityID)

	if resultEntity.Error != nil {
		log.Panic(resultEntity.Error)
	}

	if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	// DEBUG: log.Println(entity)

	return entity
}

// EntityFindByAttribute finds an entity by attribute
func EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) *Entity {
	entityAttribute := &EntityAttribute{}

	subQuery := GetDb().Model(&Entity{}).Select("id").Where("type = ?", entityType)
	result := GetDb().First(&entityAttribute, "entity_id IN (?) AND attribute_key = ? AND attribute_value = ?", subQuery, attributeKey, attributeValue)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	// DEBUG: log.Println(entityAttribute)

	entity := &Entity{}

	resultEntity := GetDb().First(&entity, "id = ?", entityAttribute.EntityID)

	if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if resultEntity.Error != nil {
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	return entity
}

// EntityList lists entities
func EntityList(entityType string, offset uint64, perPage uint64, search string, orderBy string, sort string) []Entity {
	entityList := []Entity{}
	result := GetDb().Where("type = ?", entityType).Order(orderBy + " " + sort).Offset(int(offset)).Limit(int(perPage)).Find(&entityList)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return entityList
}

// EntityCount counts entities
func EntityCount(entityType string) uint64 {
	var count int64
	GetDb().Model(&Entity{}).Where("type = ?", entityType).Count(&count)
	return uint64(count)
	// sqlStr, args, _ := squirrel.Select("COUNT(*) AS count").From(TableArticle).Limit(1).ToSql()

	// entities := Query(sqlStr, args...)

	// count, _ := strconv.ParseUint(entities[0]["count"], 10, 64)

	// return count
}
