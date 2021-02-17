package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/kozlov-d/mx-api-trainee/common"
	"github.com/kozlov-d/mx-api-trainee/data"
	"github.com/tealeg/xlsx/v3"
)

type Server struct {
	Router *mux.Router
	Data   *data.Data
	Config common.Config
}

//NewServer returns server initialized with given config and other internals
func NewServer(c common.Config) Server {
	s := Server{Config: c, Router: mux.NewRouter().StrictSlash(true)}
	s.initStorages()
	s.setupRouter()
	return s
}

func (s *Server) download(link url.URL, taskID, merchantID int) {
	start := time.Now()
	//maybe worth to separate for testing
	resp, err := http.Get(link.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	wb, err := xlsx.OpenBinary([]byte(string(body)))
	if err != nil {
		log.Fatal(err)
		return
	}
	sh := wb.Sheets[0]
	if sh == nil {
		log.Fatal("Workbook does not contain any sheet")
		return
	}
	//aka create if not exists
	m, err := s.Data.UpsertMerchant(merchantID)
	if err != nil {
		log.Fatal(err)
	}
	err = sh.ForEachRow(func(r *xlsx.Row) error {
		oR := &common.OfferRow{}
		if err := r.ReadStruct(oR); err != nil {
			//can't read row: e.g. invalid type conversion
			//also skips if any column value is empty
			//might ran into bigger amount of @missed than expected
			//won't miss in case of conversion `empty -> bool` and empty string
			//returning @error would cause end of ForEachRow
			//invalid type conversion won't lead to panic, just skips the row
			s.Data.Cache.UpdateTask(taskID, common.Task{Missed: 1})
			return nil
		}
		oC := common.Convert(*oR)
		//add @missed if fields not valid
		if ok := oC.Validate(); ok != true {
			s.Data.Cache.UpdateTask(taskID, common.Task{Missed: 1})
			return nil
		}
		//check if offer already exist
		o, err := s.Data.SelectOffer(oC.OfferID, m.MerchantID)
		if err != nil {
			log.Fatal(err)
		}
		if o != nil {
			if oC.Available {
				//add @updated if offer does exist and @available is 1
				s.Data.UpdateOffer(oC)
				s.Data.Cache.UpdateTask(taskID, common.Task{Updated: 1})
				return nil
			}
			{ //add @deleted if merchant and offer does exists and @available is 0
				s.Data.DeleteOffer(o.OfferID, m.MerchantID)
				s.Data.Cache.UpdateTask(taskID, common.Task{Deleted: 1})
				return nil
			}
		}
		//add @created if offer doesn't exist and @available is 1
		if oC.Available {
			if err := s.Data.InsertOffer(oC, m.MerchantID); err != nil {
				log.Fatal(err)
			}
			s.Data.Cache.UpdateTask(taskID, common.Task{Created: 1})
			return nil
		}
		return nil
	}, xlsx.SkipEmptyRows)

	if err != nil {
		log.Fatal(err)
	}
	s.Data.Cache.CompleteTask(taskID, start)
}

//Start the server
func (s Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%s", s.Config.AppConfig.Port), s.Router)
}

func (s *Server) initStorages() {
	s.Data = &data.Data{Cache: data.CreateCache()}
	s.Data.InitializeDB(s.Config.DBConfig)
}
