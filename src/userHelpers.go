package main

func chunkArray(arr []Item, chunkSize int) [][]Item {
	var divided [][]Item
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		divided = append(divided, arr[i:end])
	}
	return divided
}
