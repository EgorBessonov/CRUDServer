package service

import (
	"CRUDServer/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	log "github.com/sirupsen/logrus"
)

func Registration(form repository.RegistrationForm, rps repository.Repository, ctx context.Context) error {
	hPassword, err := hashPassword(form.Password)
	if err != nil {
		servicesOperationError(err, "Registration()")
		return err
	}
	form.Password = hPassword
	err = rps.CreateAuthUser(ctx, form)
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	if password == "" {
		log.WithFields(log.Fields{
			"method": "hashPassword()",
			"err":    errors.New("no input supplied"),
		}).Info("Authentication info.")
		return "", errors.New("no input supplied")
	}
	h := sha256.New()
	h.Write([]byte(password))
	hashedPassword := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hashedPassword, nil
}

func servicesOperationError(err error, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"err":    err,
		"status": "operation failed.",
	}).Info("Services info.")
}
