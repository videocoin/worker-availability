package stats

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func reportHandler(appctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(rw, "Use GET.")
			return
		}
		appctx.Log.Debugf("received request %v", req.URL)
		values := req.URL.Query()

		var (
			startT, endT time.Time
			start        = values.Get("start")
			end          = values.Get("end")
			duration     = values.Get("duration")
			err          error
		)
		if len(duration) > 0 {
			d, err := time.ParseDuration(duration)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "Failed to parse duration %s: %v", duration, err)
				return
			}
			endT = time.Now()
			startT = endT.Add(-d)
		} else {
			startT, err = time.Parse(time.RFC3339, start)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "Failed to parse start %s: %v", start, err)
				return
			}
			endT, err = time.Parse(time.RFC3339, end)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "Failed to parse end %s: %v", end, err)
				return
			}
		}
		report, err := CreateReport(appctx, startT, endT)
		if err != nil {
			appctx.Log.Errorf("failed to create a report: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(rw, "Temporary internal error.")
			return
		}
		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=report-%d.csv", time.Now().Unix()))
		rw.WriteHeader(http.StatusOK)
		_, err = report.WriteTo(rw)
		if err != nil {
			appctx.Log.Debugf("can't send report to the user %v", err)
		}
	}
}

func Serve(appctx Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/report", reportHandler(appctx))
	server := &http.Server{
		Addr:           appctx.C.Address,
		Handler:        mux,
		ReadTimeout:    appctx.C.ReadHTTPTimeout,
		WriteTimeout:   appctx.C.WriteHTTPTimeout,
		MaxHeaderBytes: appctx.C.MaxHeaderBytes,
	}
	appctx.Log.Infof("http server is listening on %v", appctx.C.Address)
	go func() {
		<-appctx.Done()
		appctx.Log.Info("application received interrupt")
		err := server.Close()
		if err != nil {
			appctx.Log.Debugf("server closed with error %v", err)
		}
	}()
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			appctx.Log.Errorf("http server crashed with %v", err)
			return err
		}
	}
	appctx.Log.Info("server exited")
	return nil
}
