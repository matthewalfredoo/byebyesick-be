package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/util"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func WrapError(err error, customCode ...int) error {
	errWrapper := &apperror.Wrapper{}

	if ok := errors.As(err, &errWrapper); !ok {
		errWrapper.ErrorStored = err
	}

	if len(customCode) > 0 {
		errWrapper.Code = customCode[0]
		return errWrapper
	}

	var (
		errJsonSyntax     *json.SyntaxError
		errJsonUnmarshall *json.UnmarshalTypeError
		errTimeParse      *time.ParseError
		errValidation     validator.ValidationErrors
		errNotFound       *apperror.NotFound
		errAlreadyExist   *apperror.AlreadyExist
		errNotMatch       *apperror.NotMatch
		errForbidden      *apperror.Forbidden
		errAuth           *apperror.AuthError
	)

	switch {
	case errors.Is(errWrapper.ErrorStored, context.DeadlineExceeded):
		errWrapper.Code = http.StatusGatewayTimeout
		errWrapper.Message = "server timeout, took too long to process request"

	case errors.Is(errWrapper.ErrorStored, io.EOF):
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, io.ErrUnexpectedEOF):
		fallthrough

	case errors.As(errWrapper.ErrorStored, &errJsonSyntax):
		fallthrough

	case errors.As(errWrapper.ErrorStored, &errJsonUnmarshall):
		errWrapper.ErrorStored = fmt.Errorf("invalid JSON syntax or format")
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errTimeParse):
		errWrapper.Code = http.StatusBadRequest

	case errWrapper.ErrorStored.Error() == "invalid request":
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, apperror.ErrForgotPasswordTokenInvalid), errors.Is(errWrapper.ErrorStored, apperror.ErrForgotPasswordTokenExpired):
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, apperror.ErrRegisterTokenInvalid), errors.Is(errWrapper.ErrorStored, apperror.ErrRegisterTokenExpired):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errValidation):
		errWrapper.Code = http.StatusBadRequest
		errWrapper.Message = handleErrValidation(errValidation)

	case errors.As(errWrapper.ErrorStored, &errForbidden):
		errWrapper.Code = http.StatusForbidden

	case errors.As(errWrapper.ErrorStored, &errNotFound):
		errWrapper.Code = http.StatusNotFound

	case errors.As(errWrapper.ErrorStored, &errAlreadyExist):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errNotMatch):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errAuth):
		errWrapper.Code = http.StatusUnauthorized

	case errors.Is(err, apperror.ErrUnauthorized):
		errWrapper.Code = http.StatusUnauthorized

	case errors.Is(errWrapper.ErrorStored, apperror.ErrForbiddenViewEntity):
		errWrapper.Code = http.StatusForbidden

	case errors.Is(errWrapper.ErrorStored, apperror.ErrForbiddenModifyEntity):
		errWrapper.Code = http.StatusForbidden

	case errors.Is(errWrapper.ErrorStored, apperror.ErrDeleteAlreadyAssignedAdmin):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrInvalidDecimal):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrPharmacyProductUniqueConstraint):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrProductCategoryUniqueConstraint):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrProductCategoryStillUsedByProducts):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrProductUniqueConstraint):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrInvalidLatLong):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrInvalidRegisterRole):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrRecordNotFound):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrWrongCredentials):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrPaymentSent):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrPaymentConfirmed):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrInvalidCityProvinceCombi):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrPasswordTooLong):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrStartDateAfterEndDate):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrBadConfirmStatus):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrInsufficientProductStock):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrProductImageDoesNotExistInContext):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrProductStockNotEnoughToAddToCart):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrProductAddedToCartMustHaveAtLeastOne):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrMainAddressNotFound):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrRequestStockMutationFromOwnPharmacy):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrRequestStockMutationDifferentProduct):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrAlreadyFinishedRequest):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrGetShipmentMethodNoItems):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrGetShipmentMethodDifferentPharmacy):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrBadTransactionCancelStatus):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrPaymentNotSent):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrBadRejectStatus):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrBadShipStatus):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrBadReceiveStatus):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrInvalidIntInString):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrChatStillOngoing):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrConsultationSessionAlreadyHasPrescription):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrConsultationSessionPrescriptionMustExistBeforeIssuingSickLeave):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrSickLeaveStartingDateShouldBeBeforeEndingDate):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrConsultationSessionAlreadyHasSickLeaveForm):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrChatAlreadyEnded):
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, apperror.ErrNoPharmacyToStockTransfer):
		fallthrough

	case errors.Is(errWrapper.ErrorStored, apperror.ErrGetShipmentCost):
		errWrapper.Code = http.StatusServiceUnavailable

	default:
		errWrapper.Code = http.StatusInternalServerError
	}

	return errWrapper
}

func handleErrValidation(ve validator.ValidationErrors) string {
	buff := bytes.NewBufferString("")

	for i := range ve {
		buff.WriteString(createErrValidationMsgTag(ve[i]))
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func createErrValidationMsgTag(fieldError validator.FieldError) string {
	fieldName := util.PascalToSnake(fieldError.Field())
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' is required", fieldName)
	case "email":
		return fmt.Sprintf("field '%s' must be in the format of an email", fieldName)
	case "datetime":
		return fmt.Sprintf("field '%s' must follow the format %s", fieldName, fieldError.Param())
	case "number":
		return fmt.Sprintf("field '%s' must be a number", fieldName)
	case "numbergt":
		return fmt.Sprintf("field '%s' must be a number greater than %s", fieldName, fieldError.Param())
	case "numeric":
		return fmt.Sprintf("field '%s' must be numeric", fieldName)
	case "numericgt":
		return fmt.Sprintf("field '%s' must be a numeric greater than %s", fieldName, fieldError.Param())
	case "len":
		return fmt.Sprintf("field '%s' must have exactly %s characters long", fieldName, fieldError.Param())
	case "min":
		switch fieldError.Type().Kind() {
		case reflect.String:
			return fmt.Sprintf("field '%s' must be at least %s characters long", fieldName, fieldError.Param())
		case reflect.Slice:
			return fmt.Sprintf("field '%s' must have at least %s item", fieldName, fieldError.Param())
		default:
			return fmt.Sprintf("field '%s' have a minimum value of %s", fieldName, fieldError.Param())
		}
	case "max":
		switch fieldError.Type().Kind() {
		case reflect.String:
			return fmt.Sprintf("field '%s' must be at maximum %s characters long", fieldName, fieldError.Param())
		case reflect.Slice:
			return fmt.Sprintf("field '%s' must have at maximum %s item", fieldName, fieldError.Param())
		default:
			return fmt.Sprintf("field '%s' have maximum value of %s", fieldName, fieldError.Param())
		}
	case "oneof":
		params := strings.ReplaceAll(fieldError.Param(), " ", ", ")
		return fmt.Sprintf("item '%v' on field '%s' must be one of %s", fieldError.Value(), fieldName, params)
	case "latitude":
		return fmt.Sprintf("field '%s' must be a valid latitude", fieldName)
	case "longitude":
		return fmt.Sprintf("field '%s' must be a valid longitude", fieldName)
	case "gtfield":
		otherFieldName := util.PascalToSnake(fieldError.Param())
		return fmt.Sprintf("field '%s' must be greater than '%s' field", fieldName, otherFieldName)
	case "filesize":
		return fmt.Sprintf("field '%s' must have a file size at maximum %v KB", fieldName, fieldError.Param())
	case "filetype":
		param := strings.ReplaceAll(fieldError.Param(), " ", ", ")
		return fmt.Sprintf("field '%s' must have a file with type as one of %s", fieldName, param)
	case "comma_separated":
		return fmt.Sprintf("field '%s' does not contain valid comma separated values", fieldName)
	default:
		msg := fmt.Sprintf("field '%s' failed on validation %s %s", fieldName, fieldError.Tag(), fieldError.Param())
		return strings.TrimSpace(msg)
	}
}
