package createtypemispobjects

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
)

type OptionsAddNewObject struct {
	Host        string
	AuthKey     string
	UserAuthKey string
}

func AddNewObject(ctx context.Context, data mispapi.InputSettings, opts OptionsAddNewObject) {
	fmt.Println("func 'specialObject.SetFunc', START... |||")

	rmisp, err := mispapi.NewMispRequest(
		mispapi.WithHost(opts.Host),
		mispapi.WithUserAuthKey(opts.UserAuthKey),
		mispapi.WithMasterAuthKey(opts.AuthKey))
	if err != nil {
		log.Println(err)

		return
	}

	log.Println("func 'specialObject.SetFunc', send event --->")

	//отправляет в API MISP событие в виде типа Event и возвращает результат который содержит
	//id события в MISP, у MISP свой уникальный id для событий
	//только с использованием этого id в MISP добавляются все остальные объекты
	_, resBodyByte, err := rmisp.SendEvent_ForTest(ctx, data.Data.GetEvent())
	if err != nil {
		log.Println("func 'specialObject.SetFunc', EVENT ERROR:", err)

		return
	}

	fmt.Println("func 'specialObject.SetFunc', get event response <---")

	//Все ошибки которые могут возникнуть при дальнейшем взаимодействии с MISP
	//будут попрежнему логироватся.
	//Однако, статус выполнения для функции будет ставится в TRUE, что бы не досить
	//MISP, так как все последующие попытки будут начинатся с добавления 'event', а
	//добавить 'event' с таким id нельзя.
	//Необходимо удалить предыдущий.

	resMisp := mispapi.MispResponse{}
	if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
		log.Println(err)

		return
	}

	log.Println("func 'specialObject.SetFunc', MISP response:", resMisp)

	//получаем уникальный id MISP
	var eventId string
	for key, value := range resMisp.Event {
		if key == "id" {
			if str, ok := value.(string); ok {
				eventId = str

				break
			}
		}
	}

	log.Println("func 'specialObject.SetFunc', MISP eventId:", eventId)

	if eventId == "" {
		log.Println("the formation of events of the 'Attributes' type was not performed because the EventID is empty")

		return
	}

	// добавляем event_reports
	if err := rmisp.SendEventReports_ForTest(ctx, eventId, data.Data.GetReports()); err != nil {
		log.Println(err)

		return
	}

	//
	// в тесте это пока придётся исключить!!!!
	//
	//отправляем запрос для добавления в БД Redis, id кейса и нового события
	//m.SendDataOutput(OutputSetting{
	//	Command: "set new event id",
	//	CaseId:  fmt.Sprint(data.CaseId),
	//	EventId: eventId,
	//})
	//

	//добавляем атрибуты
	_, _, warning, err := rmisp.SendAttribytes_ForTest(ctx, eventId, data.Data.GetAttributes())
	if err != nil {
		// тут ошибка может быть при добавлении только одного из многих объектов
		// соответственно тормозить весь процесс только из-за того что была ошибка
		// при добавлении одного или нескольких объектов не стоит
		// если же не был добавлен ни один из объектов, то это возможно глабальная
		// ошибка доступа, следовательно, при добавлении следующих объектов она также
		// может вылезти, тогда там и будет выполнен останов всей цепочки
		log.Println(err)
	}
	if warning != "" {
		log.Println("warning", warning)
	}

	// добавляем объекты
	if _, _, err = rmisp.SendObjects_ForTest(ctx, eventId, data.Data.GetObjects()); err != nil {
		//тут такая же ситуация что и с ошибками при выполнении метода rmisp.sendAttribytes
		log.Println(err)
	}

	// берем небольшой таймаут, нужен для того что бы MISP успел обработать и добавить в БД
	// всё ранее ему переданное, если обработка переданных объектов не была завершена
	// возможны накладки или сбои при добавлении данных
	// это недостаток MISP, с этим я ничего не могу поделать
	time.Sleep(4 * time.Second)

	// добавляем event_tags
	if err := rmisp.SendEventTags_ForTest(ctx, eventId, data.Data.GetObjectTags()); err != nil {
		log.Println(err)
	}

	//публикуем добавленное событие
	//masterKey нужен для публикации события так как пользователь
	//должен иметь более расшириные права чем могут иметь некоторые
	//обычные пользователи
	resMsg, err := rmisp.SendRequestPublishEvent_ForTest(ctx, eventId)
	if err != nil {
		log.Println(err)
	}
	if resMsg != "" {
		log.Panicln("warning", resMsg)
	}

	// отправляем в ядро информацию по event Id
	//m.SendDataOutput(OutputSetting{
	//	Command: "send event id",
	//	EventId: eventId,
	//	CaseId:  fmt.Sprint(data.CaseId),
	//	RootId:  data.RootId,
	//	TaskId:  data.TaskId,
	//})

	fmt.Println("func 'specialObject.SetFunc', STOP |||")
}
