package models

// модель запроса на получение короткой ссылки
type ShortenAPIRequest struct {
	Url *string `json:"url"` //обязательное
}

// модель ответа с короткой ссылкой
type ShortenAPIResponse struct {
	Result string `json:"result"`
}
