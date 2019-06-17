package bank

type SortOrderType int8

const (
	SortAsc SortOrderType = iota
	SortDesc
)

func (t SortOrderType) Valid() bool {
	return t == SortAsc || t == SortDesc
}

type Filters struct {
	TransactionTypes []TransactionType
	SortOrder        SortOrderType
	Limit            int
	Offset           int
}

type Filter func(*Filters)

func FilterByTransactionType(transactionTypes ...TransactionType) Filter {
	return func(filters *Filters) {
		filters.TransactionTypes = transactionTypes
	}
}

func SortOrder(sortOrder SortOrderType) Filter {
	return func(filters *Filters) {
		if !sortOrder.Valid() {
			return
		}
		filters.SortOrder = sortOrder
	}
}

func Limit(limit int) Filter {
	return func(filters *Filters) {
		if limit <= 0 {
			return
		}
		filters.Limit = limit
	}
}

func Offset(offset int) Filter {
	return func(filters *Filters) {
		if offset <= 0 {
			return
		}
		filters.Offset = offset
	}
}

func getFilters(filterSetters ...Filter) Filters {
	filters := Filters{
		TransactionTypes: make([]TransactionType, 0),
		SortOrder:        SortAsc,
		Limit:            -1,
		Offset:           0,
	}
	for _, filter := range filterSetters {
		filter(&filters)
	}
	return filters
}
