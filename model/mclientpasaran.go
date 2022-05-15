package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"sync"
	"time"

	"bitbucket.org/isbtotogroup/apisavetransaksi_go/config"
	"bitbucket.org/isbtotogroup/apisavetransaksi_go/db"
	"bitbucket.org/isbtotogroup/apisavetransaksi_go/entities"
	"bitbucket.org/isbtotogroup/apisavetransaksi_go/helpers"
	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

var mutex sync.RWMutex

func Get_Domain() (helpers.Response, error) {
	var obj entities.Model_domain
	var arraobj []entities.Model_domain
	var res helpers.Response
	msg := "Data Not Found"
	render_page := time.Now()
	ctx := context.Background()
	con := db.CreateCon()
	sql_select := `SELECT 
		nmdomain    
		FROM ` + config.DB_tbl_mst_domain + `  
		WHERE tipedomain = 'FRONTEND'
		AND statusdomain ='RUNNING'  
	`
	rowdomain, err := con.QueryContext(ctx, sql_select)
	defer rowdomain.Close()
	helpers.ErrorCheck(err)

	for rowdomain.Next() {
		var nmdomain_db string
		err = rowdomain.Scan(&nmdomain_db)
		if err != nil {
			return res, err
		}
		obj.Domain = nmdomain_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Totalrecord = len(arraobj)
	res.Record = arraobj
	res.Time = time.Since(render_page).String()

	return res, nil
}

type datajobs struct {
	Idtrxkeluarandetail      string
	Idtrxkeluaran            string
	Datetimedetail           string
	Ipaddress                string
	Idcompany                string
	Username                 string
	Typegame                 string
	Nomortogel               string
	Posisitogel              string
	Bet                      string
	Diskon                   string
	Win                      string
	Kei                      string
	Browsertogel             string
	Devicetogel              string
	Statuskeluarandetail     string
	Createkeluarandetail     string
	Createdatekeluarandetail string
}
type dataresult struct {
	Idtrxkeluarandetail string
	Status              string
}
type temp_data struct {
	Idtrxkeluaran string
	Ipaddress     string
	Idcompany     string
	Username      string
	Typegame      string
	Nomortogel    string
	Posisitogel   string
	Bet           string
	Diskon        string
	Win           string
	Kei           string
	Browsertogel  string
	Devicetogel   string
}

func Savetransaksi(client_username, client_company, idtrxkeluaran, idcomppasaran, devicemember, formipaddress, timezone string, clientdompet, totalbayarbet int, list4d string) (helpers.ResponseSaveTransaksi, error) {
	var res helpers.ResponseSaveTransaksi
	var obj temp_data
	var arraobj []temp_data
	con := db.CreateCon()
	ctx := context.Background()
	render_page := time.Now()
	tglnow, _ := goment.New()
	flag_loop := false
	flag_next := false
	msg := "Failed"
	msg_nomor := ""
	totalbelanja := totalbayarbet
	dompet := clientdompet
	pasaran_code := ""
	limittotal_togel4d := 0
	limittotal_togel3d := 0
	limittotal_togel3dd := 0
	limittotal_togel2d := 0
	limittotal_togel2dd := 0
	limittotal_togel2dt := 0
	limittotal_togelcolokbebas := 0
	limittotal_togelcolokmacau := 0
	limittotal_togelcoloknaga := 0
	limittotal_togelcolokjitu := 0
	limittotal_togel5050umum := 0
	limittotal_togel5050special := 0
	limittotal_togel5050kombinasi := 0
	limittotal_togelkombinasi := 0
	limittotal_togeldasar := 0
	limittotal_togelshio := 0
	limit_togel4d := 0
	limit_togel3d := 0
	limit_togel3dd := 0
	limit_togel2d := 0
	limit_togel2dd := 0
	limit_togel2dt := 0

	limittotal_togel4d_fullbb := 0
	limittotal_togel3d_fullbb := 0
	limittotal_togel3dd_fullbb := 0
	limittotal_togel2d_fullbb := 0
	limittotal_togel2dd_fullbb := 0
	limittotal_togel2dt_fullbb := 0

	limit_togel4d_fullbb := 0
	limit_togel3d_fullbb := 0
	limit_togel3dd_fullbb := 0
	limit_togel2d_fullbb := 0
	limit_togel2dd_fullbb := 0
	limit_togel2dt_fullbb := 0

	limit_togelcolokbebas := 0
	limit_togelcolokmacau := 0
	limit_togelcoloknaga := 0
	limit_togelcolokjitu := 0
	limit_togel5050umum := 0
	limit_togel5050special := 0
	limit_togel5050kombinasi := 0
	limit_togelkombinasi := 0
	limit_togeldasar := 0
	limit_togelshio := 0
	if int(dompet) < int(totalbelanja) {
		msg = "Balance Anda Tidak Cukup"
		flag_loop = true
	}
	// _, _, view_client_invoice := Get_mappingdatabase(client_company)
	_, trx_keluarantogel_detail, view_client_invoice := Get_mappingdatabase(client_company)

	sql_select := `SELECT 
		idpasarantogel, 
		limittotal_togel_4d, limittotal_togel_3d, limittotal_togel_3dd, limittotal_togel_2d, limittotal_togel_2dd, limittotal_togel_2dt, 
		limit_togel_4d, limit_togel_3d, limit_togel_3dd, limit_togel_2d, limit_togel_2dd, limit_togel_2dt, 
		limit_togel_4d_fullbb, limit_togel_3d_fullbb, limit_togel_3dd_fullbb, limit_togel_2d_fullbb, limit_togel_2dd_fullbb, limit_togel_2dt_fullbb, 
		limittotal_togel_4d_fullbb, limittotal_togel_3d_fullbb, limittotal_togel_3dd_fullbb, limittotal_togel_2d_fullbb, limittotal_togel_2dd_fullbb, limittotal_togel_2dt_fullbb, 
		limittotal_togel_colokbebas, limit_togel_colokbebas, 
		limittotal_togel_colokmacau, limit_togel_colokmacau, 
		limittotal_togel_coloknaga, limit_togel_coloknaga, 
		limittotal_togel_colokjitu, limit_togel_colokjitu, 
		limittotal_togel_5050umum, limit_togel_5050umum, 
		limittotal_togel_5050special, limit_togel_5050special, 
		limittotal_togel_5050kombinasi, limit_togel_5050kombinasi, 
		limittotal_togel_kombinasi, limit_togel_kombinasi, 
		limittotal_togel_dasar, limit_togel_dasar, 
		limittotal_togel_shio, limit_togel_shio
		FROM ` + config.DB_VIEW_CLIENT_VIEW_PASARAN + `  
		WHERE idcompany = ? 
		AND idcomppasaran = ? 
	`
	row, err := con.QueryContext(ctx, sql_select, client_company, idcomppasaran)

	helpers.ErrorCheck(err)
	nolimit := 0
	for row.Next() {
		nolimit = nolimit + 1
		var (
			idpasarantogel_db                                                                                                                                                         string
			limittotal_togel_4d_db, limittotal_togel_3d_db, limittotal_togel_3dd_db, limittotal_togel_2d_db, limit_togeltotal_2dd_db, limittotal_togel_2dt_db                         float32
			limit_togel_4d_db, limit_togel_3d_db, limit_togel_3dd_db, limit_togel_2d_db, limit_togel_2dd_db, limit_togel_2dt_db                                                       float32
			limit_togel_4d_fullbb, limit_togel_3d_fullbb, limit_togel_3dd_fullbb, limit_togel_2d_fullbb, limit_togel_2dd_fullbb, limit_togel_2dt_fullbb                               float32
			limittotal_togel_4d_fullbb, limittotal_togel_3d_fullbb, limittotal_togel_3dd_fullbb, limittotal_togel_2d_fullbb, limittotal_togel_2dd_fullbb, limittotal_togel_2dt_fullbb float32
			limittotal_togel_colokbebas_db, limittotal_togel_colokmacau_db, limittotal_togel_coloknaga_db, limittotal_togel_colokjitu_db                                              float32
			limit_togel_colokbebas_db, limit_togel_colokmacau_db, limit_togel_coloknaga_db, limit_togel_colokjitu_db                                                                  float32
			limittotal_togel_5050umum_db, limittotal_togel_5050special_db, limittotal_togel_5050kombinasi_db, limittotal_togel_kombinasi_db                                           float32
			limit_togel_5050umum_db, limit_togel_5050special_db, limit_togel_5050kombinasi_db, limit_togel_kombinasi_db                                                               float32
			limittotal_togel_dasar_db, limittotal_togel_shio_db                                                                                                                       float32
			limit_togel_dasar_db, limit_togel_shio_db                                                                                                                                 float32
		)
		err = row.Scan(
			&idpasarantogel_db,
			&limittotal_togel_4d_db, &limittotal_togel_3d_db, &limittotal_togel_3dd_db, &limittotal_togel_2d_db, &limit_togeltotal_2dd_db, &limittotal_togel_2dt_db,
			&limit_togel_4d_db, &limit_togel_3d_db, &limit_togel_3dd_db, &limit_togel_2d_db, &limit_togel_2dd_db, &limit_togel_2dt_db,
			&limit_togel_4d_fullbb, &limit_togel_3d_fullbb, &limit_togel_3dd_fullbb, &limit_togel_2d_fullbb, &limit_togel_2dd_fullbb, &limit_togel_2dt_fullbb,
			&limittotal_togel_4d_fullbb, &limittotal_togel_3d_fullbb, &limittotal_togel_3dd_fullbb, &limittotal_togel_2d_fullbb, &limittotal_togel_2dd_fullbb, &limittotal_togel_2dt_fullbb,
			&limittotal_togel_colokbebas_db, &limit_togel_colokbebas_db,
			&limittotal_togel_colokmacau_db, &limit_togel_colokmacau_db,
			&limittotal_togel_coloknaga_db, &limit_togel_coloknaga_db,
			&limittotal_togel_colokjitu_db, &limit_togel_colokjitu_db,
			&limittotal_togel_5050umum_db, &limit_togel_5050umum_db,
			&limittotal_togel_5050special_db, &limit_togel_5050special_db,
			&limittotal_togel_5050kombinasi_db, &limit_togel_5050kombinasi_db,
			&limittotal_togel_kombinasi_db, &limit_togel_kombinasi_db,
			&limittotal_togel_dasar_db, &limit_togel_dasar_db, &limittotal_togel_shio_db, &limit_togel_shio_db)

		helpers.ErrorCheck(err)
		pasaran_code = idpasarantogel_db
		limittotal_togel4d = int(limittotal_togel_4d_db)
		limittotal_togel3d = int(limittotal_togel_3d_db)
		limittotal_togel3dd = int(limittotal_togel_3dd_db)
		limittotal_togel2d = int(limittotal_togel_2d_db)
		limittotal_togel2dd = int(limit_togeltotal_2dd_db)
		limittotal_togel2dt = int(limittotal_togel_2dt_db)
		limit_togel4d = int(limit_togel_4d_db)
		limit_togel3d = int(limit_togel_3d_db)
		limit_togel3dd = int(limit_togel_3dd_db)
		limit_togel2d = int(limit_togel_2d_db)
		limit_togel2dd = int(limit_togel_2dd_db)
		limit_togel2dt = int(limit_togel_2dt_db)

		limittotal_togel4d_fullbb = int(limittotal_togel_4d_fullbb)
		limittotal_togel3d_fullbb = int(limittotal_togel_3d_fullbb)
		limittotal_togel3dd_fullbb = int(limittotal_togel_3dd_fullbb)
		limittotal_togel2d_fullbb = int(limittotal_togel_2d_fullbb)
		limittotal_togel2dd_fullbb = int(limittotal_togel_2dd_fullbb)
		limittotal_togel2dt_fullbb = int(limittotal_togel_2dt_fullbb)

		limit_togel4d_fullbb = int(limit_togel_4d_fullbb)
		limit_togel3d_fullbb = int(limit_togel_3d_fullbb)
		limit_togel3dd_fullbb = int(limit_togel_3dd_fullbb)
		limit_togel2d_fullbb = int(limit_togel_2d_fullbb)
		limit_togel2dd_fullbb = int(limit_togel_2dd_fullbb)
		limit_togel2dt_fullbb = int(limit_togel_2dt_fullbb)

		limittotal_togelcolokbebas = int(limittotal_togel_colokbebas_db)
		limittotal_togelcolokmacau = int(limittotal_togel_colokmacau_db)
		limittotal_togelcoloknaga = int(limittotal_togel_coloknaga_db)
		limittotal_togelcolokjitu = int(limittotal_togel_colokjitu_db)
		limittotal_togel5050umum = int(limittotal_togel_5050umum_db)
		limittotal_togel5050special = int(limittotal_togel_5050special_db)
		limittotal_togel5050kombinasi = int(limittotal_togel_5050kombinasi_db)
		limittotal_togelkombinasi = int(limittotal_togel_kombinasi_db)
		limittotal_togeldasar = int(limittotal_togel_dasar_db)
		limittotal_togelshio = int(limittotal_togel_shio_db)

		limit_togelcolokbebas = int(limit_togel_colokbebas_db)
		limit_togelcolokmacau = int(limit_togel_colokmacau_db)
		limit_togelcoloknaga = int(limit_togel_coloknaga_db)
		limit_togelcolokjitu = int(limit_togel_colokjitu_db)
		limit_togel5050umum = int(limit_togel_5050umum_db)
		limit_togel5050special = int(limit_togel_5050special_db)
		limit_togel5050kombinasi = int(limit_togel_5050kombinasi_db)
		limit_togelkombinasi = int(limit_togel_kombinasi_db)
		limit_togeldasar = int(limit_togel_dasar_db)
		limit_togelshio = int(limit_togel_shio_db)
	}
	defer row.Close()
	if nolimit > 0 {
		statuspasaran := _checkpasaran(client_company, pasaran_code)

		if statuspasaran == "OFFLINE" {
			msg = "Maaf, Pasaran Sudah Tutup"
			flag_loop = true
		}
	}
	if !flag_loop {
		permainan := ""
		var limit_total_togel int = 0
		var limit_global_togel int = 0
		var limittotal_sum int = 0
		var limit_sum int = 0
		var totalbet_all int = 0
		var totalbet_all_limit int = 0
		var totalbayar int = 0
		flag_save := false
		json := []byte(list4d)

		jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			nomor_DD, _, _, _ := jsonparser.Get(value, "nomor")
			permainan_DD, _, _, _ := jsonparser.Get(value, "permainan")
			tipetoto_DD, _, _, _ := jsonparser.Get(value, "tipetoto")
			bet_DD, _, _, _ := jsonparser.Get(value, "bet")
			diskonpercen_DD, _, _, _ := jsonparser.Get(value, "diskonpercen")
			kei_percen_DD, _, _, _ := jsonparser.Get(value, "kei_percen")
			win_DD, _, _, _ := jsonparser.Get(value, "win")
			switch string(permainan_DD) {
			case "4D":
				permainan = "4D/3D/2D"
				if string(tipetoto_DD) == "DISC" {
					limit_total_togel = limittotal_togel4d
					limit_global_togel = limit_togel4d
				} else {
					limit_total_togel = limittotal_togel4d_fullbb
					limit_global_togel = limit_togel4d_fullbb
				}
			case "3D":
				permainan = "4D/3D/2D"
				if string(tipetoto_DD) == "DISC" {
					limit_total_togel = limittotal_togel3d
					limit_global_togel = limit_togel3d
				} else {
					limit_total_togel = limittotal_togel3d_fullbb
					limit_global_togel = limit_togel3d_fullbb
				}
			case "3DD":
				permainan = "4D/3D/2D"
				if string(tipetoto_DD) == "DISC" {
					limit_total_togel = limittotal_togel3dd
					limit_global_togel = limit_togel3dd
				} else {
					limit_total_togel = limittotal_togel3dd_fullbb
					limit_global_togel = limit_togel3dd_fullbb
				}
			case "2D":
				permainan = "4D/3D/2D"
				if string(tipetoto_DD) == "DISC" {
					limit_total_togel = limittotal_togel2d
					limit_global_togel = limit_togel2d
				} else {
					limit_total_togel = limittotal_togel2d_fullbb
					limit_global_togel = limit_togel2d_fullbb
				}
			case "2DD":
				permainan = "4D/3D/2D"
				if string(tipetoto_DD) == "DISC" {
					limit_total_togel = limittotal_togel2dd
					limit_global_togel = limit_togel2dd
				} else {
					limit_total_togel = limittotal_togel2dd_fullbb
					limit_global_togel = limit_togel2dd_fullbb
				}
			case "2DT":
				permainan = "4D/3D/2D"
				if string(tipetoto_DD) == "DISC" {
					limit_total_togel = limittotal_togel2dt
					limit_global_togel = limit_togel2dt
				} else {
					limit_total_togel = limittotal_togel2dt_fullbb
					limit_global_togel = limit_togel2dt_fullbb
				}
			case "COLOK_BEBAS":
				permainan = "COLOK BEBAS"
				limit_total_togel = limittotal_togelcolokbebas
				limit_global_togel = limit_togelcolokbebas
			case "COLOK_MACAU":
				permainan = "COLOK MACAU"
				limit_total_togel = limittotal_togelcolokmacau
				limit_global_togel = limit_togelcolokmacau
			case "COLOK_NAGA":
				permainan = "COLOK NAGA"
				limit_total_togel = limittotal_togelcoloknaga
				limit_global_togel = limit_togelcoloknaga
			case "COLOK_JITU":
				permainan = "COLOK JITU"
				limit_total_togel = limittotal_togelcolokjitu
				limit_global_togel = limit_togelcolokjitu
			case "50_50_UMUM":
				permainan = "50 - 50 UMUM"
				limit_total_togel = limittotal_togel5050umum
				limit_global_togel = limit_togel5050umum
			case "50_50_SPECIAL":
				permainan = "50 - 50 SPECIAL"
				limit_total_togel = limittotal_togel5050special
				limit_global_togel = limit_togel5050special
			case "50_50_KOMBINASI":
				permainan = "50 - 50 KOMBINASI"
				limit_total_togel = limittotal_togel5050kombinasi
				limit_global_togel = limit_togel5050kombinasi
			case "MACAU_KOMBINASI":
				permainan = "MACAU / KOMBINASI"
				limit_total_togel = limittotal_togelkombinasi
				limit_global_togel = limit_togelkombinasi
			case "DASAR":
				permainan = "DASAR"
				limit_total_togel = limittotal_togeldasar
				limit_global_togel = limit_togeldasar
			case "SHIO":
				permainan = "SHIO"
				limit_total_togel = limittotal_togelshio
				limit_global_togel = limit_togelshio
			}
			bet := string(bet_DD)
			diskon := string(diskonpercen_DD)
			kei := string(kei_percen_DD)
			bet2, _ := strconv.Atoi(bet)
			diskon2, _ := strconv.ParseFloat(diskon, 64)
			kei2, _ := strconv.ParseFloat(kei, 64)
			diskonvalue := math.Ceil(float64(bet2) * diskon2)
			keivalue := math.Ceil(float64(bet2) * kei2)
			sqllimitsum := `SELECT
					CAST(COALESCE(SUM(bet), 0) as UNSIGNED) AS total 
					FROM ` + view_client_invoice + ` 
					WHERE idtrxkeluaran = ?
					AND typegame = ?
					AND nomortogel = ? 
					AND posisitogel = ? 
				`

			row := con.QueryRowContext(ctx, sqllimitsum, idtrxkeluaran, string(permainan_DD), string(nomor_DD), string(tipetoto_DD))
			switch e := row.Scan(&limit_sum); e {
			case sql.ErrNoRows:
				log.Println("No rows were returned!")
			case nil:
				// log.Println(iddoc)
			default:
				log.Println("ERROR LIMIT GLOBAL :", e)
			}

			totalbet_all = limit_sum + bet2
			if limit_global_togel < totalbet_all {
				flag_save = true
				msg_nomor += "Nomor ini : " + string(nomor_DD) + ", sudah melebihi LIMIT GLOBAL<br>"
			}

			sqllimittotalsum := `SELECT
					CAST(COALESCE(SUM(bet), 0) as UNSIGNED) AS total 
					FROM ` + view_client_invoice + ` 
					WHERE idtrxkeluaran = ?
					AND typegame = ?
					AND nomortogel = ? 
					AND username = ?  
					AND posisitogel = ? 
				`

			row_limit := con.QueryRowContext(ctx, sqllimittotalsum, idtrxkeluaran, string(permainan_DD), string(nomor_DD), client_username, string(tipetoto_DD))
			switch e := row_limit.Scan(&limittotal_sum); e {
			case sql.ErrNoRows:
				log.Println("No rows were returned!")
			case nil:
				// log.Println(iddoc)
			default:
				log.Println("ERROR LIMIT TOTAL :", e)
			}

			totalbet_all_limit = int(limittotal_sum) + bet2
			if limit_total_togel < totalbet_all_limit {
				flag_save = true
				msg_nomor += "Nomor ini : " + string(nomor_DD) + ", sudah melebihi LIMIT TOTAL<br>"
			}

			log.Printf("Limit GLOBAL : %d", limit_sum)
			log.Printf("Limit Total : %d", limit_total_togel)
			log.Printf("Limit Total  check number by username : %d", limittotal_sum)
			log.Printf("Limit Total  sum : %d", totalbet_all_limit)
			log.Printf("bet : %d", bet2)
			log.Printf("FLAG SAVE : %t", flag_save)
			log.Printf("PERMAINAN : %s", permainan)
			log.Printf("TIPE : %s", tipetoto_DD)

			if !flag_save { //VALID

				bayar := bet2 - int(diskonvalue) - int(keivalue)
				totalbayar = totalbayar + int(bayar)

				log.Printf("NOMOR : %s", nomor_DD)
				log.Printf("BET : %d", bet2)
				log.Printf("DISKON PERSEN : %.2f", diskon2)
				log.Printf("DISKON AMOUNT : %d", int(diskonvalue))
				log.Printf("KEI PERSEN : %.2f", kei2)
				log.Printf("KEI AMOUNT : %d", int(keivalue))
				log.Printf("BAYAR : %d", bayar)
				log.Printf("TOTAL BAYAR : %d", totalbayar)

				obj.Idtrxkeluaran = idtrxkeluaran
				obj.Ipaddress = formipaddress
				obj.Idcompany = client_company
				obj.Username = client_username
				obj.Typegame = string(permainan_DD)
				obj.Nomortogel = string(nomor_DD)
				obj.Posisitogel = string(tipetoto_DD)
				obj.Bet = string(bet_DD)
				obj.Diskon = string(diskonpercen_DD)
				obj.Win = string(win_DD)
				obj.Kei = string(kei_percen_DD)
				obj.Browsertogel = timezone
				obj.Devicetogel = devicemember
				arraobj = append(arraobj, obj)
			}
			flag_save = false
		})
		log.Println("Total keranjang valid :", len(arraobj))
		totals_bet := len(arraobj)
		log.Println()

		if totals_bet > 0 {
			runtime.GOMAXPROCS(8)
			totalWorker := 100
			buffer_bet := totals_bet + 1
			jobs_bet := make(chan datajobs, buffer_bet)
			results_bet := make(chan dataresult, buffer_bet)

			wg := &sync.WaitGroup{}
			for w := 0; w < totalWorker; w++ {
				wg.Add(1)
				mutex.Lock()
				go _doJobInsert(trx_keluarantogel_detail, jobs_bet, results_bet, con, wg)
				mutex.Unlock()
			}
			for i := 0; i < len(arraobj); i++ {
				mutex.Lock()
				year := tglnow.Format("YY")
				month := tglnow.Format("MM")
				field_column_counter := trx_keluarantogel_detail + tglnow.Format("YYYY") + month
				idrecord_counter := Get_counter(field_column_counter)

				idrecord_counter2 := strconv.Itoa(idrecord_counter)
				idrecord := string(year) + string(month) + idrecord_counter2

				jobs_bet <- datajobs{
					Idtrxkeluarandetail:      idrecord,
					Idtrxkeluaran:            arraobj[i].Idtrxkeluaran,
					Datetimedetail:           tglnow.Format("YYYY-MM-DD HH:mm:ss"),
					Ipaddress:                arraobj[i].Ipaddress,
					Idcompany:                arraobj[i].Idcompany,
					Username:                 arraobj[i].Username,
					Typegame:                 arraobj[i].Typegame,
					Nomortogel:               arraobj[i].Nomortogel,
					Posisitogel:              arraobj[i].Posisitogel,
					Bet:                      arraobj[i].Bet,
					Diskon:                   arraobj[i].Diskon,
					Win:                      arraobj[i].Win,
					Kei:                      arraobj[i].Kei,
					Browsertogel:             arraobj[i].Browsertogel,
					Devicetogel:              arraobj[i].Devicetogel,
					Statuskeluarandetail:     "RUNNING",
					Createkeluarandetail:     client_username,
					Createdatekeluarandetail: tglnow.Format("YYYY-MM-DD HH:mm:ss")}
				mutex.Unlock()
				fmt.Println("Data valid : ", arraobj[i].Nomortogel)
			}
			close(jobs_bet)
			flag_next = true
			for a := 1; a <= totals_bet; a++ { //RESULT
				flag_result := <-results_bet
				if flag_result.Status == "Failed" {
					flag_next = false
				}
			}
			close(results_bet)
			wg.Wait()
			flag_next = true
		}

		log.Println(time.Since(render_page).String())
		if flag_next {
			msg = "Success"
		} else {
			msg = "Failed"
		}
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Messageerror = msg_nomor
		res.Totalbayar = totalbayar
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Messageerror = nil
		res.Totalbayar = 0
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}

func _doJobInsert(fieldtable string, jobs <-chan datajobs, results chan<- dataresult, con *sql.DB, wg *sync.WaitGroup) {
	ctx := context.Background()
	for capture := range jobs {
		for {
			var outerError error
			func(outerError *error) {
				defer func() {
					if err := recover(); err != nil {
						*outerError = fmt.Errorf("%v", err)
					}
				}()
				sql_insert := `
					INSERT INTO ` + fieldtable + ` 
					(
						idtrxkeluarandetail, idtrxkeluaran, datetimedetail,
						ipaddress, idcompany, username, typegame, nomortogel,posisitogel, bet,
						diskon, win, kei, browsertogel, devicetogel, statuskeluarandetail, 
						createkeluarandetail, createdatekeluarandetail
					) values (
						?, ?, ?, 
						?, ?, ?, ?, ?, ?,?, 
						?, ?, ?, ?, ?, ?,
						?, ?
					)
				`

				stmt, err := con.PrepareContext(ctx, sql_insert)
				helpers.ErrorCheck(err)
				defer stmt.Close()
				res, err := stmt.ExecContext(ctx,
					capture.Idtrxkeluarandetail,
					capture.Idtrxkeluaran,
					capture.Datetimedetail,
					capture.Ipaddress,
					capture.Idcompany,
					capture.Username,
					capture.Typegame,
					capture.Nomortogel,
					capture.Posisitogel,
					capture.Bet,
					capture.Diskon,
					capture.Win,
					capture.Kei,
					capture.Browsertogel,
					capture.Devicetogel,
					capture.Statuskeluarandetail,
					capture.Createkeluarandetail,
					capture.Createdatekeluarandetail)
				helpers.ErrorCheck(err)
				id_insert, err_insert := res.RowsAffected()
				helpers.ErrorCheck(err_insert)
				if id_insert > 0 {
					results <- dataresult{Idtrxkeluarandetail: capture.Idtrxkeluarandetail, Status: "Success"}
				} else {
					results <- dataresult{Idtrxkeluarandetail: capture.Idtrxkeluarandetail, Status: "Failed"}
				}
			}(&outerError)
			if outerError == nil {
				break
			}
		}
	}
	wg.Done()
}
func _checkpasaran(client_company, pasaran_code string) string {
	var myDays = []string{"minggu", "senin", "selasa", "rabu", "kamis", "jumat", "sabtu"}
	statuspasaran := "ONLINE"

	con := db.CreateCon()
	ctx := context.Background()

	tglnow, _ := goment.New()
	daynow := tglnow.Format("d")
	intVar, _ := strconv.ParseInt(daynow, 0, 8)
	daynowhari := myDays[intVar]

	tbl_trx_keluaran, _, _ := Get_mappingdatabase(client_company)

	sqlpasaran := `SELECT 
		idcomppasaran, nmpasarantogel, 
		jamtutup, jamopen  
		FROM ` + config.DB_VIEW_CLIENT_VIEW_PASARAN + `  
		WHERE idcompany = ? 
		AND idpasarantogel = ? 
	`

	rowpasaran, err := con.QueryContext(ctx, sqlpasaran, client_company, pasaran_code)
	defer rowpasaran.Close()
	helpers.ErrorCheck(err)
	for rowpasaran.Next() {
		var (
			idcomppasaran, nmpasarantogel, jamtutup, jamopen string
			idtrxkeluaran, keluaranperiode, haripasaran      string
		)

		err = rowpasaran.Scan(&idcomppasaran, &nmpasarantogel, &jamtutup, &jamopen)
		helpers.ErrorCheck(err)

		sqlkeluaran := `
			SELECT 
			idtrxkeluaran, keluaranperiode
			FROM ` + tbl_trx_keluaran + `  
			WHERE idcompany = ?
			AND idcomppasaran = ?
			AND keluarantogel = ''
			LIMIT 1
		`
		err := con.QueryRowContext(ctx, sqlkeluaran, client_company, idcomppasaran).Scan(&idtrxkeluaran, &keluaranperiode)
		helpers.ErrorCheck(err)

		sqlpasaranonline := `
			SELECT
				haripasaran
			FROM ` + config.DB_tbl_mst_company_game_pasaran_offline + ` 
			WHERE idcomppasaran = ?
			AND idcompany = ? 
			AND haripasaran = ? 
		`

		errpasaranonline := con.QueryRowContext(ctx, sqlpasaranonline, idcomppasaran, client_company, daynowhari).Scan(&haripasaran)

		if errpasaranonline != sql.ErrNoRows {
			tglskrg := tglnow.Format("YYYY-MM-DD HH:mm:ss")
			jamtutup := tglnow.Format("YYYY-MM-DD") + " " + jamtutup
			jamopen := tglnow.Format("YYYY-MM-DD") + " " + jamopen

			if tglskrg >= jamtutup && tglskrg <= jamopen {
				statuspasaran = "OFFLINE"
			}
		}
	}

	return statuspasaran
}
