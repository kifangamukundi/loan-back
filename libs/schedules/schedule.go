package schedules

// Schedules defines common cron schedule expressions
var Schedules = map[string]string{
	"EVERY_MINUTE":     "* * * * *",
	"EVERY_5_MINUTES":  "*/5 * * * *",
	"EVERY_10_MINUTES": "*/10 * * * *",
	"EVERY_15_MINUTES": "*/15 * * * *",
	"EVERY_30_MINUTES": "*/30 * * * *",
	"HOURLY":           "0 * * * *",
	"HOURLY_AT_15":     "15 * * * *",
	"HOURLY_AT_30":     "30 * * * *",
	"HOURLY_AT_45":     "45 * * * *",
	"DAILY":            "0 0 * * *",
	"DAILY_AT_1":       "0 1 * * *",
	"DAILY_AT_2":       "0 2 * * *",
	"WEEKLY":           "0 3 * * 0",
	"WEEKLY_AT_5":      "0 5 * * 1",
	"MONTHLY":          "0 6 1 * *",
	"MONTHLY_LAST_DAY": "0 23 L * *",
	"EVERY_2_HOURS":    "0 */2 * * *",
	"EVERY_4_HOURS":    "0 */4 * * *",
	"DAILY_NOON":       "0 12 * * *",
	"DAILY_AT_5":       "0 17 * * *",
	"EVERY_6_HOURS":    "0 */6 * * *",
}
