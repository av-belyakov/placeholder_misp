package memorytemporarystorage_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"

	"github.com/av-belyakov/placeholder_misp/memorytemporarystorage"
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

	Context("Тест 1. Тестируем добавление, редактирование и удаление информации в HiveFormatMessage", func() {
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

			Expect(at).ShouldNot(Equal(BeFalse()))
			Expect(ok).ShouldNot(Equal(BeFalse()))
		})

		It("Должен быть успешно удален элемент", func() {
			testDelUuid := uuid.NewString()
			cst.SetProcessedDataHiveFormatMessage(uuid.NewString(), map[string]interface{}{"test message delete_23": 23})
			cst.SetProcessedDataHiveFormatMessage(testDelUuid, map[string]interface{}{"test message delete_4": 4})

			Expect(len(cst.HiveFormatMessage.Storages)).Should(Equal(5))

			cst.HiveFormatMessage.Delete(testDelUuid)
			Expect(len(cst.HiveFormatMessage.Storages)).Should(Equal(4))
		})
	})

	Context("Тест 2. Тестируем хранилище temporaryInputCase", func() {
		It("Должно быть успешно добавлено несколько значений", func() {
			cst.SetTemporaryCase(64375, memorytemporarystorage.SettingsInputCase{EventId: "4532"})

			Expect(len(cst.GetListTemporaryCases())).Should(Equal(1))

			cst.SetTemporaryCase(788435, memorytemporarystorage.SettingsInputCase{EventId: "4533"})
			cst.SetTemporaryCase(134355, memorytemporarystorage.SettingsInputCase{EventId: "4534"})

			Expect(len(cst.GetListTemporaryCases())).Should(Equal(3))

			tc, ok := cst.GetTemporaryCase(788435)

			Expect(ok).Should(BeTrue())
			Expect(tc.EventId).Should(Equal("4533"))
		})
	})
})
