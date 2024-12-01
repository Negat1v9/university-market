package static

import (
	"fmt"

	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
)

var (
	ErrorPreCheckOutAmount = "Что-то пошло не так"
	NoTgCmdFinded          = "<b>Хмм...</b> что то пошло не так"
	ErrBot                 = "Кажется что-то пошло не так. Попробуйте позже"
	FilesExpected          = "Пришлите файл или фото"
	WaitingFiles           = "Пришлите фото (<u>только одно в сообщении</u>) или файл для работыв чат либо опубликуйте работу сейчас.\n\nФайлы будут видны <u>всем работникам</u>, которые откликнулись на вашу работу."
	SuccessPublichTask     = "Работа успешно опубликована!\nОжидайте отклики от работников."
	HelpCmd                = "Как это работает:"
	WorkerSelected         = "🥳 Мы рады сообщить вам, что пользователь выбрал вас для выполнения работы.\nПожалуйста, свяжитесь с ним для обсуждения деталей."
	ErrNoUserName          = "Чтобы поделился контактом с пользователем у вас должен быть <u>username</u>.\nЕго можно создать в настройках профиля <b>телеграм</b>."
)

func MsgStartCommand(nickname string) string {
	return fmt.Sprintf("Добро пожаловать <u>%s</u>", nickname)
}
func MsgErrOnSuccessPayment(tgPaymentID string) string {
	return fmt.Sprintf("Кажется с нашей стороны произошла ошибка\n<b>Обратитесь пожалуйста в поддержку</b>\nПерешлите это сообщение в чат поддержки, мы обязательно разберемся с вашей проблемой.\n<code>%s</code>", tgPaymentID)
}
func BalancePayment(amount int) string {
	return fmt.Sprintf("Баланс вашего аккаунта успешно пополнен на %d звезд", amount)
}

func OnRespondFromWorker(worker *usermodel.WorkerInfo) string {
	return fmt.Sprintf("<b>Новый отклик!</b>\nО работнике:\n<code>%s</code>", worker.Description)
}

func SuccessAttachFiles() string {
	return fmt.Sprintf("<b>Файл успешно добавлен</b>\n%s", WaitingFiles)
}

func MsgShareContact(username string) string {
	return fmt.Sprintf("Вы выбрали <code>%s</code> работка для своей работы, можете связаться с ним для уточнения деталей! @%s", username, username)
}
