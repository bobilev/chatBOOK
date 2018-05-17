package util

import "strconv"

func RKL(num int) string{// ReversKeyboardLayout
	arrayAlfa := []string{"к","а","б","в","г","д","е","ж","з","п"}

	return arrayAlfa[num]
}
func DeRKL(letter string) (string,bool){
	mapAlfa := map[string]int{"к":0,"а":1,"б":2,"в":3,"г":4,"д":5,"е":6,"ж":7,"з":8,"п":9}
	res,ok := mapAlfa[letter]
	return strconv.Itoa(res),ok
}