package utils

import (
	"context"
	"net/http"

	"github.com/fengzhongzhu1621/xgo/db/db"
	"github.com/fengzhongzhu1621/xgo/network/constant"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/rs/xid"
)

// NewHeader take out the required header value and create a new header
func NewHeader(header http.Header) http.Header {
	newHeader := http.Header{}
	newHeader.Set("Content-Type", "application/json")

	nethttp.SetUser(newHeader, nethttp.GetUser(header))

	nethttp.SetSupplierAccount(newHeader, nethttp.GetSupplierAccount(header))

	nethttp.SetRid(newHeader, nethttp.GetRid(header))

	nethttp.SetLanguage(newHeader, nethttp.GetLanguage(header))

	nethttp.SetAppCode(newHeader, nethttp.GetAppCode(header))

	db.SetTXId(newHeader, db.GetTXId(header))

	db.SetTXTimeout(newHeader, db.GetTXTimeout(header))

	if nethttp.IsReqFromWeb(header) {
		nethttp.SetReqFromWeb(newHeader)
	}

	return newHeader
}

func NewHeaderFromContext(ctx context.Context) http.Header {
	rid := ctx.Value(constant.ContextRequestIDField)
	ridValue, _ := rid.(string)

	user := ctx.Value(constant.ContextRequestUserField)
	userValue, _ := user.(string)

	owner := ctx.Value(constant.ContextRequestOwnerField)
	ownerValue, _ := owner.(string)

	return GenCommonHeader(userValue, ownerValue, ridValue)
}

func GenCommonHeader(user, supplierAccount, rid string) http.Header {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	if user == "" {
		user = constant.SystemOperatorUserName
	}

	if supplierAccount == "" {
		supplierAccount = constant.DefaultOwnerID
	}

	if rid == "" {
		rid = xid.New().String()
	}

	nethttp.SetUser(header, user)
	nethttp.SetSupplierAccount(header, supplierAccount)
	nethttp.SetRid(header, rid)
	return header
}

func GenDefaultHeader() http.Header {
	return GenCommonHeader("", "", "")
}

// BuildHeader build header by user & supplier account
func BuildHeader(user string, supplierAccount string) http.Header {
	return GenCommonHeader(user, supplierAccount, "")
}

func Header(header http.Header) http.Header {
	newHeader := make(http.Header)

	nethttp.SetRid(newHeader, nethttp.GetRid(header))
	nethttp.SetUser(newHeader, nethttp.GetUser(header))
	nethttp.SetUserToken(newHeader, nethttp.GetUserToken(header))
	nethttp.SetUserTicket(newHeader, nethttp.GetUserTicket(header))
	nethttp.SetLanguage(newHeader, nethttp.GetLanguage(header))
	nethttp.SetSupplierAccount(newHeader, nethttp.GetSupplierAccount(header))
	nethttp.SetAppCode(newHeader, nethttp.GetAppCode(header))
	nethttp.SetReqRealIP(newHeader, nethttp.GetReqRealIP(header))
	if nethttp.IsReqFromWeb(header) {
		nethttp.SetReqFromWeb(newHeader)
	}
	newHeader.Add(constant.ReadReferenceKey, header.Get(constant.ReadReferenceKey))

	return newHeader
}
