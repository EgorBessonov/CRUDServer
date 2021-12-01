package repository

import "github.com/google/uuid"


type PostgreRepository struct{}

func(rps PostgreRepository) CreateUser(u User)error{

	return nil
}

func(rps PostgreRepository) ReadUser(){

}

func(rps PostgreRepository) UpdateUser(u User)error{
	return nil
}

func(rps PostgreRepository) DeleteUser(userID uuid.UUID)error{
	return nil
}

func(rps PostgreRepository) AddImage(){

}

func(rps PostgreRepository) GetImage(){

}