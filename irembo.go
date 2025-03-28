package irembopay

// IremboPay is the main client for interacting with the IremboPay API
type IremboPay struct {
	Config  *Config
	Invoice *InvoiceService
	Batch   *BatchService
	Payment *PaymentService
}

// NewIremboPay creates a new IremboPay client
func NewIremboPay(config *Config) *IremboPay {
	client := NewClient(config)

	return &IremboPay{
		Config:  config,
		Invoice: NewInvoiceService(client, config),
		Batch:   NewBatchService(client, config),
		Payment: NewPaymentService(client, config),
	}
}

// NewSandboxClient creates a new IremboPay client for the sandbox environment
func NewSandboxClient(secretKey string, opts ...ConfigOption) (*IremboPay, error) {
	config, err := NewConfig(Sandbox, secretKey, opts...)
	if err != nil {
		return nil, err
	}

	return NewIremboPay(config), nil
}

// NewProductionClient creates a new IremboPay client for the production environment
func NewProductionClient(secretKey string, opts ...ConfigOption) (*IremboPay, error) {
	config, err := NewConfig(Production, secretKey, opts...)
	if err != nil {
		return nil, err
	}

	return NewIremboPay(config), nil
}
