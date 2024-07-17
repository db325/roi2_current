package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// This is the account struct. There are 3 roles: OIC, Admin, and Drive Staff. Only OIC and Admin can create accounts. The Admin
// can create accounts for OIC. The OIC can create accounts for Drive Staff. The  Admin can view all accounts.

const (
	OIC    = "Officer in Charge"
	Admin  = "Administrator"
	DStaff = "Drive Staff"
)

type Account struct {
	mu         *sync.RWMutex
	ID         string
	Name       string
	Email      string
	Role       string
	HashedPwd  string
	pwd        string
	Center     int //center is the center number, and the only one they can access
	CenterName string
}

// EncodePassword hashes a password using bcrypt.
func EncodePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies if a password matches its hashed counterpart.

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (acc *Account) SetPassword(password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	acc.HashedPwd = string(hashedPassword)
	return nil
}

func (acc *Account) GetPassword() string {

	return acc.pwd
}

// GenerateRandomID generates a random ID based on email and name properties
func GenerateRandomID(email, name string) string {
	// Concatenate email and name
	data := email + name

	// Hash the concatenated data using MD5
	hash := md5.Sum([]byte(data))

	// Encode the hashed data to a hexadecimal string
	id := hex.EncodeToString(hash[:])

	// Return the ID
	return id
}

// function to encode a string to hex
func EncodeStringToHex(str string) string {
	encoded := hex.EncodeToString([]byte(str))
	return encoded
}

// function to decode from hex to string
func DecodeHexToString(hexString string) string {
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("error decoding hex to string", err)
	}
	return string(decoded)
}

var Administrator = CreateAdminAccount()

func CreateAdminAccount() Account {
	adminPwd := "SharkJump1234!!"
	adminEmail := "db.coder.523@gmail.com"
	adminName := "DaJovan B. Coder"
	adminCenterName := "All"

	// Generate a random ID for the admin account
	adminID := GenerateRandomID(adminEmail, adminName)

	// Hash the admin password
	adminHashedPwd, err := EncodePassword(adminPwd)
	if err != nil {
		fmt.Println("error encoding password", err)
	}

	// Create the admin account with the specified properties
	adminAccount := Account{
		ID:         adminID,
		Name:       adminName,
		Email:      adminEmail,
		Role:       Admin,
		HashedPwd:  adminHashedPwd,
		pwd:        adminPwd,
		Center:     1000, // Assuming center 0 for admin
		CenterName: adminCenterName,
	}

	return adminAccount
}

// CreateAccount creates an account with the given properties
func CreateAccount(role string, center int, pwd string, email string, name string) Account {
	// Generate a random ID
	id := GenerateRandomID(email, name)

	// Check if the role is "Admin"
	if role == Admin {
		// If the role is "Admin", set it to "Drive Staff" instead
		role = DStaff
	}
	n := strings.Split(name, " ")[0]
	p, e := EncodePassword(pwd)

	if e != nil {
		fmt.Println("error encoding password", e)
	}
	cn := c.Centers[center].Name

	return Account{ID: id, Name: n, Email: email, Role: role, Center: center, pwd: pwd, HashedPwd: p, CenterName: cn}
}

// AccountStore is a store for accounts
type AccountStore struct {
	accounts []Account
	mutex    sync.Mutex
}

// AddAccount adds an account to the store
func (as *AccountStore) AddAccount(account Account) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	// Check if the store already has an account with the same role and center
	for _, a := range as.accounts {
		if a.Role == account.Role && a.Center == account.Center {
			fmt.Println("Account already exists", a.Role, a.Center, account.Role, account.Center)
			return
		}
	}

	fmt.Println("Account does not exist, adding now", account.Role, account.Center, account.Name, account.Email, account.pwd)
	as.accounts = append(as.accounts, account)
}

// GetAccounts returns all accounts in the store
func (as *AccountStore) GetAccounts() []Account {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	return as.accounts
}

func (as *AccountStore) GetAccountsByCenter(center int) []Account {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	var accounts []Account
	for _, a := range as.accounts {
		if a.Center == center {
			accounts = append(accounts, a)
		}
	}
	return accounts
}

// GetAccountByRole returns an account by role
func (as *AccountStore) GetAccountByRole(role string) *Account {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for _, account := range as.accounts {
		if account.Role == role {
			return &account
		}
	}
	return nil
}

// GetAccountByID returns an account by ID
func (as *AccountStore) GetAccountByID(id string) *Account {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for _, account := range as.accounts {
		if account.ID == id {
			return &account
		}
	}
	return nil
}

// GetAccountByEmail returns an account by email
func (as *AccountStore) GetAccountByEmail(email string) *Account {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for _, account := range as.accounts {
		if account.Email == email {
			return &account
		}
	}
	return nil
}

// GetAccountByName returns an account by name
func (as *AccountStore) GetAccountByName(name string) *Account {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for _, account := range as.accounts {
		if account.Name == name {
			return &account
		}
	}
	return nil
}

// DeleteAccountByID deletes an account by ID
func (as *AccountStore) DeleteAccountByID(id string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for i, account := range as.accounts {
		if account.ID == id {
			as.accounts = append(as.accounts[:i], as.accounts[i+1:]...)
			break
		}
	}
}

// DeleteAccountByEmail deletes an account by email
func (as *AccountStore) DeleteAccountByEmail(email string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for i, account := range as.accounts {
		if account.Email == email {
			as.accounts = append(as.accounts[:i], as.accounts[i+1:]...)
			break
		}
	}
}

// DeleteAccountByName deletes an account by name
func (as *AccountStore) DeleteAccountByName(name string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for i, account := range as.accounts {
		if account.Name == name {
			as.accounts = append(as.accounts[:i], as.accounts[i+1:]...)
			break
		}
	}
}

// UpdateAccount updates an account
func (as *AccountStore) UpdateAccount(account Account) {
	as.mutex.Lock()

	for i, a := range as.accounts {
		if a.ID == account.ID {
			as.accounts[i] = account
			as.mutex.Unlock()
			break
		}
	}
}

// AccountExistsByID checks if an account exists by ID
func (as *AccountStore) AccountExistsByID(id string) bool {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for _, account := range as.accounts {
		if account.ID == id {
			return true
		}
	}
	return false
}

// AccountExistsByEmail checks if an account exists by email
func (as *AccountStore) AccountExistsByEmail(email string) bool {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for _, account := range as.accounts {
		if account.Email == email {
			return true
		}
	}
	return false
}

// AccountExistsByName checks if an account exists by name
func (as *AccountStore) AccountExistsByName(name string) bool {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	for _, account := range as.accounts {
		if account.Name == name {
			return true
		}
	}
	return false
}

// create a func that generates dummy accounts for testing. it should create 2 accounts for each center. one OIC and one DStaff
func GenerateDummyAccounts() {
	// Use the global account store

	// Create accounts for each center
	for i := 0; i < len(c.Centers); i++ {
		// Create an OIC account
		oicName := fmt.Sprintf("oic-%d", i+1)

		//use GenerateRandomName to create a random name for the password value
		oic_random_name := GenerateRandomName()
		fmt.Println("oic_random_name used for password", oic_random_name)

		oic := CreateAccount(OIC, i, oic_random_name, "oictest", oicName)
		fmt.Println("oic created: ", oic.pwd, oic.HashedPwd)
		// Add the OIC account to the store
		AS.AddAccount(oic)
		updateUsers(oic.Center, oic)
		// Create a Drive Staff account
		dstaffName := fmt.Sprintf("dstaff-%d", i+1)
		newname := strings.Split(dstaffName, " ")[0]

		ds_random_name := GenerateRandomName()
		//fmt.Println("ds_random_name used for password", ds_random_name)
		dstaff := CreateAccount(DStaff, i, ds_random_name, "dstafftest", newname)
		fmt.Println("dstaff created: ", dstaff.pwd, dstaff.HashedPwd)
		// Add the Drive Staff account to the store
		AS.AddAccount(dstaff)
		updateUsers(dstaff.Center, dstaff)
	}
}

// List of 150 animals
var animals = []string{
	"Lion", "Tiger", "Elephant", "Giraffe", "Zebra",
	"Gorilla", "Hippopotamus", "Rhinoceros", "Cheetah", "Leopard",
	"Jaguar", "Polar Bear", "Grizzly Bear", "Panda", "Koala",
	"Kangaroo", "Wallaby", "Platypus", "Tasmanian Devil", "Chimpanzee",
	"Orangutan", "Bonobo", "Baboon", "Mandrill", "Proboscis Monkey",
	"Rhesus Monkey", "Capuchin Monkey", "Spider Monkey", "Howler Monkey", "Squirrel Monkey",
	"Lemur", "Marmoset", "Tamarin", "Gibbon", "Sloth",
	"Anteater", "Armadillo", "Pangolin", "Aardvark", "Hyena",
	"Wolf", "Fox", "Coyote", "Jackal", "Dingo",
	"Dog", "Cat", "Lynx", "Bobcat", "Cougar",
	"Panther", "Puma", "Caracal", "Ocelot", "Serval",
	"Snow Leopard", "Clouded Leopard", "Asian Leopard", "African Leopard", "Persian Leopard",
	"Amur Leopard", "Jaguarundi", "Maned Wolf", "Red Fox", "Gray Fox",
	"Fennec Fox", "Swift Fox", "Arctic Fox", "Raccoon", "Red Panda",
	"Fossa", "Civet", "Binturong", "Tiger Quoll", "Spotted Hyena",
	"Striped Hyena", "Brown Hyena", "African Wild Dog", "Gray Wolf", "Red Wolf",
	"Ethiopian Wolf", "Coywolf", "Golden Jackal", "Black-Backed Jackal", "Side-Striped Jackal",
	"European Jackal", "Domestic Cat", "Lioness", "Tigress", "Leopardess",
	"Jaguaress", "Lion Cub", "Tiger Cub", "Leopard Cub", "Jaguar Cub",
	"Lion Pride", "Tiger Stripe", "Leopard Spot", "Jaguar Print", "Lion Roar",
	"Tiger Growl", "Leopard Hiss", "Jaguar Snarl", "Wolf Pack", "Fox Den",
	"Coyote Pack", "Jackal Pack", "Dingo Pack", "Dog Pack", "Cat Clan",
	"Pride of Lions", "Tiger Family", "Leopard Group", "Jaguar Troop", "Elephant Herd",
	"Giraffe Tower", "Zebra Crossing", "Hippo Pod", "Rhinoceros Crash", "Cheetah Coalition",
	"Bear Family", "Panda Group", "Koala Mob", "Kangaroo Mob", "Platypus Paddle",
	"Tasmanian Devil Pack", "Chimpanzee Troop", "Bonobo Group", "Baboon Troop", "Monkey Troop",
	"Lemur Troop", "Sloth Bed", "Anteater Parade", "Armadillo Roll", "Pangolin Puddle",
	"Aardvark Army",
}

// List of 150 verbs
var verbs = []string{
	"accept", "achieve", "act", "add", "adjust",
	"admire", "advise", "agree", "announce", "apologize",
	"appear", "applaud", "appreciate", "approve", "argue",
	"arrive", "ask", "assist", "attach", "attend",
	"attract", "avoid", "bake", "balance", "ban",
	"bathe", "be", "bear", "beat", "become",
	"beg", "behave", "believe", "belong", "bend",
	"bite", "blame", "bless", "blink", "blot",
	"blow", "boast", "boil", "borrow", "bounce",
	"bow", "brake", "brush", "bubble", "build",
	"burn", "bury", "calculate", "call", "camp",
	"care", "carry", "carve", "cause", "challenge",
	"change", "charge", "chase", "cheat", "check",
	"cheer", "chew", "choose", "chop", "claim",
	"clap", "clean", "clear", "climb", "clip",
	"close", "collect", "color", "comb", "comfort",
	"command", "communicate", "compare", "compete", "complain",
	"complete", "concentrate", "concern", "confess", "confuse",
	"connect", "consider", "consist", "contain", "continue",
	"copy", "correct", "cough", "count", "cover",
	"crack", "crash", "crawl", "cross", "crush",
	"cry", "cure", "curl", "curve", "cut",
	"damage", "dance", "dare", "deal", "decay",
	"deceive", "decide", "decorate", "delay", "delight",
	"deliver", "depend", "describe", "desert", "deserve",
	"destroy", "develop", "die", "disagree", "disappear",
	"discover", "dislike", "divide", "do", "donate",
	"draw", "dream", "drink", "drive", "drop",
	"dry", "dust", "earn", "eat", "educate",
	"embarrass", "embrace", "employ", "empty", "encourage",
	"end", "enjoy", "enter", "entertain", "escape",
	"evaporate", "examine", "excite", "exercise", "exist",
	"expand", "expect", "explain", "explode", "explore",
	"extend", "face",
}

// GenerateRandomName generates a random name based on the animals and verbs lists. It concatenates a random animal and a random verb. It also adds a random number to the end of the name.
func GenerateRandomName() string {
	// Generate a random number
	num := rand.Intn(1000)

	// Concatenate a random animal and a random verb
	name := fmt.Sprintf("%s%s%d", animals[rand.Intn(len(animals))], verbs[rand.Intn(len(verbs))], num)

	// Remove spaces from the name
	name = strings.ReplaceAll(name, " ", "")

	// Convert the name to title case
	name = strings.Title(name)

	// Add random non-letter characters
	name = addRandomNonLetterCharacters(name)

	// Truncate the name to 16 characters
	if len(name) > 16 {
		name = name[:16]
	}

	// Return the name
	return name
}

func addRandomNonLetterCharacters(name string) string {
	nonLetterChars := []string{"!", "@", "#", "$", "%", "&", "*", "?"}

	// Add random non-letter characters to the name
	for i := 0; i < 3; i++ {
		randomChar := nonLetterChars[rand.Intn(len(nonLetterChars))]
		randomIndex := rand.Intn(len(name) + 1)
		name = name[:randomIndex] + randomChar + name[randomIndex:]
	}

	return name
}
