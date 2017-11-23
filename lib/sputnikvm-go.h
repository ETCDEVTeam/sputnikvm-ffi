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
