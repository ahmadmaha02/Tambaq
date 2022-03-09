package user

type Service interface {
	FindAll() ([]User, error)
	FindByID(ID int) (User, error)
	Register(user User) (User, error)
	RegisterMember(user User) (User, error)
	Login() ([]User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) Register(input PostRegisterBody) (User, error) {
	user := User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       input.Password,
		Alamat:         input.Alamat,
		Jenis_Budidaya: input.Jenis_Budidaya,
		Lokasi_Tambak:  input.Lokasi_Tambak,
		Luas_Kolam:     input.Luas_Kolam,
	}

	registerdUser, err := s.repository.Register(user)
	if err != nil {
		return registerdUser, err
	}
	return registerdUser, nil
}

func (s *service) FindAll() ([]User, error) {
	users, err := s.repository.FindAll()
	return users, err
}

