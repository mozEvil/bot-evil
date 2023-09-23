package main

const (
	TOKEN_NAME     = "TOKEN_BOT_EVIL"
	TOKEN_VALUE    = TOKEN_SECRET
	UPDATE_TIMEOUT = 60

	MAX_USER_GEMS uint16 = 500

	EMOJI_GEM          = "💎"          // "\U0001F48E"
	EMOJI_SMILE        = "\U0001F642" // 🙂
	EMOJI_SUNGLASSES   = "\U0001F60E" // 😎
	EMOJI_WOW          = "\U0001F604" // 😄
	EMOJI_DONT_KNOW    = "\U0001F937" // 🤷
	EMOJI_SAD          = "\U0001F63F" // 😿
	EMOJI_BICEPS       = "\U0001F4AA" // 💪
	EMOJI_BUTTON_START = "\U000025B6" // ▶
	EMOJI_BUTTON_END   = "\U000025C0" // ◀

	BUTTON_TEXT_PRINT_INTRO       = EMOJI_BUTTON_START + "Да, мне очень интересно!" + EMOJI_BUTTON_END
	BUTTON_TEXT_SKIP_INTRO        = EMOJI_BUTTON_START + "Нет, я уже все знаю." + EMOJI_BUTTON_END
	BUTTON_TEXT_BALANCE           = EMOJI_BUTTON_START + "Баланс" + EMOJI_BUTTON_END
	BUTTON_TEXT_USEFUL_ACTIVITIES = EMOJI_BUTTON_START + "Полезные дела" + EMOJI_BUTTON_END
	BUTTON_TEXT_REWARDS           = EMOJI_BUTTON_START + "Награды" + EMOJI_BUTTON_END
	BUTTON_TEXT_PRINT_MENU        = EMOJI_BUTTON_START + "МЕНЮ" + EMOJI_BUTTON_END

	BUTTON_CODE_PRINT_INTRO       = "print_intro"
	BUTTON_CODE_SKIP_INTRO        = "skip_intro"
	BUTTON_CODE_BALANCE           = "show_balance"
	BUTTON_CODE_USEFUL_ACTIVITIES = "show_useful_activities"
	BUTTON_CODE_REWARDS           = "show_rewards"
	BUTTON_CODE_PRINT_MENU        = "print_menu"
)
