package repository


type PostgreRepository struct{}

func(rps PostgreRepository) CreateUser(u User)error{

	return nil
}

func(rps PostgreRepository) ReadUser(u string) (*User, error){
	return &User{}, nil
}

func(rps PostgreRepository) UpdateUser(u User)error{
	return nil
}

func(rps PostgreRepository) DeleteUser(userID string)error{
	return nil
}

func(rps PostgreRepository) AddImage(){

}

func(rps PostgreRepository) GetImage(){

}