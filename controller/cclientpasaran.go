package controller

import (
	"context"
	"log"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/apisavetransaksi_go/entities"
	"bitbucket.org/isbtotogroup/apisavetransaksi_go/helpers"
	"bitbucket.org/isbtotogroup/apisavetransaksi_go/model"
	"github.com/buger/jsonparser"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

const fielddomain_redis = "LISTDOMAIN"
const fieldsetting_redis = "LISTSETTING_MASTER"
const fieldallpasaran_redis = "listpasaran_"
const fieldresult_redis = "listresult_"
const fieldconfig_redis = "config_"
const fieldinvoice_redis = "listinvoice_"
const fieldlimit_redis = "limitpasaran_"

type responseaxios struct {
	Status          int    `json:"status"`
	Token           string `json:"token"`
	Member_company  string `json:"member_company"`
	Member_username string `json:"member_username"`
	Member_credit   int    `json:"member_credit"`
}
type responseaxioserror struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SaveTogel(c *fiber.Ctx) error {
	client := new(entities.Controller_clientSaveTogel)
	if err := c.BodyParser(client); err != nil {
		panic(err.Error())
	}
	flag_domain := _domainsecurity(client.Hostname)
	if !flag_domain {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "NOT REGISTER",
			"record":  nil,
		})
	}
	_, _, clientdompet := Fetch_apimoney("qC5YmBvXzabGp34jJlKvnC6wCrr3pLCwBzsLoSzl4k=")

	result, err := model.Savetransaksi(
		client.Client_Username,
		client.Client_Company, client.Idtrxkeluaran, client.Idcomppasaran, client.Devicemember, client.Formipaddress, client.Timezone, clientdompet, client.Totalbayarbet, client.List4d)
	if err != nil {
		// panic(err.Error())
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusAccepted,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredisclient(client.Client_Company, client.Idtrxkeluaran, client.Client_Username, client.Pasarancode, client.Pasaranperiode)
	return c.JSON(result)
}
func Fetch_apimoney(token string) (string, string, int) {
	axios := resty.New()
	resp, err := axios.R().
		SetResult(responseaxios{}).
		SetError(responseaxioserror{}).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"token": token,
		}).
		Post("http://128.199.241.112:6061/api/servicetoken")
	if err != nil {
		log.Println(err.Error())
	}
	result := resp.Result().(*responseaxios)
	return result.Member_company, result.Member_username, result.Member_credit
}
func _domainsecurity(nmdomain string) bool {
	log.Printf("Domain Client : %s", nmdomain)
	resultredis, flag_domain := helpers.GetRedis(fielddomain_redis)
	flag := false
	if len(nmdomain) > 0 {
		if !flag_domain {
			result, err := model.Get_Domain()
			if err != nil {
				flag = false
			}
			log.Println("DOMAIN MYSQL")
			helpers.SetRedis(fielddomain_redis, result, 24*time.Hour)
			flag = true
		} else {
			jsonredis_domain := []byte(resultredis)
			record_RD, _, _, _ := jsonparser.Get(jsonredis_domain, "record")
			jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				domain, _ := jsonparser.GetString(value, "domain")
				if nmdomain == domain {
					flag = true
					log.Println("DOMAIN CACHE")
				}
			})
		}
	}
	return flag
}
func _deleteredisclient(company, idtrxkeluaran, username, pasarancode, pasaranperiode string) {
	field_redis := "listinvoice_" + strings.ToLower(company) + "_" + idtrxkeluaran + "_" + strings.ToLower(username)
	val := helpers.DeleteRedis(field_redis)
	log.Printf("DELETE REDIS INVOICE %d\n", val)
	val_limit := helpers.DeleteRedis(
		fieldlimit_redis + strings.ToLower(company) + "_" +
			strings.ToLower(username) + "_" +
			strings.ToLower(pasarancode) + "_" +
			strings.ToLower(pasaranperiode))
	log.Printf("DELETE REDIS LIMIT %d\n", val_limit)

	//AGEN
	val_agen_dashboard := helpers.DeleteRedis("LISTDASHBOARDPASARAN_AGENT_" + strings.ToLower(company))
	val_agen := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran)
	val_agenlistmember := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTMEMBER")
	val_agenlistbettable := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTBETTABLE")
	val_agen4d := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_4D")
	val_agen3d := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_3D")
	val_agen2d := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_2D")
	val_agen2dd := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_2DD")
	val_agen2dt := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_2DT")
	val_agencolokbebas := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_COLOK_BEBAS")
	val_agencolokmacau := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_COLOK_MACAU")
	val_agencoloknaga := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_COLOK_NAGA")
	val_agencolokjitu := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_COLOK_JITU")
	val_agen5050umum := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_50_50_UMUM")
	val_agen5050special := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_50_50_SPECIAL")
	val_agen5050kombinasi := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_50_50_KOMBINASI")
	val_agenmacaukombinasi := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_MACAU_KOMBINASI")
	val_agendasar := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_DASAR")
	val_agenshio := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTPERMAINAN_SHIO")
	val_agenall := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTBET_all")
	val_agenwinner := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTBET_winner")
	val_agencancel := helpers.DeleteRedis("LISTPERIODE_AGENT_" + strings.ToLower(company) + "_INVOICE_" + idtrxkeluaran + "_LISTBET_cancel")
	log.Printf("DELETE REDIS AGEN DASHBOARD %d\n", val_agen_dashboard)
	log.Printf("DELETE REDIS AGEN PERIODE %d\n", val_agen)
	log.Printf("DELETE REDIS AGEN PERIODE LISTMEMBER %d\n", val_agenlistmember)
	log.Printf("DELETE REDIS AGEN PERIODE LISTBETTABEL %d\n", val_agenlistbettable)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN 4D %d\n", val_agen4d)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN 3D %d\n", val_agen3d)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN 2D %d\n", val_agen2d)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN 2DD %d\n", val_agen2dd)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN 2DT %d\n", val_agen2dt)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN COLOK BEBAS %d\n", val_agencolokbebas)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN COLOK MACAU %d\n", val_agencolokmacau)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN COLOK NAGA %d\n", val_agencoloknaga)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN COLOK JITU %d\n", val_agencolokjitu)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN 5050UMUM %d\n", val_agen5050umum)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN 5050SPECIAL %d\n", val_agen5050special)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN 5050KOMBINASI %d\n", val_agen5050kombinasi)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN MACAU KOMBINASI %d\n", val_agenmacaukombinasi)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN DASAR %d\n", val_agendasar)
	log.Printf("DELETE REDIS AGEN PERIODE PERMAINAN SHIO %d\n", val_agenshio)
	log.Printf("DELETE REDIS AGEN PERIODE LIST BET ALL %d\n", val_agenall)
	log.Printf("DELETE REDIS AGEN PERIODE LIST BET WINNER %d\n", val_agenwinner)
	log.Printf("DELETE REDIS AGEN PERIODE LIST BET CANCEL %d\n", val_agencancel)
}
