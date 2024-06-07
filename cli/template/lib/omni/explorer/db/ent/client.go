// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/google/uuid"
	"github.com/omni-network/omni/explorer/db/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/omni-network/omni/explorer/db/ent/block"
	"github.com/omni-network/omni/explorer/db/ent/chain"
	"github.com/omni-network/omni/explorer/db/ent/msg"
	"github.com/omni-network/omni/explorer/db/ent/receipt"
	"github.com/omni-network/omni/explorer/db/ent/xprovidercursor"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Block is the client for interacting with the Block builders.
	Block *BlockClient
	// Chain is the client for interacting with the Chain builders.
	Chain *ChainClient
	// Msg is the client for interacting with the Msg builders.
	Msg *MsgClient
	// Receipt is the client for interacting with the Receipt builders.
	Receipt *ReceiptClient
	// XProviderCursor is the client for interacting with the XProviderCursor builders.
	XProviderCursor *XProviderCursorClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Block = NewBlockClient(c.config)
	c.Chain = NewChainClient(c.config)
	c.Msg = NewMsgClient(c.config)
	c.Receipt = NewReceiptClient(c.config)
	c.XProviderCursor = NewXProviderCursorClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:             ctx,
		config:          cfg,
		Block:           NewBlockClient(cfg),
		Chain:           NewChainClient(cfg),
		Msg:             NewMsgClient(cfg),
		Receipt:         NewReceiptClient(cfg),
		XProviderCursor: NewXProviderCursorClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:             ctx,
		config:          cfg,
		Block:           NewBlockClient(cfg),
		Chain:           NewChainClient(cfg),
		Msg:             NewMsgClient(cfg),
		Receipt:         NewReceiptClient(cfg),
		XProviderCursor: NewXProviderCursorClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Block.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Block.Use(hooks...)
	c.Chain.Use(hooks...)
	c.Msg.Use(hooks...)
	c.Receipt.Use(hooks...)
	c.XProviderCursor.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Block.Intercept(interceptors...)
	c.Chain.Intercept(interceptors...)
	c.Msg.Intercept(interceptors...)
	c.Receipt.Intercept(interceptors...)
	c.XProviderCursor.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *BlockMutation:
		return c.Block.mutate(ctx, m)
	case *ChainMutation:
		return c.Chain.mutate(ctx, m)
	case *MsgMutation:
		return c.Msg.mutate(ctx, m)
	case *ReceiptMutation:
		return c.Receipt.mutate(ctx, m)
	case *XProviderCursorMutation:
		return c.XProviderCursor.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// BlockClient is a client for the Block schema.
type BlockClient struct {
	config
}

// NewBlockClient returns a client for the Block from the given config.
func NewBlockClient(c config) *BlockClient {
	return &BlockClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `block.Hooks(f(g(h())))`.
func (c *BlockClient) Use(hooks ...Hook) {
	c.hooks.Block = append(c.hooks.Block, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `block.Intercept(f(g(h())))`.
func (c *BlockClient) Intercept(interceptors ...Interceptor) {
	c.inters.Block = append(c.inters.Block, interceptors...)
}

// Create returns a builder for creating a Block entity.
func (c *BlockClient) Create() *BlockCreate {
	mutation := newBlockMutation(c.config, OpCreate)
	return &BlockCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Block entities.
func (c *BlockClient) CreateBulk(builders ...*BlockCreate) *BlockCreateBulk {
	return &BlockCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *BlockClient) MapCreateBulk(slice any, setFunc func(*BlockCreate, int)) *BlockCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &BlockCreateBulk{err: fmt.Errorf("calling to BlockClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*BlockCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &BlockCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Block.
func (c *BlockClient) Update() *BlockUpdate {
	mutation := newBlockMutation(c.config, OpUpdate)
	return &BlockUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *BlockClient) UpdateOne(b *Block) *BlockUpdateOne {
	mutation := newBlockMutation(c.config, OpUpdateOne, withBlock(b))
	return &BlockUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *BlockClient) UpdateOneID(id int) *BlockUpdateOne {
	mutation := newBlockMutation(c.config, OpUpdateOne, withBlockID(id))
	return &BlockUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Block.
func (c *BlockClient) Delete() *BlockDelete {
	mutation := newBlockMutation(c.config, OpDelete)
	return &BlockDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *BlockClient) DeleteOne(b *Block) *BlockDeleteOne {
	return c.DeleteOneID(b.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *BlockClient) DeleteOneID(id int) *BlockDeleteOne {
	builder := c.Delete().Where(block.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &BlockDeleteOne{builder}
}

// Query returns a query builder for Block.
func (c *BlockClient) Query() *BlockQuery {
	return &BlockQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeBlock},
		inters: c.Interceptors(),
	}
}

// Get returns a Block entity by its id.
func (c *BlockClient) Get(ctx context.Context, id int) (*Block, error) {
	return c.Query().Where(block.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BlockClient) GetX(ctx context.Context, id int) *Block {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryMsgs queries the msgs edge of a Block.
func (c *BlockClient) QueryMsgs(b *Block) *MsgQuery {
	query := (&MsgClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := b.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(block.Table, block.FieldID, id),
			sqlgraph.To(msg.Table, msg.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, block.MsgsTable, block.MsgsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(b.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryReceipts queries the receipts edge of a Block.
func (c *BlockClient) QueryReceipts(b *Block) *ReceiptQuery {
	query := (&ReceiptClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := b.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(block.Table, block.FieldID, id),
			sqlgraph.To(receipt.Table, receipt.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, block.ReceiptsTable, block.ReceiptsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(b.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *BlockClient) Hooks() []Hook {
	return c.hooks.Block
}

// Interceptors returns the client interceptors.
func (c *BlockClient) Interceptors() []Interceptor {
	return c.inters.Block
}

func (c *BlockClient) mutate(ctx context.Context, m *BlockMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&BlockCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&BlockUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&BlockUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&BlockDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Block mutation op: %q", m.Op())
	}
}

// ChainClient is a client for the Chain schema.
type ChainClient struct {
	config
}

// NewChainClient returns a client for the Chain from the given config.
func NewChainClient(c config) *ChainClient {
	return &ChainClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `chain.Hooks(f(g(h())))`.
func (c *ChainClient) Use(hooks ...Hook) {
	c.hooks.Chain = append(c.hooks.Chain, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `chain.Intercept(f(g(h())))`.
func (c *ChainClient) Intercept(interceptors ...Interceptor) {
	c.inters.Chain = append(c.inters.Chain, interceptors...)
}

// Create returns a builder for creating a Chain entity.
func (c *ChainClient) Create() *ChainCreate {
	mutation := newChainMutation(c.config, OpCreate)
	return &ChainCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Chain entities.
func (c *ChainClient) CreateBulk(builders ...*ChainCreate) *ChainCreateBulk {
	return &ChainCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *ChainClient) MapCreateBulk(slice any, setFunc func(*ChainCreate, int)) *ChainCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &ChainCreateBulk{err: fmt.Errorf("calling to ChainClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*ChainCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &ChainCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Chain.
func (c *ChainClient) Update() *ChainUpdate {
	mutation := newChainMutation(c.config, OpUpdate)
	return &ChainUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ChainClient) UpdateOne(ch *Chain) *ChainUpdateOne {
	mutation := newChainMutation(c.config, OpUpdateOne, withChain(ch))
	return &ChainUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ChainClient) UpdateOneID(id int) *ChainUpdateOne {
	mutation := newChainMutation(c.config, OpUpdateOne, withChainID(id))
	return &ChainUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Chain.
func (c *ChainClient) Delete() *ChainDelete {
	mutation := newChainMutation(c.config, OpDelete)
	return &ChainDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ChainClient) DeleteOne(ch *Chain) *ChainDeleteOne {
	return c.DeleteOneID(ch.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ChainClient) DeleteOneID(id int) *ChainDeleteOne {
	builder := c.Delete().Where(chain.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ChainDeleteOne{builder}
}

// Query returns a query builder for Chain.
func (c *ChainClient) Query() *ChainQuery {
	return &ChainQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeChain},
		inters: c.Interceptors(),
	}
}

// Get returns a Chain entity by its id.
func (c *ChainClient) Get(ctx context.Context, id int) (*Chain, error) {
	return c.Query().Where(chain.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ChainClient) GetX(ctx context.Context, id int) *Chain {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *ChainClient) Hooks() []Hook {
	return c.hooks.Chain
}

// Interceptors returns the client interceptors.
func (c *ChainClient) Interceptors() []Interceptor {
	return c.inters.Chain
}

func (c *ChainClient) mutate(ctx context.Context, m *ChainMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ChainCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ChainUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ChainUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ChainDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Chain mutation op: %q", m.Op())
	}
}

// MsgClient is a client for the Msg schema.
type MsgClient struct {
	config
}

// NewMsgClient returns a client for the Msg from the given config.
func NewMsgClient(c config) *MsgClient {
	return &MsgClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `msg.Hooks(f(g(h())))`.
func (c *MsgClient) Use(hooks ...Hook) {
	c.hooks.Msg = append(c.hooks.Msg, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `msg.Intercept(f(g(h())))`.
func (c *MsgClient) Intercept(interceptors ...Interceptor) {
	c.inters.Msg = append(c.inters.Msg, interceptors...)
}

// Create returns a builder for creating a Msg entity.
func (c *MsgClient) Create() *MsgCreate {
	mutation := newMsgMutation(c.config, OpCreate)
	return &MsgCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Msg entities.
func (c *MsgClient) CreateBulk(builders ...*MsgCreate) *MsgCreateBulk {
	return &MsgCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *MsgClient) MapCreateBulk(slice any, setFunc func(*MsgCreate, int)) *MsgCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &MsgCreateBulk{err: fmt.Errorf("calling to MsgClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*MsgCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &MsgCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Msg.
func (c *MsgClient) Update() *MsgUpdate {
	mutation := newMsgMutation(c.config, OpUpdate)
	return &MsgUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *MsgClient) UpdateOne(m *Msg) *MsgUpdateOne {
	mutation := newMsgMutation(c.config, OpUpdateOne, withMsg(m))
	return &MsgUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *MsgClient) UpdateOneID(id int) *MsgUpdateOne {
	mutation := newMsgMutation(c.config, OpUpdateOne, withMsgID(id))
	return &MsgUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Msg.
func (c *MsgClient) Delete() *MsgDelete {
	mutation := newMsgMutation(c.config, OpDelete)
	return &MsgDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *MsgClient) DeleteOne(m *Msg) *MsgDeleteOne {
	return c.DeleteOneID(m.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *MsgClient) DeleteOneID(id int) *MsgDeleteOne {
	builder := c.Delete().Where(msg.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &MsgDeleteOne{builder}
}

// Query returns a query builder for Msg.
func (c *MsgClient) Query() *MsgQuery {
	return &MsgQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeMsg},
		inters: c.Interceptors(),
	}
}

// Get returns a Msg entity by its id.
func (c *MsgClient) Get(ctx context.Context, id int) (*Msg, error) {
	return c.Query().Where(msg.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *MsgClient) GetX(ctx context.Context, id int) *Msg {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryBlock queries the block edge of a Msg.
func (c *MsgClient) QueryBlock(m *Msg) *BlockQuery {
	query := (&BlockClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := m.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(msg.Table, msg.FieldID, id),
			sqlgraph.To(block.Table, block.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, msg.BlockTable, msg.BlockPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(m.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryReceipts queries the receipts edge of a Msg.
func (c *MsgClient) QueryReceipts(m *Msg) *ReceiptQuery {
	query := (&ReceiptClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := m.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(msg.Table, msg.FieldID, id),
			sqlgraph.To(receipt.Table, receipt.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, msg.ReceiptsTable, msg.ReceiptsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(m.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *MsgClient) Hooks() []Hook {
	hooks := c.hooks.Msg
	return append(hooks[:len(hooks):len(hooks)], msg.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *MsgClient) Interceptors() []Interceptor {
	return c.inters.Msg
}

func (c *MsgClient) mutate(ctx context.Context, m *MsgMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&MsgCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&MsgUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&MsgUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&MsgDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Msg mutation op: %q", m.Op())
	}
}

// ReceiptClient is a client for the Receipt schema.
type ReceiptClient struct {
	config
}

// NewReceiptClient returns a client for the Receipt from the given config.
func NewReceiptClient(c config) *ReceiptClient {
	return &ReceiptClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `receipt.Hooks(f(g(h())))`.
func (c *ReceiptClient) Use(hooks ...Hook) {
	c.hooks.Receipt = append(c.hooks.Receipt, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `receipt.Intercept(f(g(h())))`.
func (c *ReceiptClient) Intercept(interceptors ...Interceptor) {
	c.inters.Receipt = append(c.inters.Receipt, interceptors...)
}

// Create returns a builder for creating a Receipt entity.
func (c *ReceiptClient) Create() *ReceiptCreate {
	mutation := newReceiptMutation(c.config, OpCreate)
	return &ReceiptCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Receipt entities.
func (c *ReceiptClient) CreateBulk(builders ...*ReceiptCreate) *ReceiptCreateBulk {
	return &ReceiptCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *ReceiptClient) MapCreateBulk(slice any, setFunc func(*ReceiptCreate, int)) *ReceiptCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &ReceiptCreateBulk{err: fmt.Errorf("calling to ReceiptClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*ReceiptCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &ReceiptCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Receipt.
func (c *ReceiptClient) Update() *ReceiptUpdate {
	mutation := newReceiptMutation(c.config, OpUpdate)
	return &ReceiptUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ReceiptClient) UpdateOne(r *Receipt) *ReceiptUpdateOne {
	mutation := newReceiptMutation(c.config, OpUpdateOne, withReceipt(r))
	return &ReceiptUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ReceiptClient) UpdateOneID(id int) *ReceiptUpdateOne {
	mutation := newReceiptMutation(c.config, OpUpdateOne, withReceiptID(id))
	return &ReceiptUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Receipt.
func (c *ReceiptClient) Delete() *ReceiptDelete {
	mutation := newReceiptMutation(c.config, OpDelete)
	return &ReceiptDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ReceiptClient) DeleteOne(r *Receipt) *ReceiptDeleteOne {
	return c.DeleteOneID(r.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ReceiptClient) DeleteOneID(id int) *ReceiptDeleteOne {
	builder := c.Delete().Where(receipt.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ReceiptDeleteOne{builder}
}

// Query returns a query builder for Receipt.
func (c *ReceiptClient) Query() *ReceiptQuery {
	return &ReceiptQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeReceipt},
		inters: c.Interceptors(),
	}
}

// Get returns a Receipt entity by its id.
func (c *ReceiptClient) Get(ctx context.Context, id int) (*Receipt, error) {
	return c.Query().Where(receipt.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ReceiptClient) GetX(ctx context.Context, id int) *Receipt {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryBlock queries the block edge of a Receipt.
func (c *ReceiptClient) QueryBlock(r *Receipt) *BlockQuery {
	query := (&BlockClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(receipt.Table, receipt.FieldID, id),
			sqlgraph.To(block.Table, block.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, receipt.BlockTable, receipt.BlockPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryMsgs queries the msgs edge of a Receipt.
func (c *ReceiptClient) QueryMsgs(r *Receipt) *MsgQuery {
	query := (&MsgClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(receipt.Table, receipt.FieldID, id),
			sqlgraph.To(msg.Table, msg.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, receipt.MsgsTable, receipt.MsgsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ReceiptClient) Hooks() []Hook {
	hooks := c.hooks.Receipt
	return append(hooks[:len(hooks):len(hooks)], receipt.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *ReceiptClient) Interceptors() []Interceptor {
	return c.inters.Receipt
}

func (c *ReceiptClient) mutate(ctx context.Context, m *ReceiptMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ReceiptCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ReceiptUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ReceiptUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ReceiptDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Receipt mutation op: %q", m.Op())
	}
}

// XProviderCursorClient is a client for the XProviderCursor schema.
type XProviderCursorClient struct {
	config
}

// NewXProviderCursorClient returns a client for the XProviderCursor from the given config.
func NewXProviderCursorClient(c config) *XProviderCursorClient {
	return &XProviderCursorClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `xprovidercursor.Hooks(f(g(h())))`.
func (c *XProviderCursorClient) Use(hooks ...Hook) {
	c.hooks.XProviderCursor = append(c.hooks.XProviderCursor, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `xprovidercursor.Intercept(f(g(h())))`.
func (c *XProviderCursorClient) Intercept(interceptors ...Interceptor) {
	c.inters.XProviderCursor = append(c.inters.XProviderCursor, interceptors...)
}

// Create returns a builder for creating a XProviderCursor entity.
func (c *XProviderCursorClient) Create() *XProviderCursorCreate {
	mutation := newXProviderCursorMutation(c.config, OpCreate)
	return &XProviderCursorCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of XProviderCursor entities.
func (c *XProviderCursorClient) CreateBulk(builders ...*XProviderCursorCreate) *XProviderCursorCreateBulk {
	return &XProviderCursorCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *XProviderCursorClient) MapCreateBulk(slice any, setFunc func(*XProviderCursorCreate, int)) *XProviderCursorCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &XProviderCursorCreateBulk{err: fmt.Errorf("calling to XProviderCursorClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*XProviderCursorCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &XProviderCursorCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for XProviderCursor.
func (c *XProviderCursorClient) Update() *XProviderCursorUpdate {
	mutation := newXProviderCursorMutation(c.config, OpUpdate)
	return &XProviderCursorUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *XProviderCursorClient) UpdateOne(xc *XProviderCursor) *XProviderCursorUpdateOne {
	mutation := newXProviderCursorMutation(c.config, OpUpdateOne, withXProviderCursor(xc))
	return &XProviderCursorUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *XProviderCursorClient) UpdateOneID(id uuid.UUID) *XProviderCursorUpdateOne {
	mutation := newXProviderCursorMutation(c.config, OpUpdateOne, withXProviderCursorID(id))
	return &XProviderCursorUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for XProviderCursor.
func (c *XProviderCursorClient) Delete() *XProviderCursorDelete {
	mutation := newXProviderCursorMutation(c.config, OpDelete)
	return &XProviderCursorDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *XProviderCursorClient) DeleteOne(xc *XProviderCursor) *XProviderCursorDeleteOne {
	return c.DeleteOneID(xc.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *XProviderCursorClient) DeleteOneID(id uuid.UUID) *XProviderCursorDeleteOne {
	builder := c.Delete().Where(xprovidercursor.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &XProviderCursorDeleteOne{builder}
}

// Query returns a query builder for XProviderCursor.
func (c *XProviderCursorClient) Query() *XProviderCursorQuery {
	return &XProviderCursorQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeXProviderCursor},
		inters: c.Interceptors(),
	}
}

// Get returns a XProviderCursor entity by its id.
func (c *XProviderCursorClient) Get(ctx context.Context, id uuid.UUID) (*XProviderCursor, error) {
	return c.Query().Where(xprovidercursor.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *XProviderCursorClient) GetX(ctx context.Context, id uuid.UUID) *XProviderCursor {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *XProviderCursorClient) Hooks() []Hook {
	return c.hooks.XProviderCursor
}

// Interceptors returns the client interceptors.
func (c *XProviderCursorClient) Interceptors() []Interceptor {
	return c.inters.XProviderCursor
}

func (c *XProviderCursorClient) mutate(ctx context.Context, m *XProviderCursorMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&XProviderCursorCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&XProviderCursorUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&XProviderCursorUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&XProviderCursorDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown XProviderCursor mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Block, Chain, Msg, Receipt, XProviderCursor []ent.Hook
	}
	inters struct {
		Block, Chain, Msg, Receipt, XProviderCursor []ent.Interceptor
	}
)
