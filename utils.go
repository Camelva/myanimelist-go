package myanimelist

// fieldsList is small util function to avoid importing additional libraries
//func fieldsList(fields ...string) map[string]struct{} {
//	var fieldsMap = make(map[string]struct{}, len(fields))
//	if len(fields) < 1 {
//		return fieldsMap
//	}
//	for _, v := range fields {
//		fieldsMap[v] = struct{}{}
//	}
//	return fieldsMap
//}

func makeList(objects []string) map[string]struct{} {
	var list = make(map[string]struct{}, len(objects))
	if len(objects) < 1 {
		return list
	}
	for _, v := range objects {
		list[v] = struct{}{}
	}
	return list
}

func makeListInt(objects []int) map[int]struct{} {
	var list = make(map[int]struct{}, len(objects))
	if len(objects) < 1 {
		return list
	}
	for _, v := range objects {
		list[v] = struct{}{}
	}
	return list
}

// fixSorting is small internal helper function.
// Main purpose is to avoid creating duplicate variables
// such as SortAnimeListByStartTime && SortMangaListByStartTime
func fixSorting(sort string, kind string) string {
	var prefix string
	if kind == "anime" {
		prefix = "anime_"
	} else if kind == "manga" {
		prefix = "manga_"
	} else {
		// undefined kind
		return string(sort)
	}
	if sort == SortListByTitle || sort == SortListByStartDate || sort == SortListByID {
		return prefix + string(sort)
	}
	return string(sort)
}
