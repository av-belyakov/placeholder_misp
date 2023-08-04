package memorytemporarystorage_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"

	"placeholder_misp/memorytemporarystorage"
)

var _ = Describe("StorageTemporary", Ordered, func() {
	var (
		uuidMsg1, uuidMsg2, uuidMsg3 string
		cst                          *memorytemporarystorage.CommonStorageTemporary
	)

	BeforeAll(func() {
		uuidMsg1, uuidMsg2, uuidMsg3 = uuid.NewString(), uuid.NewString(), uuid.NewString()

		cst = memorytemporarystorage.NewTemporaryStorage()
	})

	Context("Тест 1. Тестируем добавление и редактирование информации в HiveFormatMessage", func() {
		It("Должны быть успешно добавлены RAW данные", func() {
			b1 := []byte("test message 1")
			cst.SetRawDataHiveFormatMessage(uuidMsg1, b1)
			cst.SetRawDataHiveFormatMessage(uuidMsg2, []byte("test message 2"))

			Expect(len(cst.HiveFormatMessage.Storages)).Should(Equal(2))

			rd, ok := cst.GetRawDataHiveFormatMessage(uuidMsg1)

			Expect(rd).Should(Equal(b1))
			Expect(ok).ShouldNot(Equal(BeFalse()))
		})

		It("Должно быть успешно добавлено обрабатываемое сообщение", func() {
			pd1 := map[string]interface{}{"test message": 1}

			cst.SetProcessedDataHiveFormatMessage(uuidMsg1, pd1)
			cst.SetProcessedDataHiveFormatMessage(uuidMsg3, map[string]interface{}{"test message": 3})

			Expect(len(cst.HiveFormatMessage.Storages)).Should(Equal(3))

			pd, ok := cst.GetProcessedDataHiveFormatMessage(uuidMsg1)

			Expect(pd).Should(Equal(pd1))
			Expect(ok).ShouldNot(Equal(BeFalse()))
		})

		It("Должно быть успешно изменено состояние информирующее что можно пропустить сообщение", func() {
			cst.SetAllowedTransferTrueHiveFormatMessage(uuidMsg3)
			cst.SetRawDataHiveFormatMessage(uuidMsg3, []byte("dfd gfg hfghtg hh"))

			at, ok := cst.GetAllowedTransferHiveFormatMessage(uuidMsg3)
			pd, _ := cst.GetProcessedDataHiveFormatMessage(uuidMsg3)
			rd, _ := cst.GetRawDataHiveFormatMessage(uuidMsg3)

			fmt.Println("___=== uuid:", uuidMsg3, " pd: ", pd, " rd: ", rd, " at: ", at)

			Expect(at).ShouldNot(Equal(BeFalse()))
			Expect(ok).ShouldNot(Equal(BeFalse()))
		})
	})

	/*Context("Тест 2. Проверяем возможность автоматического удаления информации", func() {
			It("Должно быть успешно изменено состояние предназначенное для автоматического удаления временное информации", func(){
	cst.
			})

			It("При проверке наличия информации он должна быть удалена", func() {

			})
		})*/
})
