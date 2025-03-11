package sqlite3api

import (
	"context"
	"fmt"
)

// Ping проверка соединения с БД
func (module *ApiSqlite3Module) Ping(ctx context.Context) error {
	return module.db.PingContext(ctx)
}

// ConnectionClose закрывает соединение с БД
func (module *ApiSqlite3Module) ConnectionClose() {
	module.db.Close()
}

// GetChRequest канал для запросов к модулю
func (module *ApiSqlite3Module) GetChRequest() <-chan Request {
	return module.chRequest
}

// SendDataToModule передача данных в модуль
func (module *ApiSqlite3Module) SendDataToModule(req Request) {
	module.chRequest <- req
}

// SearchCaseId поиск информации по caseId
func (module *ApiSqlite3Module) SearchCaseId(ctx context.Context, caseId int) (result int, err error) {
	rows, err := module.db.QueryContext(ctx, "SELECT eventId FROM placeholder_misp WHERE caseId=?", caseId)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		if err := rows.Scan(&result); err != nil {
			break
		}
	}

	return result, nil
}

// UpdateCaseId обнавляет содержимое таблицы по caseId, если caseId не найдено, добавляет
func (module *ApiSqlite3Module) UpdateCaseId(ctx context.Context, caseId, eventId int) error {
	id, err := module.SearchCaseId(ctx, caseId)
	if err != nil {
		fmt.Println("func 'ApiSqlite3Module.UpdateCaseId' Search ERROR:", err)

		return err
	}

	fmt.Printf("func 'ApiSqlite3Module.UpdateCaseId' found eventId:%d\n", id)

	if id == 0 {
		fmt.Println("func 'ApiSqlite3Module.UpdateCaseId' INSERT")

		if _, err := module.db.ExecContext(ctx, "INSERT INTO placeholder_misp (caseId, eventId) VALUES (?,?)", caseId, eventId); err != nil {
			fmt.Println("func 'ApiSqlite3Module.UpdateCaseId' Insert ERROR:", err)

			return err
		}
	} else {
		fmt.Println("func 'ApiSqlite3Module.UpdateCaseId' UPDATE")

		if _, err := module.db.ExecContext(ctx, "UPDATE placeholder_misp SET eventId=? WHERE caseId=?", eventId, caseId); err != nil {
			fmt.Println("func 'ApiSqlite3Module.UpdateCaseId' Update ERROR:", err)

			return err
		}
	}

	return nil
}

// DeleteCaseId удаляет информацию по кейсу
func (module *ApiSqlite3Module) DeleteCaseId(ctx context.Context, caseId int) error {
	if _, err := module.db.ExecContext(ctx, "DELETE FROM placeholder_misp WHERE caseId=?", caseId); err != nil {
		return err
	}

	return nil
}
