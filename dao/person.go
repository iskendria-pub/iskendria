package dao

type Person struct {
	Id              string
	CreatedOn       int64
	ModifiedOn      int64
	PublicKey       string
	Name            string
	Email           string
	IsMajor         bool
	IsSigned        bool
	Balance         int32
	BiographyHash   string
	BiographyFormat string
	Organization    string
	Telephone       string
	Address         string
	PostalCode      string
	Country         string
	ExtraInfo       string
}

func SearchPersonByKey(key string) ([]*Person, error) {
	// TODO: Implement
	return nil, nil
}

func GetPersonById(id string) (*Person, error) {
	// TODO: Implement
	return nil, nil
}

type PersonUpdate struct {
	PublicKey     string
	Name          string
	Email         string
	BiographyHash string
	Organization  string
	Telephone     string
	Address       string
	PostalCode    string
	Country       string
	ExtraInfo     string
}

func PersonToPersonUpdate(p *Person) *PersonUpdate {
	return &PersonUpdate{
		PublicKey:     p.PublicKey,
		Name:          p.Name,
		Email:         p.Email,
		BiographyHash: p.BiographyHash,
		Organization:  p.Organization,
		Telephone:     p.Telephone,
		Address:       p.Address,
		PostalCode:    p.PostalCode,
		Country:       p.Country,
		ExtraInfo:     p.ExtraInfo,
	}
}
