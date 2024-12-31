// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/OpenToolkitLab/GTerm/backend/dal/model"
)

func newHost(db *gorm.DB, opts ...gen.DOOption) host {
	_host := host{}

	_host.hostDo.UseDB(db, opts...)
	_host.hostDo.UseModel(&model.Host{})

	tableName := _host.hostDo.TableName()
	_host.ALL = field.NewAsterisk(tableName)
	_host.ID = field.NewUint(tableName, "id")
	_host.CreatedAt = field.NewTime(tableName, "created_at")
	_host.UpdatedAt = field.NewTime(tableName, "updated_at")
	_host.DeletedAt = field.NewField(tableName, "deleted_at")
	_host.Label = field.NewString(tableName, "label")
	_host.Address = field.NewString(tableName, "address")
	_host.Port = field.NewUint32(tableName, "port")
	_host.Comment = field.NewString(tableName, "comment")
	_host.CredentialID = field.NewUint(tableName, "credential_id")
	_host.GroupID = field.NewUint(tableName, "group_id")
	_host.Credential = hostBelongsToCredential{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Credential", "model.Credential"),
	}

	_host.fillFieldMap()

	return _host
}

type host struct {
	hostDo

	ALL          field.Asterisk
	ID           field.Uint
	CreatedAt    field.Time
	UpdatedAt    field.Time
	DeletedAt    field.Field
	Label        field.String
	Address      field.String
	Port         field.Uint32
	Comment      field.String
	CredentialID field.Uint
	GroupID      field.Uint
	Credential   hostBelongsToCredential

	fieldMap map[string]field.Expr
}

func (h host) Table(newTableName string) *host {
	h.hostDo.UseTable(newTableName)
	return h.updateTableName(newTableName)
}

func (h host) As(alias string) *host {
	h.hostDo.DO = *(h.hostDo.As(alias).(*gen.DO))
	return h.updateTableName(alias)
}

func (h *host) updateTableName(table string) *host {
	h.ALL = field.NewAsterisk(table)
	h.ID = field.NewUint(table, "id")
	h.CreatedAt = field.NewTime(table, "created_at")
	h.UpdatedAt = field.NewTime(table, "updated_at")
	h.DeletedAt = field.NewField(table, "deleted_at")
	h.Label = field.NewString(table, "label")
	h.Address = field.NewString(table, "address")
	h.Port = field.NewUint32(table, "port")
	h.Comment = field.NewString(table, "comment")
	h.CredentialID = field.NewUint(table, "credential_id")
	h.GroupID = field.NewUint(table, "group_id")

	h.fillFieldMap()

	return h
}

func (h *host) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := h.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (h *host) fillFieldMap() {
	h.fieldMap = make(map[string]field.Expr, 11)
	h.fieldMap["id"] = h.ID
	h.fieldMap["created_at"] = h.CreatedAt
	h.fieldMap["updated_at"] = h.UpdatedAt
	h.fieldMap["deleted_at"] = h.DeletedAt
	h.fieldMap["label"] = h.Label
	h.fieldMap["address"] = h.Address
	h.fieldMap["port"] = h.Port
	h.fieldMap["comment"] = h.Comment
	h.fieldMap["credential_id"] = h.CredentialID
	h.fieldMap["group_id"] = h.GroupID

}

func (h host) clone(db *gorm.DB) host {
	h.hostDo.ReplaceConnPool(db.Statement.ConnPool)
	return h
}

func (h host) replaceDB(db *gorm.DB) host {
	h.hostDo.ReplaceDB(db)
	return h
}

type hostBelongsToCredential struct {
	db *gorm.DB

	field.RelationField
}

func (a hostBelongsToCredential) Where(conds ...field.Expr) *hostBelongsToCredential {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a hostBelongsToCredential) WithContext(ctx context.Context) *hostBelongsToCredential {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a hostBelongsToCredential) Session(session *gorm.Session) *hostBelongsToCredential {
	a.db = a.db.Session(session)
	return &a
}

func (a hostBelongsToCredential) Model(m *model.Host) *hostBelongsToCredentialTx {
	return &hostBelongsToCredentialTx{a.db.Model(m).Association(a.Name())}
}

type hostBelongsToCredentialTx struct{ tx *gorm.Association }

func (a hostBelongsToCredentialTx) Find() (result *model.Credential, err error) {
	return result, a.tx.Find(&result)
}

func (a hostBelongsToCredentialTx) Append(values ...*model.Credential) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a hostBelongsToCredentialTx) Replace(values ...*model.Credential) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a hostBelongsToCredentialTx) Delete(values ...*model.Credential) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a hostBelongsToCredentialTx) Clear() error {
	return a.tx.Clear()
}

func (a hostBelongsToCredentialTx) Count() int64 {
	return a.tx.Count()
}

type hostDo struct{ gen.DO }

type IHostDo interface {
	gen.SubQuery
	Debug() IHostDo
	WithContext(ctx context.Context) IHostDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IHostDo
	WriteDB() IHostDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IHostDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IHostDo
	Not(conds ...gen.Condition) IHostDo
	Or(conds ...gen.Condition) IHostDo
	Select(conds ...field.Expr) IHostDo
	Where(conds ...gen.Condition) IHostDo
	Order(conds ...field.Expr) IHostDo
	Distinct(cols ...field.Expr) IHostDo
	Omit(cols ...field.Expr) IHostDo
	Join(table schema.Tabler, on ...field.Expr) IHostDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IHostDo
	RightJoin(table schema.Tabler, on ...field.Expr) IHostDo
	Group(cols ...field.Expr) IHostDo
	Having(conds ...gen.Condition) IHostDo
	Limit(limit int) IHostDo
	Offset(offset int) IHostDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IHostDo
	Unscoped() IHostDo
	Create(values ...*model.Host) error
	CreateInBatches(values []*model.Host, batchSize int) error
	Save(values ...*model.Host) error
	First() (*model.Host, error)
	Take() (*model.Host, error)
	Last() (*model.Host, error)
	Find() ([]*model.Host, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Host, err error)
	FindInBatches(result *[]*model.Host, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Host) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IHostDo
	Assign(attrs ...field.AssignExpr) IHostDo
	Joins(fields ...field.RelationField) IHostDo
	Preload(fields ...field.RelationField) IHostDo
	FirstOrInit() (*model.Host, error)
	FirstOrCreate() (*model.Host, error)
	FindByPage(offset int, limit int) (result []*model.Host, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IHostDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (h hostDo) Debug() IHostDo {
	return h.withDO(h.DO.Debug())
}

func (h hostDo) WithContext(ctx context.Context) IHostDo {
	return h.withDO(h.DO.WithContext(ctx))
}

func (h hostDo) ReadDB() IHostDo {
	return h.Clauses(dbresolver.Read)
}

func (h hostDo) WriteDB() IHostDo {
	return h.Clauses(dbresolver.Write)
}

func (h hostDo) Session(config *gorm.Session) IHostDo {
	return h.withDO(h.DO.Session(config))
}

func (h hostDo) Clauses(conds ...clause.Expression) IHostDo {
	return h.withDO(h.DO.Clauses(conds...))
}

func (h hostDo) Returning(value interface{}, columns ...string) IHostDo {
	return h.withDO(h.DO.Returning(value, columns...))
}

func (h hostDo) Not(conds ...gen.Condition) IHostDo {
	return h.withDO(h.DO.Not(conds...))
}

func (h hostDo) Or(conds ...gen.Condition) IHostDo {
	return h.withDO(h.DO.Or(conds...))
}

func (h hostDo) Select(conds ...field.Expr) IHostDo {
	return h.withDO(h.DO.Select(conds...))
}

func (h hostDo) Where(conds ...gen.Condition) IHostDo {
	return h.withDO(h.DO.Where(conds...))
}

func (h hostDo) Order(conds ...field.Expr) IHostDo {
	return h.withDO(h.DO.Order(conds...))
}

func (h hostDo) Distinct(cols ...field.Expr) IHostDo {
	return h.withDO(h.DO.Distinct(cols...))
}

func (h hostDo) Omit(cols ...field.Expr) IHostDo {
	return h.withDO(h.DO.Omit(cols...))
}

func (h hostDo) Join(table schema.Tabler, on ...field.Expr) IHostDo {
	return h.withDO(h.DO.Join(table, on...))
}

func (h hostDo) LeftJoin(table schema.Tabler, on ...field.Expr) IHostDo {
	return h.withDO(h.DO.LeftJoin(table, on...))
}

func (h hostDo) RightJoin(table schema.Tabler, on ...field.Expr) IHostDo {
	return h.withDO(h.DO.RightJoin(table, on...))
}

func (h hostDo) Group(cols ...field.Expr) IHostDo {
	return h.withDO(h.DO.Group(cols...))
}

func (h hostDo) Having(conds ...gen.Condition) IHostDo {
	return h.withDO(h.DO.Having(conds...))
}

func (h hostDo) Limit(limit int) IHostDo {
	return h.withDO(h.DO.Limit(limit))
}

func (h hostDo) Offset(offset int) IHostDo {
	return h.withDO(h.DO.Offset(offset))
}

func (h hostDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IHostDo {
	return h.withDO(h.DO.Scopes(funcs...))
}

func (h hostDo) Unscoped() IHostDo {
	return h.withDO(h.DO.Unscoped())
}

func (h hostDo) Create(values ...*model.Host) error {
	if len(values) == 0 {
		return nil
	}
	return h.DO.Create(values)
}

func (h hostDo) CreateInBatches(values []*model.Host, batchSize int) error {
	return h.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (h hostDo) Save(values ...*model.Host) error {
	if len(values) == 0 {
		return nil
	}
	return h.DO.Save(values)
}

func (h hostDo) First() (*model.Host, error) {
	if result, err := h.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Host), nil
	}
}

func (h hostDo) Take() (*model.Host, error) {
	if result, err := h.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Host), nil
	}
}

func (h hostDo) Last() (*model.Host, error) {
	if result, err := h.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Host), nil
	}
}

func (h hostDo) Find() ([]*model.Host, error) {
	result, err := h.DO.Find()
	return result.([]*model.Host), err
}

func (h hostDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Host, err error) {
	buf := make([]*model.Host, 0, batchSize)
	err = h.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (h hostDo) FindInBatches(result *[]*model.Host, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return h.DO.FindInBatches(result, batchSize, fc)
}

func (h hostDo) Attrs(attrs ...field.AssignExpr) IHostDo {
	return h.withDO(h.DO.Attrs(attrs...))
}

func (h hostDo) Assign(attrs ...field.AssignExpr) IHostDo {
	return h.withDO(h.DO.Assign(attrs...))
}

func (h hostDo) Joins(fields ...field.RelationField) IHostDo {
	for _, _f := range fields {
		h = *h.withDO(h.DO.Joins(_f))
	}
	return &h
}

func (h hostDo) Preload(fields ...field.RelationField) IHostDo {
	for _, _f := range fields {
		h = *h.withDO(h.DO.Preload(_f))
	}
	return &h
}

func (h hostDo) FirstOrInit() (*model.Host, error) {
	if result, err := h.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Host), nil
	}
}

func (h hostDo) FirstOrCreate() (*model.Host, error) {
	if result, err := h.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Host), nil
	}
}

func (h hostDo) FindByPage(offset int, limit int) (result []*model.Host, count int64, err error) {
	result, err = h.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = h.Offset(-1).Limit(-1).Count()
	return
}

func (h hostDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = h.Count()
	if err != nil {
		return
	}

	err = h.Offset(offset).Limit(limit).Scan(result)
	return
}

func (h hostDo) Scan(result interface{}) (err error) {
	return h.DO.Scan(result)
}

func (h hostDo) Delete(models ...*model.Host) (result gen.ResultInfo, err error) {
	return h.DO.Delete(models)
}

func (h *hostDo) withDO(do gen.Dao) *hostDo {
	h.DO = *do.(*gen.DO)
	return h
}
