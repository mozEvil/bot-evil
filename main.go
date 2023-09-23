package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var gBot *tgbotapi.BotAPI
var gToken string
var gChatId int64

var gUsersInChat Users

type User struct {
	id   int64
	name string
	gems uint16
}
type Users []*User

type Activity struct {
	code, name string
	gems       uint16
}
type Activities []*Activity

var gUsefulActivities = Activities{
	// Self-Development
	{"yoga", "Yoga (15 minutes)", 1},
	{"meditation", "Meditation (15 minutes)", 1},
	{"language", "Learning a foreign language (15 minutes)", 1},
	{"swimming", "Swimming (15 minutes)", 1},
	{"walk", "Walk (15 minutes)", 1},
	{"chores", "Chores", 1},

	// Work
	{"work_learning", "Studying work materials (15 minutes)", 1},
	{"portfolio_work", "Working on a portfolio project (15 minutes)", 1},
	{"resume_edit", "Resume editing (15 minutes)", 1},

	// Creativity
	{"creative", "Creative creation (15 minutes)", 1},
	{"reading", "Reading fiction literature (15 minutes)", 1},
}

var gRewards = Activities{
	// Entertainment
	{"watch_series", "Watching a series (1 episode)", 10},
	{"watch_movie", "Watching a movie (1 item)", 30},
	{"social_nets", "Browsing social networks (30 minutes)", 10},

	// Food
	{"eat_sweets", "300 kcal of sweets", 60},
}

func init() {
	os.Setenv(TOKEN_NAME, TOKEN_VALUE)

	if gToken = os.Getenv(TOKEN_NAME); gToken == "" {
		panic(fmt.Errorf(`failed to load environment variable "%s"`, TOKEN_NAME))
	}

	var err error

	if gBot, err = tgbotapi.NewBotAPI(gToken); err != nil {
		log.Panic(err)
	}

	// gBot.Debug = true
}

func main() {
	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updConfig := tgbotapi.NewUpdate(0)
	updConfig.Timeout = UPDATE_TIMEOUT

	for update := range gBot.GetUpdatesChan(updConfig) {
		if isCallbackQuery(&update) {
			updateProcessing(&update)
		} else if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			gChatId = update.Message.Chat.ID
			askToPrintIntro()
		}
	}
}

// --- some helpful

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func isCallbackQuery(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data != ""
}

func callbackQueryFromIsMissing(update *tgbotapi.Update) bool {
	return update.CallbackQuery == nil || update.CallbackQuery.From == nil
}

func getUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if callbackQueryFromIsMissing(update) {
		return
	}

	userId := update.CallbackQuery.From.ID

	for _, userInChat := range gUsersInChat {
		if userId == userInChat.id {
			return userInChat, true
		}
	}

	return
}

func storeUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if callbackQueryFromIsMissing(update) {
		return
	}

	from := update.CallbackQuery.From
	user = &User{id: from.ID, name: strings.TrimSpace(from.FirstName + " " + from.LastName), gems: 0}
	gUsersInChat = append(gUsersInChat, user)

	return user, true
}

func showActivities(activities Activities, message string, isUseful bool) {
	activitiesButtonsRows := make([]([]tgbotapi.InlineKeyboardButton), 0, len(activities)+1)
	for _, activity := range activities {
		activityDescription := ""
		if isUseful {
			activityDescription = fmt.Sprintf("+ %d %s: %s", activity.gems, EMOJI_GEM, activity.name)
		} else {
			activityDescription = fmt.Sprintf("- %d %s: %s", activity.gems, EMOJI_GEM, activity.name)
		}
		activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(activityDescription, activity.code))
	}
	activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(BUTTON_TEXT_PRINT_MENU, BUTTON_CODE_PRINT_MENU))

	msg := tgbotapi.NewMessage(gChatId, message)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(activitiesButtonsRows...)
	gBot.Send(msg)
}

func findActivity(activities Activities, choiceCode string) (activity *Activity, found bool) {
	for _, activity := range activities {
		if choiceCode == activity.code {
			return activity, true
		}
	}
	return
}

// --- some useful

func sendMessageWithDelay(delayInSec uint8, message string) {
	sendStringMessage(message)
	delay(delayInSec)
}

func sendStringMessage(message string) {
	gBot.Send(tgbotapi.NewMessage(gChatId, message))
}

func delay(seconds uint8) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func getKeyboardRow(buttonText, buttonCode string) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode))
}

// some process actions

func updateProcessing(update *tgbotapi.Update) {
	user, found := getUserFromUpdate(update)
	if !found {
		if user, found = storeUserFromUpdate(update); !found {
			sendStringMessage("Unable to identify the user")
			return
		}
	}

	choiceCode := update.CallbackQuery.Data
	log.Printf("[%T] %s", time.Now(), choiceCode)

	switch choiceCode {
	case BUTTON_CODE_BALANCE:
		showBalance(user)
	case BUTTON_CODE_USEFUL_ACTIVITIES:
		showUsefulActivities()
	case BUTTON_CODE_REWARDS:
		showRewards()
	case BUTTON_CODE_PRINT_INTRO:
		printIntro(update)
		showMenu()
	case BUTTON_CODE_SKIP_INTRO:
		showMenu()
	case BUTTON_CODE_PRINT_MENU:
		showMenu()
	default:
		if usefulActivity, found := findActivity(gUsefulActivities, choiceCode); found {
			processUsefulActivity(usefulActivity, user)

			delay(2)
			showUsefulActivities()
			return
		}

		if reward, found := findActivity(gRewards, choiceCode); found {
			processReward(reward, user)

			delay(2)
			showRewards()
			return
		}

		log.Printf(`[%T] !!!!!!!!! ERROR: Unknown code "%s"`, time.Now(), choiceCode)
		msg := fmt.Sprintf("%s, I'm sorry, I don't recognize code '%s' %s Please report this error to my creator.", user.name, choiceCode, EMOJI_SAD)
		sendStringMessage(msg)
	}
}

func processUsefulActivity(activity *Activity, user *User) {
	errorMsg := ""
	if activity.gems == 0 {
		errorMsg = fmt.Sprintf(`Активность "%s" не имеет вознаграждения`, activity.name)
	} else if user.gems+activity.gems > MAX_USER_GEMS {
		errorMsg = fmt.Sprintf("Кажется, ты трудоголик (или врунишка). Достигнут максимум %d %s \nПотрать хоть не много, развлекаться тоже бывает полезно!", MAX_USER_GEMS, EMOJI_GEM)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, I'm sorry, but %s %s Your balance remains unchanged.", user.name, errorMsg, EMOJI_SAD)
	} else {
		user.gems += activity.gems
		resultMessage = fmt.Sprintf(`%s, the "%s" activity is completed! %d %s has been added to your account. Keep it up! %s%s Now you have %d %s`,
			user.name, activity.name, activity.gems, EMOJI_GEM, EMOJI_BICEPS, EMOJI_SUNGLASSES, user.gems, EMOJI_GEM)
	}
	sendStringMessage(resultMessage)
}

func processReward(activity *Activity, user *User) {
	errorMsg := ""
	if activity.gems == 0 {
		errorMsg = fmt.Sprintf(`Награда "%s" абсолютно бесплатна!`, activity.name)
	} else if user.gems < activity.gems {
		errorMsg = fmt.Sprintf(`you currently have %d %s. You cannot afford "%s" for %d %s`, user.gems, EMOJI_GEM, activity.name, activity.gems, EMOJI_GEM)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, денег нема %s %s Активность не доступна %s", user.name, errorMsg, EMOJI_SAD, EMOJI_DONT_KNOW)
	} else {
		user.gems -= activity.gems
		resultMessage = fmt.Sprintf(`%s, the reward "%s" has been paid for, get started! %d %s has been deducted from your account. Now you have %d %s`, user.name, activity.name, activity.gems, EMOJI_GEM, user.gems, EMOJI_GEM)
	}
	sendStringMessage(resultMessage)
}

// ---- some actions

func askToPrintIntro() {
	msg := tgbotapi.NewMessage(gChatId, "Привет! Я чат-бот, хочешь расскажу о своей миссии?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_PRINT_INTRO, BUTTON_CODE_PRINT_INTRO),
		getKeyboardRow(BUTTON_TEXT_SKIP_INTRO, BUTTON_CODE_SKIP_INTRO),
	)
	gBot.Send(msg)
}

func printIntro(update *tgbotapi.Update) {
	sendMessageWithDelay(4, "Моя задача помочь всему человечеству и тебе в частности прокачать себя "+EMOJI_SUNGLASSES)
	sendMessageWithDelay(7, "Я буду отслеживать твою активность.\nПолезные дела вознагражу, а за вредные возьму плату "+EMOJI_GEM)
}

func showMenu() {
	msg := tgbotapi.NewMessage(gChatId, "Выберите опцию:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_BALANCE, BUTTON_CODE_BALANCE),
		getKeyboardRow(BUTTON_TEXT_USEFUL_ACTIVITIES, BUTTON_CODE_USEFUL_ACTIVITIES),
		getKeyboardRow(BUTTON_TEXT_REWARDS, BUTTON_CODE_REWARDS),
	)
	gBot.Send(msg)
}

func showBalance(user *User) {
	msg := fmt.Sprintf("%s, у тебя 0 %s \nМожет быть сделаешь уже что-то полезное?", user.name, EMOJI_GEM)
	if gems := user.gems; gems > 0 {
		msg = fmt.Sprintf("%s, у тебя %d %s", user.name, gems, EMOJI_GEM)
	}
	sendStringMessage(msg)
	showMenu()
}

func showUsefulActivities() {
	showActivities(gUsefulActivities, "Track a useful activity or return to the main menu:", true)
}

func showRewards() {
	showActivities(gRewards, "Purchase a reward or return to the main menu:", false)
}
