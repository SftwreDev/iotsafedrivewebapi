package structs

type TrustedContacts struct {
	ID      uint   `json:"id"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Contact string `json:"contact" validate:"required"`
}

type AllTrustedContacts struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
	Owner   string `json:"owner"`
}
