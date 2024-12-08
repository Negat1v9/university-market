package static

import (
	"fmt"

	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
)

var (
	ErrorPreCheckOutAmount = "Что-то пошло не так"
	NoTgCmdFinded          = "<b>Хмм...</b> что то пошло не так"
	ErrBot                 = "Кажется что-то пошло не так. Попробуйте позже"
	FilesExpected          = "Пришлите <b>файл</b>\nЕсли у вас есть <u>только фото</u>, то вы можете прислать его как файл.\n<b>Максимум 5 файлов для одной задачи</b>"
	WaitingFiles           = "Пришлите файл (<u>только одно в сообщении</u>) для работы в чат либо опубликуйте работу сейчас.\n\nФайлы будут видны <u>всем работникам</u>, которые откликнулись на вашу работу.\n<b>Максимум 5 файлов для одной задачи</b>"
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

func OnRespondFromWorker() string {
	return "🎉 <b>Новый отклик!</b>"
}

func AddInformationTask(meta *taskmodel.TaskMeta) string {
	return fmt.Sprintf("\n\n<b>ℹ️ О работе:</b>\n⚪️ <u>Институт</u>: %s\n⚪️ <u>Задание</u>: %s\n⚪️ <u>Предмет</u>: %s\n", meta.TaskType, meta.TaskType, meta.Subject)
}
func SuccessAttachFiles() string {
	return fmt.Sprintf("<b>Файл успешно добавлен</b>\n%s", WaitingFiles)
}

func MsgShareContact(username string) string {
	return fmt.Sprintf("Вы выбрали <code>%s</code> работника для своей работы, можете связаться с ним для уточнения деталей! @%s", username, username)
}
