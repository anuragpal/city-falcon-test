package models

type ListParams struct {
    Query 			string 					`json:"query,omitempty"`
    Page 			int64 					`json:"page"`
    RecordPerPage	int64 					`json:"rpp"`
    OrderBy 		string 					`json:"order_by"`
    SortBy 			string 					`json:"sort_by"`
    Record 			string 					`json:"record,omitempty"`
}