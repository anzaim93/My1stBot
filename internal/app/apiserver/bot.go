package apiserver

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

type spending struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Spending int    `json:"spending"`
	Category string `json:"category"`
	Date     string `json:"date"`
}

func (s *APIServer) bot() {
	conStr := "user=anzaim password=12345 dbname=tgbot sslmode=disable"
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI("6178163320:AAFiq1wXfRRwmSA0Dx-zxn_OyR8cFbF6OQY")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			command := strings.Split(update.Message.Text, " ")

			if _, err := strconv.Atoi(command[0]); err == nil {
				if len(command) > 2 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не верный запрос"))
					if err != nil {
						return
					}
				}
				if len(command) == 1 {
					db.Exec("insert into spendings (username, spending, category, DateAndTime) values ($1, $2, $3, $4)",
						update.Message.Chat.ID, command[0], "не указано", update.Message.Time())
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "покупка добавлена"))
				}

				if len(command) == 2 {
					db.Exec("insert into spendings (username, spending, category, DateAndTime) values ($1, $2, $3, $4)",
						update.Message.Chat.ID, command[0], command[1], update.Message.Time())
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "покупка добавлена"))
				}
			} else {
				switch command[0] {
				case "week":
					if len(command) != 1 {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не верный запрос 1"))
					} else {
						rows, err := db.Query("SELECT spending, category, dateandtime from spendings WHERE username = $1 AND dateandtime > current_timestamp::date - interval '7 day'",
							update.Message.Chat.ID)
						if err != nil {
							panic(err)
						}
						defer rows.Close()

						var spendings []spending

						for rows.Next() {
							p := spending{}
							err := rows.Scan(&p.Spending, &p.Category, &p.Date)
							if err != nil {
								fmt.Println(err)
								continue
							}
							spendings = append(spendings, p)
						}
						for _, p := range spendings {
							x := strconv.Itoa(p.Spending) + " " + p.Category + " " + p.Date
							bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, x))
						}
					}

				case "month":
					if len(command) != 1 {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не верный запрос 2"))
					} else {

						rows, err := db.Query("SELECT spending, category, dateandtime from spendings WHERE username = $1 AND dateandtime > current_timestamp::date - interval '30 day' ",
							update.Message.Chat.ID)
						if err != nil {
							panic(err)
						}
						defer rows.Close()

						var spendings []spending

						for rows.Next() {
							p := spending{}
							err := rows.Scan(&p.Spending, &p.Category, &p.Date)
							if err != nil {
								fmt.Println(err)
								continue
							}
							spendings = append(spendings, p)
						}
						for _, p := range spendings {
							var x string
							x = strconv.Itoa(p.Spending) + " " + p.Category + " " + p.Date
							bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, x))
						}
					}
				case "last":
					if len(command) != 2 {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не верный запрос 3"))
					} else {

						rows, err := db.Query("SELECT spending, category, dateandtime from spendings WHERE username = $1 ORDER BY dateandtime DESC LIMIT $2",
							update.Message.Chat.ID, command[1])
						if err != nil {
							panic(err)
						}
						defer rows.Close()

						var spendings []spending

						for rows.Next() {
							p := spending{}
							err := rows.Scan(&p.Spending, &p.Category, &p.Date)
							if err != nil {
								fmt.Println(err)
								continue
							}
							spendings = append(spendings, p)
						}
						for _, p := range spendings {
							var x string
							x = strconv.Itoa(p.Spending) + " " + p.Category + " " + p.Date
							bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, x))
						}
					}
				case "/week":
					if len(command) != 1 {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не верный запрос 1"))
					} else {
						rows, err := db.Query("SELECT ID, username, spending, category, dateandtime from spendings WHERE username = $1 AND dateandtime > current_timestamp::date - interval '7 day'",
							update.Message.Chat.ID)
						if err != nil {
							panic(err)
						}
						defer rows.Close()

						var spendings []spending

						for rows.Next() {
							p := spending{}
							err := rows.Scan(&p.ID, &p.Username, &p.Spending, &p.Category, &p.Date)
							if err != nil {
								fmt.Println(err)
								continue
							}
							spendings = append(spendings, p)
						}
						s.configureRouter(spendings, "/week")
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Данные доступны по ссылке http://localhost:8080/week"))
					}
				case "/month":
					if len(command) != 1 {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не верный запрос 1"))
					} else {
						rows, err := db.Query("SELECT ID, username, spending, category, dateandtime from spendings WHERE username = $1 AND dateandtime > current_timestamp::date - interval '30 day'",
							update.Message.Chat.ID)
						if err != nil {
							panic(err)
						}
						defer rows.Close()

						var spendings []spending

						for rows.Next() {
							p := spending{}
							err := rows.Scan(&p.ID, &p.Username, &p.Spending, &p.Category, &p.Date)
							if err != nil {
								fmt.Println(err)
								continue
							}
							spendings = append(spendings, p)
						}
						s.configureRouter(spendings, "/month")
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Данные доступны по ссылке http://localhost:8080/month"))
					}
				case "/last":
					if len(command) != 2 {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не верный запрос 1"))
					} else {
						rows, err := db.Query("SELECT ID, username, spending, category, dateandtime from spendings WHERE username = $1 ORDER BY dateandtime DESC LIMIT $2",
							update.Message.Chat.ID, command[1])
						if err != nil {
							panic(err)
						}
						defer rows.Close()

						var spendings []spending

						for rows.Next() {
							p := spending{}
							err := rows.Scan(&p.ID, &p.Username, &p.Spending, &p.Category, &p.Date)
							if err != nil {
								fmt.Println(err)
								continue
							}
							spendings = append(spendings, p)
						}
						s.configureRouter(spendings, "/last")
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Данные доступны по ссылке http://localhost:8080/last"))
					}

				default:
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Чтобы добавить расходы введите целое число и название через пробел (название не является обязательным), так же Вы можете использовать следующие запросы : week - покажет траты за последние 7 дней; /week - покажет траты за последние 7 дней в формате json; month - покажет траты за последние 30 дней; /month - покажет траты за последние 30 дней в формате json; last N - покажет последние N трат; /lastN - покажет последние N трат в формате json"))

				}
			}
		}
	}
}
