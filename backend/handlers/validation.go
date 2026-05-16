package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"project-manager/helper"
)

// validate プロセス全体で共有するバリデータ。
// json タグ名をフィールド名として使うように初期化する。
var validate = func() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}()

// Bind リクエストボディの JSON デコードと validate タグに基づく検証を行う。
// 失敗時はクライアントへ 400 を書き込み false を返す。呼び出し側は false なら即 return。
func Bind(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := DecodeJSON(r, dst); err != nil {
		slog.DebugContext(r.Context(), "failed to decode body", "error", err)
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	if err := validate.Struct(dst); err != nil {
		slog.DebugContext(r.Context(), "validation failed", "error", err)
		helper.ErrorResponse(w, http.StatusBadRequest, formatValidationError(err, dst))
		return false
	}
	return true
}

// formatValidationError validator.ValidationErrors を整形してクライアントに返す文字列にする。
// dst は対象 struct。Param に含まれる Go フィールド名を JSON タグ名へ解決するために使う。
func formatValidationError(err error, dst any) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return "invalid request"
	}
	fieldMap := buildFieldNameMap(dst)
	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		msgs = append(msgs, fmt.Sprintf("%s: %s", fe.Field(), describeRule(fe, fieldMap)))
	}
	return strings.Join(msgs, "; ")
}

// describeRule validator のタグから人間が読めるルール説明を返す。
// gtefield / gtfield など他フィールドを参照するルールでは fieldMap で JSON タグ名に変換する。
func describeRule(fe validator.FieldError, fieldMap map[string]string) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", fe.Param())
	case "gtefield":
		return fmt.Sprintf("must be greater than or equal to %s", lookupJSONName(fe.Param(), fieldMap))
	case "gtfield":
		return fmt.Sprintf("must be greater than %s", lookupJSONName(fe.Param(), fieldMap))
	case "ltefield":
		return fmt.Sprintf("must be less than or equal to %s", lookupJSONName(fe.Param(), fieldMap))
	case "ltfield":
		return fmt.Sprintf("must be less than %s", lookupJSONName(fe.Param(), fieldMap))
	default:
		return fmt.Sprintf("failed on %q", fe.Tag())
	}
}

// buildFieldNameMap struct の Go フィールド名から JSON タグ名へのマップを作る。
// ネストした struct は対象外（必要になったら拡張する）。
func buildFieldNameMap(dst any) map[string]string {
	t := reflect.TypeOf(dst)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	m := make(map[string]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name != "" && name != "-" {
			m[fld.Name] = name
		}
	}
	return m
}

// lookupJSONName Go フィールド名を JSON タグ名へ変換。マップにない場合は元の名前を返す。
func lookupJSONName(goName string, m map[string]string) string {
	if v, ok := m[goName]; ok {
		return v
	}
	return goName
}
