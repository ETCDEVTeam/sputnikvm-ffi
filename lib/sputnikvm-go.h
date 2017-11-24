typedef struct {
  unsigned char data[20];
} sputnikvm_address;

typedef struct {
  unsigned char data[32];
} sputnikvm_gas;

typedef struct {
  unsigned char data[32];
} sputnikvm_u256;

typedef unsigned char sputnikvm_action;
extern const unsigned char CALL_ACTION;
extern const unsigned char CREATE_ACTION;

typedef struct {
  sputnikvm_address caller;
  sputnikvm_gas gas_price;
  sputnikvm_gas gas_limit;
  sputnikvm_action action;
  sputnikvm_address action_address;
  sputnikvm_u256 value;
  unsigned char *input;
  unsigned int input_len;
  sputnikvm_u256 nonce;
} sputnikvm_transaction;

typedef struct {
  sputnikvm_address beneficiary;
  unsigned long long int timestamp;
  sputnikvm_u256 number;
  sputnikvm_u256 difficulty;
  sputnikvm_gas gas_limit;
} sputnikvm_header_params;

typedef enum {
  none, account, account_code, account_storage, blockhash
} sputnikvm_require_type;

typedef struct {
  sputnikvm_address address;
  sputnikvm_u256 key;
} sputnikvm_require_value_account_storage;

typedef union {
  sputnikvm_address account;
  sputnikvm_require_value_account_storage account_storage;
  sputnikvm_u256 blockhash;
} sputnikvm_require_value;

typedef struct {
  sputnikvm_require_type type;
  sputnikvm_require_value value;
} sputnikvm_require;

typedef struct sputnikvm_vm_S sputnikvm_vm_t;

extern sputnikvm_vm_t *
sputnikvm_new_frontier(sputnikvm_transaction transaction, sputnikvm_header_params header);

extern sputnikvm_vm_t *
sputnikvm_new_homestead(sputnikvm_transaction transaction, sputnikvm_header_params header);

extern sputnikvm_vm_t *
sputnikvm_new_eip150(sputnikvm_transaction transaction, sputnikvm_header_params header);

extern sputnikvm_vm_t *
sputnikvm_new_eip160(sputnikvm_transaction transaction, sputnikvm_header_params header);

extern sputnikvm_require
sputnikvm_fire(sputnikvm_vm_t *vm);

extern void
sputnikvm_free(sputnikvm_vm_t *vm);

extern sputnikvm_transaction
sputnikvm_default_transaction(void);

extern sputnikvm_header_params
sputnikvm_default_header_params(void);
