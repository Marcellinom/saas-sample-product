package session

import "github.com/gin-gonic/gin"

type Data struct {
	ctx     *gin.Context
	id      string
	data    map[string]interface{}
	storage Storage
}

func (d Data) Get(key string) interface{} {
	return d.data[key]
}

func (d Data) Set(key string, value interface{}) {
	d.data[key] = value
}

func (d Data) Delete(key string) {
	delete(d.data, key)
}

func (d Data) Clear() {
	for key := range d.data {
		delete(d.data, key)
	}
}

func (d Data) Save() error {
	return d.storage.Save(d.ctx, d.id, d.data)
}

func NewData(ctx *gin.Context, id string, data map[string]interface{}, storage Storage) Data {
	return Data{
		ctx:     ctx,
		id:      id,
		data:    data,
		storage: storage,
	}
}

func Default(ctx *gin.Context) Data {
	dataIf, exists := ctx.Get("session")
	if !exists {
		panic("session not found in context, make sure you have called session.StartSession middleware")
	}
	data, ok := dataIf.(Data)
	if !ok {
		panic("session is not of type session.Data")
	}

	return data
}
