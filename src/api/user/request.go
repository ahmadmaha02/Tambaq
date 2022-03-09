package user

type PostRegisterBody struct {
	Name           string `json:"nama" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Alamat         string `json:"alamat" binding:"required"`
	Jenis_Budidaya string `json:"jenis_budidaya" binding:"required"`
	Lokasi_Tambak  string `json:"lokasi_tambak" binding:"required"`
	Luas_Kolam     string `json:"luas_kolam" binding:"required"`
}

type PostLoginBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}