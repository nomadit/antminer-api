package util

type Page struct {
    List       *[]interface{} `json:"list"`
    PageNum    int            `json:"pageNum"`
    PerPage    int            `json:"perPage"`
    TotalCount int            `json:"totalCount"`
}
