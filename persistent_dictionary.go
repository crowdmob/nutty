package nutty

import (
  "time"
  "github.com/crowdmob/goamz/dynamodb"
)

type PersistentModelWithDictionaryKey interface {
  DictionaryKey()           (string, string) // name and value
  ModelName()               string
  SetCreatedAt(time.Time)
  SetUpdatedAt(time.Time)
}

func (nuttyApp *App) AddToDynamoDB(m PersistentModelWithDictionaryKey) error {
  m.SetCreatedAt(time.Now())
  
  marshalledAttrs, err := dynamodb.MarshalAttributes(m)
  if err != nil { return err }
  
  _, hashKeyValue := m.DictionaryKey()
  nuttyApp.DynamoDBTableForModel(m).PutItem(
    hashKeyValue, 
    "", 
    marshalledAttrs,
  )
  return nil
}

func (nuttyApp *App) GetFromDynamoDB(key string, dest PersistentModelWithDictionaryKey) error {
  attrs, err := nuttyApp.DynamoDBTableForModel(dest).GetItem(key, "")
  if err != nil { return err }
  
  err = dynamodb.UnmarshalAttributes(&attrs, dest)
  return err
}

func (nuttyApp *App) ExistsInDynamoDB(key string, m PersistentModelWithDictionaryKey) (bool, error) {
  attrs, err := nuttyApp.DynamoDBTableForModel(m).GetItem(key, "")
  if err != nil { return false, err }
  
  hashKeyName, _ := m.DictionaryKey()
  if attrs[hashKeyName] != nil {
    return true, nil
  } 
  return false, nil
}

func (nuttyApp *App) UpdateInDynamoDB(m PersistentModelWithDictionaryKey) error {
  m.SetUpdatedAt(time.Now())
  panic("UpdateInDynamoDB Not Yet Implemented")
}

func (nuttyApp *App) RemoveFromDynamoDB(m PersistentModelWithDictionaryKey) error {
  panic("RemoveFromDynamoDB Not Yet Implemented")
}

func (nuttyApp *App) DynamoDBTableForModel(m PersistentModelWithDictionaryKey) *dynamodb.Table {
  hashKeyName, _ := m.DictionaryKey()
  return (&dynamodb.Server{nuttyApp.AwsAuth, nuttyApp.AwsRegion}).NewTable(
    nuttyApp.DynamoDBTableNameForModel(m), 
    dynamodb.PrimaryKey{
      dynamodb.NewStringAttribute(hashKeyName, ""), 
      nil,
    },
  )
}

func (nuttyApp *App) DynamoDBTableNameForModel(m PersistentModelWithDictionaryKey) string {
  return m.ModelName() + "s-" + nuttyApp.Name + "-" + nuttyApp.Env
}

