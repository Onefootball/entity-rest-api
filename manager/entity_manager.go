package manager

import (
	"fmt"
	"strings"
	"time"
	"errors"
	"strconv"
	"database/sql"
)

type EntityDbManager struct {
	Db *sql.DB
}

func NewEntityDbManager(db *sql.DB) *EntityDbManager {
	return &EntityDbManager{
		db,
	}
}

func (em *EntityDbManager) GetEntities(entity string, filterParams map[string]string, limit string, offset string, orderBy string, orderDir string) ([]map[string]interface{}, int, error) {

	var whereClause string

	if (len(filterParams) > 0) {
		whereClause = " WHERE "
		paramCount := 0

		r := strings.NewReplacer("*", "%")

		for filterParamKey, filterParamVal := range filterParams {

			filterParamVal = r.Replace(filterParamVal)

			if paramCount == 0 {
				whereClause = fmt.Sprintf("%s `%s` LIKE '%s'", whereClause, filterParamKey, filterParamVal)
			} else {
				whereClause = fmt.Sprintf("%s AND `%s` LIKE '%s'", whereClause, filterParamKey, filterParamVal)
			}
			paramCount++
		}
	} else {
		whereClause = ""
	}

	query := fmt.Sprintf(
		"SELECT * FROM `%s`%s ORDER BY %s %s LIMIT %s, %s",
		entity,
		whereClause,
		orderBy,
		orderDir,
		offset,
		limit,
	)

	allResults, err := em.retrieveAllResultsByQuery(query)

	if err != nil {
		return make([]map[string]interface{}, 0), 0, err
	}

	var countResult string

	countQuery := fmt.Sprintf(
		"SELECT count(id) FROM `%s`%s",
		entity,
		whereClause,
	)

	countErr := em.Db.QueryRow(countQuery).Scan(&countResult)

	if countErr != nil {
		return make([]map[string]interface{}, 0), 0, countErr
	}

	count, _ := strconv.Atoi(countResult)

	return allResults, count, nil
}

func (em *EntityDbManager) GetEntity(entity string, id string) (map[string]interface{}, error) {

	result, err := em.retrieveSingleResultById(entity, id)

	if err != nil {
		return make(map[string]interface{}), err
	}

	return result, nil
}

func (em *EntityDbManager) PostEntity(entity string, postData map[string]string) (int64, error) {

	columnsQuery := fmt.Sprintf(
		"SHOW COLUMNS FROM `%s`",
		entity,
	)

	columnsResult, err := em.retrieveAllResultsByQuery(columnsQuery)

	if err != nil {
		return 0, err
	}

	newEntity := make(map[string]string)

	for _, columnsRow := range columnsResult {
		column := columnsRow["Field"].(string)

		if column == "id" {
			continue
		}

		_, ok := postData[column]

		if ok {
			newEntity[column] = postData[column]
		} else {
			newEntity[column] = ""
		}
	}

	var insertColumnsString string
	var insertValuesString string

	propertyCount := 0

	for entityKey, entityVal := range newEntity {

		if propertyCount == 0 {
			insertColumnsString = fmt.Sprintf("`%s`", entityKey)
			insertValuesString = fmt.Sprintf("'%s'", entityVal)
		} else {
			insertColumnsString = fmt.Sprintf("%s, `%s`", insertColumnsString, entityKey)
			insertValuesString = fmt.Sprintf("%s, '%s'", insertValuesString, entityVal)
		}
		propertyCount++
	}

	insertQuery := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES(%s)",
		entity,
		insertColumnsString,
		insertValuesString,
	)

	res, err := em.Db.Exec(insertQuery)

	if err != nil {
		return 0, err
	}

	newId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (api *EntityDbManager) UpdateEntity(entity string, id string, updateData map[string]string) (int64, map[string]interface{}, error) {

	entityToUpdate, err := api.retrieveSingleResultById(entity, id)

	if err != nil {
		//rest.Error(w, err.Error(), http.StatusInternalServerError)
		return 0, make(map[string]interface{}), err
	}

	for updKey, _ := range entityToUpdate {
		_, ok := updateData[updKey]

		if ok {
			entityToUpdate[updKey] = updateData[updKey]
		}
	}

	var updateString string
	propertyCount := 0

	for entityKey, entityVal := range entityToUpdate {

		if propertyCount == 0 {
			updateString = fmt.Sprintf("`%s` = '%s'", entityKey, entityVal)
		} else {
			updateString = fmt.Sprintf("%s, `%s` = '%s'", updateString, entityKey, entityVal)
		}
		propertyCount++
	}

	updQuery := fmt.Sprintf(
		"UPDATE `%s` SET %s WHERE id = %s",
		entity,
		updateString,
		id,
	)

	res, err := api.Db.Exec(updQuery)

	if err != nil {
		return 0, make(map[string]interface{}), err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, make(map[string]interface{}), err
	}

	return rowsAffected, entityToUpdate, nil
}

func (api *EntityDbManager) DeleteEntity(entity string, id string) (int64, error) {

	query := fmt.Sprintf(
		"DELETE FROM `%s` WHERE id = %s",
		entity,
		id,
	)

	res, err := api.Db.Exec(query)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (api *EntityDbManager) retrieveAllResultsByQuery(query string) ([]map[string]interface{}, error) {

	allResults := make([]map[string]interface{}, 0)

	rows, err := api.Db.Query(query)
	defer rows.Close()

	fmt.Printf("%#v\n", query)

	if err != nil {
		return allResults, err
	}

	cols, err := rows.Columns()

	if err != nil {
		return allResults, err
	}

	rawResult := make([]interface{}, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice

	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make(map[string]interface{})

		err = rows.Scan(dest...)

		if err != nil {
			return allResults, err
		}

		for i, raw := range rawResult {
			result[cols[i]] = api.convertDbValue(raw)
		}

		allResults = append(allResults, result)
	}

	return allResults, nil
}

func (api *EntityDbManager) retrieveSingleResultById(entity string, id string) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	query := fmt.Sprintf(
		"SELECT * FROM `%s` WHERE id = %s",
		entity,
		id,
	)

	rows, err := api.Db.Query(query)
	defer rows.Close()

	if err != nil {
		return result, err
	}

	cols, err := rows.Columns()

	if err != nil {
		return result, err
	}

	rawResult := make([]interface{}, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice

	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	resultCount := 0

	for rows.Next() {

		resultCount++
		err = rows.Scan(dest...)

		if err != nil {
			return result, err
		}

		for i, raw := range rawResult {
			result[cols[i]] = api.convertDbValue(raw)
		}

		if resultCount > 1 {
			continue
		}
	}

	if resultCount != 1 {
		return result, errors.New("Id query returned inappropriate result")
	}

	return result, nil
}

func (api *EntityDbManager) convertDbValue(dbValue interface{}) interface{} {

	switch t := dbValue.(type) {
		default:
		fmt.Printf("[EntityDbManager] Unexpected type %T: %#v\n", t, dbValue)
		return ""
		case bool:
		return dbValue.(bool)
		case int:
		return dbValue.(int)
		case int64:
		return dbValue.(int64)
		case []byte:
		return string(dbValue.([]byte))
		case string:
		return dbValue.(string)
		case time.Time:
		return dbValue.(time.Time).String()
		case nil:
		return nil
	}
}

