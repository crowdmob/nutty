package nutty

import (
	"github.com/crowdmob/goamz/dynamodb"
	"strconv"
	"time"
)

type PersistentModelWithDictionaryKey interface {
	DictionaryKey() (string, string) // name and value
	ModelName() string
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}

func (nuttyApp *App) AddToDynamoDB(m PersistentModelWithDictionaryKey, optionalAppName string) error {
	m.SetCreatedAt(time.Now())

	marshalledAttrs, err := dynamodb.MarshalAttributes(m)
	if err != nil {
		return err
	}

	_, hashKeyValue := m.DictionaryKey()
	_, err = nuttyApp.DynamoDBTableForModel(m, optionalAppName).PutItem(
		hashKeyValue,
		"",
		marshalledAttrs,
	)
	return err
}

func (nuttyApp *App) GetFromDynamoDB(key string, dest PersistentModelWithDictionaryKey, optionalAppName string) error {
	attrs, err := nuttyApp.DynamoDBTableForModel(dest, optionalAppName).GetItem(&dynamodb.Key{HashKey: key})
	if err != nil {
		return err
	}

	err = dynamodb.UnmarshalAttributes(&attrs, dest)
	return err
}

func (nuttyApp *App) ExistsInDynamoDB(key string, m PersistentModelWithDictionaryKey, optionalAppName string) (bool, error) {
	attrs, err := nuttyApp.DynamoDBTableForModel(m, optionalAppName).GetItem(&dynamodb.Key{HashKey: key})
	if err != nil {
		return false, nil
	} // treat an erroneous response as empty

	hashKeyName, _ := m.DictionaryKey()
	if attrs[hashKeyName] != nil {
		return true, nil
	}
	return false, nil
}

func (nuttyApp *App) UpdateInDynamoDB(m PersistentModelWithDictionaryKey, optionalAppName string) error {
	m.SetUpdatedAt(time.Now())

	marshalledAttrs, err := dynamodb.MarshalAttributes(m)
	if err != nil {
		return err
	}

	// TODO once dynamodb has a proper UpdateItem we should use that instead of PutItem
	_, hashKeyValue := m.DictionaryKey()
	_, err = nuttyApp.DynamoDBTableForModel(m, optionalAppName).PutItem(
		hashKeyValue,
		"",
		marshalledAttrs,
	)
	return err
}


func (nuttyApp *App) AppendToStringSetsInDynamoDB(m PersistentModelWithDictionaryKey, attributeAppends map[string]string, optionalAppName string) error {
	attrs := make([]dynamodb.Attribute, len(attributeAppends))

	i := 0
	for key, val := range attributeAppends {
		attrs[i] = dynamodb.Attribute{
			Type:  dynamodb.TYPE_STRING_SET,
			Name:  key,
			Value: val,
		}
		i++
	}

	_, hashKeyValue := m.DictionaryKey()
	_, err := nuttyApp.DynamoDBTableForModel(m, optionalAppName).AddItem(
		&dynamodb.Key{HashKey: hashKeyValue},
		attrs,
	)
	return err
}


func intsToDynamoAttrs(attributeIncrements map[string]int64) []dynamodb.Attribute {
	attrs := make([]dynamodb.Attribute, len(attributeIncrements))
	i := 0
	for key, val := range attributeIncrements {
		attrs[i] = dynamodb.Attribute{
			Type:  dynamodb.TYPE_NUMBER,
			Name:  key,
			Value: strconv.FormatInt(val, 10),
		}
		i++
	}
	return attrs
}
func (nuttyApp *App) SetIntsInDynamoDB(m PersistentModelWithDictionaryKey, attributeIncrements map[string]int64, optionalAppName string) error {
	_, hashKeyValue := m.DictionaryKey()
	_, err := nuttyApp.DynamoDBTableForModel(m, optionalAppName).UpdateItem(
		&dynamodb.Key{HashKey: hashKeyValue},
		intsToDynamoAttrs(attributeIncrements),
	)
	return err
}

func (nuttyApp *App) IncrementIntsInDynamoDB(m PersistentModelWithDictionaryKey, attributeIncrements map[string]int64, optionalAppName string) error {
	_, hashKeyValue := m.DictionaryKey()
	_, err := nuttyApp.DynamoDBTableForModel(m, optionalAppName).AddItem(
		&dynamodb.Key{HashKey: hashKeyValue},
		intsToDynamoAttrs(attributeIncrements),
	)
	return err
}

func (nuttyApp *App) IncrementNumericsInDynamoDB(m PersistentModelWithDictionaryKey, intAttributeIncrements map[string]int64, floatAttributeIncrements map[string]float64, optionalAppName string) error {
	attrs := make([]dynamodb.Attribute, len(intAttributeIncrements)+len(floatAttributeIncrements))

	i := 0
	for key, val := range floatAttributeIncrements {
		attrs[i] = dynamodb.Attribute{
			Type:  dynamodb.TYPE_NUMBER,
			Name:  key,
			Value: strconv.FormatFloat(val, 'f', -1, 64),
		}
		i++
	}

	for key, val := range intAttributeIncrements {
		attrs[i] = dynamodb.Attribute{
			Type:  dynamodb.TYPE_NUMBER,
			Name:  key,
			Value: strconv.FormatInt(val, 10),
		}
		i++
	}
	_, hashKeyValue := m.DictionaryKey()
	_, err := nuttyApp.DynamoDBTableForModel(m, optionalAppName).AddItem(
		&dynamodb.Key{HashKey: hashKeyValue},
		attrs,
	)
	return err
}

func (nuttyApp *App) IncrementFloatsInDynamoDB(m PersistentModelWithDictionaryKey, attributeIncrements map[string]float64, optionalAppName string) error {
	attrs := make([]dynamodb.Attribute, len(attributeIncrements))

	i := 0
	for key, val := range attributeIncrements {
		attrs[i] = dynamodb.Attribute{
			Type:  dynamodb.TYPE_NUMBER,
			Name:  key,
			Value: strconv.FormatFloat(val, 'f', -1, 64),
		}
		i++
	}

	_, hashKeyValue := m.DictionaryKey()
	_, err := nuttyApp.DynamoDBTableForModel(m, optionalAppName).AddItem(
		&dynamodb.Key{HashKey: hashKeyValue},
		attrs,
	)
	return err
}

func (nuttyApp *App) RemoveFromDynamoDB(m PersistentModelWithDictionaryKey) error {
	panic("RemoveFromDynamoDB Not Yet Implemented")
}

func (nuttyApp *App) DynamoDBTableForModel(m PersistentModelWithDictionaryKey, optionalAppName string) *dynamodb.Table {
	hashKeyName, _ := m.DictionaryKey()
	return (&dynamodb.Server{nuttyApp.AwsAuth, nuttyApp.AwsRegion}).NewTable(
		nuttyApp.DynamoDBTableNameForModel(m, optionalAppName),
		dynamodb.PrimaryKey{
			dynamodb.NewStringAttribute(hashKeyName, ""),
			nil,
		},
	)
}

func (nuttyApp *App) DynamoDBTableNameForModel(m PersistentModelWithDictionaryKey, optionalAppName string) string {
	if optionalAppName == "" {
		optionalAppName = nuttyApp.Name
	}
	return m.ModelName() + "s-" + optionalAppName + "-" + nuttyApp.Env
}
