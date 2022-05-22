package pocketsmith

// based on https://github.com/Medium/medium-sdk-go/blob/master/medium.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// endpoint is the default endpoint  for Pocketsmith's API.
	endpoint = "https://api.pocketsmith.com/v2"
	// defaultTimeout is the default timeout duration used on HTTP requests.
	defaultTimeout = time.Second * 300
	// defaultCode is the default error code for failures.
	defaultCode = -1
)

// CreateTransactionAccountTransactionOptions defines the options for creating
// a transaction in a transaction account in PocketSmith.
type CreateTransactionAccountTransactionOptions struct {
	Date         string  `json:"date"`
	Payee        string  `json:"payee"`
	Amount       float64 `json:"amount"`
	Labels       string  `json:"labels,omitempty"` // must be comma seperated list
	CategoryID   int32   `json:"category_id,omitempty"`
	Note         string  `json:"note,omitempty"`
	Memo         string  `json:"memo,omitempty"`
	IsTransfer   bool    `json:"is_transfer,omitempty"`
	ChequeNumber string  `json:"cheque_number,omitempty"`
}

// CreateAccountForUserOptions defines the options for creating an account for
// a user in PocketSmith.
type CreateAccountForUserOptions struct {
	InstitutionID int    `json:"institution_id"`
	Title         string `json:"title"`
	CurrencyCode  string `json:"currency_code"`
	Type          string `json:"type"`
}

// ListTransactionsInTransactionAccountOptions defines the options for listing
// transactions in a transaction account in PocketSmith.
type ListTransactionsInTransactionAccountOptions struct {
	StartDate         string `json:"start_date,omitempty"`
	EndDate           string `json:"end_date,omitempty"`
	OnlyUncategorised int32  `json:"only_uncategorized,omitempty"`
	Type              string `json:"type,omitempty"`
}

// UpdateTransactionOptions defines the options for updating a transaction
// in PocketSmith.
type UpdateTransactionOptions struct {
	ID           int32   `json:"-"`
	Labels       string  `json:"labels,omitempty"` // must be comma seperated list.
	Payee        string  `json:"payee,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	Date         string  `json:"date,omitempty"`
	IsTransfer   bool    `json:"is_transfer,omitempty"`
	CategoryID   int32   `json:"category_id,omitempty"`
	Note         string  `json:"note,omitempty"`
	Memo         string  `json:"memo,omitempty"`
	ChequeNumber string  `json:"cheque_number,omitempty"`
}

// ListAttachmentsForUserOptions defines the options for listing attachments for
// a given user in PocketSmith.
type ListAttachmentsForUserOptions struct {
	Unassigned int `json:"unassigned"`
}

// CreateAttachmentsForUserOptions defines the options for creating attachments
// for a given user in PocketSmith.
type CreateAttachmentsForUserOptions struct {
	Title    string `json:"title"`
	FileName string `json:"file_name"`
	FileData string `json:"file_data"`
}

// AssignAttachmentToTransactionOptions defines the options for assigning an
// attachment to a transaction in PocketSmith.
type AssignAttachmentToTransactionOptions struct {
	TransactionID int32 `json:"-"`
	AttachmentID  int   `json:"attachment_id"`
}

// CreateInstitutionForAuthedUserOptions defines the options for creating a
// institution in PocketSmith.
type CreateInstitutionForAuthedUserOptions struct {
	Title        string `json:"title"`
	CurrencyCode string `json:"currency_code"`
}

// User defines a PocketSmith user.
// https://developers.pocketsmith.com/reference#get_users-id
type User struct {
	ID                      int       `json:"id"`
	Login                   string    `json:"login"`
	Name                    string    `json:"name"`
	Email                   string    `json:"email"`
	AvatarURL               string    `json:"avatar_url"`
	TimeZone                string    `json:"time_zone"`
	WeekStartDay            int       `json:"week_start_day"`
	BaseCurrencyCode        string    `json:"base_currency_code"`
	AlwaysShowBaseCurrency  bool      `json:"always_show_base_currency"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	UsingMultipleCurrencies bool      `json:"using_multiple_currencies"`
	LastLoggedInAt          time.Time `json:"last_logged_in_at"`
	LastActivityAt          time.Time `json:"last_activity_at"`
}

// Institution defines a PocketSmith institution.
// https://developers.pocketsmith.com/reference#get_institutions-id
type Institution struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	CurrencyCode string    `json:"currency_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TransactionAccount defines a PocketSmith transaction account.
// https://developers.pocketsmith.com/reference#get_transaction-accounts-id
type TransactionAccount struct {
	ID                           int         `json:"id"`
	Name                         string      `json:"name"`
	Number                       string      `json:"number"`
	Type                         string      `json:"type"`
	CurrencyCode                 string      `json:"currency_code"`
	CurrentBalance               float64     `json:"current_balance"`
	CurrentBalanceInBaseCurrency float64     `json:"current_balance_in_base_currency"`
	CurrentBalanceExchangeRate   float64     `json:"current_balance_exchange_rate"`
	CurrentBalanceDate           string      `json:"current_balance_date"`
	StartingBalance              float64     `json:"starting_balance"`
	StartingBalanceDate          string      `json:"starting_balance_date"`
	Institution                  Institution `json:"institution"`
	CreatedAt                    time.Time   `json:"created_at"`
	UpdatedAt                    time.Time   `json:"updated_at"`
}

// Account defines a PocketSmith account.
// https://developers.pocketsmith.com/reference#get_accounts-id
type Account struct {
	ID                           int                  `json:"id"`
	Title                        string               `json:"title"`
	Type                         string               `json:"type"`
	IsNetWorth                   bool                 `json:"is_net_worth"`
	CurrencyCode                 string               `json:"currency_code"`
	CurrentBalance               float64              `json:"current_balance"`
	CurrentBalanceInBaseCurrency float64              `json:"current_balance_in_base_currency"`
	CurrentBalanceExchangeRate   float64              `json:"current_balance_exchange_rate"`
	CurrentBalanceDate           string               `json:"current_balance_date"`
	PrimaryTransactionAccount    TransactionAccount   `json:"primary_transaction_account"`
	TransactionAccounts          []TransactionAccount `json:"transaction_accounts"`
	CreatedAt                    time.Time            `json:"created_at"`
	UpdatedAt                    time.Time            `json:"updated_at"`
}

// Category defines a PocketSmith category.
// https://developers.pocketsmith.com/reference#get_categories-id
type Category struct {
	ID         int32       `json:"id"`
	Title      string      `json:"title"`
	Colour     string      `json:"colour"`
	Children   []*Category `json:"children"`
	ParentID   *int        `json:"parent_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	IsTransfer bool        `json:"is_transfer"`
}

// Transaction defines a PocketSmith transaction.
// https://developers.pocketsmith.com/reference#get_transactions-id
type Transaction struct {
	ID                   int32              `json:"id"`
	Date                 string             `json:"date"`
	Payee                string             `json:"payee"`
	OriginalPayee        string             `json:"original_payee"`
	Amount               float64            `json:"amount"`
	UploadSource         string             `json:"upload_source"`
	ClosingBalance       float64            `json:"closing_balance"`
	Memo                 string             `json:"memo"`
	Note                 string             `json:"note"`
	Labels               []string           `json:"labels"`
	Type                 string             `json:"type"`
	Status               string             `json:"status"`
	IsTransfer           bool               `json:"is_transfer"`
	NeedsReview          bool               `json:"needs_review"`
	ChequeNumber         string             `json:"cheque_number"`
	AmountInBaseCurrency float64            `json:"amount_in_base_currency"`
	Category             Category           `json:"category"`
	TransactionAccount   TransactionAccount `json:"transaction_account"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

// Attachment defines a PocketSmith attachment.
// https://developers.pocketsmith.com/reference#get_attachments-id
type Attachment struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	FileName        string `json:"file_name"`
	Type            string `json:"type"`
	ContentType     string `json:"content_type"`
	ContentTypeMeta struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Extension   string `json:"extension"`
	} `json:"content_type_meta"`
	OriginalURL string `json:"original_url"`
	Variants    struct {
		ThumbURL string `json:"thumb_url"`
		LargeURL string `json:"large_url"`
	} `json:"variants"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Error defines an error received when making a request to the API.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error returns a string representing the error, satisfying the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("pocketsmith: %s (%d)", e.Message, e.Code)
}

// PocketSmith defines the PocketSmith client.
type PocketSmith struct {
	httpClient *http.Client // the HTTP client used to query against the API.
	token      string       // the PocketSmith API token.
	user       User         // the PocketSmith user associated with the token.
}

// NewClient returns a new PocketSmith API client which can be used to make
// requests to the PocketSmith API.
func NewClient(token string, retrieveTokenUser bool) (*PocketSmith, error) {
	ps := &PocketSmith{
		token:      token,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	if retrieveTokenUser {
		user, err := ps.GetAuthedUser()
		if err != nil {
			return nil, err
		}
		ps.user = user
	}
	return ps, nil
}

// clientRequest defines information that can be used to make a request to PocketSmith.
type clientRequest struct {
	method string
	path   string
	data   interface{}
}

// GetAuthedUser returns the authed PocketSmith user associated with the given token.
// https://developers.pocketsmith.com/reference#get_me
func (ps *PocketSmith) GetAuthedUser() (User, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   "/me",
	}
	user := User{}
	_, err := ps.request(cr, &user)
	return user, err
}

// ListTransactionAccountsForAuthedUser returns a list of transaction accounts
// for the authed user associated with the given token.
func (ps *PocketSmith) ListTransactionAccountsForAuthedUser() ([]TransactionAccount, error) {
	return ps.ListTransactionAccountsForUser(ps.user.ID)
}

// ListTransactionAccountsForUser returns a list of transaction accounts for a given user id.
// https://developers.pocketsmith.com/reference#get_users-id-transaction-accounts
func (ps *PocketSmith) ListTransactionAccountsForUser(userID int) ([]TransactionAccount, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/transaction_accounts", userID),
	}
	accounts := []TransactionAccount{}
	_, err := ps.request(cr, &accounts)
	return accounts, err
}

// DeleteAccount uses the given account id to delete a PocketSmith account.
// https://developers.pocketsmith.com/reference#delete_accounts-id
func (ps *PocketSmith) DeleteAccount(accountID int) error {
	cr := clientRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/accounts/%v", accountID),
	}
	_, err := ps.request(cr, nil)
	return err
}

// ListAccountsForAuthedUser returns a list of accounts for the authed user
// associated with the given token.
func (ps *PocketSmith) ListAccountsForAuthedUser() ([]Account, error) {
	return ps.ListAccountsForUser(ps.user.ID)
}

// ListAccountsForUser returns a list of accounts for a given user id.
// https://developers.pocketsmith.com/reference#get_users-id-accounts
func (ps *PocketSmith) ListAccountsForUser(userID int) ([]Account, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/accounts", userID),
	}
	accounts := []Account{}
	_, err := ps.request(cr, &accounts)
	return accounts, err
}

// ListInstitutionsForAuthedUser returns a list of institutions for the authed
// user associated with the given token.
func (ps *PocketSmith) ListInstitutionsForAuthedUser() ([]Institution, error) {
	return ps.ListInstitutionsForUser(ps.user.ID)
}

// ListInstitutionsForUser returns a list of institutions for a given user id.
// https://developers.pocketsmith.com/reference#get_users-id-institutions
func (ps *PocketSmith) ListInstitutionsForUser(userID int) ([]Institution, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/institutions", userID),
	}
	institutions := []Institution{}
	_, err := ps.request(cr, &institutions)
	return institutions, err
}

// CreateTransactionAccountTransaction creates a transaction in the given account id.
// https://developers.pocketsmith.com/reference#post_transaction-accounts-id-transactions
func (ps *PocketSmith) CreateTransactionAccountTransaction(accountID int, options CreateTransactionAccountTransactionOptions) (Transaction, error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transaction_accounts/%v/transactions", accountID),
		data:   options,
	}
	transaction := Transaction{}
	_, err := ps.request(cr, &transaction)
	return transaction, err
}

// CreateAccountForAuthedUser creates an account for the authed user associated
// with the given token.
func (ps *PocketSmith) CreateAccountForAuthedUser(options CreateAccountForUserOptions) (Account, error) {
	return ps.CreateAccountForUser(ps.user.ID, options)
}

// CreateAccountForUser creates an account for the given user id.
// https://developers.pocketsmith.com/reference#post_users-id-accounts
func (ps *PocketSmith) CreateAccountForUser(userID int, options CreateAccountForUserOptions) (Account, error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/accounts", userID),
		data:   options,
	}
	account := Account{}
	_, err := ps.request(cr, &account)
	return account, err
}

// ListTransactionsInTransactionAccount returns a list of transactions for a given transaction account id.
// https://developers.pocketsmith.com/reference#get_transaction-accounts-id-transactions
func (ps *PocketSmith) ListTransactionsInTransactionAccount(accountID int, options ListTransactionsInTransactionAccountOptions) ([]Transaction, error) {
	transactions := []Transaction{}
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/transaction_accounts/%v/transactions?per_page=100", accountID),
	}
	for {
		batch := []Transaction{}
		resp, err := ps.request(cr, &batch)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, batch...)
		next := getHeaderLink(resp.Header, "next")
		if next == "" {
			break
		}
		cr.path = strings.Replace(next, endpoint, "", -1)
	}
	return transactions, nil
}

// ListCategoriesForAuthedUser returns a list of categories for the authed user
// associated with the given token.
func (ps *PocketSmith) ListCategoriesForAuthedUser() ([]Category, error) {
	return ps.ListCategoriesForUser(ps.user.ID)
}

// ListCategoriesForUser returns a list of categories for a given user id.
// https://developers.pocketsmith.com/reference#get_users-id-categories
func (ps *PocketSmith) ListCategoriesForUser(userID int) ([]Category, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/categories", userID),
	}
	categories := []Category{}
	_, err := ps.request(cr, &categories)
	return categories, err
}

// GetCategoryForAuthedUserByName returns a category for a given category name for
// the authed user associated with the given token.
func (ps *PocketSmith) GetCategoryForAuthedUserByName(category string) (Category, error) {
	return ps.GetCategoryForUserByName(ps.user.ID, category)
}

// GetCategoryForUserByName returns a category for a given category name for the given user.
func (ps *PocketSmith) GetCategoryForUserByName(userID int, category string) (Category, error) {
	categories, err := ps.ListCategoriesForUser(userID)
	if err != nil {
		return Category{}, Error{Message: fmt.Sprintf("failed to get categories list for given user: %v", err), Code: defaultCode}
	}
	for _, c := range categories {
		if c.Title == category {
			return c, nil
		}
	}
	return Category{}, Error{Message: fmt.Sprintf("failed to find %s in categories list for given user", category), Code: defaultCode}
}

// UpdateTransaction updates a PocketSmith transaction.
// https://developers.pocketsmith.com/reference#put_transactions-id
func (ps *PocketSmith) UpdateTransaction(options UpdateTransactionOptions) (Transaction, error) {
	if options.ID == 0 {
		return Transaction{}, Error{Message: "missing transaction id", Code: defaultCode}
	}
	cr := clientRequest{
		method: http.MethodPut,
		path:   fmt.Sprintf("/transactions/%v", options.ID),
		data:   options,
	}
	transaction := Transaction{}
	_, err := ps.request(cr, &transaction)
	return transaction, err
}

// ListAttachmentsForAuthedUser lists the attachments for the authed user
// associated with the given token.
func (ps *PocketSmith) ListAttachmentsForAuthedUser(options ListAttachmentsForUserOptions) ([]Attachment, error) {
	return ps.ListAttachmentsForUser(ps.user.ID, options)
}

// ListAttachmentsForUser lists the attachments for a given user id.
// https://developers.pocketsmith.com/reference#get_users-id-attachments
func (ps *PocketSmith) ListAttachmentsForUser(userID int, options ListAttachmentsForUserOptions) ([]Attachment, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/attachments", userID),
		data:   options,
	}
	attachments := []Attachment{}
	_, err := ps.request(cr, &attachments)
	return attachments, err
}

// CreateAttachmentsForAuthedUser creates an attachment for the authed user
// associated with the given token.
func (ps *PocketSmith) CreateAttachmentsForAuthedUser(options CreateAttachmentsForUserOptions) (*Attachment, error) {
	return ps.CreateAttachmentsForUser(ps.user.ID, options)
}

// CreateAttachmentsForUser creates an attachment for a given user id.
// https://developers.pocketsmith.com/reference#post_users-id-attachments
func (ps *PocketSmith) CreateAttachmentsForUser(userID int, options CreateAttachmentsForUserOptions) (*Attachment, error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/attachments", userID),
		data:   options,
	}
	attachment := &Attachment{}
	_, err := ps.request(cr, &attachment)
	return attachment, err
}

// DeleteAttachment deletes a PocketSmith attachment.
// https://developers.pocketsmith.com/reference#delete_attachments-id
func (ps *PocketSmith) DeleteAttachment(attachmentID int) error {
	if attachmentID == 0 {
		return Error{Message: "missing attachment id", Code: defaultCode}
	}
	cr := clientRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/attachments/%v", attachmentID),
	}
	_, err := ps.request(cr, nil)
	return err
}

// AssignAttachmentToTransaction assigns a given attachment to a transaction.
// https://developers.pocketsmith.com/reference#post_transactions-id-attachments
func (ps *PocketSmith) AssignAttachmentToTransaction(options AssignAttachmentToTransactionOptions) (Attachment, error) {
	if options.TransactionID == 0 {
		return Attachment{}, fmt.Errorf("missing transaction id")
	}
	if options.AttachmentID == 0 {
		return Attachment{}, fmt.Errorf("missing attachment id")
	}
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%v/attachments", options.TransactionID),
		data:   options,
	}
	attachment := Attachment{}
	_, err := ps.request(cr, attachment)
	return attachment, err
}

// CreateInstitutionForAuthedUser creates an institution in PocketSmith
// for the authed user associated with the given token.
func (ps *PocketSmith) CreateInstitutionForAuthedUser(options CreateInstitutionForAuthedUserOptions) (*Institution, error) {
	return ps.CreateInstitutionForUser(ps.user.ID, options)
}

// CreateInstitutionForUser creates an institution in PocketSmith
// for the given user.
// https://developers.pocketsmith.com/reference#post_users-id-institutions
func (ps *PocketSmith) CreateInstitutionForUser(userID int, options CreateInstitutionForAuthedUserOptions) (*Institution, error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/institutions", userID),
		data:   options,
	}
	institution := &Institution{}
	_, err := ps.request(cr, institution)
	return institution, err
}

// generateJSONRequestData returns the body and content type for a JSON request.
func (ps *PocketSmith) generateJSONRequestData(cr clientRequest) ([]byte, string, error) {
	body, err := json.Marshal(cr.data)
	if err != nil {
		return nil, "", Error{fmt.Sprintf("Failed to marshal JSON: %s", err), defaultCode}
	}
	return body, "application/json", nil
}

// envelope defines errors returned from the PocketSmith API.
type envelope struct {
	Error string `json:"error"`
}

// request makes a request to PocketSmith's API.
func (ps *PocketSmith) request(cr clientRequest, result interface{}) (*http.Response, error) {

	// Get the body and content type.
	body, contentType, err := ps.generateJSONRequestData(cr)
	if err != nil {
		return nil, err
	}

	// Construct the request.
	req, err := http.NewRequest(cr.method, endpoint+cr.path, bytes.NewReader(body))
	if err != nil {
		return nil, Error{fmt.Sprintf("Failed to construct request: %s", err), defaultCode}
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Set("X-Developer-Key", ps.token)

	// Make the request.
	resp, err := ps.httpClient.Do(req)
	if err != nil {
		return nil, Error{fmt.Sprintf("Failed to make request: %s", err), defaultCode}
	}
	defer resp.Body.Close()

	// Parse the response.
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, Error{fmt.Sprintf("Failed to read response: %d %s", resp.StatusCode, err), defaultCode}
	}
	// fmt.Println(string(b))
	switch {
	case http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices:
		if len(b) == 0 {
			return resp, nil
		}
		return resp, json.Unmarshal(b, &result)
	}
	var env envelope
	if err := json.Unmarshal(b, &env); err != nil {
		return nil, Error{fmt.Sprintf("Failed to parse error response: %d %s ", resp.StatusCode, err), defaultCode}
	}
	return nil, missingErr{resp.StatusCode, env.Error}
}

// getHeaderLink returns the url of a rel in the Link header.
func getHeaderLink(header http.Header, rel string) string {
	var rgx = regexp.MustCompile(`<(.+?)>;\s*rel="(.+?)"`)
	for _, link := range header["Link"] {
		for _, m := range rgx.FindAllStringSubmatch(link, -1) {
			if len(m) != 3 {
				continue
			}
			if m[2] == rel {
				return m[1]
			}
		}
	}
	return ""
}

type missingErr struct {
	code int
	body string
}

func (e missingErr) Error() string {
	return fmt.Sprintf("%v %s", e.code, e.body)
}

// IsMissing is a helper function to help determine if the error returned is a missing error.
func IsMissing(err error) bool {
	_, ok := err.(missingErr)
	return ok
}
