package user

type User struct {
	ID             uint   `gorm:"primarykey" json:"id"`
	Name           string `json:"nama"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Alamat         string `json:"alamat"`
	Jenis_Budidaya string `json:"jenis_budidaya"`
	Lokasi_Tambak  string `json:"lokasi_tambak"`
	Luas_Kolam     string `json:"luas_kolam"`
}

