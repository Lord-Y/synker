// Package processing provide all requirements to process change data capture
package processing

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"testing"
	"time"

	"github.com/Lord-Y/synker/commons"
	"github.com/Lord-Y/synker/logger"
	"github.com/jackc/fake"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestValidate_fail_falseconfig(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		var c Validate
		c.Logger = logger.NewLogger()
		c.ConfigDir = "examples/falseconfig"
		assert.Nil(c.ParseAndValidateConfig())
		return
	}
	cmd := exec.Command(
		os.Args[0],
		"synker",
		"validate",
		"-test.run=TestValidate_fail_falseconfig",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}

func TestValidate_processing(t *testing.T) {
	assert := assert.New(t)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		var c Validate
		c.Logger = logger.NewLogger()
		c.ConfigDir = "examples/schemas"
		assert.Nil(c.ParseAndValidateConfig())
		go func() {
			time.Sleep(60 * time.Second)
			<-sigc
			signal.Stop(sigc)
		}()
		c.processing()
	}()

	err = proc.Signal(os.Interrupt)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
}

func (c *Validate) testProcessing(stop chan struct{}) {
	var i int
	time.AfterFunc(100*time.Second, func() {
		stop <- struct{}{}
	})

	for {
		select {
		case <-stop:
			return
		default:
			if i == 0 {
				i++
				c.processing()
			}
		}
	}
}

func TestValidate_processing_with_new_id(t *testing.T) {
	assert := assert.New(t)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	var c Validate
	c.Logger = logger.NewLogger()
	c.ConfigDir = "examples/schemas"
	assert.Nil(c.ParseAndValidateConfig())
	var (
		userPromotionCodeValue string
		fakeCharacters         []string
	)
	for i := range 10 {
		fakeCharacters = append(fakeCharacters, fake.CharactersN(5))

		_, err = manageUserPromoCodesForUnitTesting("add", "", fakeCharacters[i])
		assert.Nil(err)
		c.Logger.Info().Msgf("manageUserPromoCodesForUnitTesting add")
	}

	userPromotionCode := "xxxxx"
	userPromotionCodeValue, err = manageUserPromoCodesForUnitTesting("update_promotion_codes", "", userPromotionCode)
	assert.Nil(err)
	c.Logger.Info().Msgf("manageUserPromoCodesForUnitTesting update")

	ridesRevenue := 60.00
	ridesRevenueOld, err := manageRidesForUnitTesting("update_rides", 0, ridesRevenue)
	assert.Nil(err)

	time.Sleep(time.Second)
	_, err = manageRidesForUnitTesting("revert_rides", ridesRevenueOld, ridesRevenue)
	assert.Nil(err)

	stop := make(chan struct{})
	go c.testProcessing(stop)

	time.Sleep(20 * time.Second)
	_, err = manageUserPromoCodesForUnitTesting("revert_update_promotion_codes", userPromotionCode, userPromotionCodeValue)
	assert.Nil(err)
	c.Logger.Info().Msgf("manageUserPromoCodesForUnitTesting revert")

	for _, fakeCharacter := range fakeCharacters {
		_, err = manageUserPromoCodesForUnitTesting("delete", "", fakeCharacter)
		assert.Nil(err)
		c.Logger.Info().Msgf("manageUserPromoCodesForUnitTesting delete")
	}

	_, err = manageRidesForUnitTesting("delete", ridesRevenueOld, ridesRevenue)
	if err != nil {
		c.Logger.Info().Msgf("err %s", err.Error())
	}
	assert.Nil(err)

	time.Sleep(40 * time.Second)
	<-stop

	if err := proc.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
}

func manageUserPromoCodesForUnitTesting(action, oldValue, fakeCharacter string) (value string, err error) {
	cfg, err := pgxpool.ParseConfig(commons.GetPGURI())
	if err != nil {
		return
	}

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	//golangci-lint fail on this check while the transaction error is checked
	defer tx.Rollback(ctx) //nolint

	switch action {
	case "delete":
		_, err = tx.Exec(
			ctx,
			"DELETE FROM user_promo_codes WHERE city = $1 AND usage_count >= $2 LIMIT 1",
			"amsterdam",
			10,
		)
	case "update_promotion_codes":
		err = tx.QueryRow(
			ctx,
			"SELECT code FROM user_promo_codes LIMIT 1",
		).Scan(&value)
		if err != nil && err.Error() != pgx.ErrNoRows.Error() {
			return
		}

		_, err = tx.Exec(
			ctx,
			"UPDATE user_promo_codes SET code = $1 WHERE code = $2 LIMIT 1",
			value,
			fakeCharacter,
		)
		if err != nil {
			return
		}
	case "revert_update_promotion_codes":
		_, err = tx.Exec(
			ctx,
			"UPDATE user_promo_codes SET code = $1 WHERE code = $2 LIMIT 1",
			fakeCharacter,
			oldValue,
		)
		if err != nil {
			return
		}
	default:
		_, err = tx.Exec(
			ctx,
			"INSERT INTO user_promo_codes VALUES($1,$2,$3,NOW(),$4)",
			"amsterdam",
			"ae147ae1-47ae-4800-8000-000000000022",
			fakeCharacter,
			10,
		)
	}
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		return
	}
	return
}

func manageRidesForUnitTesting(action string, oldValue, fakeCharacter float64) (value float64, err error) {
	cfg, err := pgxpool.ParseConfig(commons.GetPGURI())
	if err != nil {
		return
	}

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	//golangci-lint fail on this check while the transaction error is checked
	defer tx.Rollback(ctx) //nolint

	switch action {
	case "delete":
		_, err = tx.Exec(
			ctx,
			"DELETE FROM vehicle_location_histories WHERE city = $1 LIMIT 1",
			"amsterdam",
		)
	case "update_rides":
		err = tx.QueryRow(
			ctx,
			"SELECT revenue FROM rides LIMIT 1",
		).Scan(&value)
		if err != nil && err.Error() != pgx.ErrNoRows.Error() {
			return
		}

		_, err = tx.Exec(
			ctx,
			"UPDATE rides SET revenue = $1 WHERE revenue = $2 LIMIT 1",
			value,
			fakeCharacter,
		)
		if err != nil {
			return
		}
	case "revert_rides":
		_, err = tx.Exec(
			ctx,
			"UPDATE rides SET revenue = $1 WHERE revenue = $2 LIMIT 1",
			fakeCharacter,
			oldValue,
		)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		return
	}
	return
}

func TestIsMap(t *testing.T) {
	assert := assert.New(t)

	var c Validate
	m := make(map[string]interface{})

	m["a"] = "a"
	ok, multiKey := c.isMap(m)
	assert.Equal(true, ok)
	assert.Equal(false, multiKey)

	m["b"] = "b"
	ok, multiKey = c.isMap(m)
	assert.Equal(true, ok)
	assert.Equal(true, multiKey)
}

func TestLoadFiles(t *testing.T) {
	assert := assert.New(t)

	_, err := loadFiles("examples/emptydir/xx.yaml")
	assert.Error(err)
}
