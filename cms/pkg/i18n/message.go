package i18n

import (
	"fmt"
	"strconv"
)

type Message struct {
	ID      string
	Content string
}

func parseMessageFileBytes(buf []byte, unmarshalFunc UnmarshalFunc) ([]*Message, error) {
	var data map[string]interface{}
	err := unmarshalFunc(buf, &data)
	if err != nil {
		return nil, err
	}

	// 递归解析data
	return getMessages(data)
}

func getMessages(data map[string]interface{}) (messageList []*Message, err error) {
	for key, value := range data {
		if valueStr, ok := value.(string); ok {
			message := &Message{ID: key, Content: valueStr}
			messageList = append(messageList, message)
		} else {
			childMessageList, err := addChildMessages(value)
			if err != nil {
				return nil, err
			}

			for ind := range childMessageList {
				if childMessageList[ind].ID != "" {
					childMessageList[ind].ID = key + NestedSeparator + childMessageList[ind].ID
				} else {
					childMessageList[ind].ID = key
				}
			}

			messageList = append(messageList, childMessageList...)
		}

	}

	return
}

func addChildMessages(data interface{}) ([]*Message, error) {
	switch data := data.(type) {
	case string:
		return []*Message{{ID: "", Content: data}}, nil
	case map[string]interface{}:
		return getMessages(data)
	case []map[string]interface{}:
		var childMessagesList []*Message
		for index, item := range data {
			messageList, err := getMessages(item)
			if err != nil {
				return nil, err
			}

			for ind := range messageList {
				messageList[ind].ID = strconv.Itoa(index) + NestedSeparator + messageList[ind].ID
			}

			childMessagesList = append(childMessagesList, messageList...)
		}
		return childMessagesList, nil
	case []interface{}:
		var childMessagesList []*Message
		for index, item := range data {
			messageList, err := addChildMessages(item)
			if err != nil {
				return nil, err
			}

			for ind := range messageList {
				if messageList[ind].ID == "" {
					messageList[ind].ID = strconv.Itoa(index)
				} else {
					messageList[ind].ID = strconv.Itoa(index) + NestedSeparator + messageList[ind].ID
				}
			}

			childMessagesList = append(childMessagesList, messageList...)
		}
		return childMessagesList, nil
	default:
		return nil, fmt.Errorf("%w. %T: %+v", errInvalidTranslationValue, data, data)
	}
}
