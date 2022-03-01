package main

import (
	"database/sql"
	"net/http"

	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
)

type CampaignSmsWrap struct {
	Results         []models.CampaignSms `json:"results"`
	Sent            int                  `db:"sent" json:"sent"`
	Delivered       int                  `db:"delivered" json:"delivered"`
	Failed          int                  `db:"failed" json:"failed"`
	CampaignId      int                  `db:"campaign_id" json:"campaignId"`
	CampaignName    string               `db:"name" json:"campaignName"`
	CampaignAltBody string               `db:"altbody" json:"campaignAltBody"`
}

type CampaignSmsWrapByUserId struct {
	Results   []models.CampaignSms `json:"results"`
	Sent      int                  `db:"sent" json:"sent"`
	Delivered int                  `db:"delivered" json:"delivered"`
	Failed    int                  `db:"failed" json:"failed"`
}

// handleGetSmsLogsByCampaignId retrieves lists of campaign sms
func handleGetSmsLogsByCampaignId(c echo.Context) error {
	var (
		app        = c.Get("app").(*App)
		campaignId = c.Param("campaignId")
		out        []CampaignSmsWrap
		outResults []models.CampaignSms
	)

	if err := app.queries.GetCampaignSmsCounts.Select(&out, campaignId); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, okResp{[]struct{}{}})
		}

		app.log.Printf("error fetching campaign sms counts: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			app.i18n.Ts("globals.messages.errorFetching",
				"name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}

	if err := app.queries.GetCampaignSmsLogs.Select(&outResults, campaignId); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, okResp{[]struct{}{}})
		}

		app.log.Printf("error fetching campaign sms logs: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			app.i18n.Ts("globals.messages.errorFetching",
				"name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}

	if len(out) != 1 {
		app.log.Printf("error fetching campaign sms with id: %v", campaignId)
		return c.JSON(http.StatusOK, okResp{[]struct{}{}})
	}

	out[0].Results = outResults

	return c.JSON(http.StatusOK, okResp{out[0]})
}

// handleGetSmsLogsByUserId retrieves lists of campaign sms
func handleGetSmsLogsByUserId(c echo.Context) error {
	var (
		app        = c.Get("app").(*App)
		userid     = c.Param("userid")
		out        []CampaignSmsWrapByUserId
		outResults []models.CampaignSms
	)

	if err := app.queries.GetCampaignSmsCountsByUserId.Select(&out, userid); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, okResp{[]struct{}{}})
		}

		app.log.Printf("error fetching campaign sms counts: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			app.i18n.Ts("globals.messages.errorFetching",
				"name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}

	if err := app.queries.GetCampaignSmsLogsByUserId.Select(&outResults, userid); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, okResp{[]struct{}{}})
		}

		app.log.Printf("error fetching campaign sms logs: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			app.i18n.Ts("globals.messages.errorFetching",
				"name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}

	if len(out) != 1 {
		app.log.Printf("error fetching campaign sms with id: %v", userid)
		return c.JSON(http.StatusOK, okResp{[]struct{}{}})
	}

	out[0].Results = outResults

	return c.JSON(http.StatusOK, okResp{out[0]})
}
